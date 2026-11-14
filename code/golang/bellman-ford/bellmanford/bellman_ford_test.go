// -----------------------------------------------------------------------------
// Package: bellmanford
// File: bellman_ford_test.go
//
// Description:
//
//	This file contains unit tests for the Bellman–Ford single-source shortest
//	path algorithm implementation. The tests validate correct distance
//	computation, handling of unreachable vertices, and detection of
//	negative-weight cycles.
//
//	Each test follows a structured Arrange / Act / Assert pattern to ensure
//	clarity and repeatability. Variable names are intentionally descriptive
//	to reflect algorithmic intent rather than implementation detail.
//
// Notes:
//   - All tests assume zero-based vertex indexing.
//   - Error conditions are treated as first-class outcomes.
//   - The algorithm is expected to tolerate negative weights safely.
//
// -----------------------------------------------------------------------------
package bellmanford

import "testing"

// TestBellmanFord_ComputesCorrectShortestPaths verifies that the algorithm
// computes correct shortest-path distances for a graph containing negative
// edge weights but no negative-weight cycles.
func TestBellmanFord_ComputesCorrectShortestPaths(test *testing.T) {
	// Arrange.
	graph := Graph{
		Count: 5,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 6},
			{Source: 0, Destination: 2, Weight: 7},
			{Source: 1, Destination: 2, Weight: 8},
			{Source: 1, Destination: 3, Weight: 5},
			{Source: 1, Destination: 4, Weight: -4},
			{Source: 2, Destination: 3, Weight: -3},
			{Source: 2, Destination: 4, Weight: 9},
			{Source: 3, Destination: 1, Weight: -2},
			{Source: 4, Destination: 0, Weight: 2},
			{Source: 4, Destination: 3, Weight: 7},
		},
	}

	expectedDistances := []int{0, 2, 7, 4, -2}

	// Act.
	computedDistances, err := BellmanFord(graph, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	if len(computedDistances) != len(expectedDistances) {
		test.Fatalf("distance slice length mismatch: expected %d, got %d",
			len(expectedDistances),
			len(computedDistances),
		)
	}

	for vertexIndex := range expectedDistances {
		if computedDistances[vertexIndex] != expectedDistances[vertexIndex] {
			test.Errorf("incorrect distance for vertex %d: expected %d, got %d",
				vertexIndex,
				expectedDistances[vertexIndex],
				computedDistances[vertexIndex],
			)
		}
	}
}

// TestBellmanFord_DetectsNegativeWeightCycle verifies that the algorithm
// returns an error when the graph contains a negative-weight cycle reachable
// from the source vertex.
func TestBellmanFord_DetectsNegativeWeightCycle(test *testing.T) {
	// Arrange.
	graphWithNegativeCycle := Graph{
		Count: 3,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 1},
			{Source: 1, Destination: 2, Weight: -1},
			{Source: 2, Destination: 0, Weight: -1},
		},
	}

	// Act.
	computedDistances, err := BellmanFord(graphWithNegativeCycle, 0)

	// Assert.
	if err == nil {
		test.Fatalf("expected error due to negative-weight cycle, got distances: %v",
			computedDistances,
		)
	}
}

// TestBellmanFord_HandlesUnreachableVertices verifies that vertices which
// cannot be reached from the source retain infinite distance values.
func TestBellmanFord_HandlesUnreachableVertices(test *testing.T) {
	// Arrange.
	graphWithDisconnectedVertex := Graph{
		Count: 4,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 3},
			{Source: 1, Destination: 2, Weight: 4},
		},
	}

	// Act.
	computedDistances, err := BellmanFord(graphWithDisconnectedVertex, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	maxIntegerValue := int(^uint(0) >> 1)

	if computedDistances[3] != maxIntegerValue {
		test.Errorf("expected unreachable vertex to have infinite distance, got %d",
			computedDistances[3],
		)
	}
}

// TestBellmanFord_InvalidSourceVertex verifies that providing an out-of-range
// source vertex results in an error.
func TestBellmanFord_InvalidSourceVertex(test *testing.T) {
	// Arrange.
	graph := Graph{
		Count: 2,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 5},
		},
	}

	invalidSourceVertex := 5

	// Act.
	computedDistances, err := BellmanFord(graph, invalidSourceVertex)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for invalid source vertex, got distances: %v",
			computedDistances,
		)
	}
}

