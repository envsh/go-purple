package plugin

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
extern char* goprpl_blist_icon(PurpleAccount *a, PurpleBuddy *b);
extern void goprpl_login(PurpleAccount *account);
extern void goprpl_close(PurpleConnection *gc);
extern GList* goprpl_status_types(PurpleAccount *a);

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
	Login     func()
	Close     func()

	// optional

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
	this.ppi.Login()
}

//export goprpl_close
func goprpl_close(gc *C.PurpleConnection) {
	var this = _plugin_instance
	this.ppi.Close()
}

//export goprpl_status_types
func goprpl_status_types(a *C.PurpleAccount) *C.GList {
	var stype *C.PurpleStatusType
	var types *C.GList

	stype = C.purple_status_type_new(C.PURPLE_STATUS_AVAILABLE, nil, nil, C.TRUE)
	types = C.g_list_append(types, stype)

	stype = C.purple_status_type_new_full(C.PURPLE_STATUS_AWAY, nil, nil, C.TRUE, C.TRUE, C.FALSE)
	types = C.g_list_append(types, stype)

	stype = C.purple_status_type_new_full(C.PURPLE_STATUS_UNAVAILABLE, C.CString("ToxBusy"), C.CString("Busy"), C.TRUE, C.TRUE, C.FALSE)
	types = C.g_list_append(types, stype)

	stype = C.purple_status_type_new(C.PURPLE_STATUS_OFFLINE, nil, nil, C.TRUE)
	types = C.g_list_append(types, stype)

	return types
}

var _plugin_instance *Plugin = nil

// when call go's init() and purple's purple_init_plugin
//export purple_init_plugin
func purple_init_plugin(plugin *C.PurplePlugin) C.gboolean {
	log.Println(plugin)

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
