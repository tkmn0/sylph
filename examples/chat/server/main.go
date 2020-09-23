package main

import (
	"fmt"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/examples/util"
	"github.com/tkmn0/sylph/pkg/channel"
)

func main() {
	hub := util.NewHub()
	s := sylph.NewServer("127.0.0.1", 4444)

	s.OnTransport(func(t sylph.Transport) {
		t.OnChannel(func(c channel.Channel) {
			hub.Register(c)
		})
		t.OnClose(func() {
			fmt.Println("transport closed")
		})
	})

	go s.Run()
	hub.Chat()
}
