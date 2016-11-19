/*
helper go wrapper plugin, not include raw plugin.h function
*/
package purple

/*
#include <glib.h>

#include "notify.h"
#include "plugin.h"
#include "version.h"
#include "prpl.h"
#include <string.h>

static void _set_plugin_type(PurplePluginInfo *pi) {
    pi->type = PURPLE_PLUGIN_STANDARD;
    pi->type = PURPLE_PLUGIN_PROTOCOL;
   //  pi->type = PURPLE_PLUGIN_LOADER;
}
extern gboolean goprpl_plugin_load(PurplePlugin *p);
extern gboolean goprpl_plugin_unload(PurplePlugin *p);
extern void goprpl_plugin_destroy(PurplePlugin *p);
extern char* goprpl_blist_icon(PurpleAccount *ac, PurpleBuddy *b);
extern void goprpl_login(PurpleAccount *account);
extern void goprpl_close(PurpleConnection *gc);
extern GList* goprpl_status_types(PurpleAccount *ac);
extern GList *goprpl_chat_info(PurpleConnection *gc);
extern GHashTable *goprpl_chat_info_defaults(PurpleConnection *, char *chat_name);

extern int goprpl_send_im(PurpleConnection *gc, char *who, char *msg, PurpleMessageFlags flags);
// === group chat ===
extern void goprpl_join_chat(PurpleConnection *gc, GHashTable *comp);
extern void goprpl_reject_chat(PurpleConnection *gc, GHashTable *comp);
extern char *goprpl_get_chat_name(GHashTable *comp);
extern void goprpl_chat_invite(PurpleConnection *gc, int id, char *message, char *who);
extern void goprpl_chat_leave(PurpleConnection *gc, int id);
extern void goprpl_chat_whisper(PurpleConnection *gc, int id, char *who, char *message);
extern int  goprpl_chat_send(PurpleConnection *gc, int id, char *message, PurpleMessageFlags flags);
extern PurpleRoomlist *goprpl_roomlist_get_list(PurpleConnection *gc);
extern void goprpl_add_buddy_with_invite(PurpleConnection *gc, PurpleBuddy *buddy, PurpleGroup *group, char *message);
extern void goprpl_remove_buddy(PurpleConnection *gc, PurpleBuddy *buddy, PurpleGroup *group);
extern void goprpl_add_permit(PurpleConnection *gc, char *name);
extern void goprpl_add_deny(PurpleConnection *gc, char *name);
extern void goprpl_rem_permit(PurpleConnection *gc, char *name);
extern void goprpl_rem_deny(PurpleConnection *gc, char *name);
extern void goprpl_get_info(PurpleConnection *gc, char *name);
extern char* goprpl_status_text(PurpleBuddy *buddy);

static void _set_plugin_funcs(PurplePluginInfo *pi, PurplePluginProtocolInfo *pppi) {
    pi->load = goprpl_plugin_load;
    pi->unload = goprpl_plugin_unload;
    pi->destroy = goprpl_plugin_destroy;

    //
    pi->extra_info = pppi;
    // pppi->icon_spec = NO_BUDDY_ICONS;
    PurpleBuddyIconSpec ispec = {0};
    memcpy(&pppi->icon_spec, &ispec, sizeof(ispec));
    pppi->list_icon = (const char*(*)(PurpleAccount*, PurpleBuddy*))goprpl_blist_icon;
    pppi->login = goprpl_login;
    pppi->close = goprpl_close;
    pppi->status_types = goprpl_status_types;
    // optional callbacks
    pppi->chat_info = goprpl_chat_info;
    pppi->chat_info_defaults = goprpl_chat_info_defaults;
    pppi->send_im = (int(*)(PurpleConnection *gc, const char*, const char*, PurpleMessageFlags))goprpl_send_im;
    pppi->join_chat = goprpl_join_chat;
    pppi->reject_chat = goprpl_reject_chat;
    pppi->get_chat_name = goprpl_get_chat_name;
    pppi->chat_invite = goprpl_chat_invite;
    pppi->chat_leave = goprpl_chat_leave;
    pppi->chat_whisper = goprpl_chat_whisper;
    pppi->chat_send = goprpl_chat_send;
    pppi->roomlist_get_list = goprpl_roomlist_get_list;
    pppi->add_buddy_with_invite = goprpl_add_buddy_with_invite;
    pppi->remove_buddy = goprpl_remove_buddy;
    pppi->add_permit = goprpl_add_permit;
    pppi->add_deny = goprpl_add_deny;
    pppi->rem_permit = goprpl_rem_permit;
    pppi->rem_deny = goprpl_rem_deny;
    pppi->get_info = goprpl_get_info;
    pppi->status_text = goprpl_status_text;
    // TODO fix compile warnings
}

// utils
static GHashTable *goprpl_hash_table_new_full() {
    return g_hash_table_new_full(g_str_hash, g_str_equal, NULL, g_free);
}
*/
import "C"
import "unsafe"

