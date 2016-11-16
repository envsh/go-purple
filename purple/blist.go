package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type Buddy struct {
	buddy *C.PurpleBuddy
}

func newBuddyFrom(buddy *C.PurpleBuddy) *Buddy {
	this := &Buddy{buddy}
	return this
}

func NewBuddy(a *Account, name string, alias string) *Buddy {
	buddy := C.purple_buddy_new(a.account, C.CString(name), C.CString(alias))
	this := &Buddy{}
	this.buddy = buddy
	return this
}

func (this *Buddy) SetAlias(alias string) {
	C.purple_blist_alias_buddy(this.buddy, C.CString(alias))
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

func (this *Buddy) GetName() string {
	cstr := C.purple_buddy_get_name(this.buddy)
	return C.GoString(cstr)
}

func (this *Buddy) GetAliasOnly() string {
	cstr := C.purple_buddy_get_alias_only(this.buddy)
	return C.GoString(cstr)
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

func newGroupFrom(group *C.PurpleGroup) *Group {
	this := &Group{}
	this.group = group
	return this
}
