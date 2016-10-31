
#include "misc.h"

GHashTable*
gopurple_connection_get_chat_info_defaults(PurpleConnection* conn, char* name)
{
    PurplePlugin* prpl = purple_connection_get_prpl(conn);
    PurplePluginProtocolInfo* info = (PurplePluginProtocolInfo*)prpl->info->extra_info;
    return info->chat_info_defaults(conn, name);
}
