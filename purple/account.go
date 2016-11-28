package purple

/*
#include <libpurple/purple.h>
*/
import "C"
import "unsafe"

type Account struct {
	// private
	account *C.PurpleAccount
}

func newAccountFrom(acc *C.PurpleAccount) *Account {
	this := &Account{}
	this.account = acc
	return this
}

func NewAccountCreate(account string, protocol string, password string) *Account {
	this := &Account{}

	acc := C.purple_account_new(CCString(account).Ptr, CCString(protocol).Ptr)
	C.purple_account_set_password(acc, CCString(password).Ptr)
	C.purple_account_set_remember_password(acc, C.TRUE)
	C.purple_account_set_enabled(acc, CCString(UI_ID).Ptr, C.TRUE)
	C.purple_accounts_add(acc)

	this.account = acc
	return this
}

func NewAccount(args ...interface{}) *Account {
	return nil
}

func (this *Account) Destroy() {
	C.purple_account_destroy(this.account)
	this.account = nil
}

func (this *Account) Connect() {
	C.purple_account_connect(this.account)
}

func (this *Account) GetConnection() *Connection {
	cgc := C.purple_account_get_connection(this.account)
	if cgc == nil {
		return nil
	}
	return newConnectionFrom(cgc)
}

func (this *Account) GetUserName() string {
	cstr := C.purple_account_get_username(this.account)
	return C.GoString(cstr)
}

func (this *Account) GetAlias() string {
	cstr := C.purple_account_get_alias(this.account)
	return C.GoString(cstr)
}

func (this *Account) AddBuddy(b *Buddy) {
	C.purple_account_add_buddy(this.account, b.buddy)
}

func (this *Account) FindBuddy(name string) *Buddy {
	buddy := C.purple_find_buddy(this.account, CCString(name).Ptr)
	if buddy == nil {
		return nil
	}
	return newBuddyFrom(buddy)
}

func (this *Account) FindBuddies(name string) []*Buddy {
	var slst *C.GSList
	if name == "" {
		slst = C.purple_find_buddies(this.account, nil)
	} else {
		slst = C.purple_find_buddies(this.account, CCString(name).Ptr)
	}

	if slst != nil {
		tlen := C.g_slist_length(slst)
		buddies := make([]*Buddy, tlen)
		for idx := 0; idx < int(tlen); idx++ {
			cbuddy := C.g_slist_nth_data(slst, C.guint(idx))
			buddies[idx] = newBuddyFrom((*C.PurpleBuddy)(cbuddy))
		}
		C.g_slist_free(slst)
		return buddies
	}

	return nil
}

func (this *Account) SetEnabled(enable bool) {
	C.purple_account_set_enabled(this.account, CCString(UI_ID).Ptr, C.TRUE)
}
func (this *Account) GetEnabled() bool {
	rc := C.purple_account_get_enabled(this.account, CCString(UI_ID).Ptr)
	if rc == C.TRUE {
		return true
	}
	return false
}
func (this *Account) GetString(name string) string {
	cstr := C.purple_account_get_string(this.account, CCString(name).Ptr, nil)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}
func (this *Account) SetString(name string, value string) {
	C.purple_account_set_string(this.account, CCString(name).Ptr, CCString(value).Ptr)
}
func (this *Account) SetInt(name string, value int) {
	C.purple_account_set_int(this.account, CCString(name).Ptr, C.int(value))
}
func (this *Account) GetInt(name string) int {
	rc := C.purple_account_get_int(this.account, CCString(name).Ptr, C.int(0))
	return int(rc)
}
func (this *Account) SetBool(name string, value bool) {
	if value {
		C.purple_account_set_bool(this.account, CCString(name).Ptr, C.TRUE)
	} else {
		C.purple_account_set_bool(this.account, CCString(name).Ptr, C.FALSE)
	}
}
func (this *Account) GetBool(name string) bool {
	rc := C.purple_account_get_bool(this.account, CCString(name).Ptr, C.FALSE)
	if rc == C.TRUE {
		return true
	}
	return false
}

