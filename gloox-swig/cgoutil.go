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
