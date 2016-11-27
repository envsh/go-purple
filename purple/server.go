package purple

/*
#include <libpurple/purple.h>

*/
import "C"

import (
	"time"
)

////
func (this *Connection) ServChatInvite(id int, name string, who string) {
	C.serv_chat_invite(this.conn, C.int(id), CCString(name).Ptr, CCString(who).Ptr)
}

func (this *Connection) ServChatLeave(id int) {
	C.serv_chat_leave(this.conn, C.int(id))
}

func (this *Connection) ServChatWhisper(id int, name string, who string) {
	C.serv_chat_whisper(this.conn, C.int(id), CCString(name).Ptr, CCString(who).Ptr)
}

func (this *Connection) ServChatSend(id int, msg string, flags int) int {
	rc := C.serv_chat_send(this.conn, C.int(id), CCString(msg).Ptr, C.PurpleMessageFlags(flags))
	return int(rc)
}

func (this *Buddy) ServAliasBuddy() {
	C.serv_alias_buddy(this.buddy)
}

func (this *Connection) ServGotAlias(who string, alias string) {
	C.serv_got_alias(this.conn, CCString(who).Ptr, CCString(alias).Ptr)
}

////
func (this *Connection) ServGotIM(who string, msg string, mtype int) {
	samsg := C.strdup(CCString(msg).Ptr)
	C.serv_got_im(this.conn, CCString(who).Ptr, samsg,
		C.PurpleMessageFlags(mtype), C.time_t(time.Now().Unix()))
}

func (this *Connection) ServGotTyping(name string, timeout int, state int) {
	C.serv_got_typing(this.conn, CCString(name).Ptr, C.int(timeout), C.PurpleTypingState(state))
}

func (this *Connection) ServGotTypingStopped(name string) {
	C.serv_got_typing_stopped(this.conn, CCString(name).Ptr)
}

func (this *Connection) ServJoinChat(data *GHashTable) {
	C.serv_join_chat(this.conn, data.ht)
}

func (this *Connection) ServRejectChat(data *GHashTable) {
	C.serv_reject_chat(this.conn, data.ht)
}

func (this *Connection) ServGotChatInvite(name string, who string, msg string, data *GHashTable) {
	C.serv_got_chat_invite(this.conn, CCString(name).Ptr, CCString(who).Ptr,
		CCString(msg).Ptr, data.ht)
}

func (this *Connection) ServGotJoinedChat(id int, name string) *Conversation {
	conv := C.serv_got_joined_chat(this.conn, C.int(id), CCString(name).Ptr)
	return newConversationFrom(conv)
}

func (this *Connection) ServGotJoinChatFailed(data *GHashTable) {
	C.purple_serv_got_join_chat_failed(this.conn, data.ht)
}

func (this *Connection) ServGotChatLeft(id int) {
	C.serv_got_chat_left(this.conn, C.int(id))
}

func (this *Connection) ServGotChatIn(id int, who string, flags int, msg string) {
	C.serv_got_chat_in(this.conn, C.int(id), CCString(who).Ptr,
		C.PurpleMessageFlags(flags), CCString(msg).Ptr,
		C.time_t(time.Now().Unix()))
}
