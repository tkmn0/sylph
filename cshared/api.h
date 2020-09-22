// Callbacks
// handle OnTransport with transport id
typedef void (*onTransportCallback)(const char *);
// handle OnChannel with trasnsport id
typedef void (*onChannelCallback)(const char *);
// handle OnChannelClosed with transport id
typedef void (*onChannelClosedCallback)();

// Invokes
// invoke OnTransport with transport id
static inline void invokeOnTransport(const char *id, onTransportCallback f){
    f(id);
}

static inline void invokeOnChannel(const char *id, onChannelCallback f){
    f(id);
}

static inline void invokeOnChannelClosed(onChannelClosedCallback f){
    f();
}