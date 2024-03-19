package newgame

import (
	"encoding/json"
	"fmt"
	"goprojects/httptordle/internal/api"
	"goprojects/httptordle/internal/problem"
	"goprojects/httptordle/internal/session"
	"goprojects/httptordle/internal/tordle"
	"log/slog"
	"net/http"

	"github.com/oklog/ulid/v2"
)

type gameAdder interface {
	Add(game session.Game) error
}

// global variable that references each corpus
var corpora = map[string]string{
	"en": "./../../../corpus/english.txt",
	"he": "./../../../corpus/greek.txt",
	"cr": "./../../../corpus/cree.txt",
}

// Handler returns the handler handler for the game creation endpoint.
func Handler(db gameAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.URL.Query().Get(api.Lang)
		corpusPath, ok := corpora[lang]
		if !ok {
			corpusPath = corpora["en"]
		}

		game, err := createGame(db, corpusPath)
		if err != nil {
			slog.Error(fmt.Sprintf("unable to create a new game: %s", err))
			problem.Of(http.StatusInternalServerError).
				Append(problem.Detail("failed to create a new game")).
				WriteTo(w)
			return
		}

		w.Header().Set(api.HeaderContentType, api.HeaderApplicationJson)
		w.WriteHeader(http.StatusCreated)

		apiGame := api.ToGameResponse(game)
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to write response: %s", err))
		}
	}
}

const maxAttempts = 5

func createGame(db gameAdder, corpusPath string) (session.Game, error) {
	corpus, err := tordle.ReadCorpus(corpusPath)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to read corpus: %w", err)
	}

	if len(corpus) == 0 {
		return session.Game{}, tordle.ErrEmptyCorpus
	}

	game, err := tordle.New(corpus)
	if err != nil {
		return session.Game{}, fmt.Errorf("failed to create a new tordle game")
	}

	g := session.Game{
		ID:           session.GameID(ulid.Make().String()),
		Tordle:       *game,
		AttemptsLeft: maxAttempts,
		Guesses:      []session.Guess{},
		Status:       session.StatusPlaying,
	}

	if err = db.Add(g); err != nil {
		return session.Game{}, fmt.Errorf("failed to save the new game")
	}

	return g, nil
}
