package icmpcat

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	secret   = "Ezv27ceMoBruteP4gh1c6Kebs69J4F5KhJNIewmGJxY="
	icmpIPv4 = "ip4:icmp"
	localIfc = "0.0.0.0"
	seqInit  = 1
	mtu      = 1400
)

// ICMPCat allows for read/write access to a remote host
// over ICMP to a node also running ICMPCat.
type ICMPCat interface {

	// Send a slice of bytes to the remote host.
	Send(ipv4.ICMPType, []byte, string) error

	// OnReceive registers a callback to invoke with received messages.
	OnReceive(func(*net.IPAddr, []byte))

	// Listen blocks to receive incoming messages.
	Listen()
}

// New returns an object for sending/receiving data over ICMP.
func New() (ICMPCat, error) {
	conn, err := icmp.ListenPacket(icmpIPv4, localIfc)
	if err != nil {
		return nil, fmt.Errorf("failed to establish ICMP connection: %v", err)
	}
	crypter, err := NewCrypter(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypter: %v", err)
	}
	return &icmpCat{
		conn:    conn,
		crypter: crypter,
		seq:     seqInit,
	}, nil
}

type icmpCat struct {
	conn     *icmp.PacketConn
	crypter  Crypter
	seq      int
	callback func(*net.IPAddr, []byte)
}

func (c *icmpCat) Send(typ ipv4.ICMPType, b []byte, hostIP string) error {
	host := net.ParseIP(hostIP)
	if host == nil {
		return fmt.Errorf("failed to parse IP: %v", hostIP)
	}
	ip := &net.IPAddr{IP: host}

	for i := 0; i <= len(b)/mtu; i++ {
		start := i * mtu
		end := (i*mtu + mtu)
		if end > len(b) {
			end = len(b)
		}
		data := c.crypter.Encrypt(b[start:end])
		// log.Printf("sent %x", data)
		msg, err := newEcho(typ, data, c.seq)
		if err != nil {
			return err
		}

		if _, err := c.conn.WriteTo(msg, ip); err != nil {
			return err
		}
		c.seq++
	}
	return nil
}

func (c *icmpCat) OnReceive(callback func(*net.IPAddr, []byte)) {
	c.callback = callback
}

func (c *icmpCat) Listen() {
	for {
		buf := make([]byte, 1500)
		n, peer, err := c.conn.ReadFrom(buf)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		msg, err := parseEcho(buf, n)
		if err != nil {
			continue
		}
		// log.Printf("got %x", msg)
		res, err := c.crypter.Decrypt(msg)
		if err != nil {
			continue
		}
		ipPeer, _ := peer.(*net.IPAddr)
		c.callback(ipPeer, res)
	}
}
