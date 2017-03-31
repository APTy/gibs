package main

import (
	"errors"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// IANA protocol number assigned to ICMP IPv4
const protoICMP = 1

func newEcho(data []byte, seq int) ([]byte, error) {
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEchoReply, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: seq,
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
	case ipv4.ICMPTypeEchoReply:
		echo, ok := rm.Body.(*icmp.Echo)
		if !ok {
			return nil, errors.New("failed to parse echo reply")
		}
		return echo.Data, nil
	default:
		return nil, errors.New("unknown message type")
	}
}
