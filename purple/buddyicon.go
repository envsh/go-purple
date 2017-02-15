package purple

/*
#include <libpurple/purple.h>
*/
import "C"

// import "unsafe"

type BuddyIcon struct {
	icon *C.PurpleBuddyIcon
}

func newBuddyIconFrom(icon *C.PurpleBuddyIcon) *BuddyIcon { return &BuddyIcon{icon} }

func (this *Account) BuddyIconsSetForUser(username string, icon_data []byte) {
	icon_data_ := C.CBytes(icon_data) // purple's ownership
	C.purple_buddy_icons_set_for_user(this.account, CCString(username).Ptr,
		icon_data_, C.size_t(len(icon_data)), nil)
}
