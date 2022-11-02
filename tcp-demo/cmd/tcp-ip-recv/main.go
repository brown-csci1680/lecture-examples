/*
 * TCP-in-IP-in-UDP listener example
 * To run:
 *  ./ip-tcp-recv <bind port>
 *  where <bind port> is the port on which to receive packets.
 */
package main

import (
	"fmt"
	"ip-demo/pkg/iptcp_utils"
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

		// ***************** PARSE IP HEADER *****************************
		// Marshal the received byte array into a UDP header
		// NOTE:  This does not validate the checksum or check any fields
		// (You'll need to do this part yourself)
		hdr, err := ipv4.ParseHeader(buffer)

		if err != nil {
			fmt.Println("Error parsing header", err)
			continue
		}

		ipHeaderSize := hdr.Len

		// Validate the IP checksum
		ipHeaderBytes := buffer[:ipHeaderSize]

		ipChecksumFromHeader := uint16(hdr.Checksum)
		ipComputedChecksum := iptcp_utils.ValidateIPChecksum(ipHeaderBytes,
			ipChecksumFromHeader)

		var ipChecksumState string
		if ipComputedChecksum == ipChecksumFromHeader {
			ipChecksumState = "OK"
		} else {
			ipChecksumState = "FAIL"
		}

		// This is just a demo, so we should only be seeing TCP packets,
		// drop everything else
		if hdr.Protocol != int(header.TCPProtocolNumber) {
			fmt.Println("Packet is not a TCP packet, skipping")
			continue
		}

		// ******************** PARSE TCP HEADER ************************
		// Next, get the TCP header

		// **** IMPORTANT ****:  The total length of the data is included
		// in the **IP header**.  This is very important because
		// ReadFromUDP reads into a buffer of size 1400, but the actual
		// message may be smaller!
		// Therefore, to get the correct-sized payload, we need
		// to slice it out of buffer
		tcpHeaderAndData := buffer[ipHeaderSize:hdr.TotalLen]

		// Parse the TCP header into a struct
		tcpHdr := iptcp_utils.ParseTCPHeader(tcpHeaderAndData)

		// Get the payload
		tcpPayload := tcpHeaderAndData[tcpHdr.DataOffset:]

		// Now that we have all the pieces, we can verify the TCP checksum
		// In general, the checksum function expects the checksum field to be
		// set to 0, which allows us to verify it by checking against the
		// value sent in the header.
		// An alternative is to *not* clear this value and then compare
		// tcpComputedChecksum == 0 (for details, see EdStem #208)
		tcpChecksumFromHeader := tcpHdr.Checksum // Save original
		tcpHdr.Checksum = 0
		tcpComputedChecksum := iptcp_utils.ComputeTCPChecksum(&tcpHdr, hdr.Src, hdr.Dst, tcpPayload)

		var tcpChecksumState string
		if tcpComputedChecksum == tcpChecksumFromHeader {
			tcpChecksumState = "OK"
		} else {
			tcpChecksumState = "FAIL"
		}
		// Finally, print everything out
		fmt.Printf("Received TCP packet from %s\nIP Header:  %v\nIP Checksum:  %s\nTCP header:  %+v\nFlags:  %s\nTCP Checksum:  %s\nPayload (%d bytes):  %s\n",
			sourceAddr.String(), hdr, ipChecksumState, tcpHdr, iptcp_utils.TCPFlagsAsString(tcpHdr.Flags), tcpChecksumState, len(tcpPayload), string(tcpPayload))
	}
}
