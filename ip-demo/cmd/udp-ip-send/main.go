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
	"os"

	"github.com/google/netstack/tcpip/header"
	"golang.org/x/net/ipv4"
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
	conn, err := net.ListenUDP("udp4", bindLocalAddr)
	if err != nil {
		log.Panicln("Dial: ", err)
	}

	// Start filling in the header
	// NOTE:  This example uses hard-coded values for the
	// source, destination, and protocol--you will need to
	// do something different!
	hdr := ipv4.Header{
		Version:  4,
		Len:      20, // Header length is always 20 when no IP options
		TOS:      0,
		TotalLen: ipv4.HeaderLen + len(message),
		ID:       0,
		Flags:    0,
		FragOff:  0,
		TTL:      32,
		Protocol: 0,
		Checksum: 0, // Should be 0 until checksum is computed
		Src:      net.ParseIP("192.168.0.1"),
		Dst:      net.ParseIP("192.168.0.2"),
		Options:  []byte{},
	}

	// Assemble the header into a byte array
	headerBytes, err := hdr.Marshal()
	if err != nil {
		log.Fatalln("Error marshalling header:  ", err)
	}

	// Compute the checksum (see below)
	// Cast back to an int, which is what the Header structure expects
	hdr.Checksum = int(ComputeChecksum(headerBytes))

	headerBytes, err = hdr.Marshal()
	if err != nil {
		log.Fatalln("Error marshalling header:  ", err)
	}

	bytesToSend := make([]byte, 0, len(headerBytes)+len(message))
	bytesToSend = append(bytesToSend, headerBytes...)
	bytesToSend = append(bytesToSend, []byte(message)...)

	// Send the message to the "link-layer" addr:port on UDP
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
	// The checksum function in the library we're using seems
	// to have been built to plug into some other software that expects
	// to receive the complement of this value.
	// The reasons for this are unclear to me at the moment, but for now
	// take my word for it.  =)
	checksumInv := checksum ^ 0xffff

	return checksumInv
}
