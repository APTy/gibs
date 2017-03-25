package main

import (
	"os/exec"
	"strings"
)

func runCmd(cmdString string) string {
	parsed := parseCmdString(cmdString)
	name, args := parsed[0], parsed[1:]

	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func parseCmdString(str string) []string {
	return strings.Fields(str)
}
