package tordle

import (
	"slices"
	"testing"
)

func TestPickWord(t *testing.T) {
	corpus := []string{"HELLO", "SALUT", "ПРИВЕТ", "ΧΑΙΡΕ"}
	word := pickWord(corpus)

	if !slices.Contains(corpus, word) {
		t.Errorf("expected a word in the corpus, got %q", word)
	}
}
