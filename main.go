package main

import "log"

func main() {
	flags := parseFlags()
	cli, err := NewCLI()
	if err != nil {
		log.Fatalf("failed to create CLI: %v", err)
	}

	if flags.isBindShell {
		cli.BindShell()
		return
	}

	if flags.host == "" {
		log.Fatal("missing host or cmd")
	}
	cli.SendCmd(flags.cmd, flags.host)
}
