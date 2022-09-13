package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang-sockets/pkg/protocol"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("Usage:  %s <address> <port number>", os.Args[0])
	}

	address := os.Args[1]
	portNumber := os.Args[2]

	addrToUse := fmt.Sprintf("%s:%s", address, portNumber)
	conn, err := net.Dial("tcp4", addrToUse)
	defer conn.Close()

	if err != nil {
		log.Fatalln("Error connecting:  ", err)
	}

	msg := &protocol.GuessMessage{MessageType: protocol.MessageTypeGuess, Number: 42}
	b, err := conn.Write(msg.Marshal())

	if err != nil {
		log.Fatalln("Error writing:  ", err)
	}

	fmt.Printf("Wrote %d bytes\n", b)
}
