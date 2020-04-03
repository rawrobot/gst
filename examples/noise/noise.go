package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/bksworm/gst"
)

const (
	SIZE      = 320 * 240 * 3
	frameRate = 25
	movieLen  = 3
)

func main() {
	vrc := NewVideoRec()
	err := vrc.Assamble()
	errorCheck(err)
	// Launch the application
	ml := gst.MainLoopNew()
	go ml.Run()
	vrc.Record()
	log.Println("Start!")
	err = makeNoise(vrc, time.Second/frameRate)
	errorCheck(err)
	vrc.Pause()
	ml.Quit()
}

func errorCheck(e error) {
	if e != nil {
		// panic for any errors.
		log.Panic(e)
	}
}
func makeNoise(vrec *VideoRec, interval time.Duration) (err error) {
	b := make([]byte, SIZE)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	nOfFrames := frameRate * movieLen
	for n := 0; n < nOfFrames; n++ {
		select {
		case <-ticker.C:
			//log.Printf("Sent at %v", t)
		}

		for i := 0; i < SIZE; i++ {
			b[i] = byte(rand.Intn(255))
		}

		err = vrec.PushBuffer(b)
		//log.Println("+")
		if err != nil {
			break
		}
	}
	return err
}
