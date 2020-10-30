package gst

/*
#cgo pkg-config: gstreamer-1.0
#include "gst.h"
*/
import "C"

import (
    "fmt"
    "runtime"
    "unsafe"
)

type Structure struct {
	C *C.GstStructure
}

func NewStructure(name string) (structure *Structure) {
	CName := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	CGstStructure := C.gst_structure_new_empty(CName)

	structure = &Structure{
		C: CGstStructure,
	}

	runtime.SetFinalizer(structure, func(structure *Structure) {
		C.gst_structure_free(structure.C)
	})

	return structure
}

func errNoSuchField(t, name string) error {
    return fmt.Errorf("structure does not have a %s named %s", t, name)
}

func (s *Structure) GetName() string {
    return C.GoString(C.gst_structure_get_name(s.C))
}

func (s *Structure) GetBool(name string) (bool, error) {
    var out C.gboolean

    if C.FALSE == C.gst_structure_get_boolean(s.C, C.CString(name), &out) {
        return false, errNoSuchField("bool", name)
    }

    if out == C.TRUE {
        return true, nil
    }

    return false, nil
}

func (s *Structure) GetInt(name string) (int, error) {
    var out C.gint

    if C.FALSE == C.gst_structure_get_int(s.C, C.CString(name), &out) {
        return 0, errNoSuchField("int", name)
    }

    return int(out), nil
}

func (s *Structure) GetInt64(name string) (int64, error) {
    var out C.gint64

    if C.FALSE == C.gst_structure_get_int64(s.C, C.CString(name), &out) {
        return 0, errNoSuchField("int64", name)
    }

    return int64(out), nil
}

func (s *Structure) GetUint(name string) (uint, error) {
    var out C.guint

    if C.FALSE == C.gst_structure_get_uint(s.C, C.CString(name), &out) {
        return 0, errNoSuchField("uint", name)
    }

    return uint(out), nil
}

func (s *Structure) GetString(name string) (string, error) {
    out := C.gst_structure_get_string(s.C, C.CString(name))
    if out == nil {
        return "", errNoSuchField("string", name)
    }

    return C.GoString(out), nil
}

func (s *Structure) SetValue(name string, value interface{}) {

	CName := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.g_free(C.gpointer(unsafe.Pointer(CName)))

	switch val := value.(type) {
	case string:
		str := (*C.gchar)(unsafe.Pointer(C.CString(val)))
		defer C.g_free(C.gpointer(unsafe.Pointer(str)))
		C.X_gst_structure_set_string(s.C, CName, str)
	case int:
		C.X_gst_structure_set_int(s.C, CName, C.gint(val))
	case uint32:
		C.X_gst_structure_set_uint(s.C, CName, C.guint(val))
	case bool:
		var v int
		if val {
			v = 1
		} else {
			v = 0
		}
		C.X_gst_structure_set_bool(s.C, CName, C.gboolean(v))
	}
}

func (s *Structure) ToString() (str string) {
	Cstr := C.gst_structure_to_string(s.C)
	str = C.GoString((*C.char)(unsafe.Pointer(Cstr)))
	C.g_free((C.gpointer)(unsafe.Pointer(Cstr)))

	return
}
