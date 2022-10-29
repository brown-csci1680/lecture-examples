/*
 * TCP-in-IP-in-UDP listener example
 * To run:
 *  ./ip-tcp-recv <bind port>
 *  where <bind port> is the port on which to receive packets.
 */
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/netstack/tcpip/header"
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

		if err != nil {
			fmt.Println("Error parsing header", err)
			continue
		}

		headerSize := hdr.Len

		// Validate the checksum.
		headerBytes := buffer[:headerSize]

		checksumFromHeader := uint16(hdr.Checksum)
		computedChecksum := ValidateChecksum(headerBytes, checksumFromHeader)
		var checksumState string
		if computedChecksum == checksumFromHeader {
			checksumState = "OK"
		} else {
			checksumState = "FAIL"
		}

		if hdr.Protocol != int(header.TCPProtocolNumber) {
			fmt.Println("Packet is not a TCP packet, skipping")
			continue
		}

		// Next, get the TCP header
		// NOTE:  This does NOT validate the TCP checksum
		tcpHeaderAndData := buffer[headerSize:]
		tcpHdr := ParseTCPHeader(tcpHeaderAndData)

		// Finally, get the message
		message := tcpHeaderAndData[tcpHdr.DataOffset:]

		// Finally, print everything out
		fmt.Printf("Received IP packet from %s\nIP Header:  %v\nIP Checksum:  %s\nTCP header:  %+v\nMessage:  %s\n",
			sourceAddr.String(), hdr, checksumState, tcpHdr, string(message))
	}
}

// Build a TCPFields struct from the TCP byte array
// NOTE: the netstack package might have other options for parsing the header
// that you may like better--this example is most similar to our other class examples.
// Your mileage may vary!
func ParseTCPHeader(b []byte) header.TCPFields {
	td := header.TCP(b)
	return header.TCPFields{
		SrcPort:    td.SourcePort(),
		DstPort:    td.DestinationPort(),
		SeqNum:     td.SequenceNumber(),
		AckNum:     td.AckNumber(),
		DataOffset: td.DataOffset(),
		Flags:      td.Flags(),
		WindowSize: td.WindowSize(),
		Checksum:   td.Checksum(),
	}
}

// Validate the checksum using the netstack package
// Here, we provide both the byte array for the header AND
// the initial checksum value that was stored in the header
//
// "Why don't we need to set the checksum value to 0 first?"
//
// Normally, the checksum is computed with the checksum field
// of the header set to 0.  This library creatively avoids
// this step by instead subtracting the initial value from
// the computed checksum.
// If you use a different language or checksum function, you may
// need to handle this differently.
func ValidateChecksum(b []byte, fromHeader uint16) uint16 {
	checksum := header.Checksum(b, fromHeader)

	return checksum
}
