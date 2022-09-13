package main

import (
	"fmt"
	"golang-sockets/pkg/protocol"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s <port number>", os.Args[0])
	}

	portNumber := os.Args[1]

	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", portNumber))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := conn.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	for {
		msg := protocol.ReadGuessMessage(conn)
		log.Printf("Received guess:  %d", msg.Number)
	}
}
