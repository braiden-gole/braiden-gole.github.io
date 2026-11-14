// -----------------------------------------------------------------------------
// Package: floydwarshall
// File: floyd_warshall_test.go
//
// Description:
//
//	This file provides unit tests for the Floyd–Warshall all-pairs shortest
//	path algorithm implementation. The tests validate correctness across
//	common scenarios, including basic graphs, unreachable vertices, and
//	defensive input handling.
//
//	Each test follows a strict Arrange / Act / Assert structure to ensure
//	clarity of intent and ease of maintenance.
//
// Notes:
//   - Tests use math.Inf(1) to represent unreachable edges.
//   - The testing parameter name is intentionally "test" for clarity.
//   - Variable names favor explicit meaning over brevity.
//   - Floating-point comparisons use exact values where safe.
//
// -----------------------------------------------------------------------------
package floydwarshall

import (
	"math"
	"testing"
)

// TestFloydWarshall_BasicConnectedGraph verifies shortest paths in a small,
// connected weighted graph.
func TestFloydWarshall_BasicConnectedGraph(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyList := [][]float64{
		{0, 3, infinity, 7},
		{8, 0, 2, infinity},
		{5, infinity, 0, 1},
		{2, infinity, infinity, 0},
	}

	expectedShortestDistances := [][]float64{
		{0, 3, 5, 6},
		{5, 0, 2, 3},
		{3, 6, 0, 1},
		{2, 5, 7, 0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	for source := 0; source < len(shortestDistances); source++ {
		for destination := 0; destination < len(shortestDistances[source]); destination++ {
			if shortestDistances[source][destination] != expectedShortestDistances[source][destination] {
				test.Fatalf("incorrect distance from vertex %d to %d: expected %f, got %f",
					source,
					destination,
					expectedShortestDistances[source][destination],
					shortestDistances[source][destination])
			}
		}
	}
}

// TestFloydWarshall_UnreachableVertices verifies that unreachable paths remain
// infinite after algorithm execution.
func TestFloydWarshall_UnreachableVertices(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyList := [][]float64{
		{0, 1, infinity},
		{infinity, 0, infinity},
		{infinity, infinity, 0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	for source := 0; source < len(shortestDistances); source++ {
		for destination := 0; destination < len(shortestDistances[source]); destination++ {
			// Skip known reachable paths and self-distances.
			if (source == 0 && destination == 1) || source == destination {
				continue
			}

			if shortestDistances[source][destination] != infinity {
				test.Fatalf("expected unreachable path from vertex %d to %d to remain infinite",
					source,
					destination,
				)
			}
		}
	}
}

// TestFloydWarshall_EmptyMatrix verifies that an empty graph returns an empty
// result matrix.
func TestFloydWarshall_EmptyMatrix(test *testing.T) {
	// Arrange.
	adjacencyList := [][]float64{}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if len(shortestDistances) != 0 {
		test.Fatalf("expected empty result matrix, got length %d",
			len(shortestDistances))
	}
}

// TestFloydWarshall_NilMatrix verifies that a nil matrix input is handled
// defensively and results in an empty matrix.
func TestFloydWarshall_NilMatrix(test *testing.T) {
	// Arrange.
	var adjacencyList [][]float64

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if shortestDistances == nil || len(shortestDistances) != 0 {
		test.Fatalf("expected empty non-nil matrix for nil input")
	}
}

// TestFloydWarshall_IntermediatePathDiscovery verifies that shorter paths
// are discovered via intermediate vertices.
func TestFloydWarshall_IntermediatePathDiscovery(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyList := [][]float64{
		{0, 5, infinity},
		{infinity, 0, 2},
		{infinity, infinity, 0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if shortestDistances[0][2] != 7 {
		test.Fatalf("expected path from vertex 0 to 2 to have distance 7, got %f",
			shortestDistances[0][2])
	}
}

// TestFloydWarshall_DirectedGraph verifies that directionality is preserved
// in the computed shortest paths.
func TestFloydWarshall_DirectedGraph(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyList := [][]float64{
		{0, 1, infinity},
		{infinity, 0, infinity},
		{infinity, 2, 0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if shortestDistances[2][0] != infinity {
		test.Fatalf("expected no path from vertex 2 to 0, got %f",
			shortestDistances[2][0])
	}

	if shortestDistances[0][1] != 1 {
		test.Fatalf("expected path from vertex 0 to 1 to have distance 1, got %f",
			shortestDistances[0][1])
	}
}

// TestFloydWarshall_ZeroWeightEdges verifies that zero-weight edges are
// handled correctly.
func TestFloydWarshall_ZeroWeightEdges(test *testing.T) {
	// Arrange.
	adjacencyList := [][]float64{
		{0, 0, 5},
		{0, 0, 1},
		{5, 1, 0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if shortestDistances[0][2] != 1 {
		test.Fatalf("expected shortest path from vertex 0 to 2 to be 1, got %f",
			shortestDistances[0][2])
	}
}

// TestFloydWarshall_NegativeEdgesWithoutCycle verifies that negative edge
// weights are handled correctly when no negative cycle exists.
func TestFloydWarshall_NegativeEdgesWithoutCycle(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyList := [][]float64{
		{0, -2, infinity},
		{infinity, 0, 3},
		{infinity, infinity, 0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if shortestDistances[0][2] != 1 {
		test.Fatalf("expected shortest path from vertex 0 to 2 to be 1, got %f",
			shortestDistances[0][2])
	}
}

// TestFloydWarshall_NegativeCycle verifies that the algorithm panics when
// a negative cycle is present.
func TestFloydWarshall_NegativeCycle(test *testing.T) {
	// Arrange.
	defer func() {
		if recover() == nil {
			test.Fatalf("expected panic due to negative cycle, but no panic occurred")
		}
	}()

	adjacencyList := [][]float64{
		{0, -1, math.Inf(1)},
		{math.Inf(1), 0, -1},
		{-1, math.Inf(1), 0},
	}

	// Act.
	FloydWarshall(adjacencyList)

	// Assert.
	// Panic is expected before reaching this point.
}

// TestFloydWarshall_SingleVertex verifies that a graph with a single vertex
// returns the correct distance of zero.
func TestFloydWarshall_SingleVertex(test *testing.T) {
	// Arrange.
	adjacencyList := [][]float64{
		{0},
	}

	// Act.
	shortestDistances := FloydWarshall(adjacencyList)

	// Assert.
	if len(shortestDistances) != 1 || len(shortestDistances[0]) != 1 {
		test.Fatalf("expected 1x1 result matrix, got dimensions %dx%d",
			len(shortestDistances), len(shortestDistances[0]))
	}

	if shortestDistances[0][0] != 0 {
		test.Fatalf("expected self-distance to be 0, got %f", shortestDistances[0][0])
	}
}

// TestFloydWarshall_Immutability verifies that the input adjacency matrix
// is treated as immutable and is not modified during computation.
func TestFloydWarshall_Immutability(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyList := [][]float64{
		{0, 3, infinity},
		{infinity, 0, 1},
		{2, infinity, 0},
	}

	// Create a structural snapshot to compare against after execution.
	originalMatrixSnapshot := [][]float64{
		{0, 3, infinity},
		{infinity, 0, 1},
		{2, infinity, 0},
	}

	// Act.
	_ = FloydWarshall(adjacencyList)

	// Assert.
	for sourceIndex := 0; sourceIndex < len(adjacencyList); sourceIndex++ {
		for destinationIndex := 0; destinationIndex < len(adjacencyList[sourceIndex]); destinationIndex++ {
			if adjacencyList[sourceIndex][destinationIndex] != originalMatrixSnapshot[sourceIndex][destinationIndex] {
				test.Fatalf("input matrix was mutated at [%d][%d]: expected %f, got %f",
					sourceIndex,
					destinationIndex,
					originalMatrixSnapshot[sourceIndex][destinationIndex],
					adjacencyList[sourceIndex][destinationIndex],
				)
			}
		}
	}
}

// TestFloydWarshall_NilRowPanic verifies that a panic occurs when an adjacency
// matrix contains a nil row.
func TestFloydWarshall_NilRowPanic(test *testing.T) {
	// Arrange.
	defer func() {
		if recoveredPanic := recover(); recoveredPanic == nil {
			test.Fatalf("expected panic due to nil row in adjacency matrix, but no panic occurred")
		}
	}()

	adjacencyList := [][]float64{
		{0, 1},
		nil,
	}

	// Act.
	FloydWarshall(adjacencyList)

	// Assert.
	// Panic is expected before reaching this point.
}

// TestFloydWarshall_NonSquareMatrixPanic verifies that a panic occurs when
// the adjacency matrix is not square (e.g., row length mismatch).
func TestFloydWarshall_NonSquareMatrixPanic(test *testing.T) {
	// Arrange.
	defer func() {
		if recoveredPanic := recover(); recoveredPanic == nil {
			test.Fatalf("expected panic due to non-square adjacency matrix, but no panic occurred")
		}
	}()

	adjacencyList := [][]float64{
		{0, 1, 2},
		{3, 4, 5},
	}

	// Act.
	FloydWarshall(adjacencyList)

	// Assert.
	// Panic is expected before reaching this point.
}

// TestFloydWarshall_DisconnectedComponents verifies that the algorithm correctly
// computes shortest paths across completely disconnected components within a graph.
func TestFloydWarshall_DisconnectedComponents(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyMatrix := [][]float64{
		{0, 2, infinity, infinity},
		{2, 0, infinity, infinity},
		{infinity, infinity, 0, 5},
		{infinity, infinity, 5, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	if shortestDistanceMatrix[0][2] != infinity {
		test.Fatalf("expected infinite distance between disconnected components from vertex 0 to 2, got %f",
			shortestDistanceMatrix[0][2])
	}

	if shortestDistanceMatrix[0][1] != 2 {
		test.Fatalf("expected distance between connected vertex 0 and 1 to be 2, got %f",
			shortestDistanceMatrix[0][1])
	}
}

// TestFloydWarshall_SelfLoopPositive verifies that positive self-loops on a vertex
// are handled correctly and normalized properly by the algorithm.
func TestFloydWarshall_SelfLoopPositive(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyMatrix := [][]float64{
		{5, 1, infinity},
		{infinity, 0, 2},
		{infinity, infinity, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	if shortestDistanceMatrix[0][0] != 0 {
		test.Fatalf("expected positive self-loop distance to normalize to 0, got %f",
			shortestDistanceMatrix[0][0])
	}
}

// TestFloydWarshall_LargeSymmetricGraph verifies path correctness on a larger
// fully connected symmetric graph.
func TestFloydWarshall_LargeSymmetricGraph(test *testing.T) {
	// Arrange.
	numberOfVertices := 5
	adjacencyMatrix := make([][]float64, numberOfVertices)

	for sourceIndex := 0; sourceIndex < numberOfVertices; sourceIndex++ {
		adjacencyMatrix[sourceIndex] = make([]float64, numberOfVertices)
		for destinationIndex := 0; destinationIndex < numberOfVertices; destinationIndex++ {
			if sourceIndex == destinationIndex {
				adjacencyMatrix[sourceIndex][destinationIndex] = 0
			} else {
				adjacencyMatrix[sourceIndex][destinationIndex] = float64(sourceIndex + destinationIndex + 1)
			}
		}
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	for sourceIndex := 0; sourceIndex < numberOfVertices; sourceIndex++ {
		for destinationIndex := 0; destinationIndex < numberOfVertices; destinationIndex++ {
			if shortestDistanceMatrix[sourceIndex][destinationIndex] > adjacencyMatrix[sourceIndex][destinationIndex] {
				test.Fatalf("shortest distance exceeds direct edge weight at [%d][%d]",
					sourceIndex, destinationIndex)
			}
		}
	}
}

// TestFloydWarshall_LinearChainGraph verifies shortest paths in a linear chain graph topology.
func TestFloydWarshall_LinearChainGraph(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyMatrix := [][]float64{
		{0, 2, infinity, infinity},
		{infinity, 0, 3, infinity},
		{infinity, infinity, 0, 4},
		{infinity, infinity, infinity, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	if shortestDistanceMatrix[0][3] != 9 {
		test.Fatalf("expected cumulative linear chain path distance from vertex 0 to 3 to be 9, got %f",
			shortestDistanceMatrix[0][3])
	}
}

// TestFloydWarshall_StarTopologyGraph verifies shortest paths in a star network topology graph.
func TestFloydWarshall_StarTopologyGraph(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyMatrix := [][]float64{
		{0, 1, 1, 1},
		{1, 0, infinity, infinity},
		{1, infinity, 0, infinity},
		{1, infinity, infinity, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	if shortestDistanceMatrix[1][2] != 2 {
		test.Fatalf("expected shortest path through hub from vertex 1 to 2 to be 2, got %f",
			shortestDistanceMatrix[1][2])
	}
}

// TestFloydWarshall_AllInfiniteEdges verifies behavior when all non-diagonal edges are infinite.
func TestFloydWarshall_AllInfiniteEdges(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyMatrix := [][]float64{
		{0, infinity, infinity},
		{infinity, 0, infinity},
		{infinity, infinity, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	for sourceIndex := 0; sourceIndex < len(shortestDistanceMatrix); sourceIndex++ {
		for destinationIndex := 0; destinationIndex < len(shortestDistanceMatrix[sourceIndex]); destinationIndex++ {
			if sourceIndex == destinationIndex {
				if shortestDistanceMatrix[sourceIndex][destinationIndex] != 0 {
					test.Fatalf("expected self-distance to be 0 at index %d", sourceIndex)
				}
			} else {
				if shortestDistanceMatrix[sourceIndex][destinationIndex] != infinity {
					test.Fatalf("expected unreachable edge to remain infinite at [%d][%d]", sourceIndex, destinationIndex)
				}
			}
		}
	}
}

// TestFloydWarshall_NegativeSelfLoopPanic verifies that a negative self-loop triggers the negative cycle panic.
func TestFloydWarshall_NegativeSelfLoopPanic(test *testing.T) {
	// Arrange.
	defer func() {
		if recoveredPanic := recover(); recoveredPanic == nil {
			test.Fatalf("expected panic due to negative self-loop cycle, but no panic occurred")
		}
	}()

	adjacencyMatrix := [][]float64{
		{-1, 2},
		{3, 0},
	}

	// Act.
	FloydWarshall(adjacencyMatrix)

	// Assert.
	// Panic is expected before reaching this point.
}

// TestFloydWarshall_MultipleComponentsDisconnected verifies isolation between multiple disconnected clusters.
func TestFloydWarshall_MultipleComponentsDisconnected(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)

	adjacencyMatrix := [][]float64{
		{0, 3, infinity, infinity, infinity},
		{3, 0, infinity, infinity, infinity},
		{infinity, infinity, 0, 1, 2},
		{infinity, infinity, 1, 0, infinity},
		{infinity, infinity, 2, infinity, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	if shortestDistanceMatrix[0][4] != infinity {
		test.Fatalf("expected no path between separate disconnected components, got %f", shortestDistanceMatrix[0][4])
	}
}

// TestFloydWarshall_MultipleAlternativePaths verifies that the algorithm correctly
// selects the optimal path when multiple competing routes exist between vertices.
func TestFloydWarshall_MultipleAlternativePaths(test *testing.T) {
	// Arrange.
	infinity := math.Inf(1)
	adjacencyMatrix := [][]float64{
		{0, 10, 3, infinity},
		{infinity, 0, 1, 2},
		{infinity, 4, 0, 8},
		{infinity, infinity, infinity, 0},
	}

	// Act.
	shortestDistanceMatrix := FloydWarshall(adjacencyMatrix)

	// Assert.
	expectedDistanceZeroToOne := float64(7)

	if shortestDistanceMatrix[0][1] != expectedDistanceZeroToOne {
		test.Fatalf("expected optimal path distance from vertex 0 to 1 to be %f, got %f",
			expectedDistanceZeroToOne, shortestDistanceMatrix[0][1])
	}
}

// FuzzFloydWarshall performs native Go fuzz testing on the Floyd-Warshall algorithm
// using randomized byte streams converted into edge weights to verify stability and prevent panics.
func FuzzFloydWarshall(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 1, 2, 3})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawInputBytes []byte) {
		if len(rawInputBytes) == 0 {
			return
		}

		dimensionSize := int(math.Sqrt(float64(len(rawInputBytes))))

		if dimensionSize == 0 {
			dimensionSize = 1
		}

		adjacencyMatrix := make([][]float64, dimensionSize)
		byteIndex := 0

		// Act.
		for rowIndex := 0; rowIndex < dimensionSize; rowIndex++ {
			adjacencyMatrix[rowIndex] = make([]float64, dimensionSize)
			for colIndex := 0; colIndex < dimensionSize; colIndex++ {
				if byteIndex < len(rawInputBytes) {
					currentByteValue := rawInputBytes[byteIndex]

					byteIndex++

					if currentByteValue%25 == 0 {
						adjacencyMatrix[rowIndex][colIndex] = math.Inf(1)
					} else {
						adjacencyMatrix[rowIndex][colIndex] = float64(currentByteValue)
					}
				} else {
					adjacencyMatrix[rowIndex][colIndex] = 1.0
				}
			}
		}

		// Execute algorithm wrapped in recovery to catch expected panics safely without crashing the fuzzer.
		func() {
			defer func() {
				_ = recover()
			}()

			_ = FloydWarshall(adjacencyMatrix)
		}()
	})
}

// BenchmarkFloydWarshall_SmallGraph benchmarks the performance of the algorithm
// with a small graph of 10 vertices.
func BenchmarkFloydWarshall_SmallGraph(bench *testing.B) {
	// Arrange.
	numberOfVertices := 10
	infinity := math.Inf(1)
	adjacencyMatrix := make([][]float64, numberOfVertices)

	for sourceIndex := 0; sourceIndex < numberOfVertices; sourceIndex++ {
		adjacencyMatrix[sourceIndex] = make([]float64, numberOfVertices)

		for destinationIndex := 0; destinationIndex < numberOfVertices; destinationIndex++ {
			if sourceIndex == destinationIndex {
				adjacencyMatrix[sourceIndex][destinationIndex] = 0
			} else if (sourceIndex+1)%numberOfVertices == destinationIndex {
				adjacencyMatrix[sourceIndex][destinationIndex] = 2.5
			} else {
				adjacencyMatrix[sourceIndex][destinationIndex] = infinity
			}
		}
	}

	// Act.
	for bench.Loop() {
		_ = FloydWarshall(adjacencyMatrix)
	}
}

// BenchmarkFloydWarshall_MediumGraph benchmarks the performance of the algorithm
// with a medium-sized graph of 50 vertices.
func BenchmarkFloydWarshall_MediumGraph(bench *testing.B) {
	// Arrange.
	numberOfVertices := 50
	infinity := math.Inf(1)
	adjacencyMatrix := make([][]float64, numberOfVertices)

	for sourceIndex := 0; sourceIndex < numberOfVertices; sourceIndex++ {
		adjacencyMatrix[sourceIndex] = make([]float64, numberOfVertices)

		for destinationIndex := 0; destinationIndex < numberOfVertices; destinationIndex++ {
			if sourceIndex == destinationIndex {
				adjacencyMatrix[sourceIndex][destinationIndex] = 0
			} else if (sourceIndex+destinationIndex)%3 == 0 {
				adjacencyMatrix[sourceIndex][destinationIndex] = float64((sourceIndex + destinationIndex) % 10)
			} else {
				adjacencyMatrix[sourceIndex][destinationIndex] = infinity
			}
		}
	}

	// Act.
	for bench.Loop() {
		_ = FloydWarshall(adjacencyMatrix)
	}
}

// BenchmarkFloydWarshall_LargeGraph benchmarks the performance of the algorithm
// with a larger graph of 100 vertices.
func BenchmarkFloydWarshall_LargeGraph(bench *testing.B) {
	// Arrange.
	numberOfVertices := 100
	infinity := math.Inf(1)
	adjacencyMatrix := make([][]float64, numberOfVertices)

	for sourceIndex := 0; sourceIndex < numberOfVertices; sourceIndex++ {
		adjacencyMatrix[sourceIndex] = make([]float64, numberOfVertices)

		for destinationIndex := 0; destinationIndex < numberOfVertices; destinationIndex++ {
			if sourceIndex == destinationIndex {
				adjacencyMatrix[sourceIndex][destinationIndex] = 0
			} else if (sourceIndex+destinationIndex)%7 == 0 {
				adjacencyMatrix[sourceIndex][destinationIndex] = float64((sourceIndex+destinationIndex)%5) + 1
			} else {
				adjacencyMatrix[sourceIndex][destinationIndex] = infinity
			}
		}
	}

	// Act.
	for bench.Loop() {
		_ = FloydWarshall(adjacencyMatrix)
	}
}
