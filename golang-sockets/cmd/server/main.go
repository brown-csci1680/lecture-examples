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
	"os/signal"
	"syscall"
	"time"
)

var GameState *game.GameInfo

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s <port number>", os.Args[0])
	}
	//log.Default().SetOutput(io.Discard) //Equivalent of writing logs to /dev/null

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
	defer conn.Close()

	// Another way to do this:
	// conn, err := net.Listen("tcp", fmt.Sprintf(":%s", portNumber))

	// Initialize the guessing game
	rand.Seed(time.Now().Unix())
	GameState = game.InitializeGame()
	fmt.Println("Target number is:  ", GameState.TargetNumber)

	// Instead of adding a REPL to our server (like Snowcast)
	// Catch Ctrl+C and use this to have the server close all connections
	ctrlCChan := make(chan os.Signal, 1)
	signal.Notify(ctrlCChan, os.Interrupt, syscall.SIGINT)

	go waitForConnections(conn)

	<-ctrlCChan
	fmt.Println("Caught Ctrl+C, closing clients...")
	GameState.TerminateClients()
	fmt.Println("All clients closed!")
}

func waitForConnections(listenConn *net.TCPListener) {
	for {
		// Wait for new connections (returns a new conn object for each client)
		conn, err := listenConn.Accept()
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
	defer conn.Close()
	defer GameState.RemoveClient(ci)

	log.Printf("New client:  %s\n", conn.RemoteAddr().String())

	// Our client handler needs to do two things:
	// 1. Respond to guesses from the client
	// 2. Send out a message when the game resets

	socketChan := make(chan protocol.GuessMessage, 1)
	go func() {
		defer conn.Close() // Ensure the socket is closed when this goroutine exits

		for {
			msg, err := protocol.ReadGuessMessage(conn, false)
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
				fmt.Printf("Received guess:  %d", msg.Number)

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

		case <-ci.ServerCloseChan:
			log.Printf("Server closing, removing client")
			return
		}
	}
}
