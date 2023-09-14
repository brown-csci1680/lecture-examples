package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sockets-demo/pkg/game"
	"sockets-demo/pkg/protocol"
)

const MaxMessageSize = 1500

type ClientInfo struct {
	Conn       net.Conn
	NumGuesses int
}

var gameInfo *game.GameInfo

func main() {
	log.SetOutput(os.Stderr)

	gameInfo = game.InitializeGame()

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:6666")
	if err != nil {
		panic(err)
	}

	// Create listen socket (bind)
	listenConn, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	for {
		// Wait for a client to connect
		// Every client gets a new conn
		clientConn, err := listenConn.Accept()
		if err != nil {
			// TODO:  Probably shouldn't quit here...
			panic(err)
		}

		clientInfo := &ClientInfo{
			Conn:       clientConn,
			NumGuesses: 0,
		}
		go handleClient(clientInfo)

	}

}

func handleClient(ci *ClientInfo) {
	fmt.Printf("Starting client handler for %s", ci.Conn.RemoteAddr())
	for {
		//buffer := make([]byte, MaxMessageSize)
		//b, err := conn.Read(buffer)
		msg, err := protocol.ReadGuessMessage(ci.Conn)
		//time.Sleep(5 * time.Second) // Be annoying

		if err != nil {
			// BAD!  Don't want to crash server something goes wrong with one client
			// TODO:  Should handle graceful disconnect (err == io.EOF) vs. an actual error
			// Can also do other stuff with logging (like turn on logs with flags)
			log.Println("Error with client", err)
			break
		}

		switch msg.MessageType {
		case protocol.MessageTypeGuess:
			log.Printf("%s:  Client guessed:  %d\n",
				ci.Conn.RemoteAddr(), msg.Number)
			response := gameInfo.DoGuess(msg.Number)

			respMessage := protocol.GuessMessage{
				MessageType: protocol.MessageTypeResponse,
				Number:      response,
			}

			responseBytes, err := respMessage.Marshal()
			if err != nil {
				panic(err)
			}
			ci.Conn.Write(responseBytes)
		case protocol.MessageTypeResponse:
			// Invalid!
		default:
			// Also invalid!
		}

		// toPrint := string(buffer)
		// fmt.Printf("%s:  read %d bytes:  %s\n",
		// 	conn.RemoteAddr(), b, toPrint)
	}
}
