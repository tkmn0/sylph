package main

import (
	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/examples/util"
	"github.com/tkmn0/sylph/pkg/channel"
)

func main() {
	client := sylph.NewClient()

	client.OnTransport(func(t sylph.Transport) {
		t.OnChannel(func(ch channel.Channel) {
			util.Chat(ch)
		})

		c := channel.ChannelConfig{
			Unordered:        false,
			ReliabliityType:  channel.ReliabilityTypeReliable,
			ReliabilityValue: 0,
		}
		t.OpenChannel(c)
	})

	client.Connect("127.0.0.1", 4444, sylph.TransportConfig{
		HeartbeatRateMillisec:  1000,
		TimeOutDurationMilliSec: 300,
	})
}