import (
	"log"
	"runtime"

	"github.com/kitech/colog"
)

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()

	colog.Register()
	colog.SetFlags(log.LstdFlags | log.Lshortfile | colog.Flags())
}

// the fromc field means wrapper C directly
type PluginInfo struct {
	Id          string
	Name        string
	Version     string
	Summary     string
	Description string
	Author      string
	Homepage    string

	/**
	 * If a plugin defines a 'load' function, and it returns FALSE,
	 * then the plugin will not be loaded.
	 */
	Load    func(plugin *Plugin) bool
	Unload  func(plugin *Plugin) bool
	Destroy func(plugin *Plugin)

	// private
	pi    *C.PurplePluginInfo
	fromc bool
}

type PluginProtocolInfo struct {
	// must
	BlistIcon   func() string
	StatusTypes func(*Account) []*StatusType
	Login       func(*Account)
	Close       func(*Connection)

	// optional, might by Proirity high to low
	ChatInfo           func(gc *Connection) []*ProtoChatEntry
	ChatInfoDefaults   func(gc *Connection, chat_name string) map[string]string
	SendIM             func(gc *Connection, who string, message string) int
	JoinChat           func(gc *Connection, comp *GHashTable)
	RejectChat         func(gc *Connection, comp *GHashTable)
	GetChatName        func(comp *GHashTable) string
	ChatInvite         func(gc *Connection, id int, message string, who string)
	ChatLeave          func(gc *Connection, id int)
	ChatWhisper        func(gc *Connection, id int, who string, message string)
	ChatSend           func(gc *Connection, id int, message string, flags int) int
	RoomlistGetList    func(gc *Connection)
	AddBuddyWithInvite func(gc *Connection, buddy *Buddy, group *Group, message string)
	RemoveBuddy        func(gc *Connection, buddy *Buddy, group *Group)
	AddPermit          func(gc *Connection, name string)
	AddDeny            func(gc *Connection, name string)
	RemPermit          func(gc *Connection, name string)
	RemDeny            func(gc *Connection, name string)
	GetInfo            func(gc *Connection, name string)
	StatusText         func(*Buddy) string

	// private
	ppi   *C.PurplePluginProtocolInfo
	fromc bool
}

type Plugin struct {
	cpp   *C.PurplePlugin
	cppi  *C.PurplePluginInfo
	cpppi *C.PurplePluginProtocolInfo
	pi    *PluginInfo
	ppi   *PluginProtocolInfo

	init_func func(*Plugin)

	// private
	fromc bool
}

func NewPlugin(pi *PluginInfo, ppi *PluginProtocolInfo, init_func func(*Plugin)) *Plugin {
	if pi == nil {
		log.Panicln("nil PluginInfo")
	}
	if _plugin_instance != nil {
		log.Panicln("already exists")
	}

	this := &Plugin{pi: pi, ppi: ppi, init_func: init_func}
	_plugin_instance = this
	this.cpppi = new(C.PurplePluginProtocolInfo)
	this.cppi = new(C.PurplePluginInfo)

	this.convertInfo()
	return this
}

