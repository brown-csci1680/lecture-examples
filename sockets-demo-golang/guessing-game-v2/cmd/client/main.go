package main

import (
	"bufio"
	"fmt"
	"go-lecture-demo/pkg/game"
	"go-lecture-demo/pkg/protocol"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {

	address := os.Args[1]
	portNumber := os.Args[2]

	addrString := fmt.Sprintf("%s:%s", address, portNumber)

	conn, err := net.Dial("tcp4", addrString)
	if err != nil {
		log.Fatalln("connect", err)
	}
	defer conn.Close()

	// We would like to be able to read from the socket and take keyboard input
	// at the same time--this way, the server can send us messages even while
	// we're waiting for the user to enter a guess
	// One way to do this is to create separate goroutines to watch each input source,
	// and then use channels to signal the main loop to act on the data
	keyboardChan := make(chan int, 1)
	responseChan := make(chan protocol.GuessMessage, 1)

	// Goroutine to read from keyboard for a new guess
	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		// Wait for a line from stdin, convert it to an int
		for scanner.Scan() {
			fmt.Println("Enter a guess:  ")
			line := scanner.Text()
			guess, err := strconv.Atoi(line)

			if err != nil {
				fmt.Printf("Invalid guess:  %s\n", line)
				continue
			}

			// Send an integer to the channel
			keyboardChan <- guess
		}
	}()

	go func() {
		for {
			// Wait for a message from the server
			response := protocol.ReadGuessMessage(conn)

			responseChan <- response
		}
	}()

	for {
		// Watch both channels, act on one when something happens
		select {
		case guess := <-keyboardChan:
			// Got a new guess, create a message and send it
			msg := &protocol.GuessMessage{
				MessageType: protocol.MessageTypeGuess,
				Number:      int32(guess),
			}

			b, err := conn.Write(msg.Marshal())
			if err != nil {
				log.Fatalln("write", err)
			}
			log.Printf("Sent %d bytes\n", b)
		case response := <-responseChan:
			PrintResponses(response)
		}
	}

}
func PrintResponses(msg protocol.GuessMessage) {
	if msg.MessageType == protocol.MessageTypeResponse {
		switch msg.Number {
		case game.GuessTooHigh:
			fmt.Println("Too high!")
		case game.GuessTooLow:
			fmt.Println("Too low!")
		case game.GuessCorrect:
			fmt.Println("YAY!")
		default:
			fmt.Println("Invalid response:  ", msg.Number)
		}
	} else {
		fmt.Println("Invalid message type:  ", msg.MessageType)
	}
}
