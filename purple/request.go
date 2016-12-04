package purple

/*
#include <libpurple/purple.h>

// fix _ not defined and hidden definition
#include <libintl.h>
#define _(str) gettext((str))

extern void gopurple_request_action_cb(char *user_data, int choice);
extern int get_reqlen();

static void gopurple_request_yes_no(void *handle, char *title, char *primary, char *secondary,
                                    int default_action, PurpleAccount *ac, char *who,
                                    PurpleConversation *conv, char *user_data)
{
     purple_request_yes_no(handle, title, primary, secondary, default_action,
		 ac, who, conv, user_data, gopurple_request_action_cb, gopurple_request_action_cb);
}

static void gopurple_request_ok_cancel(void *handle, char *title, char *primary, char *secondary,
                                    int default_action, PurpleAccount *ac, char *who,
                                    PurpleConversation *conv, char *user_data)
{
     purple_request_ok_cancel(handle, title, primary, secondary, default_action,
		 ac, who, conv, user_data, gopurple_request_action_cb, gopurple_request_action_cb);
}
static void *gopurple_request_accept_cancel(void *handle, char *title, char *primary, char *secondary,
                                    int default_action, PurpleAccount *ac, char *who,
                                    PurpleConversation *conv, char *user_data)
{
     void *r = purple_request_accept_cancel(handle, title, primary, secondary, default_action,
		 ac, who, conv, user_data, gopurple_request_action_cb, gopurple_request_action_cb);
     return r;
}

static void gopurple_request_accept_cancel_with_icon(void *handle, char *title, char *primary,
                                    char *secondary,
                                    int default_action, PurpleAccount *ac, char *who,
                                    PurpleConversation *conv,
                                    gconstpointer icon_data, gsize icon_size,
                                    char *user_data)
{
     purple_request_accept_cancel_with_icon(handle, title, primary, secondary, default_action,
		 ac, who, conv, icon_data, icon_size, user_data,
         gopurple_request_action_cb, gopurple_request_action_cb);
}

extern void *gopurple_request_action_fn(char *title, char *primary,
	                        char *secondary, int default_action,
	                        PurpleAccount *account, char *who,
	                        PurpleConversation *conv, void *user_data,
	                        size_t action_count, va_list *actions);

static void *gopurple_request_action_fn_bridge(const char *title, const char *primary,
	                        const char *secondary, int default_action,
	                        PurpleAccount *account, const char *who,
	                        PurpleConversation *conv, void *user_data,
	                        size_t action_count, va_list actions) {
    // strcpy(0x1, "abccd");
    void *r = gopurple_request_action_fn(title, primary, secondary, default_action,
                            account, who, conv, user_data, action_count, &actions);
    return r;
}

static void gopurple_request_ui_ops_setfns(PurpleRequestUiOps *ops) {
    printf("111len: %d \n", get_reqlen());
    ops->request_action = gopurple_request_action_fn_bridge;
}

*/
import "C"
import "unsafe"

// import "runtime"

import (
	"log"
	//	"sync"
	"reflect"
	"time"
)

// 利用操作系统的内存地址作为标识，让OS维护其值，
// 利用一个内存地址不能同时分配2次及以上的特性
type Request struct {
	seq      *C.char
	userData interface{}
	yescb    func(interface{})
	nocb     func(interface{})
}

var requests = make(map[*C.char]*Request)
var reqno uint64 = 0

func nxtReqNo() uint64 {
	reqno += 1
	return reqno
}

func RequestYesNo(userData interface{}, gc *Connection, title, primary string,
	yescb func(interface{}), nocb func(interface{})) {

	cseq := (*C.char)(C.calloc(1, 1))
	req := &Request{seq: cseq, userData: userData, yescb: yescb, nocb: nocb}
	requests[cseq] = req

	var handle unsafe.Pointer
	var ac *Account
	if gc != nil {
		handle = unsafe.Pointer(gc.conn)
		ac = gc.ConnGetAccount()
	}
	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	who := "who's"
	C.gopurple_request_yes_no(handle, CCString(title).Ptr, CCString(primary).Ptr, nil,
		C.int(default_action), ac.account, CCString(who).Ptr, nil, cseq)
}

