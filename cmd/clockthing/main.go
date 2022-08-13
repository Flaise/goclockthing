package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func hoursToTally(hours int) (fours int, ones int) {
	fours = hours / 4
	ones = hours - (fours * 4)
	return
}

func nowToTally() (fours int, ones int) {
	hours := time.Now().Hour() % 12
	if hours == 0 {
		hours = 12
	}
	return hoursToTally(hours)
}

func main() {
	scheduler := gocron.NewScheduler(time.Local)

	scheduler.Every(1).Hour().StartAt(time.Unix(0, 0)).Do(func() {
		fours, ones := nowToTally()

		for i := 0; i < fours; i += 1 {

		}
		for i := 0; i < ones; i += 1 {

		}
	})

	scheduler.StartBlocking()
}
