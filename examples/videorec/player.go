package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/bksworm/gst"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const wigetName = "sinkWidget"
const photoSink = "photoSink"

type VideoRec struct {
	pipe      *gst.Pipeline
	widget    *gtk.Widget
	photoSink *gst.Element
	playing   bool
	frames    chan []byte
	shutter   chan int
}

func NewVideoRec() (p *VideoRec) {
	p = &VideoRec{}
	p.frames = make(chan []byte, 10) //this is queue for 10 jpeg images
	p.shutter = make(chan int, 1)
	return p
}

func (p *VideoRec) Assamble() (err error) {
	p.pipe, err = gst.ParseLaunch("v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! " +
		"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=" + wigetName +
		" t. ! queue !   jpegenc !  appsink name= " + photoSink)

	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	sink := p.pipe.GetByName(wigetName)
	if sink == nil {
		return fmt.Errorf("pipeline: %w", errors.New("sink with name "+wigetName+" not found"))
	}

	p.widget, err = getWidget(sink)
	if err != nil {
		log.Println("Cann't get move area widget!")
		return fmt.Errorf("pipeline: %s", err.Error())
	}
	p.photoSink = p.pipe.GetByName(photoSink)
	if p.photoSink == nil {
		err = fmt.Errorf("pipeline:  sink %s not found ", photoSink)
		log.Println(err.Error())
		return
	}
	return err
}

func (p *VideoRec) Play() {
	p.pipe.SetState(gst.StatePlaying)
	p.playing = true
}

func (p *VideoRec) Pause() {
	p.pipe.SetState(gst.StatePaused)
	p.playing = false
}

func (p *VideoRec) TakePicture() {
	if p.playing {
		//it will save 3 samples in a row
		p.shutter <- 3
	}
}

//This routine pulls gstreamer pipeline and save a number of images on demand
func (p *VideoRec) PictureTaker(saveToDir string) (err error) {

	var (
		s *gst.Sample
		n int
	)
	go p.jpegSaver(saveToDir) //start save in separate thread to balance load

	for {
		//we have to read number of shutts here to avoid races
		select {
		case n = <-p.shutter:
			log.Printf("%d samples to take ", n)
		default:
			break
		}

		s, err = p.photoSink.PullSampleOrSkip(n == 0) // if no samples to take just skip
		if err != nil {
			if err == gst.EOS && p.playing == false { //if pipeline is paused
				//TODO:  we should not call pull!
				continue
			} else {
				break
			}
		}

		if n != 0 {
			//log.Printf("samples %d", n)
			//log.Printf("image size %d", len(s.Data))
			select {
			case p.frames <- s.Data: //send image to jpegSaver()
				n -= 1
			default:
				err = errors.New("Something bad in PictureTaker")
				log.Println(err.Error())
			}
		}
	}
	return
}

func (p *VideoRec) jpegSaver(dir string) error {
	for jpg := range p.frames {
		fullpath := filepath.Join(dir, buildFileName())
		fd, err := os.Create(fullpath)
		if err != nil {
			log.Print(err)
			return err
		}
		fd.Write(jpg)
		fd.Close()
	}
	return nil
}

// with my restect to:
// see: thttps://codereview.stackexchange.com/questions/132025/create-a-new-unique-file
var debugCounter uint64

func nextDebugId() string {
	return fmt.Sprintf("%d", atomic.AddUint64(&debugCounter, 1))
}

func buildFileName() string {
	return time.Now().Format("20060102150405") + "_" + nextDebugId() + ".jpg"
}

//the most time  of coding is spent here due to  go type system and memory model
func getWidget(e *gst.Element) (w *gtk.Widget, err error) {
	var ok bool
	obj := glib.Take(e.AsObj())
	p, err := obj.GetProperty("widget")
	if err != nil {
		return w, errors.New("Element doesn't have  widget property!")
	}

	ip := p.(interface{})
	w, ok = ip.(*gtk.Widget)
	if !ok {
		return w, errors.New("It is not a *gtk.Widget")
	}
	return w, err
}
