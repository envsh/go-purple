package purple

/*
#include <libpurple/purple.h>

static void gopurple_account_set_status(PurpleAccount *account, const char *status_id,
	gboolean active) {
    purple_account_set_status(account, status_id, active, NULL);
}
*/
import "C"

type AccountOption struct {
	// private
	ao    *C.PurpleAccountOption
	fromc bool
}

type AccountUserSplit struct {
	// private
	aus   *C.PurpleAccountUserSplit
	fromc bool
}

func newAccountOptionFrom(ao *C.PurpleAccountOption) *AccountOption {
	this := &AccountOption{}
	this.ao = ao
	this.fromc = true
	return this
}

func newAccountUserSplitFrom(aus *C.PurpleAccountUserSplit) *AccountUserSplit {
	this := &AccountUserSplit{}
	this.aus = aus
	this.fromc = true
	return this
}

func NewAccountOption(otype PrefType, text, prefName string) *AccountOption {
	r := C.purple_account_option_new(C.PurplePrefType(otype),
		CCString(text).Ptr, CCString(prefName).Ptr)
	return newAccountOptionFrom(r)
}

/**
 * Creates a new boolean account option.
 *
 * @param text          The text of the option.
 * @param pref_name     The account preference name for the option.
 * @param default_value The default value.
 *
 * @return The account option.
 */
func NewAccountOptionBool(text, prefName string, defaultValue bool) *AccountOption {
	r := C.purple_account_option_bool_new(CCString(text).Ptr,
		CCString(prefName).Ptr, go2cBool(defaultValue))
	return newAccountOptionFrom(r)
}

/**
 * Creates a new integer account option.
 *
 * @param text          The text of the option.
 * @param pref_name     The account preference name for the option.
 * @param default_value The default value.
 *
 * @return The account option.
 */
func NewAccountOptionInt(text, prefName string, defaultValue int) *AccountOption {
	r := C.purple_account_option_int_new(CCString(text).Ptr,
		CCString(prefName).Ptr, C.int(defaultValue))
	return newAccountOptionFrom(r)
}

/**
 * Creates a new string account option.
 *
 * @param text          The text of the option.
 * @param pref_name     The account preference name for the option.
 * @param default_value The default value.
 *
 * @return The account option.
 */
func NewAccountOptionString(text, prefName string, defaultValue string) *AccountOption {
	r := C.purple_account_option_string_new(CCString(text).Ptr,
		CCString(prefName).Ptr, CCString(defaultValue).Ptr)
	return newAccountOptionFrom(r)
}

/**
 * Creates a new list account option.
 *
 * The list passed will be owned by the account option, and the
 * strings inside will be freed automatically.
 *
 * The list is a list of #PurpleKeyValuePair items. The key is the label that
 * should be displayed to the user, and the <tt>(const char *)</tt> value is
 * the internal ID that should be passed to purple_account_set_string() to
 * choose that value.
 *
 * @param text      The text of the option.
 * @param pref_name The account preference name for the option.
 * @param list      The key, value list.
 *
 * @return The account option.
 */
func NewAccountOptionList(text, prefName string, list *GList) *AccountOption {
	r := C.purple_account_option_list_new(CCString(text).Ptr,
		CCString(prefName).Ptr, list.lst)
	return newAccountOptionFrom(r)
}

/**
 * Destroys an account option.
 *
 * @param option The option to destroy.
 */
func (this *AccountOption) Destroy() {
	C.purple_account_option_destroy(this.ao)
}

/**
 * Sets the default boolean value for an account option.
 *
 * @param option The account option.
 * @param value  The default boolean value.
 */
func (this *AccountOption) SetDefaultBool(value bool) {
	C.purple_account_option_set_default_bool(this.ao, go2cBool(value))
}

/**
 * Sets the default integer value for an account option.
 *
 * @param option The account option.
 * @param value  The default integer value.
 */
func (this *AccountOption) SetDefaultInt(value int) {
	C.purple_account_option_set_default_int(this.ao, C.int(value))
}

/**
 * Sets the default string value for an account option.
 *
 * @param option The account option.
 * @param value  The default string value.
 */
func (this *AccountOption) SetDefaultString(value string) {
	C.purple_account_option_set_default_string(this.ao, CCString(value).Ptr)
}

/**
 * Sets the masking for an account option. Setting this to %TRUE acts
 * as a hint to the UI that the option's value should be obscured from
 * view, like a password.
 *
 * @param option The account option.
 * @param masked The masking.
 */
func (this *AccountOption) SetMasked(masked bool) {
	C.purple_account_option_set_masked(this.ao, go2cBool(masked))
}

/**
 * Sets the list values for an account option.
 *
 * The list passed will be owned by the account option, and the
 * strings inside will be freed automatically.
 *
 * The list is in key, value pairs. The key is the ID stored and used
 * internally, and the value is the label displayed.
 *
 * @param option The account option.
 * @param values The default list value.
 */
func (this *AccountOption) SetList(values *GList) {
	C.purple_account_option_set_list(this.ao, values.lst)
}

/**
 * Adds an item to a list account option.
 *
 * @param option The account option.
 * @param key    The key.
 * @param value  The value.
 */
