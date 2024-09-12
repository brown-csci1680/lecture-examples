package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	//conn, err := net.Dial("tcp4", "127.0.0.1:6666")
	// if err != nil {
	// 	// Abort program with stack trace (good enough for now)
	// 	panic(err)
	// }

	// Continuously read lines from stdin to send as guesses
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter a guess: ")
	for scanner.Scan() {
		line := scanner.Text()

		guess, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Invalid input")
			fmt.Printf("Enter a guess: ")
			continue
		}

		fmt.Println("User guessed ", guess)

		//b, err := conn.Write(...)

		// fmt.Printf("Sent %d bytes\n", b)

		// For next loop
		fmt.Printf("Enter a guess: ")
	}
}

// func PrintResponses(msg *protocol.GuessMessage) {
// 	if msg.MessageType == protocol.MessageTypeResponse {
// 		switch msg.Number {
// 		case game.GuessTooHigh:
// 			fmt.Println("Too high!")
// 		case game.GuessTooLow:
// 			fmt.Println("Too low!")
// 		case game.GuessCorrect:
// 			fmt.Println("YAY!")
// 		default:
// 			fmt.Println("Invalid response:  ", msg.Number)
// 		}
// 	} else {
// 		fmt.Println("Invalid message type:  ", msg.MessageType)
// 	}
// }
