// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file provides the execution driver for the persistent segment tree
//	package. It initializes a versioned tree, performs range queries, creates
//	a new historical version via point updates, and verifies data immutability
//	and path-copying persistence across versions.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/braiden-gole/persistent-segment-tree/persistentsegmenttree"
)

func main() {
	// Initialize the source integer slice for the initial tree version.
	initialValues := []int{1, 2, 3, 4, 5}

	// Construct the base persistent segment tree holding Version 0.
	segmentTree, err := persistentsegmenttree.NewSegmentTree(initialValues)

	if err != nil {
		fmt.Printf("Error creating persistent segment tree: %v\n", err)

		return
	}

	// Query the range sum over indices 1 through 3 on the initial version.
	initialSum, queryError := segmentTree.QueryRange(0, 1, 3)

	if queryError != nil {
		fmt.Printf("Query error: %v\n", queryError)

		return
	}

	// Print the result of the range query on Version 0.
	fmt.Printf("Version 0 Sum [1, 3]: %d\n", initialSum)

	// Create a new historical tree version (Version 1) by updating index 2 to 10.
	versionOneIndex, insertError := segmentTree.InsertVersion(0, 2, 10)

	if insertError != nil {
		fmt.Printf("Insert error: %v\n", insertError)

		return
	}

	// Query Version 0 again to confirm it remains completely unmodified (persistence guarantee).
	versionZeroSumAgain, versionZeroQueryError := segmentTree.QueryRange(0, 1, 3)

	if versionZeroQueryError != nil {
		fmt.Printf("Query error: %v\n", versionZeroQueryError)

		return
	}

	// Print the immutable result of Version 0 after the update.
	fmt.Printf("Version 0 Sum [1, 3] after update: %d\n", versionZeroSumAgain)

	// Query the updated range sum over Version 1.
	versionOneSum, versionOneQueryError := segmentTree.QueryRange(versionOneIndex, 1, 3)

	if versionOneQueryError != nil {
		fmt.Printf("Query error: %v\n", versionOneQueryError)

		return
	}

	// Print the updated result for Version 1.
	fmt.Printf("Version 1 (index 2 updated to 10) Sum [1, 3]: %d\n", versionOneSum)
}
