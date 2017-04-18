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

static void _set_plugin_type(PurplePluginInfo *pi, int type_) {
    pi->type = PURPLE_PLUGIN_STANDARD;
    pi->type = PURPLE_PLUGIN_PROTOCOL;
   //  pi->type = PURPLE_PLUGIN_LOADER;
    pi->type = type_;
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
extern void goprpl_set_chat_topic(PurpleConnection *gc, int id, char *topic);
extern char* goprpl_normalize(PurpleConnection *gc, char *who);
// not tested
extern char* goprpl_list_emblem(PurpleBuddy *buddy);
extern void goprpl_tooltip_text(PurpleBuddy *buddy, PurpleNotifyUserInfo *userInfo, gboolean full);
extern guint goprpl_send_typing(PurpleConnection *gc, char *name, PurpleTypingState state);
extern void goprpl_keepalive(PurpleConnection *gc);
extern void goprpl_register_user(PurpleAccount *ac);
// UnregisterUser func(ac *Account, cb PurpleAccountUnregistrationCb, ...)
extern gboolean goprpl_offline_message(PurpleBuddy *buddy);
extern gint goprpl_send_raw(PurpleConnection *gc, char* buf, int len);
// RoomlistRoomSerialize func(room *RoomlistRoom) string
extern gboolean goprpl_send_attention(PurpleConnection *gc, char *userName, guint atype);
extern GList* goprpl_get_attension_types(PurpleAccount *ac);
extern gboolean goprpl_can_receive_file(PurpleConnection *gc, char *who);
extern void goprpl_send_file(PurpleConnection *gc, char *who, char *filename);
extern PurpleXfer* goprpl_new_xfer(PurpleConnection *gc, char *who);

static void _set_plugin_funcs(PurplePluginInfo *pi, PurplePluginProtocolInfo *pppi) {
    pi->load = goprpl_plugin_load;
    pi->unload = goprpl_plugin_unload;
    pi->destroy = goprpl_plugin_destroy;

    // TODO can not set call unconditionally. should check the PRPL's real setting
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
    pppi->set_chat_topic = goprpl_set_chat_topic;
    pppi->normalize = goprpl_normalize;
    // not tested
    pppi->list_emblem = goprpl_list_emblem;
    pppi->tooltip_text = goprpl_tooltip_text;
    pppi->send_typing = goprpl_send_typing;
    pppi->keepalive = goprpl_send_typing;
    pppi->register_user = goprpl_register_user;
    // UnregisterUser func(ac *Account, cb PurpleAccountUnregistrationCb, ...)
    pppi->offline_message = goprpl_offline_message;
    pppi->send_raw = goprpl_send_raw;
    // RoomlistRoomSerialize func(room *RoomlistRoom) string
    pppi->send_attention = goprpl_send_attention;
    pppi->get_attention_types = goprpl_get_attension_types;

    // file transfer
    pppi->can_receive_file = goprpl_can_receive_file;
    pppi->send_file = goprpl_send_file;
    pppi->new_xfer = goprpl_new_xfer;

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
	Type int

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
	// optional fields
	Options  ProtocolOptions
	IconSpec BuddyIconSpec

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
	SetChatTopic       func(gc *Connection, id int, topic string)
	Normalize          func(gc *Connection, who string) string
	// not tested
	ListEmblem   func(buddy *Buddy) string
	TooltipText  func(buddy *Buddy, userInfo *NotifyUserInfo, full bool)
	SendTyping   func(gc *Connection, name string, state int) uint
	KeepAlive    func(gc *Connection)
	RegisterUser func(ac *Account)
	// UnregisterUser func(ac *Account, cb PurpleAccountUnregistrationCb, ...)
	OfflineMessage func(buddy *Buddy) bool
	SendRaw        func(gc *Connection, buf string, len int) int
	// RoomlistRoomSerialize func(room *RoomlistRoom) string
	SendAttention     func(gc *Connection, userName string, atype uint) bool
	GetAttentionTypes func(ac *Account) *GList

	// file transfer
	CanReceiveFile func(gc *Connection, who string) bool
	SendFile       func(gc *Connection, who string, filename string)
	NewXfer        func(gc *Connection, who string) *Xfer

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

	// because type is golang's keyword, shit
	C._set_plugin_type(this.cppi, C.int(this.pi.Type))
	if !(this.pi.Type >= PLUGIN_UNKNOWN && this.pi.Type <= PLUGIN_PROTOCOL) {
		log.Panicln("Unsupported plugin type:", this.pi.Type)
	}

	this.cppi.id = C.CString(this.pi.Id)
	this.cppi.name = C.CString(this.pi.Name)
	this.cppi.version = C.CString(this.pi.Version)
	this.cppi.summary = C.CString(this.pi.Summary)
	this.cppi.description = C.CString(this.pi.Description)
	this.cppi.author = C.CString(this.pi.Author)

	// this will set all without check
	C._set_plugin_funcs(this.cppi, this.cpppi) // c version
	this.set_plugin_funcs()                    // go version
	// this will check and unset nil callback functions
	this._unset_plugin_funcs()

	// protocol info
	this.cppi.extra_info = unsafe.Pointer(this.cpppi)
	if this.cpppi.options == 0 {
		this.cpppi.options = C.OPT_PROTO_CHAT_TOPIC | C.OPT_PROTO_PASSWORD_OPTIONAL |
			C.OPT_PROTO_INVITE_MESSAGE
	}
	this.cpppi.struct_size = C.sizeof_PurplePluginProtocolInfo
}

// go version of C._set_plugin_funcs()
func (this *Plugin) set_plugin_funcs() {
	this.cpppi.options = C.PurpleProtocolOptions(this.ppi.Options)
	var ispec C.PurpleBuddyIconSpec
	ispec.format = C.CString(this.ppi.IconSpec.Format)
	ispec.min_width = C.int(this.ppi.IconSpec.MinWidth)
	ispec.min_height = C.int(this.ppi.IconSpec.MinHeight)
	ispec.max_width = C.int(this.ppi.IconSpec.MaxWidth)
	ispec.max_height = C.int(this.ppi.IconSpec.MinHeight)
	ispec.max_filesize = C.size_t(this.ppi.IconSpec.MaxFilesize)
	ispec.scale_rules = C.PurpleIconScaleRules(this.ppi.IconSpec.ScaleRules)
	this.cpppi.icon_spec = ispec

	// must
	this.cpppi.list_icon = go2cfn(C.goprpl_blist_icon)
	this.cpppi.status_types = go2cfn(C.goprpl_status_types)
	this.cpppi.login = go2cfn(C.goprpl_login)
	this.cpppi.close = go2cfn(C.goprpl_close)

	// optional, might by Proirity high to low
	this.cpppi.chat_info = go2cfn(C.goprpl_chat_info)
	this.cpppi.chat_info_defaults = go2cfn(C.goprpl_chat_info_defaults)
	this.cpppi.send_im = go2cfn(C.goprpl_send_im)
	this.cpppi.join_chat = go2cfn(C.goprpl_join_chat)
	this.cpppi.reject_chat = go2cfn(C.goprpl_reject_chat)
	this.cpppi.get_chat_name = go2cfn(C.goprpl_get_chat_name)
	this.cpppi.chat_invite = go2cfn(C.goprpl_chat_invite)
	this.cpppi.chat_leave = go2cfn(C.goprpl_chat_leave)
	this.cpppi.chat_whisper = go2cfn(C.goprpl_chat_whisper)
	this.cpppi.chat_send = go2cfn(C.goprpl_chat_send)
	this.cpppi.roomlist_get_list = go2cfn(C.goprpl_roomlist_get_list)
	this.cpppi.add_buddy_with_invite = go2cfn(C.goprpl_add_buddy_with_invite)
	this.cpppi.remove_buddy = go2cfn(C.goprpl_remove_buddy)
	this.cpppi.add_permit = go2cfn(C.goprpl_add_permit)
	this.cpppi.add_deny = go2cfn(C.goprpl_add_deny)
	this.cpppi.rem_permit = go2cfn(C.goprpl_rem_permit)
	this.cpppi.rem_deny = go2cfn(C.goprpl_rem_deny)
	this.cpppi.get_info = go2cfn(C.goprpl_get_info)
	this.cpppi.status_text = go2cfn(C.goprpl_status_text)
	this.cpppi.set_chat_topic = go2cfn(C.goprpl_set_chat_topic)
	this.cpppi.normalize = go2cfn(C.goprpl_normalize)

	// not tested
	this.cpppi.list_emblem = go2cfn(C.goprpl_list_emblem)
	this.cpppi.tooltip_text = go2cfn(C.goprpl_tooltip_text)
	this.cpppi.send_typing = go2cfn(C.goprpl_send_typing)
	this.cpppi.keepalive = go2cfn(C.goprpl_keepalive)
	this.cpppi.register_user = go2cfn(C.goprpl_register_user)
	// UnregisterUser func(ac *Account, cb PurpleAccountUnregistrationCb, ...)

	this.cpppi.offline_message = go2cfn(C.goprpl_offline_message)
	this.cpppi.send_raw = go2cfn(C.goprpl_send_raw)
	// RoomlistRoomSerialize func(room *RoomlistRoom) string
	this.cpppi.send_attention = go2cfn(C.goprpl_send_attention)
	this.cpppi.get_attention_types = go2cfn(C.goprpl_get_attension_types)

	// file transfer
	this.cpppi.can_receive_file = go2cfn(C.goprpl_can_receive_file)
	this.cpppi.send_file = go2cfn(C.goprpl_send_file)
	this.cpppi.new_xfer = go2cfn(C.goprpl_new_xfer)
}

// this will check and unset nil callback functions
func (this *Plugin) _unset_plugin_funcs() {
	// must
	lackofmust := false
	if this.ppi.BlistIcon == nil {
		this.cpppi.list_icon = nil
		lackofmust = true
		log.Println("BlistIcon method must set.")
	}
	if this.ppi.StatusTypes == nil {
		this.cpppi.status_types = nil
		lackofmust = true
		log.Println("StatusTypes method must set.")
	}
	if this.ppi.Login == nil {
		this.cpppi.login = nil
		lackofmust = true
		log.Println("Login method must set.")
	}
	if this.ppi.Close == nil {
		this.cpppi.close = nil
		lackofmust = true
		log.Println("Close method must set.")
	}
	if lackofmust {
		log.Fatalln("Lack of some must set methods.")
	}

	// optional, might by Proirity high to low
	if this.ppi.ChatInfo == nil {
		this.cpppi.chat_info = nil
	}
	if this.ppi.ChatInfoDefaults == nil {
		this.cpppi.chat_info_defaults = nil
	}
	if this.ppi.SendIM == nil {
		this.cpppi.send_im = nil
	}
	if this.ppi.JoinChat == nil {
		this.cpppi.join_chat = nil
	}
	if this.ppi.RejectChat == nil {
		this.cpppi.reject_chat = nil
	}
	if this.ppi.GetChatName == nil {
		this.cpppi.get_chat_name = nil
	}
	if this.ppi.ChatInvite == nil {
		this.cpppi.chat_invite = nil
	}
	if this.ppi.ChatLeave == nil {
		this.cpppi.chat_leave = nil
	}
	if this.ppi.ChatWhisper == nil {
		this.cpppi.chat_whisper = nil
	}
	if this.ppi.ChatSend == nil {
		this.cpppi.chat_send = nil
	}
	if this.ppi.RoomlistGetList == nil {
		this.cpppi.roomlist_get_list = nil
	}
	if this.ppi.AddBuddyWithInvite == nil {
		this.cpppi.add_buddy_with_invite = nil
	}
	if this.ppi.RemoveBuddy == nil {
		this.cpppi.remove_buddy = nil
	}
	if this.ppi.AddPermit == nil {
		this.cpppi.add_permit = nil
	}
	if this.ppi.AddDeny == nil {
		this.cpppi.add_deny = nil
	}
	if this.ppi.RemPermit == nil {
		this.cpppi.rem_permit = nil
	}
	if this.ppi.RemDeny == nil {
		this.cpppi.rem_deny = nil
	}
	if this.ppi.GetInfo == nil {
		this.cpppi.get_info = nil
	}
	if this.ppi.StatusText == nil {
		this.cpppi.status_text = nil
	}
	if this.ppi.SetChatTopic == nil {
		this.cpppi.set_chat_topic = nil
	}
	if this.ppi.Normalize == nil {
		this.cpppi.normalize = nil
	}
	// not tested
	if this.ppi.ListEmblem == nil {
		this.cpppi.list_emblem = nil
	}
	if this.ppi.TooltipText == nil {
		this.cpppi.tooltip_text = nil
	}
	if this.ppi.SendTyping == nil {
		this.cpppi.send_typing = nil
	}
	if this.ppi.KeepAlive == nil {
		this.cpppi.keepalive = nil
	}
	if this.ppi.RegisterUser == nil {
		this.cpppi.register_user = nil
	}
	// UnregisterUser func(ac *Account, cb PurpleAccountUnregistrationCb, ...)
	if this.ppi.OfflineMessage == nil {
		this.cpppi.offline_message = nil
	}
	if this.ppi.SendRaw == nil {
		this.cpppi.send_raw = nil
	}
	// RoomlistRoomSerialize func(room *RoomlistRoom) string
	if this.ppi.SendAttention == nil {
		this.cpppi.send_attention = nil
	}
	if this.ppi.GetAttentionTypes == nil {
		this.cpppi.get_attention_types = nil
	}

	// file transfer
	if this.ppi.CanReceiveFile == nil {
		this.cpppi.can_receive_file = nil
	}
	if this.ppi.SendFile == nil {
		this.cpppi.send_file = nil
	}
	if this.ppi.NewXfer == nil {
		this.cpppi.new_xfer = nil
	}
}

// this is not purple's function
func (this *Plugin) AddProtocolOption(ao *AccountOption) {
	this.cpppi.protocol_options = C.g_list_append(this.cpppi.protocol_options, (C.gpointer)(ao.ao))
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
			types = C.g_list_append(types, (C.gpointer)(sty.sty))
		}
	}
	if types == nil { // add default online status types
		var stype *C.PurpleStatusType

		stype = C.purple_status_type_new(C.PURPLE_STATUS_AVAILABLE,
			CCString("tox_online").Ptr, CCString("Online").Ptr, C.TRUE)
		types = C.g_list_append(types, (C.gpointer)(stype))

		stype = C.purple_status_type_new(C.PURPLE_STATUS_OFFLINE,
			CCString("tox_offline").Ptr, CCString("Offline").Ptr, C.TRUE)
		types = C.g_list_append(types, (C.gpointer)(stype))
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
			m = C.g_list_append(m, (C.gpointer)(pce.get()))
		}

		// another for storage
		if false {
			pce = NewProtoChatEntry("GroupNumber", "GroupNumber", false)
			m = C.g_list_append(m, (C.gpointer)(pce.get()))
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
		dchan := C.g_strdup((*C.gchar)(CCString("ToxChannel").Ptr))
		C.g_hash_table_insert(defaults, (C.gpointer)(dchan),
			(C.gpointer)(C.g_strdup((*C.gchar)(chatName))))
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
		if len(stxt) > 0 {
			return C.CString(stxt)
		}
	}
	return nil
}