func (this *Plugin) convertInfo() {
	this.cppi.magic = C.PURPLE_PLUGIN_MAGIC
	this.cppi.major_version = C.PURPLE_MAJOR_VERSION
	this.cppi.minor_version = C.PURPLE_MINOR_VERSION

	C._set_plugin_type(this.cppi)

	this.cppi.id = C.CString(this.pi.Id)
	this.cppi.name = C.CString(this.pi.Name)
	this.cppi.version = C.CString(this.pi.Version)
	this.cppi.summary = C.CString(this.pi.Summary)
	this.cppi.description = C.CString(this.pi.Description)
	this.cppi.author = C.CString(this.pi.Author)

	C._set_plugin_funcs(this.cppi, this.cpppi)

	// protocol info
	this.cppi.extra_info = unsafe.Pointer(this.cpppi)
	this.cpppi.options = C.OPT_PROTO_CHAT_TOPIC | C.OPT_PROTO_PASSWORD_OPTIONAL |
		C.OPT_PROTO_INVITE_MESSAGE
	this.cpppi.struct_size = C.sizeof_PurplePluginProtocolInfo
}

func newPluginInfoFrom(plugInfo *C.PurplePluginInfo) *PluginInfo {
	this := &PluginInfo{fromc: true, pi: plugInfo}

	this.Id = C.GoString(plugInfo.id)
	this.Name = C.GoString(plugInfo.name)
	this.Version = C.GoString(plugInfo.version)
	this.Summary = C.GoString(plugInfo.summary)
	this.Description = C.GoString(plugInfo.description)
	this.Author = C.GoString(plugInfo.author)
	this.Homepage = C.GoString(plugInfo.homepage)

	return this
}

func newPluginFrom(plug *C.PurplePlugin) *Plugin {
	this := &Plugin{fromc: true, cpp: plug}
	this.cpp = plug
	return this
}

func newPluginProtocolInfoFrom(protoInfo *C.PurplePluginProtocolInfo) *PluginProtocolInfo {
	this := &PluginProtocolInfo{fromc: true, ppi: protoInfo}
	return this
}

// callbacks/events
//func init_plugin(plugin *C.PurplePlugin) {}
func (this *Plugin) init(plugin *C.PurplePlugin) {
	log.Println("")
	this.cpp = plugin

	if this.init_func != nil {
		this.init_func(this)
	}
}

//export goprpl_plugin_load
func goprpl_plugin_load(plugin *C.PurplePlugin) C.gboolean {
	var this = _plugin_instance
	this.pi.Load(this)
	return C.TRUE
}

//export goprpl_plugin_unload
func goprpl_plugin_unload(plugin *C.PurplePlugin) C.gboolean {
	var this = _plugin_instance
	this.pi.Unload(this)
	return C.TRUE
}

//export goprpl_plugin_destroy
func goprpl_plugin_destroy(plugin *C.PurplePlugin) {
	var this = _plugin_instance
	this.pi.Destroy(this)
}

//export goprpl_blist_icon
func goprpl_blist_icon(ac *C.PurpleAccount, b *C.PurpleBuddy) *C.char {
	var this = _plugin_instance
	icon := this.ppi.BlistIcon()
	return C.CString(icon)
}

//export goprpl_login
func goprpl_login(ac *C.PurpleAccount) {
	var this = _plugin_instance
	if this.ppi.Login != nil {
		this.ppi.Login(newAccountFrom(ac))
	}
}

//export goprpl_close
func goprpl_close(gc *C.PurpleConnection) {
	var this = _plugin_instance
	if this.ppi.Close != nil {
		this.ppi.Close(newConnectionFrom(gc))
	}
}

//export goprpl_status_types
func goprpl_status_types(ac *C.PurpleAccount) *C.GList {
	var this = _plugin_instance

	var types *C.GList
	if this.ppi.StatusTypes != nil {
		stys := this.ppi.StatusTypes(newAccountFrom(ac))
		for _, sty := range stys {
			types = C.g_list_append(types, sty.sty)
		}
	}
	if types == nil { // add default online status types
		var stype *C.PurpleStatusType

		stype = C.purple_status_type_new(C.PURPLE_STATUS_AVAILABLE,
			C.CString("tox_online"), C.CString("Online"), C.TRUE)
		types = C.g_list_append(types, stype)

		stype = C.purple_status_type_new(C.PURPLE_STATUS_OFFLINE,
			C.CString("tox_offline"), C.CString("Offline"), C.TRUE)
		types = C.g_list_append(types, stype)
	}

	if types == nil {
		log.Panicln("wtf")
	}
	return types
}

