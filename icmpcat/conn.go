package icmpcat

import (
	"io"
	"log"
	"net"

	"golang.org/x/net/icmp"
)

const seqInit = 1

type conn struct {
	wire      *wire
	packetSeq uint32
	acks      chan uint32
	sendReqs  chan bool
}

func newClientConn(_conn *icmp.PacketConn, crypter Crypter, hostIP string) (*conn, error) {
	wire, err := newWire(clientMode, _conn, crypter)
	if err != nil {
		return nil, err
	}
	if wire.setHost(hostIP); err != nil {
		return nil, err
	}
	return &conn{
		wire:      wire,
		packetSeq: seqInit,
		acks:      make(chan uint32, 1024),
		sendReqs:  make(chan bool, 1),
	}, nil
}

func newServerConn(_conn *icmp.PacketConn, crypter Crypter) (*conn, error) {
	wire, err := newWire(serverMode, _conn, crypter)
	if err != nil {
		return nil, err
	}
	return &conn{
		wire:      wire,
		packetSeq: seqInit,
		acks:      make(chan uint32, 1024),
		sendReqs:  make(chan bool, 1),
	}, nil
}

func (c *conn) send(p packet) error {
	p.ProtoID = protoID
	p.Seq = c.packetSeq
	if err := c.wire.write(&p); err != nil {
		return err
	}
	return nil
}

func (c *conn) recv() (*packet, net.Addr, error) {
	return c.wire.read()
}

func (c *conn) sendHello() error {
	log.Printf("Initializing connection to %s", c.wire.host)
	if err := c.send(packet{
		TypeID: connType,
		Code:   connClientHelloCode,
	}); err != nil {
		return err
	}
	log.Printf("Sent client hello")
	return nil
}

func (c *conn) recvHello() error {
	log.Printf("Listening for server reply")
	for {
		p, _, err := c.recv()
		if err != nil {
			continue
		}
		if p.TypeID != connType && p.Code != connServerHelloCode {
			log.Printf("Got packet(Type: %d Code %d)", p.TypeID, p.Code)
			continue
		}
		log.Printf("Got server hello")
		return nil
	}
}

func (c *conn) acceptHello() error {
	log.Printf("Listening for client hello")
	for {
		p, peer, err := c.recv()
		if err != nil {
			continue
		}
		log.Printf("Received packet from client")
		if p.TypeID != connType && p.Code != connClientHelloCode {
			log.Printf("Not a client hello, ignoring")
			continue
		}
		if err := c.wire.setHost(peer.String()); err != nil {
			return err
		}
		return nil
	}
}

func (c *conn) ackHello() error {
	if err := c.send(packet{
		TypeID: connType,
		Code:   connServerHelloCode,
	}); err != nil {
		return err
	}
	log.Printf("Sent server hello")
	return nil
}

func (c *conn) sendDataStream(data []byte) error {
	log.Printf("Sending data stream packet, len: %d", len(data))
	if err := c.send(packet{
		TypeID: dataType,
		Code:   dataStreamCode,
		Data:   data,
	}); err != nil {
		log.Printf("Send err: %v", err)
		return err
	}
	<-c.acks
	return nil
}

func (c *conn) sendDataEOF() error {
	log.Printf("Sending data EOF packet")
	if err := c.send(packet{
		TypeID: dataType,
		Code:   dataEOFCode,
	}); err != nil {
		log.Printf("Send err: %v", err)
		return err
	}
	<-c.acks
	if c.isServer() {
		return nil
	}
	log.Printf("Sending data request packet")
	if err := c.send(packet{
		TypeID: dataType,
		Code:   dataReqCode,
	}); err != nil {
		log.Printf("Send err: %v", err)
		return err
	}
	return nil
}

func (c *conn) isServer() bool {
	return c.wire.mode == serverMode
}

func (c *conn) waitForSendRequest() {
	<-c.sendReqs
}

func (c *conn) sendAck() error {
	if err := c.send(packet{
		TypeID: ackType,
		Code:   ackCode,
	}); err != nil {
		log.Printf("Send err: %v", err)
		return err
	}
	return nil
}

func (c *conn) onData(callback func(io.Reader)) {
	stream := newStream()
	receivingData := false
	for {
		p, _, err := c.recv()
		if err != nil {
			continue
		}
		if p.TypeID == connType {
			log.Printf("ignoring connection type")
			continue
		}
		if p.TypeID == ackType {
			log.Printf("Received ACK")
			c.acks <- p.Seq
			continue
		}
		if !receivingData {
			receivingData = true
			go callback(stream)
		}
		if p.Code == dataStreamCode {
			log.Printf("Received data stream")
			stream.Write(p.Data)
			c.sendAck()
			continue
		}
		if p.Code == dataReqCode {
			log.Printf("Received data request")
			c.sendReqs <- true
			continue
		}
		log.Printf("Received data EOF")
		stream.Close()
		receivingData = false
		stream = newStream()
		c.sendAck()
	}
}
