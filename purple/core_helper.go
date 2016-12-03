package purple

/*
// core.c encapse libpurple's core init
#include <libpurple/purple.h>
#include <glib.h>

#include "core_helper.h"
extern PurpleEventLoopUiOps *gopurple_get_loopops();
extern void gopurple_connect_to_signals(void);
*/
import "C"

import "unsafe"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/kitech/colog"
)

var purple_debug = false

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()

	log.SetFlags(log.Flags() | log.Lshortfile)
	colog.Register()

	flag.BoolVar(&purple_debug, "purple-debug", purple_debug, "enable purple debug.")
}

type CoreCallbacks struct {
	ConnectProgress        func(gc *Connection, text string, step int, stepCount int)
	Connected              func(gc *Connection)
	Disconnected           func(gc *Connection)
	Notice                 func(gc *Connection, text string)
	ReportDisconnect       func(gc *Connection, text string)
	NetworkConnected       func()
	NetworkDisconnected    func()
	ReportDisconnectReason func(gc *Connection, reason int, text string)
}

var ccbs CoreCallbacks

type ConnSignals struct {
	SignedOn  func(gc *Connection /*, data unsafe.Pointer*/)
	SignedOff func(gc *Connection)
}
type BlistSignals struct {
	BuddySignedOn  func(buddy *Buddy)
	BuddySignedOff func(buddy *Buddy)
}
type ConvSignals struct {
	ReceivedIMMsg   func(ac *Account, sender, message string, conv *Conversation, flags int)
	ReceivedChatMsg func(ac *Account, sender, message string, conv *Conversation, flags int)
}
type AccountSignals struct{}

type CoreSignals struct {
	ConnSignals
	BlistSignals
	ConvSignals
	AccountSignals
}

var csigs CoreSignals

type PurpleCore struct {
	accountUiOps C.PurpleAccountUiOps
	connUiOps    C.PurpleConnectionUiOps
	convUiOps    C.PurpleConversationUiOps

	loopUiOps C.PurpleEventLoopUiOps
	coreUiOps C.PurpleCoreUiOps

	loop  *C.GMainLoop
	ccbs  CoreCallbacks
	csigs CoreSignals
}

func NewPurpleCore() *PurpleCore {
	this := &PurpleCore{}
	this.setupCore()
	return this
}

func (this *PurpleCore) SetCallbacks(cbs CoreCallbacks) {
	this.ccbs = cbs
	ccbs = cbs
}

func (this *PurpleCore) SetSignals(sigs CoreSignals) {
	this.csigs = sigs
	csigs = sigs
}

func (this *PurpleCore) initUiOps() {

	if true {
		// wow, done
		this.connUiOps.connect_progress = (*[0]byte)((unsafe.Pointer)(C.gopurple_connect_progress))
		this.connUiOps.connected = (*[0]byte)((unsafe.Pointer)(C.gopurple_connected))
		this.connUiOps.disconnected = (*[0]byte)((unsafe.Pointer)(C.gopurple_disconnected))
		this.connUiOps.notice = (*[0]byte)((unsafe.Pointer)(C.gopurple_notice))
		this.connUiOps.report_disconnect = (*[0]byte)((unsafe.Pointer)(C.gopurple_report_disconnect))
		this.connUiOps.network_connected = (*[0]byte)((unsafe.Pointer)(C.gopurple_network_connected))
		this.connUiOps.network_disconnected = (*[0]byte)((unsafe.Pointer)(C.gopurple_network_disconnected))
		this.connUiOps.report_disconnect_reason = (*[0]byte)((unsafe.Pointer)(C.gopurple_report_disconnect_reason))
		// log.Printf("%+v\n", this.connUiOps)
	}

	if true {
		this.accountUiOps.request_authorize = go2cfn(C.gopurple_request_authorize)
	}

	C.purple_connections_set_ui_ops(&this.connUiOps)
	C.purple_accounts_set_ui_ops(&this.accountUiOps)
	C.purple_conversations_set_ui_ops(&this.convUiOps)

	if false {
		this.accountUiOps.request_authorize = nil
	}
}

func (this *PurpleCore) initCoreOps() {
	if true {
		this.loopUiOps.timeout_add = go2cfn(C.g_timeout_add)
		this.loopUiOps.timeout_remove = go2cfn(C.g_source_remove)
		this.loopUiOps.input_add = nil
		this.loopUiOps.input_remove = go2cfn(C.g_source_remove)
		this.loopUiOps.input_get_error = nil
		this.loopUiOps.timeout_add_seconds = go2cfn(C.g_timeout_add_seconds)
	}
	C.purple_core_set_ui_ops(&this.coreUiOps)
	C.purple_eventloop_set_ui_ops(&this.loopUiOps)
}

