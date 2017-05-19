package gloox

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func freeAll(args []*C.char) {
	for _, a := range args {
		C.free((unsafe.Pointer)(a))
	}
}

func c2gobool(ok C.int) bool {
	if ok == 0 {
		return false
	}
	return true
}

func go2cbool(ok bool) C.int {
	if ok {
		return 1
	}
	return 0
}

func c2goboolx(ok interface{}) bool {
	return false
}

func SwigEnum(gov uintptr) int {
	r := int(*(*C.int)((unsafe.Pointer)(gov)))
	return r
}

func SwigEnumAsInt(gov SwigWrapType) int {
	r := int(*(*C.int)((unsafe.Pointer)(gov.Swigcptr())))
	return r
}

type SwigWrapType interface {
	Swigcptr() uintptr
}

// 也可以用于其他的非类指针
func DeleteSwigEnum(gov SwigWrapType) {
	enump := (*C.int)((unsafe.Pointer)(gov.Swigcptr()))
	C.free((unsafe.Pointer)(enump))
}
