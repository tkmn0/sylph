package engine

type MessageType uint8

const (
	MessageTypeHeartBeat MessageType = iota
	MessageTypeBody
	MessageTypeChunk
	MessageTypeUnknown
	MessageTypeInitialize
)

type InitializeMessage struct {
	StreamType  uint8  `json:"stream_type"`
	TransportId string `json:"transport_id"`
}
