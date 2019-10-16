package main

import (
	"context"
	"os"
)

// ./3700send <recv_host>:<recv_port>
func main() {
	if len(os.Args) != 2 {
		panic("Must pass in a single argument: <recv_host>:<recv_port>")
	}

	ctx := context.TODO()
	target := os.Args[1]
	source := os.Stdin
	if err := sender(ctx, target, source); err != nil {
		panic(err)
	}
}
