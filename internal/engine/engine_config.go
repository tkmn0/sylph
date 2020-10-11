package engine

import (
	"time"
)

type EngineConfig struct {
	HeartbeatRateMillisec  time.Duration
	TimeOutDurationMilliSec time.Duration
}
