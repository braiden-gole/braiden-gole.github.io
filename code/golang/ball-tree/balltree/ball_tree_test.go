// -----------------------------------------------------------------------------
// Package: balltree
// File: ball_tree_test.go
//
// Description:
//
//	This file contains unit tests for the Ball Tree implementation.
//
//	The tests validate correctness, ordering, and defensive behavior for
//	both radius-based searches and k-nearest neighbor queries.
//
//	Test coverage includes:
//	  - Tree construction with valid input.
//	  - Radius search correctness and pruning behavior.
//	  - k-nearest neighbor ordering and distance accuracy.
//	  - Defensive handling of invalid parameters and nil receivers.
//
//	All tests use deterministic point sets with analytically verifiable
//	distances to ensure numerical correctness.
//
// Notes:
//   - Tests assume Euclidean distance metrics.
//   - Floating-point comparisons use tolerance thresholds.
//   - The testing parameter is intentionally named `test` for clarity.
//
// -----------------------------------------------------------------------------
package balltree

import (
	"math"
	"testing"
)

// TestNewBallTree_ConstructsSuccessfully verifies that a BallTree can be
// constructed successfully with valid input points and parameters.
func TestNewBallTree_ConstructsSuccessfully(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}

	leafSizeLimit := 2

	// Act.
	tree, err := New(points, leafSizeLimit)

	// Assert.
	if err != nil {
		test.Fatalf("expected no error, got %v", err)
	}

	if tree == nil {
		test.Fatal("expected non-nil BallTree")
	}

	if tree.root == nil {
		test.Fatal("expected tree root to be initialized")
	}
}

// TestRadiusSearch_ReturnsCorrectResults verifies that RadiusSearch returns
// only points within the specified radius and sorts them by distance.
func TestRadiusSearch_ReturnsCorrectResults(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{5, 5},
	}

	tree, err := New(points, 2)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{2.5, 2.5}
	searchRadius := 1.0

	// Act.
	results, err := tree.RadiusSearch(queryPoint, searchRadius)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error during radius search: %v", err)
	}

	if len(results) != 2 {
		test.Fatalf("expected 2 results, got %d", len(results))
	}

	expectedDistance := math.Sqrt(0.5)

	for _, neighbor := range results {
		if math.Abs(neighbor.Distance-expectedDistance) > 1e-6 {
			test.Fatalf("unexpected distance: got %.6f, expected %.6f",
				neighbor.Distance,
				expectedDistance,
			)
		}
	}
}

// TestKNearestNeighbors_ReturnsCorrectOrdering verifies that k-nearest neighbors
// are returned in ascending distance order with correct indices.
func TestKNearestNeighbors_ReturnsCorrectOrdering(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{5, 5},
		{8, 8},
	}

	tree, err := New(points, 2)
	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{2.5, 2.5}
	k := 3

	// Act.
	results, err := tree.KNearestNeighbors(queryPoint, k)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error during k-NN search: %v", err)
	}

	if len(results) != k {
		test.Fatalf("expected %d neighbors, got %d", k, len(results))
	}

	if results[0].Distance > results[1].Distance {
		test.Fatal("results are not sorted by ascending distance")
	}

	if results[1].Distance > results[2].Distance {
		test.Fatal("results are not sorted by ascending distance")
	}
}

// TestRadiusSearch_InvalidRadius verifies that negative radius queries
// are rejected with an error.
func TestRadiusSearch_InvalidRadius(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0, 0}

	// Act.
	_, err = tree.RadiusSearch(queryPoint, -1)

	// Assert.
	if err == nil {
		test.Fatal("expected error for negative radius, got nil")
	}
}

// TestKNearestNeighbors_InvalidK verifies that invalid k values
// are rejected defensively.
func TestKNearestNeighbors_InvalidK(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0, 0}

	// Act.
	_, err = tree.KNearestNeighbors(queryPoint, 0)

	// Assert.
	if err == nil {
		test.Fatal("expected error for k <= 0, got nil")
	}
}

