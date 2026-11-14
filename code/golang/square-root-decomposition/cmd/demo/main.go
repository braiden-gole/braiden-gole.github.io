// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file provides a demonstration executable showcasing the usage of the
//	Square Root Decomposition package, including initialization, range queries,
//	point updates, and post-update range verification.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/braiden-gole/square-root-decomposition/squarerootdecomposition"
)

func main() {
	// Initialize the source numbers slice for the demonstration.
	numbers := []int{1, 4, 3, 2, 8, 5, 6, 7}

	// Instantiate the square root decomposition data structure.
	decomposition, err := squarerootdecomposition.NewSqrtDecomposition(numbers)

	// Handle any initialization errors defensively.
	if err != nil {
		fmt.Printf("Error creating decomposition: %v\n", err)

		return
	}

	// Query the sum of elements within the initial range.
	result, err := decomposition.QueryRange(1, 5)

	// Print the range sum result or handle query errors.
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
	} else {
		fmt.Printf("Sum of range [1, 5]: %d\n", result)
	}

	// Perform a point update on a specific element index.
	updateError := decomposition.UpdateItem(4, 10)

	// Report whether the item update succeeded or failed.
	if updateError != nil {
		fmt.Printf("Update error: %v\n", updateError)
	} else {
		fmt.Println("Updated element successfully.")
	}

	// Query the range sum again to verify dynamic cache updates.
	result, err = decomposition.QueryRange(1, 5)

	// Print the updated range sum result or handle errors.
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
	} else {
		fmt.Printf("New sum of range [1, 5] after update: %d\n", result)
	}
}
