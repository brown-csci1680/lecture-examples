package main

import (
	"fmt"
	"go-lecture-demo/pkg/game"
	"go-lecture-demo/pkg/protocol"
	"log"
	"net"
	"os"
)

// Global pointer to our game state
var GameState *game.GameInfo = nil

func main() {

	portNumber := os.Args[1]

	// Create a socket and listen on port 8888
	conn, err := net.Listen("tcp4", fmt.Sprintf(":%s", portNumber))
	if err != nil {
		log.Fatalln("Error binding port ", err)
	}

	// Initialize our game state
	clientIndex := 0
	GameState = game.InitializeGame()
	log.Printf("Target number is %d.  Shhhh...\n", GameState.TargetNumber)

	for {
		// Wait for new connetions
		// this will block until someone connects
		clientConn, err := conn.Accept()

		if err != nil {
			log.Fatalln("accept", err)
		}

		// Create a new data structure representing our client state
		clientInfo := &game.ClientInfo{Conn: clientConn,
			Id: clientIndex}
		clientIndex++

		// Start a goroutine for this client
		go handleClient(clientInfo)
	}
}

// Runs once for each client
func handleClient(clientInfo *game.ClientInfo) {
	conn := clientInfo.Conn
	defer conn.Close() // When this funcion returns, call conn.Close()

	log.Printf("Client %d connected:  %s\n",
		clientInfo.Id,
		conn.RemoteAddr().String())

	for {
		// Wait for a message from the client
		guess := protocol.ReadGuessMessage(conn)
		log.Printf("Client guessed %d\n", guess.Number)

		// Decide if the client's guess was correct
		responseValue := GameState.DoGuess(guess.Number)

		// If the user guessed correctly, reset the game
		// (ie, pick a new number)
		if responseValue == game.GuessCorrect {
			GameState.ResetGame()
		}

		// Make a packet to send
		response := protocol.GuessMessage{
			MessageType: protocol.MessageTypeResponse,
			Number:      responseValue, // -1, 0, 1
		}
		conn.Write(response.Marshal()) // Send it out to the client
	}

}
