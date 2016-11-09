package purple

/*
#cgo pkg-config: glib-2.0 purple
#cgo CFLAGS: -g -O0 -DHAVE_DBUS -DPURPLE_PLUGINS -DGOMAXPROCS=1

#include <unistd.h>
#include <sys/syscall.h>
#include <stdint.h>
#include <pthread.h>
static uint64_t MyTid() { return pthread_self(); }
static uint64_t MyTid2() { return syscall(sizeof(void*)==4?224:186); }
*/
import "C"

func MyTid() uint64 {
	return uint64(C.MyTid())
}

func MyTid2() uint64 {
	return uint64(C.MyTid2())
}
