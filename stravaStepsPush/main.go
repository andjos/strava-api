package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/strava/go.strava"
	pushbullet "github.com/xconstruct/go-pushbullet"

	"github.com/gin-gonic/gin"
)

type stravaE struct {
	StravaObjectType string `json:"object_type"`
	StravaObjectID   int    `json:"object_id"`
	StravaOwnerID    int    `json:"owner_id"`
	StravaAspectType string `json:"aspect_type"`
}

var stravaToken string
var pushBulletToken string

func main() {

	stravaToken = os.Getenv("stravaToken")
	if stravaToken == "" {
		log.Println("Strava token missing!!! :(")
		os.Exit(1)
	}

	pushBulletToken = os.Getenv("pushBulletToken")
	if pushBulletToken == "" {
		log.Println("PushBullet token missing!!! :(")
		os.Exit(1)
	}

	log.SetFlags(log.LstdFlags) // Turn on timestamps for logging...
	router := gin.Default()

	router.GET("/ping", ping)
	router.GET("/event", hubChallengeEvent)
	router.POST("/event", postEvent)
	router.Run()
}

// Simple Alive function
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})

}

// Response to hub.challenge
func hubChallengeEvent(c *gin.Context) {
	hubChallenge := c.Query("hub.challenge")
	c.JSON(http.StatusOK, gin.H{
		"hub.challenge": hubChallenge,
	})
	log.Println("Hub Challenge Token:", hubChallenge)
}

//Respond to Strava Post, take care of OwnerID, Objecttype and ObjectID and AspectType
func postEvent(c *gin.Context) {
	d := &stravaE{}
	c.Bind(d)
	c.String(http.StatusOK, "")
	log.Println("OwnerID:", d.StravaOwnerID, "ObjectType:", d.StravaObjectType, "AspectType:", d.StravaAspectType, "ObjectID:", d.StravaObjectID)

	if d.StravaAspectType != "delete" {

		sSteps := stravaSteps(d)
		log.Println("TotalSteps taken:", sSteps)
		pushSteps(sSteps)

	}

}

//Calculate total steps taken
func stravaSteps(d *stravaE) int {

	client := strava.NewClient(stravaToken)
	service := strava.NewActivitiesService(client)
	activity, err := service.Get(int64(d.StravaObjectID)).Do()
	if err != nil {
		log.Println("Problem to get activity, ID:", d.StravaObjectID)
		panic(err)
	}

	var totalStep int
	if activity.AverageCadence == 0 {
		log.Println("No cadence data found!!")
		totalStep = 0
	} else {
		totalStep = (int(activity.AverageCadence) * 2) * (activity.MovingTime / 60)
	}

	return totalStep
}

// Push steps to devices via Pushbullet
func pushSteps(s int) {
	var message string
	pb := pushbullet.New(pushBulletToken)
	devs, err := pb.Devices()
	if err != nil {
		log.Println("Problem with pushdevices.")
		panic(err)
	}

	if s != 0 {
		message = "Total steps taken: " + strconv.Itoa(s)
	} else {
		message = "No cadence data found, total steps not possible to calculate :("
	}

	for dev := range devs {

		if devs[dev].Active == true {
			err = pb.PushNote(devs[dev].Iden, "StravaSteps", message)
			if err != nil {
				log.Println("Problem to push.")
				panic(err)
			}

			log.Println("Msg Pushed to:", devs[dev].Nickname)
		}
	}

}
