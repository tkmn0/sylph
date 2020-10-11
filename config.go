package sylph

import (
	"time"
)

type TransportConfig struct {
	HeartbeatRateMillisec   time.Duration
	TimeOutDurationMilliSec time.Duration
}
