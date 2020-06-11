package main

import (
	"log"
	"time"

	"github.com/bksworm/gst"
)

type LinePlugin struct {
	gst.VideoIPTransformPlugin
}

func NewLinePlugin(e *gst.Element) *LinePlugin {
	lp := &LinePlugin{}
	lp.VideoIPTransformPlugin.Element = *e
	return lp
}

//draws horithontal black line at the midle of frame
func (lp *LinePlugin) TransformIP(vf *gst.VideoFrame) error {

	for i := 0; i < vf.NPlanes; i++ {
		p := vf.Plane(i)
		log.Printf("plane %d  %dx%d  size %d", i, p.Width, p.Height, p.Size)
	}

	y := vf.Plane(0)
	p := y.Height / 2 * y.Stride * y.PixelStride

	for w := 0; w < y.Width; w++ {
		y.Pixels[p+w] = 0
	}
	//gst.MemSet(y.Pixels[p:p+y.Width], 0) is better way :)
	return nil
}

func main() {
	pipeline, err := gst.ParseLaunch("videotestsrc  pattern=white num-buffers=1 " +
		" !  video/x-raw,format=I420,width=320,height=240,framerate=25/1 " +
		" !  govideocallback name=gofilter ! jpegenc ! filesink location=line.jpeg")

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
	gst.SetVideoIPTransformCallback(plugin)

	pipeline.SetState(gst.StatePlaying)
	time.Sleep(time.Second)
	pipeline.SetState(gst.StateNull)

}
