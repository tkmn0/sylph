package engine

type MessageType uint8

const (
	MessageTypeHeartBeat MessageType = iota
	MessageTypeBody
	MessageTypeChunk
	MessageTypeUnknown
	MessageTypeInitialize
)
