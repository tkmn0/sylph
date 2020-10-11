package stream

type Stream interface {
	Error(e error)
	Close()
	CloseStream(notify bool)
	Message(s string)
	Data(b []byte)
	Read(buffer []byte) (int, error, bool)
	WriteData(buffer []byte) (int, error)
	WriteMessage(buffer []byte) (int, error)
	StreamId() string
	OnDataSendHandler(handler func(data []byte) (int, error))
	OnMessageHandler(handler func(message string) (int, error))
	OnCloseHandler(handler func())
}
