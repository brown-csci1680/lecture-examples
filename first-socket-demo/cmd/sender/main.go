package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Hello world!")
	// Would like to write
	//send("Hello world!")

	conn, err := net.Dial("udp4", "127.0.0.1:6666")
	if err != nil {
		// Abort program with stack trace (good enough for now)
		panic(err)
	}

	// Read a line from the terminal and then
	fmt.Printf("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Wait until user has typed a line
	line := scanner.Text()

	b, err := conn.Write([]byte(line))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sent %d bytes", b)
}
