package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#cgo LDFLAGS: -L/.local/lib/x86_64-linux-gnu
#include "gst.h"
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

var (
	mutex         sync.Mutex
	callbackStore = map[uint64]*Element{}
	ErrEOS        = errors.New("EOS")
)

type PadAddedCallback func(element *Element, pad *Pad)

type StateOptions int

const (
	StateVoidPending StateOptions = C.GST_STATE_VOID_PENDING
	StateNull        StateOptions = C.GST_STATE_NULL
	StateReady       StateOptions = C.GST_STATE_READY
	StatePaused      StateOptions = C.GST_STATE_PAUSED
	StatePlaying     StateOptions = C.GST_STATE_PLAYING
)

type Element struct {
	GstElement *C.GstElement
	onPadAdded PadAddedCallback
	callbackID uint64
}

func (e *Element) Name() (name string) {
	n := (*C.char)(unsafe.Pointer(C.gst_object_get_name((*C.GstObject)(unsafe.Pointer(e.GstElement)))))
	if n != nil {
		name = string(nonCopyCString(n, C.int(C.strlen(n))))
	}
	return
}

func (e *Element) Link(dst *Element) bool {

	result := C.gst_element_link(e.GstElement, dst.GstElement)
	return result == C.TRUE
}

func (e *Element) GetPadTemplate(name string) (padTemplate *PadTemplate) {

	n := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.g_free(C.gpointer(unsafe.Pointer(n)))

	CPadTemplate := C.gst_element_class_get_pad_template(C.X_GST_ELEMENT_GET_CLASS(e.GstElement), n)
	padTemplate = &PadTemplate{
		C: CPadTemplate,
	}

	return
}

func (e *Element) GetRequestPad(padTemplate *PadTemplate, name string, caps *Caps) (pad *Pad) {

	var n *C.gchar
	var c *C.GstCaps

	if name == "" {
		n = nil
	} else {
		n = (*C.gchar)(unsafe.Pointer(C.CString(name)))
		defer C.g_free(C.gpointer(unsafe.Pointer(n)))
	}
	if caps == nil {
		c = nil
	} else {
		c = caps.caps
	}
	CPad := C.gst_element_request_pad(e.GstElement, padTemplate.C, n, c)
	pad = &Pad{
		pad: CPad,
	}

	return
}

func (e *Element) GetStaticPad(name string) (pad *Pad) {

	n := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.g_free(C.gpointer(unsafe.Pointer(n)))
	CPad := C.gst_element_get_static_pad(e.GstElement, n)
	pad = &Pad{
		pad: CPad,
	}

	return
}

func (e *Element) AddPad(pad *Pad) bool {
	Cret := C.gst_element_add_pad(e.GstElement, pad.pad)
	return Cret == 1
}

func (e *Element) GetClockBaseTime() uint64 {
	CClockTime := C.gst_element_get_base_time(e.GstElement)
	return uint64(CClockTime)
}

func (e *Element) GetClock() (gstClock *Clock) {
	CElementClock := C.gst_element_get_clock(e.GstElement)
	gstClock = &Clock{
		C: CElementClock,
	}

	runtime.SetFinalizer(gstClock, func(gstClock *Clock) {
		C.gst_object_unref(C.gpointer(unsafe.Pointer(gstClock.C)))
	})
	return gstClock
}

func (e *Element) PushBuffer(data []byte) error {

	b := C.CBytes(data)
	defer C.free(b)

	gstReturn := C.X_gst_app_src_push_buffer(e.GstElement, b, C.int(len(data)))
	if gstReturn != C.GST_FLOW_OK {
		return errors.New("could not push buffer on appsrc element")
	}

	return nil
}

func (e *Element) PullSample() (sample *Sample, err error) {

	CGstSample := C.gst_app_sink_pull_sample((*C.GstAppSink)(unsafe.Pointer(e.GstElement)))
	if CGstSample == nil {
		err = errors.New("could not pull a sample from appsink")
		return
	}

	gstBuffer := C.gst_sample_get_buffer(CGstSample)

	if gstBuffer == nil {
		err = errors.New("could not pull a sample from appsink")
		return
	}

	mapInfo := (*C.GstMapInfo)(unsafe.Pointer(C.malloc(C.sizeof_GstMapInfo)))
	defer C.free(unsafe.Pointer(mapInfo))

	if int(C.X_gst_buffer_map(gstBuffer, mapInfo)) == 0 {
		return sample, fmt.Errorf("could not map gstBuffer %#v", gstBuffer)
	}

	CData := (*[1 << 30]byte)(unsafe.Pointer(mapInfo.data))
	data := make([]byte, int(mapInfo.size))
	copy(data, CData[:])

	duration := uint64(C.X_gst_buffer_get_duration(gstBuffer))

	sample = &Sample{
		Data:     data,
		Duration: duration,
	}

	C.gst_buffer_unmap(gstBuffer, mapInfo)
	C.gst_sample_unref(CGstSample)

	return
}

