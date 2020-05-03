package gst

/*
#cgo pkg-config: gstreamer-video-1.0

#include "callback.h"

*/
import "C"
import (
	"unsafe"
)

var videoFilterStore = map[uintptr]VideoTransformIp{}

const propertyName = "transform-ip-callback"
const propCallerId = "caller-id"

type Elementer interface {
	CallbackID() uintptr
	GetElement() Element
}

type VideoTransformIp interface {
	Elementer
	TransformIP(frame *VideoFrame) error
}

type VideoFrame struct {
	GstVideoFrame *C.GstVideoFrame
	Format        int
	NPlanes       int
	Size          int
}

func (vfr *VideoFrame) InitData() bool {
	if vfr.GstVideoFrame == nil {
		return false
	}
	frame := vfr.GstVideoFrame

	vfr.NPlanes = int(C.get_frame_n_planes(frame))
	vfr.Size = int(frame.info.size)
	vfr.Format = int(C.get_frame_format(frame))

	return true
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

type VideoFilterPlugin struct {
	Element
}

func (vfp *VideoFilterPlugin) GetElement() Element {
	return vfp.Element
}

//"Pure virtual method :)
func (vfp *VideoFilterPlugin) TransformIP(vf *VideoFrame) error {
	return nil
}

func (vfp *VideoFilterPlugin) CallbackID() uintptr {
	return uintptr(unsafe.Pointer(vfp.GstElement))
}

//export go_transform_frame_ip
func go_transform_frame_ip(filter *C.GstVideoFilter, frame *C.GstVideoFrame) (ret C.GstFlowReturn) {
	ret = C.GST_FLOW_OK
	callbackID := uintptr(unsafe.Pointer(filter))
	mutex.Lock()
	plugin := videoFilterStore[callbackID]
	mutex.Unlock()

	if plugin == nil {
		return C.GST_FLOW_ERROR
	}

	vf := &VideoFrame{
		GstVideoFrame: frame,
	}
	vf.InitData()

	err := plugin.TransformIP(vf)

	if err != nil {
		ret = C.GST_FLOW_ERROR
	}

	return ret
}

func SetVideoTransformIpCallback(vti VideoTransformIp) {

	callbackID := vti.CallbackID()
	mutex.Lock()
	videoFilterStore[callbackID] = vti
	mutex.Unlock()
	//log.Printf("object %T", vti)
	e := vti.GetElement().GstElement
	//C.X_go_set_callback_id(e, C.guint64(callbackID))
	C.X_go_set_callback_transform_ip(e)
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
