package sylph

import (
	"fmt"

	"github.com/tkmn0/sylph/internal/engine"

	"github.com/google/uuid"
	"github.com/tkmn0/sylph/internal/transport"
)

type Server struct {
	listenner          *Listener
	listenerConfig     ListenerConfig
	transports         []Transport
	onTransportHandler func(transport Transport)
	close              chan bool
}

func NewServer() *Server {
	return &Server{
		listenner:  NewListenner(),
		transports: []Transport{},
		close:      make(chan bool),
	}
}

func (s *Server) obserbeClose() {
	<-s.close
	s.close = nil
	s.listenner.Close()
}

func (s *Server) Run(address string, port int, tc TransportConfig) {
	go s.obserbeClose()

	c := ListenerConfig{
		address: address,
		port:    port,
	}
	s.listenner.Listen(c)
	for {
		conn := <-s.listenner.connection
		id, err := s.createId()
		if err != nil {
			fmt.Println("id creation error")
			return
		}

		sctp := transport.NewSctpTransport(id)
		err = sctp.Init(conn, false, engine.EngineConfig{
			HeartbeatRateMillisec:   tc.HeartbeatRateMillisec,
			TimeOutDurationMilliSec: tc.TimeOutDurationMilliSec,
		})

		if err != nil {
			fmt.Println("sctp initialize error")
		}

		s.transports = append(s.transports, sctp)
		if s.onTransportHandler != nil {
			s.onTransportHandler(sctp)
		}

		go sctp.AcceptStreamLoop()

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
	if s.close != nil {
		s.close <- true
	}
}
