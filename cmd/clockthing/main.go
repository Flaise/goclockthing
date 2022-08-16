package main

import (
	"io"
	"os/exec"
	"time"

	"github.com/Flaise/goclockthing/assets"
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

	cmd := exec.Command("aplay")
	stream, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = io.Copy(stream, file)
	if err != nil {
		return err
	}

	err = stream.Close()
	if err != nil {
		return err
	}

	time.Sleep(290 * time.Millisecond)

	return nil
}

func playTally(fours int, ones int) error {
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

func playCurrentTally() error {
	fours, ones := nowToTally()

	return playTally(fours, ones)
}

func playHalf() error {
	err := nextTone("half.wav")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	scheduler := gocron.NewScheduler(time.Local)

	scheduler.Every(1).Hour().StartAt(time.Unix(0, 0)).Do(func() {
		err := playCurrentTally()
		if err != nil {
			panic(err)
		}
	})

	scheduler.Every(1).Hour().StartAt(time.Unix(60*30, 0)).Do(func() {
		err := playHalf()
		if err != nil {
			panic(err)
		}
	})

	scheduler.StartBlocking()
}
