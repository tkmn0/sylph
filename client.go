// Package sylph implements API to hadle sctp transports.
package sylph

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/tkmn0/sylph/internal/engine"
	"github.com/tkmn0/sylph/internal/transport"
)

type ConnectionState int

const (
	None ConnectionState = iota
	Calling
	Canceled
	TimeOut
	ErrorToReady
	Connected
)

func (s ConnectionState) String() string {
	switch s {
	case None:
		return "None"
	case Calling:
		return "Calling"
	case Canceled:
		return "Canceled"
	case TimeOut:
		return "TimeOut"
	case ErrorToReady:
		return "ErrorToReady"
	case Connected:
		return "Connected"
	}
	return ""
}

// Client handles base connections. (udp, dtls, sctp)
// The relationship Client and Transport is one to one.
type Client struct {
	addr                     *net.UDPAddr
	cancel                   context.CancelFunc
	onTransportHandler       func(t Transport)
	onConnectionStateChanged func(state ConnectionState)
	transports               map[string]Transport
	conn                     *dtls.Conn
	connectionState          ConnectionState
}

// NewClient creates a new Client
func NewClient() *Client {
	return &Client{
		transports: map[string]Transport{},
	}
}

// Connect tries to connect with sylph Server.
// After connection established, OnTransport will be called.
func (c *Client) Connect(address string, port int, tc TransportConfig) {
	// Prepare the IP to connect to
	addr := &net.UDPAddr{IP: net.ParseIP(address), Port: port}

	// Generate a certificate and private key to secure the connection
	certificate, genErr := selfsign.GenerateSelfSigned()

	if genErr != nil {
		c.connectionState = ErrorToReady
		if c.onConnectionStateChanged != nil {
			c.onConnectionStateChanged(c.connectionState)
		}
		return
	}

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	// Connect to a DTLS server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	c.cancel = cancel
	c.connectionState = Calling

	dtlsConn, err := dtls.DialWithContext(ctx, "udp", addr, config)

	if err != nil {
		c.connectionState = TimeOut
		if c.onConnectionStateChanged != nil {
			c.onConnectionStateChanged(c.connectionState)
		}
		return
	}

	c.conn = dtlsConn

	t := transport.NewSctpTransport("")
	t.Init(dtlsConn, true, engine.EngineConfig{
		HeartbeatRateMillisec:   tc.HeartbeatRateMillisec,
		TimeOutDurationMilliSec: tc.TimeOutDurationMilliSec,
	})
	t.OnTransportInitialized = func() {
		c.onTransportHandler(t)
	}
	c.transports[t.Id()] = t
	c.connectionState = Connected
	if c.onConnectionStateChanged != nil {
		c.onConnectionStateChanged(c.connectionState)
	}
}

func (c *Client) ConnectAsync(address string, port int, tc TransportConfig) {
	go c.Connect(address, port, tc)
}

func (c *Client) ConnectionState() ConnectionState {
	return c.connectionState
}

// Transport returns Transport corresponded with id
func (c *Client) Transport(id string) Transport {
	if t, exists := c.transports[id]; exists && !t.IsClosed() {
		return t
	} else {
		return nil
	}
}

// OnTransport handles callback when connection established.
// Callback parameter has Transport.
func (c *Client) OnTransport(handler func(t Transport)) {
	c.onTransportHandler = handler
}

// Close closes dtls connection.
func (c *Client) Close() {
	if c.connectionState == Calling && c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			fmt.Println("dtls closed:", err.Error())
		}
	}
}

func (c *Client) OnConnectionStateChanged(handler func(s ConnectionState)) {
	c.onConnectionStateChanged = handler
}
