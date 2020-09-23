package main

import (
	"fmt"
	"time"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/pkg/channel"
)

func main() {
	c := sylph.NewClient()
	c.OnTransport(func(t sylph.Transport) {
		fmt.Println("client on transport")
		t.OnChannel(func(c channel.Channel) {
			fmt.Println("client on channel")
			c.OnClose(func() {
				fmt.Println("client channel on close")
			})

			c.OnError(func(err error) {
				fmt.Println("channel on error", err)
			})

			c.OnMessage(func(m string) {
				fmt.Println("channel on message", m)
			})

			c.OnData(func(d []byte) {
				fmt.Println("channel on data", d)
			})

			counter := 0

			go func() {
				for {
					time.Sleep(time.Second)
					c.SendData([]byte("hello from client data"))
					c.SendMessage("hello from client message")
					if counter == 5 {
						c.Close()
						// t.Close()
						// t.OpenChannel(channel.ChannelConfig{})
						break
					}
					counter++
				}
			}()
		})

		t.OnClose(func() {
			fmt.Println("client transport on close")
		})
		c := channel.ChannelConfig{
			Unordered:        false,
			ReliabliityType:  channel.ReliabilityTypeReliable,
			ReliabilityValue: 0,
		}
		t.OpenChannel(c)
		// t.OpenChannel()
	})

	c.Connect("127.0.0.1", 4444)
	for {
	}
}