// optional callbacks
//export goprpl_chat_info
func goprpl_chat_info(gc *C.PurpleConnection) *C.GList {
	var this = _plugin_instance

	var m *C.GList
	if this.ppi.ChatInfo != nil {
		infos := this.ppi.ChatInfo(newConnectionFrom(gc))
		if infos == nil {
			log.Panicln("need chat info")
		}

		var pce *ProtoChatEntry
		for _, info := range infos {
			// pce = NewProtoChatEntry(info, info, true)
			pce = info
			m = C.g_list_append(m, pce.get())
		}

		// another for storage
		if false {
			pce = NewProtoChatEntry("GroupNumber", "GroupNumber", false)
			m = C.g_list_append(m, pce.get())
		}
	}
	return m
}

//export goprpl_chat_info_defaults
func goprpl_chat_info_defaults(gc *C.PurpleConnection, chatName *C.char) *C.GHashTable {
	var this = _plugin_instance
	if this.ppi.ChatInfoDefaults != nil {
		this.ppi.ChatInfoDefaults(newConnectionFrom(gc), C.GoString(chatName))
	}

	var defaults *C.GHashTable
	defaults = C.goprpl_hash_table_new_full()
	if chatName != nil {
		dchan := C.g_strdup((*C.gchar)(C.CString("ToxChannel")))
		C.g_hash_table_insert(defaults, dchan,
			(*C.gchar)(C.g_strdup((*C.gchar)(chatName))))
	}
	return defaults
}

//export goprpl_send_im
func goprpl_send_im(gc *C.PurpleConnection, who *C.char, msg *C.char, flags C.PurpleMessageFlags) C.int {
	var this = _plugin_instance
	if this.ppi.SendIM != nil {
		ret := this.ppi.SendIM(newConnectionFrom(gc), C.GoString(who), C.GoString(msg))
		return C.int(ret)
	}

	return C.int(-1)
}

//export goprpl_join_chat
func goprpl_join_chat(gc *C.PurpleConnection, comp *C.GHashTable) {
	var this = _plugin_instance
	if this.ppi.JoinChat != nil {
		this.ppi.JoinChat(newConnectionFrom(gc), newGHashTableFrom(comp))
	}
}

//export goprpl_reject_chat
func goprpl_reject_chat(gc *C.PurpleConnection, comp *C.GHashTable) {
	var this = _plugin_instance
	if this.ppi.RejectChat != nil {
		this.ppi.RejectChat(newConnectionFrom(gc), newGHashTableFrom(comp))
	}
}

//export goprpl_get_chat_name
func goprpl_get_chat_name(comp *C.GHashTable) *C.char {
	var this = _plugin_instance
	if this.ppi.GetChatName != nil {
		this.ppi.GetChatName(newGHashTableFrom(comp))
	}
	return nil
}

//export goprpl_chat_invite
func goprpl_chat_invite(gc *C.PurpleConnection, id C.int, message *C.char, who *C.char) {
	var this = _plugin_instance
	if this.ppi.ChatInvite != nil {
		this.ppi.ChatInvite(newConnectionFrom(gc), int(id), C.GoString(message), C.GoString(who))
	}
}

//export goprpl_chat_leave
func goprpl_chat_leave(gc *C.PurpleConnection, id C.int) {
	var this = _plugin_instance
	if this.ppi.ChatLeave != nil {
		this.ppi.ChatLeave(newConnectionFrom(gc), int(id))
	}
}

//export goprpl_chat_whisper
func goprpl_chat_whisper(gc *C.PurpleConnection, id C.int, who *C.char, message *C.char) {
	var this = _plugin_instance
	if this.ppi.ChatWhisper != nil {
		this.ppi.ChatWhisper(newConnectionFrom(gc), int(id), C.GoString(who), C.GoString(message))
	}
}

