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

// cgo don't support variadic parameter list, so wrapper it
func PrplGotUserStatus(ac *Account, name, statusId string) {
	C.gopurple_prpl_got_user_status(ac.account, CCString(name).Ptr, CCString(statusId).Ptr)
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
