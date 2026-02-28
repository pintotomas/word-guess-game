package game

import (
	"errors"
	"strings"
	"unicode"
)

var ErrNoGuessesRemaining = errors.New("no guesses remaining")

type Game struct {
	word             string
	guessesRemaining int
	guessed          map[int]rune
}

func New(word string, maxGuesses int) *Game {
	word = strings.ToLower(word)
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

// Guess looks for "ch" in the word and updates the guessed word accordingly. If no matches, remaining attempts are decremented
func (g *Game) Guess(ch rune) error {
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
