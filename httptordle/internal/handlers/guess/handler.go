package guess

import (
	"encoding/json"
	"errors"
	"fmt"
	"goprojects/httptordle/internal/api"
	"goprojects/httptordle/internal/problem"
	"goprojects/httptordle/internal/repository"
	"goprojects/httptordle/internal/session"
	"goprojects/httptordle/internal/tordle"
	"log/slog"
	"net/http"
)

type gameGuesser interface {
	Find(id session.GameID) (session.Game, error)
	Update(game session.Game) error
}

// Handler returns the handler for the guess endpoint.
func Handler(db gameGuesser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(api.GameID)
		if id == "" {
			problem.Of(http.StatusBadRequest).
				Append(problem.Detail("missing the id of the game")).
				WriteTo(w)
			return
		}

		req := api.GuessRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			problem.Of(http.StatusBadRequest).Append(problem.Wrap(err)).WriteTo(w)
			return
		}

		game, err := guess(session.GameID(id), req.Guess, db)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, repository.ErrNotFound):
				status = http.StatusNotFound
			case errors.Is(err, tordle.ErrInvalidGuessLength):
				status = http.StatusBadRequest
			case errors.Is(err, session.ErrGameOver):
				status = http.StatusForbidden
			default:
				status = http.StatusInternalServerError
			}

			problem.Of(status).Append(problem.Wrap(err)).WriteTo(w)
			return
		}

		apiGame := api.ToGameResponse(game)

		w.Header().Set(api.HeaderContentType, api.HeaderApplicationJson)
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to write response: %s", err))
		}
	}
}

func guess(id session.GameID, guess string, db gameGuesser) (session.Game, error) {
	// does the game exist?
	game, err := db.Find(id)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to find game: %w", err)
	}

	// are plays still allowed?
	if game.AttemptsLeft == 0 || game.Status == session.StatusWon {
		return session.Game{}, session.ErrGameOver
	}

	// what does tordle say about this guess?
	feedback, err := game.Tordle.Play(guess)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to play move: %w", err)
	}

	slog.Info("guess %v is valid in game %s", guess, id)

	// record the play
	game.Guesses = append(game.Guesses, session.Guess{
		Word:     guess,
		Feedback: feedback.String(),
	})

	game.AttemptsLeft -= 1

	switch {
	case feedback.GameWon():
		game.Status = session.StatusWon
	case game.AttemptsLeft == 0:
		game.Status = session.StatusLost
	default:
		game.Status = session.StatusPlaying
	}

	// update the game
	err = db.Update(game)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to save play: %w", err)
	}

	return game, nil
}
