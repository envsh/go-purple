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

func (this *Conversation) SetName(title string) {
	C.purple_conversation_set_name(this.conv, C.CString(title))
}

func (this *Conversation) GetName() string {
	name := C.purple_conversation_get_name(this.conv)
	return C.GoString(name)
}

func (this *Conversation) SetData(key, data string) {
	C.purple_conversation_set_data(this.conv, C.CString(key), C.CString(data))
}

func (this *Conversation) GetData(key string) string {
	data := C.purple_conversation_get_data(this.conv, C.CString(key))
	return C.GoString((*C.char)(data))
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
	newGListFrom(lst).Each(func(item C.gpointer) {
		ccbudy := (*C.PurpleConvChatBuddy)(item)
		goccbudy := newConvChatBuddyFrom(ccbudy)
		str := goccbudy.GetName()
		res = append(res, str)
	})

	return res
}

func (this *ConvChat) Send(message string) {
	C.purple_conv_chat_send(this.chat, C.CString(message))
}

func (this *ConvChat) Write(who string, message string, flags int) {
	// C.purple_conv_chat_write(this.chat, C.CString(message))
}

func (this *ConvChat) SendWithFlag(message string, flags int) {
	C.purple_conv_chat_send_with_flags(this.chat, C.CString(message), C.PurpleMessageFlags(flags))
}

func (this *ConvChat) FindBuddy(name string) *ConvChatBuddy {
	cbbudy := C.purple_conv_chat_cb_find(this.chat, C.CString(name))
	return newConvChatBuddyFrom(cbbudy)
}

func (this *ConvChatBuddy) GetName() string {
	return C.GoString(this.buddy.name)
}

func (this *ConvChatBuddy) GetAlias() string {
	return C.GoString(this.buddy.alias)
}

func (this *ConvChatBuddy) SetAlias(alias string) {
	if this.buddy != nil {
		if this.buddy.alias != nil {
			// TODO free it first
		}
		this.buddy.alias = C.CString(alias)
	}
}

func (this *ConvChatBuddy) Destroy() {
	buddy := this.buddy
	this.buddy = nil
	C.purple_conv_chat_cb_destroy(buddy)
}

func GetConversations() []*Conversation {
	convs := make([]*Conversation, 0)
	lst := C.purple_get_conversations()
	newGListFrom(lst).Each(func(item C.gpointer) {
		conv := (*C.PurpleConversation)(item)
		convs = append(convs, newConversationFrom(conv))
	})
	return convs
}

func GetIMs() []*Conversation {
	return nil
}

func GetChats() []*Conversation {
	convs := make([]*Conversation, 0)
	lst := C.purple_get_chats()
	newGListFrom(lst).Each(func(item C.gpointer) {
		conv := (*C.PurpleConversation)(item)
		convs = append(convs, newConversationFrom(conv))
	})
	return convs
}

func (this *Connection) FindChat(id int) *Conversation {
	conv := C.purple_find_chat(this.conn, C.int(id))
	if conv == nil {
		return nil
	}
	return newConversationFrom(conv)
}
