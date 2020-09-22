package engine

import (
	"io"
	"sync"
	"time"

	"github.com/tkmn0/sylph/internal/stream"
)

type StreamEngine struct {
	close         chan bool
	err           chan error
	builder       *MessageBuilder
	parcer        *MessageParcer
	lastHeartbeat time.Time
	isStarted     bool
	isClosed      bool
	sync.Mutex
	heartbeatRateMillisec time.Duration
	OnStreamClosed        func(stream stream.Stream)
	OnStream              func(stream stream.Stream, streamType stream.StreamType)
}

func NewStreamEngine() *StreamEngine {
	return &StreamEngine{
		close:                 make(chan bool),
		err:                   make(chan error),
		builder:               NewMessageBuilder(),
		parcer:                NewMessageParcer(),
		heartbeatRateMillisec: 1000,
	}
}

func (e *StreamEngine) Run(s stream.Stream, t stream.StreamType) {
	s.OnDataSendHandler(func(data []byte) (int, error) {
		return s.WriteData(e.builder.BuildMessage(data, MessageTypeBody))
	})
	s.OnMessageHandler(func(message string) (int, error) {
		return s.WriteMessage(e.builder.BuildMessage([]byte(message), MessageTypeBody))
	})
	s.OnCloseHandler(func() {
		if !e.isClosed {
			e.close <- true
		}
	})

	e.setupStream(s, t)
	go e.handleHeartBeat(s)
	go e.readStream(s)
}

func (e *StreamEngine) Stop() {
	if !e.isClosed {
		e.close <- true
	}
}

func (e *StreamEngine) setupStream(s stream.Stream, t stream.StreamType) {
	_, err := s.WriteData(e.builder.InitializeMessage(t))
	e.checkError(err)
}

func (e *StreamEngine) handleHeartBeat(s stream.Stream) {
	ticker := time.NewTicker(time.Millisecond * e.heartbeatRateMillisec)
loop:
	for {
		select {
		case <-e.close:
			e.isClosed = true

			if e.OnStreamClosed != nil {
				e.OnStreamClosed(s)
			}
			break loop
		case err := <-e.err:
			s.Error(err)
			break loop
		case <-ticker.C:
			if e.isClosed {
				break loop
			}
			_, err := s.WriteData(e.builder.HeartBeatmessage())
			e.checkError(err)
		}
	}
}

func (e *StreamEngine) readStream(s stream.Stream) {
loop:
	for {
		select {
		case <-e.close:
			if e.OnStreamClosed != nil {
				e.OnStreamClosed(s)
			}
			break loop
		case err := <-e.err:
			s.Error(err)
			break loop
		default:
			buffer := make([]byte, 1024)
			l, err, isString := s.Read(buffer)
			buffer = buffer[:l]
			invalid := e.checkError(err)
			if invalid {
				return
			}

			if !e.isStarted {
				e.lastHeartbeat = time.Now()
				go e.handleHealthCheck()
				e.isStarted = true
			}
			e.Lock()
			e.lastHeartbeat = time.Now()
			e.Unlock()

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
			}
		}
	}
}

func (e *StreamEngine) handleHealthCheck() {
loop:
	for {
		select {
		case <-e.close:
			break loop
		case <-e.err:
			break loop
		default:
			e.Lock()
			if time.Since(e.lastHeartbeat) > time.Millisecond*(e.heartbeatRateMillisec+300) {
				if !e.isClosed {
					e.close <- true
				}
				break loop
			}
			e.Unlock()
		}
	}
}

func (e *StreamEngine) checkError(err error) bool {
	invalid := false
	if err != nil {
		if err == io.EOF {
			if !e.isClosed {
				e.close <- true
			}
		} else {
			e.err <- err
		}
		invalid = true
	}
	return invalid
}
