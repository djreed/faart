package main

// ./3700recv
func main() {
	if err := receiver(); err != nil {
		panic(err)
	}
}
