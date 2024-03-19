package guess_test

import (
	"goprojects/httptordle/internal/api"
	"goprojects/httptordle/internal/handlers/guess"
	"goprojects/httptordle/internal/session"
	"goprojects/httptordle/internal/tordle"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	game, _ := tordle.New([]string{"pocket"})

	req, err := http.NewRequest(
		http.MethodPut,
		"/games/123456",
		strings.NewReader(`{"guess":"pocket"}`),
	)
	require.NoError(t, err)

	req.SetPathValue(api.GameID, "123456")

	recorder := httptest.NewRecorder()

	handleFunc := guess.Handler(gameGuesserStub{
		session.Game{
			ID:           "123456",
			Tordle:       *game,
			AttemptsLeft: 5,
			Guesses:      nil,
			Status:       session.StatusPlaying,
		},
	})
	handleFunc(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, api.HeaderApplicationJson, recorder.Header().Get(api.HeaderContentType))
	assert.JSONEq(
		t,
		`{"id":"123456","attempts_left":4,"guesses":[{"word":"pocket", "feedback":"💚💚💚💚💚💚"}],"word_length":6,"status":"Won"}`,
		recorder.Body.String(),
	)

}

type gameGuesserStub struct {
	game session.Game
}

func (g gameGuesserStub) Find(_ session.GameID) (session.Game, error) {
	return g.game, nil
}

func (g gameGuesserStub) Update(_ session.Game) error {
	return nil
}
