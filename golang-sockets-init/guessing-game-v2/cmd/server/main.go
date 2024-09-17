package main

import (
	"io"
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
	// Instead of cluttering our program printing to stdout, we can use the log package to write
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
	log.Printf("Client connected from %s\n", ci.Conn.RemoteAddr())
	for {
		// Try to read a guess message from the socket
		msg, err := protocol.ReadGuessMessage(ci.Conn)

		if err == io.EOF { // Client hung up
			log.Printf("%s:  Client exiting", ci.Conn.RemoteAddr())
			return // Exits the thread (socket will be closed)
		} else if err != nil { // Other errors
			log.Printf("%s:  Client exiting", ci.Conn.RemoteAddr())
			return
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
			// Invalid!  Client shouldn't send a response
		default:
			// Other message types are also invalid!
		}
	}
}
