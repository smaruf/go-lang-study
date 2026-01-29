package main

import "fmt"

func main() {
	fmt.Println("Hello world!")
}

// run:
// go build -o hello ./hello.go
// tinygo build -o hello ./hello.go
// tinygo flash -target=microbit ./hello.go
