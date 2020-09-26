package main

/*
#include "api.h"
#include <stdint.h>
*/
import "C"
import (
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/pkg/channel"
)

var client *sylph.Client
var server *sylph.Server
var channels map[string]unsafe.Pointer
var transports map[string]unsafe.Pointer

func init() {
	channels = map[string]unsafe.Pointer{}
	transports = map[string]unsafe.Pointer{}
}

// TODO
// Close Transport and realse pointer
// Release method for unity

//export InitializeClient
func InitializeClient() uintptr {
	client = sylph.NewClient()
	return uintptr(unsafe.Pointer(client))
}

//export Connect
func Connect(p unsafe.Pointer, address *C.char, port C.int) {
	c := (*sylph.Client)(p)
	c.Connect(C.GoString(address), int(port))
}

//export InitializeServer
func InitializeServer(address *C.char, port C.int) uintptr {
	server = sylph.NewServer(C.GoString(address), int(port))
	return uintptr(unsafe.Pointer(server))
}

//export RunServer
func RunServer(p unsafe.Pointer) {
	s := (*sylph.Server)(p)
	go s.Run()
}

//export StopServer
func StopServer(p unsafe.Pointer) {
	s := (*sylph.Server)(p)
	s.Close()
}

//export OpenChannel
func OpenChannel(p unsafe.Pointer, unordered bool, reliableType byte, reliableValue uint32) {
	t := *(*sylph.Transport)(p)
	config := channel.ChannelConfig{
		Unordered:        unordered,
		ReliabliityType:  channel.ReliabilityType(reliableType),
		ReliabilityValue: uint32(reliableValue),
	}
	t.OpenChannel(config)
}

//export CloseChannel
func CloseChannel(p unsafe.Pointer) {
	c := *(*channel.Channel)(p)
	delete(channels, c.Id())
	c.Close()
}

//export SendMessage
func SendMessage(p unsafe.Pointer, message *C.char) bool {
	c := *(*channel.Channel)(p)
	_, err := c.SendMessage(C.GoString(message))
	return err == nil
}

//export SendData
func SendData(p unsafe.Pointer, ptr unsafe.Pointer, length C.int) bool {
	c := *(*channel.Channel)(p)
	_, err := c.SendData(C.GoBytes(ptr, length))
	return err == nil
}

//export RegisterOnTransportCallback
func RegisterOnTransportCallback(p unsafe.Pointer, callback C.onTransportCallback, isServer bool) {
	if isServer {
		s := (*sylph.Server)(p)
		s.OnTransport(func(t sylph.Transport) {
			ptr := unsafe.Pointer(&t)
			transports[t.Id()] = ptr
			C.invokeOnTransport(C.uintptr_t(uintptr(ptr)), callback)
		})
	} else {
		c := (*sylph.Client)(p)
		c.OnTransport(func(t sylph.Transport) {
			ptr := unsafe.Pointer(&t)
			transports[t.Id()] = ptr
			C.invokeOnTransport(C.uintptr_t(uintptr(ptr)), callback)
		})
	}
}

//export RegisterOnChannelCallback
func RegisterOnChannelCallback(p unsafe.Pointer, callback C.onChannelCallback) {
	t := *(*sylph.Transport)(p)
	t.OnChannel(func(c channel.Channel) {
		ptr := unsafe.Pointer(&c)
		channels[c.Id()] = ptr
		C.invokeOnChannel(C.uintptr_t(uintptr(ptr)), callback)
	})
}

//export RegisterOnChannelClosedCallback
func RegisterOnChannelClosedCallback(p unsafe.Pointer, callback C.onChannelClosedCallback) {
	c := *(*channel.Channel)(p)
	c.OnClose(func() {
		C.invokeOnChannelClosed(callback)
	})
}

//export RegisterOnChannelErrorCallback
func RegisterOnChannelErrorCallback(p unsafe.Pointer, callback C.onChannelErrorCallback) {
	c := *(*channel.Channel)(p)
	c.OnError(func(err error) {
		C.invokeOnChannelError(C.CString(err.Error()), callback)
	})
}

//export RegisterOnMessageCallback
func RegisterOnMessageCallback(p unsafe.Pointer, callback C.onMessageCallback) {
	c := *(*channel.Channel)(p)
	c.OnMessage(func(message string) {
		C.invokeOnMessage(C.CString(message), callback)
	})
}

//export RegisterOnDataCallback
func RegisterOnDataCallback(p unsafe.Pointer, callback C.onDataCallback) {
	c := *(*channel.Channel)(p)
	c.OnData(func(data []byte) {
		C.invokeOnData(unsafe.Pointer(&data[0]), C.int(len(data)), callback)
	})
}

func main() {}
