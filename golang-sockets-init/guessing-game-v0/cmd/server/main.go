package main

import (
	"log"
	"net"
	"os"
)

const MaxMessageSize = 1500

func main() {
	log.SetOutput(os.Stderr)

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:6666")
	if err != nil {
		panic(err)
	}

	// Create listen socket (bind)
	listenConn, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	clientConn, err := listenConn.Accept()
	if err != nil {
		panic(err) // Good enough???????????
	}

	// TODO:  Do something with clientConn
	_ = clientConn // (this silences the unused var warning)

}
