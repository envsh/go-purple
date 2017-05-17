#include "gloox_combine.h"
#include "gloox_extras.h"

using namespace gloox;

///
uint64_t RefillDataForm(uint64_t dfx) {
    DataForm *form = (DataForm*)(dfx);
    DataForm* pNewFm = new DataForm(TypeSubmit);
    const DataForm::FieldList& fl = form->fields();
    for (const DataFormField* pFld : fl) {
        DataFormField* pNewFld = pNewFm->addField
            (pFld->type(), pFld->name(), pFld->value());
        if (pFld->name() == "muc#roomconfig_roomdesc")
            pNewFld->setValue("RoomDesc_JinqTest");
        else if (pFld->name() == "muc#roomconfig_publicroom")
            pNewFld->setValue("0");
        else if (pFld->name() == "muc#roomconfig_enablelogging")
            pNewFld->setValue("0");
        else if (pFld->name() == "muc#roomconfig_passwordprotectedroom")
            pNewFld->setValue("1");
        else if (pFld->name() == "muc#roomconfig_roomsecret")
            pNewFld->setValue("123456");
        else if (pFld->name() == "public_list")
            pNewFld->setValue("0");
        else {}
    }
    return (uint64_t)pNewFm;
}
