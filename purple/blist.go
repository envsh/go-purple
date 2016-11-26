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

func NewBuddy(ac *Account, name string, alias string) *Buddy {
	buddy := C.purple_buddy_new(ac.account, C.CString(name), C.CString(alias))
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

func BlistGetBuddies() []*Buddy {
	lst := C.purple_blist_get_buddies()
	newGSListFrom(lst).Each(func(item C.gpointer) {
	})
	return nil
}

type Chat struct {
	// private
	chat *C.PurpleChat
}

func newChatFrom(chat *C.PurpleChat) *Chat {
	this := &Chat{}
	this.chat = chat
	return this
}

func (this *Chat) GetComponents() *GHashTable {
	ht := C.purple_chat_get_components(this.chat)
	if ht == nil {
		return nil
	}
	return newGHashTableFrom(ht)
}

func (this *Account) BlistFindChat(name string) *Chat {
	chat := C.purple_blist_find_chat(this.account, C.CString(name))
	if chat == nil {
		return nil
	}
	return newChatFrom(chat)
}

func (this *Account) ChatNew(alias string, components *GHashTable) *Chat {
	chat := C.purple_chat_new(this.account, C.CString(alias), components.ht)
	return newChatFrom(chat)
}

func (this *Chat) Destory() {
	C.purple_chat_destroy(this.chat)
	this.chat = nil
}

func (this *Chat) Node() *BlistNode {
	node := &this.chat.node
	return newBlistNodeFrom(node)
}

func (this *Chat) BlistAddChat(g *Group, node *BlistNode) {
	C.purple_blist_add_chat(this.chat, nil, nil)
}

type BlistNode struct {
	node *C.PurpleBlistNode
}

func newBlistNodeFrom(node *C.PurpleBlistNode) *BlistNode {
	this := &BlistNode{}
	this.node = node
	return this
}

func (this *BlistNode) GetBool(key string) bool {
	rc := C.purple_blist_node_get_bool(this.node, C.CString(key))
	if rc == C.TRUE {
		return true
	}
	return false
}
func (this *BlistNode) SetBool(key string, value bool) {
	C.purple_blist_node_set_bool(this.node, C.CString(key), go2cBool(value))
}

func (this *BlistNode) GetString(key string) string {
	s := C.purple_blist_node_get_string(this.node, C.CString(key))
	return C.GoString(s)
}
func (this *BlistNode) SetString(key, value string) {
	C.purple_blist_node_set_string(this.node, C.CString(key), C.CString(value))
}

func (this *BlistNode) Settings() *GHashTable {
	return newGHashTableFrom(this.node.settings)
}

func (this *BlistNode) SetFlags(flags int) {
	C.purple_blist_node_set_flags(this.node, C.PurpleBlistNodeFlags(flags))
}
func (this *BlistNode) GetFlags() (flags int) {
	cflags := C.purple_blist_node_get_flags(this.node)
	flags = int(cflags)
	return
}

type BuddyList struct {
	// private
	blst *C.PurpleBuddyList
}

func newBuddyListFrom(blst *C.PurpleBuddyList) *BuddyList {
	this := &BuddyList{}
	this.blst = blst
	return this
}

func BlistNew() *BuddyList {
	blst := C.purple_blist_new()
	return newBuddyListFrom(blst)
}

func BlistSet(blst *BuddyList) {
	C.purple_set_blist(blst.blst)
}

func BlistInit() {
	C.purple_blist_init()
}

func BlistGet() *BuddyList {
	blst := C.purple_get_blist()
	return newBuddyListFrom(blst)
}

func BlistLoad() {
	C.purple_blist_load()
}

func BlistScheduleSave() {
	C.purple_blist_schedule_save()
}
