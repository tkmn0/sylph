package api

import (
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/cshared/handler"
	"github.com/tkmn0/sylph/pkg/channel"
)

var clients []*sylph.Client
var servers []*sylph.Server
var callbackHandler *handler.CallbackHandler

func Initialize() {
	clients = make([]*sylph.Client, 0)
	servers = make([]*sylph.Server, 0)
	callbackHandler = handler.NewCallbackHandler()
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
	c.Connect(address, port, config)
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
	t := *(*sylph.Transport)(p)
	if t == nil {
		return
	}
	t.OpenChannel(config)
}

func CloseTransport(p unsafe.Pointer) {
	t := *(*sylph.Transport)(p)
	if t == nil {
		return
	}
	t.Close()
}

func CloseChannel(p unsafe.Pointer) {
	c := *(*channel.Channel)(p)
	if c == nil {
		return
	}
	c.Close()
}

func SendMessage(p unsafe.Pointer, message string) bool {
	c := *(*channel.Channel)(p)
	if c == nil {
		return false
	}
	_, err := c.SendMessage(message)
	return err == nil
}

func SendData(p unsafe.Pointer, data []byte) bool {
	c := *(*channel.Channel)(p)
	if c == nil {
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

	clients = []*sylph.Client{}
	servers = []*sylph.Server{}
	callbackHandler = handler.NewCallbackHandler()
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
