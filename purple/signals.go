package purple

/*
#include <libpurple/purple.h>

static gulong gopurple_signal_connect(char *signal) {
    static int handle;
    void *data = NULL;
    void *instance = purple_connections_get_handle();
    return purple_signal_connect(instance, signal, &handle, PURPLE_CALLBACK(NULL), data);
}

*/
import "C"

func SignalConnect(signal string, cbfn func()) {
	C.gopurple_signal_connect(CCString(signal).Ptr)
}
