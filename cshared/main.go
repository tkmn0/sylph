package main

import "C"
import (
	"time"
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/cshared/api"
	"github.com/tkmn0/sylph/pkg/channel"
)

//export Initialize
func Initialize() {
	// initializeInternal()
	api.Initialize()
}

//export InitializeClient
func InitializeClient() uintptr {
	return api.InitializeClient()
}

//export Connect
func Connect(ptr unsafe.Pointer, address *C.char, port C.int, heartbeatRate int64, timeOutDuration int64) {
	api.Connect(ptr, C.GoString(address), int(port), sylph.TransportConfig{
		HeartbeatRateMillisec:   time.Duration(heartbeatRate),
		TimeOutDurationMilliSec: time.Duration(timeOutDuration),
	})
}

//export InitializeServer
func InitializeServer() uintptr {
	return api.InitializeServer()
}

//export RunServer
func RunServer(ptr unsafe.Pointer, address *C.char, port C.int, heartbeatRate int64, timeOutDuration int64) {
	api.RunServer(ptr, C.GoString(address), int(port), sylph.TransportConfig{
		HeartbeatRateMillisec:   time.Duration(heartbeatRate),
		TimeOutDurationMilliSec: time.Duration(timeOutDuration),
	})
}

//export StopServer
func StopServer(ptr unsafe.Pointer) {
	api.StopServer(ptr)
}

//export OpenChannel
func OpenChannel(ptr unsafe.Pointer, unordered bool, reliableType byte, reliableValue uint32) {
	config := channel.ChannelConfig{
		Unordered:        unordered,
		ReliabliityType:  channel.ReliabilityType(reliableType),
		ReliabilityValue: uint32(reliableValue),
	}
	api.OpenChannel(ptr, config)
}

//export CloseTransport
func CloseTransport(ptr unsafe.Pointer) {
	api.CloseTransport(ptr)
}

//export CloseChannel
func CloseChannel(ptr unsafe.Pointer) {
	api.CloseChannel(ptr)
}

//export SendMessage
func SendMessage(ptr unsafe.Pointer, message *C.char) bool {
	return api.SendMessage(ptr, C.GoString(message))
}

//export SendData
func SendData(ptr unsafe.Pointer, dataPtr unsafe.Pointer, length C.int) bool {
	return api.SendData(ptr, C.GoBytes(dataPtr, length))
}

//export Dispose
func Dispose() {
	api.Dispose()
}

//export ObserveOnTransport
func ObserveOnTransport(ptr unsafe.Pointer) uintptr {
	return api.ReadOnTransport(ptr)
}

//export ObserveOnTransportClosed
func ObserveOnTransportClosed(ptr unsafe.Pointer) uintptr {
	return api.ReadOnTransportClosed(ptr)
}

//export ObserveOnChannel
func ObserveOnChannel(ptr unsafe.Pointer) uintptr {
	return api.ReadOnChannel(ptr)
}

//export ObserveOnChannelClosed
func ObserveOnChannelClosed(ptr unsafe.Pointer) uintptr {
	return api.ReadOnChannelClosed(ptr)
}

//export ObserveOnChannelError
func ObserveOnChannelError(ptr unsafe.Pointer) uintptr {
	return api.ReadOnChannelError(ptr)
}

//export ObserveOnChannelMessage
func ObserveOnChannelMessage(ptr unsafe.Pointer) uintptr {
	return api.ReadOnChannelMessage(ptr)
}

//export ObserveOnChannelData
func ObserveOnChannelData(ptr unsafe.Pointer) uintptr {
	return api.ReadOnChannelData(ptr)
}

func main() {}
