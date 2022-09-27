package game

import (
	"math/rand"
	"net"
	"sync"
)

type ClientInfo struct {
	Id              int
	Conn            net.Conn
	GameResetChan   chan bool
	ServerCloseChan chan bool
}

type GameInfo struct {
	GameLock      sync.Mutex
	TotalGuesses  int
	TargetNumber  int32
	nextClientIdx int // Counter to increment each time we add a new client

	ClientListLock  sync.Mutex
	Clients         []*ClientInfo
	ClientWaitGroup sync.WaitGroup
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
	clientIndex := g.nextClientIdx
	g.nextClientIdx++

	ci := &ClientInfo{
		Id:              clientIndex,
		Conn:            conn,
		GameResetChan:   make(chan bool, 1),
		ServerCloseChan: make(chan bool, 1),
	}
	g.ClientWaitGroup.Add(1)

	g.ClientListLock.Lock()
	g.Clients = append(g.Clients, ci)
	g.ClientListLock.Unlock()

	return ci
}

func (g *GameInfo) RemoveClient(target *ClientInfo) {
	g.ClientListLock.Lock()
	var idx int

	// Find the client by its pointer and remove it from the list
	// TODO:  There might be a better way to handle this process,
	// comments welcome!

	for i, ci := range g.Clients {
		if ci == target {
			idx = i
		}
	}

	g.Clients = append(g.Clients[:idx], g.Clients[idx+1:]...)

	g.ClientListLock.Unlock()

	g.ClientWaitGroup.Done()
}

func (g *GameInfo) TerminateClients() {
	g.ClientListLock.Lock()
	for _, ci := range g.Clients {
		ci.ServerCloseChan <- true
	}
	g.ClientListLock.Unlock()

	// Wait for all clients to be done
	g.ClientWaitGroup.Wait()
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
