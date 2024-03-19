package repository

import (
	"fmt"
	"goprojects/httptordle/internal/session"
	"log/slog"
	"sync"
)

// GameRepository holds all the current games.
type GameRepository struct {
	mutex   sync.Mutex
	storage map[session.GameID]session.Game
}

// New creates an empty game repository.
func New() *GameRepository {
	return &GameRepository{
		storage: make(map[session.GameID]session.Game),
	}
}

// Add inserts for the first time a game in memory.
func (gr *GameRepository) Add(game session.Game) error {
	slog.Info("Adding a game...")

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	_, ok := gr.storage[game.ID]
	if ok {
		return fmt.Errorf("%w (%s)", ErrConflictingID, game.ID)
	}

	gr.storage[game.ID] = game
	return nil
}

// Find a game based on its ID. If nothing is found, return a nil pointer and an ErrNotFound error.
func (gr *GameRepository) Find(id session.GameID) (session.Game, error) {
	slog.Info("Finding a game...")

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	game, ok := gr.storage[id]
	if !ok {
		return session.Game{}, fmt.Errorf("can't find game %s: %w", id, ErrNotFound)
	}

	return game, nil
}

// Update a game in the database, overwriting it.
func (gr *GameRepository) Update(game session.Game) error {
	slog.Info("Updating a game...")

	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	_, ok := gr.storage[game.ID]
	if !ok {
		return fmt.Errorf("can't find game %s: %w", game.ID, ErrNotFound)
	}

	gr.storage[game.ID] = game
	return nil
}
