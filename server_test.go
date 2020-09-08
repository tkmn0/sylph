package sylph_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/pkg/channel"
)

func TestServerConnection(test *testing.T) {
	fmt.Println("TestServerConnection")
	closeCh := make(chan bool)

	s := sylph.NewServer("127.0.0.1", 4444)
	defer s.Close()
	s.OnTransport = func(t sylph.Transport) {
		fmt.Println("server on transport")
		t.OnChannel(func(c channel.Channel) {
			fmt.Println("server on channel")
			c.OnMessage(func(m string) {
				fmt.Println("server on message", m)
			})
		})
	}
	go s.Run()

	c := sylph.NewClient()
	c.OnTransport = func(t sylph.Transport) {
		fmt.Println("client on transport")
		t.OnChannel(func(c channel.Channel) {
			fmt.Println("client on channel")
			counter := 0
		loop:
			for {
				time.Sleep(time.Millisecond)
				c.SendMessage("hello")
				counter++
				if counter > 3 {
					close(closeCh)
					break loop
				}
			}
		})
		t.OpenChannel()
	}
	c.Connect("127.0.0.1", 4444)

	{
		<-closeCh
		s.Close()
		time.Sleep(time.Second * 5)
		fmt.Println("closed test")
	}
}

func TestChannelClose(test *testing.T) {
	fmt.Println("TestChannelClose")
	s := sylph.NewServer("127.0.0.1", 4444)
	s.OnTransport = func(t sylph.Transport) {
		fmt.Println("server on transport")
		t.OnChannel(func(c channel.Channel) {
			fmt.Println("server on channel")
			c.OnMessage(func(m string) {
				fmt.Println("server on message", m)
			})
		})
	}
	go s.Run()

	c := sylph.NewClient()
	c.OnTransport = func(t sylph.Transport) {
		fmt.Println("client on transport")
		t.OnChannel(func(c channel.Channel) {
			fmt.Println("client on channel")
			counter := 0
		loop:
			for {
				time.Sleep(time.Millisecond)
				c.SendMessage("hello")
				counter++
				if counter > 3 {
					c.Close()
					break loop
				}
			}
		})
		if err := t.OpenChannel(); err != nil {
			fmt.Println("open first channel erorr")
		}

		if err := t.OpenChannel(); err != nil {
			fmt.Println("open second channel erorr", err)
		}
	}
	c.Connect("127.0.0.1", 4444)

}
