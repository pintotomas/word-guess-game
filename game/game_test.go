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
