package engine

import (
	"encoding/json"
)

type MessageBuilder struct{}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (b *MessageBuilder) BuildMessage(buffer []byte, t MessageType) []byte {
	payload := []byte{uint8(t)}
	return append(payload[:], buffer[:]...)
}

func (b *MessageBuilder) HeartBeatmessage() []byte {
	return []byte{uint8(MessageTypeHeartBeat)}
}

func (b *MessageBuilder) InitializeMessage(m InitializeMessage) []byte {
	bytes, _ := json.Marshal(&m)
	return append([]byte{uint8(MessageTypeInitialize)}, bytes[:]...)
}
