package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/bksworm/gst"
)

const frameSrc = "frameSrc"

var EBF = errors.New("Queue is full!")

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

func (vrec *VideoRec) Assamble() (err error) {
	cmd := fmt.Sprintf("appsrc name=%s  stream-type=0 format=time is-live=true do-timestamp=true "+
		" !  video/x-raw,width=320,height=240,format=RGB,framerate=25/1 "+
		" !  videoconvert ! video/x-raw,format=I420 "+
		" ! x264enc !  matroskamux ! filesink location=appsrc.mkv", frameSrc)
	//" ! jpegenc ! multipartmux  ! filesink location=multipart.mjpeg ", frameSrc)
	log.Println(cmd)
	vrec.pipe, err = gst.ParseLaunch(cmd)

	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	vrec.frameSrc = vrec.pipe.GetByName(frameSrc)
	if vrec.frameSrc == nil {
		err = fmt.Errorf("pipeline:  src %s not found ", frameSrc)
		log.Println(err.Error())
		return
	}
	return err
}

func (vrec *VideoRec) Record() {
	vrec.pipe.SetState(gst.StatePlaying)
	vrec.recording = true
}

func (vrec *VideoRec) Pause() {
	vrec.pipe.SetState(gst.StatePaused)
	vrec.recording = false
}

func (vrec *VideoRec) Recoder() (err error) {
	for frame := range vrec.frames {
		err = vrec.frameSrc.PushBufferAsync(frame)
		if err != nil {
			break
		}
	}
	return nil
}

func (vrec *VideoRec) PushBuffer(frame []byte) error {
	return vrec.frameSrc.PushBuffer(frame)
}
func (vrec *VideoRec) PutFrame(b []byte) (err error) {
	select {
	case vrec.frames <- b:
	default:
		err = EBF
	}
	return
}
