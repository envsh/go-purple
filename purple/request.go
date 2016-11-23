package purple

/*
#include <libpurple/purple.h>

// fix _ not defined and hidden definition
#include <libintl.h>
#define _(str) gettext((str))

extern void gopurple_request_action_cb(char *user_data, int choice);

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

static void gopurple_request_accept_cancel(void *handle, char *title, char *primary, char *secondary,
                                    int default_action, PurpleAccount *ac, char *who,
                                    PurpleConversation *conv, char *user_data)
{
     purple_request_accept_cancel(handle, title, primary, secondary, default_action,
		 ac, who, conv, user_data, gopurple_request_action_cb, gopurple_request_action_cb);
}

*/
import "C"
import "unsafe"

import (
	"log"
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

func RequestYesNo(userData interface{}, gc *Connection, title, primary string,
	yescb func(interface{}), nocb func(interface{})) {

	cseq := (*C.char)(C.calloc(1, 1))
	req := &Request{seq: cseq, userData: userData, yescb: yescb, nocb: nocb}
	requests[cseq] = req

	handle := unsafe.Pointer(gc.conn)
	ac := gc.ConnGetAccount()
	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	who := "who's"
	C.gopurple_request_yes_no(handle, C.CString(title), C.CString(primary), nil,
		C.int(default_action), ac.account, C.CString(who), nil, cseq)
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

	handle := unsafe.Pointer(gc.conn)
	ac := gc.ConnGetAccount()
	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	who := "who's"
	C.gopurple_request_ok_cancel(handle, C.CString(title), C.CString(primary), nil,
		C.int(default_action), ac.account, C.CString(who), nil, cseq)
}

func RequestOkCancelDemo(userData interface{}, gc *Connection,
	yescb func(interface{}), nocb func(interface{})) {
	title := "it's title"
	primary := "it's primary"
	RequestOkCancel(userData, gc, title, primary, yescb, nocb)
}

func RequestAcceptCancel(userData interface{}, gc *Connection, title, primary string,
	yescb func(interface{}), nocb func(interface{})) {

	cseq := (*C.char)(C.calloc(1, 1))
	req := &Request{seq: cseq, userData: userData, yescb: yescb, nocb: nocb}
	requests[cseq] = req

	handle := unsafe.Pointer(gc.conn)
	ac := gc.ConnGetAccount()
	// title := "it's title"
	// primary := "it's primary"
	// secondary := "it's secondary"
	default_action := 0
	who := "who's"
	C.gopurple_request_accept_cancel(handle, C.CString(title), C.CString(primary), nil,
		C.int(default_action), ac.account, C.CString(who), nil, cseq)
}
func RequestAcceptCancelDemo(userData interface{}, gc *Connection,
	yescb func(interface{}), nocb func(interface{})) {
	title := "it's title"
	primary := "it's primary"
	RequestAcceptCancel(userData, gc, title, primary, yescb, nocb)
}

//export gopurple_request_action_cb
func gopurple_request_action_cb(userData *C.char, choice C.int) {
	log.Println(userData, choice)

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
		log.Println("wtf", userData, uint64(*userData))
	}

	// clean up
	if _, ok := requests[userData]; ok {
		delete(requests, userData)
		C.free(unsafe.Pointer(userData))
	}
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
