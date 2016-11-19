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
import (
	"log"
	"strconv"
	"syscall"
)

// TODO unix/linux only
func MyTid() uint64 {
	return uint64(C.MyTid())
}

func MyTid2() uint64 {
	return uint64(C.MyTid2())
}

const PtrSize = 32 << uintptr(^uintptr(0)>>63)
const IntSize = strconv.IntSize
const CIntSize = C.sizeof_int

var archs = map[int]uintptr{
	32: 224, 64: 186,
}

func MyTid3() uint64 {
	r1, r2, err := syscall.Syscall(archs[PtrSize], 0, 0, 0)
	if false {
		log.Println(r1, r2, err)
	}
	return uint64(r1)
}
