package purple

/*
#include <libpurple/purple.h>

static void gopurple_prpl_got_user_status(PurpleAccount *account,
     const char *name, const char *status_id)
{ purple_prpl_got_user_status(account, name, status_id, NULL); }

*/
import "C"
import "unsafe"

import (
	"fmt"
)

type IconScaleRules int

const (
	/**< We scale the icon when we display it */
	ICON_SCALE_DISPLAY = IconScaleRules(C.PURPLE_ICON_SCALE_DISPLAY)
	/**< We scale the icon before we send it to the server */
	ICON_SCALE_SEND = IconScaleRules(C.PURPLE_ICON_SCALE_SEND)
)

type BuddyIconSpec struct {
	Format      string
	MinWidth    int            /**< Minimum width of this icon  */
	MinHeight   int            /**< Minimum height of this icon */
	MaxWidth    int            /**< Maximum width of this icon  */
	MaxHeight   int            /**< Maximum height of this icon */
	MaxFilesize int            /**< Maximum size in bytes */
	ScaleRules  IconScaleRules /**< How to stretch this icon */

	// private
	fromc bool
	bis   *C.PurpleBuddyIconSpec
}

type ProtocolOptions int

const (
	OPT_PROTO_UNIQUE_CHATNAME       = ProtocolOptions(C.OPT_PROTO_UNIQUE_CHATNAME)
	OPT_PROTO_CHAT_TOPIC            = ProtocolOptions(C.OPT_PROTO_CHAT_TOPIC)
	OPT_PROTO_NO_PASSWORD           = ProtocolOptions(C.OPT_PROTO_NO_PASSWORD)
	OPT_PROTO_MAIL_CHECK            = ProtocolOptions(C.OPT_PROTO_MAIL_CHECK)
	OPT_PROTO_IM_IMAGE              = ProtocolOptions(C.OPT_PROTO_IM_IMAGE)
	OPT_PROTO_PASSWORD_OPTIONAL     = ProtocolOptions(C.OPT_PROTO_PASSWORD_OPTIONAL)
	OPT_PROTO_USE_POINTSIZE         = ProtocolOptions(C.OPT_PROTO_USE_POINTSIZE)
	OPT_PROTO_REGISTER_NOSCREENNAME = ProtocolOptions(C.OPT_PROTO_REGISTER_NOSCREENNAME)
	OPT_PROTO_SLASH_COMMANDS_NATIVE = ProtocolOptions(C.OPT_PROTO_SLASH_COMMANDS_NATIVE)
	OPT_PROTO_INVITE_MESSAGE        = ProtocolOptions(C.OPT_PROTO_INVITE_MESSAGE)
)

// cgo don't support variadic parameter list, so wrapper it
func PrplGotUserStatus(ac *Account, name, statusId string) {
	C.gopurple_prpl_got_user_status(ac.account, CCString(name).Ptr, CCString(statusId).Ptr)
}
func (this *Account) GotUserStatus(name, statusId string) {
	PrplGotUserStatus(this, name, statusId)
}

type ProtoChatEntry struct {
	pce *C.struct_proto_chat_entry
}

func NewProtoChatEntry(label, identifier string, required bool) *ProtoChatEntry {
	this := &ProtoChatEntry{}

	var pce *C.struct_proto_chat_entry
	// pce = new(C.struct_proto_chat_entry)  // crash!!!
	pce = (*C.struct_proto_chat_entry)(C.calloc(C.size_t(1), C.sizeof_struct_proto_chat_entry))
	pce.label = C.CString(fmt.Sprintf("_%s", label))
	pce.identifier = C.CString(identifier)
	pce.required = C.FALSE
	if required {
		pce.required = C.TRUE
	}

	this.pce = pce
	return this
}

func (this *ProtoChatEntry) get() *C.struct_proto_chat_entry {
	return this.pce
}

func (this *ProtoChatEntry) Destroy() {
	C.free(unsafe.Pointer(this.pce))
}
