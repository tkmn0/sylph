package channel

// PayloadProtocolIdentifier is an enum for DataChannel payload types
type StreamType uint32

const (
	StreamTypeString StreamType = 51
	StreamTypeBinary StreamType = 53
)
