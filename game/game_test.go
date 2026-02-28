package game

import (
	"errors"
	"strings"
	"testing"
)

func newTestGame(t *testing.T, word string, maxGuesses int) *Game {
	t.Helper()
	g, err := New(word, maxGuesses)
	if err != nil {
		t.Fatalf("unexpected error creating game: %v", err)
	}
	return g
}

func TestNewGame(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)

	if g.word != "apple" {
		t.Errorf("expected word %q, got %q", "apple", g.word)
	}
	if g.guessesRemaining != 6 {
		t.Errorf("expected guessesRemaining %d, got %d", 6, g.guessesRemaining)
	}
}

func TestNewGameInvalidMaxGuesses(t *testing.T) {
	_, err := New("APPLE", 0)
	if !errors.Is(err, ErrInvalidMaxGuesses) {
		t.Errorf("expected ErrInvalidMaxGuesses, got %v", err)
	}

	_, err = New("APPLE", -1)
	if !errors.Is(err, ErrInvalidMaxGuesses) {
		t.Errorf("expected ErrInvalidMaxGuesses for negative value, got %v", err)
	}
}

func TestNewGameWordTooLong(t *testing.T) {
	word := strings.Repeat("a", 101)
	_, err := New(word, 6)
	if !errors.Is(err, ErrWordTooLong) {
		t.Errorf("expected ErrWordTooLong, got %v", err)
	}

	// 100 characters should be fine
	word = strings.Repeat("a", 100)
	_, err = New(word, 6)
	if err != nil {
		t.Errorf("expected no error for 100-char word, got %v", err)
	}
}

func TestNewGameInvalidWord(t *testing.T) {
	_, err := New("APP123", 6)
	if !errors.Is(err, ErrInvalidWord) {
		t.Errorf("expected ErrInvalidWord for digits, got %v", err)
	}

	_, err = New("APP LE", 6)
	if !errors.Is(err, ErrInvalidWord) {
		t.Errorf("expected ErrInvalidWord for spaces, got %v", err)
	}

	_, err = New("APP-LE", 6)
	if !errors.Is(err, ErrInvalidWord) {
		t.Errorf("expected ErrInvalidWord for hyphens, got %v", err)
	}
}

func TestStatusNewGame(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)

	current, remaining := g.Status()
	if current != "_____" {
		t.Errorf("expected current %q, got %q", "_____", current)
	}
	if remaining != 6 {
		t.Errorf("expected remaining %d, got %d", 6, remaining)
	}
}

func TestGuessCorrect(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)

	g.Guess('P')
	current, remaining := g.Status()
	if current != "_pp__" {
		t.Errorf("expected current %q, got %q", "_pp__", current)
	}
	if remaining != 6 {
		t.Errorf("expected remaining %d, got %d", 6, remaining)
	}
}

func TestGuessWrong(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)

	g.Guess('Z')
	current, remaining := g.Status()
	if current != "_____" {
		t.Errorf("expected current %q, got %q", "_____", current)
	}
	if remaining != 5 {
		t.Errorf("expected remaining %d, got %d", 5, remaining)
	}
}

func TestGuessCaseInsensitive(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)

	g.Guess('p') // lowercase
	current1, _ := g.Status()

	g2 := newTestGame(t, "APPLE", 6)
	g2.Guess('P') // uppercase
	current2, _ := g2.Status()

	if current1 != current2 {
		t.Errorf("expected same result for 'p' and 'P', got %q and %q", current1, current2)
	}
}

func TestGuessNoRemaining(t *testing.T) {
	g := newTestGame(t, "APPLE", 1)

	_ = g.Guess('Z') // wrong, drops to 0
	err := g.Guess('A')
	if !errors.Is(err, ErrNoGuessesRemaining) {
		t.Errorf("expected ErrNoGuessesRemaining, got %v", err)
	}
}

func TestWonNewGame(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)
	if g.Won() {
		t.Error("expected Won() to be false for a new game")
	}
}

func TestWonPartiallyGuessed(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)
	g.Guess('A')
	g.Guess('P')
	if g.Won() {
		t.Error("expected Won() to be false when not all letters guessed")
	}
}

func TestWonFullyGuessed(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)
	g.Guess('A')
	g.Guess('P')
	g.Guess('L')
	g.Guess('E')
	if !g.Won() {
		t.Error("expected Won() to be true when all letters guessed")
	}
}

func TestGuessAfterWon(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)
	g.Guess('A')
	g.Guess('P')
	g.Guess('L')
	g.Guess('E')

	err := g.Guess('Z')
	if !errors.Is(err, ErrGameWon) {
		t.Errorf("expected ErrGameWon, got %v", err)
	}
}

func TestNewGameGuessedMap(t *testing.T) {
	g := newTestGame(t, "APPLE", 6)

	if len(g.guessed) != 5 {
		t.Fatalf("expected guessed map length 5, got %d", len(g.guessed))
	}
	for i := 0; i < 5; i++ {
		if g.guessed[i] != '_' {
			t.Errorf("expected '_' at position %d, got %c", i, g.guessed[i])
		}
	}
}
