package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	words, err := loadWords("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	s := store.New()

	r := mux.NewRouter()

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":                id,
			"current":           current,
			"guesses_remaining": guessesRemaining,
		})
	}).Methods(http.MethodPost)

	r.HandleFunc("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		g, err := s.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		current, guessesRemaining := g.Status()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":                id,
			"current":           current,
			"guesses_remaining": guessesRemaining,
		})
	}).Methods(http.MethodGet)

	r.HandleFunc("/games/{id}/guess", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "guess for game: %s\n", id)
	}).Methods(http.MethodPost)

	log.Printf("Starting server on http://%s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatal(err)
	}
}
