package main

import (
	"bytes"
	"context"

	"github.com/bksworm/gst"
)

type FrameSrc struct {
	src     *gst.Element
	frames  chan *bytes.Buffer
	shutter chan int
	bp      *gst.BufferPool
}

func NewFrameSrc(src *gst.Element, bp *gst.BufferPool) (fs *FrameSrc) {
	fs = &FrameSrc{}
	fs.src = src
	fs.frames = make(chan *bytes.Buffer, bp.NumPooled()) //this is queue for frames
	fs.bp = bp
	return fs
}

func (fs *FrameSrc) Close() {
	close(fs.frames)
}

func (fs *FrameSrc) Return(bb *bytes.Buffer) {
	fs.bp.Put(bb)
}

//returns a frame or error
func (fs *FrameSrc) PutNoBlock(bb *bytes.Buffer) (err error) {
	select {
	case fs.frames <- bb:
	default:
		err = EFULL
	}
	return err
}

//returns a frame or blocks
func (fs *FrameSrc) Put(bb *bytes.Buffer) {
	fs.frames <- bb
}

//Pulls frames from a pipe
func (fs *FrameSrc) Send(ctx context.Context) (err error) {

	for frame := range fs.frames {
		//we have to read number of shutts here to avoid races
		select {
		case <-ctx.Done():
			return ECANCEL
		default:
			break
		}

		err = fs.src.PushBuffer(frame.Bytes())
		fs.Return(frame)
		if err != nil {
			break
		}
	}
	return err
}
