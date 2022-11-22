package iptcp_utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/checksum"
	"gvisor.dev/gvisor/pkg/tcpip/header"
)

const (
	IpHeaderLen          = ipv4.HeaderLen
	TcpHeaderLen         = header.TCPMinimumSize
	TcpPseudoHeaderLen   = 12
	IpProtoTcp           = header.TCPProtocolNumber
	MaxVirtualPacketSize = 1400
)

// Build a TCPFields struct from the TCP byte array
//
// NOTE: the netstack package might have other options for parsing the header
// that you may like better--this example is most similar to our other class
// examples.  Your mileage may vary!
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

// The TCP checksum is computed based on a "pesudo-header" that
// combines the (virtual) IP source and destination address, protocol value,
// as well as the TCP header and payload
//
// This is one example one is way to combine all of this information
// and compute the checksum leveraging the netstack package.
//
// For more details, see the "Checksum" component of RFC9293 Section 3.1,
// https://www.rfc-editor.org/rfc/rfc9293.txt
func ComputeTCPChecksum(tcpHdr *header.TCPFields,
	sourceIP net.IP, destIP net.IP, payload []byte) uint16 {

	// Fill in the pseudo header
	pseudoHeaderBytes := make([]byte, TcpPseudoHeaderLen)

	// First are the source and dest IPs.  This function only supports
	// IPv4, so make sure the IPs are IPv4 addresses
	if ip := sourceIP.To4(); ip != nil {
		copy(pseudoHeaderBytes[0:4], ip)
	} else {
		// This error shouldn't ever occur in our project
		// If it did, would it be appropriate to call panic()?
		// No.  If we encounter a packet that has a processing
		// error, we should really just drop the packet, not
		// crash the node!
		panic("Invalid source IP length, only IPv4 supported")
	}

	if ip := destIP.To4(); ip != nil {
		copy(pseudoHeaderBytes[4:8], ip)
	} else {
		panic("Invalid dest IP length, only IPv4 supported")
	}

	// Next, add the protocol number and header length
	pseudoHeaderBytes[8] = uint8(0)
	pseudoHeaderBytes[9] = uint8(IpProtoTcp)

	totalLength := TcpHeaderLen + len(payload)
	binary.BigEndian.PutUint16(pseudoHeaderBytes[10:12], uint16(totalLength))

	// Turn the TcpFields struct into a byte array
	headerBytes := header.TCP(make([]byte, TcpHeaderLen))
	headerBytes.Encode(tcpHdr)

	// Compute the checksum for each individual part and combine To combine the
	// checksums, we leverage the "initial value" argument of the netstack's
	// checksum package to carry over the value from the previous part
	pseudoHeaderChecksum := checksum.Checksum(pseudoHeaderBytes, 0)
	headerChecksum := checksum.Checksum(headerBytes, pseudoHeaderChecksum)
	fullChecksum := checksum.Checksum(payload, headerChecksum)

	// Return the inverse of the computed value,
	// which seems to be the convention of the checksum algorithm
	// in the netstack package's implementation
	return fullChecksum ^ 0xffff
}

// Compute the checksum using the netstack package
func ComputeIPChecksum(b []byte) uint16 {
	checksum := checksum.Checksum(b, 0)

	// Invert the checksum value.  Why is this necessary?
	// This function returns the inverse of the checksum
	// on an initial computation.  While this may seem weird,
	// it makes it easier to use this same function
	// to validate the checksum on the receiving side.
	// See ValidateChecksum in the receiver file for details.
	checksumInv := checksum ^ 0xffff

	return checksumInv
}

// Validate the checksum using the netstack package Here, we provide both the
// byte array for the header AND the initial checksum value that was stored in
// the header
//
// "Why don't we need to set the checksum value to 0 first?"
//
// Normally, the checksum is computed with the checksum field of the header set
// to 0.  This library creatively avoids this step by instead subtracting the
// initial value from the computed checksum.  If you use a different language or
// checksum function, you may need to handle this differently.
func ValidateIPChecksum(b []byte, fromHeader uint16) uint16 {
	checksum := checksum.Checksum(b, fromHeader)

	return checksum
}

// Pretty-print TCP flags value as a string
func TCPFlagsAsString(flags header.TCPFlags) string {
	strMap := map[header.TCPFlags]string{
		header.TCPFlagAck: "ACK",
		header.TCPFlagFin: "FIN",
		header.TCPFlagPsh: "PSH",
		header.TCPFlagRst: "RST",
		header.TCPFlagSyn: "SYN",
		header.TCPFlagUrg: "URG",
	}

	matches := make([]string, 0)

	for b, str := range strMap {
		if flags.Contains(b) {
			matches = append(matches, str)
		}
	}

	ret := strings.Join(matches, "+")

	return ret
}

// Pretty-print a TCP header (with pretty-printed flags)
// Otherwise, using %+v in format strings is a good enough view in most cases
func TCPFieldsToString(hdr *header.TCPFields) string {
	return fmt.Sprintf("{SrcPort:%d DstPort:%d, SeqNum:%d AckNum:%d DataOffset:%d Flags:%s WindowSize:%d Checksum:%x UrgentPointer:%d}",
		hdr.SrcPort, hdr.DstPort, hdr.SeqNum, hdr.AckNum, hdr.DataOffset, TCPFlagsAsString(hdr.Flags), hdr.WindowSize, hdr.Checksum, hdr.UrgentPointer)
}
