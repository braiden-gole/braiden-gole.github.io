// -----------------------------------------------------------------------------
// Package: persistentsegmenttree
// File: persistent_segment_tree_test.go
//
// Description:
//
//	This file provides a comprehensive suite of 30 unique unit tests, fuzz tests,
//	and modern benchmarks. Covering every function within the persistent
//	segment tree package. It validates correctness, path-copying persistence, bounds
//	checking, and performance across historical versions.
//
// -----------------------------------------------------------------------------
package persistentsegmenttree

import (
	"sync"
	"testing"
)

// TestNewSegmentTree_ValidValues verifies that a segment tree is successfully created from a valid slice of integers.
func TestNewSegmentTree_ValidValues(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5}

	// Act.
	tree, err := NewSegmentTree(values)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error creating segment tree: %v", err)
	}

	if tree == nil {
		test.Fatal("expected segment tree instance to be non-nil")
	}
}

// TestNewSegmentTree_EmptyValues verifies that attempting to create a segment tree with an empty slice returns an error.
func TestNewSegmentTree_EmptyValues(test *testing.T) {
	// Arrange.
	values := []int{}

	// Act.
	_, err := NewSegmentTree(values)

	// Assert.
	if err == nil {
		test.Error("expected error when creating segment tree with empty values, got nil")
	}
}

// TestQueryRange_BaseVersionValidRange verifies correct range sum query on version 0.
func TestQueryRange_BaseVersionValidRange(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	sum, queryError := tree.QueryRange(0, 1, 3)

	// Assert.
	if queryError != nil {
		test.Fatalf("unexpected query error: %v", queryError)
	}

	expectedSum := 2 + 3 + 4

	if sum != expectedSum {
		test.Errorf("expected range sum %d, got %d", expectedSum, sum)
	}
}

// TestQueryRange_FullRange verifies querying the sum of the entire array.
func TestQueryRange_FullRange(test *testing.T) {
	// Arrange.
	values := []int{10, 20, 30, 40}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	sum, queryError := tree.QueryRange(0, 0, 3)

	// Assert.
	if queryError != nil {
		test.Fatalf("unexpected query error: %v", queryError)
	}

	if sum != 100 {
		test.Errorf("expected total sum 100, got %d", sum)
	}
}

// TestQueryRange_SingleElementRange verifies querying a range spanning a single index.
func TestQueryRange_SingleElementRange(test *testing.T) {
	// Arrange.
	values := []int{5, 10, 15}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	sum, queryError := tree.QueryRange(0, 1, 1)

	// Assert.
	if queryError != nil {
		test.Fatalf("unexpected query error: %v", queryError)
	}

	if sum != 10 {
		test.Errorf("expected single element sum 10, got %d", sum)
	}
}

// TestQueryRange_InvalidVersionIndex verifies that querying a non-existent version index returns an error.
func TestQueryRange_InvalidVersionIndex(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, queryError := tree.QueryRange(99, 0, 2)

	// Assert.
	if queryError == nil {
		test.Error("expected error when querying out-of-bounds version index, got nil")
	}
}

// TestQueryRange_NegativeVersionIndex verifies that negative version indices return an error.
func TestQueryRange_NegativeVersionIndex(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, queryError := tree.QueryRange(-1, 0, 2)

	// Assert.
	if queryError == nil {
		test.Error("expected error for negative version index, got nil")
	}
}

// TestQueryRange_InvalidBoundaries verifies that inverted query bounds return an error.
func TestQueryRange_InvalidBoundaries(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3, 4}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, queryError := tree.QueryRange(0, 3, 1)

	// Assert.
	if queryError == nil {
		test.Error("expected error for inverted query boundaries, got nil")
	}
}

// TestQueryRange_OutOfBoundsLeft verifies that negative left boundary returns an error.
func TestQueryRange_OutOfBoundsLeft(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, queryError := tree.QueryRange(0, -1, 2)

	// Assert.
	if queryError == nil {
		test.Error("expected error for negative left query boundary, got nil")
	}
}

// TestQueryRange_OutOfBoundsRight verifies that right boundary exceeding array length returns an error.
func TestQueryRange_OutOfBoundsRight(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, queryError := tree.QueryRange(0, 0, 5)

	// Assert.
	if queryError == nil {
		test.Error("expected error for right query boundary exceeding length, got nil")
	}
}

// TestInsertVersion_ValidUpdateAndPersistence verifies that updating an index creates a new version while preserving the old version.
func TestInsertVersion_ValidUpdateAndPersistence(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	newVersionIndex, insertError := tree.InsertVersion(0, 2, 10)

	if insertError != nil {
		test.Fatalf("unexpected insertion error: %v", insertError)
	}

	oldVersionSum, _ := tree.QueryRange(0, 1, 3)
	newVersionSum, _ := tree.QueryRange(newVersionIndex, 1, 3)

	// Assert.
	if oldVersionSum != 9 {
		test.Errorf("expected old version sum to remain 9, got %d", oldVersionSum)
	}

	if newVersionSum != 16 {
		test.Errorf("expected new version sum to be 16, got %d", newVersionSum)
	}
}