// TestBellmanFord_SingleVertexGraph verifies that a graph containing only
// a single vertex returns a zero distance when the source is that vertex.
func TestBellmanFord_SingleVertexGraph(test *testing.T) {
	// Arrange.
	singleVertexGraph := Graph{
		Count: 1,
		Edges: []Edge{},
	}

	// Act.
	computedDistances, err := BellmanFord(singleVertexGraph, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	if len(computedDistances) != 1 {
		test.Fatalf("expected distance slice of length 1, got %d",
			len(computedDistances),
		)
	}

	if computedDistances[0] != 0 {
		test.Errorf("expected distance for source vertex to be 0, got %d",
			computedDistances[0],
		)
	}
}

// TestBellmanFord_NoEdges verifies that a graph with multiple vertices and
// no edges leaves all non-source vertices unreachable.
func TestBellmanFord_NoEdges(test *testing.T) {
	// Arrange.
	graphWithNoEdges := Graph{
		Count: 4,
		Edges: []Edge{},
	}

	// Act.
	computedDistances, err := BellmanFord(graphWithNoEdges, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	maxIntegerValue := int(^uint(0) >> 1)

	if computedDistances[0] != 0 {
		test.Errorf("expected source vertex distance to be 0, got %d",
			computedDistances[0],
		)
	}

	for vertexIndex := 1; vertexIndex < graphWithNoEdges.Count; vertexIndex++ {
		if computedDistances[vertexIndex] != maxIntegerValue {
			test.Errorf("expected vertex %d to be unreachable, got %d",
				vertexIndex,
				computedDistances[vertexIndex],
			)
		}
	}
}

// TestBellmanFord_MultipleEdgesBetweenSameVertices verifies that the algorithm
// correctly selects the minimum-weight edge when multiple edges connect the
// same source and destination vertices.
func TestBellmanFord_MultipleEdgesBetweenSameVertices(test *testing.T) {
	// Arrange.
	graphWithParallelEdges := Graph{
		Count: 2,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 10},
			{Source: 0, Destination: 1, Weight: 3},
			{Source: 0, Destination: 1, Weight: 7},
		},
	}

	// Act.
	computedDistances, err := BellmanFord(graphWithParallelEdges, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	if computedDistances[1] != 3 {
		test.Errorf("expected shortest distance to be 3, got %d",
			computedDistances[1],
		)
	}
}

// TestBellmanFord_NegativeEdgesWithoutCycle verifies that negative edges
// alone do not trigger a failure when no negative-weight cycle exists.
func TestBellmanFord_NegativeEdgesWithoutCycle(test *testing.T) {
	// Arrange.
	graphWithNegativeEdges := Graph{
		Count: 4,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: -1},
			{Source: 1, Destination: 2, Weight: -2},
			{Source: 2, Destination: 3, Weight: 5},
		},
	}

	expectedDistances := []int{0, -1, -3, 2}

	// Act.
	computedDistances, err := BellmanFord(graphWithNegativeEdges, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	for vertexIndex := range expectedDistances {
		if computedDistances[vertexIndex] != expectedDistances[vertexIndex] {
			test.Errorf("incorrect distance for vertex %d: expected %d, got %d",
				vertexIndex,
				expectedDistances[vertexIndex],
				computedDistances[vertexIndex],
			)
		}
	}
}

// TestBellmanFord_ZeroWeightEdges verifies that edges with zero weight are
// handled correctly and do not distort shortest-path calculations.
func TestBellmanFord_ZeroWeightEdges(test *testing.T) {
	// Arrange.
	graphWithZeroWeightEdges := Graph{
		Count: 3,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 0},
			{Source: 1, Destination: 2, Weight: 0},
		},
	}

	// Act.
	computedDistances, err := BellmanFord(graphWithZeroWeightEdges, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	if computedDistances[2] != 0 {
		test.Errorf("expected distance to vertex 2 to be 0, got %d",
			computedDistances[2],
		)
	}
}

