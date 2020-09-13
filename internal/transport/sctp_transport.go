package transport

import (
	"fmt"
	"net"

	"github.com/pion/logging"
	"github.com/pion/sctp"
	"github.com/tkmn0/sylph/internal/engine"
	"github.com/tkmn0/sylph/internal/stream"
	"github.com/tkmn0/sylph/pkg/channel"
)

type SctpTransport struct {
	id              string
	assosiation     *sctp.Association
	streams         []stream.Stream
	baseStream      stream.Stream
	onStreamHandler func(c channel.Channel)
	onCloseHanlder  func()
	engine          *engine.StreamEngine
	close           chan bool
	isClosed        bool
	streamCount     uint16
}

func NewSctpTransport(id string) *SctpTransport {
	return &SctpTransport{
		id:          id,
		streams:     []stream.Stream{},
		engine:      engine.StreamNewEngine(),
		close:       make(chan bool),
		isClosed:    false,
		streamCount: 0,
	}
}

func (t *SctpTransport) Init(conn net.Conn, isClient bool) error {
	t.engine.OnStreamClosed = t.onStreamClosed

	config := sctp.Config{
		NetConn:       conn,
		LoggerFactory: logging.NewDefaultLoggerFactory(),
	}

	if isClient {
		a, err := sctp.Client(config)
		if err != nil {
			fmt.Println("sctp client creation error", err)
			return err
		}
		t.assosiation = a
		t.openBaseChannel()
	} else {
		t.engine.OnStream = t.onStreamWithType
		a, err := sctp.Server(config)
		if err != nil {
			fmt.Println("sctp server creation error", err)
			return err
		}
		t.assosiation = a
	}

	go func() {
		<-t.close
		t.isClosed = true
		if t.onCloseHanlder != nil {
			t.onCloseHanlder()
		}
	}()

	return nil
}

func (t *SctpTransport) AcceptStreamLoop() {
	for {
		select {
		case <-t.close:
			t.assosiation.Close()
		default:
			st, err := t.assosiation.AcceptStream()
			if err != nil {
				fmt.Println(t.id, "stream accept error")
				return
			}
			sctpStream := stream.NewSctpStream(st)
			t.engine.Run(sctpStream, stream.StreamTypeUnKnown)
		}
	}
}

func (t *SctpTransport) openBaseChannel() {
	st, err := t.assosiation.OpenStream(t.streamCount, sctp.PayloadTypeWebRTCBinary)
	t.streamCount++

	if err != nil {
		fmt.Println("error open stream", err)
	}

	sctpStream := stream.NewSctpStream(st)
	t.engine.Run(sctpStream, stream.StreamTypeBase)
	t.streams = append(t.streams, sctpStream)
	t.baseStream = sctpStream
}

func (t *SctpTransport) OpenChannel() error {
	st, err := t.assosiation.OpenStream(t.streamCount, sctp.PayloadTypeWebRTCBinary)
	t.streamCount++

	if err != nil {
		return err
	}

	sctpStream := stream.NewSctpStream(st)
	t.engine.Run(sctpStream, stream.StreamTypeApp)
	t.streams = append(t.streams, sctpStream)
	if t.onStreamHandler != nil {
		t.onStreamHandler(sctpStream)
	}
	return nil
}

func (t *SctpTransport) onStreamClosed(s stream.Stream) {
	if t.baseStream.StreamId() == s.StreamId() {
		if !t.isClosed {
			t.close <- true
		}
	}

	s.CloseStream(t.baseStream.StreamId() != s.StreamId())

	streams := []stream.Stream{}

	for _, st := range t.streams {
		if st.StreamId() == s.StreamId() {
			// change remove logic (when array length is 2, outof range)
			// t.streams = append(t.streams[:i], t.streams[i+1:]...)
			continue
		}
		streams = append(streams, st)
	}
	t.streams = streams
}

func (t *SctpTransport) onStreamWithType(st stream.Stream, streamType stream.StreamType) {
	if streamType == stream.StreamTypeBase {
		t.baseStream = st
	}
	t.streams = append(t.streams, st)

	if streamType == stream.StreamTypeApp {
		if t.onStreamHandler != nil {
			sctpStream := t.changeStreamToSctpStream(st)
			if sctpStream != nil {
				t.onStreamHandler(sctpStream)
			}
		}
	}
}

func (t *SctpTransport) changeStreamToSctpStream(i interface{}) *stream.SctpStream {
	sctpStream, ok := i.(*stream.SctpStream)
	if ok {
		return sctpStream
	} else {
		return nil
	}
}

func (t *SctpTransport) Id() string {
	return t.id
}
func (t *SctpTransport) OnChannel(handler func(channel channel.Channel)) {
	t.onStreamHandler = handler
}

func (t *SctpTransport) OnClose(handler func()) {
	t.onCloseHanlder = handler
}

func (t *SctpTransport) Channel(id uint16) {}
func (t *SctpTransport) SetConfig()        {}
func (t *SctpTransport) Close() {
	t.engine.Stop()

	if !t.isClosed {
		t.close <- true
	}
}
