/*
 * Webserver demo
 */
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage:  %s <port> <directory>\n", os.Args[0])
		os.Exit(1)
	}

	port := os.Args[1]
	dir := os.Args[2]

	// https://pkg.go.dev/net/http#FileServer
	fileHandler := http.FileServer(http.Dir(dir))
	addr := fmt.Sprintf(":%s", port)

	err := http.ListenAndServe(addr, fileHandler)
	if err != nil {
		log.Fatal(err)
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
