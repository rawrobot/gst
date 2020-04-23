package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/bksworm/gst"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

const (
	wigetName     = "sinkWidget"
	photoSinkName = "photo_sink"
	movieSinkName = "movie_sink"
	QUEUE_SIZE    = 10
)

type Player struct {
	pipe      *gst.Pipeline
	widget    *gtk.Widget
	photoSink *FrameSink
	movieSink *FrameSink
	playing   bool
}

func NewPlayer() (p *Player) {
	p = &Player{}
	return p
}

func (p *Player) Assemble() (err error) {
	//	p.pipe, err = gst.ParseLaunch("v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! " +
	//		"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=" + wigetName +
	//		" t. ! queue !  videoconvert ! video/x-raw,format=I420  " +
	//		" !  tee name=u ! queue ! jpegenc !  appsink name= " + photoSink + " u. " +
	//		" ! queue !   appsink name= " + movieSink)

	// pipelineStr := fmt.Sprintf("v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! "+
	// 	"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=\"%s\" "+
	// 	" t. ! queue !   jpegenc !  appsink name=\"%s\"", wigetName, photoSink)

	p.pipe, err = gst.ParseLaunch("v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! " +
		"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=" + wigetName +
		" t. ! queue !   jpegenc !  appsink name=photo_sink ") //+ photoSinkName)

	if err != nil {
		return errors.Wrap(err, "player")
	}

	sink := p.pipe.GetByName(wigetName)
	if sink == nil {
		err = errors.Wrap(errors.Errorf(" sink %s not found ", wigetName), "player")
		return err
	}

	p.widget, err = getWidget(sink)
	if err != nil {
		//log.Println("Cann't get move area widget!")
		return errors.Wrap(err, "player")
	}

	gstPhotoSink := p.pipe.GetByName(photoSinkName)
	if p.photoSink == nil {
		err = errors.Wrap(errors.Errorf(" sink %s not found ", photoSinkName), "player")
		//log.Println(err.Error())
		return
	}
	p.photoSink = NewFrameSink(gstPhotoSink, QUEUE_SIZE)

	// gstMovieSink := p.pipe.GetByName(movieSink)
	// if p.movieSink == nil {
	// 	err = fmt.Errorf("pipeline:  sink %s not found ", movieSink)
	// 	log.Println(err.Error())
	// 	return
	// }
	// p.movieSink = NewFrameSink(gstMovieSink, QUEUE_SIZE)

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
		p.photoSink.PullCtrl(3)
	}
}

//This routine pulls gstreamer pipeline and save a number of images on demand
func (p *Player) PictureTaker(ctx context.Context, saveToDir string) (err error) {

	go p.jpegSaver(ctx, saveToDir)
	return p.photoSink.Pull(ctx)
}

func (p *Player) jpegSaver(ctx context.Context, dir string) error {
	for jpg := range p.photoSink.frames {
		fullpath := filepath.Join(dir, buildFileName())
		fd, err := os.Create(fullpath)
		log.Println(fullpath)
		if err != nil {
			log.Print(err)
			return err
		}
		fd.Write(jpg.Bytes())
		fd.Close()
		p.photoSink.Return(jpg)
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
		return w, errors.Wrap(err, "player")
	}

	ip := p.(interface{})
	w, ok = ip.(*gtk.Widget)
	if !ok {
		return w, errors.Wrap(errors.New("It is not a *gtk.Widget"), "player")
	}
	return w, err
}
