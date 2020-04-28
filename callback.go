package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-video-1.0
#cgo LDFLAGS: -L/.local/lib/x86_64-linux-gnu

#include <gst/gst.h>
#include <gst/video/video.h>
#include <gst/base/gstbasetransform.h>
#include <gst/video/gstvideofilter.h>
*/
import "C"
import (
	"log"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	filterAsErrCause = "videofilter"
	propName         = "callback"
)

type TransformIPCallback func(filter *C.GstVideoFilter, frame *C.GstVideoFrame) C.GstFlowReturn

// funcAddr returns function value fn executable code address.
func funcAddr(fn interface{}) uintptr {
	// emptyInterface is the header for an interface{} value.
	type emptyInterface struct {
		typ   uintptr
		value *uintptr
	}
	e := (*emptyInterface)(unsafe.Pointer(&fn))
	return *e.value
}

func SetCallBack(name string, pl *Pipeline, fn interface{}) error {
	e := pl.GetByName(name)
	if e == nil {
		return errors.Wrap(errors.Errorf("element %s not found", name), filterAsErrCause)
	}
	v := funcAddr(fn)
	e.SetObject(propName, v)
	return nil
}

func TestIPCallback(filter *C.GstVideoFilter, frame *C.GstVideoFrame) C.GstFlowReturn {
	log.Panicln("++++++++++++++++++++++++++++")
	return 0
}
