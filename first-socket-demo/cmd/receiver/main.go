package main

import (
	"fmt"
	"net"
)

const MaxMessageSize = 1500

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:6666")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err) // Good enough for now
	}

	buffer := make([]byte, MaxMessageSize)
	b, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		panic(err) // Good enough for now
	}

	toPrint := string(buffer)
	fmt.Printf("Read %d bytes:  %s\n", b, toPrint)

}
