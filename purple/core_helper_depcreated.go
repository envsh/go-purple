package purple

/*
#include <libpurple/purple.h>
extern void init_libpurple(void);
typedef GHashTable *(*chat_info_defaults_func)(PurpleConnection *, const char *chat_name);

//
#include "misc.h"
*/
import "C"

import (
	"log"
)

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
