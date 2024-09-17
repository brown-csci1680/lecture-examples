package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sockets-demo/pkg/game"
	"sockets-demo/pkg/protocol"
	"strconv"
)

func main() {
	// Would like to write
	//send("Hello world!")

	// Dial now does two things
	// - Creates a TCP socket
	// - Establishes connection to server (connect() syscall)
	conn, err := net.Dial("tcp4", "127.0.0.1:6666")
	if err != nil {
		// Abort program with stack trace (good enough for now)
		panic(err)
	}

	// Continuously read lines from stdin to send as guesses
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a guess: ") // No newline so it makes a prompt
	for scanner.Scan() {
		line := scanner.Text()

		guess, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Invalid input")
			fmt.Print("Enter a guess")
			continue
		}

		// Send the guess
		protocol.SendGuessMessage(conn, guess)

		// Wait for a response
		resp, err := protocol.ReadGuessMessage(conn)
		if err != nil {
			panic(err)
		}

		PrintResponses(resp)

		fmt.Print("Enter a guess: ")
	}
}

func PrintResponses(msg *protocol.GuessMessage) {
	if msg.MessageType == protocol.MessageTypeResponse {
		switch msg.Number {
		case game.GuessTooHigh:
			fmt.Println(">>>>> Too high!")
		case game.GuessTooLow:
			fmt.Println(">>>>> Too low!")
		case game.GuessCorrect:
			fmt.Println(">>>>> YAY!")
		default:
			fmt.Println("Invalid response:  ", msg.Number)
		}
	} else {
		fmt.Println("Invalid message type:  ", msg.MessageType)
	}
}
