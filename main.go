package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "new game")
	}).Methods(http.MethodPost)

	r.HandleFunc("/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "game status: %s\n", id)
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
