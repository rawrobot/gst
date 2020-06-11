package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bksworm/gst"
)

func main() {
	pipeline, err := gst.ParseLaunch(
		" matroskamux name=mux ! filesink location=test.mkv  " +
			"videotestsrc pattern=ball num-buffers=150 ! video/x-raw,format=I420,width=320,height=240,framerate=30/1 " +
			" ! queue ! x264enc ! mux.video_0  " +
			" appsrc name=texter format=time is-live=true do-timestamp=true ! text/x-raw,format=utf8  !  queue ! mux.subtitle_0 ")
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
	go appsrc.PushTime(2)

	lp := gst.MainLoopNew()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// Run Cleanup
		pipeline.SetState(gst.StateReady)
		pipeline.SetState(gst.StateNull)
		lp.Quit()
	}()

	pipeline.SetState(gst.StatePlaying)
	go lp.Run()
	time.Sleep(5 * time.Second)
	pipeline.SetState(gst.StateReady)
	pipeline.SetState(gst.StateNull)
	lp.Quit()
}

type TextPusher struct {
	appsrc *gst.Element
}

func (tp *TextPusher) PushTime(framerate int) (err error) {
	interval := time.Second / time.Duration(framerate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	var frames int
	for {
		select {
		case <-ticker.C:
		}
		currentTime := time.Now()
		str := currentTime.Format("15:04:05")

		err = tp.appsrc.PushBuffer([]byte(str)) //Async([]byte(str), framerate, frames)
		//log.Println(b)
		if err != nil {
			break
		}
		frames += 1
	}
	return err
}
