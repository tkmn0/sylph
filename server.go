package sylph

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tkmn0/sylph/internal/engine"
	"github.com/tkmn0/sylph/internal/listener"
	"github.com/tkmn0/sylph/internal/transport"
)

// Server handles base connections. (udp, dtls, sctp)
// A Server includes a listener.
// A Server handles a bundle of Transports.
// After a Client connected, OnTransport will be called.
type Server struct {
	listener           *listener.Listener
	listenerConfig     listener.ListenerConfig
	transports         []Transport
	onTransportHandler func(transport Transport)
	close              chan bool
}

// NewServer creates a Server.
func NewServer() *Server {
	return &Server{
		listener:   listener.NewListener(),
		transports: []Transport{},
		close:      make(chan bool),
	}
}

// obserbeClose obserbes and handles closing
func (s *Server) obserbeClose() {
	<-s.close
	s.close = nil
	s.listener.Close()
}

// Run runs server with address, port, and TransportConfig.
// This will block process, call this with goroutine when necessary.
func (s *Server) Run(address string, port int, tc TransportConfig) {
	go s.obserbeClose()

	c := listener.ListenerConfig{
		Address: address,
		Port:    port,
	}
	s.listener.Listen(c)
	for {
		conn := <-s.listener.Connection
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

// createId creates id for Transport.
// This id will be used server side and client.
// The id in client side has suffix "-client" with server side id.
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

// OnTransport will be called when Client connected.
func (s *Server) OnTransport(handler func(t Transport)) {
	s.onTransportHandler = handler
}

// Close closes server.
func (s *Server) Close() {
	if s.close != nil {
		s.close <- true
	}
}
