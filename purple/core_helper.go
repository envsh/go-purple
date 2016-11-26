package purple

/*
// core.c encapse libpurple's core init
#include <libpurple/purple.h>

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

type PurpleCore struct {
	accountUiOps C.PurpleAccountUiOps
	connUiOps    C.PurpleConnectionUiOps
	convUiOps    C.PurpleConversationUiOps

	loopUiOps C.PurpleEventLoopUiOps
	coreUiOps C.PurpleCoreUiOps

	loop *C.GMainLoop
}

func NewPurpleCore() *PurpleCore {
	this := &PurpleCore{}
	this.setupCore()
	return this
}

func (this *PurpleCore) initUiOps() {
	C.purple_connections_set_ui_ops(&this.connUiOps)
	C.purple_accounts_set_ui_ops(&this.accountUiOps)
	C.purple_conversations_set_ui_ops(&this.convUiOps)

	if false {
		this.accountUiOps.request_authorize = nil
	}
}

func (this *PurpleCore) initCoreOps() {
	C.purple_core_set_ui_ops(&this.coreUiOps)
	C.purple_eventloop_set_ui_ops(&this.loopUiOps)
}

func (this *PurpleCore) initLibpurple() {
	if false {
		home := fmt.Sprintf("%s/%s", os.Getenv("HOME"), CUSTOM_USER_DIRECTORY)
		C.purple_util_set_user_dir(C.CString(home))
		C.purple_debug_set_enabled(C.FALSE)
	}

	// this.initCoreOps()
	C.purple_core_set_ui_ops(&this.coreUiOps)
	// C.purple_eventloop_set_ui_ops(&this.loopUiOps)
	C.purple_eventloop_set_ui_ops(C.gopurple_get_loopops())

	C.purple_plugins_add_search_path(C.CString(CUSTOM_PLUGIN_PATH))

	if cok := C.purple_core_init(C.CString(UI_ID)); cok != C.TRUE {
		log.Println("libpurple initialization failed. Dumping core." +
			"Please report this!\n")
		os.Exit(-1)
	}

	C.purple_set_blist(C.purple_blist_new())
	C.purple_blist_load()
	C.purple_prefs_load()
	C.purple_plugins_load_saved(C.CString(PLUGIN_SAVE_PREF))
	C.purple_pounces_load()

	//
	C.gopurple_connect_to_signals()
}

// 在这个setup之后，才能够调用purple函数。
func (this *PurpleCore) setupCore() {
	this.loop = C.g_main_loop_new(nil, C.FALSE)

	this.initUiOps()

	this.initLibpurple()
}

func (this *PurpleCore) MainLoop() {

	log.Println("hhhh")
	go func() { C.g_main_loop_run(this.loop) }()
	log.Println("looped")
	select {}
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
func gopurple_request_authorize() { log.Println("hehhe") }

//export gopurple_connect_progress
func gopurple_connect_progress() { log.Println("hehhe") }

//export gopurple_notice
func gopurple_notice() { log.Println("hehhe") }

//export gopurple_network_disconnected
func gopurple_network_disconnected() { log.Println("hehhe") }

//export gopurple_report_disconnect_reason
func gopurple_report_disconnect_reason() { log.Println("hehhe") }

// TODO nice way
var SignedOn func(*Connection)

//export gopurple_signed_on
func gopurple_signed_on(gc *C.PurpleConnection, data unsafe.Pointer) {
	log.Println("hehhe", gc, data)
	if SignedOn != nil {
		SignedOn(newConnectionFrom(gc))
	}
}

//export gopurple_buddy_signed_on
func gopurple_buddy_signed_on() { log.Println("hehhe") }

//export gopurple_buddy_signed_off
func gopurple_buddy_signed_off() { log.Println("hehhe") }

//export gopurple_buddy_away
func gopurple_buddy_away() { log.Println("hehhe") }

//export gopurple_buddy_idle
func gopurple_buddy_idle() { log.Println("hehhe") }

//export gopurple_received_im_msg
func gopurple_received_im_msg() { log.Println("hehhe") }

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

//export gopurple_received_chat_msg
func gopurple_received_chat_msg() { log.Println("hehhe") }
