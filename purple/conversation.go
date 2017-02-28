package purple

/*
#include <libpurple/purple.h>
*/
import "C"
import "unsafe"

import _ "log"
import (
	"time"
)

type TypingState int

const (
	NOT_TYPING = int(C.PURPLE_NOT_TYPING)
	TYPING     = int(C.PURPLE_TYPING)
	TYPED      = int(C.PURPLE_TYPED)
)

type ConvType int

const (
	CONV_TYPE_UNKNOWN = int(C.PURPLE_CONV_TYPE_UNKNOWN) /**< Unknown conversation type. */
	CONV_TYPE_IM      = int(C.PURPLE_CONV_TYPE_IM)      /**< Instant Message.           */
	CONV_TYPE_CHAT    = int(C.PURPLE_CONV_TYPE_CHAT)    /**< Chat room.                 */
	CONV_TYPE_MISC    = int(C.PURPLE_CONV_TYPE_MISC)    /**< A misc. conversation.      */
	CONV_TYPE_ANY     = int(C.PURPLE_CONV_TYPE_ANY)     /**< Any type of conversation.  */
)

type MessageFlags int

const (
	MESSAGE_SEND        MessageFlags = C.PURPLE_MESSAGE_SEND
	MESSAGE_RECV                     = C.PURPLE_MESSAGE_RECV
	MESSAGE_SYSTEM                   = C.PURPLE_MESSAGE_SYSTEM
	MESSAGE_AUTO_RESP                = C.PURPLE_MESSAGE_AUTO_RESP
	MESSAGE_ACTIVE_ONLY              = C.PURPLE_MESSAGE_ACTIVE_ONLY
	MESSAGE_NICK                     = C.PURPLE_MESSAGE_NICK
	MESSAGE_NO_LOG                   = C.PURPLE_MESSAGE_NO_LOG
	MESSAGE_WHISPER                  = C.PURPLE_MESSAGE_WHISPER
	MESSAGE_ERROR                    = C.PURPLE_MESSAGE_ERROR
	MESSAGE_DELAYED                  = C.PURPLE_MESSAGE_DELAYED
	MESSAGE_RAW                      = C.PURPLE_MESSAGE_RAW
	MESSAGE_IMAGES                   = C.PURPLE_MESSAGE_IMAGES
	MESSAGE_NOTIFY                   = C.PURPLE_MESSAGE_NOTIFY
	MESSAGE_NO_LINKIFY               = C.PURPLE_MESSAGE_NO_LINKIFY
	MESSAGE_INVISIBLE                = C.PURPLE_MESSAGE_INVISIBLE
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

type ConvIm struct {
	im *C.PurpleConvIm
}

func newConvImFrom(im *C.PurpleConvIm) *ConvIm { return &ConvIm{im} }

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

func NewConversation(ctype int, ac *Account, name string) *Conversation {
	this := &Conversation{}
	conv := C.purple_conversation_new(C.PurpleConversationType(ctype), ac.account, CCString(name).Ptr)
	this.conv = conv
	return this
}

func (this *Conversation) GetChatData() *ConvChat {
	chat := C.purple_conversation_get_chat_data(this.conv)
	return newConvChatFrom(chat)
}

func (this *Conversation) SetName(title string) {
	C.purple_conversation_set_name(this.conv, CCString(title).Ptr)
}

func (this *Conversation) GetName() string {
	name := C.purple_conversation_get_name(this.conv)
	return C.GoString(name)
}

func (this *Conversation) SetData(key, data string) {
	C.purple_conversation_set_data(this.conv, CCString(key).Ptr, CCString(data).gpointer())
}

func (this *Conversation) GetData(key string) string {
	data := C.purple_conversation_get_data(this.conv, CCString(key).Ptr)
	return C.GoString((*C.char)(data))
}

func (this *Conversation) GetAccount() *Account {
	ac := C.purple_conversation_get_account(this.conv)
	if ac == nil {
		return nil
	}
	return newAccountFrom(ac)
}

func (this *Conversation) GetConnection() *Connection {
	gc := C.purple_conversation_get_gc(this.conv)
	if gc == nil {
		return nil
	}
	return newConnectionFrom(gc)
}

func (this *Conversation) SetLogging(logging bool) {
	C.purple_conversation_set_logging(this.conv, go2cBool(logging))
}
func (this *Conversation) IsLogging() bool {
	return c2goBool(C.purple_conversation_is_logging(this.conv))
}

func (this *Conversation) Destroy() {
	C.purple_conversation_destroy(this.conv)
	this.conv = nil
	this = nil
}

func (this *ConvChat) AddUser(user string) {
	C.purple_conv_chat_add_user(this.chat, CCString(user).Ptr, nil, 0, C.TRUE)
}

func (this *ConvChat) RenameUser(oldUser string, newUser string) {
	C.purple_conv_chat_rename_user(this.chat, CCString(oldUser).Ptr, CCString(newUser).Ptr)
}

func (this *ConvChat) RemoveUser(user string, reason string) {
	C.purple_conv_chat_remove_user(this.chat, CCString(user).Ptr, CCString(reason).Ptr)
}

func (this *ConvChat) GetUsers() []*ConvChatBuddy {
	lst := C.purple_conv_chat_get_users(this.chat)

	res := make([]*ConvChatBuddy, 0)
	newGListFrom(lst).Each(func(item C.gpointer) {
		ccbudy := (*C.PurpleConvChatBuddy)(item)
		goccbudy := newConvChatBuddyFrom(ccbudy)
		res = append(res, goccbudy)
	})

	return res
}

func (this *ConvChat) Send(message string) {
	C.purple_conv_chat_send(this.chat, CCString(message).Ptr)
}

func (this *ConvChat) Write(who string, message string, flags int) {
	C.purple_conv_chat_write(this.chat, CCString(who).Ptr,
		CCString(message).Ptr, C.PurpleMessageFlags(flags), C.time_t(time.Now().Unix()))
}

func (this *ConvChat) SendWithFlag(message string, flags int) {
	C.purple_conv_chat_send_with_flags(this.chat, CCString(message).Ptr, C.PurpleMessageFlags(flags))
}

func (this *ConvChat) FindBuddy(name string) *ConvChatBuddy {
	cbbudy := C.purple_conv_chat_cb_find(this.chat, CCString(name).Ptr)
	return newConvChatBuddyFrom(cbbudy)
}

func (this *ConvChat) SetTopic(who, topic string) {
	C.purple_conv_chat_set_topic(this.chat, CCString(who).Ptr, CCString(topic).Ptr)
}
func (this *ConvChat) GetTopic() string {
	r := C.purple_conv_chat_get_topic(this.chat)
	return C.GoString(r)
}
func (this *ConvChat) SetId(id int) {
	C.purple_conv_chat_set_id(this.chat, C.int(id))
}
func (this *ConvChat) GetId() int {
	r := C.purple_conv_chat_get_id(this.chat)
	return int(r)
}

func (this *ConvChat) HasLeft() bool {
	r := C.purple_conv_chat_has_left(this.chat)
	return c2goBool(r)
}

func NewConvChatBuddy(name, alias string) *ConvChatBuddy {
	r := C.purple_conv_chat_cb_new(C.CString(name), C.CString(alias), 0)
	return newConvChatBuddyFrom(r)
}

func (this *ConvChatBuddy) GetName() string {
	return C.GoString(this.buddy.name)
}

func (this *ConvChatBuddy) GetAlias() string {
	return C.GoString(this.buddy.alias)
}

/* no corresponding method in libpurple
func (this *ConvChatBuddy) SetAlias(alias string) {
	if this.buddy != nil {
		if this.buddy.alias != nil {
			// TODO free it first
		}
		if this.buddy.alias != nil {
			// C.free(unsafe.Pointer(this.buddy.alias))
		}
		this.buddy.alias = C.CString(alias)
	}
}
*/

func (this *ConvChatBuddy) Destroy() {
	buddy := this.buddy
	this.buddy = nil

	C.purple_conv_chat_cb_destroy(buddy)
	this = nil
}

func (this *ConvChatBuddy) SetAttribute(chat *ConvChat, key, value string) {
	C.purple_conv_chat_cb_set_attribute(chat.chat, this.buddy, CCString(key).Ptr, CCString(value).Ptr)
}
func (this *ConvChatBuddy) GetAttribute(key string) string {
	r := C.purple_conv_chat_cb_get_attribute(this.buddy, CCString(key).Ptr)
	return C.GoString(r)
}

// conv im
func (this *ConvIm) GetConversation() *Conversation {
	r := C.purple_conv_im_get_conversation(this.im)
	return newConversationFrom(r)
}

func (this *ConvIm) SetIcon(icon *BuddyIcon) {
	C.purple_conv_im_set_icon(this.im, icon.icon)
}

func (this *ConvIm) GetIcon() *BuddyIcon {
	r := C.purple_conv_im_get_icon(this.im)
	return newBuddyIconFrom(r)
}

func (this *ConvIm) SetTypingState(state int) {
	C.purple_conv_im_set_typing_state(this.im, C.PurpleTypingState(state))
}

func (this *ConvIm) GetTypingState() int {
	r := C.purple_conv_im_get_typing_state(this.im)
	return int(r)
}

func (this *ConvIm) StartTypingTimeout(timeout int) {
	C.purple_conv_im_start_typing_timeout(this.im, C.int(timeout))
}

func (this *ConvIm) StopTypingTimeout() {
	C.purple_conv_im_stop_typing_timeout(this.im)
}

func (this *ConvIm) GetTypingTimeout() uint {
	r := C.purple_conv_im_get_typing_timeout(this.im)
	return uint(r)
}

func (this *ConvIm) SetTypeAgain(val uint) {
	C.purple_conv_im_set_type_again(this.im, C.uint(val))
}

func (this *ConvIm) GetTypeAgain() uint64 {
	r := C.purple_conv_im_get_type_again(this.im)
	return uint64(r)
}

func (this *ConvIm) StartSendTypedTimeout() {
	C.purple_conv_im_start_send_typed_timeout(this.im)
}

func (this *ConvIm) StopSendTypedTimeout() {
	C.purple_conv_im_stop_send_typed_timeout(this.im)
}

func (this *ConvIm) GetSendTypedTimeout() uint {
	r := C.purple_conv_im_get_send_typed_timeout(this.im)
	return uint(r)
}

func (this *ConvIm) UpdateTyping() {
	C.purple_conv_im_update_typing(this.im)
}

func (this *ConvIm) Write(who, message string, flags int) {
	C.purple_conv_im_write(this.im, CCString(who).Ptr,
		CCString(message).Ptr, C.PurpleMessageFlags(flags),
		C.time_t(time.Now().Unix()))
}

func (this *ConvIm) Send(message string) {
	C.purple_conv_im_send(this.im, CCString(message).Ptr)
}

func (this *ConvIm) SendWithFlags(message string, flags int) {
	C.purple_conv_im_send_with_flags(this.im, CCString(message).Ptr, C.PurpleMessageFlags(flags))
}

//
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

func ConversationsGetHandle() unsafe.Pointer {
	return C.purple_conversations_get_handle()
}

func FindConversationWithAccount(ctype int, name string, ac *Account) *Conversation {
	conv := C.purple_find_conversation_with_account(C.PurpleConversationType(ctype),
		CCString(name).Ptr, ac.account)
	if conv == nil {
		return nil
	}
	return newConversationFrom(conv)
}
