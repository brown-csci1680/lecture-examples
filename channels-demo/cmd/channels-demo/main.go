/*
 * IP-in-UDP listener example
 * To run:
 *  ./udp-ip-recv <bind port>
 *  where <bind port> is the port on which to receive packets.
 */
package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

var bytesReady uint32 = 0

var sendChan chan bool

func SendThread() {
	for {
		<-sendChan // Why is this necessary?
		if bytesReady != 0 {
			fmt.Printf("\"Sending\" %d bytes\n", bytesReady)
			bytesReady = 0
		}

		// Why SHOULDN'T you do this??
		//time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	repl := ReplInitialize()
	defer repl.Close()

	sendChan = make(chan bool, 1)
	go SendThread()

	for {
		line, done := ReplGetLine(repl)
		if done {
			break
		}

		numBytes, err := strconv.ParseUint(line, 10, 32)
		if err != nil {
			fmt.Println(err)
			continue
		}

		bytesReady += uint32(numBytes)
		sendChan <- true // Could send any type
	}

}

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
