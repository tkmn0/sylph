package main

import (
	"fmt"
	"time"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/pkg/channel"
)

func main() {
	s := sylph.NewServer("127.0.0.1", 4444)
	s.OnTransport = func(t sylph.Transport) {
		fmt.Println("transport opened")
		t.OnChannel(func(c channel.Channel) {
			fmt.Println("server on channel")

			c.OnClose(func() {
				fmt.Println("channel on close")
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
					c.SendData([]byte("hello data from server"))
					c.SendMessage("hello message from server")
					if counter == 5 {
						// c.Close()
						// t.Close()
						break
					}
					counter++
				}
			}()
		})
		t.OnClose(func() {
			fmt.Println("transport closed")
		})
	}
	s.Run()
}
