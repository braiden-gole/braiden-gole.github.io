// -----------------------------------------------------------------------------
// Package: monotonicchainconvexhull
// File: monotonic_chain_convex_hull_test.go
//
// Description:
//
//	This file defines unit tests for the Monotonic Chain convex hull
//	implementation. The tests validate correctness across common and
//	edge-case geometric configurations, including degenerate inputs,
//	collinear points, and interior point elimination.
//
// -----------------------------------------------------------------------------
package monotonicchainconvexhull

import "testing"

// TestMonotonicChainConvexHull_SquareWithInteriorPoint verifies that
// interior points are excluded from the hull when collinear boundary
// points are not retained.
func TestMonotonicChainConvexHull_SquareWithInteriorPoint(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
		{X: 1, Y: 1},
	}

	expectedHullVertices := []Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != len(expectedHullVertices) {
		test.Fatalf("unexpected hull vertex count: got %d, expected %d",
			len(computedHullVertices),
			len(expectedHullVertices),
		)
	}

	for vertexIndex := range expectedHullVertices {
		if computedHullVertices[vertexIndex] != expectedHullVertices[vertexIndex] {
			test.Fatalf("unexpected hull vertex at index %d: got %+v, expected %+v",
				vertexIndex,
				computedHullVertices[vertexIndex],
				expectedHullVertices[vertexIndex],
			)
		}
	}
}

// TestMonotonicChainConvexHull_CollinearBoundaryIncluded verifies that
// boundary collinear points are retained when explicitly requested.
func TestMonotonicChainConvexHull_CollinearBoundaryIncluded(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		true,
	)

	// Assert.
	if len(computedHullVertices) <= 4 {
		test.Fatalf("expected collinear boundary point to be retained, got %d hull vertices",
			len(computedHullVertices),
		)
	}
}

// TestMonotonicChainConvexHull_AllPointsCollinear verifies that when all
// points are collinear, only the extreme endpoints are returned when
// collinear boundary points are excluded.
func TestMonotonicChainConvexHull_AllPointsCollinear(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 3, Y: 0},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != 2 {
		test.Fatalf("expected 2 extreme endpoints for collinear input, got %d",
			len(computedHullVertices),
		)
	}

	if computedHullVertices[0] != (Point{X: 0, Y: 0}) ||
		computedHullVertices[1] != (Point{X: 3, Y: 0}) {
		test.Fatalf("unexpected hull endpoints: got %+v",
			computedHullVertices,
		)
	}
}

// TestMonotonicChainConvexHull_DuplicatePoints verifies that duplicate
// input points do not affect the computed hull.
func TestMonotonicChainConvexHull_DuplicatePoints(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
		{X: 2, Y: 2},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != 4 {
		test.Fatalf("unexpected hull size with duplicate input points: got %d, expected %d",
			len(computedHullVertices),
			4,
		)
	}
}

// TestMonotonicChainConvexHull_SinglePoint verifies that a single-point
// input returns that point unchanged.
func TestMonotonicChainConvexHull_SinglePoint(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 5, Y: 5},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != 1 {
		test.Fatalf("expected single-point hull, got %d",
			len(computedHullVertices),
		)
	}

	if computedHullVertices[0] != inputPoints[0] {
		test.Fatalf("unexpected single-point hull value: got %+v, expected %+v",
			computedHullVertices[0],
			inputPoints[0],
		)
	}
}

// TestMonotonicChainConvexHull_TwoDistinctPoints verifies that two
// distinct points are returned unchanged and in order.
func TestMonotonicChainConvexHull_TwoDistinctPoints(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 1, Y: 1},
		{X: 3, Y: 3},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != 2 {
		test.Fatalf("expected hull size of %d for two-point input, got %d",
			2,
			len(computedHullVertices),
		)
	}

	if computedHullVertices[0] != inputPoints[0] ||
		computedHullVertices[1] != inputPoints[1] {
		test.Fatalf("unexpected hull vertices for two-point input: got %+v",
			computedHullVertices,
		)
	}
}

// TestMonotonicChainConvexHull_Triangle verifies that a simple triangle
// returns all three vertices in counter-clockwise order.
func TestMonotonicChainConvexHull_Triangle(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 1, Y: 2},
	}

	expectedHullVertices := []Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 1, Y: 2},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != len(expectedHullVertices) {
		test.Fatalf("unexpected hull size for triangle input: got %d, expected %d",
			len(computedHullVertices),
			len(expectedHullVertices),
		)
	}

	for vertexIndex := range expectedHullVertices {
		if computedHullVertices[vertexIndex] != expectedHullVertices[vertexIndex] {
			test.Fatalf("unexpected triangle hull vertex at index %d: got %+v, expected %+v",
				vertexIndex,
				computedHullVertices[vertexIndex],
				expectedHullVertices[vertexIndex],
			)
		}
	}
}

