package purple

/*
#include <libpurple/purple.h>
*/
import "C"
import "unsafe"

type PrefsUiOps struct {
	// private
	uiops *C.PurplePrefsUiOps
}

func newPrefsUiOpsFrom(uiops *C.PurplePrefsUiOps) *PrefsUiOps {
	this := &PrefsUiOps{}
	this.uiops = uiops
	return this
}

func PrefsSetUiOps(uiops *PrefsUiOps) {
	C.purple_prefs_set_ui_ops(uiops.uiops)
}

func PrefsGetUiOps() *PrefsUiOps {
	return newPrefsUiOpsFrom(C.purple_prefs_get_ui_ops())
}

func PrefsGetHandle() unsafe.Pointer {
	return C.purple_prefs_get_handle()
}

/**
 * Initialize core prefs
 */
func PrefsInit() {
	C.purple_prefs_init()
}

/**
 * Uninitializes the prefs subsystem.
 */
func PrefsUninit() {
	C.purple_prefs_uninit()
}

/**
 * Add a new typeless pref.
 *
 * @param name  The name of the pref
 */
func PrefsAddNone(name string) {
	C.purple_prefs_add_none(CCString(name).Ptr)
}

/**
 * Add a new boolean pref.
 *
 * @param name  The name of the pref
 * @param value The initial value to set
 */
func PrefsAddBool(name string, value bool) {
	C.purple_prefs_add_bool(CCString(name).Ptr, go2cBool(value))
}

/**
 * Add a new integer pref.
 *
 * @param name  The name of the pref
 * @param value The initial value to set
 */
func PrefsAddInt(name string, value int) {
	C.purple_prefs_add_int(CCString(name).Ptr, C.int(value))
}

/**
 * Add a new string pref.
 *
 * @param name  The name of the pref
 * @param value The initial value to set
 */
func PrefsAddString(name, value string) {
	C.purple_prefs_add_string(CCString(name).Ptr, CCString(value).Ptr)
}

/**
 * Add a new string list pref.
 *
 * @param name  The name of the pref
 * @param value The initial value to set
 * @note This function takes a copy of the strings in the value list. The list
 *       itself and original copies of the strings are up to the caller to
 *       free.
 */
func PrefsAddStringList(name string, value *GList) {
	C.purple_prefs_add_string_list(CCString(name).Ptr, value.lst)
}

/**
 * Add a new path pref.
 *
 * @param name  The name of the pref
 * @param value The initial value to set
 */
func PrefsAddPath(name, value string) {
	C.purple_prefs_add_path(CCString(name).Ptr, CCString(value).Ptr)
}

/**
 * Add a new path list pref.
 *
 * @param name  The name of the pref
 * @param value The initial value to set
 * @note This function takes a copy of the strings in the value list. The list
 *       itself and original copies of the strings are up to the caller to
 *       free.
 */
func PrefsAddPathList(name string, value *GList) {
	C.purple_prefs_add_path_list(CCString(name).Ptr, value.lst)
}

/**
 * Remove a pref.
 *
 * @param name The name of the pref
 */
func PrefsRemove(name string) {
	C.purple_prefs_remove(CCString(name).Ptr)
}

/**
 * Rename a pref
 *
 * @param oldname The old name of the pref
 * @param newname The new name for the pref
 */
func PrefsRename(oldname, newname string) {
	C.purple_prefs_rename(CCString(oldname).Ptr, CCString(newname).Ptr)
}

/**
 * Rename a boolean pref, toggling it's value
 *
 * @param oldname The old name of the pref
 * @param newname The new name for the pref
 */
func PrefsRenameBooleanToggle(oldname, newname string) {
	C.purple_prefs_rename_boolean_toggle(CCString(oldname).Ptr, CCString(newname).Ptr)
}

/**
 * Remove all prefs.
 */
func PrefsDestroy() {
	C.purple_prefs_destroy()
}

/**
 * Set raw pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 *
 * @deprecated We're not really sure what purpose this function serves, so it
 *             will be removed in 3.0.0.  Preferences values set using this
 *             function aren't serialized to prefs.xml, which could be
 *             misleading.  There is also no purple_prefs_get_generic, which
 *             means that if you can't really get the value (other in a
 *             connected callback).  If you think you have a use for this then
 *             please let us know.
 */
/* TODO: When this is removed, also remove struct purple_pref->value.generic */
func PrefsSetGeneric(name string, value unsafe.Pointer) {
	C.purple_prefs_set_generic(CCString(name).Ptr, value)
}

/**
 * Set boolean pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 */
func PrefsSetBool(name string, value bool) {
	C.purple_prefs_set_bool(CCString(name).Ptr, go2cBool(value))
}

/**
 * Set integer pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 */
func PrefsSetInt(name string, value int) {
	C.purple_prefs_set_int(CCString(name).Ptr, C.int(value))
}

/**
 * Set string pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 */
func PrefsSetString(name string, value string) {
	C.purple_prefs_set_string(CCString(name).Ptr, CCString(value).Ptr)
}

