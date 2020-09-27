#include <stdint.h>

// Callbacks
// handle OnTransport 
typedef void (*onTransportCallback)(const uintptr_t);
// handle OnTransportClosed 
typedef void (*onTransportClosedCallback)();
// handle OnChannel with trasnsport id
typedef void (*onChannelCallback)(const uintptr_t);
// handle OnChannelClosed 
typedef void (*onChannelClosedCallback)();
// handle OnChannelError 
typedef void (*onChannelErrorCallback)(const char *);
// handle OnMessage with trasnport id
typedef void (*onMessageCallback)(const char *);
// handle OnData with trasnport id
typedef void (*onDataCallback)(const void *, const int length);

// Invokes
// invoke OnTransport 
static inline void invokeOnTransport(const uintptr_t ptr, onTransportCallback f){
    f(ptr);
}

static inline void invokeOnTransportClosed(onTransportClosedCallback f){
    f();
}

static inline void invokeOnChannel(const uintptr_t ptr, onChannelCallback f){
    f(ptr);
}

static inline void invokeOnChannelClosed(onChannelClosedCallback f){
    f();
}

static inline void invokeOnChannelError(const char *message, onChannelErrorCallback f){
    f(message);
}

static inline void invokeOnMessage(const char *message, onMessageCallback f){
    f(message);
}

static inline void invokeOnData(const void *data, const int length, onDataCallback f){
    f(data, length);
}