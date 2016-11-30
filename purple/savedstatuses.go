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
		C.purple_savedstatus_new(CCString(title).Ptr, C.PurpleStatusPrimitive(stype)))
}
func (this *SavedStatus) SetTitle(title string) {
	C.purple_savedstatus_set_title(this.ss, CCString(title).Ptr)
}
func (this *SavedStatus) SetType(stype int) {
	C.purple_savedstatus_set_type(this.ss, C.PurpleStatusPrimitive(stype))
}
func (this *SavedStatus) SetMessage(message string) {
	C.purple_savedstatus_set_message(this.ss, CCString(message).Ptr)
}

func SavedStatusFindTransientByTypeAndMessage(stype int, message string) *SavedStatus {
	ss := C.purple_savedstatus_find_transient_by_type_and_message(
		C.PurpleStatusPrimitive(stype), CCString(message).Ptr)
	return newSavedStatusFrom(ss)
}

func (this *SavedStatus) Activate() {
	C.purple_savedstatus_activate(this.ss)
}
func (this *SavedStatus) ActivateForAccount(ac *Account) {
	C.purple_savedstatus_activate_for_account(this.ss, ac.account)
}

func SavedStatusSetIdleAway(idleaway bool) {
	C.purple_savedstatus_set_idleaway(go2cBool(idleaway))
}
