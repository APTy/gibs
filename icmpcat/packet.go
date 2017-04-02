package icmpcat

import (
	"bytes"
	"encoding/gob"
)

const (
	protoID uint16 = 0x4655

	ackType  uint8 = 0
	connType uint8 = 1
	dataType uint8 = 2

	ackCode uint8 = 0

	connClientHelloCode uint8 = 0
	connServerHelloCode uint8 = 1
	connGoodbyeCode     uint8 = 2

	dataStreamCode uint8 = 0
	dataEOFCode    uint8 = 1
	dataReqCode    uint8 = 2
)

// packet is the structure of the ICMP body data
// it consists of an 8-byte header followed by data
type packet struct {
	ProtoID uint16
	TypeID  uint8
	Code    uint8
	Seq     uint32
	Data    []byte
}

func fromBytes(data []byte) (*packet, error) {
	p := new(packet)
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *packet) toBytes() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(*p)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
