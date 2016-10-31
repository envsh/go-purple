#ifndef _GOPURPLE_MISC_H_
#define _GOPURPLE_MISC_H_

#include <libpurple/purple.h>

#include <glib.h>
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <string.h>
#include <unistd.h>

/*
  for function pointer that can not called from go directly
 */


GHashTable*
gopurple_connection_get_chat_info_defaults(PurpleConnection* conn, char* name);

#endif

