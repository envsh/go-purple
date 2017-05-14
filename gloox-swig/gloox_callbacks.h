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
    uint64_t MUCInvitationHandlerRCB_new(uint64_t objno);
    void MUCInvitationHandlerRCB_delete(uint64_t objno);
    extern void MUCInvitationHandlerRCB_handleMUCInvitation
    (uint64_t objno, char*room,char*from,char*reason,char*body,
     char*password, int cont, char*thread);

    uint64_t PresenceHandlerRCB_new(uint64_t objno);
    void PresenceHandlerRCB_delete(uint64_t objno);
    extern void PresenceHandlerRCB_handlePresence
    (uint64_t objno, int ptype, char*from, char*to, char*status);

    uint64_t SubscriptionHandlerRCB_new(uint64_t objno);
    void SubscriptionHandlerRCB_delete(uint64_t objno);
    extern void SubscriptionHandlerRCB_handleSubscription
    (uint64_t objno, int ptype, char*from, char*to, char*status);

    uint64_t StatisticsHandlerRCB_new(uint64_t objno);
    void StatisticsHandlerRCB_delete(uint64_t objno);
    extern void StatisticsHandlerRCB_handleStatistics
    (uint64_t objno,
     long int totalBytesSent,
     long int totalBytesReceived,
     long int compressedBytesSent,
     long int compressedBytesReceived,
     long int uncompressedBytesSent,
     long int uncompressedBytesReceived,
     long int totalStanzasSent,
     long int totalStanzasReceived,
     long int iqStanzasSent,
     long int iqStanzasReceived,
     long int messageStanzasSent,
     long int messageStanzasReceived,
     long int s10nStanzasSent,
     long int s10nStanzasReceived,
     long int presenceStanzasSent,
     long int presenceStanzasReceived,
     int encryption,
     int compression
);


    uint64_t TagHandlerRCB_new(uint64_t objno);
    void TagHandlerRCB_delete(uint64_t objno);
    extern void TagHandlerRCB_handleTag
    (uint64_t objno);

#ifdef __cplusplus
};
#endif

#endif
