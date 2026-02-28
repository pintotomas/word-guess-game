package game

import (
	"errors"
	"testing"
)

func TestNewGame(t *testing.T) {
	g := New("APPLE", 6)

	if g.word != "APPLE" {
		t.Errorf("expected word %q, got %q", "APPLE", g.word)
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
	if current != "_PP__" {
		t.Errorf("expected current %q, got %q", "_PP__", current)
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

func TestGuessNoRemaining(t *testing.T) {
	g := New("APPLE", 1)

	_ = g.Guess('Z') // wrong, drops to 0
	err := g.Guess('A')
	if !errors.Is(err, ErrNoGuessesRemaining) {
		t.Errorf("expected ErrNoGuessesRemaining, got %v", err)
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
