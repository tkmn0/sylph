package sylph

import (
	"github.com/tkmn0/sylph/pkg/channel"
)

type Transport interface {
	OpenChannel(config channel.ChannelConfig) error
	OnChannel(handler func(channel channel.Channel))
	OnClose(handler func())
	Close()
	Id() string
	Channel(id uint16)
	SetConfig()
}
