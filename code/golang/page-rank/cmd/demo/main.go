// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file provides an executable demonstration entry point for the
//	pagerank package. It constructs a sample directed web graph containing
//	interlinked pages, cycles, and a dangling node, configures standard
//	algorithm parameters (damping factor, iteration caps, and convergence
//	tolerance), invokes the GooglePageRank computation, and displays the
//	resulting steady-state importance scores.
//
// Notes:
//   - Serves as a runnable reference implementation of package usage.
//   - Outputs formatted rank results directly to standard output.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/braiden-gole/pagerank/pagerank"
)

func main() {
	// Define a sample web graph representing pages and their outgoing links.
	// - Page "A" links to "B" and "C".
	// - Page "B" links to "C".
	// - Page "C" links back to "A".
	// - Page "D" links to "C" (Note: "D" is a dangling node with no outgoing links).
	graph := pagerank.Graph{
		"A": {"B", "C"},
		"B": {"C"},
		"C": {"A"},
		"D": {"C"},
	}

	// Configure algorithm parameters:
	// - Damping factor: 0.85 (standard probability of a user following a link vs. random teleportation).
	// - Maximum iterations: 100 (upper bound for convergence).
	// - Tolerance: 1e-6 (threshold where smaller delta stops iteration early).
	dampingFactor := 0.85
	maximumIterations := 100
	tolerance := 1e-6

	// Execute the Google PageRank algorithm.
	ranks := pagerank.GooglePageRank(graph, dampingFactor, maximumIterations, tolerance)

	// Guard against a failed computation.
	if ranks == nil {
		fmt.Println("Error: PageRank calculation failed due to invalid inputs.")
		return
	}

	// Display the final computed PageRank scores for each node.
	fmt.Println("Final Google PageRank Scores:")
	fmt.Println("-----------------------------")

	for node, score := range ranks {
		fmt.Printf("Page %s: %.4f\n", node, score)
	}
}
