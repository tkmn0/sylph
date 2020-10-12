package sylph

import (
	"time"
)

// TransportConfig is cofig for transport.
// HeartbeatRateMillisec is rate for heart beat. Heartbeat sends 1 byte header.
// Server and Client send heartbeat to detect the ohter side is running.
// TimeOutDurationMillisec is time out duration to detect the other side is living.
type TransportConfig struct {
	HeartbeatRateMillisec   time.Duration
	TimeOutDurationMilliSec time.Duration
}
