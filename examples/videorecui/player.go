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
	wigetName        = "sinkWidget"
	photoSinkName    = "photo_sink"
	movieSinkName    = "movie_sink"
	QUEUE_SIZE       = 10
	playerAsErrCause = "player"
)

type Player struct {
	pipe      *gst.Pipeline
	widget    *gtk.Widget
	photoSink *FrameSink
	rec       *VideoRec
	playing   bool
}

func NewPlayer() (p *Player) {
	p = &Player{}
	return p
}

const (
	CMD = "v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! " +
		"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=" + wigetName +
		" t. ! queue !  videoconvert ! video/x-raw,format=I420  " +
		" !  tee name=u ! queue ! jpegenc !  appsink name=" + photoSinkName + " u. " +
		" ! queue !   appsink name=" + movieSinkName +
		" appsrc name=" + movieSrcName + "  stream-type=0 format=time is-live=true do-timestamp=true " +
		"  ! video/x-raw,width=640,height=480,format=I420,framerate=30/1" +
		"  ! x264enc !  matroskamux ! filesink location=video.mkv"
		//" ! jpegenc !  multipartmux ! filesink location=video.mjpeg "
		//" ! x264enc !  splitmuxsink muxer=matroskamux location=video%02d.mkv " +
		//" max-size-time=10000000000 max-size-bytes=1000000"

	CMD0 = "v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! " +
		"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=" + wigetName +
		" t. ! queue !   jpegenc !  appsink name=" + photoSinkName
)

func (p *Player) Assemble() (err error) {
	cmd := CMD
	p.pipe, err = gst.ParseLaunch(cmd)

	// p.pipe, err = gst.ParseLaunch("v4l2src device=/dev/video0  ! video/x-raw,width=640,height=480 ! " +
	// 	"tee name=t !  queue !  videoconvert ! video/x-raw,format=BGRA ! gtksink name=" + wigetName +
	// 	" t. ! queue !   jpegenc !  appsink name= " + photoSink)

	if err != nil {
		return errors.Wrap(err, playerAsErrCause)
	}
	log.Println(cmd)

	sink := p.pipe.GetByName(wigetName)
	if sink == nil {
		return errors.Wrap(errors.Errorf("element %s not found", wigetName), playerAsErrCause)
	}

	p.widget, err = getWidget(sink)
	if err != nil {
		return errors.Wrap(err, playerAsErrCause)
	}

	e := p.pipe.GetByName(photoSinkName)
	if e == nil {
		return errors.Wrap(errors.Errorf("element %s not found", photoSinkName), playerAsErrCause)
	}
	p.photoSink = NewFrameSink(e, QUEUE_SIZE)

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

func (p *Player) StartRecording() {
	if p.playing {
		p.rec.PullCtrl(-1)
	}
}

func (p *Player) StopRecording() {
	if p.playing {
		p.rec.PullCtrl(0)
	}
}

func (p *Player) Close() {
	p.pipe.SetState(gst.StateNull)
	p.playing = false
	p.photoSink.Close()
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

//This routine pulls gstreamer pipeline and save a movie
func (p *Player) MovieMaker(ctx context.Context, saveToDir string) (err error) {
	p.rec = NewVideoRec()
	return p.rec.Start(ctx, p.pipe)
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