// appsink
func (e *Element) IsEOS() bool {
	Cbool := C.gst_app_sink_is_eos((*C.GstAppSink)(unsafe.Pointer(e.GstElement)))
	return Cbool == 1
}

func (e *Element) SetObject(name string, value interface{}) {

	cname := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.g_free(C.gpointer(unsafe.Pointer(cname)))
	switch val := value.(type) {
	case string:
		str := (*C.gchar)(unsafe.Pointer(C.CString(val)))
		defer C.g_free(C.gpointer(unsafe.Pointer(str)))
		C.X_gst_g_object_set_string(e.GstElement, cname, str)
	case int:
		C.X_gst_g_object_set_int(e.GstElement, cname, C.gint(val))
	case uint32:
		C.X_gst_g_object_set_uint(e.GstElement, cname, C.guint(val))
	case uint64:
		C.X_gst_g_object_set_uint64(e.GstElement, cname, C.guint64(val))
	case bool:
		var cvalue int
		if val {
			cvalue = 1
		} else {
			cvalue = 0
		}
		C.X_gst_g_object_set_bool(e.GstElement, cname, C.gboolean(cvalue))
	case *Caps:
		caps := val
		C.X_gst_g_object_set_caps(e.GstElement, cname, caps.caps)
	case *Structure:
		structure := val
		C.X_gst_g_object_set_structure(e.GstElement, cname, structure.C)
	}
}

func (e *Element) cleanCallback() {

	if e.onPadAdded == nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	delete(callbackStore, e.callbackID)
}

//export go_callback_new_pad
func go_callback_new_pad(CgstElement *C.GstElement, CgstPad *C.GstPad, callbackID C.guint64) {

	mutex.Lock()
	element := callbackStore[uint64(callbackID)]
	mutex.Unlock()

	if element == nil {
		return
	}

	callback := element.onPadAdded

	pad := &Pad{
		pad: CgstPad,
	}

	callback(element, pad)
}

func (e *Element) SetPadAddedCallback(callback PadAddedCallback) {
	e.onPadAdded = callback

	var callbackID uint64
	mutex.Lock()
	for {
		callbackID = rand.Uint64()
		if callbackStore[callbackID] != nil {
			continue
		}
		callbackStore[callbackID] = e
		break
	}
	mutex.Unlock()

	e.callbackID = callbackID

	detailedSignal := (*C.gchar)(C.CString("pad-added"))
	defer C.free(unsafe.Pointer(detailedSignal))

	runtime.SetFinalizer(e, func(e *Element) {
		e.cleanCallback()
	})

	C.X_g_signal_connect(e.GstElement, detailedSignal, C.guint64(callbackID))
}

func ElementFactoryMake(factoryName string, name string) (e *Element, err error) {
	var pName *C.gchar

	pFactoryName := (*C.gchar)(unsafe.Pointer(C.CString(factoryName)))
	defer C.g_free(C.gpointer(unsafe.Pointer(pFactoryName)))
	if name == "" {
		pName = nil
	} else {
		pName = (*C.gchar)(unsafe.Pointer(C.CString(name)))
		defer C.g_free(C.gpointer(unsafe.Pointer(pName)))
	}
	gstElt := C.gst_element_factory_make(pFactoryName, pName)

	if gstElt == nil {
		err = fmt.Errorf("could not create a GStreamer element factoryName %s, name %s", factoryName, name)
		return
	}

	e = &Element{
		GstElement: gstElt,
	}

	return
}

func nonCopyGoBytes(ptr uintptr, length int) []byte {
	var slice []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = length
	header.Len = length
	header.Data = ptr
	return slice
}

func nonCopyCString(data *C.char, size C.int) []byte {
	return nonCopyGoBytes(uintptr(unsafe.Pointer(data)), int(size))
}

//Added by BKSWORM
func (e *Element) AsObj() unsafe.Pointer {
	return unsafe.Pointer(e.GstElement)
}

