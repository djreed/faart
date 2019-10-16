package main

import "context"

// ./3700recv
func main() {
	ctx := context.TODO()
	if err := receiver(ctx); err != nil {
		panic(err)
	}
}
