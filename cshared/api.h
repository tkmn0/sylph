// Callbacks
// handle OnTransport with transport id
typedef void (*onTransportCallback)(const char *);
// handle OnChannel with trasnsport id
typedef void (*onChannelCallback)(const char *);
// handle OnChannelClosed with transport id
typedef void (*onChannelClosedCallback)();
// handle OnChannelError with transport id
typedef void (*onChannelErrorCallback)(const char *);
// handle OnMessage with trasnport id
typedef void (*onMessageCallback)(const char *);
// handle OnData with trasnport id
typedef void (*onDataCallback)(const void *, const int length);

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

static inline void invokeOnChannelError(const char *message, onChannelErrorCallback f){
    f(message);
}

static inline void invokeOnMessage(const char *message, onMessageCallback f){
    f(message);
}

static inline void invokeOnData(const void *data, const int length, onDataCallback f){
    f(data, length);
}