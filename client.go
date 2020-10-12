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
	"github.com/tkmn0/sylph/pkg/util"
)

// Client handles base connections. (udp, dtls, sctp)
// The relationship Client and Transport is one to one.
type Client struct {
	addr               *net.UDPAddr
	cancel             context.CancelFunc
	onTransportHandler func(t Transport)
	transports         map[string]Transport
	conn               *dtls.Conn
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
	util.Check(genErr)
	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	// Connect to a DTLS server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	c.cancel = cancel

	dtlsConn, err := dtls.DialWithContext(ctx, "udp", addr, config)
	util.Check(err)
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
	c.cancel()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			fmt.Println("dtls closed:", err.Error())
		}
	}
}
