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
//*** End of the eventloop functions.



static void *request_authorize(PurpleAccount *account, const char *remote_user, const char *id, const char *alias, const char *message, gboolean on_list,
	PurpleAccountRequestAuthorizationCb authorize_cb, PurpleAccountRequestAuthorizationCb deny_cb, void *user_data)
{
	printf("Buddy authorization request from \"%s\" as \"%s\" (%s): %s\n", remote_user, alias, purple_account_get_protocol_id(account), message);

	authorize_cb(user_data);
	//deny_cb(user_data);

    gopurple_request_authorize();

	return NULL;
}

static void connect_progress(PurpleConnection *gc, const char *text, size_t step, size_t step_count)
{
	printf("Connection in progress: %i/%i: %s\n", step, step_count, text);
    gopurple_connect_progress();
}

static void notice(PurpleConnection *gc, const char *text)
{
	printf("Connection notice: %s\n", text);
    gopurple_notice();
}

static void network_disconnected(void)
{
	printf("This machine has been disconnected from the internet\n");
    gopurple_network_disconnected();
}

static void report_disconnect_reason(PurpleConnection *gc, PurpleConnectionError reason, const char *text)
{
	PurpleAccount *account = purple_connection_get_account(gc);
	printf("Connection disconnected: \"%s\" (%s)\n  >Error: %d\n  >Reason: %s\n", purple_account_get_username(account), purple_account_get_protocol_id(account), reason, text);
    gopurple_report_disconnect_reason();
}

//*** Conversation uiops
// static void write_conv(PurpleConversation *conv, const char *who, const char *alias,
			// const char *message, PurpleMessageFlags flags, time_t mtime)
// {
	// const char *name;
	// if (alias && *alias)
		// name = alias;
	// else if (who && *who)
		// name = who;
	// else
		// name = NULL;

	// printf("write_conv: (%s) %s %s: %s\n", purple_conversation_get_name(conv),
			// purple_utf8_strftime("(%H:%M:%S)", localtime(&mtime)),
			// name, message);
// }

static PurpleAccountUiOps account_uiops =
{
	NULL,
	NULL,
	NULL,
	request_authorize,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL
};

static PurpleConnectionUiOps connection_uiops =
{
	connect_progress,
	NULL,
	NULL,
	notice,
	NULL,
	NULL,
	network_disconnected,
	report_disconnect_reason,
	NULL,
	NULL,
	NULL
};

static PurpleConversationUiOps conv_uiops =
{
	NULL,
	NULL,
	NULL,
	NULL,
	NULL, //write_conv,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL
};

void ui_init(void)
{

	purple_connections_set_ui_ops(&connection_uiops);
	purple_accounts_set_ui_ops(&account_uiops);
	purple_conversations_set_ui_ops(&conv_uiops);

}

static PurpleCoreUiOps core_uiops =
{
	NULL,
	NULL,
	ui_init,
	NULL,

	NULL,
	NULL,
	NULL,
	NULL
};

void init_libpurple(void)
{
	char home[1024];
	sprintf(home, "%s/%s", getenv("HOME"), CUSTOM_USER_DIRECTORY);
	purple_util_set_user_dir(home); // CUSTOM_USER_DIRECTORY

	purple_debug_set_enabled(FALSE);
	// purple_debug_set_enabled(TRUE);

	purple_core_set_ui_ops(&core_uiops);

	purple_eventloop_set_ui_ops(&glib_eventloops);

	purple_plugins_add_search_path(CUSTOM_PLUGIN_PATH);

	if (!purple_core_init(UI_ID)) {
		fprintf(stderr,
				"libpurple initialization failed. Dumping core.\n"
				"Please report this!\n");
		abort();
	}

	purple_set_blist(purple_blist_new());
	purple_blist_load();

	purple_prefs_load();

	purple_plugins_load_saved(PLUGIN_SAVE_PREF);

	purple_pounces_load();
}

