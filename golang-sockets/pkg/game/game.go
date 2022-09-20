package game

import (
	"math/rand"
	"net"
	"sync"
)

type ClientInfo struct {
	Conn          net.Conn
	GameResetChan chan bool
}

type GameInfo struct {
	GameLock     sync.Mutex
	TotalGuesses int
	TargetNumber int32

	ClientListLock sync.Mutex
	Clients        []*ClientInfo
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

func (g *GameInfo) NewClient(conn net.Conn) *ClientInfo {
	ci := &ClientInfo{
		Conn:          conn,
		GameResetChan: make(chan bool, 1),
	}

	g.ClientListLock.Lock()
	g.Clients = append(g.Clients, ci)
	g.ClientListLock.Unlock()

	return ci
}

func (g *GameInfo) ResetGame() {
	g.GameLock.Lock()
	defer g.GameLock.Unlock()

	g.TargetNumber = rand.Int31()
	g.TotalGuesses = 0

	for _, c := range g.Clients {
		c.GameResetChan <- true
	}
}

func (g *GameInfo) DoGuess(n int32) int32 {
	g.GameLock.Lock()
	defer g.GameLock.Unlock()

	g.TotalGuesses++

	if n < g.TargetNumber {
		return GuessTooLow
	} else if n > g.TargetNumber {
		return GuessTooHigh
	} else {
		return GuessCorrect
	}
}
