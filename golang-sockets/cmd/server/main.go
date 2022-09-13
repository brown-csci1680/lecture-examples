package main

import (
	"fmt"
	"golang-sockets/pkg/game"
	"golang-sockets/pkg/protocol"
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

	portNumber := os.Args[1]

	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%s", portNumber))
	if err != nil {
		log.Fatalln("Error translating address:  ", err)
	}

	//conn, err := net.Listen("tcp", fmt.Sprintf(":%s", portNumber))
	conn, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalln(err)
	}

	rand.Seed(time.Now().Unix())
	GameState = game.InitializeGame()
	log.Println("Target number:  ", GameState.TargetNumber)

	for {
		conn, err := conn.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	log.Printf("New client:  %s\n", conn.RemoteAddr().String())
	for {
		msg := protocol.ReadGuessMessage(conn)
		log.Printf("Received guess:  %d", msg.Number)

		responseValue := GameState.DoGuess(msg.Number)

		response := protocol.GuessMessage{
			MessageType: protocol.MessageTypeResponse,
			Number:      responseValue,
		}

		conn.Write(response.Marshal())
	}
}
