package purple

/*
#include <glib.h>

#include "notify.h"
#include "plugin.h"
#include "version.h"
#include "prpl.h"
#include <string.h>
*/
import "C"

// import "unsafe"

func (this *Plugin) GetId() string {
	ret := C.purple_plugin_get_id(this.cpp)
	return C.GoString((*C.char)(ret))
}
func (this *Plugin) GetName() string {
	ret := C.purple_plugin_get_name(this.cpp)
	return C.GoString((*C.char)(ret))
}
func (this *Plugin) GetVersion() string {
	ret := C.purple_plugin_get_version(this.cpp)
	return C.GoString((*C.char)(ret))
}
func (this *Plugin) GetAuthor() string {
	ret := C.purple_plugin_get_author(this.cpp)
	return C.GoString((*C.char)(ret))
}
func (this *Plugin) GetSummary() string {
	ret := C.purple_plugin_get_summary(this.cpp)
	return C.GoString((*C.char)(ret))
}
func (this *Plugin) GetDescription() string {
	ret := C.purple_plugin_get_description(this.cpp)
	return C.GoString((*C.char)(ret))
}
func (this *Plugin) GetHomepage() string {
	ret := C.purple_plugin_get_homepage(this.cpp)
	return C.GoString((*C.char)(ret))
}

func PluginsFindWithName(name string) *Plugin {
	plugin := C.purple_plugins_find_with_name(CCString(name).Ptr)
	return newPluginFrom(plugin)
}

func PluginsFindWithFilename(filename string) *Plugin {
	plugin := C.purple_plugins_find_with_filename(CCString(filename).Ptr)
	return newPluginFrom(plugin)
}
func PluginsFindWithBasename(basename string) *Plugin {

	plugin := C.purple_plugins_find_with_basename(CCString(basename).Ptr)
	return newPluginFrom(plugin)
}
func PluginsFindWithId(id string) *Plugin {
	plugin := C.purple_plugins_find_with_id(CCString(id).Ptr)
	return newPluginFrom(plugin)
}

func PluginsGetLoaded() []*Plugin {
	lst := C.purple_plugins_get_loaded()

	plugs := make([]*Plugin, 0)
	newGListFrom(lst).Each(func(item C.gpointer) {
		cplug := (*C.PurplePlugin)(item)
		plugs = append(plugs, newPluginFrom(cplug))
	})
	return plugs
}
func PluginsGetProtocols() []*Plugin {
	lst := C.purple_plugins_get_protocols()

	plugs := make([]*Plugin, 0)
	newGListFrom(lst).Each(func(item C.gpointer) {
		cplug := (*C.PurplePlugin)(item)
		plugs = append(plugs, newPluginFrom(cplug))
	})
	return plugs
}
func PluginsGetAll() []*Plugin {
	lst := C.purple_plugins_get_all()

	plugs := make([]*Plugin, 0)
	newGListFrom(lst).Each(func(item C.gpointer) {
		cplug := (*C.PurplePlugin)(item)
		plugs = append(plugs, newPluginFrom(cplug))
	})
	return plugs
}
