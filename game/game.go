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
