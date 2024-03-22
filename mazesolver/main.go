package main

import (
	"fmt"
	"goprojects/mazesolver/internal/solver"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	sol, err := solver.New(inputFile)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if err := sol.Solve(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if err := sol.SaveSolution(outputFile); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	log.Printf("Solving maze %q and saving it as %q", inputFile, outputFile)
}

// usage displays the usage of the program and exits the program
func usage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: maze_solver input.png output.png")
	os.Exit(1)
}
