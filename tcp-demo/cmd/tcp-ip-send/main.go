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
	"net/netip"
	"os"
	"tcp-demo/pkg/iptcp_utils"

	ipv4header "github.com/brown-csci1680/iptcp-headers"
	"github.com/google/netstack/tcpip/header"
	//"golang.org/x/net/ipv4"
)

// Send a TCP packet inside a virtual IP packet on our IP network
//
// NOTE: This is just an example function which hard-codes all of the values in
// the virtual IP header.  In the project, you should have a function like
// send_ip(...), but it will look VERY different from this!!!  For example, you
// shouldn't be passing in the UDP conn and addr as arguments.  Instead, you may
// want to specify the virtual source/dest address, or an interface name
func SendFakeTCPPacket(conn *net.UDPConn, linkLayerRemoteAddr *net.UDPAddr,
	sourceIp netip.Addr, destIp netip.Addr,
	payload []byte) (int, error) {

	// Start filling in the TCP header
	// WARNING:  This example uses hard-coded values for the port numbers, seq
	// and ack numbers, flags, and window size--you will want to do something
	// VERY different in your project!
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

	checksum := iptcp_utils.ComputeTCPChecksum(&tcpHdr, sourceIp, destIp, payload)
	tcpHdr.Checksum = checksum

	// Serialize the TCP header
	tcpHeaderBytes := make(header.TCP, iptcp_utils.TcpHeaderLen)
	tcpHeaderBytes.Encode(&tcpHdr)

	// Combine the TCP header + payload into one byte array, which
	// becomes the payload of the IP packet
	ipPacketPayload := make([]byte, 0, len(tcpHeaderBytes)+len(payload))
	ipPacketPayload = append(ipPacketPayload, tcpHeaderBytes...)
	ipPacketPayload = append(ipPacketPayload, []byte(payload)...)

	bytesWritten, err := SendFakeIPPacket(conn, linkLayerRemoteAddr,
		sourceIp, destIp, int(iptcp_utils.IpProtoTcp),
		ipPacketPayload)

	return bytesWritten, err
}

// Send an IP packet on our virtual network
//
// NOTE: This is just an example function which hard-codes all of the values in
// the virtual IP header.  In the project, you should have a function like
// send_ip(...), but it will look VERY different from this!!!  For example, you
// shouldn't be passing in the UDP conn and addr as arguments.  Instead, you may
// want to specify the virtual source/dest address, or an interface name
func SendFakeIPPacket(conn *net.UDPConn, linkLayerRemoteAddr *net.UDPAddr,
	sourceIp netip.Addr, destIp netip.Addr,
	protocol int, payload []byte) (int, error) {

	// FIll in the IP header
	// NOTE:  This example uses hard-coded values for the
	// source, destination, and protocol--you will need to
	// do something different!
	hdr := ipv4header.IPv4Header{
		Version:  4,
		Len:      20, // Header length is always 20 when no IP options
		TOS:      0,
		TotalLen: ipv4header.HeaderLen + len(payload),
		ID:       0,
		Flags:    0,
		FragOff:  0,
		TTL:      32,
		Protocol: protocol,
		Checksum: 0, // Should be 0 until checksum is computed
		Src:      sourceIp,
		Dst:      destIp,
		Options:  []byte{},
	}

	// Assemble the IP header into a byte array
	ipHeaderBytes, err := hdr.Marshal()
	if err != nil {
		return 0, err
	}

	// Compute the IP checksum
	// Cast back to an int, which is what the Header structure expects
	hdr.Checksum = int(iptcp_utils.ComputeIPChecksum(ipHeaderBytes))

	ipHeaderBytes, err = hdr.Marshal()
	if err != nil {
		return 0, err
	}

	// Assemble everything into a single byte array
	bytesToSend := make([]byte, 0, len(ipHeaderBytes)+len(payload))
	bytesToSend = append(bytesToSend, ipHeaderBytes...)
	bytesToSend = append(bytesToSend, payload...)

	// Send the message to the "link-layer" addr:port on UDP
	bytesWritten, err := conn.WriteToUDP(bytesToSend, linkLayerRemoteAddr)
	if err != nil {
		return 0, err
	}

	return bytesWritten, nil

}

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

	fakeSourceIp := netip.MustParseAddr("10.0.0.1")
	fakeDestIp := netip.MustParseAddr("10.1.0.2")

	bytesWritten, err := SendFakeTCPPacket(conn, remoteAddr, fakeSourceIp, fakeDestIp, []byte(message))
	if err != nil {
		log.Fatalln("Error sending packet:  ", err)
	}
	fmt.Printf("Sent %d bytes\n", bytesWritten)
}
