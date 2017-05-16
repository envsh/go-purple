#ifndef _GLOOX_CALLBACKS_H_
#define _GLOOX_CALLBACKS_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C"{
#endif
    uint64_t MessageHandlerRCB_new(uint64_t gobjno);
    void MessageHandlerRCB_delete(uint64_t cobjno);
    extern void MessageHandlerRCB_handleMessage
    (uint64_t gobjno, uint64_t msg, uint64_t session);

    uint64_t ConnectionHandlerRCB_new(uint64_t gobjno);
    void ConnectionHandlerRCB_delete(uint64_t cobjno);
    uint64_t ConnectionListenerRCB_new(uint64_t gobjno);
    void ConnectionListenerRCB_delete(uint64_t cobjno);
    extern void ConnectionListenerRCB_onConnect_go(uint64_t gobjno);
    extern void ConnectionListenerRCB_onDisconnect_go(uint64_t gobjno, int error);
    extern void ConnectionListenerRCB_onTLSConnect_go(uint64_t gobjno, uint64_t info);
    uint64_t LogHandlerRCB_new(uint64_t gobjno);
    void LogHandlerRCB_delete(uint64_t cobjno);
    extern void LogHandlerRCB_handleLog(uint64_t gobjno, int level, int area, char *l);
    uint64_t MUCInvitationHandlerRCB_new(uint64_t gobjno);
    void MUCInvitationHandlerRCB_delete(uint64_t cobjno);
    extern void MUCInvitationHandlerRCB_handleMUCInvitation
    (uint64_t gobjno, char*room,char*from,char*reason,char*body,
     char*password, int cont, char*thread);

    uint64_t PresenceHandlerRCB_new(uint64_t gobjno);
    void PresenceHandlerRCB_delete(uint64_t cobjno);
    extern void PresenceHandlerRCB_handlePresence
    (uint64_t gobjno, int ptype, char*from, char*to, char*status);

    uint64_t SubscriptionHandlerRCB_new(uint64_t gobjno);
    void SubscriptionHandlerRCB_delete(uint64_t cobjno);
    extern void SubscriptionHandlerRCB_handleSubscription
    (uint64_t gobjno, int ptype, char*from, char*to, char*status);

    uint64_t StatisticsHandlerRCB_new(uint64_t gobjno);
    void StatisticsHandlerRCB_delete(uint64_t cobjno);
    extern void StatisticsHandlerRCB_handleStatistics
    (uint64_t gobjno,
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

    uint64_t TagHandlerRCB_new(uint64_t gobjno);
    void TagHandlerRCB_delete(uint64_t cobjno);
    extern void TagHandlerRCB_handleTag
    (uint64_t gobjno);

    uint64_t MUCRoomHandlerRCB_new(uint64_t gobjno);
    void MUCRoomHandlerRCB_delete(uint64_t cobjno);

    extern void MUCRoomHandlerRCB_handleMUCParticipantPresence
    (uint64_t gobjno, uint64_t room, uint64_t participant, uint64_t presence/* MUCRoom* room, const MUCRoomParticipant participant, const Presence& presence */);
    extern void MUCRoomHandlerRCB_handleMUCMessage
    (uint64_t gobjno, uint64_t room, uint64_t msg, int priv /* MUCRoom* room, const Message& msg, bool priv */);
    extern int MUCRoomHandlerRCB_handleMUCRoomCreation
    (uint64_t gobjno, uint64_t room/* MUCRoom* room*/ );
    extern void MUCRoomHandlerRCB_handleMUCSubject
    (uint64_t gobjno, uint64_t room, char*nick, char*subject /*MUCRoom* room, const std::string& nick, const std::string& subject*/ ) ;
    extern void MUCRoomHandlerRCB_handleMUCInviteDecline
    (uint64_t gobjno, uint64_t room, char*invitee, char*reason/* MUCRoom* room, const JID& invitee, const std::string& reason */);
    extern void MUCRoomHandlerRCB_handleMUCError
    (uint64_t gobjno, uint64_t room, int error/* MUCRoom* room, StanzaError error */);
    extern void MUCRoomHandlerRCB_handleMUCInfo
    (uint64_t gobjno, uint64_t room, int features, char*name, uint64_t infoForm/* MUCRoom* room, int features, const std::string& name, const DataForm* infoForm*/ );
    extern void MUCRoomHandlerRCB_handleMUCItems
    (uint64_t gobjno, uint64_t room/* MUCRoom* room, const Disco::ItemList& items */);

    uint64_t MUCRoomConfigHandlerRCB_new(uint64_t gobjno);
    void MUCRoomConfigHandlerRCB_delete(uint64_t cobjno);
    extern void MUCRoomConfigHandlerRCB_handleMUCConfigList
    (uint64_t gobjno, uint64_t room, uint64_t items, int operation);
    extern void MUCRoomConfigHandlerRCB_handleMUCConfigForm
    (uint64_t gobjno, uint64_t room, uint64_t form);
    extern void MUCRoomConfigHandlerRCB_handleMUCConfigResult
    (uint64_t gobjno, uint64_t room, int success, int operation);
    extern void MUCRoomConfigHandlerRCB_handleMUCRequest
    (uint64_t gobjno, uint64_t room, uint64_t form);


    extern uint64_t RefillDataForm(uint64_t df);
#ifdef __cplusplus
};
#endif

#endif
