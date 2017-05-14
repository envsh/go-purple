#include <iostream>
#include <cstring>
#include <memory>
#include <cstdlib>

#include "gloox_combine.h"
#include "gloox_callbacks.h"

using namespace gloox;

class BaseRCB {
public:
    uint64_t objno = 0;
};

// RCB reverse callback
class MessageHandlerRCB : public MessageHandler, public BaseRCB {
public:
    virtual void handleMessage( const Message& msg, MessageSession* session = 0 ) {
        // std::cout<<msg.body( )<<"--"<<msg.subject( )<<std::endl;
        int subType = msg.subtype();
        char *body = strdup((char*)msg.body( ).c_str());
        char *subject = strdup((char*)msg.subject( ).c_str());
        char *from = strdup((char*)msg.from().full().c_str());
        char *to = strdup((char*)msg.to().full().c_str());
        MessageHandlerRCB_handleMessage
            (objno, subType, from, to, subject, body);
    }
};


uint64_t MessageHandlerRCB_new(uint64_t objno){
    auto a = new MessageHandlerRCB();
    a->objno = objno;
    return (uint64_t)a;
}
void MessageHandlerRCB_delete(uint64_t objno) {
    delete((MessageHandlerRCB*)objno);
}

class ConnectionListenerRCB : public ConnectionListener, public BaseRCB {
public:
    virtual void onConnect(){
        ConnectionListenerRCB_onConnect_go(objno);
    }
    virtual void onDisconnect( ConnectionError e ){
        ConnectionListenerRCB_onDisconnect_go(objno, int(e));
    }
    virtual bool onTLSConnect( const CertInfo& info ){
        ConnectionListenerRCB_onTLSConnect_go(objno);
        return true;
    }
};

uint64_t ConnectionListenerRCB_new(uint64_t objno){
    auto a = new ConnectionListenerRCB();
    a->objno = objno;
    return (uint64_t)a;
}
void ConnectionListenerRCB_delete(uint64_t objno){
    delete((ConnectionListenerRCB*)objno);
}

class LogHandlerRCB : public LogHandler, public BaseRCB {
public:
    virtual void handleLog(LogLevel level, LogArea area, const std::string& message) {
        LogHandlerRCB_handleLog(objno, int(level), int(area), (char*)message.c_str());
    }
};
uint64_t LogHandlerRCB_new(uint64_t objno){
    auto a = new LogHandlerRCB();
    a->objno = objno;
    return (uint64_t)a;
}
void LogHandlerRCB_delete(uint64_t objno){
    delete((LogHandlerRCB*)objno);
}

class MUCInvitationHandlerRCB:public MUCInvitationHandler,public BaseRCB{
public:
    virtual void handleMUCInvitation( const JID& room, const JID& from, const std::string& reason,
                                      const std::string& body, const std::string& password,
                                      bool cont, const std::string& thread ){
        MUCInvitationHandlerRCB_handleMUCInvitation
            (objno,
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

uint64_t MUCInvitationHandlerRCB_new(uint64_t objno){
    // TODO howhow???
    // auto a = new MUCInvitationHandlerRCB(cb);
    // a->objno = objno;
    return (uint64_t)0;
}
void MUCInvitationHandlerRCB_delete(uint64_t objno){
    delete((MUCInvitationHandlerRCB*)objno);
}

class PresenceHandlerRCB : public PresenceHandler, public BaseRCB{
public:
    virtual void handlePresence( const Presence& presence ) {

        PresenceHandlerRCB_handlePresence
            (objno,
             int(presence.presence()),
             strdup(presence.from().full().c_str()),
             strdup(presence.to().full().c_str()),
             strdup(presence.status( ).c_str())
             );
    }
};

uint64_t PresenceHandlerRCB_new(uint64_t objno){
    auto a = new PresenceHandlerRCB();
    a->objno = objno;
    return (uint64_t)a;
}
void PresenceHandlerRCB_delete(uint64_t objno){
    delete((PresenceHandlerRCB*)objno);
}

class SubscriptionHandlerRCB : public SubscriptionHandler, public BaseRCB{
public:
    virtual void handleSubscription(const Subscription& sub) {
        SubscriptionHandlerRCB_handleSubscription
            (objno,
             int(sub.subtype()),
             strdup(sub.from().full().c_str()),
             strdup(sub.to().full().c_str()),
             strdup(sub.status().c_str())
             );
    }
};

uint64_t SubscriptionHandlerRCB_new(uint64_t objno){
    auto a = new SubscriptionHandlerRCB();
    a->objno = objno;
    return (uint64_t)a;
}
void SubscriptionHandlerRCB_delete(uint64_t objno){
    delete((SubscriptionHandlerRCB*)objno);
}

class StatisticsHandlerRCB : public StatisticsHandler, public BaseRCB{
public:
    virtual void handleStatistics( const StatisticsStruct stats ){
        StatisticsHandlerRCB_handleStatistics
            (objno,
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

uint64_t StatisticsHandlerRCB_new(uint64_t objno){
    auto a = new StatisticsHandlerRCB();
    a->objno = objno;
    return (uint64_t)a;
    
}
void StatisticsHandlerRCB_delete(uint64_t objno){
    delete((StatisticsHandlerRCB*)objno);
    
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
            (objno
             );
    }
};

uint64_t TagHandlerRCB_new(uint64_t objno){
    auto a = new TagHandlerRCB();
    a->objno = objno;
    return (uint64_t)a;
    
}
void TagHandlerRCB_delete(uint64_t objno){
    delete((StatisticsHandlerRCB*)objno);
    
}

