package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	serverAddress = "localhost:1337"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	words, err := loadWords("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	// TODO use words in your implementation
	_ = words

	mux := http.NewServeMux()

	mux.HandleFunc("POST /games", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "new game")
	})

	mux.HandleFunc("GET /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "game status: %s\n", id)
	})

	mux.HandleFunc("POST /games/{id}/guess", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "guess for game: %s\n", id)
	})

	log.Printf("Starting server on http://%s", serverAddress)
	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		log.Fatal(err)
	}
}
