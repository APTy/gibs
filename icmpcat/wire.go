package icmpcat

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type connMode int

const (
	clientMode connMode = iota
	serverMode
)

type wire struct {
	mode    connMode
	host    *net.IPAddr
	conn    *icmp.PacketConn
	crypter Crypter
	seq     int
}

func newWire(mode connMode, _conn *icmp.PacketConn, crypter Crypter) (*wire, error) {
	return &wire{
		mode:    mode,
		conn:    _conn,
		crypter: crypter,
		seq:     seqInit,
	}, nil
}

func (w *wire) setHost(hostIP string) error {
	host := net.ParseIP(hostIP)
	if host == nil {
		return fmt.Errorf("failed to parse IP: %v", hostIP)
	}
	w.host = &net.IPAddr{IP: host}
	return nil
}

func (w *wire) write(p *packet) error {
	bytes, err := p.toBytes()
	if err != nil {
		return err
	}

	echoType := ipv4.ICMPTypeEcho
	if w.mode == serverMode {
		echoType = ipv4.ICMPTypeEchoReply
	}

	data := w.crypter.Encrypt(bytes)
	msg, err := newEcho(echoType, data, w.seq)
	if err != nil {
		return err
	}

	// log.Printf("send: %x", msg)
	if _, err := w.conn.WriteTo(msg, w.host); err != nil {
		return err
	}
	w.seq++
	return nil
}

func (w *wire) read() (*packet, net.Addr, error) {
	buf := make([]byte, 1500)
	n, peer, err := w.conn.ReadFrom(buf)
	if err != nil {
		log.Printf("sock read error: %v", err)
		return nil, nil, err
	}
	msg, err := parseEcho(buf, n)
	if err != nil {
		return nil, nil, err
	}
	res, err := w.crypter.Decrypt(msg)
	if err != nil {
		log.Printf("decrypt error: %v", err)
		return nil, nil, err
	}
	// log.Printf("recv: %x", msg)
	p, err := fromBytes(res)
	if err != nil {
		log.Printf("parse error: %v", err)
		return nil, nil, err
	}
	return p, peer, nil
}
