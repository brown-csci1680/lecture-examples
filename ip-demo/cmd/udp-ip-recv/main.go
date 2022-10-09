/*
 * IP-in-UDP listener example
 * To run:
 *  ./udp-ip-recv <bind port>
 *  where <bind port> is the port on which to receive packets.
 */
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/ipv4"
)

const (
	MaxMessageSize = 1400
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:  %s <port>\n", os.Args[0])
		os.Exit(1)
	}

	port := os.Args[1]

	// To read from a UDP socket, we need to bind it to the port
	// on which we want to receive data

	// Get the address structure for the address on which we want to listen
	listenString := fmt.Sprintf(":%s", port)
	listenAddr, err := net.ResolveUDPAddr("udp4", listenString)
	if err != nil {
		log.Panicln("Error resolving address:  ", err)
	}

	// Create a socket and bind it to the port on which we want to receive data
	conn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		log.Panicln("Could not bind to UDP port: ", err)
	}

	for {
		buffer := make([]byte, MaxMessageSize)

		// Read on the UDP port
		_, sourceAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Panicln("Error reading from UDP socket ", err)
		}

		// Marshal the received byte array into a UDP header
		// NOTE:  This does not validate the checksum or check any fields
		// (You'll need to do this part yourself)
		hdr, err := ipv4.ParseHeader(buffer)

		// Find the start of the actual message (which occurs after the header)
		headerSize := hdr.Len
		message := buffer[headerSize:]

		if err != nil {
			fmt.Println("Error parsing header", err)
		}

		fmt.Printf("Received IP packetfrom %s\nHeader:  %v\nMessage:  %s\n",
			sourceAddr.String(), hdr, string(message))
	}
}
