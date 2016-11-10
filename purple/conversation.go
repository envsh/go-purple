package purple

/*
#include <libpurple/purple.h>
*/
import "C"

const (
	NOT_TYPING = int(C.PURPLE_NOT_TYPING)
	TYPING     = int(C.PURPLE_TYPING)
	TYPED      = int(C.PURPLE_TYPED)
)

type Conversation struct {
	conv *C.PurpleConversation
}

func newConversationFrom(conv *C.PurpleConversation) *Conversation {
	this := &Conversation{}
	this.conv = conv
	return this
}
