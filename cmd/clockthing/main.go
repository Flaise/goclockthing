package main

import (
	"io"
	"time"

	"github.com/Flaise/goclockthing/assets"
	"github.com/Flaise/playwav"
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

func nextTone(path string) error {
	file, err := assets.Content.Open(path)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	reader := file.(io.ReadSeeker)

	err = playwav.FromReader(reader, size)
	if err != nil {
		return err
	}

	time.Sleep(110 * time.Millisecond)

	return nil
}

func playTally() error {
	fours, ones := nowToTally()

	for i := 0; i < fours; i += 1 {
		err := nextTone("four.wav")
		if err != nil {
			return err
		}
	}
	for i := 0; i < ones; i += 1 {
		err := nextTone("one.wav")
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	scheduler := gocron.NewScheduler(time.Local)

	scheduler.Every(1).Hour().StartAt(time.Unix(0, 0)).Do(func() {
		err := playTally()
		if err != nil {
			panic(err)
		}
	})

	scheduler.StartBlocking()
}
