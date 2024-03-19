package tordle_test

import (
	"errors"
	"goprojects/httptordle/internal/tordle"
	"testing"
)

func TestReadCorpus(t *testing.T) {
	tt := map[string]struct {
		file   string
		length int
		err    error
	}{
		"English corpus": {
			file:   "../../corpus/english.txt",
			length: 34,
			err:    nil,
		},
		"empty corpus": {
			file:   "../../corpus/empty.txt",
			length: 0,
			err:    tordle.ErrEmptyCorpus,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			words, err := tordle.ReadCorpus(tc.file)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected err %v, got %v", tc.err, err)
			}

			if tc.length != len(words) {
				t.Errorf("expected %d, got %d", tc.length, len(words))
			}
		})
	}
}
