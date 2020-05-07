package gst

/*
#cgo pkg-config: gstreamer-base-1.0

#include "callback.h"
*/
import "C"
import (
	"log"
	"unsafe"
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
