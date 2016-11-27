package purple

/*
#include <libpurple/purple.h>

*/
import "C"

const (
	STATUS_UNSET          = int(C.PURPLE_STATUS_UNSET)
	STATUS_OFFLINE        = int(C.PURPLE_STATUS_OFFLINE)
	STATUS_AVAILABLE      = int(C.PURPLE_STATUS_AVAILABLE)
	STATUS_UNAVAILABLE    = int(C.PURPLE_STATUS_UNAVAILABLE)
	STATUS_INVISIBLE      = int(C.PURPLE_STATUS_INVISIBLE)
	STATUS_AWAY           = int(C.PURPLE_STATUS_AWAY)
	STATUS_EXTENDED_AWAY  = int(C.PURPLE_STATUS_EXTENDED_AWAY)
	STATUS_MOBILE         = int(C.PURPLE_STATUS_MOBILE)
	STATUS_TUNE           = int(C.PURPLE_STATUS_TUNE)
	STATUS_MOOD           = int(C.PURPLE_STATUS_MOOD)
	STATUS_NUM_PRIMITIVES = int(C.PURPLE_STATUS_NUM_PRIMITIVES)
)

type StatusType struct {
	id   string
	name string

	// private
	sty *C.PurpleStatusType
}

func newStatusTypeFrom(sty *C.PurpleStatusType) *StatusType {
	this := &StatusType{}
	this.sty = sty
	return this
}

func NewStatusTypeFull(primitive int, id, name string, saveable bool,
	settable bool, independent bool) *StatusType {
	this := &StatusType{}
	var csaveable C.gboolean = C.FALSE
	if saveable {
		csaveable = C.TRUE
	}
	var csettable C.gboolean = C.FALSE
	if settable {
		csettable = C.TRUE
	}
	var cindependent C.gboolean = C.FALSE
	if independent {
		cindependent = C.TRUE
	}

	this.sty = C.purple_status_type_new_full(C.PurpleStatusPrimitive(primitive),
		CCString(id).Ptr, CCString(name).Ptr, csaveable, csettable, cindependent)
	return this
}

func NewStatusType(primitive int, id, name string, settable bool) *StatusType {
	this := &StatusType{}
	var csettable C.gboolean = C.FALSE
	if settable {
		csettable = C.TRUE
	}

	this.sty = C.purple_status_type_new(C.PurpleStatusPrimitive(primitive),
		CCString(id).Ptr, CCString(name).Ptr, csettable)
	return this
}

func (this *StatusType) Destroy() {
	C.purple_status_type_destroy(this.sty)
}