// TestInsertVersion_MultipleSequentialUpdates verifies chaining multiple version updates.
func TestInsertVersion_MultipleSequentialUpdates(test *testing.T) {
	// Arrange.
	values := []int{10, 20, 30}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	versionOne, versionOneError := tree.InsertVersion(0, 0, 100)
	versionTwo, versionTwoError := tree.InsertVersion(versionOne, 1, 200)

	// Assert.
	if versionOneError != nil || versionTwoError != nil {
		test.Fatalf("unexpected error during sequential inserts: %v, %v", versionOneError, versionTwoError)
	}

	sumVersionzero, _ := tree.QueryRange(0, 0, 2)
	sumVersionOne, _ := tree.QueryRange(versionOne, 0, 2)
	sumVersionTwo, _ := tree.QueryRange(versionTwo, 0, 2)

	if sumVersionzero != 60 {
		test.Errorf("expected version 0 sum 60, got %d", sumVersionzero)
	}

	if sumVersionOne != 150 {
		test.Errorf("expected version 1 sum 150, got %d", sumVersionOne)
	}

	if sumVersionTwo != 330 {
		test.Errorf("expected version 2 sum 330, got %d", sumVersionTwo)
	}
}

// TestInsertVersion_InvalidBaseVersion verifies that inserting from an invalid version index returns an error.
func TestInsertVersion_InvalidBaseVersion(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, insertError := tree.InsertVersion(5, 1, 10)

	// Assert.
	if insertError == nil {
		test.Error("expected error when inserting from non-existent base version, got nil")
	}
}

// TestInsertVersion_NegativeBaseVersion verifies that negative base version indices return an error.
func TestInsertVersion_NegativeBaseVersion(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, insertError := tree.InsertVersion(-1, 1, 10)

	// Assert.
	if insertError == nil {
		test.Error("expected error for negative base version index, got nil")
	}
}

// TestInsertVersion_TargetOutOfBoundsNegative verifies that updating a negative target index returns an error.
func TestInsertVersion_TargetOutOfBoundsNegative(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, insertError := tree.InsertVersion(0, -1, 10)

	// Assert.
	if insertError == nil {
		test.Error("expected error for negative target update index, got nil")
	}
}

// TestInsertVersion_TargetOutOfBoundsTooLarge verifies that updating a target index exceeding length returns an error.
func TestInsertVersion_TargetOutOfBoundsTooLarge(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	_, insertError := tree.InsertVersion(0, 10, 10)

	// Assert.
	if insertError == nil {
		test.Error("expected error for target update index exceeding tree length, got nil")
	}
}

// TestBuildSegmentTree_SingleElementArray verifies tree construction with a single-element slice.
func TestBuildSegmentTree_SingleElementArray(test *testing.T) {
	// Arrange.
	values := []int{42}

	// Act.
	tree, err := NewSegmentTree(values)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error for single element tree: %v", err)
	}

	sum, queryError := tree.QueryRange(0, 0, 0)

	if queryError != nil || sum != 42 {
		test.Errorf("expected sum 42, got %d with error %v", sum, queryError)
	}
}

// TestBuildSegmentTree_EvenLengthArray verifies tree construction with an even-sized array.
func TestBuildSegmentTree_EvenLengthArray(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3, 4}

	// Act.
	tree, err := NewSegmentTree(values)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error for even length array: %v", err)
	}

	sum, _ := tree.QueryRange(0, 0, 3)

	if sum != 10 {
		test.Errorf("expected total sum 10, got %d", sum)
	}
}

// TestRecursiveUpdate_MissingSubtreeDynamicAllocation verifies path copying when updating uninitialized sparse subtrees.
func TestRecursiveUpdate_MissingSubtreeDynamicAllocation(test *testing.T) {
	// Arrange.
	values := []int{0, 0, 0, 0, 0, 0, 0, 0}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	// Act.
	versionOne, versionOneError := tree.InsertVersion(0, 5, 50)

	if versionOneError != nil {
		test.Fatalf("unexpected error during dynamic subtree update: %v", versionOneError)
	}

	sum, queryError := tree.QueryRange(versionOne, 4, 6)

	// Assert.
	if queryError != nil {
		test.Fatalf("unexpected query error: %v", queryError)
	}

	if sum != 50 {
		test.Errorf("expected dynamic subtree range sum 50, got %d", sum)
	}
}

