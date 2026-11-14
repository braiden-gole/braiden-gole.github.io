// -----------------------------------------------------------------------------
// Package: closestpairdivideandconquer
// File: closest_pair_test.go
//
// Description:
//
//	This file contains unit tests for the Closest Pair of Points
//	divide-and-conquer implementation.
//
//	The tests validate correct behavior across typical inputs,
//	degenerate cases, and defensive boundary conditions, including
//	empty collections, singleton collections, invalid numeric values,
//	and known geometric configurations with deterministic outcomes.
//
//	Each test follows a clear Arrange / Act / Assert structure to
//	emphasize intent and correctness over incidental implementation
//	details.
//
// Notes:
//   - Tests rely exclusively on the public API.
//   - Floating-point comparisons account for numerical precision.
//   - Positive infinity is treated as a valid sentinel result.
//
// -----------------------------------------------------------------------------
package closestpairdivideandconquer

import (
	"math"
	"testing"
)

// TestFindClosestPairDistance verifies correctness of the public API
// across a range of valid, degenerate, and invalid input scenarios.
func TestFindClosestPairDistance(test *testing.T) {
	// Arrange.
	type testCase struct {
		name            string
		collection      []Point
		minimumDistance float64
		tolerant        bool
	}

	testCollection := []testCase{
		{
			name: "two_points_basic_distance",
			collection: []Point{
				{X: 0.0, Y: 0.0},
				{X: 3.0, Y: 4.0},
			},
			minimumDistance: 5.0,
		},
		{
			name: "known_closest_pair",
			collection: []Point{
				{X: 2.0, Y: 3.0},
				{X: 12.0, Y: 30.0},
				{X: 40.0, Y: 50.0},
				{X: 5.0, Y: 1.0},
				{X: 12.0, Y: 10.0},
				{X: 3.0, Y: 4.0},
			},
			minimumDistance: math.Sqrt(2),
			tolerant:        true,
		},
		{
			name:            "empty_collection",
			collection:      []Point{},
			minimumDistance: math.Inf(1),
		},
		{
			name: "single_point_collection",
			collection: []Point{
				{X: 1.0, Y: 1.0},
			},
			minimumDistance: math.Inf(1),
		},
		{
			name: "collection_with_nan",
			collection: []Point{
				{X: math.NaN(), Y: 0.0},
				{X: 1.0, Y: 1.0},
			},
			minimumDistance: math.Inf(1),
		},
		{
			name: "collection_with_infinity",
			collection: []Point{
				{X: 0.0, Y: math.Inf(1)},
				{X: 1.0, Y: 1.0},
			},
			minimumDistance: math.Inf(1),
		},
	}

	for index := 0; index < len(testCollection); index++ {
		currentCase := testCollection[index]

		test.Run(currentCase.name, func(test *testing.T) {
			// Act.
			computedMinimumDistance := FindClosestPairDistance(currentCase.collection)

			// Assert.
			if currentCase.tolerant {
				floatingPointTolerance := 1e-9

				if math.Abs(computedMinimumDistance-currentCase.minimumDistance) > floatingPointTolerance {
					test.Fatalf("expected distance approximately %f, but got %f",
						currentCase.minimumDistance,
						computedMinimumDistance)
				}
			} else {
				if computedMinimumDistance != currentCase.minimumDistance {
					test.Fatalf("expected distance %f, but got %f",
						currentCase.minimumDistance,
						computedMinimumDistance)
				}
			}
		})
	}
}

// TestFindClosestPairDistance_DuplicatePoints verifies that identical points
// correctly yield a distance of zero.
func TestFindClosestPairDistance_DuplicatePoints(test *testing.T) {
	// Arrange.
	pointsCollection := []Point{
		{X: 1.5, Y: 2.5},
		{X: 1.5, Y: 2.5},
		{X: 10.0, Y: 10.0},
	}
	expectedDistance := 0.0

	// Act.
	computedDistance := FindClosestPairDistance(pointsCollection)

	// Assert.
	if computedDistance != expectedDistance {
		test.Fatalf("expected distance %.1f for duplicate points, got %.6f",
			expectedDistance,
			computedDistance,
		)
	}
}