func RequestYesNoDemo(userData interface{}, gc *Connection,
	yescb func(interface{}), nocb func(interface{})) {
	title := "it's title"
	primary := "it's primary"
	RequestYesNo(userData, gc, title, primary, yescb, nocb)
}

func RequestOkCancel(userData interface{}, gc *Connection, title, primary string,
	yescb func(interface{}), nocb func(interface{})) {

	cseq := (*C.char)(C.calloc(1, 1))
	req := &Request{seq: cseq, userData: userData, yescb: yescb, nocb: nocb}
	requests[cseq] = req

	var handle unsafe.Pointer
	var ac *Account
	if gc != nil {
		handle = unsafe.Pointer(gc.conn)
		ac = gc.ConnGetAccount()
	}
	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	who := "who's"
	C.gopurple_request_ok_cancel(handle, CCString(title).Ptr, CCString(primary).Ptr, nil,
		C.int(default_action), ac.account, CCString(who).Ptr, nil, cseq)
}

func RequestOkCancelDemo(userData interface{}, gc *Connection,
	yescb func(interface{}), nocb func(interface{})) {
	title := "it's title"
	primary := "it's primary"
	RequestOkCancel(userData, gc, title, primary, yescb, nocb)
}

//export get_reqlen
func get_reqlen() C.int {
	log.Println(len(requests))
	log.Printf("%p, %d, %d\n", requests, MyTid3(), GoID())
	return (C.int)(len(requests))
}

func RequestAcceptCancel(userData interface{}, gc *Connection, title, primary string,
	who string, yescb func(interface{}), nocb func(interface{})) {

	cseq := (*C.char)(C.calloc(1, 1))
	req := &Request{seq: cseq, userData: userData, yescb: yescb, nocb: nocb}
	requests[cseq] = req
	log.Printf("%p, %d, %d\n", requests, MyTid3(), GoID())
	// log.Panicln(123)

	var handle unsafe.Pointer
	var ac *Account
	if gc != nil {
		handle = unsafe.Pointer(gc.conn)
		ac = gc.ConnGetAccount()
	}

	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	// who := "who's"
	time.Sleep(0 * time.Second)
	r := C.gopurple_request_accept_cancel(handle, CCString(title).Ptr, CCString(primary).Ptr, nil,
		C.int(default_action), ac.account, CCString(who).Ptr, nil, cseq)
	// log.Println(len(requests), r, MyTid3()) // maybe invalid pointer 0x1, can not print directly
	// log.Printf("%+v, %p,%p\n", r, cseq, requests)
	rv := (int)((uintptr)(r))
	if rv < 128 {
		// RequestActionCB((unsafe.Pointer)(cseq), rv)
		gopurple_request_action_cb(cseq, C.int(rv))
	}
	// TODO TODO TODO
}
func RequestAcceptCancelDemo(userData interface{}, gc *Connection,
	yescb func(interface{}), nocb func(interface{})) {
	title := "it's title"
	primary := "it's primary"
	who := "who's"
	RequestAcceptCancel(userData, gc, title, primary, who, yescb, nocb)
}

func RequestAcceptCancelWithIcon(userData interface{}, gc *Connection, title, primary string,
	who string, iconData []byte,
	yescb func(interface{}), nocb func(interface{})) {

	cseq := (*C.char)(C.calloc(1, 1))
	req := &Request{seq: cseq, userData: userData, yescb: yescb, nocb: nocb}
	requests[cseq] = req

	var handle unsafe.Pointer
	var ac *Account
	if gc != nil {
		handle = unsafe.Pointer(gc.conn)
		ac = gc.ConnGetAccount()
	}
	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	// who := "who's"

	cicon := (&iconData[0])
	C.gopurple_request_accept_cancel_with_icon(handle, CCString(title).Ptr, CCString(primary).Ptr, nil,
		C.int(default_action), ac.account, CCString(who).Ptr, nil,
		cicon, C.gsize(len(iconData)), cseq)
}
func RequestAcceptCancelWithIconDemo(userData interface{}, gc *Connection,
	iconData []byte,
	yescb func(interface{}), nocb func(interface{})) {
	title := "it's title"
	primary := "it's primary\n\n\n\n\n\n\nend..."
	who := "who's"
	RequestAcceptCancelWithIcon(userData, gc, title, primary, who, iconData, yescb, nocb)
}

