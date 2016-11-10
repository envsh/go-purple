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
    pppi->struct_size = sizeof(PurplePluginProtocolInfo);
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
	"fmt"
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
}

type PluginProtocolInfo struct {
	// must
	BlistIcon func() string
	Login     func(*Account)
	Close     func(*Connection)

	// optional
	ChatInfo         func(gc *Connection) []string
	ChatInfoDefaults func(gc *Connection, chat_name string) map[string]string
	SendIM           func(gc *Connection, who string, message string) int
	JoinChat         func(gc *Connection, comp interface{})
	RejectChat       func(gc *Connection, comp interface{})
	GetChatName      func(comp interface{}) string
	ChatInvite       func(gc *Connection, id int, message string, who string)
	ChatLeave        func(gc *Connection, id int)
	ChatWhisper      func(gc *Connection, id int, who string, message string)
	ChatSend         func(gc *Connection, id int, message string, flags int) int
}

type Plugin struct {
	cpp   *C.PurplePlugin
	cppi  *C.PurplePluginInfo
	cpppi *C.PurplePluginProtocolInfo
	pi    *PluginInfo
	ppi   *PluginProtocolInfo

	init_func func(*Plugin)
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
}

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
func goprpl_blist_icon(a *C.PurpleAccount, b *C.PurpleBuddy) *C.char {
	var this = _plugin_instance
	icon := this.ppi.BlistIcon()
	return C.CString(icon)
}

//export goprpl_login
func goprpl_login(a *C.PurpleAccount) {
	var this = _plugin_instance
	if this.ppi.Login != nil {
		this.ppi.Login(newAccountWrapper(a))
	}
}

//export goprpl_close
func goprpl_close(gc *C.PurpleConnection) {
	var this = _plugin_instance
	if this.ppi.Close != nil {
		this.ppi.Close(newConnectWrapper(gc))
	}
}

//export goprpl_status_types
func goprpl_status_types(a *C.PurpleAccount) *C.GList {
	var stype *C.PurpleStatusType
	var types *C.GList

	stype = C.purple_status_type_new(C.PURPLE_STATUS_AVAILABLE,
		C.CString("tox_online"), C.CString("Online"), C.TRUE)
	types = C.g_list_append(types, stype)

	stype = C.purple_status_type_new_full(C.PURPLE_STATUS_AWAY,
		C.CString("tox_away"), C.CString("Away"), C.TRUE, C.TRUE, C.FALSE)
	types = C.g_list_append(types, stype)

	stype = C.purple_status_type_new_full(C.PURPLE_STATUS_UNAVAILABLE,
		C.CString("tox_busy"), C.CString("Busy"), C.TRUE, C.TRUE, C.FALSE)
	types = C.g_list_append(types, stype)

	stype = C.purple_status_type_new(C.PURPLE_STATUS_OFFLINE,
		C.CString("tox_offline"), C.CString("Offline"), C.TRUE)
	types = C.g_list_append(types, stype)

	return types
}

// optional callbacks
//export goprpl_chat_info
func goprpl_chat_info(gc *C.PurpleConnection) *C.GList {
	var this = _plugin_instance

	var m *C.GList
	if this.ppi.ChatInfo != nil {
		infos := this.ppi.ChatInfo(newConnectWrapper(gc))
		if infos == nil {
			log.Panicln("need chat info")
		}

		var pce *C.struct_proto_chat_entry
		// pce = new(C.struct_proto_chat_entry)  // crash!!!
		pce = (*C.struct_proto_chat_entry)(C.calloc(C.size_t(1), C.sizeof_struct_proto_chat_entry))
		pce.label = C.CString("_GoToxChannel")
		pce.identifier = C.CString("GoToxChannel")
		pce.required = C.TRUE
		pce.label = C.CString(fmt.Sprintf("_%s", infos[0]))
		pce.identifier = C.CString(infos[0])
		pce.required = C.TRUE
		m = C.g_list_append(m, pce)
	}
	return m
}

//export goprpl_chat_info_defaults
func goprpl_chat_info_defaults(gc *C.PurpleConnection, chatName *C.char) *C.GHashTable {
	var this = _plugin_instance
	if this.ppi.ChatInfoDefaults != nil {
		this.ppi.ChatInfoDefaults(newConnectWrapper(gc), C.GoString(chatName))
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
		ret := this.ppi.SendIM(newConnectWrapper(gc), C.GoString(who), C.GoString(msg))
		return C.int(ret)
	}

	return C.int(-1)
}

//export goprpl_join_chat
func goprpl_join_chat(gc *C.PurpleConnection, comp *C.GHashTable) {
	var this = _plugin_instance
	if this.ppi.JoinChat != nil {
		this.ppi.JoinChat(newConnectWrapper(gc), comp)
	}
}

//export goprpl_reject_chat
func goprpl_reject_chat(gc *C.PurpleConnection, comp *C.GHashTable) {
	var this = _plugin_instance
	if this.ppi.RejectChat != nil {
		this.ppi.RejectChat(newConnectWrapper(gc), comp)
	}
}

//export goprpl_get_chat_name
func goprpl_get_chat_name(comp *C.GHashTable) *C.char {
	var this = _plugin_instance
	if this.ppi.GetChatName != nil {
		this.ppi.GetChatName(comp)
	}
	return nil
}

//export goprpl_chat_invite
func goprpl_chat_invite(gc *C.PurpleConnection, id C.int, message *C.char, who *C.char) {
	var this = _plugin_instance
	if this.ppi.ChatInvite != nil {
		this.ppi.ChatInvite(newConnectWrapper(gc), int(id), C.GoString(message), C.GoString(who))
	}
}

//export goprpl_chat_leave
func goprpl_chat_leave(gc *C.PurpleConnection, id C.int) {
	var this = _plugin_instance
	if this.ppi.ChatLeave != nil {
		this.ppi.ChatLeave(newConnectWrapper(gc), int(id))
	}
}

//export goprpl_chat_whisper
func goprpl_chat_whisper(gc *C.PurpleConnection, id C.int, who *C.char, message *C.char) {
	var this = _plugin_instance
	if this.ppi.ChatWhisper != nil {
		this.ppi.ChatWhisper(newConnectWrapper(gc), int(id), C.GoString(who), C.GoString(message))
	}
}

//export goprpl_chat_send
func goprpl_chat_send(gc *C.PurpleConnection, id C.int, message *C.char, flags C.PurpleMessageFlags) C.int {
	var this = _plugin_instance
	if this.ppi.ChatSend != nil {
		ret := this.ppi.ChatSend(newConnectWrapper(gc), int(id), C.GoString(message), int(flags))
		return C.int(ret)
	}
	return C.int(0)
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
