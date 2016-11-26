package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type SavedStatus struct {
	// private
	ss *C.PurpleSavedStatus
}

func newSavedStatusFrom(ss *C.PurpleSavedStatus) *SavedStatus {
	this := &SavedStatus{}
	this.ss = ss
	return this
}

func SavedStatusNew(title string, stype int) *SavedStatus {
	return newSavedStatusFrom(
		C.purple_savedstatus_new(C.CString(title), C.PurpleStatusPrimitive(stype)))
}
func (this *SavedStatus) SetTitle(title string) {
	C.purple_savedstatus_set_title(this.ss, C.CString(title))
}
func (this *SavedStatus) SetType(stype int) {
	C.purple_savedstatus_set_type(this.ss, C.PurpleStatusPrimitive(stype))
}
func (this *SavedStatus) SetMessage(message string) {
	C.purple_savedstatus_set_message(this.ss, C.CString(message))
}

func SavedStatusFindTransientByTypeAndMessage(stype int, message string) *SavedStatus {
	ss := C.purple_savedstatus_find_transient_by_type_and_message(
		C.PurpleStatusPrimitive(stype), C.CString(message))
	return newSavedStatusFrom(ss)
}

func (this *SavedStatus) Activate() {
	C.purple_savedstatus_activate(this.ss)
}
func (this *SavedStatus) ActivateForAccount(ac *Account) {
	C.purple_savedstatus_activate_for_account(this.ss, ac.account)
}
