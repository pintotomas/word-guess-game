package main

type guessRequest struct {
	ID    string `json:"id"`
	Guess string `json:"guess"`
}

type gameResponse struct {
	ID               string `json:"id"`
	Current          string `json:"current"`
	GuessesRemaining int    `json:"guesses_remaining"`
}
