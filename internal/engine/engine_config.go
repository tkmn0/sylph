package engine

import (
	"time"
)

type EngineConfig struct {
	HeartbeatRateMisslisec  time.Duration
	TimeOutDurationMilliSec time.Duration
}