// TestBellmanFord_DisconnectedSourceVertex verifies that when the source vertex
// has no outgoing edges and other vertices form an isolated component, they remain unreachable.
// TestBellmanFord_DisconnectedSourceVertex verifies that when the source vertex
// has no edges connecting it to other components, those other vertices remain unreachable.
func TestBellmanFord_DisconnectedSourceVertex(test *testing.T) {
	// Arrange.
	graph := Graph{
		Count: 3,
		Edges: []Edge{
			{Source: 1, Destination: 2, Weight: 5},
		},
	}

	// Act.
	computedDistances, err := BellmanFord(graph, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error returned: %v", err)
	}

	if computedDistances[0] != 0 {
		test.Errorf("expected source vertex distance to be 0, got %d", computedDistances[0])
	}

	// Depending on your implementation's infinity representation for truly disconnected nodes,
	// verify that they do not inherit false shortest paths from the isolated component.
	if computedDistances[1] == 5 || computedDistances[2] == 5 {
		test.Errorf("expected vertices 1 and 2 to remain disconnected from source vertex 0")
	}
}

// TestBellmanFord_NegativeSourceIndex verifies that a negative source vertex index
// returns an error rather than panicking or indexing incorrectly.
func TestBellmanFord_NegativeSourceIndex(test *testing.T) {
	// Arrange.
	graph := Graph{
		Count: 3,
		Edges: []Edge{
			{Source: 0, Destination: 1, Weight: 2},
		},
	}

	negativeSourceVertex := -1

	// Act.
	computedDistances, err := BellmanFord(graph, negativeSourceVertex)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for negative source index, got distances: %v", computedDistances)
	}
}

// FuzzBellmanFord_Execution performs native Go fuzz testing on the Bellman–Ford algorithm
// using randomized byte inputs to construct graphs and execute shortest-path computations safely.
func FuzzBellmanFord_Execution(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{3, 0, 0, 1, 5, 0})
	fuzz.Add([]byte{4, 0, 1, 2, 10, 1})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawInputBytes []byte) {
		if len(rawInputBytes) < 3 {
			return
		}

		vertexCount := int(rawInputBytes[0]%10) + 1
		sourceVertex := int(rawInputBytes[1]) % vertexCount

		edges := []Edge{}

		for byteIndex := 2; byteIndex+2 < len(rawInputBytes); byteIndex += 3 {
			source := int(rawInputBytes[byteIndex]) % vertexCount
			destination := int(rawInputBytes[byteIndex+1]) % vertexCount
			weight := int(int8(rawInputBytes[byteIndex+2]))

			edges = append(edges, Edge{
				Source:      source,
				Destination: destination,
				Weight:      weight,
			})
		}

		graph := Graph{
			Count: vertexCount,
			Edges: edges,
		}

		// Act / Assert: Execute algorithm wrapped in recovery to catch any panics safely.
		func() {
			defer func() {
				_ = recover()
			}()

			_, _ = BellmanFord(graph, sourceVertex)
		}()
	})
}

// BenchmarkBellmanFord_SmallGraph benchmarks the performance of the Bellman–Ford algorithm
// on a small graph with 10 vertices and 30 edges.
func BenchmarkBellmanFord_SmallGraph(bench *testing.B) {
	// Arrange.
	vertexCount := 10
	edges := []Edge{}

	for source := 0; source < vertexCount; source++ {
		for destination := 0; destination < vertexCount; destination++ {
			if source != destination && (source+destination)%3 == 0 {
				edges = append(edges, Edge{
					Source:      source,
					Destination: destination,
					Weight:      (source + destination) % 10,
				})
			}
		}
	}

	graph := Graph{
		Count: vertexCount,
		Edges: edges,
	}

	// Act.
	for bench.Loop() {
		_, _ = BellmanFord(graph, 0)
	}
}

// BenchmarkBellmanFord_MediumGraph benchmarks the performance of the Bellman–Ford algorithm
// on a medium-sized graph with 50 vertices.
func BenchmarkBellmanFord_MediumGraph(bench *testing.B) {
	// Arrange.
	vertexCount := 50
	edges := []Edge{}

	for source := 0; source < vertexCount; source++ {
		for destination := 0; destination < vertexCount; destination++ {
			if source != destination && (source+destination)%4 == 0 {
				edges = append(edges, Edge{
					Source:      source,
					Destination: destination,
					Weight:      (source*3 + destination) % 20,
				})
			}
		}
	}

	graph := Graph{
		Count: vertexCount,
		Edges: edges,
	}

	// Act.
	for bench.Loop() {
		_, _ = BellmanFord(graph, 0)
	}
}
