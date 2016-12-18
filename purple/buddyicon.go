package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type BuddyIcon struct {
	icon *C.PurpleBuddyIcon
}

func newBuddyIconFrom(icon *C.PurpleBuddyIcon) *BuddyIcon { return &BuddyIcon{icon} }