// TestFindClosestPairDistance_CollinearPoints verifies correctness when all points
// lie on a single straight horizontal or vertical line.
func TestFindClosestPairDistance_CollinearPoints(test *testing.T) {
	// Arrange.
	pointsCollection := []Point{
		{X: 0.0, Y: 0.0},
		{X: 3.0, Y: 0.0},
		{X: 1.0, Y: 0.0},
		{X: 10.0, Y: 0.0},
	}

	expectedDistance := 1.0

	// Act.
	computedDistance := FindClosestPairDistance(pointsCollection)

	// Assert.
	if computedDistance != expectedDistance {
		test.Fatalf("expected distance %.1f for collinear points, got %.6f",
			expectedDistance,
			computedDistance,
		)
	}
}

// TestFindClosestPairDistance_NegativeCoordinates verifies correctness when
// points span across all four quadrants with negative coordinates.
func TestFindClosestPairDistance_NegativeCoordinates(test *testing.T) {
	// Arrange.
	pointsCollection := []Point{
		{X: -5.0, Y: -5.0},
		{X: -2.0, Y: -1.0},
		{X: -2.1, Y: -1.0},
		{X: 5.0, Y: 5.0},
	}

	expectedDistance := 0.1
	floatingPointTolerance := 1e-9

	// Act.
	computedDistance := FindClosestPairDistance(pointsCollection)

	// Assert.
	if math.Abs(computedDistance-expectedDistance) > floatingPointTolerance {
		test.Fatalf("expected distance approximately %.2f, got %.6f",
			expectedDistance,
			computedDistance,
		)
	}
}

// FuzzFindClosestPairDistance performs native Go fuzz testing on the Closest Pair algorithm
// using randomized coordinate values to ensure algorithmic stability without panicking.
func FuzzFindClosestPairDistance(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 0, 10, 0, 20, 0})
	fuzz.Add([]byte{1, 1, 2, 2, 3, 3})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawCoordinateBytes []byte) {
		if len(rawCoordinateBytes) < 4 || len(rawCoordinateBytes)%4 != 0 {
			return
		}

		numberOfPoints := len(rawCoordinateBytes) / 4
		pointsCollection := make([]Point, numberOfPoints)
		byteIndex := 0

		for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
			xCoordinate := float64(int8(rawCoordinateBytes[byteIndex]))
			byteIndex++

			yCoordinate := float64(int8(rawCoordinateBytes[byteIndex]))
			byteIndex++

			pointsCollection[pointIndex] = Point{X: xCoordinate, Y: yCoordinate}
		}

		// Act & Assert: Execute algorithm wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			computedDistance := FindClosestPairDistance(pointsCollection)

			if !math.IsInf(computedDistance, 1) && computedDistance < 0.0 {
				fuzzContext.Fatalf("invalid negative distance computed: %f", computedDistance)
			}
		}()
	})
}

// BenchmarkFindClosestPairDistance_SmallDataset benchmarks the performance of the algorithm
// on a small point set containing 50 points.
func BenchmarkFindClosestPairDistance_SmallDataset(bench *testing.B) {
	// Arrange.
	numberOfPoints := 50
	pointsCollection := make([]Point, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		pointsCollection[pointIndex] = Point{
			X: float64(pointIndex * 3),
			Y: float64((pointIndex * 7) % 50),
		}
	}

	// Act.
	for bench.Loop() {
		_ = FindClosestPairDistance(pointsCollection)
	}
}

// BenchmarkFindClosestPairDistance_MediumDataset benchmarks the performance of the algorithm
// on a medium-sized point set containing 500 points.
func BenchmarkFindClosestPairDistance_MediumDataset(bench *testing.B) {
	// Arrange.
	numberOfPoints := 500
	pointsCollection := make([]Point, numberOfPoints)

	for pointIndex := 0; pointIndex < numberOfPoints; pointIndex++ {
		pointsCollection[pointIndex] = Point{
			X: float64(pointIndex * 2),
			Y: float64((pointIndex * 13) % 300),
		}
	}

	// Act.
	for bench.Loop() {
		_ = FindClosestPairDistance(pointsCollection)
	}
}
