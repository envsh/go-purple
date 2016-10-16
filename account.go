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
