package main

/*
#include "api.h"
#include <stdbool.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/pkg/channel"
)

// cache


// callbacks
var onTransportCallback C.onTransportCallback = nil
var onTransportClosedCallback C.onTransportClosedCallback = nil
var onChannelCallback C.onChannelCallback = nil
var onChannelClosedCallback C.onChannelClosedCallback = nil
var onChannelErrorCallback C.onChannelErrorCallback = nil
var onMessageCallback C.onMessageCallback = nil
var onDataCallback C.onDataCallback = nil

func init() {
	runtime.LockOSThread()
	clients = map[int]*sylph.Client{}
	servers = map[int]*sylph.Server{}
	channels = map[string]channel.Channel{}
	transports = map[string]sylph.Transport{}
}

func setupChannel(c channel.Channel) {
	c.OnClose(func() {
		// if onChannelCallback != nil {
		// C.invokeOnChannelClosed(C.CString(c.Id()), onChannelClosedCallback)
		// }
		delete(channels, c.Id())
	})

	c.OnError(func(err error) {
		// if onChannelErrorCallback != nil {
		// C.invokeOnChannelError(C.CString(c.Id()), C.CString(err.Error()), onChannelErrorCallback)
		// }
		delete(channels, c.Id())
	})

	c.OnMessage(func(m string) {
		// if onMessageCallback != nil {
		// C.invokeOnMessage(C.CString(c.Id()), C.CString(m), onMessageCallback)
		// }
	})

	c.OnData(func(data []byte) {
		// if onDataCallback != nil {
		// C.invokeOnData(C.CString(c.Id()), unsafe.Pointer(&data[0]), C.int(len(data)), onDataCallback)
		// }
	})
}

func setupTransport(t sylph.Transport, isServer bool) {
	t.OnChannel(func(c channel.Channel) {
		setupChannel(c)
		channels[c.Id()] = c
		// if onChannelCallback != nil {
		// C.invokeOnChannel(C.CString(t.Id()), C.CString(c.Id()), onChannelCallback)
		// }
	})

	t.OnClose(func() {
		// if onTransportClosedCallback != nil {
		// C.invokeOnTransportClosed(C.CString(t.Id()), C.bool(isServer), onTransportClosedCallback)
		// }
		delete(transports, t.Id())
	})
}

//export Initialize
func Initialize() {
	// initializeInternal()
}

//export InitializeClient
func InitializeClient() int {
	id := len(clients)
	client := sylph.NewClient()
	client.OnTransport(func(t sylph.Transport) {
		setupTransport(t, false)
		transports[t.Id()] = t
		// C.invokeOnTransport(C.int(id), C.CString(t.Id()), false, onTransportCallback)
	})

	clients[id] = client
	return id
}

//export Connect
func Connect(id int, address *C.char, port C.int, heartbeatRate int64, timeOutDuration int64) {
	c, exists := clients[id]
	if exists {
		c.Connect(C.GoString(address), int(port), sylph.TransportConfig{
			HeartbeatRateMillisec:   time.Duration(heartbeatRate),
			TimeOutDurationMilliSec: time.Duration(timeOutDuration),
		})
	}
}

//export InitializeServer
func InitializeServer() int {
	id := len(servers)
	server := sylph.NewServer()
	server.OnTransport(func(t sylph.Transport) {
		setupTransport(t, true)
		transports[t.Id()] = t
		if onTransportCallback != nil {
			C.invokeOnTransport(C.int(id), C.CString(t.Id()), true, onTransportCallback)
		}
	})
	servers[id] = server
	return id
}

//export RunServer
func RunServer(id int, address *C.char, port C.int, heartbeatRate int64, timeOutDuration int64) {
	s, exists := servers[id]
	if exists {
		go s.Run(C.GoString(address), int(port), sylph.TransportConfig{
			HeartbeatRateMillisec:   time.Duration(heartbeatRate),
			TimeOutDurationMilliSec: time.Duration(timeOutDuration),
		})
	}
}

//export StopServer
func StopServer(id int) {
	s, exists := servers[id]
	if exists {
		delete(servers, id)
		s.Close()
	}
}

//export OpenChannel
func OpenChannel(id *C.char, unordered bool, reliableType byte, reliableValue uint32) {
	t, exists := transports[C.GoString(id)]
	if exists {
		config := channel.ChannelConfig{
			Unordered:        unordered,
			ReliabliityType:  channel.ReliabilityType(reliableType),
			ReliabilityValue: uint32(reliableValue),
		}
		t.OpenChannel(config)
	}
}

//export CloseTransport
func CloseTransport(id *C.char) {
	t, exists := transports[C.GoString(id)]
	if exists {
		delete(transports, C.GoString(id))
		t.Close()
	}
}

//export CloseChannel
func CloseChannel(id *C.char) {
	c, exists := channels[C.GoString(id)]
	if exists {
		delete(channels, C.GoString(id))
		c.Close()
	}
}

//export SendMessage
func SendMessage(id *C.char, message *C.char) bool {
	c, exists := channels[C.GoString(id)]
	if !exists {
		return false
	}
	_, err := c.SendMessage(C.GoString(message))
	return err == nil
}

//export SendData
func SendData(id *C.char, ptr unsafe.Pointer, length C.int) bool {
	c, exists := channels[C.GoString(id)]
	if !exists {
		return false
	}
	_, err := c.SendData(C.GoBytes(ptr, length))
	return err == nil
}

//export DestroyAll
func DestroyAll() {
	fmt.Println("call destroy all")
	// for _, c := range channels {
	// 	c.Close()
	// }
	// channels = map[string]channel.Channel{}
	fmt.Println("close channels")

	for _, t := range transports {
		t.Close()
	}

	// fmt.Println("close transports")

	for _, c := range clients {
		c.Close()
	}

	fmt.Println("close clients")

	for _, s := range servers {
		s.Close()
	}

	// fmt.Println("close servers")

	// // initializeInternal()

	// fmt.Println("free pointers")
}

//export Panic
func Panic() {
	var panicChan chan bool
	close(panicChan)
}

//export RegisterOnTransportCallback
func RegisterOnTransportCallback(callback C.onTransportCallback) {
	onTransportCallback = callback
}

//export RegisterOnTransportClosedCallback
func RegisterOnTransportClosedCallback(callback C.onTransportClosedCallback) {
	onTransportClosedCallback = callback
}

//export RegisterOnChannelCallback
func RegisterOnChannelCallback(callback C.onChannelCallback) {
	onChannelCallback = callback
}

//export RegisterOnChannelClosedCallback
func RegisterOnChannelClosedCallback(callback C.onChannelClosedCallback) {
	onChannelClosedCallback = callback
}

//export RegisterOnChannelErrorCallback
func RegisterOnChannelErrorCallback(callback C.onChannelErrorCallback) {
	onChannelErrorCallback = callback
}

//export RegisterOnMessageCallback
func RegisterOnMessageCallback(callback C.onMessageCallback) {
	onMessageCallback = callback
}

//export RegisterOnDataCallback
func RegisterOnDataCallback(callback C.onDataCallback) {
	onDataCallback = callback
}

func main() {}
