package game

import (
	"math/rand"
	"sync"
)

type GameInfo struct {
	GameLock     sync.Mutex
	TotalGuesses int
	TargetNumber int32
}

const (
	GuessTooHigh = 1
	GuessCorrect = 0
	GuessTooLow  = -1
)

func InitializeGame() GameInfo {
	return GameInfo{
		// Other fields initialized to zero
		TargetNumber: rand.Int31(),
	}
}

func (GameInfo g) DoGuess(n int32) int32 {
	g.GameLock.Lock()
	defer g.GameLock.Unlock()

	if n < g.TargetNumber {
		return GuessTooLow
	} else if n > g.TargetNumber {
		return GuessTooHigh
	} else {
		return GuessCorrect
	}
}
