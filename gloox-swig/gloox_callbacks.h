#ifndef _GLOOX_CALLBACKS_H_
#define _GLOOX_CALLBACKS_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C"{
#endif
    uint64_t MessageHandlerRCB_new(uint64_t objno);
    void MessageHandlerRCB_delete(uint64_t objno);
    extern void MessageHandlerRCB_handleMessage
    (uint64_t objno,int subType,char*from,char*to,char*subject,char*body);
    uint64_t ConnectionHandlerRCB_new(uint64_t objno);
    void ConnectionHandlerRCB_delete(uint64_t objno);
    uint64_t ConnectionListenerRCB_new(uint64_t objno);
    void ConnectionListenerRCB_delete(uint64_t objno);
    extern void ConnectionListenerRCB_onConnect_go(uint64_t objno);
    extern void ConnectionListenerRCB_onDisconnect_go(uint64_t objno, int error);
    extern void ConnectionListenerRCB_onTLSConnect_go(uint64_t objno);
    uint64_t LogHandlerRCB_new(uint64_t objno);
    void LogHandlerRCB_delete(uint64_t objno);
    extern void LogHandlerRCB_handleLog(uint64_t objno, int level, int area, char *l);
#ifdef __cplusplus
};
#endif

#endif
