package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/bksworm/gst"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const wigetName = "sinkWidget"
const photoSink = "photoSink"

type Player struct {
	pipe      *gst.Pipeline
	widget    *gtk.Widget
	photoSink *gst.Element
	playing   bool
	jpegs     chan *bytes.Buffer
	shutter   chan int
	bp        *gst.BufferPool
}

func NewPlayer() (p *Player) {
	p = &Player{}
	p.jpegs = make(chan *bytes.Buffer, 10) //this is queue for 10 jpeg images
	p.shutter = make(chan int, 1)
	p.bp = gst.NewBufferPool(3)
	return p
}

func (p *Player) Assemble() (err error) {
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

func (p *Player) Play() {
	p.pipe.SetState(gst.StatePlaying)
	p.playing = true
}

func (p *Player) Pause() {
	p.pipe.SetState(gst.StatePaused)
	p.playing = false
}

func (p *Player) Stop() {
	p.pipe.SetState(gst.StateReady)
	p.playing = false
}

func (p *Player) Close() {
	p.pipe.SetState(gst.StateNull)
	p.playing = false
}
func (p *Player) TakePicture() {
	if p.playing {
		//it will save 3 samples in a row
		p.shutter <- 3
	}
}

//This routine pulls gstreamer pipeline and save a number of images on demand
func (p *Player) PictureTaker(saveToDir string) (err error) {
	var (
		sampleData *bytes.Buffer
		n          int
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

		if n == 0 {
			sampleData = nil // if no samples we need to take, just skip
		} else {
			sampleData = p.bp.Get()
		}

		err = p.photoSink.PullSampleBB(sampleData)
		if err != nil {
			if err == gst.EOS && p.playing == false { //if pipeline is paused
				//we should not call too often in such case
				runtime.Gosched()
				continue
			} else {
				break
			}
		}

		if n != 0 {
			select {
			case p.jpegs <- sampleData: //send image to jpegSaver()
				n -= 1
			default:
				err = errors.New("Picture queue if fl!")
				log.Println(err.Error())
			}
		}
	}
	return
}

func (p *Player) jpegSaver(dir string) error {
	for jpg := range p.jpegs {
		fullpath := filepath.Join(dir, buildFileName())
		fd, err := os.Create(fullpath)
		log.Println(fullpath)
		if err != nil {
			log.Print(err)
			return err
		}
		fd.Write(jpg.Bytes())
		fd.Close()
		p.bp.Put(jpg)
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