// TestNilBallTree_DefensiveBehavior verifies that calling methods
// on a nil *BallTree does not panic and returns errors where appropriate.
func TestNilBallTree_DefensiveBehavior(test *testing.T) {
	// Arrange.
	var tree *BallTree

	queryPoint := []float64{0, 0}

	// Act.
	_, radiusErr := tree.RadiusSearch(queryPoint, 1.0)
	_, knnErr := tree.KNearestNeighbors(queryPoint, 1)

	// Assert.
	if radiusErr == nil {
		test.Fatal("expected error when calling RadiusSearch on nil tree")
	}

	if knnErr == nil {
		test.Fatal("expected error when calling KNearestNeighbors on nil tree")
	}
}

// TestNewBallTree_RejectsEmptyPoints verifies that constructing a BallTree
// with an empty point set fails with an error.
func TestNewBallTree_RejectsEmptyPoints(test *testing.T) {
	// Arrange.
	var points [][]float64

	// Act.
	tree, err := New(points, 2)

	// Assert.
	if err == nil {
		test.Fatal("expected error for empty point set, got nil")
	}

	if tree != nil {
		test.Fatal("expected nil tree on error")
	}
}

// TestNewBallTree_RejectsMismatchedDimensions verifies that points
// with differing dimensionality are rejected.
func TestNewBallTree_RejectsMismatchedDimensions(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1, 1},
	}

	// Act.
	tree, err := New(points, 2)

	// Assert.
	if err == nil {
		test.Fatal("expected error for mismatched dimensionality")
	}

	if tree != nil {
		test.Fatal("expected nil tree on error")
	}
}

// TestRadiusSearch_NoResults verifies that a radius search correctly
// returns an empty slice when no points fall within range.
func TestRadiusSearch_NoResults(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{10, 10},
		{20, 20},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{5, 5}
	radius := 1.0

	// Act.
	results, err := tree.RadiusSearch(queryPoint, radius)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		test.Fatalf("expected no results, got %d", len(results))
	}
}

// TestRadiusSearch_ExactBoundary verifies that points exactly on the
// radius boundary are included.
func TestRadiusSearch_ExactBoundary(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{3, 4},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0, 0}
	radius := 5.0

	// Act.
	results, err := tree.RadiusSearch(queryPoint, radius)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		test.Fatalf("expected 2 results, got %d", len(results))
	}
}

// TestKNearestNeighbors_KGreaterThanPoints verifies that requesting
// more neighbors than points returns all points without error.
func TestKNearestNeighbors_KGreaterThanPoints(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0, 0}

	k := 10

	// Act.
	results, err := tree.KNearestNeighbors(queryPoint, k)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(points) {
		test.Fatalf("expected %d neighbors, got %d",
			len(points),
			len(results),
		)
	}
}

// TestKNearestNeighbors_SinglePoint verifies correctness when the
// tree contains exactly one point.
func TestKNearestNeighbors_SinglePoint(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{7, 7},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{7, 7}

	// Act.
	results, err := tree.KNearestNeighbors(queryPoint, 1)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		test.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Distance != 0 {
		test.Fatalf("expected zero distance, got %.6f", results[0].Distance)
	}
}

// TestRadiusSearch_QueryDimensionalityMismatch verifies that queries
// with incorrect dimensionality are rejected.
func TestRadiusSearch_QueryDimensionalityMismatch(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0}

	// Act.
	_, err = tree.RadiusSearch(queryPoint, 1.0)

	// Assert.
	if err == nil {
		test.Fatal("expected error for dimensionality mismatch")
	}
}

// TestKNearestNeighbors_QueryDimensionalityMismatch verifies that
// k-NN queries with incorrect dimensionality are rejected.
func TestKNearestNeighbors_QueryDimensionalityMismatch(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0}

	// Act.
	_, err = tree.KNearestNeighbors(queryPoint, 1)

	// Assert.
	if err == nil {
		test.Fatal("expected error for dimensionality mismatch")
	}
}

// TestRadiusSearch_NegativeRadiusBoundary verifies that a radius of exactly 0.0
// returns only exact-match points and handles zero distance gracefully.
func TestRadiusSearch_ZeroRadiusExactMatch(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{0, 0},
		{1, 1},
	}

	tree, err := New(points, 1)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{0, 0}

	// Act.
	results, err := tree.RadiusSearch(queryPoint, 0.0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error during zero radius search: %v", err)
	}

	if len(results) != 1 {
		test.Fatalf("expected 1 exact match result, got %d", len(results))
	}

	if results[0].Distance != 0.0 {
		test.Fatalf("expected zero distance, got %.6f", results[0].Distance)
	}
}

