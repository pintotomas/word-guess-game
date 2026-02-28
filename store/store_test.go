package store

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/fleetdm/wordgame/game"
)

func newTestGame(t *testing.T) *game.Game {
	t.Helper()
	g, err := game.New("APPLE", 6)
	if err != nil {
		t.Fatalf("unexpected error creating game: %v", err)
	}
	return g
}

func TestSaveAndGet(t *testing.T) {
	s := New()
	g := newTestGame(t)

	s.Save("abc", g)

	got, err := s.Get("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != g {
		t.Error("expected same game pointer back")
	}
}

func TestGetNotFound(t *testing.T) {
	s := New()

	_, err := s.Get("nonexistent")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	s := New()
	g := newTestGame(t)

	s.Save("abc", g)
	s.Delete("abc")

	_, err := s.Get("abc")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound after delete, got %v", err)
	}
}

func TestDeleteNonexistent(t *testing.T) {
	s := New()
	// should not raise any error
	s.Delete("nonexistent")
}

// TestConcurrentAccess runs Save, Get, and Delete concurrently to detect data races.
// Run with: go test -race ./store/
func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(3)
		id := fmt.Sprintf("game-%d", i)
		g := newTestGame(t)

		go func() {
			defer wg.Done()
			s.Save(id, g)
		}()

		go func() {
			defer wg.Done()
			s.Get(id) // may or may not find it, that's fine
		}()

		go func() {
			defer wg.Done()
			s.Delete(id)
		}()
	}

	wg.Wait()
}