//export goprpl_set_chat_topic
func goprpl_set_chat_topic(gc *C.PurpleConnection, id C.int, topic *C.char) {
	var this = _plugin_instance
	if this.ppi.SetChatTopic != nil {
		this.ppi.SetChatTopic(newConnectionFrom(gc), int(id), C.GoString(topic))
	}
}

//export goprpl_normalize
func goprpl_normalize(gc *C.PurpleConnection, who *C.char) *C.char {
	var this = _plugin_instance
	if this.ppi.Normalize != nil {
		norm := this.ppi.Normalize(newConnectionFrom(gc), C.GoString(who))
		return C.CString(norm)
	}
	return nil
}

//export goprpl_list_emblem
func goprpl_list_emblem(buddy *C.PurpleBuddy) *C.char {
	var this = _plugin_instance
	if this.ppi.ListEmblem != nil {
		emblem := this.ppi.ListEmblem(newBuddyFrom(buddy))
		return C.CString(emblem)
	}
	return nil
}

//export goprpl_tooltip_text
func goprpl_tooltip_text(buddy *C.PurpleBuddy, userInfo *C.PurpleNotifyUserInfo, full C.gboolean) {
	var this = _plugin_instance
	if this.ppi.TooltipText != nil {
		/*
			this.ppi.TooltipText(newBuddyFrom(buddy),
				newNotifyUserInfoFrom(userInfo), c2goBool(full))
		*/
	}
}

