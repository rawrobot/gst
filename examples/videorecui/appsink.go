package main

import (
	"bytes"
	"context"
	"log"
	"runtime"

	"github.com/pkg/errors"

	"github.com/bksworm/gst"
)

var (
	EBUSY   = errors.New("Busy!")
	EEMPTY  = errors.New("Empty!")
	EFULL   = errors.New("Frame queue is full!")
	ECANCEL = errors.New("Cancel!")
)

type FrameSink struct {
	sink    *gst.Element
	frames  chan *bytes.Buffer
	shutter chan int
	bp      *gst.BufferPool
}

func NewFrameSink(sink *gst.Element, queueSize int) (fs *FrameSink) {
	fs = &FrameSink{}
	fs.sink = sink
	fs.frames = make(chan *bytes.Buffer, queueSize) //this is queue for frames
	fs.shutter = make(chan int, 1)
	fs.bp = gst.NewBufferPool(queueSize + 2) //  +2 just to avoid races
	return fs
}

func (fs *FrameSink) Close() {
	close(fs.frames)
	close(fs.shutter)
}

//Controls pulling
// n== 0 stop, n >0 get n frames, n< 0 get frames untill stop
func (fs *FrameSink) PullCtrl(n int) (err error) {
	select {
	case fs.shutter <- n:
	default:
		err = EBUSY
	}
	return err
}

func (fs *FrameSink) Return(bb *bytes.Buffer) {
	fs.bp.Put(bb)
}

//returns a frame or error
func (fs *FrameSink) GetNoBlock() (bb *bytes.Buffer, err error) {
	select {
	case bb = <-fs.frames:
	default:
		err = EEMPTY
	}
	return bb, err
}

//returns a frame or blocks
func (fs *FrameSink) Get() (bb *bytes.Buffer) {
	bb = <-fs.frames
	return bb
}

//Pulls frames from a pipe
func (fs *FrameSink) Pull(ctx context.Context) (err error) {
	var (
		sampleData *bytes.Buffer
		n          int
	)

	for {
		//we have to read number of shutts here to avoid races
		select {
		case n = <-fs.shutter:
			log.Printf("%d samples to take ", n)
		case <-ctx.Done():
			return ECANCEL
		default:
			break
		}

		if n == 0 {
			sampleData = nil // if no samples we need to take, just skip
		} else {
			sampleData = fs.bp.Get()
		}

		err = fs.sink.PullSampleBB(sampleData)
		if err != nil {
			if err == gst.EOS && fs.sink.GetState() != gst.StatePlaying { //if pipeline is paused
				//we should not call too often in such case
				runtime.Gosched()
				continue
			} else {
				break
			}
		}

		if n != 0 {
			select {
			case fs.frames <- sampleData: //send image to jpegSaver()
				if n > 0 {
					// n > 0 means a number of shuts we need
					// n == -1 we make a movie
					n -= 1 // so we have taken one
				}
			default:
				err = EFULL
				log.Println(err.Error())
				fs.bp.Put(sampleData)
			}
		}
	}
	return
}
