package game

type Game struct {
	word             string
	guessesRemaining int
}

func New(word string, maxGuesses int) *Game {
	return &Game{
		word:             word,
		guessesRemaining: maxGuesses,
	}
}
