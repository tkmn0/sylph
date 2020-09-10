package sylph

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

type ListenerConfig struct {
	address string
	port    int
}

type Listener struct {
	addr       *net.UDPAddr
	connection chan net.Conn
	closeCh    chan bool
	cancel     context.CancelFunc
	isClosed   bool
}

func NewListenner() *Listener {
	return &Listener{
		connection: make(chan net.Conn),
		closeCh:    make(chan bool),
	}
}

func (l *Listener) Listen(c *ListenerConfig) {
	// Prepare the IP to connect to
	l.addr = &net.UDPAddr{IP: net.ParseIP(c.address), Port: c.port}

	// Generate a certificate and private key to secure the connection
	certificate, genErr := selfsign.GenerateSelfSigned()
	if genErr != nil {
		log.Fatalln("generate self signed is failed")
	}

	// Create parent context to cleanup handshaking connections on exit.
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
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
	if err != nil {
		fmt.Println("dtls liten error", err)
	}

	// defer func() {
	// 	util.Check(listener.Close())
	// }()

	go func() {
	loop:
		for {
			select {
			case <-l.closeCh:
				fmt.Println("close listener in channel")
				l.isClosed = true
				l.cancel()
				err := listener.Close()
				fmt.Println(err)
				break loop
			default:
				// Wait for a connection.
				if listener == nil {
					return
				}
				conn, err := listener.Accept()
				// util.Check(err)
				if err != nil && !l.isClosed {
					close(l.closeCh)
				}
				l.connection <- conn
			}
		}
	}()
}

func (l *Listener) Close() {
	fmt.Println("close litener", l.isClosed)
	close(l.closeCh)
}
