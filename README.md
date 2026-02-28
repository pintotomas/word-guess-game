# Word game challenge

## The game

When a new game begins, the server chooses a random word from the list of words and sets the guesses remaining to 6. Let's say the server has chosen `APPLE`. The server will return `_____` (6 guesses remaining) to the player, indicating that the chosen word is 5 characters long. Now play progresses as follows:

1. The player guesses a character that they believe to be in the chosen word.

2. a) The player's guess matches a letter in the chosen word: The server returns a new string revealing the location(s) of that character.

   b) The player's guess does not match a letter in the chosen word: The server returns the existing string, and decrements the remaining guesses. 

In the above example, imagine the player guesses `P`, now the server returns `_PP__` (6 guesses remaining). Now the player guesses `I`, the server returns `_PP__` (5 guesses remaining).

The game ends either when there are no guesses remaining (player loses), or the player has guessed all of the characters in the word (player wins).

## Challenge

Please implement the API as specified below.

You may use the provided `main.go` file as the start of your implementation. This file loads the word list and starts an HTTP listener. The `generateIdentifier()` function in `identifier.go` may be used for generating game IDs.

The server can be started with `go run .` and is configured to listen on [http://localhost:1337](http://localhost:1337).


## API

Below, the API endpoints are described:

### New game

Start a new game.

#### Request

```
POST /new
```

Request body is ignored.

#### Response

The response contains the initial state for the newly started game.

* `id` - Identifier for this game. Use this to make guesses in the game.
* `current` - The current board state. Always consists of only `_` characters at start of game.
* `guesses_remaining` - How many guesses remain before the player loses.

#### Example

``` sh
curl -X POST https://wordgame.herokuapp.com/new
{"id":"f8302916-69f1-462b-b640-e503faa94397","current":"________","guesses_remaining":6}
```

### Make guess

Guess a letter in an ongoing game. Any game that is completed (when all letters are guessed, or no guesses remain) can be cleared from the data store.

#### Request

```
POST /guess
{"id":"<game_id>","guess":"<[A-Z]>"}
```

* `id` - Identifier for this game. Get a new game identifier from the `/new` endpoint.
* `guess` - The character to guess. Must be a single ASCII character (A-Z). 

#### Response

The response contains the updated game state.

* `id` - Identifier for this game. This will be unchanged from the request.
* `current` - The current board state.
* `guesses_remaining` - How many guesses remain before the player loses.

#### Example

``` sh
 curl -X POST https://wordgame.herokuapp.com/guess -d '{"id":"f8302916-69f1-462b-b640-e503faa94397","guess":"A"}'
{"id":"f8302916-69f1-462b-b640-e503faa94397","current":"______A_","guesses_remaining":6}
```

## Running locally

### With Go

```sh
go run .
```

### With Docker

Build and run:

```sh
docker build -t wordgame .
docker run -p 1337:1337 wordgame
```

### Testing the API

Start a new game:

```sh
curl -X POST http://localhost:1337/new
```

Make a guess:

```sh
curl -X POST http://localhost:1337/guess -d '{"id":"<game_id>","guess":"A"}'
```

### Running tests

```sh
go test -race ./...
```