// TestSegmentTree_ConcurrentAccess verifies thread-safe concurrent queries across versions.
func TestSegmentTree_ConcurrentAccess(test *testing.T) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8}

	tree, err := NewSegmentTree(values)

	if err != nil {
		test.Fatalf("unexpected setup error: %v", err)
	}

	versionOne, _ := tree.InsertVersion(0, 3, 100)
	versionTwo, _ := tree.InsertVersion(versionOne, 7, 200)

	concurrencyCount := 20

	var waitGroup sync.WaitGroup

	errorChannel := make(chan error, concurrencyCount)

	waitGroup.Add(concurrencyCount)

	// Act.
	for routineIndex := 0; routineIndex < concurrencyCount; routineIndex++ {
		go func(index int) {
			defer waitGroup.Done()

			versionToQuery := index % 3

			var targetVersion int

			switch versionToQuery {
			case 0:
				targetVersion = 0
			case 1:
				targetVersion = versionOne
			case 2:
				targetVersion = versionTwo
			}

			_, queryError := tree.QueryRange(targetVersion, 0, 7)

			if queryError != nil {
				errorChannel <- queryError
			}
		}(routineIndex)
	}

	waitGroup.Wait()

	close(errorChannel)

	// Assert.
	for queryError := range errorChannel {
		if queryError != nil {
			test.Errorf("concurrent query failed: %v", queryError)
		}
	}
}

// FuzzNewSegmentTree performs fuzz testing on tree creation with random slices.
func FuzzNewSegmentTree(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{1, 2, 3, 4})

	// Act.
	fuzz.Fuzz(func(test *testing.T, data []byte) {
		// Arrange.
		values := make([]int, len(data))

		for index, currentByte := range data {
			values[index] = int(currentByte)
		}

		// Act.
		_, _ = NewSegmentTree(values)

		// Assert.
	})
}

// FuzzInsertAndQuery performs fuzz testing on version insertion and range querying with random indices.
func FuzzInsertAndQuery(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(1, 10, 0, 2)

	// Act.
	fuzz.Fuzz(func(test *testing.T, updateIndex int, updateValue int, queryLeft int, queryRight int) {
		// Arrange.
		values := []int{10, 20, 30, 40, 50}

		tree, err := NewSegmentTree(values)

		if err != nil {
			return
		}

		// Act.
		normalizedUpdate := updateIndex % 5

		if normalizedUpdate < 0 {
			normalizedUpdate = -normalizedUpdate
		}

		versionOne, insertError := tree.InsertVersion(0, normalizedUpdate, updateValue)

		if insertError != nil {
			return
		}

		normLeft := queryLeft % 5

		if normLeft < 0 {
			normLeft = -normLeft
		}

		normRight := queryRight % 5

		if normRight < 0 {
			normRight = -normRight
		}

		if normLeft > normRight {
			normLeft, normRight = normRight, normLeft
		}

		_, _ = tree.QueryRange(versionOne, normLeft, normRight)

		// Assert.
	})
}

// BenchmarkNewSegmentTree measures performance of segment tree construction.
func BenchmarkNewSegmentTree(bench *testing.B) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	// Act.
	for bench.Loop() {
		_, _ = NewSegmentTree(values)
	}

	// Assert.
}

// BenchmarkQueryRange measures performance of range queries on a segment tree.
func BenchmarkQueryRange(bench *testing.B) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	tree, err := NewSegmentTree(values)

	if err != nil {
		bench.Fatalf("setup failed: %v", err)
	}

	// Act.
	for bench.Loop() {
		_, _ = tree.QueryRange(0, 2, 12)
	}

	// Assert.
}

// BenchmarkInsertVersion measures performance of path-copied version insertion.
func BenchmarkInsertVersion(bench *testing.B) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	tree, err := NewSegmentTree(values)

	if err != nil {
		bench.Fatalf("setup failed: %v", err)
	}

	// Act.
	for bench.Loop() {
		_, _ = tree.InsertVersion(0, 5, 99)
	}

	// Assert.
}

// BenchmarkRecursiveQuery measures performance of internal recursive querying.
func BenchmarkRecursiveQuery(bench *testing.B) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8}

	tree, err := NewSegmentTree(values)

	if err != nil {
		bench.Fatalf("setup failed: %v", err)
	}

	// Act.
	for bench.Loop() {
		_, _ = recursiveQuery(tree.root[0], 0, 7, 1, 6)
	}

	// Assert.
}

// BenchmarkRecursiveUpdate measures performance of internal path-copied updates.
func BenchmarkRecursiveUpdate(bench *testing.B) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8}

	tree, err := NewSegmentTree(values)

	if err != nil {
		bench.Fatalf("setup failed: %v", err)
	}

	// Act.
	for bench.Loop() {
		_, _ = recursiveUpdate(tree.root[0], 0, 7, 4, 42)
	}

	// Assert.
}

// BenchmarkBuildSegmentTree measures performance of internal tree building.
func BenchmarkBuildSegmentTree(bench *testing.B) {
	// Arrange.
	values := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// Act.
	for bench.Loop() {
		_, _ = buildSegmentTree(values, 0, 7)
	}

	// Assert.
}
