package purple

/*
#include <libpurple/purple.h>

extern void goprpl_notify_user_info_close_cb(gpointer userData);

static void gopurple_notify_userinfo(PurpleConnection *gc, char *who,
                                     PurpleNotifyUserInfo *info, gpointer userData) {
    purple_notify_userinfo(gc, who, info, goprpl_notify_user_info_close_cb, userData);
}

*/
import "C"
import "unsafe"

import (
	"log"
)

type NotifyUserInfoEntry struct {
	// private
	nuie *C.PurpleNotifyUserInfoEntry
}

func newNotifyUserInfoEntryFrom(nuie *C.PurpleNotifyUserInfoEntry) *NotifyUserInfoEntry {
	this := &NotifyUserInfoEntry{}
	this.nuie = nuie
	return this
}

type NotifyUserInfo struct {
	// private
	nui *C.PurpleNotifyUserInfo
}

type notifyUserInfoCallbackData struct {
	closecb  func(interface{})
	userData interface{}
}

var notifyUserInfoCallbackDatas = make(map[*C.char]*notifyUserInfoCallbackData, 0)

func (this *Connection) NotifyUserInfo(who string, nui *NotifyUserInfo, closecb func(interface{}), userData interface{}) {
	// C.purple_notify_userinfo(this.conn, C.CString(who), nui.nui, nil, nil)
	cuserData := C.calloc(1, 1)
	notifyUserInfoCallbackDatas[(*C.char)(cuserData)] = &notifyUserInfoCallbackData{closecb, userData}
	C.gopurple_notify_userinfo(this.conn, C.CString(who), nui.nui, cuserData)
}

func NewNotifyUserInfo() *NotifyUserInfo {
	this := &NotifyUserInfo{}
	this.nui = C.purple_notify_user_info_new()
	return this
}

func (this *NotifyUserInfo) Destroy() {
	C.purple_notify_user_info_destroy(this.nui)
}

func (this *NotifyUserInfo) AddPair(label, value string) {
	C.purple_notify_user_info_add_pair(this.nui, C.CString(label), C.CString(value))
}

func (this *NotifyUserInfo) AddPairPlaintext(label, value string) {
	C.purple_notify_user_info_add_pair_plaintext(this.nui, C.CString(label), C.CString(value))
}

func (this *NotifyUserInfo) PrependPair(label, value string) {
	C.purple_notify_user_info_prepend_pair(this.nui, C.CString(label), C.CString(value))
}

func (this *NotifyUserInfo) RemoveEntry(entry *NotifyUserInfoEntry) {
	C.purple_notify_user_info_remove_entry(this.nui, entry.nuie)
}

//export goprpl_notify_user_info_close_cb
func goprpl_notify_user_info_close_cb(userData C.gpointer) {
	log.Println(userData)
	cbe := notifyUserInfoCallbackDatas[(*C.char)(userData)]
	cbe.closecb(cbe.userData)
	delete(notifyUserInfoCallbackDatas, (*C.char)(userData))
	C.free(unsafe.Pointer(userData))
}
