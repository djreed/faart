package main

import (
	"context"
	"os"
)

func main() {
	ctx := context.TODO()
	go server(ctx, "localhost:8080")
	go client(ctx, "localhost:8080", os.Stdin)
	for {
	}
}
