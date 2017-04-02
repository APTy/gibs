package icmpcat

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/icmp"
)

const (
	secret   = "Ezv27ceMoBruteP4gh1c6Kebs69J4F5KhJNIewmGJxY="
	icmpIPv4 = "ip4:icmp"
	localIfc = "0.0.0.0"
)

// ICMPCat allows for read/write access to a remote host
// over ICMP to a node also running ICMPCat.
type ICMPCat interface {

	// Open a connection to the remote host.
	Connect(string) error

	// Accept will listen for a new connection.
	Accept() error

	// Close destroys the connection.
	Close() error

	// Write data over the connection.
	Send(io.Reader) error

	// OnReceive registers a callback to invoke with received messages.
	OnReceive(func(io.Reader))

	// Listen blocks to receive incoming messages.
	Listen()
}

// New returns an object for sending/receiving data over ICMP.
func NewV2() (ICMPCat, error) {
	conn, err := icmp.ListenPacket(icmpIPv4, localIfc)
	if err != nil {
		return nil, fmt.Errorf("failed to establish ICMP connection: %v", err)
	}
	crypter, err := NewCrypter(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypter: %v", err)
	}
	return &icmpCatV2{
		_conn:   conn,
		crypter: crypter,
	}, nil
}

type icmpCatV2 struct {
	_conn    *icmp.PacketConn
	conn     *conn
	crypter  Crypter
	callback func(io.Reader)
}

func (c *icmpCatV2) Connect(hostIP string) error {
	conn, err := newClientConn(c._conn, c.crypter, hostIP)
	if err != nil {
		return err
	}
	if conn.sendHello(); err != nil {
		return err
	}
	if err := conn.recvHello(); err != nil {
		return err
	}
	c.conn = conn
	return nil

}

func (c *icmpCatV2) Accept() error {
	conn, err := newServerConn(c._conn, c.crypter)
	if err != nil {
		return err
	}
	if err := conn.acceptHello(); err != nil {
		return err
	}
	if err := conn.ackHello(); err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *icmpCatV2) Close() error {
	return nil
}

func (c *icmpCatV2) Send(r io.Reader) error {
	if c.conn.isServer() {
		c.conn.waitForSendRequest()
	}
	buf := make([]byte, 1350)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Send err: %v", err)
			return err
		}
		c.conn.sendDataStream(buf[:n])
	}
	c.conn.sendDataEOF()
	return nil
}

func (c *icmpCatV2) OnReceive(callback func(io.Reader)) {
	c.callback = callback
}

func (c *icmpCatV2) Listen() {
	log.Printf("Listening for data streams")
	c.conn.onData(func(r io.Reader) {
		c.callback(r)
	})
	select {}
}
