#ifndef _CORE_HELPER_H_
#define _CORE_HELPER_H_

#include <libpurple/purple.h>


void gopurple_connect_progress(PurpleConnection *gc,
                         char *text, size_t step,size_t step_count);
void gopurple_connected(PurpleConnection *gc);
void gopurple_disconnected(PurpleConnection *gc);
void gopurple_notice(PurpleConnection *gc, char *text);
void gopurple_report_disconnect(PurpleConnection *gc, char *text);
void gopurple_network_connected(void);
void gopurple_network_disconnected(void);
void gopurple_report_disconnect_reason(PurpleConnection *gc,
                                 PurpleConnectionError reason,
                                 char *text);

void *gopurple_request_authorize(PurpleAccount *account, char *remote_user, char *id, char *alias, char *message, gboolean on_list,
                        PurpleAccountRequestAuthorizationCb authorize_cb, PurpleAccountRequestAuthorizationCb deny_cb, void *user_data);

// signals
void gopurple_signed_on(PurpleConnection *gc, void *data);
void gopurple_signed_off(PurpleConnection *gc, void *data);
void gopurple_buddy_signed_on(PurpleBuddy *buddy);
void gopurple_buddy_signed_off(PurpleBuddy *buddy);

void gopurple_received_im_msg(PurpleAccount *account, char *sender, char *message, PurpleConversation *conv, PurpleMessageFlags flags);
void gopurple_received_chat_msg(PurpleAccount *account, char *sender, char *message, PurpleConversation *conv, PurpleMessageFlags flags, void *data);
void gopurple_chat_joined(PurpleConversation *conv);
void gopurple_chat_left(PurpleConversation *conv);

#endif
