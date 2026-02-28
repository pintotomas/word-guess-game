package game

import (
	"errors"
	"strings"
	"sync"
	"unicode"
)

// maxWordLength is a safe upper bound for word length. The longest English dictionary word is 45 characters,
// so 100 is more than enough.
const maxWordLength = 100

var (
	ErrNoGuessesRemaining = errors.New("no guesses remaining")
	ErrGameWon            = errors.New("game already won")
	ErrInvalidMaxGuesses  = errors.New("maxGuesses must be greater than 0")
	ErrInvalidWord        = errors.New("word must contain only English letters (a-z)")
	ErrWordTooLong        = errors.New("word must be at most 100 characters")
)

type Game struct {
	mu               sync.RWMutex
	word             string
	guessesRemaining int
	guessed          map[int]rune
}

// New creates a Game struct. It will lowercase the word as the game is case-insensitive
// Only English letters (a-z) are supported
func New(word string, maxGuesses int) (*Game, error) {
	if maxGuesses <= 0 {
		return nil, ErrInvalidMaxGuesses
	}

	if len(word) > maxWordLength {
		return nil, ErrWordTooLong
	}

	word = strings.ToLower(word)
	for _, r := range word {
		if r < 'a' || r > 'z' {
			return nil, ErrInvalidWord
		}
	}
	guessed := make(map[int]rune, len(word))
	for i := range word {
		guessed[i] = '_'
	}

	return &Game{
		word:             word,
		guessesRemaining: maxGuesses,
		guessed:          guessed,
	}, nil
}

// Status returns remaining guesses and the status of the guessed word as a string
func (g *Game) Status() (string, int) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	current := make([]rune, len(g.guessed))
	for i := 0; i < len(g.guessed); i++ {
		current[i] = g.guessed[i]
	}
	return string(current), g.guessesRemaining
}

// Won checks if there are any remaining characters to be guessed
func (g *Game) Won() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.won()
}

func (g *Game) won() bool {
	for _, r := range g.guessed {
		if r == '_' {
			return false
		}
	}
	return true
}

// Guess looks for "ch" in the word and updates the guessed word accordingly. If no matches, remaining attempts are decremented
func (g *Game) Guess(ch rune) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.won() {
		return ErrGameWon
	}
	if g.guessesRemaining == 0 {
		return ErrNoGuessesRemaining
	}

	ch = unicode.ToLower(ch)
	found := false
	for i, r := range g.word {
		if r == ch {
			g.guessed[i] = ch
			found = true
		}
	}
	if !found {
		g.guessesRemaining--
	}
	return nil
}
