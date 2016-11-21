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

func PluginsFindWithName(name string) *Plugin {
	plugin := C.purple_plugins_find_with_name(C.CString(name))
	return newPluginFrom(plugin)
}

func PluginsFindWithFilename(filename string) *Plugin {
	plugin := C.purple_plugins_find_with_filename(C.CString(filename))
	return newPluginFrom(plugin)
}
func PluginsFindWithBasename(basename string) *Plugin {

	plugin := C.purple_plugins_find_with_basename(C.CString(basename))
	return newPluginFrom(plugin)
}
func PluginsFindWithId(id string) *Plugin {
	plugin := C.purple_plugins_find_with_id(C.CString(id))
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
