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
	id                     string
	assosiation            *sctp.Association
	sctpStreams            map[string]*stream.SctpStream
	baseStream             stream.Stream
	onChannelHandler       func(c channel.Channel)
	onCloseHandler         func()
	OnTransportInitialized func()
	engines                map[string]*engine.StreamEngine
	close                  chan bool
	streamCount            uint16
	engineConfig           engine.EngineConfig
}

func NewSctpTransport(id string) *SctpTransport {
	return &SctpTransport{
		id:          id,
		sctpStreams: map[string]*stream.SctpStream{},
		engines:     map[string]*engine.StreamEngine{},
		close:       make(chan bool),
		streamCount: 0,
	}
}

func (t *SctpTransport) Init(conn net.Conn, isClient bool, engienConfig engine.EngineConfig) error {
	config := sctp.Config{
		NetConn:       conn,
		LoggerFactory: logging.NewDefaultLoggerFactory(),
	}
	t.engineConfig = engienConfig

	if isClient {
		a, err := sctp.Client(config)
		if err != nil {
			fmt.Println("sctp client creation error", err)
			return err
		}
		t.assosiation = a
		t.openBaseChannel()
	} else {
		a, err := sctp.Server(config)
		if err != nil {
			fmt.Println("sctp server creation error", err)
			return err
		}
		t.assosiation = a
	}

	go func() {
		<-t.close
		t.close = nil
		if t.onCloseHandler != nil {
			t.onCloseHandler()
		}
	}()

	return nil
}

func (t *SctpTransport) AcceptStreamLoop() {
loop:
	for {
		select {
		case <-t.close:
			t.assosiation.Close()
			break loop
		default:
			st, err := t.assosiation.AcceptStream()
			if err != nil {
				fmt.Println(t.id, "stream accept error")
				return
			}
			sctpStream := stream.NewSctpStream(st, t.id)
			e := engine.NewStreamEngine(t.engineConfig)
			e.OnStreamClosed = t.onStreamClosed
			e.OnStream = t.onStreamInitialized
			t.engines[sctpStream.StreamId()] = e
			e.Run(sctpStream, stream.StreamTypeUnKnown, t.engineConfig, t.id)
		}
	}
}

func (t *SctpTransport) openBaseChannel() {
	c := channel.ChannelConfig{
		Unordered:        false,
		ReliabliityType:  channel.ReliabilityTypeReliable,
		ReliabilityValue: 0,
	}
	err := t.openChannel(c, stream.StreamTypeBase)
	if err != nil {
		fmt.Println("error to open base stream", err)
	}
}

func (t *SctpTransport) openChannel(c channel.ChannelConfig, streamType stream.StreamType) error {
	st, err := t.assosiation.OpenStream(t.streamCount, sctp.PayloadTypeWebRTCBinary)
	st.SetReliabilityParams(c.Unordered, byte(c.ReliabliityType), c.ReliabilityValue)
	t.streamCount++

	if err != nil {
		fmt.Println("error open stream", err)
		return err
	}

	sctpStream := stream.NewSctpStream(st, t.id)
	e := engine.NewStreamEngine(t.engineConfig)
	e.OnStreamClosed = t.onStreamClosed
	e.OnStream = t.onStreamInitialized
	t.engines[sctpStream.StreamId()] = e
	e.Run(sctpStream, streamType, t.engineConfig, t.id)
	return nil
}

func (t *SctpTransport) OpenChannel(c channel.ChannelConfig) error {
	return t.openChannel(c, stream.StreamTypeApp)
}

func (t *SctpTransport) onStreamClosed(s stream.Stream) {

	s.CloseStream(t.baseStream.StreamId() != s.StreamId())

	if t.baseStream.StreamId() == s.StreamId() {
		if t.close != nil {
			t.close <- true
		}
	}

	if e, exists := t.engines[s.StreamId()]; exists {
		e.Stop()
		delete(t.engines, s.StreamId())
	}

	delete(t.sctpStreams, s.StreamId())
}

func (t *SctpTransport) onStreamInitialized(st stream.Stream, message engine.InitializeMessage) {
	t.sctpStreams[st.StreamId()] = t.changeStreamToSctpStream(st)
	streamType := stream.StreamType(message.StreamType)
	if streamType == stream.StreamTypeBase {
		t.baseStream = st
	} else if streamType == stream.StreamTypeApp {
		if t.onChannelHandler != nil {
			sctpStream := t.changeStreamToSctpStream(st)
			if sctpStream != nil {
				t.onChannelHandler(sctpStream)
			}
		}
	} else {
		// client recieved
		if t.id == "" {
			t.id = message.TransportId + "-client"
			// id not configurated, this is base stream
			if t.OnTransportInitialized != nil {
				t.OnTransportInitialized()
			}
			t.baseStream = st
		} else {
			// id is already configurated, this is app stream
			if t.onChannelHandler != nil {
				sctpStream := t.changeStreamToSctpStream(st)
				if sctpStream != nil {
					t.onChannelHandler(sctpStream)
				}
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
	t.onChannelHandler = handler
}

func (t *SctpTransport) OnClose(handler func()) {
	t.onCloseHandler = handler
}

func (t *SctpTransport) Channel(id string) channel.Channel {
	if channel, exists := t.sctpStreams[id]; exists {
		return channel
	} else {
		return nil
	}
}

func (t *SctpTransport) SetConfig() {}

func (t *SctpTransport) Close() {
	fmt.Println("close transport:", t.id)

	for _, s := range t.sctpStreams {
		s.Close()
	}

	// if t.baseStream != nil {
	// 	t.changeStreamToSctpStream(t.baseStream).Close()
	// }

	// fmt.Println("stream closed:", t.id)

	// if t.close != nil {
	// 	t.close <- true
	// }
	fmt.Println("channel closed:", t.id)
}

func (t *SctpTransport) IsClosed() bool {
	return t.close == nil
}
