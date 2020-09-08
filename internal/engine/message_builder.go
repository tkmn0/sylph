package engine

import "github.com/tkmn0/sylph/internal/stream"

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

func (b *MessageBuilder) InitializeMessage(t stream.StreamType) []byte {
	return []byte{uint8(MessageTypeInitialize), uint8(t)}
}
