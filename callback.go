package gst

/*
#cgo pkg-config: gstgocallback
#cgo LDFLAGS: -L/.local/lib/x86_64-linux-gnu

#include "gst.h"
*/
import "C"
import (
	"log"
	"unsafe"
)

const (
	filterAsErrCause = "videofilter"
	propName         = "callback"
)

var pluginStore = map[uintptr]*Plugin{}

type ChainCallback func(plugin *Plugin, pad *Pad, buf []byte) int

type Plugin struct {
	Element
	onChain ChainCallback
}

//export go_callback_chain
func go_callback_chain(CgstPad *C.GstPad, CgstElement *C.GstObject, buf *C.GstBuffer) C.GstFlowReturn {

	mutex.Lock()
	plugin := pluginStore[uintptr(unsafe.Pointer(CgstElement))]
	mutex.Unlock()
	if plugin == nil {
		return C.GST_FLOW_ERROR
	}

	callback := plugin.onChain
	pad := &Pad{
		pad: CgstPad,
	}

	var mspInfoBuf [C.sizeof_GstMapInfo]byte
	mapInfo := (*C.GstMapInfo)(unsafe.Pointer(&mspInfoBuf[0]))
	if int(C.X_gst_buffer_map(buf, mapInfo)) == 0 {
		err = errors.New(fmt.Sprintf("could not map gstBuffer %#v", gstBuffer))
		return C.GST_FLOW_ERROR
	}
	defer C.gst_buffer_unmap(buf, mapInfo)

	var b []byte
	ret := C.GstFlowReturn(callback(plugin, pad, b))

	ret = C.X_gst_pad_push(CgstElement, buf)
	return ret
}

func (e *Plugin) SetOnChainCallback(callback ChainCallback) {
	e.onChain = callback

	callbackID := uintptr(unsafe.Pointer(e.GstElement))
	mutex.Lock()
	pluginStore[callbackID] = e
	mutex.Unlock()
	p := e.GetStaticPad("sink")

	C.X_gst_pad_set_chain_function(p.pad)
}

func TestChainCallback(plugin *Plugin, pad *Pad, buf []byte) int {
	log.Println("++++++++++++++++++++++++++++")
	return 0
}