func (this *PurpleCore) initLibpurple() {
	if false {
		home := fmt.Sprintf("%s/%s", os.Getenv("HOME"), CUSTOM_USER_DIRECTORY)
		C.purple_util_set_user_dir(CCString(home).Ptr)
		C.purple_debug_set_enabled(C.FALSE)
	}
	C.purple_debug_set_enabled(go2cBool(purple_debug))

	// this.initCoreOps()
	C.purple_core_set_ui_ops(&this.coreUiOps)
	// C.purple_eventloop_set_ui_ops(&this.loopUiOps)
	C.purple_eventloop_set_ui_ops(C.gopurple_get_loopops())

	C.purple_plugins_add_search_path(CCString(CUSTOM_PLUGIN_PATH).Ptr)

	if cok := C.purple_core_init(CCString(UI_ID).Ptr); cok != C.TRUE {
		log.Println("libpurple initialization failed. Dumping core." +
			"Please report this!\n")
		os.Exit(-1)
	}

	C.purple_set_blist(C.purple_blist_new())
	C.purple_blist_load()
	C.purple_prefs_load()
	C.purple_plugins_load_saved(CCString(PLUGIN_SAVE_PREF).Ptr)
	C.purple_pounces_load()

	// 防止出现信号叠加，多次回调通知
	if false {
		C.gopurple_connect_to_signals() // c version
	} else {
		this.connect_to_signals() // go version
	}
}

func (this *PurpleCore) connect_to_signals() {
	if true {
		// signals
		signalConnect(C.purple_connections_get_handle(), "signed-on",
			(unsafe.Pointer)(C.gopurple_signed_on))
		signalConnect(C.purple_connections_get_handle(), "signed-off",
			(unsafe.Pointer)(C.gopurple_signed_off))
		signalConnect(C.purple_blist_get_handle(), "buddy-signed-on",
			(unsafe.Pointer)(C.gopurple_buddy_signed_on))
		signalConnect(C.purple_blist_get_handle(), "buddy-signed-off",
			(unsafe.Pointer)(C.gopurple_buddy_signed_off))
		signalConnect(C.purple_conversations_get_handle(), "received-im-msg",
			(unsafe.Pointer)(C.gopurple_received_im_msg))
		signalConnect(C.purple_conversations_get_handle(), "received-chat-msg",
			(unsafe.Pointer)(C.gopurple_received_chat_msg))
	}
}

// 在这个setup之后，才能够调用purple函数。
func (this *PurpleCore) setupCore() {
	this.loop = C.g_main_loop_new(nil, C.FALSE)

	this.initUiOps()

	this.initLibpurple()
}

func (this *PurpleCore) MainLoop() {
	C.g_main_loop_run(this.loop)
	// go func() { C.g_main_loop_run(this.loop) }()
	// select {}
}

// get all accout
func (this *PurpleCore) AccountsGetAll() []*Account {
	return AccountsGetAll()
}

func (this *PurpleCore) AccountsFind(name string, protocol string) *Account {
	return AccountsFind(name, protocol)
}

// callbacks
//export gopurple_request_authorize
func gopurple_request_authorize(ac *C.PurpleAccount, remote_user *C.char,
	id *C.char, alias *C.char, message *C.char, on_list C.gboolean,
	authorize_cb C.PurpleAccountRequestAuthorizationCb,
	deny_cb C.PurpleAccountRequestAuthorizationCb, user_data unsafe.Pointer) unsafe.Pointer {
	log.Println("hehhe")
	return nil
}

//export gopurple_connect_progress
func gopurple_connect_progress(gc *C.PurpleConnection, text *C.char,
	step C.size_t, step_count C.size_t) {
	log.Println("hehhe", C.GoString(text), step, step_count)

	if ccbs.ConnectProgress != nil {
		ccbs.ConnectProgress(newConnectionFrom(gc), C.GoString(text), int(step), int(step_count))
	}
}

//export gopurple_connected
func gopurple_connected(gc *C.PurpleConnection) {
	log.Println("hehhe")
	if ccbs.Connected != nil {
		ccbs.Connected(newConnectionFrom(gc))
	}
}

