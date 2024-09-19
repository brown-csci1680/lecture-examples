// Channels example
//
// This example demonstrates multiple versions of a main function
// to wait for a function to complete.  See the inline comments
// for details.
package main

import (
	"fmt"
	"time"
)

var doneChan = make(chan bool, 1) // Make a channel to notify when somethingAsync has finished

func somethingAsync() {
	fmt.Println("Doing something...")
	time.Sleep(3 * time.Second) // Wait 3 seconds
	fmt.Println("Done!")

	doneChan <- true // Send something to the channel (blocks if full)
}

// Initial version
func mainV0() {
	go somethingAsync()

}

// Q:  Will this work?? Will somethingAsync run to completion?
// No!  main will return right after starting somethingAsync.
// When main returns, the program exits, so the goroutine never finishes

// ***** VERSION 1 *****
// Solution:  use a channel (doneChan) to wait until somethingAsync has finished
func mainV1() {

	go somethingAsync()

	<-doneChan // Block until something to be written to the channel
}

// ***** VERSION 2 *****
// V1 works, but now our main() is blocked.
// What if we want to do something else at the same time?
// Example:  let's say we want to print a status update (the current time)
// every 250ms while somethingAsync is running
// There are many ways to write this code, here's one:
func mainV2() {
	updateChan := make(chan time.Time, 1) // Another channel to signal
	// when we should print a status update

	go somethingAsync()

	// This block runs as a goroutine (anonymous function)
	go func() {
		// Send the current time via updateChan every 250ms
		// (Note: This isn't a great way to implement periodic
		// events--see next example)
		for {
			updateChan <- time.Now()
			time.Sleep(250 * time.Millisecond)
		}
	}()

	for {
		// select:  wait on a series of channels.
		// When one channel has data available, run that case
		// When all channels are empty, the thread sleeps
		select {
		case <-doneChan:
			return
		case t := <-updateChan:
			fmt.Println("Waiting, time = ", t)
		}
	}

}

// ***** Final version *****
func main() {
	go somethingAsync()

	// Go has a builtin for scheduling perioic events called Ticker,
	// so let's use that instead (this is more precise than time.Sleep())
	ticker := time.NewTicker(250 * time.Millisecond)

	for {
		// Wait on a series of channels, run the code for the channel that
		// has something available
		select {
		case <-doneChan:
			return
		case t := <-ticker.C:
			fmt.Println("Waiting, time = ", t)
		}
	}

}
