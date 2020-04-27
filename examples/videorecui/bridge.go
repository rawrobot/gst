package main

import (
	"bytes"
	"context"
	"log"
	"runtime"

	"github.com/bksworm/gst"
)

//Does somthing with frame
type FrameProcessor func(bb *bytes.Buffer)

type Bridge struct {
	sink    *gst.Element
	src     *gst.Element
	frames  chan *bytes.Buffer
	shutter chan int
	bp      *gst.BufferPool
	fp      FrameProcessor
}

func NewBridge(in, out *gst.Element, queueSize int) (b *Bridge) {
	b = &Bridge{}
	b.sink = in
	b.src = out
	b.frames = make(chan *bytes.Buffer, queueSize) //this is queue for frames
	b.shutter = make(chan int, 1)
	b.bp = gst.NewBufferPool(queueSize + 2) //  +2 just to avoid races
	return b
}

func (b *Bridge) Close() {
	close(b.frames)
	close(b.shutter)
}

//Controls pulling
// n== 0 stop, n >0 get n frames, n< 0 get frames untill stop
func (b *Bridge) PullCtrl(n int) (err error) {
	select {
	case b.shutter <- n:
	default:
		err = EBUSY
	}
	return err
}

//Starts in and out routines
func (b *Bridge) ConnectPipes(ctx context.Context) {
	go b.in(ctx)
	go b.out(ctx)
}

//Pulls frames from a pipe
func (b *Bridge) in(ctx context.Context) (err error) {
	var (
		sampleData *bytes.Buffer
		n          int
	)

FOR_EXIT:
	for {
		//we have to read number of shutts here to avoid races
		select {
		case n = <-b.shutter:
			log.Printf("%d samples to take ", n)
		case <-ctx.Done():
			err = ECANCEL
			break FOR_EXIT
		default:
			break
		}

		if n == 0 {
			sampleData = nil // if no samples we need to take, just skip
		} else {
			sampleData = b.bp.Get()
		}

		err = b.sink.PullSampleBB(sampleData)
		if err != nil {
			if err == gst.EOS && b.sink.GetState() != gst.StatePlaying { //if pipeline is paused
				//we should not call too often in such case
				//log.Println("Skip input")
				runtime.Gosched()
				continue
			} else {
				break
			}
		}

		if n != 0 {
			select {
			case b.frames <- sampleData: //send image to jpegSaver()
				log.Printf("in  %d", sampleData.Len())
				if n > 0 {
					// n > 0 means a number of shuts we need
					// n == -1 we make a movie
					n -= 1 // so we have taken one
				}
			default:
				err = EFULL
				log.Println(err.Error())
			}
		}
	}
	log.Println(err.Error())
	return err
}

//puts frames to the second pipe
func (b *Bridge) out(ctx context.Context) (err error) {
FOR_EXIT:
	for frame := range b.frames {
		//we have to read number of shutts here to avoid races
		select {
		case <-ctx.Done():
			err = ECANCEL
			break FOR_EXIT
		default:
			break
		}

		if b.fp != nil {
			b.fp(frame)
		}

		err = b.src.PushBuffer(frame.Bytes())
		log.Printf("out %d", frame.Len())
		b.bp.Put(frame)
		if err != nil {
			log.Println(err.Error())
			break
		}
	}
	log.Println(err.Error())
	return err
}