static void signed_on(PurpleConnection *gc, void *data)
{
	PurpleAccount *account = purple_connection_get_account(gc);
	printf("Account connected: \"%s\" (%s)\n", purple_account_get_username(account), purple_account_get_protocol_id(account));
    gopurple_signed_on(gc, data);
}

static void buddy_signed_on(PurpleBuddy *buddy)
{
	printf("Buddy \"%s\" (%s) signed on\n", purple_buddy_get_name(buddy), purple_account_get_protocol_id(purple_buddy_get_account(buddy)));
    gopurple_buddy_signed_on();
}

static void buddy_signed_off(PurpleBuddy *buddy)
{
	printf("Buddy \"%s\" (%s) signed off\n", purple_buddy_get_name(buddy), purple_account_get_protocol_id(purple_buddy_get_account(buddy)));
    gopurple_buddy_signed_off();
}

static void buddy_away(PurpleBuddy *buddy, PurpleStatus *old_status, PurpleStatus *status)
{
	printf("Buddy \"%s\" (%s) changed status to %s\n", purple_buddy_get_name(buddy), purple_account_get_protocol_id(purple_buddy_get_account(buddy)), purple_status_get_id(status));
}

static void buddy_idle(PurpleBuddy *buddy, gboolean old_idle, gboolean idle)
{
	printf("Buddy \"%s\" (%s) changed idle state to %s\n", purple_buddy_get_name(buddy), purple_account_get_protocol_id(purple_buddy_get_account(buddy)), (idle) ? "idle" : "not idle");
}

static void readanswer(size_t *size, char *buffer, const char *sender, const char *message)
{
	FILE *fp;
	int status;
	char cmd[1024];

	sprintf(cmd, "purple-answering-machine \"%s\" \"%s\"", sender, message);

	fp = popen(cmd, "r");
	if (fp == NULL) {
		printf("Failed to run command\n" );
		return;
	}

	*size = fread(buffer, *size, 1, fp);

	pclose(fp);
}

static void received_im_msg(PurpleAccount *account, char *sender, char *message, PurpleConversation *conv, PurpleMessageFlags flags)
{
    gopurple_received_im_msg();

	if (conv==NULL)
  	{
  	conv = purple_conversation_new(PURPLE_CONV_TYPE_IM, account, sender);
  	}

	printf("received_im_msg(): (%s) %s (%s): %s\n", purple_utf8_strftime("%H:%M:%S", NULL), sender, purple_conversation_get_name(conv), message);

	// Autoreply from here:
	PurpleConvIm *im = PURPLE_CONV_IM(conv);
	purple_conv_im_set_typing_state(im, PURPLE_TYPING);

	// Let an external program to decide the answer:
	size_t size = 1024;
	char *response = calloc(size+2, sizeof(char));

	readanswer(&size, response, sender, message);
	purple_conv_im_send(im, response);
	purple_conv_im_set_typing_state(im, PURPLE_NOT_TYPING);

	free(response);
}

static void buddy_typing(PurpleAccount *account, const char *name)
{
    gopurple_buddy_typing();
	printf("User \"%s\" (%s) is typing...\n", name, purple_account_get_protocol_id(account));
}

static void buddy_typed(PurpleAccount *account, const char *name) //not supported on all protocols
{
    gopurple_buddy_typed();
	printf("User \"%s\" (%s) has typed something...\n", name, purple_account_get_protocol_id(account));
}

static void buddy_typing_stopped(PurpleAccount *account, const char *name)
{
    gopurple_buddy_typing_stopped();
	printf("User \"%s\" (%s) has stopped typing...\n", name, purple_account_get_protocol_id(account));
}

static int account_authorization_requested(PurpleAccount *account, const char *user)
{
    gopurple_account_authorization_requested();
	printf("User \"%s\" (%s) has sent a buddy request. Now autoaccepting.\n", user, purple_account_get_protocol_id(account));
	return 1; //authorize buddy request automatically (-1 denies it)
}

// static int dbus_method_called(DBusConnection *connection, DBusMessage *message)
// {
    // gopurple_dbus_method_called();
	// printf("DBus: %s\n", message);
	// return TRUE;
