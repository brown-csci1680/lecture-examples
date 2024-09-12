package game

const (
	GuessTooHigh = 1
	GuessCorrect = 0
	GuessTooLow  = -1
)

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

// type GameInfo struct {
// 	targetNumber int32 // Lower case fields are private, upper-case is public
// 	totalGuesses int
// }

// func InitializeGame() *GameInfo {
// 	gi := &GameInfo{
// 		// Fields initialized to zero if not specified
// 	}
// 	gi.resetGame()

// 	return gi
// }

// func (g *GameInfo) DoGuess(n int32) int32 {
// 	g.totalGuesses++

// 	if n < g.targetNumber {
// 		return GuessTooLow
// 	} else if n > g.targetNumber {
// 		return GuessTooHigh
// 	} else {
// 		g.resetGame()
// 		return GuessCorrect
// 	}
// }

// func (g *GameInfo) resetGame() {
// 	// Should only be called when lock is held
// 	g.targetNumber = rand.Int31n(8192)
// 	g.totalGuesses = 0

// 	fmt.Println("Target number is:  ", g.targetNumber)
// }
