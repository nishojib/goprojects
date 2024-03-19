package main

import (
	"fmt"
	"goprojects/httptordle/internal/handlers"
	"goprojects/httptordle/internal/repository"
	"net/http"
	"os"
)

func main() {
	db := repository.New()

	err := http.ListenAndServe(":8000", handlers.New(db))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
