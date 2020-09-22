package main

// #include "api.h"
import "C"
import (
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/pkg/channel"
)

var client *sylph.Client
var onTransportCallback C.onTransportCallback = nil

//export Initialize
func Initialize() {
	client = sylph.NewClient()
	client.OnTransport = func(t sylph.Transport) {
		C.invokeOnTransport(C.CString(t.Id()), onTransportCallback)
	}
}

//export Connect
func Connect(address *C.char, port C.int) {
	client.Connect(C.GoString(address), int(port))
}

//export OpenChannel
func OpenChannel(transportId *C.char) {
	if t := client.Transport(C.GoString(transportId)); t != nil {
		c := channel.ChannelConfig{}
		t.OpenChannel(c)
	}
}

//export CloseChannel
func CloseChannel(transportId *C.char, channelId *C.char) {
	if c := findChannel(transportId, channelId); c != nil {
		c.Close()
	}
}

//export RegisterOnTransportCallback
func RegisterOnTransportCallback(callback C.onTransportCallback) {
	onTransportCallback = callback
}

//export RegisterOnChannelCallback
func RegisterOnChannelCallback(id *C.char, callback C.onChannelCallback) {
	if t := client.Transport(C.GoString(id)); t != nil {
		t.OnChannel(func(channel channel.Channel) {
			C.invokeOnChannel(C.CString(channel.Id()), callback)
		})
	}
}

//export RegisterOnChannelClosedCallback
func RegisterOnChannelClosedCallback(transportId *C.char, channelId *C.char, callback C.onChannelClosedCallback) {
	if c := findChannel(transportId, channelId); c != nil {
		c.OnClose(func() {
			C.invokeOnChannelClosed(callback)
		})
	}
}

//export RegisterOnChannelErrorCallback
func RegisterOnChannelErrorCallback(transportId *C.char, channelId *C.char, callback C.onChannelErrorCallback) {
	if c := findChannel(transportId, channelId); c != nil {
		c.OnError(func(err error) {
			C.invokeOnChannelError(C.CString(err.Error()), callback)
		})
	}
}

//export RegisterOnMessageCallback
func RegisterOnMessageCallback(transportId *C.char, channelId *C.char, callback C.onMessageCallback) {
	if c := findChannel(transportId, channelId); c != nil {
		c.OnMessage(func(message string) {
			C.invokeOnMessage(C.CString(message), callback)
		})
	}
}

//export RegisterOnDataCallback
func RegisterOnDataCallback(trasnportId *C.char, channelId *C.char, callback C.onDataCallback) {
	if c := findChannel(trasnportId, channelId); c != nil {
		c.OnData(func(data []byte) {
			C.invokeOnData(unsafe.Pointer(&data[0]), C.int(len(data)), callback)
		})
	}
}

func findChannel(transportId *C.char, channelId *C.char) channel.Channel {
	t := client.Transport(C.GoString(transportId))
	if t == nil {
		return nil
	}
	return t.Channel(C.GoString(channelId))
}

func main() {}
