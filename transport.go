package sylph

import (
	"github.com/tkmn0/sylph/pkg/channel"
)

// Transport is interface for transport.
// Transport handles Channels.
// A Transport handles a bundle of Channels.
type Transport interface {
	OpenChannel(config channel.ChannelConfig) error
	OnChannel(handler func(channel channel.Channel))
	OnClose(handler func())
	Close()
	Id() string
	Channel(id string) channel.Channel
	SetConfig()
	IsClosed() bool
}
