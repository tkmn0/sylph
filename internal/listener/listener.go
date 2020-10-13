package listener

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

// Listener is dtls listener.
// This handles udp and dtls.
type Listener struct {
	addr       *net.UDPAddr
	Connection chan net.Conn
	closeCh    chan bool
	cancel     context.CancelFunc
	listener   net.Listener
}

func NewListener() *Listener {
	return &Listener{
		Connection: make(chan net.Conn),
		closeCh:    make(chan bool),
	}
}

func (l *Listener) obserbeClose() {
	<-l.closeCh
	l.closeCh = nil
	if l.cancel != nil {
		l.cancel()
	}
	if l.listener != nil {
		l.listener.Close()
	}
	l.listener = nil
}

func (l *Listener) Listen(c ListenerConfig) {

	go l.obserbeClose()

	// Prepare the IP to connect to
	l.addr = &net.UDPAddr{IP: net.ParseIP(c.Address), Port: c.Port}

	// Generate a certificate and private key to secure the connection
	certificate, genErr := selfsign.GenerateSelfSigned()
	if genErr != nil {
		log.Fatalln("generate self signed is failed")
	}

	// Create parent context to cleanup handshaking connections on exit.
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, 30*time.Second)
		},
	}

	// Listen
	listener, err := dtls.Listen("udp", l.addr, config)
	l.listener = listener

	if err != nil {
		fmt.Println("dtls liten error", err)
	}

	go func() {
		for {
			// Wait for a connection.
			if l.listener == nil {
				break
			}
			conn, err := l.listener.Accept()
			defer func() {
				if conn != nil {
					err := conn.Close()
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}()

			if err != nil {
				fmt.Println("listener error:", err.Error())
				if l.closeCh != nil {
					l.closeCh <- true
				}
				break
			}
			l.Connection <- conn
		}
	}()
}

func (l *Listener) Close() {
	if l.closeCh != nil {
		l.closeCh <- true
	}
}
