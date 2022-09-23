/*
 * UDP sender example
 *
 * This code examples demonstrates one way to
 * create a UDP socket and send messages to
 * a specific IP and port.
 *
 * Note:  this example just sends and receives string
 * data--you will need to do something different
 * in your listener.
 *
 * See the comments for details and notes on how
 * this might apply to your projects.
 */
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage:  %s <address> <port> <text>\n", os.Args[0])
		os.Exit(1)
	}

	address := os.Args[1]
	port := os.Args[2]
	message := os.Args[3]

	// Turn the address string into a UDPAddr for the connection
	addrString := fmt.Sprintf("%s:%s", address, port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addrString)
	if err != nil {
		log.Panicln("Error resolving address:  ", err)
	}

	fmt.Printf("Sending to %s:%d\n",
		remoteAddr.IP.String(), remoteAddr.Port)

	// Create a UDPConn to use for sending data
	// NOTE:  Unlike TCP, this doesn't actually send any packets
	// to establish a connection!
	// This just creates the socket in the OS
	conn, err := net.DialUDP("udp4", nil, remoteAddr)
	if err != nil {
		log.Panicln("Dial: ", err)
	}

	// Send the message over the socket
	// This will immediately send one UDP packet of
	// size len(message) bytes
	bytesWritten, err := conn.Write([]byte(message))
	if err != nil {
		log.Panicln("Error writing to socket: ", err)
	}
	fmt.Printf("Sent %d bytes\n", bytesWritten)
}
