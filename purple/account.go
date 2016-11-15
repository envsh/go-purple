package purple

/*
// core.c encapse libpurple's core init

#include <libpurple/purple.h>
*/
import "C"

type Account struct {
	account *C.PurpleAccount
}

func newAccountWrapper(acc *C.PurpleAccount) *Account {
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

func (this *Account) SetEnabled(enable bool) {
	C.purple_account_set_enabled(this.account, C.CString(UI_ID), C.TRUE)
}

func (this *Account) GetString(name string) string {
	cstr := C.purple_account_get_string(this.account, C.CString(name), nil)
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

func (this *Account) GetConnection() *Connection {
	cgc := C.purple_account_get_connection(this.account)
	if cgc == nil {
		return nil
	}
	return newConnectWrapper(cgc)
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
	return newBuddyWrapper(buddy)
}

func (this *Account) RequestAdd(name string) {
	C.purple_account_request_add(this.account, C.CString(name), nil, nil, nil)
}

// accounts
func (this *Account) AccountsAdd() {
}
func (this *Account) AccountsRemove() {
}
func (this *Account) AccountsDelete() {
}
func (this *Account) AccountsReorder(newIndex int) {
}
func AccountsGetAll() []*Account {
	return nil
}
func AccountsGetAllActive() []*Account {
	return nil
}
func AccountsFind(name, protocol string) *Account {
	return nil
}
func AccountsRestoreCurrentStatues() {
	C.purple_accounts_restore_current_statuses()
}
