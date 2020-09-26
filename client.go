package sylph

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/tkmn0/sylph/internal/engine"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/tkmn0/sylph/internal/transport"
	"github.com/tkmn0/sylph/pkg/util"
)

type Client struct {
	addr               *net.UDPAddr
	cancel             context.CancelFunc
	onTransportHandler func(t Transport)
	transports         map[string]Transport
}

func NewClient() *Client {
	return &Client{
		transports: map[string]Transport{},
	}
}

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

	t := transport.NewSctpTransport("")
	t.Init(dtlsConn, true, engine.EngineConfig{
		HeartbeatRateMisslisec:  tc.HeartbeatRateMisslisec,
		TimeOutDurationMilliSec: tc.TimeOutDurationMilliSec,
	})
	t.OnTransportInitialized = func() {
		c.onTransportHandler(t)
	}
	c.transports[t.Id()] = t
}

func (c *Client) Transport(id string) Transport {
	if t, exists := c.transports[id]; exists && !t.IsClosed() {
		return t
	} else {
		return nil
	}
}

func (c *Client) OnTransport(handler func(t Transport)) {
	c.onTransportHandler = handler
}

func (c *Client) Close() {
	for _, t := range c.transports {
		t.Close()
	}
	c.cancel()
}
