package engine

import (
	"fmt"
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
	streams               map[string]stream.Stream
	lock                  sync.RWMutex
}

func StreamNewEngine() *StreamEngine {
	return &StreamEngine{
		close:                 make(chan bool),
		err:                   make(chan error),
		builder:               NewMessageBuilder(),
		parcer:                NewMessageParcer(),
		heartbeatRateMillisec: 1000,
		streams:               map[string]stream.Stream{},
	}
}

func (e *StreamEngine) Run(s stream.Stream, t stream.StreamType) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.streams[s.StreamId()] = s
	s.OnDataSendHandler(func(data []byte) (int, error) {
		return s.WriteData(e.builder.BuildMessage(data, MessageTypeBody))
	})
	s.OnMessageHandler(func(message string) (int, error) {
		return s.WriteMessage(e.builder.BuildMessage([]byte(message), MessageTypeBody))
	})
	s.OnCloseHandler(func() {
		fmt.Println("OnClosedHandler in engine", s.StreamId())
		e.lock.Lock()
		defer e.lock.Unlock()
		delete(e.streams, s.StreamId())
		if e.OnStreamClosed != nil {
			e.OnStreamClosed(s)
		}
	})

	e.setupStream(s, t)
	go e.handleHeartBeat()
	go e.readStream()
}

func (e *StreamEngine) Stop() {
	fmt.Println("engine stop called")
	if !e.isClosed {
		e.close <- true
	}
}

func (e *StreamEngine) setupStream(s stream.Stream, t stream.StreamType) {
	_, err := s.WriteData(e.builder.InitializeMessage(t))
	if err != nil {
		e.lock.Lock()
		defer e.lock.Unlock()
		delete(e.streams, s.StreamId())
	}
}

func (e *StreamEngine) handleHeartBeat() {
	ticker := time.NewTicker(time.Millisecond * e.heartbeatRateMillisec)
loop:
	for {
		select {
		case <-e.close:
			e.isClosed = true

			if e.OnStreamClosed != nil {
				// fmt.Println("OnStreamClosed:", s.StreamId())
				// e.OnStreamClosed(s)
				for _, stream := range e.streams {
					e.OnStreamClosed(stream)
				}
			}
			break loop
		// case err := <-e.err:
		// 	// s.Error(err)
		// 	break loop
		case <-ticker.C:
			if e.isClosed {
				break loop
			}
			e.lock.Lock()
			defer e.lock.Unlock()

			for _, stream := range e.streams {
				_, err := stream.WriteData(e.builder.HeartBeatmessage())
				// e.checkError(err)
				if err != nil {
					delete(e.streams, stream.StreamId())
					if e.OnStreamClosed != nil {
						e.OnStreamClosed(stream)
					}
				}
			}
		}
	}
}

func (e *StreamEngine) readStream() {
loop:
	for {
		select {
		case <-e.close:
			if e.OnStreamClosed != nil {
				e.lock.Lock()
				defer e.lock.Unlock()
				for _, stream := range e.streams {
					e.OnStreamClosed(stream)
				}
			}
			break loop
		// case err := <-e.err:
		// 	// s.Error(err)
		// 	break loop
		default:
			e.lock.Lock()
			defer e.lock.Unlock()
			for _, s := range e.streams {

				buffer := make([]byte, 1024)
				l, err, isString := s.Read(buffer)
				buffer = buffer[:l]
				if err != nil {
					delete(e.streams, s.StreamId())
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
				fmt.Println("time exeeded, close")
				if !e.isClosed {
					e.close <- true
				}
				break loop
			}
			e.Unlock()
		}
	}
}

/*
func (e *StreamEngine) checkError(err error) bool {
	invalid := false
	if err != nil {
		if err == io.EOF {
			if !e.isClosed {
				fmt.Println("check error close engine")
				e.close <- true
			}
		} else {
			e.err <- err
		}
		invalid = true
	}
	return invalid
}
*/