func (this *AccountOption) AddListItem(key, value string) {
	C.purple_account_option_add_list_item(this.ao, CCString(key).Ptr, CCString(value).Ptr)
}

/**
 * Returns the specified account option's type.
 *
 * @param option The account option.
 *
 * @return The account option's type.
 */
func (this *AccountOption) GetType() PrefType {
	r := C.purple_account_option_get_type(this.ao)
	return PrefType(r)
}

/**
 * Returns the text for an account option.
 *
 * @param option The account option.
 *
 * @return The account option's text.
 */
func (this *AccountOption) GetText() string {
	r := C.purple_account_option_get_text(this.ao)
	return C.GoString(r)
}

/**
 * Returns the name of an account option.  This corresponds to the @c pref_name
 * parameter supplied to purple_account_option_new() or one of the
 * type-specific constructors.
 *
 * @param option The account option.
 *
 * @return The option's name.
 */
func (this *AccountOption) getSetting() string {
	r := C.purple_account_option_get_setting(this.ao)
	return C.GoString(r)
}

/**
 * Returns the default boolean value for an account option.
 *
 * @param option The account option.
 *
 * @return The default boolean value.
 */
func (this *AccountOption) GetDefaultBool() bool {
	r := C.purple_account_option_get_default_bool(this.ao)
	return c2goBool(r)
}

/**
 * Returns the default integer value for an account option.
 *
 * @param option The account option.
 *
 * @return The default integer value.
 */
func (this *AccountOption) GetDefaultInt() int {
	r := C.purple_account_option_get_default_int(this.ao)
	return int(r)
}

/**
 * Returns the default string value for an account option.
 *
 * @param option The account option.
 *
 * @return The default string value.
 */
func (this *AccountOption) GetDefaultString() string {
	r := C.purple_account_option_get_default_string(this.ao)
	return C.GoString(r)
}

/**
 * Returns the default string value for a list account option.
 *
 * @param option The account option.
 *
 * @return The default list string value.
 */
func (this *AccountOption) GetDefaultListValue() string {
	r := C.purple_account_option_get_default_list_value(this.ao)
	return C.GoString(r)
}

/**
 * Returns whether an option's value should be masked from view, like a
 * password.  If so, the UI might display each character of the option
 * as a '*' (for example).
 *
 * @param option The account option.
 *
 * @return %TRUE if the option's value should be obscured.
 */
func (this *AccountOption) GetMasked() bool {
	r := C.purple_account_option_get_masked(this.ao)
	return c2goBool(r)
}

/**
 * Returns the list values for an account option.
 *
 * @param option The account option.
 *
 * @constreturn A list of #PurpleKeyValuePair, mapping the human-readable
 *              description of the value to the <tt>(const char *)</tt> that
 *              should be passed to purple_account_set_string() to set the
 *              option.
 */
func (this *AccountOption) GetList() *GList {
	r := C.purple_account_option_get_list(this.ao)
	return newGListFrom(r)
}

/*@}*/

/**************************************************************************/
/** @name Account User Split API                                          */
/**************************************************************************/
/*@{*/

/**
 * Creates a new account username split.
 *
 * @param text          The text of the option.
 * @param default_value The default value.
 * @param sep           The field separator.
 *
 * @return The new user split.
 */
func NewAccountUserSplit(text, defaultValue string, sep byte) *AccountUserSplit {
	r := C.purple_account_user_split_new(CCString(text).Ptr,
		CCString(defaultValue).Ptr, C.char(sep))
	return newAccountUserSplitFrom(r)
}

/**
 * Destroys an account username split.
 *
 * @param split The split to destroy.
 */
func (this *AccountUserSplit) Destroy() {
	C.purple_account_user_split_destroy(this.aus)
}

/**
 * Returns the text for an account username split.
 *
 * @param split The account username split.
 *
 * @return The account username split's text.
 */
func (this *AccountUserSplit) GetText() string {
	r := C.purple_account_user_split_get_text(this.aus)
	return C.GoString(r)
}

/**
 * Returns the default string value for an account split.
 *
 * @param split The account username split.
 *
 * @return The default string.
 */
func (this *AccountUserSplit) GetDefaultValue() string {
	r := C.purple_account_user_split_get_default_value(this.aus)
	return C.GoString(r)
}

/**
 * Returns the field separator for an account split.
 *
 * @param split The account username split.
 *
 * @return The field separator.
 */
func (this *AccountUserSplit) GetSeparator() byte {
	r := C.purple_account_user_split_get_separator(this.aus)
	return byte(r)
}

/**
 * Returns the 'reverse' value for an account split.
 *
 * @param split The account username split.
 *
 * @return The 'reverse' value.
 */
func (this *AccountUserSplit) GetReverse() bool {
	r := C.purple_account_user_split_get_reverse(this.aus)
	return c2goBool(r)
}

/**
 * Sets the 'reverse' value for an account split.
 *
 * @param split   The account username split.
 * @param reverse The 'reverse' value
 */
func (this *AccountUserSplit) SetReverse(reverse bool) {
	C.purple_account_user_split_set_reverse(this.aus, go2cBool(reverse))
}
