package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"golang-sockets/pkg/game"
	"golang-sockets/pkg/protocol"
)

func HandleResponses(conn net.Conn, outChan chan protocol.GuessMessage) {
	for {
		msg := protocol.ReadGuessMessage(conn)
		outChan <- msg
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
	}
}

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("Usage:  %s <address> <port number>", os.Args[0])
	}

	address := os.Args[1]
	portNumber := os.Args[2]

	addrToUse := fmt.Sprintf("%s:%s", address, portNumber)

	addr, err := net.ResolveTCPAddr("tcp4", addrToUse)
	if err != nil {
		log.Fatalln("Error translating address:  ", err)
	}
	//conn, err := net.Dial("tcp4", addrToUse)
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Fatalln("Error connecting:  ", err)
	}
	defer conn.Close()

	fmt.Println("Connected!")

	guessChan := make(chan int, 1)
	msgChan := make(chan protocol.GuessMessage, 1)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			guess, err := strconv.Atoi(line)

			if err != nil {
				fmt.Printf("Invalid guess:  %s\n", line)
				continue
			}
			guessChan <- guess
		}
	}()

	go HandleResponses(conn, msgChan)

	for {
		select {
		case newGuess := <-guessChan:
			msg := &protocol.GuessMessage{
				MessageType: protocol.MessageTypeGuess,
				Number:      int32(newGuess),
			}
			_, err := conn.Write(msg.Marshal())

			if err != nil {
				log.Fatalln("Error writing:  ", err)
			}

		case response := <-msgChan:
			PrintResponses(response)
		}
	}

}