//export gopurple_disconnected
func gopurple_disconnected(gc *C.PurpleConnection) {
	log.Println("hehhe")
	if ccbs.Disconnected != nil {
		ccbs.Disconnected(newConnectionFrom(gc))
	}
}

//export gopurple_notice
func gopurple_notice(gc *C.PurpleConnection, text *C.char) {
	log.Println("hehhe")
	if ccbs.Notice != nil {
		ccbs.Notice(newConnectionFrom(gc), C.GoString(text))
	}
}

//export gopurple_report_disconnect
func gopurple_report_disconnect(gc *C.PurpleConnection, text *C.char) {
	log.Println("hehhe", C.GoString(text))
	if ccbs.ReportDisconnect != nil {
		ccbs.ReportDisconnect(newConnectionFrom(gc), C.GoString(text))
	}
}

//export gopurple_network_connected
func gopurple_network_connected() {
	log.Println("hehhe")
	if ccbs.NetworkConnected != nil {
		ccbs.NetworkConnected()
	}
}

//export gopurple_network_disconnected
func gopurple_network_disconnected() {
	log.Println("hehhe")
	if ccbs.NetworkDisconnected != nil {
		ccbs.NetworkDisconnected()
	}
}

//export gopurple_report_disconnect_reason
func gopurple_report_disconnect_reason(gc *C.PurpleConnection,
	reason C.PurpleConnectionError, text *C.char) {
	log.Println("hehhe", reason, C.GoString(text))
	if ccbs.ReportDisconnectReason != nil {
		ccbs.ReportDisconnectReason(newConnectionFrom(gc), int(reason), C.GoString(text))
	}
}

// signals
//export gopurple_signed_on
func gopurple_signed_on(gc *C.PurpleConnection, data unsafe.Pointer) {
	log.Println("hehhe", gc, data)
	if csigs.SignedOn != nil {
		csigs.SignedOn(newConnectionFrom(gc))
	}
}

//export gopurple_signed_off
func gopurple_signed_off(gc *C.PurpleConnection, data unsafe.Pointer) {
	log.Println("hehhe", gc, data)
	if csigs.SignedOff != nil {
		csigs.SignedOff(newConnectionFrom(gc))
	}
}

//export gopurple_buddy_signed_on
func gopurple_buddy_signed_on(buddy *C.PurpleBuddy) {
	log.Println("hehhe")
	if csigs.BuddySignedOn != nil {
		csigs.BuddySignedOn(newBuddyFrom(buddy))
	}
}

//export gopurple_buddy_signed_off
func gopurple_buddy_signed_off(buddy *C.PurpleBuddy) {
	log.Println("hehhe")
	if csigs.BuddySignedOff != nil {
		csigs.BuddySignedOff(newBuddyFrom(buddy))
	}
}

//export gopurple_buddy_away
func gopurple_buddy_away() { log.Println("hehhe") }

//export gopurple_buddy_idle
func gopurple_buddy_idle() { log.Println("hehhe") }

//export gopurple_buddy_typing
func gopurple_buddy_typing() { log.Println("hehhe") }

//export gopurple_buddy_typed
func gopurple_buddy_typed() { log.Println("hehhe") }

//export gopurple_buddy_typing_stopped
func gopurple_buddy_typing_stopped() { log.Println("hehhe") }

//export gopurple_account_authorization_requested
func gopurple_account_authorization_requested() { log.Println("hehhe") }

//export gopurple_dbus_method_called
func gopurple_dbus_method_called() { log.Println("hehhe") }

//export gopurple_received_im_msg
func gopurple_received_im_msg(ac *C.PurpleAccount, sender *C.char, buffer *C.char,
	chat *C.PurpleConversation, flags C.PurpleMessageFlags) {
	log.Println("hehhe")
	if csigs.ReceivedIMMsg != nil {
		csigs.ReceivedIMMsg(newAccountFrom(ac), C.GoString(sender),
			C.GoString(buffer), newConversationFrom(chat), int(flags))
	}
}

//export gopurple_received_chat_msg
func gopurple_received_chat_msg(ac *C.PurpleAccount, sender *C.char, buffer *C.char,
	chat *C.PurpleConversation, flags C.PurpleMessageFlags, data unsafe.Pointer) {
	log.Println("hehhe")
	if csigs.ReceivedChatMsg != nil {
		csigs.ReceivedChatMsg(newAccountFrom(ac), C.GoString(sender),
			C.GoString(buffer), newConversationFrom(chat), int(flags))
	}
}
