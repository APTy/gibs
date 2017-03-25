package main

import "flag"

type flags struct {
	isBindShell bool
	host        string
	cmd         string
}

func parseFlags() flags {
	var isBindShell bool
	var host string
	var cmd string
	flag.BoolVar(&isBindShell, "bind-shell", false, "whether or not to run as a bind shell")
	flag.StringVar(&host, "host", "", "the remote address of the host shell")
	flag.StringVar(&cmd, "cmd", "", "the command to run")
	flag.Parse()
	return flags{
		isBindShell: isBindShell,
		host:        host,
		cmd:         cmd,
	}
}