func (this *Account) SetUserName(name string) {
	C.purple_account_set_username(this.account, CCString(name).Ptr)
}
func (this *Account) SetAlias(alias string) {
	C.purple_account_set_alias(this.account, CCString(alias).Ptr)
}
func (this *Account) SetPassword(password string) {
	C.purple_account_set_password(this.account, CCString(password).Ptr)
}
func (this *Account) SetUserInfo(userInfo string) {
	C.purple_account_set_user_info(this.account, CCString(userInfo).Ptr)
}
func (this *Account) SetBuddyIconPath(path string) {
	C.purple_account_set_buddy_icon_path(this.account, CCString(path).Ptr)
}

func (this *Account) GetProtocolId() string {
	id := C.purple_account_get_protocol_id(this.account)
	return C.GoString(id)
}
func (this *Account) GetProtocolName() string {
	name := C.purple_account_get_protocol_name(this.account)
	return C.GoString(name)
}

func (this *Account) GetLog(create bool) *Log {
	log := C.purple_account_get_log(this.account, go2cBool(create))
	return newLogFrom(log)
}
func (this *Account) DestroyLog() {
	C.purple_account_destroy_log(this.account)
}
func (this *Account) GetCurrentError() *ConnectionErrorInfo {
	return newConnectionErrorInfoFrom(C.purple_account_get_current_error(this.account))
}
func (this *Account) ClearCurrentError() {
	C.purple_account_clear_current_error(this.account)
}

func (this *Account) NotifyAdded(remoteUser, id, alias, msg string) {
	C.purple_account_notify_added(this.account, CCString(remoteUser).Ptr,
		CCString(id).Ptr, CCString(alias).Ptr, CCString(msg).Ptr)
}
func (this *Account) RequestAdd(remoteUser, id, alias, msg string) {
	C.purple_account_request_add(this.account, CCString(remoteUser).Ptr,
		CCString(id).Ptr, CCString(alias).Ptr, CCString(msg).Ptr)
}
func (this *Account) RequestCloseWithAccount() {
	C.purple_account_request_close_with_account(this.account)
}

func (this *Account) SetProtocolId(protocolId string) {
	C.purple_account_set_protocol_id(this.account, CCString(protocolId).Ptr)
}
func (this *Account) SetConnection(gc *Connection) {
	C.purple_account_set_connection(this.account, gc.conn)
}
func (this *Account) SetRememberPassword(value bool) {
	C.purple_account_set_remember_password(this.account, go2cBool(value))
}
func (this *Account) SetCheckMail(value bool) {
	C.purple_account_set_check_mail(this.account, go2cBool(value))
}

// accounts
func (this *Account) AccountsAdd() {
	C.purple_accounts_add(this.account)
}
func (this *Account) AccountsRemove() {
	C.purple_accounts_remove(this.account)
}
func (this *Account) AccountsDelete() {
	C.purple_accounts_delete(this.account)
}
func (this *Account) AccountsReorder(newIndex int) {
	C.purple_accounts_reorder(this.account, C.gint(newIndex))
}

func AccountsGetAll() []*Account {
	acs := make([]*Account, 0)
	lst := C.purple_accounts_get_all()
	newGListFrom(lst).Each(func(item C.gpointer) {
		ac := newAccountFrom((*C.PurpleAccount)(item))
		acs = append(acs, ac)
	})
	return acs
}
func AccountsGetAllActive() []*Account {
	acs := make([]*Account, 0)
	lst := C.purple_accounts_get_all_active()
	newGListFrom(lst).Each(func(item C.gpointer) {
		ac := newAccountFrom((*C.PurpleAccount)(item))
		acs = append(acs, ac)
	})
	return acs
}
func AccountsFind(name, protocol string) *Account {
	acc := C.purple_accounts_find(CCString(name).Ptr, CCString(protocol).Ptr)
	if acc == nil {
	} else {
		return newAccountFrom(acc)
	}
	return nil
}
func AccountsRestoreCurrentStatues() {
	C.purple_accounts_restore_current_statuses()
}

func AccountsGetHandle() unsafe.Pointer {
	return C.purple_accounts_get_handle()
}
