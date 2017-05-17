#ifndef _GLOOX_CLONE_H_
#define _GLOOX_CLONE_H_

#include <cstring>

#include "gloox_combine.h"

using namespace gloox;

auto cloneString = [](const std::string&s) {
    return strdup(s.c_str());
};
auto cloneString2 = [](const std::string&s) {
    return s;
};
auto cloneJID = [](const JID&jid) {return new JID(jid);};

auto cloneMessage = [](const Message &msg) {
    Message *nmsg = new Message
    (msg.subtype(), msg.to(), cloneString2(msg.body( )),
     cloneString2(msg.subject( )),
     cloneString2(msg.thread())/*, cloneString2(msg.xmlLang())*/);
    // dont add msg.xmlLang() param, it's default here
    return nmsg;
};
auto clonePresense = [](const Presence &presence) {
    Presence *npres = new Presence
    (presence.presence(), presence.to(), presence.status(),
     presence.priority());
    return npres;
};
auto cloneParticipant = [](const MUCRoomParticipant &part) {
    MUCRoomParticipant *npart = new MUCRoomParticipant();
    *npart = part;
    return npart;
};


auto cloneDataFrom = [](const DataForm& form) {
    return new DataForm(form);
};

auto cloneMUCListItemList = [](const MUCListItemList&items) {
    return new MUCListItemList(items);
};

auto cloneCertInfo = [](const CertInfo&info) {
    CertInfo *ninfo = new CertInfo();
    *ninfo = info;
    return ninfo;
};

auto cloneEvent = [](const Event& event) {
    Event *nevent = new Event(event.eventType(), *event.stanza());
    return nevent;
};

#endif
