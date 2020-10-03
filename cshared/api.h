#include <stdbool.h>

// Callbacks
// handle OnTransport 
typedef void (*onTransportCallback)(const int id, const char* transportId, const bool isServer);
// handle OnTransportClosed 
typedef void (*onTransportClosedCallback)(const char* transportId, const bool isServer);
// handle OnChannel with trasnsport id
typedef void (*onChannelCallback)(const char* transportId, const char* channelId);
// handle OnChannelClosed 
typedef void (*onChannelClosedCallback)(const char* channelId);
// handle OnChannelError 
typedef void (*onChannelErrorCallback)(const char* channelId, const char* message);
// handle OnMessage with trasnport id
typedef void (*onMessageCallback)(const char* channelId, const char* message);
// handle OnData with trasnport id
typedef void (*onDataCallback)(const char* channelId, const void *data, const int length);

// Invokes
// invoke OnTransport 
static inline void invokeOnTransport(const int id, const char* transportId, bool isServer, onTransportCallback f){
    f(id, transportId, isServer);
}

static inline void invokeOnTransportClosed(const char* transportId, bool isServer, onTransportClosedCallback f){
    f(transportId, isServer);
}

static inline void invokeOnChannel(const char* transportId, const char* channelId, onChannelCallback f){
    f(transportId, channelId);
}

static inline void invokeOnChannelClosed(const char* channelId, onChannelClosedCallback f){
    f(channelId);
}

static inline void invokeOnChannelError(const char* channelId, const char *message, onChannelErrorCallback f){
    f(channelId, message);
}

static inline void invokeOnMessage(const char* channelId, const char *message, onMessageCallback f){
    f(channelId, message);
}

static inline void invokeOnData(const char* channelId, const void *data, const int length, onDataCallback f){
    f(channelId, data, length);
}