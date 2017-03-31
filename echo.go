package main

import (
	"errors"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// IANA protocol number assigned to ICMP IPv4
const protoICMP = 1

func newEcho(typ ipv4.ICMPType, data []byte, seq int) ([]byte, error) {
	wm := icmp.Message{
		Type: typ, Code: 0,
		Body: &icmp.Echo{
			ID: 12345 & 0xffff, Seq: 1,
			Data: data,
		},
	}
	return wm.Marshal(nil)
}

func parseEcho(msg []byte, n int) ([]byte, error) {
	rm, err := icmp.ParseMessage(protoICMP, msg[:n])
	if err != nil {
		return nil, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply, ipv4.ICMPTypeEcho:
		echo, ok := rm.Body.(*icmp.Echo)
		if !ok {
			return nil, errors.New("failed to parse echo reply")
		}
		return echo.Data, nil
	default:
		return nil, errors.New("unknown message type")
	}
}
