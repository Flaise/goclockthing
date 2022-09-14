package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

// See https://alsa.opensrc.org/Aplay#Questions for explanation of what arguments are valid.
var deviceName string

func playTone(path string) error {
	file, err := assets.Content.Open(path)
	if err != nil {
		return err
	}

	cmd := exec.Command("aplay", "-q", "-D", deviceName)
	cmd.Stderr = os.Stderr

	toCmd, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = io.Copy(toCmd, file)
	if err != nil {
		return err
	}

	err = toCmd.Close()
	if err != nil {
		return err
	}

	// Sometimes the first tone is delayed so wait on the command
	// instead of a regular sleep instruction.
	err = cmd.Wait()
	if err != nil {
		return err
	}

	// each audio clip is 200ms
	// also seems to be roughly 40ms of lag after
	time.Sleep(40 * time.Millisecond)

	return nil
}

func playTally(fours int, ones int) error {
	for i := 0; i < fours; i += 1 {
		err := playTone("four.wav")
		if err != nil {
			return err
		}
	}

	time.Sleep(110 * time.Millisecond)

	for i := 0; i < ones; i += 1 {
		err := playTone("one.wav")
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
	err := playTone("half.wav")
	if err != nil {
		return err
	}
	return nil
}

func scheduleChimes(scheduler *gocron.Scheduler) {
	scheduler.Every(1).Hour().StartAt(time.Unix(0, 0)).Tag("chime").Do(func() {
		if time.Now().Minute() != 0 {
			// computer probably woke from sleeping
			rescheduleChimes(scheduler)
			return
		}

		err := playCurrentTally()
		if err != nil {
			panic(err)
		}
	})

	scheduler.Every(1).Hour().StartAt(time.Unix(60*30, 0)).Tag("chime").Do(func() {
		if time.Now().Minute() != 30 {
			// computer probably woke from sleeping
			rescheduleChimes(scheduler)
			return
		}

		err := playHalf()
		if err != nil {
			panic(err)
		}
	})
}

func rescheduleChimes(scheduler *gocron.Scheduler) {
	scheduler.RemoveByTag("chime")
	scheduleChimes(scheduler)
}

func main() {
	flag.StringVar(&deviceName, "device-name", "default",
		"The device name to play audio on. If omitted, use the default device.")
	doTest := flag.Bool("test", false, "True to play audio immediately to test speaker setup.")
	flag.Parse()

	if *doTest {
		err := playTally(2, 3)
		if err != nil {
			panic(err)
		}
		return
	}

	scheduler := gocron.NewScheduler(time.Local)

	// Need to continually reschedule to fix timers after computer sleeps.
	scheduler.Every(30).Seconds().Do(func() {
		rescheduleChimes(scheduler)
	})

	fmt.Println("Running clock thing...")
	scheduler.StartBlocking()
}
