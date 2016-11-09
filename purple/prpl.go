package purple

/*
#include <libpurple/purple.h>

static void gopurple_prpl_got_user_status(PurpleAccount *account,
     const char *name, const char *status_id)
{ purple_prpl_got_user_status(account, name, status_id, NULL); }

*/
import "C"

func PrplGotUserStatus(ac *Account, name, statusId string) {
	C.gopurple_prpl_got_user_status(ac.account, C.CString(name), C.CString(statusId))
}
