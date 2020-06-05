package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bksworm/gst"
)

const (
	SIZE      = 128
	frameRate = 25
	movieLen  = 3
)

func main() {
	pipeline, err := gst.ParseLaunch(
		// " videotestsrc pattern=ball ! video/x-raw,format=I420,width=320,height=240,framerate=25/1 " +
		// 	" ! textoverlay name=ov text=\"nothing\" valignment=top halignment=left  ! autovideosink")
		" textoverlay name=ov valignment=top halignment=left  ! autovideosink  " +
			" videotestsrc pattern=ball ! video/x-raw,format=I420,width=320,height=240,framerate=25/1 ! ov.video_sink " +
			" appsrc name=texter ! text/x-raw,format=utf8 ! ov.text_sink")

	if err != nil {
		log.Println("pipeline create error ", err.Error())
		return
	}

	log.Println(pipeline.Name())

	element := pipeline.GetByName("texter")

	if element == nil {
		log.Println("pipe find element error")
		return
	}

	appsrc := TextPusher{element}
	go appsrc.PushTime(25)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// Run Cleanup
		pipeline.SetState(gst.StateReady)
		pipeline.SetState(gst.StateNull)
		os.Exit(0)
	}()

	pipeline.SetState(gst.StatePlaying)

	select {}

}

type TextPusher struct {
	appsrc *gst.Element
}

func (tp *TextPusher) PushTime(framerate int) (err error) {
	interval := time.Second / time.Duration(framerate)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		}
		currentTime := time.Now()
		str := currentTime.Format("15:04:05")
		//b := make([]byte, len(str))
		//copy(b[:], str)
		err = tp.appsrc.PushBuffer([]byte(str))
		//log.Println(b)
		if err != nil {
			break
		}
	}
	return err
}
