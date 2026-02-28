package game

import "testing"

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
