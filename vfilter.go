package gst

/*
#cgo pkg-config: gstreamer-video-1.0

#include "callback.h"

*/
import "C"
import (
	"log"
	"unsafe"
)

var videoFilterStore = map[uintptr]VideoFilter{}

const propertyName = "transform-ip-callback"
const propCallerId = "caller-id"

//type ChainCallback func(plugin *VideoFilter, pad *Pad, buf []byte) int
type VideoFilter interface {
	TransformIP(frame *VideoFrame) error
}

type VideoFrame struct {
	GstVideoFrame *C.GstVideoFrame
	Pixels        *C.guint8
	Stride        uint
	PixelStride   uint
	Height        int
	Width         int
}

func (vfr *VideoFrame) InitData() bool {
	if vfr.GstVideoFrame == nil {
		return false
	}
	frame := vfr.GstVideoFrame

	if C.gst_video_frame_map(frame, nil, nil, C.GST_MAP_WRITE) != 0 {
		// vfr.Pixels = C.GST_VIDEO_FRAME_PLANE_DATA(frame, 0)
		// vfr.Stride = C.GST_VIDEO_FRAME_PLANE_STRIDE(frame, 0)
		// vfr.PixelStride = C.GST_VIDEO_FRAME_COMP_PSTRIDE(frame, 0)
		// vfr.Height = C.GST_VIDEO_FRAME_COMP_HEIGHT(frame, 0)
		// vfr.Width = C.GST_VIDEO_FRAME_COMP_WIDTH(frame, 0)
	}
	return true
}

func (vfr *VideoFrame) Close() {
	C.gst_video_frame_unmap(vfr.GstVideoFrame)
}

type VideoFilterPlugin struct {
	Element
}

func (vfp *VideoFilterPlugin) TransformIP(frame *VideoFrame) error {
	return nil
}

//export go_transform_frame_ip
func go_transform_frame_ip(filter *C.GstVideoFilter, frame *C.GstVideoFrame) (ret C.GstFlowReturn) {
	ret = C.GST_FLOW_OK
	callbackID := uintptr(unsafe.Pointer(filter))
	mutex.Lock()
	plugin := videoFilterStore[callbackID]
	mutex.Unlock()
	log.Printf("object %x", callbackID)
	if plugin == nil {
		return C.GST_FLOW_ERROR
	}

	vf := &VideoFrame{
		GstVideoFrame: frame,
	}

	if false == vf.InitData() {
		return C.GST_FLOW_ERROR
	}
	defer vf.Close()

	err := plugin.TransformIP(vf)
	if err != nil {
		ret = C.GST_FLOW_ERROR
	}

	return ret
}

func (e *VideoFilterPlugin) SetCallback() {

	callbackID := uintptr(unsafe.Pointer(e.GstElement))
	mutex.Lock()
	videoFilterStore[callbackID] = e
	mutex.Unlock()
	log.Printf("object %x", callbackID)

	C.X_go_set_callback_id(e.GstElement, C.guint64(callbackID))
	C.X_go_set_callback_transform_ip(e.GstElement)
}

// funcAddr returns function value fn executable code address.
//@see https://habr.com/ru/post/482392/
func funcAddr(fn interface{}) uintptr {
	// emptyInterface is the header for an interface{} value.
	type emptyInterface struct {
		typ   uintptr
		value *uintptr
	}
	e := (*emptyInterface)(unsafe.Pointer(&fn))
	return *e.value
}
