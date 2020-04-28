package gst

import (
	"fmt"
	"testing"
	"time"
)

func TestCallbackSet(t *testing.T) {

	pipeline, err := ParseLaunch("videotestsrc num-buffers=2 ! video/x-raw,format=I420,width=320,height=240,framerate=2/1 !  govideocallback name=gofilter ! fakesink")

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

	err = SetCallBack("gofilter", pipeline, TestIPCallback)
	if err != nil {
		t.Error("set callback create error", err)
		t.FailNow()
	}
	pipeline.SetState(StatePlaying)
	time.Sleep(1000000)
	pipeline.SetState(StateNull)
}
