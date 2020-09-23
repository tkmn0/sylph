package sylph

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tkmn0/sylph/internal/transport"
)

type Server struct {
	listenner          *Listener
	ListenerConfig     *ListenerConfig
	transports         []Transport
	onTransportHandler func(transport Transport)
	close              chan bool
	isClosed           bool
}

func NewServer(address string, port int) *Server {
	l := NewListenner()
	c := &ListenerConfig{
		address: address,
		port:    port,
	}
	return &Server{
		listenner:      l,
		ListenerConfig: c,
		transports:     []Transport{},
		close:          make(chan bool),
	}
}

func (s *Server) Run() {
	s.listenner.Listen(s.ListenerConfig)
loop:
	for {
		select {
		case conn := <-s.listenner.connection:
			id, err := s.createId()
			if err != nil {
				fmt.Println("id creation error")
				return
			}

			sctp := transport.NewSctpTransport(id)
			err = sctp.Init(conn, false)

			go sctp.AcceptStreamLoop()

			if err != nil {
				fmt.Println("sctp initialize error")
			}

			s.transports = append(s.transports, sctp)
			if s.onTransportHandler != nil {
				s.onTransportHandler(sctp)
			}
		case <-s.close:
			s.listenner.Close()
			break loop
		}
	}
}

func (s *Server) createId() (string, error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuid := uuidObj.String()

	for _, sctp := range s.transports {
		if sctp.Id() == uuid {
			return s.createId()
		}
	}
	return uuid, nil
}

func (s *Server) OnTransport(handler func(t Transport)) {
	s.onTransportHandler = handler
}

func (s *Server) Close() {
	if !s.isClosed {
		s.isClosed = true
		s.close <- true
	}
}
