package main

import (
	"fmt"
	"os"

	strava "github.com/strava/go.strava"
)

func main() {

	stravaToken := os.Getenv("strava_token")
	if stravaToken == "" {
		fmt.Println("Strava token missing!!! :(")
		os.Exit(1)
	}

	client := strava.NewClient(stravaToken)
	service := strava.NewCurrentAthleteService(client)

	curUserActivites, err := service.ListActivities().PerPage(1).Do()
	if err != nil {
		os.Exit(1)
	}

	// (cadence * 2) * minuter
	totalStep := (int(curUserActivites[0].AverageCadence) * 2) * (curUserActivites[0].MovingTime / 60)
	fmt.Println(totalStep)

}
