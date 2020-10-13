package api

import (
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/cshared/handler"
	"github.com/tkmn0/sylph/pkg/channel"
)

var clients []*sylph.Client
var servers []*sylph.Server
var transports map[uintptr]sylph.Transport
var channels map[uintptr]channel.Channel
var callbackHandler *handler.CallbackHandler

func Initialize() {
	clients = make([]*sylph.Client, 0)
	servers = make([]*sylph.Server, 0)
	transports = make(map[uintptr]sylph.Transport)
	channels = make(map[uintptr]channel.Channel)
	callbackHandler = handler.NewCallbackHandler()

	callbackHandler.OnTransport(func(t sylph.Transport, p uintptr) {
		transports[p] = t
	})

	callbackHandler.OnChannel(func(c channel.Channel, p uintptr) {
		channels[p] = c
	})
}

func InitializeClient() uintptr {
	client := sylph.NewClient()
	clients = append(clients, client)
	callbackHandler.SetupClientEvents(client)
	return uintptr(unsafe.Pointer(client))
}

func Connect(p unsafe.Pointer, address string, port int, config sylph.TransportConfig) {
	c := (*sylph.Client)(p)
	if c == nil {
		return
	}
	c.ConnectAsync(address, port, config)
}

func InitializeServer() uintptr {
	server := sylph.NewServer()
	servers = append(servers, server)
	callbackHandler.SetupServerEvents(server)
	return uintptr(unsafe.Pointer(server))
}

func RunServer(p unsafe.Pointer, address string, port int, config sylph.TransportConfig) {
	s := (*sylph.Server)(p)
	if s == nil {
		return
	}
	go s.Run(address, port, config)
}

func StopServer(p unsafe.Pointer) {
	s := (*sylph.Server)(p)
	if s == nil {
		return
	}
	s.Close()
}

func OpenChannel(p unsafe.Pointer, config channel.ChannelConfig) {
	t, exists := transports[uintptr(p)]

	if !exists {
		return
	}
	t.OpenChannel(config)
}

func CloseTransport(p unsafe.Pointer) {
	t, exists := transports[uintptr(p)]

	if !exists {
		return
	}
	t.Close()
}

func CloseChannel(p unsafe.Pointer) {
	c, exists := channels[uintptr(p)]
	if !exists {
		return
	}
	c.Close()
}

func SendMessage(p unsafe.Pointer, message string) bool {
	c, exists := channels[uintptr(p)]
	if !exists {
		return false
	}
	_, err := c.SendMessage(message)
	return err == nil
}

func SendData(p unsafe.Pointer, data []byte) bool {
	c, exists := channels[uintptr(p)]
	if !exists {
		return false
	}
	_, err := c.SendData(data)
	return err == nil
}

func Dispose() {
	for _, c := range clients {
		c.Close()
	}

	for _, s := range servers {
		s.Close()
	}

	for _, c := range channels {
		c.Close()
	}

	for _, t := range transports {
		t.Close()
	}

	clients = []*sylph.Client{}
	servers = []*sylph.Server{}
}

func ReadConnectionState(p unsafe.Pointer) int {
	c := (*sylph.Client)(p)
	if c == nil {
		return -1
	}
	return int(c.ConnectionState())
}

func ReadOnTransport(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnTransport(p)
}

func ReadOnTransportClosed(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnTransportClosed(p)
}

func ReadOnChannel(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnChannel(p)
}

func ReadOnChannelClosed(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnChannelClosed(p)
}

func ReadOnChannelError(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnChannelError(p)
}

func ReadOnChannelMessage(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnChannelMessage(p)
}

func ReadOnChannelData(p unsafe.Pointer) uintptr {
	return callbackHandler.ReadOnChannelData(p)
}
