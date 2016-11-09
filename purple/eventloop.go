package purple

/*
#include <libpurple/purple.h>

extern int timer_timeout_callback(int);
static gboolean timer_timeout_callback_c(void *d)
{
    int tid = 0;
    memcpy(&tid, &d, sizeof(int));
    timer_timeout_callback(tid);
    return TRUE;
}

static guint gopurple_timeout_add(int interval, int tid)
{
    void *d = NULL;
    memcpy(&d, &tid, sizeof(int));
    guint h = purple_timeout_add(interval, timer_timeout_callback_c, d);
    return h;
}
*/
import "C"

// import "unsafe"
import (
	"log"
)

type Timer struct {
	d interface{}
	f func(interface{})
	h uint
}

//export timer_timeout_callback
func timer_timeout_callback(d C.int) C.int {
	t := timers[int(d)]
	t.f(t.d)
	return C.int(0)
}

var timers = make(map[int]*Timer, 0)
var timer_seq int = 0

func TimeoutAdd(interval int, d interface{}, f func(interface{})) int {
	t := &Timer{d: d, f: f}
	timer_seq = timer_seq + 1
	timers[timer_seq] = t
	h := C.gopurple_timeout_add(C.int(interval), C.int(timer_seq))
	t.h = uint(h)

	return timer_seq
}

func TimeoutRemove(h int) bool {
	t := timers[h]
	bret := C.purple_timeout_remove(C.guint(t.h))
	if bret == C.FALSE {
		log.Println("timeout remove failed")
		return false
	}
	delete(timers, h)
	return true
}
