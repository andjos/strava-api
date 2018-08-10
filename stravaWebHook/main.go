package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type stravaE struct {
	StravaObjectType string `json:"object_type"`
	StravaObjectID   int    `json:"object_id"`
	StravaOwnerID    int    `json:"owner_id"`
}

func main() {
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

//Respond to Strava Post, take care of OwnerID, Objecttype and ObjectID
func postEvent(c *gin.Context) {
	d := &stravaE{}
	c.Bind(d)
	c.String(http.StatusOK, "")
	log.Println(d.StravaOwnerID, d.StravaObjectType, d.StravaObjectID)

}
