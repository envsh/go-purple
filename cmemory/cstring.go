package cmemory

/*
#include <stdlib.h>
*/
import "C"

import "unsafe"
import "runtime"

import "log"

type CString struct {
	Ptr *C.char
	Str string
}

func CCString(s string) *CString {
	p := C.CString(s)
	this := &CString{p, s}

	runtime.SetFinalizer(this, freeCString)
	return this
}

func CCalloc(size uint) *CString {
	p := C.calloc(1, C.size_t(size))
	this := &CString{Ptr: (*C.char)(p)}

	runtime.SetFinalizer(this, freeCString)
	return this
}

func freeCString(cs *CString) {
	if false {
		log.Println("string freeing...", cs, cs.Str)
	}
	p := unsafe.Pointer(cs.Ptr)
	cs.Ptr = nil
	C.free(p)

	if freeCounter != nil {
		freeCounter()
	}
}

func (this *CString) PtrU8() *C.uchar {
	return (*C.uchar)((unsafe.Pointer)(this.Ptr))
}

func (this *CString) PtrVoid() unsafe.Pointer {
	return ((unsafe.Pointer)(this.Ptr))
}

func (this *CString) PtrV() unsafe.Pointer {
	return ((unsafe.Pointer)(this.Ptr))
}

var freeCounter func()
