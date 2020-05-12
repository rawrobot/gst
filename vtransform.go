package gst

/*
#cgo pkg-config: gstreamer-video-1.0

#include "callback.h"

*/
import "C"
import (
	"unsafe"
)

var videoTransformIpStore = map[uintptr]VideoIPTransformer{}
var videoTransformStore = map[uintptr]VideoTransformer{}

type Elementer interface {
	CallbackID() uintptr
	GetElement() Element
}

type VideoIPTransformer interface {
	Elementer
	TransformIP(frame *VideoFrame) error
}

type VideoTransformer interface {
	Elementer
	Transform(inframe, outframe *VideoFrame) error
}

type VideoFrame struct {
	GstVideoFrame *C.GstVideoFrame
	Format        int
	NPlanes       int
	Size          int
}

func NewVideoFrame(frame *C.GstVideoFrame) (vfr *VideoFrame) {
	vfr = &VideoFrame{GstVideoFrame: frame}

	vfr.NPlanes = int(C.get_frame_n_planes(frame))
	vfr.Size = int(frame.info.size)
	vfr.Format = int(C.get_frame_format(frame))

	return vfr
}

func (vfr *VideoFrame) Plane(n int) (plane VideoFramePlane) {

	frame := vfr.GstVideoFrame

	//FIXME: n > vfr.NPlanes check

	plane.Stride = int(C.get_frame_stride(frame, C.int(n)))
	plane.PixelStride = int(C.get_frame_pixel_stride(frame, C.int(n)))
	plane.Height = int(C.get_frame_h(frame, C.int(n)))
	plane.Width = int(C.get_frame_w(frame, C.int(n)))
	plane.Size = plane.PixelStride * plane.Width * plane.Height
	plane.Pixels = nonCopyGoBytes(uintptr(unsafe.Pointer(C.get_frame_data(frame, C.int(n)))), plane.Size)
	return plane
}

type VideoFramePlane struct {
	Pixels      []byte
	Size        int
	Stride      int
	PixelStride int
	Height      int
	Width       int
}

func (vfr *VideoFramePlane) Close() {
	vfr.Pixels = nil
}

type VideoIPTransformPlugin struct {
	Element
}

func (vfp *VideoIPTransformPlugin) GetElement() Element {
	return vfp.Element
}

//"Pure virtual method :)
func (vfp *VideoIPTransformPlugin) TransformIP(vf *VideoFrame) error {
	return nil
}

func (vfp *VideoIPTransformPlugin) CallbackID() uintptr {
	return uintptr(unsafe.Pointer(vfp.GstElement))
}

//export go_transform_frame_ip
func go_transform_frame_ip(filter *C.GstVideoFilter, frame *C.GstVideoFrame) (ret C.GstFlowReturn) {
	ret = C.GST_FLOW_OK
	callbackID := uintptr(unsafe.Pointer(filter))
	mutex.Lock()
	plugin := videoTransformIpStore[callbackID]
	mutex.Unlock()

	if plugin == nil {
		return C.GST_FLOW_ERROR
	}

	vf := NewVideoFrame(frame)
	err := plugin.TransformIP(vf)

	if err != nil {
		ret = C.GST_FLOW_ERROR
	}

	return ret
}

//export go_transform_frame
func go_transform_frame(filter *C.GstVideoFilter, inframe, outframe *C.GstVideoFrame) (ret C.GstFlowReturn) {
	ret = C.GST_FLOW_OK
	callbackID := uintptr(unsafe.Pointer(filter))
	mutex.Lock()
	plugin := videoTransformStore[callbackID]
	mutex.Unlock()

	if plugin == nil {
		return C.GST_FLOW_ERROR
	}

	inVf := NewVideoFrame(inframe)
	outVf := NewVideoFrame(outframe)
	err := plugin.Transform(inVf, outVf)

	if err != nil {
		ret = C.GST_FLOW_ERROR
	}

	return ret
}

func SetVideoIPTransformCallback(vti VideoIPTransformer) {

	callbackID := vti.CallbackID()
	mutex.Lock()
	videoTransformIpStore[callbackID] = vti
	mutex.Unlock()
	e := vti.GetElement().GstElement
	C.X_go_set_callback_transform_ip(e)
}

func SetVideoTransformCallback(vti VideoTransformer) {
	callbackID := vti.CallbackID()
	mutex.Lock()
	videoTransformStore[callbackID] = vti
	mutex.Unlock()
	e := vti.GetElement().GstElement
	C.X_go_set_callback_transform(e)
}

//https://eli.thegreenplace.net/2019/passing-callbacks-and-pointers-to-cgo/
type tSlice struct {
	addr uintptr
	len  int
	cap  int
}

func MemSet(b []byte, v byte) {
	slice := (*tSlice)(unsafe.Pointer(&b))
	C.memset(unsafe.Pointer(slice.addr), C.int(v), C.size_t(slice.len))
}
