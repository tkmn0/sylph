package sylph

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/tkmn0/sylph/internal/transport"
	"github.com/tkmn0/sylph/pkg/util"
)

type Client struct {
	addr           *net.UDPAddr
	cancel         context.CancelFunc
	OnTransport    func(t Transport)
	sctpTransports []*transport.SctpTransport
}

func NewClient() *Client {
	return &Client{
		sctpTransports: []*transport.SctpTransport{},
	}
}

func (c *Client) Connect(address string, port int) {
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

	t := transport.NewSctpTransport("testtest")
	t.Init(dtlsConn, true)
	c.sctpTransports = append(c.sctpTransports, t)

	if c.OnTransport != nil {
		c.OnTransport(t)
	}
}