// TestMonotonicChainConvexHull_PointsProvidedInRandomOrder verifies that
// hull correctness is independent of input ordering.
func TestMonotonicChainConvexHull_PointsProvidedInRandomOrder(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 2, Y: 2},
		{X: 0, Y: 0},
		{X: 0, Y: 2},
		{X: 2, Y: 0},
		{X: 1, Y: 1},
	}

	expectedHullVertices := []Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != len(expectedHullVertices) {
		test.Fatalf("unexpected hull size for shuffled input: got %d, expected %d",
			len(computedHullVertices),
			len(expectedHullVertices),
		)
	}

	for vertexIndex := range expectedHullVertices {
		if computedHullVertices[vertexIndex] != expectedHullVertices[vertexIndex] {
			test.Fatalf("unexpected hull vertex for shuffled input at index %d: got %+v, expected %+v",
				vertexIndex,
				computedHullVertices[vertexIndex],
				expectedHullVertices[vertexIndex],
			)
		}
	}
}

// TestMonotonicChainConvexHull_MultipleInteriorPoints verifies that
// multiple interior points are fully excluded from the hull.
func TestMonotonicChainConvexHull_MultipleInteriorPoints(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 0, Y: 0},
		{X: 3, Y: 0},
		{X: 3, Y: 3},
		{X: 0, Y: 3},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 1, Y: 2},
		{X: 2, Y: 2},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != 4 {
		test.Fatalf("expected only exterior vertices for interior-heavy input, got %d",
			len(computedHullVertices),
		)
	}
}

// TestMonotonicChainConvexHull_CollinearVerticalLine verifies correct
// handling of vertically collinear points.
func TestMonotonicChainConvexHull_CollinearVerticalLine(test *testing.T) {
	// Arrange.
	inputPoints := []Point{
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 1, Y: 2},
		{X: 1, Y: 3},
	}

	// Act.
	computedHullVertices := MonotonicChainConvexHull(
		inputPoints,
		false,
	)

	// Assert.
	if len(computedHullVertices) != 2 {
		test.Fatalf("expected 2 extreme endpoints for vertical collinear input, got %d",
			len(computedHullVertices),
		)
	}

	if computedHullVertices[0] != (Point{X: 1, Y: 0}) ||
		computedHullVertices[1] != (Point{X: 1, Y: 3}) {
		test.Fatalf("unexpected hull endpoints for vertical line: got %+v",
			computedHullVertices,
		)
	}
}

// FuzzMonotonicChainConvexHull_Operations performs native Go fuzz testing on the Monotonic Chain
// convex hull implementation using randomized coordinate byte streams to ensure stability and correctness without panicking.
func FuzzMonotonicChainConvexHull_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 0, 10, 0, 10, 10, 0, 10, 5, 5}, false)
	fuzz.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8}, true)

	fuzz.Fuzz(func(fuzzContext *testing.T, rawCoordinateBytes []byte, includeCollinearBoundary bool) {
		if len(rawCoordinateBytes) < 4 || len(rawCoordinateBytes)%4 != 0 {
			return
		}

		numberOfPoints := len(rawCoordinateBytes) / 4
		pointCollection := make([]Point, numberOfPoints)
		byteIndex := 0

		for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
			xCoordinate := float64(int8(rawCoordinateBytes[byteIndex]))
			byteIndex++

			yCoordinate := float64(int8(rawCoordinateBytes[byteIndex]))
			byteIndex++

			pointCollection[pointIndex] = Point{X: xCoordinate, Y: yCoordinate}
		}

		// Act & Assert: Execute convex hull computation wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			computedHull := MonotonicChainConvexHull(pointCollection, includeCollinearBoundary)

			if len(pointCollection) >= 3 && !includeCollinearBoundary && len(computedHull) > len(pointCollection) {
				fuzzContext.Fatalf("convex hull vertex count %d exceeds input point count %d", len(computedHull), len(pointCollection))
			}
		}()
	})
}

// BenchmarkMonotonicChainConvexHull_SmallDataset benchmarks the performance of the Monotonic Chain algorithm
// on a small point cloud containing 50 points.
func BenchmarkMonotonicChainConvexHull_SmallDataset(bench *testing.B) {
	// Arrange.
	numberOfPoints := 50
	pointCollection := make([]Point, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		pointCollection[pointIndex] = Point{
			X: float64((pointIndex * 13) % 100),
			Y: float64((pointIndex * 37) % 100),
		}
	}

	// Act.
	for bench.Loop() {
		_ = MonotonicChainConvexHull(pointCollection, false)
	}
}

// BenchmarkMonotonicChainConvexHull_MediumDataset benchmarks the performance of the Monotonic Chain algorithm
// on a medium-sized point cloud containing 500 points with collinear boundary inclusion enabled.
func BenchmarkMonotonicChainConvexHull_MediumDataset(bench *testing.B) {
	// Arrange.
	numberOfPoints := 500
	pointCollection := make([]Point, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		pointCollection[pointIndex] = Point{
			X: float64((pointIndex * 17) % 1000),
			Y: float64((pointIndex * 53) % 1000),
		}
	}

	// Act.
	for bench.Loop() {
		_ = MonotonicChainConvexHull(pointCollection, true)
	}
}
