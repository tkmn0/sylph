package channel

type ReliabilityType byte

const (
	// ReliabilityTypeReliable is used for reliable transmission
	ReliabilityTypeReliable ReliabilityType = iota
	// ReliabilityTypeRexmit is used for partial reliability by retransmission count
	ReliabilityTypeRexmit
	// ReliabilityTypeTimed is used for partial reliability by retransmission duration
	ReliabilityTypeTimed
)

type ChannelConfig struct {
	Unordered        bool
	ReliabliityType  ReliabilityType
	ReliabilityValue uint32
}