//Pulls sanple to keep pipe run, but if skeep doesn't  copy it just unref() and returns nil
func (e *Element) PullSampleOrSkip(buf []byte, skip bool) (result []byte, err error) {
	//All rendered buffers will be put in a queue so that the application can pull samples at its own rate.
	// Note that when the application does not pull samples fast enough, the queued buffers
	// could consume a lot of memory, especially when dealing with raw video frames.
	CGstSample := C.gst_app_sink_pull_sample((*C.GstAppSink)(unsafe.Pointer(e.GstElement)))
	if CGstSample == nil {
		//If an EOS event was received before any buffers, this function returns NULL.
		//Use gst_app_sink_is_eos () to check for the EOS condition.
		if e.IsEOS() {
			err = ErrEOS
		} else {
			err = errors.New("could not gst_app_sink_pull_sample()")
		}
		return
	}
	defer C.gst_sample_unref(CGstSample) //it should work fast for go 1.14 due to defer inlining

	//Do we need this sample right now?
	if skip {
		result = buf
		return
	}

	gstBuffer := C.gst_sample_get_buffer(CGstSample)
	if gstBuffer == nil {
		err = errors.New("could not gst_sample_get_buffer() from appsink")
		return
	}

	var mspInfoBuf [C.sizeof_GstMapInfo]byte //to avoid malloc, but use stack, it's about 2-3 times faster
	mapInfo := (*C.GstMapInfo)(unsafe.Pointer(&mspInfoBuf[0]))
	//mapInfo := (*C.GstMapInfo)(unsafe.Pointer(C.malloc(C.sizeof_GstMapInfo)))
	//defer C.free(unsafe.Pointer(mapInfo))

	if int(C.X_gst_buffer_map(gstBuffer, mapInfo)) == 0 {
		err = fmt.Errorf("could not map gstBuffer %#v", gstBuffer)
		return
	}
	defer C.gst_buffer_unmap(gstBuffer, mapInfo)

	CData := (*[1 << 30]byte)(unsafe.Pointer(mapInfo.data))
	if len(buf) < int(mapInfo.size) {
		data := make([]byte, int(mapInfo.size))
		copy(data, CData[:])
		result = data
	} else {
		buf = buf[:mapInfo.size]
		copy(buf, CData[:])
		result = buf
	}
	/*
		duration := uint64(C.X_gst_buffer_get_duration(gstBuffer))
		sample = &Sample{
			Data:     data,
			Duration: duration,
		}
	*/

	return
}

//Helper function for async pusher
func (e *Element) PushBufferAsync(buffer []byte) error {
	b := C.CBytes(buffer)
	defer C.free(unsafe.Pointer(b))
	//It pushes go buf to pipe line buffer is duped!
	Cbool := C.x_push_buffer_async(e.GstElement, b, C.int(len(buffer)))
	if Cbool == 0 {
		return errors.New("push-buffer error")
	}
	return nil
}

//Pulls sample from the pipeline into *bytes.Buffer.
// To keep pipe running and don't exost memory,  we need to pull samples. But if  we
// doesn't nee them, we may pass nil to the function. In this case sample won't be  copied,
// function  just unref() it  and returns nil error.
func (e *Element) PullSampleBB(bb *bytes.Buffer) (err error) {
	//All rendered buffers will be put in a queue so that the application can pull samples at its own rate.
	// Note that when the application does not pull samples fast enough, the queued buffers
	// could consume a lot of memory, especially when dealing with raw video frames.
	CGstSample := C.gst_app_sink_pull_sample((*C.GstAppSink)(unsafe.Pointer(e.GstElement)))
	if CGstSample == nil {
		//If an EOS event was received before any buffers, this function returns NULL.
		//Use gst_app_sink_is_eos () to check for the EOS condition.
		if e.IsEOS() {
			err = ErrEOS
		} else {
			err = errors.New("could not gst_app_sink_pull_sample()")
		}
		return
	}
	defer C.gst_sample_unref(CGstSample) //it should work fast for go 1.14 due to defer inlining

	//Do we need this sample at all?
	if bb == nil {
		return
	}

	gstBuffer := C.gst_sample_get_buffer(CGstSample)
	if gstBuffer == nil {
		err = errors.New("could not gst_sample_get_buffer() from appsink")
		return
	}
	//To avoid malloc, but use stack, it's about 2-3 times faster
	var mspInfoBuf [C.sizeof_GstMapInfo]byte
	mapInfo := (*C.GstMapInfo)(unsafe.Pointer(&mspInfoBuf[0]))
	if int(C.X_gst_buffer_map(gstBuffer, mapInfo)) == 0 {
		err = fmt.Errorf("could not map gstBuffer %#v", gstBuffer)
		return
	}
	defer C.gst_buffer_unmap(gstBuffer, mapInfo)

	CData := (*[1 << 30]byte)(unsafe.Pointer(mapInfo.data))
	bb.Write(CData[:mapInfo.size])
	return
}

func (e *Element) GetState() StateOptions {
	var gstStateBuf [C.sizeof_GstState]byte
	pState := (*C.GstState)(unsafe.Pointer(&gstStateBuf[0]))

	C.gst_element_get_state(e.GstElement, pState, nil, C.GST_CLOCK_TIME_NONE)
	return StateOptions(*pState)
}

func (e *Element) PushBB(data *bytes.Buffer) (err error) {

	buf := data.Bytes()
	b := (unsafe.Pointer(&buf[0]))

	gstReturn := C.X_gst_app_src_push_buffer(e.GstElement, b, C.int(data.Len()))
	if gstReturn != C.GST_FLOW_OK {
		return errors.New("could not push buffer on appsrc element")
	}

	return err
}

//Emits signal without parametrs
func (e *Element) EmitSignal(signal string) {
	cstr := C.CString(signal)
	defer C.free(unsafe.Pointer(cstr))
	C.x_g_signal_emit_by_name(e.GstElement, cstr)
}
