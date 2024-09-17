package game

import (
	"log"
	"math/rand"
	"sync"
)

// Warning:  fields here are shared data since they're
// accessed by multiple clients!  Can protect with a mutex.
//
// Note:  struct fields starting with lower-case letters are private
// This might be useful here since we want to make sure other parts
// of the code can't read the fields without the mutex
type GameInfo struct {
	targetNumber int32
	totalGuesses int

	gameLock sync.Mutex
}

const (
	GuessTooHigh = 1
	GuessCorrect = 0
	GuessTooLow  = -1
)

func InitializeGame() *GameInfo {
	gi := &GameInfo{
		// Fields initialized to zero if not specified
	}
	gi.resetGame()

	return gi
}

func (g *GameInfo) DoGuess(n int32) int32 {
	g.gameLock.Lock()
	defer g.gameLock.Unlock() // Defer runs this when the function returns

	g.totalGuesses++

	if n < g.targetNumber {
		return GuessTooLow
	} else if n > g.targetNumber {
		return GuessTooHigh
	} else {
		g.resetGame()
		return GuessCorrect
	}
}

func (g *GameInfo) resetGame() {
	// Should only be called when lock is held
	g.targetNumber = rand.Int31n(8192)
	g.totalGuesses = 0

	log.Println("**** STARTING NEW GAME ****")
	log.Println("Target number is:  ", g.targetNumber)
	log.Println("***************************")

}
