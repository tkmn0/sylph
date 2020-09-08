package stream

type StreamType uint8

const (
	StreamTypeBase StreamType = iota
	StreamTypeApp
	StreamTypeUnKnown
)
