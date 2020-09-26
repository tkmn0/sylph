package engine

import (
	"io"
	"sync"
	"time"

	"github.com/tkmn0/sylph/internal/stream"
)

type StreamEngine struct {
	close                   chan bool
	err                     chan error
	builder                 *MessageBuilder
	parcer                  *MessageParcer
	lastHeartbeat           time.Time
	lock                    sync.RWMutex
	heartbeatRateMillisec   time.Duration
	healthCheckRateMillisec time.Duration
	timeOutDurationMillisec time.Duration
	OnStreamClosed          func(stream stream.Stream)
	OnStream                func(stream stream.Stream, streamType stream.StreamType)
}

func NewStreamEngine(config EngineConfig) *StreamEngine {
	return &StreamEngine{
		close:                   make(chan bool),
		err:                     make(chan error),
		builder:                 NewMessageBuilder(),
		parcer:                  NewMessageParcer(),
		heartbeatRateMillisec:   1000,
		healthCheckRateMillisec: 100,
		timeOutDurationMillisec: 300,
	}
}

func (e *StreamEngine) Run(s stream.Stream, t stream.StreamType, config EngineConfig) {
	e.heartbeatRateMillisec = config.HeartbeatRateMisslisec
	e.timeOutDurationMillisec = config.TimeOutDurationMilliSec
	s.OnDataSendHandler(func(data []byte) (int, error) {
		return s.WriteData(e.builder.BuildMessage(data, MessageTypeBody))
	})
	s.OnMessageHandler(func(message string) (int, error) {
		return s.WriteMessage(e.builder.BuildMessage([]byte(message), MessageTypeBody))
	})
	s.OnCloseHandler(func() {
		if e.close != nil {
			e.close <- true
		}
	})

	e.setupStream(s, t)
	e.observeStatus(s)
	go e.handleHeartBeat(s)
	go e.handleHealthCheck()
	go e.readStream(s)
}

func (e *StreamEngine) Stop() {
	if e.close != nil {
		e.close <- true
	}
}

func (e *StreamEngine) setupStream(s stream.Stream, t stream.StreamType) {
	_, err := s.WriteData(e.builder.InitializeMessage(t))
	e.checkError(err)
}

func (e *StreamEngine) observeStatus(s stream.Stream) {
	go func() {
		select {
		case closed := <-e.close:
			if closed {
				if e.OnStreamClosed != nil {
					e.OnStreamClosed(s)
				}
				e.close = nil
			}
		case err := <-e.err:
			if err != nil {
				s.Error(err)
			}
		}
	}()
}

func (e *StreamEngine) handleHeartBeat(s stream.Stream) {
	ticker := time.NewTicker(time.Millisecond * e.heartbeatRateMillisec)
loop:
	for {
		<-ticker.C
		if e.close == nil {
			break loop
		}
		_, err := s.WriteData(e.builder.HeartBeatmessage())
		e.checkError(err)
	}
}

func (e *StreamEngine) readStream(s stream.Stream) {
	for {
		buffer := make([]byte, 1024)
		l, err, isString := s.Read(buffer)
		buffer = buffer[:l]
		invalid := e.checkError(err)
		if invalid {
			return
		}

		mt, buff := e.parcer.Parce(buffer)
		if mt == MessageTypeBody {
			if isString {
				s.Message(string(buff))
			} else {
				s.Data(buff)
			}
		} else if mt == MessageTypeInitialize {
			streamType := stream.StreamType(buff[0])
			if e.OnStream != nil {
				e.OnStream(s, streamType)
			}
		} else if mt == MessageTypeHeartBeat {
			e.lock.Lock()
			e.lastHeartbeat = time.Now()
			e.lock.Unlock()
		}
	}
}

func (e *StreamEngine) handleHealthCheck() {
	ticker := time.NewTicker(time.Millisecond * e.healthCheckRateMillisec)
loop:
	for {
		<-ticker.C
		e.lock.Lock()
		if !e.lastHeartbeat.IsZero() {
			if time.Since(e.lastHeartbeat) > time.Millisecond*(e.heartbeatRateMillisec+e.timeOutDurationMillisec) {
				if e.close != nil {
					e.close <- true
				}
				break loop
			}
		}
		e.lock.Unlock()
	}
}

func (e *StreamEngine) checkError(err error) bool {
	invalid := false
	if err != nil {
		if err == io.EOF {
			if e.close != nil {
				e.close <- true
			}
		} else {
			e.err <- err
		}
		invalid = true
	}
	return invalid
}
