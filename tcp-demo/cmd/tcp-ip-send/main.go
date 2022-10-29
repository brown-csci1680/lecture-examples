/*
 * TCP-in-IP-in-UDP sender example
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

	// Start filling in the IP header
	// NOTE:  This example uses hard-coded values for the
	// source, destination, and protocol--you will need to
	// do something different!
	hdr := ipv4.Header{
		Version:  4,
		Len:      20, // Header length is always 20 when no IP options
		TOS:      0,
		TotalLen: ipv4.HeaderLen + header.TCPMinimumSize + len(message),
		ID:       0,
		Flags:    0,
		FragOff:  0,
		TTL:      32,
		Protocol: int(header.TCPProtocolNumber),
		Checksum: 0, // Should be 0 until checksum is computed
		Src:      net.ParseIP("192.168.0.1"),
		Dst:      net.ParseIP("192.168.0.2"),
		Options:  []byte{},
	}

	// Start filling in the TCP header
	// NOTE:  This example uses hard-coded values for the
	// source, destination, and protocol--you will need to
	// do something VERY different!
	tcpHdr := header.TCPFields{
		SrcPort:       12345,
		DstPort:       80,
		SeqNum:        1,
		AckNum:        1,
		DataOffset:    20,
		Flags:         header.TCPFlagSyn | header.TCPFlagAck,
		WindowSize:    65535,
		Checksum:      0,
		UrgentPointer: 0,
	}

	// Assemble the IP header into a byte array
	ipHeaderBytes, err := hdr.Marshal()
	if err != nil {
		log.Fatalln("Error marshalling IP header:  ", err)
	}

	// Compute the IP checksum
	// Cast back to an int, which is what the Header structure expects
	hdr.Checksum = int(ComputeChecksum(ipHeaderBytes))

	ipHeaderBytes, err = hdr.Marshal()
	if err != nil {
		log.Fatalln("Error marshalling IP header:  ", err)
	}

	// Serialize the TCP header
	// NOTE:  This example skips the TCP checksum for now!
	tcpBytes := make(header.TCP, header.TCPMinimumSize)
	tcpBytes.Encode(&tcpHdr)

	// Assemble everything into a single byte array
	bytesToSend := make([]byte, 0, len(ipHeaderBytes)+len(message)+len(tcpBytes))
	bytesToSend = append(bytesToSend, ipHeaderBytes...)
	bytesToSend = append(bytesToSend, tcpBytes...)
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
	// This function returns the inverse of the checksum
	// on an initial computation.  While this may seem weird,
	// it makes it easier to use this same function
	// to validate the checksum on the receiving side.
	// See ValidateChecksum in the receiver file for details.
	checksumInv := checksum ^ 0xffff

	return checksumInv
}
