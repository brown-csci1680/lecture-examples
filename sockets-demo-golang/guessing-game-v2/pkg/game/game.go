package game

import (
	"log"
	"math/rand"
	"net"
	"sync"
)

// Individual state for each client
type ClientInfo struct {
	Id   int
	Conn net.Conn
}

type GameInfo struct {
	GameLock     sync.Mutex
	TotalGuesses int   // Increment every guess
	TargetNumber int32 // What we're trying to guess

	// List of connected clients
	// (We'll finish implementing this later)
	ClientList     []*ClientInfo
	ClientListLock sync.Mutex
}

const (
	GuessTooHigh = 1
	GuessCorrect = 0
	GuessTooLow  = -1
)

func InitializeGame() *GameInfo {
	return &GameInfo{
		// Other fields initialized to zero
		TargetNumber: rand.Int31n(8192),
	}
}

func (g *GameInfo) AddClient(ci *ClientInfo) {
	g.ClientListLock.Lock()
	defer g.ClientListLock.Unlock()

	g.ClientList = append(g.ClientList, ci)
}

func (g *GameInfo) ResetGame() {
	g.GameLock.Lock()
	defer g.GameLock.Unlock() // Unlock when function returns

	g.TargetNumber = rand.Int31n(8192)
	g.TotalGuesses = 0
	log.Printf("Target number is %d.  Shhhh...\n", g.TargetNumber)

	g.ClientListLock.Lock()
	// TODO:  Notify each client that the game has reset
	// for _, ci := range g.ClientList {
	// 	// Send a message to each client
	// }
	g.ClientListLock.Unlock()
}

func (g *GameInfo) DoGuess(n int32) int32 {
	g.GameLock.Lock()
	defer g.GameLock.Unlock() // Release lock when function returns

	g.TotalGuesses++

	if n < g.TargetNumber {
		return GuessTooLow
	} else if n > g.TargetNumber {
		return GuessTooHigh
	} else {
		return GuessCorrect
	}

}
