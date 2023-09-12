package main

import (
	"fmt"
	"net"
)

const MaxMessageSize = 1500

func main() {
	listenConn, err := net.Listen("tcp4", "127.0.0.1:6666")
	if err != nil {
		panic(err) // Good enough for now
	}

	clientConn, err := listenConn.Accept()
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, MaxMessageSize)
	b, err := clientConn.Read(buffer)
	if err != nil {
		panic(err) // Good enough for now
	}

	toPrint := string(buffer)
	fmt.Printf("Read %d bytes:  %s", b, toPrint)

}
