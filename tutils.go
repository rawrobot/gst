package gst

/*
#include <string.h>
#include <stdlib.h>

void malloc_copy_test(char *src, int size){
	char* dst = malloc(size) ;
	memcpy(dst, src, size) ;
	free(dst) ;
}
*/
import "C"

import (
	//"log"
	"unsafe"
)

func go_malloc_copy_test(src []byte, size int) {
	C.malloc_copy_test((*C.char)(unsafe.Pointer(&src[0])), C.int(size))
}

func make_buf_copy(src []byte, size int) []byte {
	dst := make([]byte, size)
	copy(dst, src)
	return dst
}

func buf_copy(src []byte, dst []byte, size int) []byte {
	if len(dst) < size {
		res := make([]byte, size)
		copy(res, src)
		//log.Println("realloc")
		return res
	}
	copy(dst, src)
	return dst
}
