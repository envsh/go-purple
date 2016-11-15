package purple

/*
// core.c encapse libpurple's core init

#include <libpurple/purple.h>
extern void ui_init(void);
extern void init_libpurple(void);
extern void connect_to_signalscc(void*);
typedef GHashTable *(*chat_info_defaults_func)(PurpleConnection *, const char *chat_name);

//
#include "misc.h"
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
	loopUiOps C.PurpleEventLoopUiOps

	accountUiOps C.PurpleAccountUiOps

	connUiOps C.PurpleConnectionUiOps

	convUiOps C.PurpleConversationUiOps

	coreUiOps C.PurpleCoreUiOps

	loop *C.GMainLoop
}

func NewPurpleCore() *PurpleCore {
	this := &PurpleCore{}
	this.loop = C.g_main_loop_new(nil, C.FALSE)
	return this
}

func (this *PurpleCore) InitUi() {
	C.ui_init()

	if false {
		this.accountUiOps.request_authorize = nil
	}
}

func (this *PurpleCore) InitPurple() {

	if true {
		C.init_libpurple()
		if purple_debug {
			C.purple_debug_set_enabled(C.TRUE)
		} else {
			C.purple_debug_set_enabled(C.FALSE)
		}
		pc := unsafe.Pointer(this)
		C.connect_to_signalscc(pc)
	}

	if false {
		home := fmt.Sprintf("%s/%s", os.Getenv("HOME"), CUSTOM_USER_DIRECTORY)
		C.purple_util_set_user_dir(C.CString(home))
		C.purple_debug_set_enabled(C.FALSE)
		C.purple_core_set_ui_ops(&this.coreUiOps)
		C.purple_eventloop_set_ui_ops(&this.loopUiOps)
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
	}
}

func (this *PurpleCore) Loop() {
	go func() { C.g_main_loop_run(this.loop) }()
	log.Println("looped")
	select {}
}

// get all accout
func (this *PurpleCore) AvaliableAccounts() []*Account {

	return nil
}

func (this *PurpleCore) AccountsGetAll() []*Account {
	acs := make([]*Account, 0)
	lst := C.purple_accounts_get_all()
	newGListFrom(lst).Each(func(item C.gpointer) {
		ac := newAccountWrapper((*C.PurpleAccount)(item))
		acs = append(acs, ac)
	})
	return acs
}

func (this *PurpleCore) AccountsFind(account string, protocol string) *Account {
	acc := C.purple_accounts_find(C.CString(account), C.CString(protocol))
	if acc == nil {
		log.Println("not found", account, protocol)
	} else {
		log.Println(acc.username, acc.user_info)
		return newAccountWrapper(acc)
	}

	return nil
}

// send with conv
func (this *PurpleCore) SendGroupMessage(acc *Account, room string, msg string) {

}

// send with conv
func (this *PurpleCore) sendGroupMessage1(acc *Account, room string, msg string) {
	gc := C.purple_account_get_connection(acc.account)
	hash := this.protoInfo(gc, C.CString(room))
	C.serv_join_chat(gc, hash)

	conv := C.purple_conversation_new(C.PURPLE_CONV_TYPE_CHAT, acc.account, C.CString(room))
	if conv != nil {
	}
	C.purple_conversation_present(conv) // important
	chat := C.purple_conversation_get_chat_data(conv)
	C.purple_conv_chat_send(chat, C.CString(msg))

	if true {
		chatid := this.convChatId(conv)
		log.Println(chatid)
		C.serv_chat_send(gc, chatid, C.CString(msg+" <- raw send"), 0)
	}
}

// send with blist
func (this *PurpleCore) sendGroupMessage2(acc *Account, room string, msg string) {
	gc := C.purple_account_get_connection(acc.account)
	hash := this.protoInfo(gc, C.CString(room))

	bchat := C.purple_chat_new(acc.account, C.CString(room), hash)
	C.serv_join_chat(gc, C.purple_chat_get_components(bchat))

	var chatid C.int = 0

	if true {
		// wtf bug
		conv := C.purple_find_conversation_with_account(C.PURPLE_CONV_TYPE_CHAT, C.CString(room), acc.account)
		if conv == nil {
			log.Println("not found")
		} else {
			chat := C.purple_conversation_get_chat_data(conv)
			chatid = C.purple_conv_chat_get_id(chat)
		}
		log.Println(chatid)
	}

	C.serv_chat_send(gc, chatid, C.CString(msg), 0)
}

func (this *PurpleCore) convChatId(conv *C.PurpleConversation) C.int {
	return C.purple_conv_chat_get_id(C.purple_conversation_get_chat_data(conv))
}

func (this *PurpleCore) protoInfo(conn *C.PurpleConnection, name *C.char) *C.GHashTable {
	prpl := C.purple_connection_get_prpl(conn)
	info := (*C.PurplePluginProtocolInfo)(prpl.info.extra_info)
	if info.chat_info_defaults != nil {
		fnptr := (C.chat_info_defaults_func)(info.chat_info_defaults)
		if fnptr != nil {
		}
		// cannot call non-function fnptr (type C.chat_info_defaults_func)
		// fnptr(conn, C.CString("aaa"))
	}

	hash := C.gopurple_connection_get_chat_info_defaults(conn, name)
	return hash
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

//export gopurple_signed_on
func gopurple_signed_on() {
	log.Println("hehhe")
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

//
func F_core_init(ui string) bool {
	C.init_libpurple()
	var pret = C.purple_accounts_find(C.CString("yournicknameu123@irc.freenode.net"),
		C.CString("prpl-irc"))
	log.Println(pret)
	if pret != nil {
		log.Println(C.GoString(pret.username), C.GoString(pret.alias))
	}

	/*
		var bret = C.purple_core_init(C.CString(ui))
		if bret == C.TRUE {
			return true
		}
	*/

	return false
}