/**
 * Set string list pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 */
func PrefsSetStringList(name string, value *GList) {
	C.purple_prefs_set_string_list(CCString(name).Ptr, value.lst)
}

/**
 * Set path pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 */
func PrefsSetPath(name, value string) {
	C.purple_prefs_set_path(CCString(name).Ptr, CCString(value).Ptr)
}

/**
 * Set path list pref value
 *
 * @param name  The name of the pref
 * @param value The value to set
 */
func PrefsSetPathList(name string, value *GList) {
	C.purple_prefs_set_path_list(CCString(name).Ptr, value.lst)
}

/**
 * Check if a pref exists
 *
 * @param name The name of the pref
 * @return TRUE if the pref exists.  Otherwise FALSE.
 */
func PrefsExists(name string) bool {
	ret := C.purple_prefs_exists(CCString(name).Ptr)
	return c2goBool(ret)
}

/**
 * Get pref type
 *
 * @param name The name of the pref
 * @return The type of the pref
 */
type PrefType int

func PrefsGetType(name string) PrefType {
	ret := C.purple_prefs_get_type(CCString(name).Ptr)
	return PrefType(ret)
}

/**
 * Get boolean pref value
 *
 * @param name The name of the pref
 * @return The value of the pref
 */
func PrefsGetBool(name string) bool {
	ret := C.purple_prefs_get_bool(CCString(name).Ptr)
	return c2goBool(ret)
}

/**
 * Get integer pref value
 *
 * @param name The name of the pref
 * @return The value of the pref
 */
func PrefsGetInt(name string) int {
	ret := C.purple_prefs_get_int(CCString(name).Ptr)
	return int(ret)
}

/**
 * Get string pref value
 *
 * @param name The name of the pref
 * @return The value of the pref
 */
func PrefsGetString(name string) string {
	ret := C.purple_prefs_get_string(CCString(name).Ptr)
	return C.GoString(ret)
}

/**
 * Get string list pref value
 *
 * @param name The name of the pref
 * @return The value of the pref
 */
func PrefsGetStringList(name string) *GList {
	ret := C.purple_prefs_get_string_list(CCString(name).Ptr)
	return newGListFrom(ret)
}

/**
 * Get path pref value
 *
 * @param name The name of the pref
 * @return The value of the pref
 */
func PrefsGetPath(name string) string {
	ret := C.purple_prefs_get_path(CCString(name).Ptr)
	return C.GoString(ret)
}

/**
 * Get path list pref value
 *
 * @param name The name of the pref
 * @return The value of the pref
 */
func PrefsGetPathList(name string) *GList {
	ret := C.purple_prefs_get_path_list(CCString(name).Ptr)
	return newGListFrom(ret)
}

/**
 * Returns a list of children for a pref
 *
 * @param name The parent pref
 * @return A list of newly allocated strings denoting the names of the children.
 *         Returns @c NULL if there are no children or if pref doesn't exist.
 *         The caller must free all the strings and the list.
 *
 * @since 2.1.0
 */
func PrefsGetChildrenNames(name string) *GList {
	ret := C.purple_prefs_get_children_names(CCString(name).Ptr)
	return newGListFrom(ret)
}

/**
 * Add a callback to a pref (and its children)
 *
 * @param handle   The handle of the receiver.
 * @param name     The name of the preference
 * @param cb       The callback function
 * @param data     The data to pass to the callback function.
 *
 * @return An id to disconnect the callback
 *
 * @see purple_prefs_disconnect_callback
 */
func PrefsConnectCallback(handle unsafe.Pointer, name string,
	cb C.PurplePrefCallback, data unsafe.Pointer) uint {
	ret := C.purple_prefs_connect_callback(handle, CCString(name).Ptr, cb, data)
	return uint(ret)
}

/**
 * Remove a callback to a pref
 */
func PrefsDisconnectCallback(callback_id uint) {
	C.purple_prefs_disconnect_callback((C.guint)(callback_id))
}

/**
 * Remove all pref callbacks by handle
 */
func PrefsDisconnectByHandle(handle unsafe.Pointer) {
	C.purple_prefs_disconnect_by_handle(handle)
}

/**
 * Trigger callbacks as if the pref changed
 */
func PrefsTriggerCallback(name string) {
	C.purple_prefs_trigger_callback(CCString(name).Ptr)
}

/**
 * Trigger callbacks as if the pref changed, taking a #PurplePrefCallbackData
 * instead of a name
 *
 * @since 2.11.0
 */
func PrefsTriggerCallbackObject(data *C.PurplePrefCallbackData) {
	C.purple_prefs_trigger_callback_object(data)
}

/**
 * Read preferences
 */
func PrefsLoad() bool {
	ret := C.purple_prefs_load()
	return c2goBool(ret)
}

/**
 * Rename legacy prefs and delete some that no longer exist.
 */
func PrefsUpdateOld() {
	C.purple_prefs_update_old()
}
