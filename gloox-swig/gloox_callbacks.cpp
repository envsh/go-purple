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
