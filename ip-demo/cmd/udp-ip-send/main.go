/*
 * IP-in-UDP sender example
 * This example sends an IP packet inside a UDP packet, producing
 * a packet similar to what you need for the project, including
 * computing the checksum.
 *
 * NOTE:  This example uses hard-coded fields for values in the IP
 * header--you will want to do something different in your project!
 *
 * To run:
 * ./udp-ip-send <bind port> <dest IP> <dest port> <message>
 * where <bind port> is the intended UDP SOURCE PORT of the packet
 *       <dest IP> is the IP address of the host receiving this UDP packet
 *       <dest potr> is the UDP port on the receiving host
 *       <message> is some string to send
 */
package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"

	ipv4header "github.com/brown-csci1680/iptcp-headers"
	"github.com/google/netstack/tcpip/header"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage:  %s <bind port> <dest IP> <dest port> <text>\n", os.Args[0])
		os.Exit(1)
	}

	bindPort := os.Args[1]
	address := os.Args[2]
	port := os.Args[3]
	message := os.Args[4]

	// Turn the address string into a UDPAddr for the connection
	bindAddrString := fmt.Sprintf(":%s", bindPort)
	bindLocalAddr, err := net.ResolveUDPAddr("udp4", bindAddrString)
	if err != nil {
		log.Panicln("Error resolving address:  ", err)
	}

	// Turn the address string into a UDPAddr for the connection
	addrString := fmt.Sprintf("%s:%s", address, port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addrString)
	if err != nil {
		log.Panicln("Error resolving address:  ", err)
	}

	fmt.Printf("Sending to %s:%d\n",
		remoteAddr.IP.String(), remoteAddr.Port)

	// Bind on the local UDP port:  this sets the source port
	// and creates a conn
	conn, err := net.ListenUDP("udp4", bindLocalAddr) // h1 listen on port 5001 for if0
	if err != nil {
		log.Panicln("Dial: ", err)
	}

	// Start filling in the header
	// NOTE:  This example uses hard-coded values for the
	// source, destination, and protocol--you will need to
	// do something different!
	hdr := ipv4header.IPv4Header{
		Version:  4,
		Len:      20, // Header length is always 20 when no IP options
		TOS:      0,
		TotalLen: ipv4header.HeaderLen + len(message),
		ID:       0,
		Flags:    0,
		FragOff:  0,
		TTL:      32,
		Protocol: 0,
		Checksum: 0, // Should be 0 until checksum is computed
		Src:      netip.MustParseAddr("10.0.0.1"),
		Dst:      netip.MustParseAddr("10.1.0.2"),
		Options:  []byte{},
	}

	// Assemble the header into a byte array
	headerBytes, err := hdr.Marshal()
	if err != nil {
		log.Fatalln("Error marshalling header:  ", err)
	}

	// Compute the checksum (see below)
	// Cast back to an int, which is what the Header structure expects
	hdr.Checksum = int(ComputeChecksum(headerBytes)) + 1

	headerBytes, err = hdr.Marshal()
	if err != nil {
		log.Fatalln("Error marshalling header:  ", err)
	}

	// Append header + message into one byte array
	bytesToSend := make([]byte, 0, len(headerBytes)+len(message))
	bytesToSend = append(bytesToSend, headerBytes...)
	bytesToSend = append(bytesToSend, []byte(message)...)

	// Send the message to the "link-layer" addr:port on UDP
	// FOr h1:  send to port 5002
	// ONE CALL TO WriteToUDP => 1 PACKET
	bytesWritten, err := conn.WriteToUDP(bytesToSend, remoteAddr)
	if err != nil {
		log.Panicln("Error writing to socket: ", err)
	}
	fmt.Printf("Sent %d bytes\n", bytesWritten)
}

// Compute the checksum using the netstack package
func ComputeChecksum(b []byte) uint16 {
	checksum := header.Checksum(b, 0)

	// Invert the checksum value.  Why is this necessary?
	// This function returns the inverse of the checksum
	// on an initial computation.  While this may seem weird,
	// it makes it easier to use this same function
	// to validate the checksum on the receiving side.
	// See ValidateChecksum in the receiver file for details.
	checksumInv := checksum ^ 0xffff

	return checksumInv
}