//export goprpl_chat_send
func goprpl_chat_send(gc *C.PurpleConnection, id C.int, message *C.char, flags C.PurpleMessageFlags) C.int {
	var this = _plugin_instance
	if this.ppi.ChatSend != nil {
		ret := this.ppi.ChatSend(newConnectionFrom(gc), int(id), C.GoString(message), int(flags))
		return C.int(ret)
	}
	return C.int(0)
}

//export goprpl_roomlist_get_list
func goprpl_roomlist_get_list(gc *C.PurpleConnection) *C.PurpleRoomlist {
	var this = _plugin_instance
	if this.ppi.RoomlistGetList != nil {
		this.ppi.RoomlistGetList(newConnectionFrom(gc))
	}
	return nil
}

//export goprpl_add_buddy_with_invite
func goprpl_add_buddy_with_invite(gc *C.PurpleConnection, buddy *C.PurpleBuddy, group *C.PurpleGroup, message *C.char) {
	var this = _plugin_instance
	if this.ppi.AddBuddyWithInvite != nil {
		this.ppi.AddBuddyWithInvite(newConnectionFrom(gc),
			newBuddyFrom(buddy), newGroupFrom(group), C.GoString(message))
	}
}

//export goprpl_remove_buddy
func goprpl_remove_buddy(gc *C.PurpleConnection, buddy *C.PurpleBuddy, group *C.PurpleGroup) {
	var this = _plugin_instance
	if this.ppi.RemoveBuddy != nil {
		this.ppi.RemoveBuddy(newConnectionFrom(gc), newBuddyFrom(buddy), newGroupFrom(group))
	}
}

//export goprpl_add_permit
func goprpl_add_permit(gc *C.PurpleConnection, name *C.char) {
	var this = _plugin_instance
	if this.ppi.AddPermit != nil {
		this.ppi.AddPermit(newConnectionFrom(gc), C.GoString(name))
	}
}

//export goprpl_add_deny
func goprpl_add_deny(gc *C.PurpleConnection, name *C.char) {
	var this = _plugin_instance
	if this.ppi.AddDeny != nil {
		this.ppi.AddDeny(newConnectionFrom(gc), C.GoString(name))
	}
}

//export goprpl_rem_permit
func goprpl_rem_permit(gc *C.PurpleConnection, name *C.char) {
	var this = _plugin_instance
	if this.ppi.RemPermit != nil {
		this.ppi.RemPermit(newConnectionFrom(gc), C.GoString(name))
	}
}

//export goprpl_rem_deny
func goprpl_rem_deny(gc *C.PurpleConnection, name *C.char) {
	var this = _plugin_instance
	if this.ppi.RemDeny != nil {
		this.ppi.RemDeny(newConnectionFrom(gc), C.GoString(name))
	}
}

//export goprpl_get_info
func goprpl_get_info(gc *C.PurpleConnection, name *C.char) {
	var this = _plugin_instance
	if this.ppi.GetInfo != nil {
		this.ppi.GetInfo(newConnectionFrom(gc), C.GoString(name))
	}
}

//export goprpl_status_text
func goprpl_status_text(buddy *C.PurpleBuddy) *C.char {
	var this = _plugin_instance
	if this.ppi.StatusText != nil {
		stxt := this.ppi.StatusText(newBuddyFrom(buddy))
		return C.CString(stxt)
	}
	return nil
}

var _plugin_instance *Plugin = nil

// when call go's init() and purple's purple_init_plugin
//export purple_init_plugin
func purple_init_plugin(plugin *C.PurplePlugin) C.gboolean {
	log.Println(plugin, MyTid2())
	runtime.LockOSThread()

	// _plugin_instance = NewPlugin()
	if _plugin_instance == nil {
		log.Panicln("failed")
	}
	plugin.info = _plugin_instance.cppi

	_plugin_instance.init(plugin)
	// var init_func = init_plugin
	// init_func(plugin)
	return C.purple_plugin_register(plugin)
}
