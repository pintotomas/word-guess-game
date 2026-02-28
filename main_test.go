package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/fleetdm/wordgame/store"
)

func setupTestServer() *httptest.Server {
	words := []string{"APPLE", "BANANA", "CHERRY"}
	s := store.New()
	r := newRouter(words, s)
	return httptest.NewServer(r)
}

func createGame(t *testing.T, serverURL string) gameResponse {
	t.Helper()
	resp, err := http.Post(serverURL+"/new", "application/json", nil)
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var gr gameResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return gr
}

func makeGuess(serverURL string, id string, guess string) (*http.Response, error) {
	body, _ := json.Marshal(guessRequest{ID: id, Guess: guess})
	return http.Post(serverURL+"/guess", "application/json", bytes.NewReader(body))
}

func TestNewGameEndpoint(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	gr := createGame(t, ts.URL)

	if gr.ID == "" {
		t.Error("expected non-empty game ID")
	}
	if gr.GuessesRemaining != 6 {
		t.Errorf("expected 6 guesses remaining, got %d", gr.GuessesRemaining)
	}
	for _, ch := range gr.Current {
		if ch != '_' {
			t.Errorf("expected all underscores, got %q", gr.Current)
			break
		}
	}
}

func TestGuessEndpoint(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	gr := createGame(t, ts.URL)

	resp, err := makeGuess(ts.URL, gr.ID, "A")
	if err != nil {
		t.Fatalf("failed to make guess: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result gameResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.ID != gr.ID {
		t.Errorf("expected ID %q, got %q", gr.ID, result.ID)
	}
}

func TestGuessInvalidBody(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	gr := createGame(t, ts.URL)

	resp, err := makeGuess(ts.URL, gr.ID, "AB")
	if err != nil {
		t.Fatalf("failed to make guess: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGuessNotFoundGame(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := makeGuess(ts.URL, "nonexistent-id", "A")
	if err != nil {
		t.Fatalf("failed to make guess: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestConcurrentNewGames(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	var wg sync.WaitGroup
	ids := make(chan string, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(ts.URL+"/new", "application/json", nil)
			if err != nil {
				t.Errorf("failed to create game: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}

			var gr gameResponse
			json.NewDecoder(resp.Body).Decode(&gr)
			ids <- gr.ID
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[string]bool)
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate game ID: %s", id)
		}
		seen[id] = true
	}
}

// TestConcurrentGuessesOnSameGame fires 6 wrong guesses concurrently on the same game.
// Exactly 6 should get a 200 (one per remaining guess), and the rest should get 404
func TestConcurrentGuessesOnSameGame(t *testing.T) {
	// Use a fixed word list so we know which letters are wrong
	s := store.New()
	r := newRouter([]string{"APPLE"}, s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	gr := createGame(t, ts.URL)

	// All wrong guesses
	wrongLetters := []string{"B", "D", "F", "G", "H", "I", "J", "K", "M", "N", "O", "Q"}

	var wg sync.WaitGroup
	var okCount int64
	var nonOkCount int64

	for _, l := range wrongLetters {
		wg.Add(1)
		go func(letter string) {
			defer wg.Done()
			resp, err := makeGuess(ts.URL, gr.ID, letter)
			if err != nil {
				t.Errorf("request failed: %v", err)
				return
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&okCount, 1)
			} else {
				atomic.AddInt64(&nonOkCount, 1)
			}
		}(l)
	}

	wg.Wait()

	// Exactly 6 guesses should succeed, the rest should fail.
	if okCount != 6 {
		t.Errorf("expected exactly 6 successful guesses, got %d", okCount)
	}
	if nonOkCount != int64(len(wrongLetters))-6 {
		t.Errorf("expected %d failed guesses, got %d", len(wrongLetters)-6, nonOkCount)
	}
}
