package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tkmn0/sylph/pkg/channel"
)

type Hub struct {
	channels map[string]channel.Channel
	lock     sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{channels: make(map[string]channel.Channel)}
}

func (h *Hub) Register(channel channel.Channel) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.channels[channel.Id()] = channel
	channel.OnMessage(func(message string) {
		h.broadCastMessage(message, channel.Id())
	})

	channel.OnData(func(data []byte) {
		h.broadCastData(data, channel.Id())
	})

	channel.OnClose(func() {
		h.unregister(channel)
	})

	channel.OnError(func(err error) {
		h.unregister(channel)
	})
}

func (h *Hub) unregister(channel channel.Channel) {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.channels, channel.Id())
}

func (h *Hub) broadCastMessage(m string, from string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	fmt.Println("broad cast message:", m, "from:", from)
	for _, channel := range h.channels {
		if channel.Id() != from {
			_, err := channel.SendMessage(m)
			if err != nil {
				fmt.Println("failed to send message to:", channel.Id())
			}
		}
	}
}

func (h *Hub) broadCastData(d []byte, from string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for _, channel := range h.channels {
		if channel.Id() != from {
			_, err := channel.SendData(d)
			if err != nil {
				fmt.Println("failed to send message to:", channel.Id())
			}
		}
	}
}

func (h *Hub) Chat() {
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read string from console")
		}

		if strings.TrimSpace(msg) == "exit" {
			return
		}
		h.broadCastMessage(msg, "server")
	}
}
