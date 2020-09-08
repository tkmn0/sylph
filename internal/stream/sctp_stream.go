package stream

import (
	"strconv"

	"github.com/pion/sctp"
)

type SctpStream struct {
	stream             *sctp.Stream
	onCloseHandler     func()
	onErrorHandler     func(err error)
	onMessageHandler   func(message string)
	onDataHandler      func(data []byte)
	dataSendHandler    func(data []byte) (int, error)
	messageSendHandler func(message string) (int, error)
	streamCloseHandler func()
	isClosed           bool
}

func NewSctpStream(stream *sctp.Stream) *SctpStream {
	return &SctpStream{stream: stream}
}

func (s *SctpStream) id() string {
	return strconv.Itoa(int(s.stream.StreamIdentifier()))
}

// ChannelInterface
func (s *SctpStream) SendData(buffer []byte) (int, error) {
	return s.dataSendHandler(buffer)
}

func (s *SctpStream) SendMessage(message string) (int, error) {
	return s.messageSendHandler(message)
}

func (s *SctpStream) Close() {
	s.streamCloseHandler()
}

func (s *SctpStream) Id() string {
	return s.id()
}

func (s *SctpStream) OnClose(f func()) {
	s.onCloseHandler = f
}

func (s *SctpStream) OnError(f func(err error)) {
	s.onErrorHandler = f
}

func (s *SctpStream) OnMessage(f func(message string)) {
	s.onMessageHandler = f
}

func (s *SctpStream) OnData(f func(data []byte)) {
	s.onDataHandler = f
}

// StreamInterface
func (s *SctpStream) WriteData(buffer []byte) (int, error) {
	return s.stream.WriteSCTP(buffer, sctp.PayloadTypeWebRTCBinary)
}

func (s *SctpStream) WriteMessage(buffer []byte) (int, error) {
	return s.stream.WriteSCTP(buffer, sctp.PayloadTypeWebRTCString)
}

func (s *SctpStream) Read(buffer []byte) (int, error, bool) {
	l, i, err := s.stream.ReadSCTP(buffer)
	isString := func(identifier sctp.PayloadProtocolIdentifier) bool {
		if identifier == sctp.PayloadTypeWebRTCString || identifier == sctp.PayloadTypeWebRTCStringEmpty {
			return true
		} else {
			return false
		}
	}(i)
	return l, err, isString
}

func (s *SctpStream) Error(e error) {
	if s.onErrorHandler != nil {
		s.onErrorHandler(e)
	}
	s.isClosed = true
}

func (s *SctpStream) CloseStream(notify bool) {
	if s.onCloseHandler != nil && notify && !s.isClosed {
		s.isClosed = true
		s.stream.Close()
		s.onCloseHandler()
	}
}

func (s *SctpStream) Message(m string) {
	if s.onMessageHandler != nil && !s.isClosed {
		s.onMessageHandler(m)
	}
}

func (s *SctpStream) Data(d []byte) {
	if s.onDataHandler != nil && !s.isClosed {
		s.onDataHandler(d)
	}
}

func (s *SctpStream) StreamId() string {
	return s.id()
}

func (s *SctpStream) OnDataSendHandler(handler func(data []byte) (int, error)) {
	s.dataSendHandler = handler
}

func (s *SctpStream) OnMessageHandler(handler func(message string) (int, error)) {
	s.messageSendHandler = handler
}

func (s *SctpStream) OnCloseHandler(handler func()) {
	s.streamCloseHandler = handler
}
