package main

import (
	"os"
)

// ./3700send <recv_host>:<recv_port>
func main() {
	if len(os.Args) != 2 {
		panic("Must pass in a single argument: <recv_host>:<recv_port>\nData will be read from STDIN")
	}

	target := os.Args[1]
	source := os.Stdin
	if err := sender(target, source); err != nil {
		panic(err)
	}
}