// }

static void
received_chat_msg(PurpleAccount *account, char *sender, char *buffer,
                  PurpleConversation *chat, PurpleMessageFlags flags, void *data)
{
    printf("received chat msg: %s / %s \n", sender, buffer);
    gopurple_received_chat_msg();
}

static void connect_to_signals(void)
{
	static int handle;

	purple_signal_connect(purple_connections_get_handle(), "signed-on", &handle,
                          PURPLE_CALLBACK(signed_on), (void*)5678);

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
}

void connect_to_signalscc(void *pc)
{
    printf("connect_to_signals...\n");
    connect_to_signals();
}

void sig_ctrlc(int sig)
{
	signal(SIGINT, sig_ctrlc);
	printf("SIGINT (Ctrl C), quitting...\n");
	purple_core_quit();

	exit(0);
}

int maincc(int argc, char **argv, char **env)
{

#ifdef HAVE_DBUS
	printf("Talking on DBus session %s\nThe dbus-daemon is running at %s\n", getenv("DBUS_SESSION_BUS_ADDRESS"), getenv("DBUS_SESSION_BUS_PID"));
#endif
	GMainLoop *loop = g_main_loop_new(NULL, FALSE);

	signal(SIGCHLD, SIG_IGN);
	signal(SIGINT, sig_ctrlc);

	init_libpurple();
	printf("libpurple initialized. Running version %s.\n", purple_core_get_version()); //I like to see the version number

	BLIST_PREF_SET(blist_choice_MERGE);

	connect_to_signals();

	PurpleSavedStatus *status = purple_savedstatus_find("Online");
	if (status == NULL)
	{
		printf("No savedstatuses found. Creating.\n");
		status = purple_savedstatus_new("Online", PURPLE_STATUS_AVAILABLE);
		purple_savedstatus_set_message(status, "I'm here!");
	} else {
		printf("Existing savedstatuses found. Using.\n");
	}
	purple_savedstatus_set_type(status, PURPLE_STATUS_AVAILABLE);
	purple_savedstatus_activate(status);

	PurpleAccount *account = purple_accounts_find("xyz@gmail.com", "prpl-jabber");
	if (account == NULL)
	{
		printf("No account found. Creating.\n");
		account = purple_account_new("xyz@gmail.com", "prpl-jabber"); //this could be prpl-jabber, prpl-aim, prpl-yahoo, prpl-msn, etc.
		printf("New'd\n");
		purple_account_set_password(account, "zzz");
		printf("Passw'd\n");
		purple_account_set_remember_password(account, TRUE);
		printf("Remember'd\n");
		purple_account_set_enabled(account, UI_ID, TRUE);
		printf("Enable'd\n");
		purple_accounts_add(account);
		printf("Ad'd\n");
	} else {
		printf("Found account. Using.\n");
	}

	printf("Trying to connect\n");
	purple_account_connect(account); // This starts the connection process.
	printf("Connecte'd\n");


	// account = purple_accounts_find("zzz", "prpl-sametime");

	// if (account == NULL)
	// {
		// printf("No account found. Creating.\n");
		// account = purple_account_new("zzz", "prpl-sametime"); //this could be prpl-jabber, prpl-aim, prpl-yahoo, prpl-msn, etc.
		// purple_account_set_string(account, MW_KEY_HOST, "sametime.zzz.com");
		// purple_account_set_string(account, MW_KEY_PORT, "1533");
		//
		// printf("New'd\n");
		// purple_account_set_password(account, "Fenring");
		// printf("Passw'd\n");
		// purple_account_set_remember_password(account, TRUE);
		// printf("Remember'd\n");
		// purple_account_set_enabled(account, UI_ID, TRUE);
		// printf("Enable'd\n");
		// purple_accounts_add(account);
		// printf("Ad'd\n");
	// } else {
		// printf("Found account. Using.\n");
	// }
	// printf("Trying to connect\n");
	// purple_account_connect(account);
	// printf("Connecte'd\n");

	g_main_loop_run(loop);

	return 0;

}
