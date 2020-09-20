package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tkmn0/sylph/pkg/channel"
)

func Chat(channel channel.Channel) {
	channel.OnMessage(func(message string) {
		fmt.Println("Got message: ", message)
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read console")
		}

		if strings.TrimSpace(text) == "exit" {
			return
		}

		_, err = channel.SendMessage(text)
		if err != nil {
			fmt.Println("error to send message")
		}
	}
}
