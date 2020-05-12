package gst

import (
	"fmt"
	"testing"
	"time"
)

func TestPluginCallbackSet(t *testing.T) {
	//pipeline, err := ParseLaunch("videotestsrc num-buffers=2 ! video/x-raw,format=I420,width=320,height=240,framerate=2/1 !  gocallback name=gofilter ! fakesink")
	pipeline, err := ParseLaunch("fakesrc num-buffers=2 !  gocallback name=gofilter silent=false ! fakesink")

	if err != nil {
		t.Error("pipeline create error", err)
		t.FailNow()
	}

	fmt.Println(pipeline.Name())

	element := pipeline.GetByName("gofilter")

	if element == nil {
		t.Error("pipe find element error")
		t.FailNow()
	}
	plugin := &Plugin{
		Element: *element,
	}

	plugin.SetOnChainCallback(TestChainCallback)

	pipeline.SetState(StatePlaying)
	time.Sleep(1000000)
	pipeline.SetState(StateNull)
}

func TestPluginCallbackWork(t *testing.T) {

	pipeline, err := ParseLaunch("videotestsrc num-buffers=3 ! video/x-raw,format=I420,width=320,height=240,framerate=25/1 " +
		" ! gocallback name=gofilter ! jpegenc ! multifilesink location=image_%06d.jpg")
	//pipeline, err := ParseLaunch("fakesrc num-buffers=2 !  gocallback name=gofilter silent=false ! fakesink")

	if err != nil {
		t.Error("pipeline create error", err)
		t.FailNow()
	}

	fmt.Println(pipeline.Name())

	element := pipeline.GetByName("gofilter")

	if element == nil {
		t.Error("pipe find element error")
		t.FailNow()
	}
	plugin := &Plugin{
		Element: *element,
	}

	plugin.SetOnChainCallback(TestChainCallback)

	pipeline.SetState(StatePlaying)
	time.Sleep(10000000)
	pipeline.SetState(StateNull)
}

func TestVFCallbackSet(t *testing.T) {

	pipeline, err := ParseLaunch("videotestsrc num-buffers=2 " +
		" !  video/x-raw,format=I420,width=320,height=240,framerate=2/1 " +
		" !  govideocallback name=gofilter ! fakesink")
	//pipeline, err := ParseLaunch("fakesrc num-buffers=2 !  govideocallback name=gofilter ! fakesink")

	if err != nil {
		t.Error("pipeline create error", err)
		t.FailNow()
	}

	fmt.Println(pipeline.Name())

	element := pipeline.GetByName("gofilter")

	if element == nil {
		t.Error("pipe find element error")
		t.FailNow()
	}
	plugin := &VideoIPTransformPlugin{
		Element: *element,
	}

	SetVideoIPTransformCallback(plugin)

	pipeline.SetState(StatePlaying)
	time.Sleep(10000000)
	pipeline.SetState(StateNull)
	PrintMemUsage()
}

func TestVFCallerIdSet(t *testing.T) {

	pipeline, err := ParseLaunch("videotestsrc num-buffers=2 " +
		" !  video/x-raw,format=I420,width=320,height=240,framerate=2/1 " +
		" !  govideocallback name=gofilter ! fakesink")
	//pipeline, err := ParseLaunch("fakesrc num-buffers=2 !  govideocallback name=gofilter ! fakesink")

	if err != nil {
		t.Error("pipeline create error", err)
		t.FailNow()
	}

	fmt.Println(pipeline.Name())

	element := pipeline.GetByName("gofilter")

	if element == nil {
		t.Error("pipe find element error")
		t.FailNow()
	}
	plugin := &VideoIPTransformPlugin{
		Element: *element,
	}

	SetVideoIPTransformCallback(plugin)

	pipeline.SetState(StatePlaying)
	time.Sleep(10000000)
	pipeline.SetState(StateNull)
}

func TestVFCallbackMem(t *testing.T) {
	PrintMemUsage()

	pipeline, err := ParseLaunch("videotestsrc num-buffers=200 " +
		" !  video/x-raw,format=I420,width=320,height=240,framerate=200/1 " +
		" !  govideocallback name=gofilter ! fakesink")

	if err != nil {
		t.Error("pipeline create error", err)
		t.FailNow()
	}

	fmt.Println(pipeline.Name())

	element := pipeline.GetByName("gofilter")

	if element == nil {
		t.Error("pipe find element error")
		t.FailNow()
	}
	plugin := &VideoIPTransformPlugin{
		Element: *element,
	}

	SetVideoIPTransformCallback(plugin)

	pipeline.SetState(StatePlaying)
	time.Sleep(1_000_000_000)
	pipeline.SetState(StateNull)
	PrintMemUsage()
}
