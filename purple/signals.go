package purple

/*
#include <libpurple/purple.h>

static gulong gopurple_signal_connect(char *signal) {
    static int handle;
    void *data = NULL;
    void *instance = purple_connections_get_handle();
    return purple_signal_connect(instance, signal, &handle, PURPLE_CALLBACK(NULL), data);
}

static gulong gopurple_signal_connect2(void *instance, char *signal, void *fn) {
    static int handle;
    void *data = NULL;
    return purple_signal_connect(instance, signal, &handle, PURPLE_CALLBACK(fn), data);
}

*/
import "C"
import "unsafe"

import (
	"log"
)

// 这个函数具有叠加效果，多次调用，则触发多次事件回调
func signalConnect(instance unsafe.Pointer, signal string, fn unsafe.Pointer) int {
	r := C.gopurple_signal_connect2(instance, CCString(signal).Ptr, fn)
	if r == 0 {
		log.Println("error:", signal, r)
	}
	return int(r)
}

func SignalConnectConn(signal string, fn unsafe.Pointer) {
	h := C.purple_connections_get_handle()
	signalConnect(h, signal, fn)
}

func SignalConnectConv(signal string, fn unsafe.Pointer) {
	h := C.purple_conversations_get_handle()
	signalConnect(h, signal, fn)
}
