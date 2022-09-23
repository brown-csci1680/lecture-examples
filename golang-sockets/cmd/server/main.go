package main

import (
	"fmt"
	"golang-sockets/pkg/game"
	"golang-sockets/pkg/protocol"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var GameState *game.GameInfo

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s <port number>", os.Args[0])
	}
	log.Default().SetOutput(io.Discard) //Equivalent of writing logs to /dev/null

	portNumber := os.Args[1]

	// Get a TCPAddr and listen on the port number we specified on the command line
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%s", portNumber))
	if err != nil {
		log.Fatalln("Error translating address:  ", err)
	}

	conn, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalln(err)
	}

	// Another way to do this:
	// conn, err := net.Listen("tcp", fmt.Sprintf(":%s", portNumber))

	// Initialize the guessing game
	rand.Seed(time.Now().Unix())
	GameState = game.InitializeGame()
	log.Println("Target number:  ", GameState.TargetNumber)

	for {
		// Wait for new connections (returns a new conn object for each client)
		conn, err := conn.Accept()
		if err != nil {
			log.Fatalln("accept:  ", err)
		}

		// Create new per-client state, and start a goroutine for this client
		ci := GameState.NewClient(conn)
		go handleClient(ci)
	}
}

func handleClient(ci *game.ClientInfo) {
	conn := ci.Conn
	defer conn.Close() // Ensure the socket is closed when this goroutine exits

	log.Printf("New client:  %s\n", conn.RemoteAddr().String())

	// Our client handler needs to do two things:
	// 1. Respond to guesses from the client
	// 2. Send out a message when the game resets

	socketChan := make(chan protocol.GuessMessage, 1)
	go func() {
		for {
			msg, err := protocol.ReadGuessMessage(conn)
			if err != nil {
				if err == io.EOF {
					log.Printf("Client closed connection")
				}
				close(socketChan)
				return
			} else {
				socketChan <- msg
			}

		}
	}()

	for {
		select {
		case msg, ok := <-socketChan:
			if !ok {
				log.Printf("Client exited")
				return
			} else {
				log.Printf("Received guess:  %d", msg.Number)

				responseValue := GameState.DoGuess(msg.Number)

				if responseValue == game.GuessCorrect {
					GameState.ResetGame()
				}

				response := protocol.GuessMessage{
					MessageType: protocol.MessageTypeResponse,
					Number:      responseValue,
				}
				conn.Write(response.Marshal())
			}

		case <-ci.GameResetChan:
			response := protocol.GuessMessage{
				MessageType: protocol.MessageTypeNewGame,
				Number:      0,
			}
			conn.Write(response.Marshal())
		}

	}
}
