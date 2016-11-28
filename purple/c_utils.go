package purple

/*
#include <libpurple/purple.h>

*/
import "C"
import "unsafe"

// import "reflect"

func c2goBool(ok C.gboolean) bool {
	if ok == C.TRUE {
		return true
	}
	return false
}

func go2cBool(ok bool) C.gboolean {
	if ok {
		return C.TRUE
	}
	return C.FALSE
}

//
type go2cfnty *[0]byte

// 参数怎么传递
func go2cfnp(fn unsafe.Pointer) *[0]byte {
	return go2cfnty(fn)
}
func go2cfn(fn interface{}) *[0]byte {
	// assert(reflect.TypeOf(fn).Kind == reflect.Ptrx)
	return go2cfnp(fn.(unsafe.Pointer))
}
