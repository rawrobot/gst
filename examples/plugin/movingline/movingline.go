package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bksworm/gst"
)

type LinePlugin struct {
	gst.VideoFilterPlugin
	line int
}

func NewLinePlugin(e *gst.Element) *LinePlugin {
	lp := &LinePlugin{}
	lp.VideoFilterPlugin.Element = *e
	return lp
}

//draws horithontal black line at the midle of frame
func (lp *LinePlugin) TransformIP(vf *gst.VideoFrame) error {

	y := vf.Plane(0)
	p := lp.line * y.Stride * y.PixelStride
	gst.MemSet(y.Pixels[p:p+y.Width], 0)
	lp.line += 1
	if lp.line >= y.Height {
		lp.line = 0
	}
	return nil
}

func main() {
	pipeline, err := gst.ParseLaunch("videotestsrc  pattern=white  " +
		" !  video/x-raw,format=I420,width=320,height=240,framerate=25/1 " +
		" !  govideocallback name=gofilter ! autovideosink")

	if err != nil {
		log.Println("pipeline create error ", err.Error())
		return
	}

	log.Println(pipeline.Name())

	element := pipeline.GetByName("gofilter")

	if element == nil {
		log.Println("pipe find element error")
		return
	}

	plugin := NewLinePlugin(element)
	gst.SetVideoTransformIpCallback(plugin)

	pipeline.SetState(gst.StatePlaying)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// Run Cleanup
		pipeline.SetState(gst.StateNull)
		os.Exit(0)
	}()

	select {}

}
