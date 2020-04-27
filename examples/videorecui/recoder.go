package main

import (
	"context"
	"log"

	"github.com/bksworm/gst"
	"github.com/pkg/errors"
)

const (
	movieSrcName      = "frameSrc"
	recoderAsErrCause = "player"
)

type VideoRec struct {
	pipe   *gst.Pipeline
	out    *gst.Element
	in     *gst.Element
	bridge *Bridge
}

func NewVideoRec() (vrec *VideoRec) {
	vrec = &VideoRec{}
	return vrec
}

func (p *VideoRec) Close() {
	p.bridge.Close()
}

func (p *VideoRec) Assemble(pipe *gst.Pipeline) (err error) {
	log.Printf("%+v", pipe)
	p.in = pipe.GetByName(movieSinkName)
	if p.in == nil {
		return errors.Wrap(errors.Errorf("element %s not found", movieSinkName), recoderAsErrCause)
	}

	p.out = pipe.GetByName(movieSrcName)
	if p.out == nil {
		return errors.Wrap(errors.Errorf("element %s not found", movieSrcName), recoderAsErrCause)
	}
	p.pipe = pipe
	return err
}

func (p *VideoRec) Start(ctx context.Context, pipe *gst.Pipeline) (err error) {
	err = p.Assemble(pipe)
	if err != nil {
		log.Println(err)
		return err
	}
	p.bridge = NewBridge(p.in, p.out, QUEUE_SIZE)
	p.bridge.ConnectPipes(ctx)
	//p.PullCtrl(-1)

	return err
}

//Controls pulling
// n== 0 stop, n >0 get n frames, n< 0 get frames untill stop
func (p *VideoRec) PullCtrl(n int) (err error) {
	return p.bridge.PullCtrl(n)
}
