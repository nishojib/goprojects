package tordle

import (
	"slices"
	"testing"
)

func TestPickRandomWord(t *testing.T) {
	words := []string{"HELLO", "SALUT", "ПРИВЕТ", "ΧΑΙΡΕ"}
	word, err := PickRandomWord(words)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)

	}

	if !slices.Contains(words, word) {
		t.Errorf("expected a word in the corpus, got %q", word)
	}
}