//export goprpl_send_typing
func goprpl_send_typing(gc *C.PurpleConnection, name *C.char, state C.PurpleTypingState) C.guint {
	var this = _plugin_instance
	if this.ppi.SendTyping != nil {
		ret := this.ppi.SendTyping(newConnectionFrom(gc), C.GoString(name), int(state))
		return C.guint(ret)
	}
	return C.guint(0)
}

//export goprpl_keepalive
func goprpl_keepalive(gc *C.PurpleConnection) {
	var this = _plugin_instance
	if this.ppi.KeepAlive != nil {
		this.ppi.KeepAlive(newConnectionFrom(gc))
	}
}

//export goprpl_register_user
func goprpl_register_user(ac *C.PurpleAccount) {
	var this = _plugin_instance
	if this.ppi.RegisterUser != nil {
		this.ppi.RegisterUser(newAccountFrom(ac))
	}
}

// UnregisterUser func(ac *Account, cb PurpleAccountUnregistrationCb, ...)
//export goprpl_offline_message
func goprpl_offline_message(buddy *C.PurpleBuddy) C.gboolean {
	var this = _plugin_instance
	if this.ppi.OfflineMessage != nil {
		ret := this.ppi.OfflineMessage(newBuddyFrom(buddy))
		return go2cBool(ret)
	}
	return C.FALSE
}

