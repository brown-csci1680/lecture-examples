package game

import "sync"

type GameInfo struct {
	GameLock     sync.Mutex
	TotalGuesses int
}
