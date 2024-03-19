package getstatus

import (
	"encoding/json"
	"errors"
	"fmt"
	"goprojects/httptordle/internal/api"
	"goprojects/httptordle/internal/problem"
	"goprojects/httptordle/internal/repository"
	"goprojects/httptordle/internal/session"
	"log/slog"
	"net/http"
)

type gameFinder interface {
	Find(id session.GameID) (session.Game, error)
}

// Handler returns the handler for the get status endpoint.
func Handler(db gameFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(api.GameID)
		if id == "" {
			problem.Of(http.StatusBadRequest).
				Append(problem.Detail("missing the id of the game")).
				WriteTo(w)
			return
		}

		game, err := db.Find(session.GameID(id))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				problem.Of(http.StatusNotFound).
					Append(problem.Detail("this game does not exist yet")).
					WriteTo(w)
				return
			}

			slog.Error(fmt.Sprintf("cannot fetch game %s: %s", id, err))
			problem.Of(http.StatusInternalServerError).
				Append(problem.Detail("failed to fetch game")).
				WriteTo(w)
			return
		}

		apiGame := api.ToGameResponse(game)

		w.Header().Set(api.HeaderContentType, api.HeaderApplicationJson)
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to write response: %s", err))
			problem.Of(http.StatusInternalServerError).
				Append(problem.Detail("failed to write response")).
				WriteTo(w)
			return
		}
	}
}
