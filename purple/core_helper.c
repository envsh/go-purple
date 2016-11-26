#include "_cgo_export.h"

#include <libpurple/purple.h>

#include <glib.h>
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <string.h>
#include <unistd.h>

#define CUSTOM_USER_DIRECTORY  ".gopurple"
#define CUSTOM_PLUGIN_PATH     ""
#define PLUGIN_SAVE_PREF       "/purple/purple-daemon/plugins/saved"
#define UI_ID                  "gopurple-daemon"

/// plugin preference names /
#define MW_PRPL_OPT_BASE          "/plugins/prpl/meanwhile"
#define MW_PRPL_OPT_BLIST_ACTION  MW_PRPL_OPT_BASE "/blist_action"
#define MW_PRPL_OPT_PSYCHIC       MW_PRPL_OPT_BASE "/psychic"
#define MW_PRPL_OPT_FORCE_LOGIN   MW_PRPL_OPT_BASE "/force_login"
#define MW_PRPL_OPT_SAVE_DYNAMIC  MW_PRPL_OPT_BASE "/save_dynamic"

/// keys to get/set purple plugin information
#define MW_KEY_HOST        "server"
#define MW_KEY_PORT        "port"
#define MW_KEY_FORCE       "force_login"
#define MW_KEY_FAKE_IT     "fake_client_id"
#define MW_KEY_CLIENT      "client_id_val"
#define MW_KEY_MAJOR       "client_major"
#define MW_KEY_MINOR       "client_minor"

///** the possible meanwhile buddy list storage settings
enum blist_choice {
  blist_choice_LOCAL = 1, //**< local only
  blist_choice_MERGE = 2, //**< merge from server
  blist_choice_STORE = 3, //**< merge from and save to server
  blist_choice_SYNCH = 4, //**< sync with server
};


//** the default blist storage option
#define BLIST_CHOICE_DEFAULT  blist_choice_MERGE

#define BLIST_PREF_SET(n) purple_prefs_set_int(MW_PRPL_OPT_BLIST_ACTION, n)

#define PURPLE_GLIB_READ_COND  (G_IO_IN  | G_IO_HUP | G_IO_ERR)
#define PURPLE_GLIB_WRITE_COND (G_IO_OUT | G_IO_HUP | G_IO_ERR | G_IO_NVAL)

typedef struct _PurpleGLibIOClosure {
	PurpleInputFunction function;
	guint result;
	gpointer data;
} PurpleGLibIOClosure;

static void purple_glib_io_destroy(gpointer data)
{
	g_free(data);
}

static gboolean purple_glib_io_invoke(GIOChannel *source, GIOCondition condition, gpointer data)
{
	PurpleGLibIOClosure *closure = data;
	PurpleInputCondition purple_cond = 0;

	if (condition & PURPLE_GLIB_READ_COND)
		purple_cond |= PURPLE_INPUT_READ;
	if (condition & PURPLE_GLIB_WRITE_COND)
		purple_cond |= PURPLE_INPUT_WRITE;

	closure->function(closure->data, g_io_channel_unix_get_fd(source),
			  purple_cond);

	return TRUE;
}

static guint glib_input_add(gint fd, PurpleInputCondition condition, PurpleInputFunction function,
							   gpointer data)
{
	PurpleGLibIOClosure *closure = g_new0(PurpleGLibIOClosure, 1);
	GIOChannel *channel;
	GIOCondition cond = 0;

	closure->function = function;
	closure->data = data;

	if (condition & PURPLE_INPUT_READ)
		cond |= PURPLE_GLIB_READ_COND;
	if (condition & PURPLE_INPUT_WRITE)
		cond |= PURPLE_GLIB_WRITE_COND;

	channel = g_io_channel_unix_new(fd);
	closure->result = g_io_add_watch_full(channel, G_PRIORITY_DEFAULT, cond,
					      purple_glib_io_invoke, closure, purple_glib_io_destroy);

	g_io_channel_unref(channel);
	return closure->result;
}

static PurpleEventLoopUiOps glib_eventloops =
{
	g_timeout_add,
	g_source_remove,
	glib_input_add,
	g_source_remove,
	NULL,
#if GLIB_CHECK_VERSION(2,14,0)
	g_timeout_add_seconds,
#else
	NULL,
#endif
//	/* padding
	NULL,
	NULL,
	NULL
};
PurpleEventLoopUiOps *gopurple_get_loopops() {
    return &glib_eventloops;
}
//*** End of the eventloop functions.


///////////
#include "core_helper.h"
void gopurple_connect_to_signals(void)
{
	static int handle;

	purple_signal_connect(purple_connections_get_handle(), "signed-on", &handle,
                          PURPLE_CALLBACK(gopurple_signed_on), (void*)&handle);
    /*
	purple_signal_connect(purple_blist_get_handle(), "buddy-signed-on", &handle,
				PURPLE_CALLBACK(buddy_signed_on), NULL);

	purple_signal_connect(purple_blist_get_handle(), "buddy-signed-off", &handle,
				PURPLE_CALLBACK(buddy_signed_off), NULL);

	purple_signal_connect(purple_blist_get_handle(), "buddy-status-changed", &handle,
				PURPLE_CALLBACK(buddy_away), NULL);

	purple_signal_connect(purple_blist_get_handle(), "buddy-idle-changed", &handle,
				PURPLE_CALLBACK(buddy_idle), NULL);

	purple_signal_connect(purple_conversations_get_handle(), "received-im-msg", &handle,
				PURPLE_CALLBACK(received_im_msg), NULL);

	purple_signal_connect(purple_conversations_get_handle(), "buddy-typing", &handle,
				PURPLE_CALLBACK(buddy_typing), NULL);

	purple_signal_connect(purple_conversations_get_handle(), "buddy-typed", &handle,
				PURPLE_CALLBACK(buddy_typed), NULL);

	purple_signal_connect(purple_conversations_get_handle(), "buddy-typing-stopped", &handle,
				PURPLE_CALLBACK(buddy_typing_stopped), NULL);

	purple_signal_connect(purple_accounts_get_handle(), "account-authorization-requested", &handle,
				PURPLE_CALLBACK(account_authorization_requested), NULL);

	// purple_signal_connect(purple_accounts_get_handle(), "dbus-method-called", &handle,
				// PURPLE_CALLBACK(dbus_method_called), NULL);

	purple_signal_connect(purple_conversations_get_handle(), "received-chat-msg", &handle,
                          PURPLE_CALLBACK(received_chat_msg), NULL);
    */
}
