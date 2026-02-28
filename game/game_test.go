package game

import (
	"errors"
	"testing"
)

func TestNewGame(t *testing.T) {
	g := New("APPLE", 6)

	if g.word != "apple" {
		t.Errorf("expected word %q, got %q", "apple", g.word)
	}
	if g.guessesRemaining != 6 {
		t.Errorf("expected guessesRemaining %d, got %d", 6, g.guessesRemaining)
	}
}

func TestStatusNewGame(t *testing.T) {
	g := New("APPLE", 6)

	current, remaining := g.Status()
	if current != "_____" {
		t.Errorf("expected current %q, got %q", "_____", current)
	}
	if remaining != 6 {
		t.Errorf("expected remaining %d, got %d", 6, remaining)
	}
}

func TestGuessCorrect(t *testing.T) {
	g := New("APPLE", 6)

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
	g := New("APPLE", 6)

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
	g := New("APPLE", 6)

	g.Guess('p') // lowercase
	current1, _ := g.Status()

	g2 := New("APPLE", 6)
	g2.Guess('P') // uppercase
	current2, _ := g2.Status()

	if current1 != current2 {
		t.Errorf("expected same result for 'p' and 'P', got %q and %q", current1, current2)
	}
}

func TestGuessNoRemaining(t *testing.T) {
	g := New("APPLE", 1)

	_ = g.Guess('Z') // wrong, drops to 0
	err := g.Guess('A')
	if !errors.Is(err, ErrNoGuessesRemaining) {
		t.Errorf("expected ErrNoGuessesRemaining, got %v", err)
	}
}

func TestWonNewGame(t *testing.T) {
	g := New("APPLE", 6)
	if g.Won() {
		t.Error("expected Won() to be false for a new game")
	}
}

func TestWonPartiallyGuessed(t *testing.T) {
	g := New("APPLE", 6)
	g.Guess('A')
	g.Guess('P')
	if g.Won() {
		t.Error("expected Won() to be false when not all letters guessed")
	}
}

func TestWonFullyGuessed(t *testing.T) {
	g := New("APPLE", 6)
	g.Guess('A')
	g.Guess('P')
	g.Guess('L')
	g.Guess('E')
	if !g.Won() {
		t.Error("expected Won() to be true when all letters guessed")
	}
}

func TestGuessAfterWon(t *testing.T) {
	g := New("APPLE", 6)
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
	g := New("APPLE", 6)

	if len(g.guessed) != 5 {
		t.Fatalf("expected guessed map length 5, got %d", len(g.guessed))
	}
	for i := 0; i < 5; i++ {
		if g.guessed[i] != '_' {
			t.Errorf("expected '_' at position %d, got %c", i, g.guessed[i])
		}
	}
}
