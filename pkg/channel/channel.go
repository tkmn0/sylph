package channel

type Channel interface {
	SendData(buffer []byte) (int, error)
	SendMessage(message string) (int, error)
	Close()
	Id() string
	OnClose(f func())
	OnError(f func(err error))
	OnMessage(f func(message string))
	OnData(f func(data []byte))
}