//export gopurple_request_action_cb
func gopurple_request_action_cb(userData *C.char, choice C.int) {
	log.Println(userData, choice, len(requests), MyTid3())

	if req, ok := requests[userData]; ok {
		if choice == 0 {
			if req.nocb != nil {
				req.nocb(req.userData)
			}
		} else {
			if req.yescb != nil {
				req.yescb(req.userData)
			}
		}
	} else {
		log.Println("wtf", userData, uint64(*userData), len(requests))
		for k, v := range requests {
			log.Println(k, v)
		}
	}

	// clean up
	if _, ok := requests[userData]; ok {
		log.Println(userData, userData)
		delete(requests, userData)
		C.free(unsafe.Pointer(userData))
	}
}
func RequestActionCB(userData interface{}, choice int) {
	userDatap := userData.(unsafe.Pointer)
	// log.Println(userDatap, len(requests))

	gopurple_request_action_cb((*C.char)(userDatap), C.int(choice))
}

//export gopurple_request_choice_cb
func gopurple_request_choice_cb(userData *C.char, choice C.int) {
	log.Println(userData, choice)
}

//export gopurple_request_fields_cb
func gopurple_request_fields_cb(userData *C.char, fields *C.PurpleRequestFields) {
	log.Println(userData, fields)
}

//export gopurple_request_file_cb
func gopurple_request_file_cb(userData *C.char, filename *C.char) {
	log.Println(userData, filename)
}

////
type RequestUiOps struct {
	RequestAction func(title string, primary string, secondary string, default_action int,
		account *Account, who string, conv *Conversation, user_data interface{},
		action_count int) interface{}

	// private
	requiops *C.PurpleRequestUiOps
}

func newRequestUiOpsFrom(requiops *C.PurpleRequestUiOps) *RequestUiOps {
	this := &RequestUiOps{}
	this.requiops = requiops

	return this
}

func NewRequestUiOpts() *RequestUiOps {
	this := &RequestUiOps{}
	this.requiops = new(C.PurpleRequestUiOps)

	return this
}

var _instance_requiops *RequestUiOps

func RequestSetUiOps(requiops *RequestUiOps) {
	_instance_requiops = requiops
	_instance_requiops.setfns()
	C.purple_request_set_ui_ops(requiops.requiops)
}

func (this *RequestUiOps) setfns() {
	// this.requiops.request_action = go2cfn(C.gopurple_request_action_fn_bridge)
	// this.requiops.request_action = C.gopurple_request_action_fn_bridge
	C.gopurple_request_ui_ops_setfns(this.requiops)
}

func RequestGetUiOps() *RequestUiOps {
	ops := C.purple_request_get_ui_ops()
	if ops == nil {
		return nil
	}
	return &RequestUiOps{requiops: ops}
}

/**
 * Returns the UI operations structure to be used when displaying a
 * request.
 *
 * @return The UI operations structure.
 */

//export gopurple_request_action_fn
func gopurple_request_action_fn(title *C.char, primary *C.char, secondary *C.char, default_action C.int,
	account *C.PurpleAccount, who *C.char, conv *C.PurpleConversation, user_data unsafe.Pointer,
	action_count C.size_t, actions *C.va_list) unsafe.Pointer {
	log.Println("herere", action_count, user_data, MyTid3())
	if _instance_requiops == nil {
		log.Println("not exists")
		return nil
	}

	_, ok := requests[(*C.char)(user_data)]
	log.Println(ok, user_data, (*C.char)(user_data), len(requests))
	log.Printf("%+v, %p\n", user_data, requests)
	// log.Panicln(123)

	this := _instance_requiops
	if this.RequestAction != nil {
		r := this.RequestAction(C.GoString(title), C.GoString(primary), C.GoString(secondary),
			int(default_action), newAccountFrom(account), C.GoString(who),
			newConversationFrom(conv), user_data, int(action_count))
		rty := reflect.TypeOf(r)
		switch rty.Kind() {
		case reflect.Int:
			return (unsafe.Pointer)((uintptr)(r.(int)))
		case reflect.Ptr:
			return r.(unsafe.Pointer)
		}
	}
	return nil
}
