package purple

/*
#include <libpurple/purple.h>

*/
import "C"

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
