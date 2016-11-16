package purple

/*
// core.c encapse libpurple's core init

#include <libpurple/purple.h>
*/
import "C"

type Account struct {
	account *C.PurpleAccount
}

func newAccountFrom(acc *C.PurpleAccount) *Account {
	this := &Account{}
	this.account = acc
	return this
}

func NewAccountCreate(account string, protocol string, password string) *Account {
	this := &Account{}

	acc := C.purple_account_new(C.CString(account), C.CString(protocol))
	C.purple_account_set_password(acc, C.CString(password))
	C.purple_account_set_remember_password(acc, C.TRUE)
	C.purple_account_set_enabled(acc, C.CString(UI_ID), C.TRUE)
	C.purple_accounts_add(acc)

	this.account = acc
	return this
}

func NewAccount(args ...interface{}) *Account {
	return nil
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
	buddy := C.purple_find_buddy(this.account, C.CString(name))
	if buddy == nil {
		return nil
	}
	return newBuddyFrom(buddy)
}

func (this *Account) RequestAdd(name string) {
	C.purple_account_request_add(this.account, C.CString(name), nil, nil, nil)
}

func (this *Account) SetEnabled(enable bool) {
	C.purple_account_set_enabled(this.account, C.CString(UI_ID), C.TRUE)
}
func (this *Account) GetEnabled() bool {
	rc := C.purple_account_get_enabled(this.account, C.CString(UI_ID))
	if rc == C.TRUE {
		return true
	}
	return false
}
func (this *Account) GetString(name string) string {
	cstr := C.purple_account_get_string(this.account, C.CString(name), nil)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}
func (this *Account) SetString(name string, value string) {
	C.purple_account_set_string(this.account, C.CString(name), C.CString(value))
}
func (this *Account) SetInt(name string, value int) {
	C.purple_account_set_int(this.account, C.CString(name), C.int(value))
}
func (this *Account) GetInt(name string) int {
	rc := C.purple_account_get_int(this.account, C.CString(name), C.int(0))
	return int(rc)
}
func (this *Account) SetBool(name string, value bool) {
	if value {
		C.purple_account_set_bool(this.account, C.CString(name), C.TRUE)
	} else {
		C.purple_account_set_bool(this.account, C.CString(name), C.FALSE)
	}
}
func (this *Account) GetBool(name string) bool {
	rc := C.purple_account_get_bool(this.account, C.CString(name), C.FALSE)
	if rc == C.TRUE {
		return true
	}
	return false
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
	acc := C.purple_accounts_find(C.CString(name), C.CString(protocol))
	if acc == nil {
	} else {
		return newAccountFrom(acc)
	}
	return nil
}
func AccountsRestoreCurrentStatues() {
	C.purple_accounts_restore_current_statuses()
}
