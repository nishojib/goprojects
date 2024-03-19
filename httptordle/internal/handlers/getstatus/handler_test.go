package getstatus_test

import (
	"goprojects/httptordle/internal/api"
	"goprojects/httptordle/internal/handlers/getstatus"
	"goprojects/httptordle/internal/session"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/games/", nil)
	require.NoError(t, err)

	req.SetPathValue(api.GameID, "123456")

	recorder := httptest.NewRecorder()

	handle := getstatus.Handler(gameFinderStub{session.Game{ID: "123456"}, nil})
	handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, api.HeaderApplicationJson, recorder.Header().Get(api.HeaderContentType))
	assert.JSONEq(
		t,
		`{"id":"123456","attempts_left":0,"guesses":[],"word_length":0,"status":""}`,
		recorder.Body.String(),
	)
}

type gameFinderStub struct {
	game session.Game
	err  error
}

func (g gameFinderStub) Find(_ session.GameID) (session.Game, error) {
	return g.game, g.err
}
