package main

// #include "api.h"
// #include "stdlib.h"
import "C"
import (
	"fmt"

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
	fmt.Println("address:", C.GoString(address), "port:", port)
	client.Connect(C.GoString(address), int(port))
}

//export OpenChannel
func OpenChannel(transportId *C.char) {
	t := client.Transport(C.GoString(transportId))
	if t == nil {
		return
	}
	c := channel.ChannelConfig{}
	t.OpenChannel(c)
}

//export CloseChannel
func CloseChannel(transportId *C.char, channelId *C.char) {
	t := client.Transport(C.GoString(transportId))
	if t == nil {
		return
	}
	c := t.Channel(C.GoString(channelId))
	if c == nil {
		return
	}

	fmt.Println("call Channel.Close with:", C.GoString(transportId), C.GoString(channelId))

	c.Close()
}

//export TestDebug
func TestDebug() {
	fmt.Println("test debug")
}

//export RegisterOnTransportCallback
func RegisterOnTransportCallback(callback C.onTransportCallback) {
	onTransportCallback = callback
}

//export RegisterOnChannelCallback
func RegisterOnChannelCallback(id *C.char, callback C.onChannelCallback) {
	t := client.Transport(C.GoString(id))
	if t == nil {
		return
	}
	t.OnChannel(func(channel channel.Channel) {
		C.invokeOnChannel(C.CString(channel.Id()), callback)
	})
}

//export RegisterOnChannelClosedCallback
func RegisterOnChannelClosedCallback(transportId *C.char, channelId *C.char, callback C.onChannelClosedCallback) {
	t := client.Transport(C.GoString(transportId))
	transportIdGo := C.GoString(transportId)
	channelIdGo := C.GoString(channelId)
	fmt.Println(transportIdGo, channelIdGo)
	if t == nil {
		return
	}
	c := t.Channel(C.GoString(channelId))
	if c == nil {
		return
	}
	c.OnClose(func() {
		fmt.Println("on channel closed:", channelIdGo)
		C.invokeOnChannelClosed(callback)
	})
}

func main() {}
