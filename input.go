package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const _bold = "\x1b[1m"
const _normal = "\x1b[0m"

var _prompt = fmt.Sprintf("%sgibshell $ %s", _bold, _normal)

type Input struct{}

func (i Input) On(callback func([]byte)) {
	go func() {
		for {
			fmt.Printf("%s", _prompt)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.Replace(input, "\n", "", -1)
			if input == "" {
				continue
			}
			callback([]byte(input))
		}
	}()
}