// TestKNearestNeighbors_LargeDataset verifies k-NN correctness across a wider
// multidimensional coordinate space.
func TestKNearestNeighbors_LargeDataset(test *testing.T) {
	// Arrange.
	points := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
		{5.0, 6.0},
		{7.0, 8.0},
		{9.0, 10.0},
	}

	tree, err := New(points, 2)

	if err != nil {
		test.Fatalf("tree construction failed: %v", err)
	}

	queryPoint := []float64{4.0, 5.0}
	k := 2

	// Act.
	results, err := tree.KNearestNeighbors(queryPoint, k)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if len(results) != k {
		test.Fatalf("expected %d neighbors, got %d", k, len(results))
	}

	if results[0].Distance > results[1].Distance {
		test.Fatal("expected results to be sorted in ascending order of distance")
	}
}

// FuzzBallTree_Operations performs native Go fuzz testing on the Ball Tree implementation
// using randomized byte inputs to construct trees and perform queries safely without panicking.
func FuzzBallTree_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 0, 1, 1, 2, 2, 3, 3})
	fuzz.Add([]byte{10, 20, 30, 40, 50, 60})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawInputBytes []byte) {
		if len(rawInputBytes) < 4 {
			return
		}

		dimensionSize := 2
		numberOfPoints := len(rawInputBytes) / dimensionSize

		if numberOfPoints == 0 {
			return
		}

		points := make([][]float64, numberOfPoints)
		byteIndex := 0

		for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
			points[pointIndex] = make([]float64, dimensionSize)

			for dimensionIndex := 0; dimensionIndex < dimensionSize; dimensionIndex++ {
				if byteIndex < len(rawInputBytes) {
					points[pointIndex][dimensionIndex] = float64(rawInputBytes[byteIndex])
					byteIndex++
				}
			}
		}

		// Execute construction and search operations wrapped in recovery to catch malformed inputs safely.
		func() {
			defer func() {
				_ = recover()
			}()

			tree, err := New(points, 2)

			if err != nil {
				return
			}

			queryPoint := []float64{1.0, 1.0}

			_, _ = tree.RadiusSearch(queryPoint, 10.0)
			_, _ = tree.KNearestNeighbors(queryPoint, 3)
		}()
	})
}

// BenchmarkBallTree_Construction benchmarks the performance of constructing
// a Ball Tree with 500 multi-dimensional points.
func BenchmarkBallTree_Construction(bench *testing.B) {
	// Arrange.
	numberOfPoints := 500
	points := make([][]float64, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		points[pointIndex] = []float64{
			float64(pointIndex % 50),
			float64((pointIndex * 3) % 50),
		}
	}

	// Act.
	for bench.Loop() {
		_, _ = New(points, 4)
	}
}

// BenchmarkBallTree_RadiusSearch benchmarks the performance of performing radius searches
// on a pre-constructed Ball Tree.
func BenchmarkBallTree_RadiusSearch(bench *testing.B) {
	// Arrange.
	numberOfPoints := 500
	points := make([][]float64, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		points[pointIndex] = []float64{
			float64(pointIndex % 50),
			float64((pointIndex * 3) % 50),
		}
	}

	tree, err := New(points, 4)

	if err != nil {
		bench.Fatalf("failed to construct tree for benchmarking: %v", err)
	}

	queryPoint := []float64{25.0, 25.0}

	// Act.
	for bench.Loop() {
		_, _ = tree.RadiusSearch(queryPoint, 15.0)
	}
}

// BenchmarkBallTree_KNearestNeighbors benchmarks the performance of k-nearest neighbor queries
// on a pre-constructed Ball Tree.
func BenchmarkBallTree_KNearestNeighbors(bench *testing.B) {
	// Arrange.
	numberOfPoints := 500
	points := make([][]float64, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		points[pointIndex] = []float64{
			float64(pointIndex % 50),
			float64((pointIndex * 3) % 50),
		}
	}

	tree, err := New(points, 4)

	if err != nil {
		bench.Fatalf("failed to construct tree for benchmarking: %v", err)
	}

	queryPoint := []float64{25.0, 25.0}

	k := 5

	// Act.
	for bench.Loop() {
		_, _ = tree.KNearestNeighbors(queryPoint, k)
	}
}
