package purple

/*
#include <stdlib.h>
*/
import "C"

import "unsafe"
import "runtime"
import "sync"

import "log"
import "time"

type CString struct {
	Ptr *C.char
	Str string
	sz  int
	t   time.Time
	pc  uintptr
	di  *debugInfoForCString
}

type debugInfoForCString struct {
	Str string
	sz  int
	t   time.Time
	pc  uintptr
}

func CCString(s string) *CString {
	p := C.CString(s)
	this := &CString{Ptr: p, Str: s, t: time.Now(), sz: len(s)}
	pc, _, _, _ := runtime.Caller(1)
	this.pc = pc

	setFinalizerForCString(this)
	return this
}

func CCalloc(size uint) *CString {
	p := C.calloc(1, C.size_t(size))
	this := &CString{Ptr: (*C.char)(p), t: time.Now(), sz: int(size)}
	pc, _, _, _ := runtime.Caller(1)
	this.pc = pc

	setFinalizerForCString(this)
	return this
}

var triggerTick sync.Once

func setFinalizerForCString(this *CString) {
	runtime.SetFinalizer(this, freeCString)

	triggerTick.Do(finalIterate)
}

func finalIterate() {
	go func() {
		tk := time.Tick(5 * time.Second)
		for {
			select {
			case <-tk:
				runtime.GC()
			}
		}
	}()
}

func freeCString(cs *CString) {
	if cs == nil {
		log.Panicln("wtf")
	}
	if cs.Ptr == nil {
		log.Panicln("wtf", cs)
	}

	if false {
		ostr := cs.Str
		if len(ostr) > 67 {
			ostr = ostr[0:67] + "..."
		}
		fn := runtime.FuncForPC(cs.pc)
		file, line := fn.FileLine(cs.pc)
		name := fn.Name()

		log.Printf("Freeing...%p, %v, %v, %v, T:%v <-@ (%v %v) %v:%v:%v\n",
			cs, cs.Ptr, len(cs.Str), ostr, freedSize+uint64(cs.sz),
			cs.pc, cs.t, file, line, name)
	}
	pp := &cs.Ptr
	p := unsafe.Pointer(cs.Ptr)
	cs.Ptr = nil
	C.free(p)
	pp = (**C.char)((unsafe.Pointer)((uintptr)(0x1)))
	if false {
		println(pp)
	}

	if freeCounter != nil {
		freeCounter()
	}
	freedSize += uint64(cs.sz)
}

func (this *CString) PtrU8() *C.uchar {
	return (*C.uchar)((unsafe.Pointer)(this.Ptr))
}

func (this *CString) PtrVoid() unsafe.Pointer {
	return ((unsafe.Pointer)(this.Ptr))
}

var freeCounter func()
var freedSize uint64
