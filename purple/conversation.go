package purple

/*
#include <libpurple/purple.h>
*/
import "C"
import _ "log"

const (
	NOT_TYPING = int(C.PURPLE_NOT_TYPING)
	TYPING     = int(C.PURPLE_TYPING)
	TYPED      = int(C.PURPLE_TYPED)
)

type Conversation struct {
	conv *C.PurpleConversation
}

type ConvChat struct {
	chat *C.PurpleConvChat
}
type ConvChatBuddy struct {
	buddy *C.PurpleConvChatBuddy
}

func newConversationFrom(conv *C.PurpleConversation) *Conversation {
	this := &Conversation{}
	this.conv = conv
	return this
}

func newConvChatFrom(chat *C.PurpleConvChat) *ConvChat {
	this := &ConvChat{}
	this.chat = chat
	return this
}

func newConvChatBuddyFrom(buddy *C.PurpleConvChatBuddy) *ConvChatBuddy {
	this := &ConvChatBuddy{}
	this.buddy = buddy
	return this
}

func (this *Conversation) GetChatData() *ConvChat {
	chat := C.purple_conversation_get_chat_data(this.conv)
	return newConvChatFrom(chat)
}

func (this *ConvChat) AddUser(user string) {
	C.purple_conv_chat_add_user(this.chat, C.CString(user), nil, 0, C.TRUE)
}

func (this *ConvChat) RenameUser(oldUser string, newUser string) {
	C.purple_conv_chat_rename_user(this.chat, C.CString(oldUser), C.CString(newUser))
}

func (this *ConvChat) RemoveUser(user string) {
	reason := "hehe"
	C.purple_conv_chat_remove_user(this.chat, C.CString(user), C.CString(reason))
}

func (this *ConvChat) GetUsers() []string {
	lst := C.purple_conv_chat_get_users(this.chat)

	res := make([]string, 0)
	newGListFrom(lst).Each(func(item C.gpointer) interface{} {
		ccbudy := (*C.PurpleConvChatBuddy)(item)
		goccbudy := newConvChatBuddyFrom(ccbudy)
		str := goccbudy.GetName()
		res = append(res, str)
		return str
	})

	return res
}

func (this *ConvChatBuddy) GetName() string {
	return C.GoString(this.buddy.name)
}

func (this *ConvChatBuddy) GetAlias() string {
	return C.GoString(this.buddy.alias)
}
