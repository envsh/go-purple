package gloox

/*
#include "gloox_extras.h"
*/
import "C"

// some const maybe

//
func RefillDataForm(df uintptr) uint64 {
	newForm := C.RefillDataForm(C.uint64_t(df))
	return uint64(newForm)
}
