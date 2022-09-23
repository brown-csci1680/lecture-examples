/*
 * UDP listener example
 *
 * This code examples demonstrates one way to
 * create a UDP socket and wait for messages.
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
		// Read from the UDP socket.  Note that we need to provide a buffer
		// as large, or larger, than the biggest message we want to receive
		// On the Internet, packets are generally < 1400 bytes (we'll learn why later)
		buffer := make([]byte, MaxMessageSize)
		bytesRead, sourceAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Panicln("Error reading from UDP socket ", err)
		}

		// Print our data to stdout
		// ***NOTE:  DO NOT DO THIS IN SNOWCAST***
		// In Snowcast, the data the listener receives is not
		// an ASCII string--instead, it's *binary* data from, eg.
		// an mp3 file.  Instead, you should write it directly
		// to stdout using something like io.Write.
		// Why?  Print/Printf/Println will interpret this data
		// like a string, and therefore might add newlines or
		// react to other formatting, which will corrupt
		// the data you are trying to output (and thus your
		// rate will be wrong!)
		message := string(buffer)
		fmt.Printf("Received %d bytes from %s:  %s\n",
			bytesRead, sourceAddr.String(), message)
	}
}
