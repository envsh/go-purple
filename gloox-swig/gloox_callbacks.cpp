#include <iostream>
#include <cstring>
#include <memory>
#include <cstdlib>

#include "gloox_combine.h"
#include "gloox_callbacks.h"
#include "gloox_clone.h"

using namespace gloox;


class BaseRCB {
public:
    uint64_t gobjno = 0;
};

// RCB reverse callback
class MessageHandlerRCB : public MessageHandler, public BaseRCB {
public:
    virtual void handleMessage( const Message& msg, MessageSession* session = 0 ) {
        // std::cout<<msg.body( )<<"--"<<msg.subject( )<<std::endl;
        // int subType = msg.subtype();
        // char *body = strdup((char*)msg.body( ).c_str());
        // char *subject = strdup((char*)msg.subject( ).c_str());
        // char *from = strdup((char*)msg.from().full().c_str());
        // char *to = strdup((char*)msg.to().full().c_str());
        Message *nmsg = cloneMessage(msg);
        MessageHandlerRCB_handleMessage
            (gobjno, (uint64_t)nmsg, (uint64_t)session);
    }
};


uint64_t MessageHandlerRCB_new(uint64_t gobjno){
    auto a = new MessageHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void MessageHandlerRCB_delete(uint64_t cobjno) {
    delete((MessageHandlerRCB*)cobjno);
}

class ConnectionListenerRCB : public ConnectionListener, public BaseRCB {
public:
    virtual void onConnect(){
        ConnectionListenerRCB_onConnect_go(gobjno);
    }
    virtual void onDisconnect( ConnectionError e ){
        ConnectionListenerRCB_onDisconnect_go(gobjno, int(e));
    }
    virtual bool onTLSConnect( const CertInfo& info ){
        ConnectionListenerRCB_onTLSConnect_go(gobjno, (uint64_t)cloneCertInfo(info));
        return true;
    }
};

uint64_t ConnectionListenerRCB_new(uint64_t gobjno){
    auto a = new ConnectionListenerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void ConnectionListenerRCB_delete(uint64_t cobjno){
    delete((ConnectionListenerRCB*)cobjno);
}

class LogHandlerRCB : public LogHandler, public BaseRCB {
public:
    virtual void handleLog(LogLevel level, LogArea area, const std::string& message) {
        LogHandlerRCB_handleLog(gobjno, int(level), int(area), (char*)message.c_str());
    }
};
uint64_t LogHandlerRCB_new(uint64_t gobjno){
    auto a = new LogHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void LogHandlerRCB_delete(uint64_t cobjno){
    delete((LogHandlerRCB*)cobjno);
}

class MUCInvitationHandlerRCB:public MUCInvitationHandler,public BaseRCB{
public:
    virtual void handleMUCInvitation( const JID& room, const JID& from, const std::string& reason,
                                      const std::string& body, const std::string& password,
                                      bool cont, const std::string& thread ){
        MUCInvitationHandlerRCB_handleMUCInvitation
            (gobjno,
             strdup(room.full().c_str()),
             strdup(from.full().c_str()),
             strdup(reason.c_str()),
             strdup(body.c_str()),
             strdup(password.c_str()),
             cont?1:0,
             strdup(thread.c_str())
             );
    }
};

uint64_t MUCInvitationHandlerRCB_new(uint64_t gobjno){
    // TODO howhow???
    // auto a = new MUCInvitationHandlerRCB(cb);
    // a->gobjno = gobjno;
    return (uint64_t)0;
}
void MUCInvitationHandlerRCB_delete(uint64_t cobjno){
    delete((MUCInvitationHandlerRCB*)cobjno);
}


class PresenceHandlerRCB : public PresenceHandler, public BaseRCB{
public:
    virtual void handlePresence( const Presence& presence ) {

        PresenceHandlerRCB_handlePresence
            (gobjno,
             int(presence.presence()),
             strdup(presence.from().full().c_str()),
             strdup(presence.to().full().c_str()),
             strdup(presence.status( ).c_str())
             );
    }
};

uint64_t PresenceHandlerRCB_new(uint64_t gobjno){
    auto a = new PresenceHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void PresenceHandlerRCB_delete(uint64_t cobjno){
    delete((PresenceHandlerRCB*)cobjno);
}

class SubscriptionHandlerRCB : public SubscriptionHandler, public BaseRCB{
public:
    virtual void handleSubscription(const Subscription& sub) {
        SubscriptionHandlerRCB_handleSubscription
            (gobjno,
             int(sub.subtype()),
             strdup(sub.from().full().c_str()),
             strdup(sub.to().full().c_str()),
             strdup(sub.status().c_str())
             );
    }
};

uint64_t SubscriptionHandlerRCB_new(uint64_t gobjno){
    auto a = new SubscriptionHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void SubscriptionHandlerRCB_delete(uint64_t cobjno){
    delete((SubscriptionHandlerRCB*)cobjno);
}

class StatisticsHandlerRCB : public StatisticsHandler, public BaseRCB{
public:
    virtual void handleStatistics( const StatisticsStruct stats ){
        StatisticsHandlerRCB_handleStatistics
            (gobjno,
             stats.totalBytesSent,
             stats.totalBytesReceived,
             stats.compressedBytesSent,
             stats.compressedBytesReceived,
             stats.uncompressedBytesSent,
             stats.uncompressedBytesReceived,
             stats.totalStanzasSent,
             stats.totalStanzasReceived,
             stats.iqStanzasSent,
             stats.iqStanzasReceived,
             stats.messageStanzasSent,
             stats.messageStanzasReceived,
             stats.s10nStanzasSent,
             stats.s10nStanzasReceived,
             stats.presenceStanzasSent,
             stats.presenceStanzasReceived,
             stats.encryption?1:0,
             stats.compression?1:0
             );
    }
};

uint64_t StatisticsHandlerRCB_new(uint64_t gobjno){
    auto a = new StatisticsHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
    
}
void StatisticsHandlerRCB_delete(uint64_t cobjno){
    delete((StatisticsHandlerRCB*)cobjno);
    
}

class TagHandlerRCB : public TagHandler, public BaseRCB{
public:
    virtual void handleTag( Tag *tag ){
        tag->xml();
        tag->prefix();
        tag->xmlns();
        tag->name();
        tag->cdata();
        tag->findAttribute("abc");
        TagHandlerRCB_handleTag
            (gobjno
             );
    }
};

uint64_t TagHandlerRCB_new(uint64_t gobjno){
    auto a = new TagHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void TagHandlerRCB_delete(uint64_t cobjno){
    delete((TagHandlerRCB*)cobjno);
}


class MUCRoomHandlerRCB : public MUCRoomHandler, public BaseRCB {
public:
    virtual void handleMUCParticipantPresence( MUCRoom* room, const MUCRoomParticipant participant, const Presence& presence ) {
        // MUCRoomHandlerRCB_handleMUCParticipantPresence
        // (gobjno, (uint64_t)room);
        Presence *npres = clonePresense(presence);
        MUCRoomParticipant *npart = cloneParticipant(participant);
        MUCRoomHandlerRCB_handleMUCParticipantPresence
            (gobjno, (uint64_t)room, (uint64_t)npres, (uint64_t)npart);
    }

    virtual void handleMUCMessage( MUCRoom* room, const Message& msg, bool priv ) {
        Message *nmsg = cloneMessage(msg);
        MUCRoomHandlerRCB_handleMUCMessage
            (gobjno, (uint64_t)room, (uint64_t)nmsg, int(priv));
    }

    virtual bool handleMUCRoomCreation( MUCRoom* room ) {
        MUCRoomHandlerRCB_handleMUCRoomCreation
            (gobjno, (uint64_t)room);
        return true;
    }

    virtual void handleMUCSubject( MUCRoom* room, const std::string& nick,
                                   const std::string& subject ) {
        MUCRoomHandlerRCB_handleMUCSubject
            (gobjno, (uint64_t)room, cloneString(nick), cloneString(subject));
    }

    virtual void handleMUCInviteDecline( MUCRoom* room, const JID& invitee,
                                         const std::string& reason ) {
        MUCRoomHandlerRCB_handleMUCInviteDecline
            (gobjno, (uint64_t)room, cloneString(invitee.full()), cloneString(reason));
    }

    virtual void handleMUCError( MUCRoom* room, StanzaError error ) {
        MUCRoomHandlerRCB_handleMUCError
            (gobjno, (uint64_t)room, int(error));
    }

    virtual void handleMUCInfo( MUCRoom* room, int features, const std::string& name,
                                const DataForm* infoForm ) {
        MUCRoomHandlerRCB_handleMUCInfo
            (gobjno, (uint64_t)room, int(features),
             cloneString(name), (uint64_t)infoForm);
    }

    virtual void handleMUCItems( MUCRoom* room, const Disco::ItemList& items ) {
        MUCRoomHandlerRCB_handleMUCItems
            (gobjno, (uint64_t)room);
    }

};

uint64_t MUCRoomHandlerRCB_new(uint64_t gobjno) {
    auto a = new MUCRoomHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void MUCRoomHandlerRCB_delete(uint64_t cobjno){
    delete((MUCRoomHandlerRCB*)cobjno);
}

/////
class MUCRoomConfigHandlerRCB : public MUCRoomConfigHandler, public BaseRCB{
public:
    virtual void handleMUCConfigList (MUCRoom *room, const MUCListItemList &items, MUCOperation operation){
        MUCRoomConfigHandlerRCB_handleMUCConfigList
            (gobjno, (uint64_t)room,
             (uint64_t)cloneMUCListItemList(items), int(operation));
    }

    virtual void handleMUCConfigForm (MUCRoom *room, const DataForm &form){
        MUCRoomConfigHandlerRCB_handleMUCConfigForm
            (gobjno, (uint64_t)room, (uint64_t)cloneDataFrom(form));
    }

    virtual void handleMUCConfigResult (MUCRoom *room, bool success, MUCOperation operation){
        MUCRoomConfigHandlerRCB_handleMUCConfigResult
            (gobjno, (uint64_t)room, success?1:0, int(operation));
    }

    virtual void handleMUCRequest (MUCRoom *room, const DataForm &form) {
        MUCRoomConfigHandlerRCB_handleMUCRequest
            (gobjno, (uint64_t)room, (uint64_t)cloneDataFrom(form));
    }
};

uint64_t MUCRoomConfigHandlerRCB_new(uint64_t gobjno){
    auto a = new MUCRoomConfigHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void MUCRoomConfigHandlerRCB_delete(uint64_t cobjno){
    delete((MUCRoomConfigHandlerRCB*)cobjno);
}

///
class MessageSessionHandlerRCB : public MessageSessionHandler, public BaseRCB {
public:
    virtual void 	handleMessageSession (MessageSession *session) {
        MessageSessionHandlerRCB_handleMessageSession
            (gobjno, (uint64_t)session);
    }
};

uint64_t MessageSessionHandlerRCB_new(uint64_t gobjno) {
    auto a = new MessageSessionHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void MessageSessionHandlerRCB_delete(uint64_t cobjno) {
    delete((MessageSessionHandlerRCB*)cobjno);
}

///
class MessageEventHandlerRCB : public MessageEventHandler, public BaseRCB {
public:
    virtual void 	handleMessageEvent (const JID &from, MessageEventType event) {
        MessageEventHandlerRCB_handleMessageEvent
            (gobjno, (uint64_t)cloneJID(from), int(event));
    }
};

uint64_t MessageEventHandlerRCB_new(uint64_t gobjno) {
    auto a = new MessageEventHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void MessageEventHandlerRCB_delete(uint64_t cobjno) {
    delete((MessageEventHandlerRCB*)cobjno);
}

///
class ChatStateHandlerRCB : public ChatStateHandler, public BaseRCB {
public:
    virtual void 	handleChatState (const JID &from, ChatStateType state) {
        ChatStateHandlerRCB_handleChatState
            (gobjno, (uint64_t)cloneJID(from), int(state));
    }
};

uint64_t ChatStateHandlerRCB_new(uint64_t gobjno) {
    auto a = new ChatStateHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void ChatStateHandlerRCB_delete(uint64_t cobjno) {
    delete((ChatStateHandlerRCB*)cobjno);
}

///
class EventHandlerRCB : public EventHandler, public BaseRCB {
public:
    virtual void 	handleEvent (const Event& event) {
        Event *nevent = cloneEvent(event);
        EventHandlerRCB_handleEvent
            (gobjno, (uint64_t)(nevent));
    }
};

uint64_t EventHandlerRCB_new(uint64_t gobjno) {
    auto a = new EventHandlerRCB();
    a->gobjno = gobjno;
    return (uint64_t)a;
}
void EventHandlerRCB_delete(uint64_t cobjno) {
    delete((EventHandlerRCB*)cobjno);
}
