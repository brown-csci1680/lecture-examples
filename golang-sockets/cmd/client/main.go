package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"golang-sockets/pkg/game"
	"golang-sockets/pkg/protocol"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage:  %s <address> <port number>",
			os.Args[0])
	}

	// Variables in golang:  if we use :=,
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

	// We would like to be able to read from the socket and take keyboard input
	// at the same time--this way, the server can send us messages even while
	// we're waiting for the user to enter a guess
	// One way to do this is to create separate goroutines to watch each input source,
	// and then use channels to signal the main loop to act on the data
	keyboardChan := make(chan int, 1)
	msgChan := make(chan protocol.GuessMessage, 1)
	doneChan := make(chan struct{}, 1)

	// Blocking operation:  read from keyboard
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			guess, err := strconv.Atoi(line)

			if err != nil {
				fmt.Printf("Invalid guess:  %s\n", line)
				continue
			}

			// Wait for a line of input, send to main loop
			keyboardChan <- guess
		}
	}()

	// Start a goroutine to wait for a message from the server
	go HandleResponses(conn, msgChan, doneChan)

	for {

		// Watch both channels, do something when an event happens
		select {
		case newGuess := <-keyboardChan: // Input from keyboard
			SendGuess(newGuess, conn)
		case response := <-msgChan: // Input from socket
			PrintResponses(response)
		case <-doneChan:
			fmt.Printf("Server closed connection")
			return
		}
	}

}

func SendGuess(num int, conn net.Conn) {
	guess := &protocol.GuessMessage{MessageType: protocol.MessageTypeGuess,
		Number: int32(num)}

	bytesWritten, err := conn.Write(guess.Marshal())
	log.Printf("Wrote %d bytes\n", bytesWritten)
	if err != nil {
		// NOTE:  This not ideal--this function should do something better here,
		// like returning the error so the client can quit gracefully
		log.Fatalln("Write error:  ", err)
	}
}

func SendGuessV2(num int, conn net.Conn) {
	buf1 := new(bytes.Buffer)
	err := binary.Write(buf1, binary.BigEndian, uint8(protocol.MessageTypeGuess))
	if err != nil {
		log.Fatalln("Marshal failed:  ", err)
	}
	_, err = conn.Write(buf1.Bytes())
	if err != nil {
		log.Fatalln("Write", err)
	}

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

func HandleResponses(conn net.Conn, outChan chan protocol.GuessMessage, doneChan chan struct{}) {
	for {
		msg, err := protocol.ReadGuessMessage(conn, false)
		if err == io.EOF {
			doneChan <- struct{}{}
		}
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
