package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fleetdm/wordgame/game"
	"github.com/fleetdm/wordgame/store"
	"github.com/gorilla/mux"
)

const (
	serverAddress = "localhost:1337"
	maxGuesses    = 6
)

func newRouter(words []string, s *store.Store) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		word := words[rand.Intn(len(words))]

		g, err := game.New(word, maxGuesses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := generateIdentifier()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Save(id, g)

		current, guessesRemaining := g.Status()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gameResponse{
			ID:               id,
			Current:          current,
			GuessesRemaining: guessesRemaining,
		})
	}).Methods(http.MethodPost)

	r.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		var req guessRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if len(req.Guess) != 1 {
			http.Error(w, "guess must be a single character", http.StatusBadRequest)
			return
		}

		g, err := s.Get(req.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ch := rune(req.Guess[0])
		if err := g.Guess(ch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		current, guessesRemaining := g.Status()

		// Clean up finished games
		if g.Won() || guessesRemaining == 0 {
			s.Delete(req.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gameResponse{
			ID:               req.ID,
			Current:          current,
			GuessesRemaining: guessesRemaining,
		})
	}).Methods(http.MethodPost)

	return r
}

func main() {
	rand.Seed(time.Now().UnixNano())

	words, err := loadWords("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	s := store.New()
	r := newRouter(words, s)

	log.Printf("Starting server on http://%s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatal(err)
	}
}
