package game

import (
	"fmt"
	"math/rand"
)

type GameInfo struct {
	TargetNumber int32
	TotalGuesses int
}

const (
	GuessTooHigh = 1
	GuessCorrect = 0
	GuessTooLow  = -1
)

func InitializeGame() *GameInfo {
	gi := &GameInfo{
		// Other fields initialized to zero
		TargetNumber: rand.Int31n(8192),
	}
	fmt.Println("Target is", gi.TargetNumber)
	return gi
}

func (g *GameInfo) resetGame() {
	g.TargetNumber = rand.Int31n(8192)
	g.TotalGuesses = 0
}

func (g *GameInfo) DoGuess(n int32) int32 {
	g.TotalGuesses++

	if n < g.TargetNumber {
		return GuessTooLow
	} else if n > g.TargetNumber {
		return GuessTooHigh
	} else {
		g.resetGame()
		return GuessCorrect
	}
}