//export goprpl_send_raw
func goprpl_send_raw(gc *C.PurpleConnection, buf *C.char, len C.int) C.int {
	var this = _plugin_instance
	if this.ppi.SendRaw != nil {
		ret := this.ppi.SendRaw(newConnectionFrom(gc), C.GoString(buf), int(len))
		return C.int(ret)
	}
	return C.int(0)
}

// RoomlistRoomSerialize func(room *RoomlistRoom) string
//export goprpl_send_attention
func goprpl_send_attention(gc *C.PurpleConnection, userName *C.char, atype C.guint) C.gboolean {
	var this = _plugin_instance
	if this.ppi.SendAttention != nil {
		ret := this.ppi.SendAttention(newConnectionFrom(gc), C.GoString(userName), uint(atype))
		return go2cBool(ret)
	}
	return C.FALSE
}

//export goprpl_get_attension_types
func goprpl_get_attension_types(ac *C.PurpleAccount) *C.GList {
	var this = _plugin_instance
	if this.ppi.GetAttentionTypes != nil {
		ret := this.ppi.GetAttentionTypes(newAccountFrom(ac))
		return ret.lst
	}
	return nil
}

// file transfer
//export goprpl_can_receive_file
func goprpl_can_receive_file(gc *C.PurpleConnection, who *C.char) C.gboolean {
	var this = _plugin_instance
	if this.ppi.CanReceiveFile != nil {
		ret := this.ppi.CanReceiveFile(newConnectionFrom(gc), C.GoString(who))
		return go2cBool(ret)
	}
	return C.FALSE
}

//export goprpl_send_file
func goprpl_send_file(gc *C.PurpleConnection, who *C.char, filename *C.char) {
	var this = _plugin_instance
	if this.ppi.SendFile != nil {
		this.ppi.SendFile(newConnectionFrom(gc), C.GoString(who), C.GoString(filename))
	}
}

//export goprpl_new_xfer
func goprpl_new_xfer(gc *C.PurpleConnection, who *C.char) *C.PurpleXfer {
	var this = _plugin_instance
	if this.ppi.NewXfer != nil {
		ret := this.ppi.NewXfer(newConnectionFrom(gc), C.GoString(who))
		return ret.xfer
	}

	return nil
}
