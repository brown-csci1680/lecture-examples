package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Hello world!")
	// Would like to write
	//send("Hello world!")

	conn, err := net.Dial("tcp4", "127.0.0.1:6666")
	if err != nil {
		// Abort program with stack trace (good enough for now)
		panic(err)
	}

	b, err := conn.Write([]byte("Hello world!"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sent %d bytes", b)
}
