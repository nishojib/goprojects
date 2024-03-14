package main

import (
	"fmt"
	"goprojects/tordle/tordle"
	"os"
)

const maxAttempts = 6

func main() {
	corpus, err := tordle.ReadCorpus("corpus/english.txt")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to read corpus: %s", err)
		os.Exit(1)
	}

	t, err := tordle.New(corpus, tordle.WithMaxAttempts(maxAttempts))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to start game: %s", err)
		os.Exit(1)
	}

	t.Play()
}
