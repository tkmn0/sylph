package sylph

import (
	"time"
)

type TransportConfig struct {
	HeartbeatRateMisslisec  time.Duration
	TimeOutDurationMilliSec time.Duration
}
