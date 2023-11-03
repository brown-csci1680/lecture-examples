/*
 * Example for why/how to use channels to signal a thread to wake up
 *
 * FOR A DETAILED EXPLANATION:  see the gearup II video
 */
package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

// In this example, we have a variable "bytesReady" that is set
// from the REPL.  When bytesReady != 0, the sending thread wakes up,
// "sends" some data (here, just prints it out), and then sets bytesReady back to 0
// In your project, you will have many more threads than this
// (and your sending thread will send data), but this minimal example
// demonstrates the interaction between threads and how we can deal with it.
var bytesReady uint32 = 0
var sendChan chan bool

func main() {
	repl := ReplInitialize()
	defer repl.Close()

	sendChan = make(chan bool, 1) // Make a channel (here declared as global, yours would be in a socket struct)

	go SendThread() // Start the sending thread

	for {
		line, done := ReplGetLine(repl) // Read from REPL
		if done {
			break
		}

		numBytes, err := strconv.ParseUint(line, 10, 32)
		if err != nil {
			fmt.Println(err)
			continue
		}

		bytesReady += uint32(numBytes)
		sendChan <- true // Send something to channel; threads waiting will wake up
	}

}

// Send thread Version 1 - busy waiting
// Why is this bad?  Continually checking a variable like this will
// keep the thread running all the time just to check a variable
// This will use 100% CPU, wasting CPU cycles!
// func SendThread() {
// 	for {
// 		if bytesReady != 0 {
// 			fmt.Printf("Sending %d bytes\n", bytesReady)
// 			bytesReady = 0
// 		}
// 	}
// }

// Send thread Version 2, still not good:  We could reduce the CPU usage by sleeping, but
// The thread still needs to wake up, even if no data is available
// this still wastes CPU cycles!
// More importantly, if we only check for new data every T millliseconds
// it means we only send data once every T milliseconds--this will negatively
// impact performance!
// func SendThread() {
// 	for {
// 		if bytesReady != 0 {
// 			fmt.Printf("Sending %d bytes\n", bytesReady)
// 			bytesReady = 0
// 		}
// 	}
// }

// Best version:  use a channel
// Waiting on a channel will block until another thread writes to it
// This means that the thread is idle until someone sends data!
// This is good because:
//   - The thread is asleep until there is data available (so no wasted cycles)
//   - The thread wake up as soon as data is available  (so no extra delays)
func SendThread() {
	for {
		<-sendChan // "Receive on the channel"; if nothing there, this will block
		if bytesReady != 0 {
			fmt.Printf("Sending %d bytes\n", bytesReady)
			bytesReady = 0
		}

	}
}

// ****************** REPL FUNCTIONS **************************
// See these for an example of how to get a REPL with history
// like in the IP/TCP reference
// This REPL relies on a go module for "readline".  To add it to your project,
// run:  "github.com/chzyer/readline"

// Initialize the repl
func ReplInitialize() *readline.Instance {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       "/tmp/readline-channels-demo.tmp",
		InterruptPrompt:   "^C",
		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}

	return l
}

// Get a line from the repl
// To keep the example clean, we abstract this into a helper.
// For better error handling, you may just want to do this in the loop that reads a line
func ReplGetLine(repl *readline.Instance) (string, bool) {
	line, err := repl.Readline()
	if err == readline.ErrInterrupt {
		return "", true
	} else if err == io.EOF {
		return "", true
	}

	line = strings.TrimSpace(line)

	return line, false
}
