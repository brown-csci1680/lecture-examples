package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"golang-sockets/pkg/game"
	"golang-sockets/pkg/protocol"
)

func HandleResponses(conn net.Conn, outChan chan protocol.GuessMessage) {
	for {
		msg, _ := protocol.ReadGuessMessage(conn)
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
	} else if msg.MessageType == protocol.MessageTypeNewGame {
		fmt.Println("New game!")
	} else {
		fmt.Println("Invalid message type:  ", msg.MessageType)
	}
}

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("Usage:  %s <address> <port number>",
			os.Args[0])
	}

	// Variables:  if we use :=,
	// the compiler will automatically determine the type
	address := os.Args[1]
	portNumber := os.Args[2]

	addrToUse := fmt.Sprintf("%s:%s", address, portNumber)

	conn, err := net.Dial("tcp4", addrToUse)
	if err != nil {
		log.Fatalln("Error connecting:  ", err)
	}
	defer conn.Close()

	// Get a net.TCPConn from a net.Conn
	// (This is called a type assertion)
	//tcpConn := conn.(*net.TCPConn)

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
			// msg := &protocol.GuessMessage{
			// 	MessageType: protocol.MessageTypeGuess,
			// 	Number:      int32(newGuess),
			// }
			// _, err := conn.Write(msg.Marshal())
			SendGuess(newGuess, conn)

			if err != nil {
				log.Fatalln("Error writing:  ", err)
			}

		case response := <-msgChan:
			PrintResponses(response)
		}
	}

}

func SendGuess(num int, conn net.Conn) {
	buf1 := new(bytes.Buffer)
	err := binary.Write(buf1, binary.BigEndian, uint8(protocol.MessageTypeGuess))
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}
	_, err = conn.Write(buf1.Bytes())
	if err != nil {
		log.Fatalln("Write", err)
	}

	time.Sleep(1 * time.Second)
	buf2 := new(bytes.Buffer)
	err = binary.Write(buf2, binary.BigEndian, int32(num))
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}
	_, err = conn.Write(buf2.Bytes())
	if err != nil {
		log.Fatalln("Write", err)
	}
}
