package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type Buddy struct {
	buddy *C.PurpleBuddy
}

func newBuddyWrapper(buddy *C.PurpleBuddy) *Buddy {
	this := &Buddy{buddy}
	return this
}

func NewBuddy(a *Account, name string, alias string) *Buddy {
	buddy := C.purple_buddy_new(a.account, C.CString(name), C.CString(alias))
	this := &Buddy{}
	this.buddy = buddy
	return this
}

// @param g can be nil
func (this *Buddy) BlistAdd(g *Group) {
	if g == nil {
		C.purple_blist_add_buddy(this.buddy, nil, nil, nil)
	} else {
		C.purple_blist_add_buddy(this.buddy, nil, g.group, nil)
	}
}

func (this *Buddy) SetProtocolData(data interface{}) {
	C.purple_buddy_set_protocol_data(this.buddy, C.CString("hehhehee"))
}

type Group struct {
	group *C.PurpleGroup
}

func NewGroup(name string) *Group {
	group := C.purple_group_new(C.CString(name))
	this := &Group{}
	this.group = group
	return this
}
