package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type Log struct {
	//private
	log *C.PurpleLog
}

func newLogFrom(log *C.PurpleLog) *Log {
	this := &Log{}
	this.log = log
	return this
}
