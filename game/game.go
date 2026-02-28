package game

type Game struct {
	word             string
	guessesRemaining int
	guessed          map[int]rune
}

func New(word string, maxGuesses int) *Game {
	guessed := make(map[int]rune, len(word))
	for i := range word {
		guessed[i] = '_'
	}

	return &Game{
		word:             word,
		guessesRemaining: maxGuesses,
		guessed:          guessed,
	}
}

// Status returns remaining guesses and the status of the guessed word as a string
func (g *Game) Status() (string, int) {
	current := make([]rune, len(g.guessed))
	for i := 0; i < len(g.guessed); i++ {
		current[i] = g.guessed[i]
	}
	return string(current), g.guessesRemaining
}
