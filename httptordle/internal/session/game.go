package session

import (
	"errors"
	"goprojects/httptordle/internal/tordle"
)

// A GameID represents the ID of a game.
type GameID string

// Status is the current status of the game and tells what operations can be made on it.
type Status string

// Guess is a pair of a word (submitted by the player) and its feedback (provided by Tordle).
type Guess struct {
	Word     string
	Feedback string
}

// Game contains the information about a game
type Game struct {
	// ID is the identified of a game.
	ID GameID

	// The game of Tordle that is being played.
	Tordle tordle.Game

	// AttemptsLeft counts the number of attempts left before the game is over.
	AttemptsLeft byte

	// Guesses is the list of past guesses, and their feedback.
	Guesses []Guess

	// Status tells whether the game is playable.
	Status Status
}

const (
	StatusPlaying = "Playing"
	StatusWon     = "Won"
	StatusLost    = "Lost"
)

// ErrGameOver is returned when a play is made but the game is over.
var ErrGameOver = errors.New("game over")
