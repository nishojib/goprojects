package api

import (
	"goprojects/httptordle/internal/session"
)

// ToGameResponse converts a session.Game into a GameResponse.
func ToGameResponse(g session.Game) GameResponse {
	solution := g.Tordle.ShowAnswer()

	apiGame := GameResponse{
		ID:           string(g.ID),
		AttemptsLeft: g.AttemptsLeft,
		Guesses:      make([]Guess, len(g.Guesses)),
		Status:       string(g.Status),
		WordLength:   byte(len(solution)),
	}

	for index := range len(g.Guesses) {
		apiGame.Guesses[index].Word = g.Guesses[index].Word
		apiGame.Guesses[index].Feedback = g.Guesses[index].Feedback
	}

	if g.AttemptsLeft == 0 {
		apiGame.Solution = solution
	}

	return apiGame
}
