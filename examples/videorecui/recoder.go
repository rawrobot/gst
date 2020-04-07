package main

import (
	"fmt"

	"github.com/bksworm/gst"
)

const frameSrc = "frameSrc"

type VideoRec struct {
	pipe      *gst.Pipeline
	frameSrc  *gst.Element
	recording bool
	frames    chan []byte
	shutter   chan int
}

func NewVideoRec() (vrec *VideoRec) {
	vrec = &VideoRec{}
	vrec.frames = make(chan []byte, 10) //this is queue for 10 jpeg images
	vrec.shutter = make(chan int, 1)
	return vrec
}

func (p *VideoRec) Assemble() (err error) {
	p.pipe, err = gst.ParseLaunch("videotestsrc  ! video/x-raw,width=640,height=480,format=I420 ! " +
		"vp9enc !  splitmuxsink muxer=matroskamux location=video%02d.mkv " +
		" max-size-time=10000000000 max-size-bytes=1000000")

	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}
	return err
}

func (p *VideoRec) Play() {
	p.pipe.SetState(gst.StatePlaying)
	p.recording = true
}

func (p *VideoRec) Pause() {
	p.pipe.SetState(gst.StatePaused)
	p.recording = false
}

func (p *VideoRec) Stop() {
	p.pipe.SetState(gst.StateReady)
	p.recording = false
}

func (p *VideoRec) Close() {
	p.pipe.SetState(gst.StateNull)
	p.recording = false
}
