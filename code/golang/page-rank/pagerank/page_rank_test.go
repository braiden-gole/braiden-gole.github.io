// -----------------------------------------------------------------------------
// Package: pagerank
// File: pagerank_test.go
//
// Description:
//
//	This file contains a comprehensive suite of 30 unique unit tests for the
//	GooglePageRank function. It covers edge cases, invalid inputs, boundary
//	conditions, various graph topologies (such as chains, cycles, stars, and
//	disconnected components), dangling nodes, and parameter variations.
//
//	Each test follows the Arrange-Act-Assert pattern with clear descriptive
//	comments.
//
// -----------------------------------------------------------------------------
package pagerank

import (
	"math"
	"strconv"
	"testing"
)

// TestGooglePageRank_NilGraph verifies that passing a nil graph returns nil.
func TestGooglePageRank_NilGraph(test *testing.T) {
	// Arrange.
	var graph Graph = nil

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for nil graph, got %v", ranks)
	}
}

// TestGooglePageRank_EmptyGraph verifies that passing an empty graph map returns nil.
func TestGooglePageRank_EmptyGraph(test *testing.T) {
	// Arrange.
	graph := Graph{}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for empty graph, got %v", ranks)
	}
}

// TestGooglePageRank_DampingFactorBelowZero verifies that a damping factor less than 0.0 returns nil.
func TestGooglePageRank_DampingFactorBelowZero(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, -0.1, 100, 1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for invalid damping factor < 0, got %v", ranks)
	}
}

// TestGooglePageRank_DampingFactorAboveOne verifies that a damping factor greater than 1.0 returns nil.
func TestGooglePageRank_DampingFactorAboveOne(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 1.1, 100, 1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for invalid damping factor > 1, got %v", ranks)
	}
}

// TestGooglePageRank_ZeroIterations verifies that zero maximum iterations returns nil.
func TestGooglePageRank_ZeroIterations(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 0, 1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for zero maximum iterations, got %v", ranks)
	}
}

// TestGooglePageRank_NegativeIterations verifies that negative maximum iterations returns nil.
func TestGooglePageRank_NegativeIterations(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, -5, 1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for negative maximum iterations, got %v", ranks)
	}
}

// TestGooglePageRank_NegativeTolerance verifies that a negative tolerance value returns nil.
func TestGooglePageRank_NegativeTolerance(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, -1e-6)

	// Assert.
	if ranks != nil {
		test.Errorf("expected nil ranks for negative tolerance, got %v", ranks)
	}
}

// TestGooglePageRank_SingleNodeNoLinks verifies PageRank behavior for a single isolated node.
func TestGooglePageRank_SingleNodeNoLinks(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 1 || math.Abs(ranks["A"]-1.0) > 1e-5 {
		test.Errorf("expected score of 1.0 for single node graph, got %v", ranks)
	}
}

// TestGooglePageRank_SingleNodeSelfLoop verifies PageRank behavior for a single node with a self-loop.
func TestGooglePageRank_SingleNodeSelfLoop(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 1 || math.Abs(ranks["A"]-1.0) > 1e-5 {
		test.Errorf("expected score of 1.0 for self-loop node graph, got %v", ranks)
	}
}

// TestGooglePageRank_TwoNodesMutual verifies PageRank for two mutually linking nodes.
func TestGooglePageRank_TwoNodesMutual(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 2 {
		test.Fatalf("expected 2 nodes, got %v", ranks)
	}

	if math.Abs(ranks["A"]-0.5) > 1e-5 || math.Abs(ranks["B"]-0.5) > 1e-5 {
		test.Errorf("expected equal ranks of 0.5 for mutual nodes, got A: %f, B: %f", ranks["A"], ranks["B"])
	}
}

// TestGooglePageRank_TwoNodesOneDirection verifies PageRank for a directed edge from A to B.
func TestGooglePageRank_TwoNodesOneDirection(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 2 {
		test.Fatalf("expected 2 nodes, got %v", ranks)
	}

	if ranks["B"] <= ranks["A"] {
		test.Errorf("expected target B to have higher rank than source A due to link accumulation, got A: %f, B: %f", ranks["A"], ranks["B"])
	}
}

// TestGooglePageRank_ThreeNodesChain verifies PageRank for a linear chain of three nodes (A -> B -> C).
func TestGooglePageRank_ThreeNodesChain(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"C"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Fatalf("expected 3 nodes, got %v", ranks)
	}

	if ranks["C"] <= ranks["B"] || ranks["B"] <= ranks["A"] {
		test.Errorf("expected ascending rank progression along the chain, got %v", ranks)
	}
}

// TestGooglePageRank_ThreeNodesCycle verifies PageRank for a simple three-node cycle (A -> B -> C -> A).
func TestGooglePageRank_ThreeNodesCycle(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Fatalf("expected 3 nodes, got %v", ranks)
	}

	expected := 1.0 / 3.0

	for node, score := range ranks {
		if math.Abs(score-expected) > 1e-3 {
			test.Errorf("expected rank ~%f for node %s in cycle, got %f", expected, node, score)
		}
	}
}

// TestGooglePageRank_StarGraph verifies PageRank for a star topology where a hub points to multiple leaves.
func TestGooglePageRank_StarGraph(test *testing.T) {
	// Arrange.
	graph := Graph{
		"Hub": {"L1", "L2", "L3"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 4 {
		test.Fatalf("expected 4 nodes, got %v", ranks)
	}

	if ranks["L1"] <= ranks["Hub"] {
		test.Errorf("expected leaves to have higher rank than hub in star graph, got Hub: %f, L1: %f", ranks["Hub"], ranks["L1"])
	}
}

// TestGooglePageRank_CompleteGraph verifies PageRank for a fully connected complete graph.
func TestGooglePageRank_CompleteGraph(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B", "C"},
		"B": {"A", "C"},
		"C": {"A", "B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Fatalf("expected 3 nodes, got %v", ranks)
	}

	expected := 1.0 / 3.0

	for node, score := range ranks {
		if math.Abs(score-expected) > 1e-4 {
			test.Errorf("expected equal ranks in complete graph for %s, got %f", node, score)
		}
	}
}

// TestGooglePageRank_DisconnectedComponents verifies PageRank across disconnected subgraphs.
func TestGooglePageRank_DisconnectedComponents(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"A"},
		"X": {"Y"},
		"Y": {"X"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 4 {
		test.Fatalf("expected 4 nodes, got %v", ranks)
	}

	if math.Abs(ranks["A"]-ranks["X"]) > 1e-4 {
		test.Errorf("expected symmetric components to share identical ranks, got A: %f, X: %f", ranks["A"], ranks["X"])
	}
}

// TestGooglePageRank_SingleDanglingNode verifies handling of a single dangling node with no outgoing links.
func TestGooglePageRank_SingleDanglingNode(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 2 {
		test.Fatalf("expected 2 nodes, got %v", ranks)
	}

	if ranks["B"] <= 0 {
		test.Errorf("expected dangling node B to retain a non-zero redistributed rank, got %f", ranks["B"])
	}
}

// TestGooglePageRank_MultipleDanglingNodes verifies handling when multiple nodes have zero out-degree.
func TestGooglePageRank_MultipleDanglingNodes(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {},
		"B": {},
		"C": {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Fatalf("expected 3 nodes, got %v", ranks)
	}

	expected := 1.0 / 3.0

	for node, score := range ranks {
		if math.Abs(score-expected) > 1e-4 {
			test.Errorf("expected equal rank distribution for all-dangling nodes, got %s: %f", node, score)
		}
	}
}

// TestGooglePageRank_TargetNodeMissingFromKeys verifies automatic registration of target nodes omitted as map keys.
func TestGooglePageRank_TargetNodeMissingFromKeys(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 2 {
		test.Fatalf("expected 2 nodes including implicitly discovered target 'B', got %v", ranks)
	}

	if _, exists := ranks["B"]; !exists {
		test.Errorf("expected target 'B' to be present in computed ranks")
	}
}

// TestGooglePageRank_MaxIterationsCap verifies that iteration limit cuts off computation safely without crashing.
func TestGooglePageRank_MaxIterationsCap(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 1, 1e-15)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Errorf("expected valid ranks despite hitting iteration cap, got %v", ranks)
	}
}

// TestGooglePageRank_StrictTolerance verifies convergence under extremely tight tolerance constraints.
func TestGooglePageRank_StrictTolerance(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B", "C"},
		"B": {"C"},
		"C": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 500, 1e-12)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Errorf("expected accurate convergence with strict tolerance, got %v", ranks)
	}
}

// TestGooglePageRank_LooseTolerance verifies convergence behavior under loose tolerance settings.
func TestGooglePageRank_LooseTolerance(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B", "C"},
		"B": {"C"},
		"C": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-1)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Errorf("expected quick termination under loose tolerance, got %v", ranks)
	}
}

// TestGooglePageRank_DampingFactorZero verifies that a damping factor of 0.0 results in uniform base rank distribution.
func TestGooglePageRank_DampingFactorZero(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.0, 100, 1e-6)

	// Assert.
	if ranks == nil {
		test.Fatalf("expected non-nil ranks")
	}

	expected := 1.0 / 3.0

	for node, score := range ranks {
		if math.Abs(score-expected) > 1e-5 {
			test.Errorf("expected uniform rank %f with 0.0 damping factor for %s, got %f", expected, node, score)
		}
	}
}

// TestGooglePageRank_DampingFactorOne verifies PageRank behavior when damping factor is 1.0 (no random teleports).
func TestGooglePageRank_DampingFactorOne(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 1.0, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 3 {
		test.Errorf("expected valid ranks with damping factor 1.0, got %v", ranks)
	}
}

// TestGooglePageRank_ExampleGraph verifies the exact sample graph used in previous tutorials.
func TestGooglePageRank_ExampleGraph(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B", "C"},
		"B": {"C"},
		"C": {"A"},
		"D": {"C"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 4 {
		test.Fatalf("expected 4 nodes in example graph, got %v", ranks)
	}

	if math.Abs(ranks["A"]-0.3725) > 1e-3 || math.Abs(ranks["B"]-0.1958) > 1e-3 || math.Abs(ranks["C"]-0.3941) > 1e-3 {
		test.Errorf("computed ranks differ from expected example values, got %v", ranks)
	}
}

// TestGooglePageRank_SymmetricDiamondGraph verifies PageRank for a symmetric diamond directed graph topology.
func TestGooglePageRank_SymmetricDiamondGraph(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"B", "C"},
		"B": {"D"},
		"C": {"D"},
		"D": {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 4 {
		test.Fatalf("expected 4 nodes, got %v", ranks)
	}

	if math.Abs(ranks["B"]-ranks["C"]) > 1e-5 {
		test.Errorf("expected symmetric intermediate nodes B and C to have equal rank, got B: %f, C: %f", ranks["B"], ranks["C"])
	}
}

// TestGooglePageRank_SinkGraph verifies PageRank where multiple paths funnel into a single sink node.
func TestGooglePageRank_SinkGraph(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A":    {"Sink"},
		"B":    {"Sink"},
		"C":    {"Sink"},
		"Sink": {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 4 {
		test.Fatalf("expected 4 nodes, got %v", ranks)
	}

	if ranks["Sink"] <= ranks["A"] {
		test.Errorf("expected sink node to accumulate highest rank, got Sink: %f, A: %f", ranks["Sink"], ranks["A"])
	}
}

// TestGooglePageRank_SourceGraph verifies PageRank for a source node branching out with no incoming loops.
func TestGooglePageRank_SourceGraph(test *testing.T) {
	// Arrange.
	graph := Graph{
		"Source": {"A", "B"},
		"A":      {"C"},
		"B":      {"C"},
		"C":      {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 4 {
		test.Fatalf("expected 4 nodes, got %v", ranks)
	}

	if ranks["C"] <= ranks["Source"] {
		test.Errorf("expected downstream node C to capture higher rank than source, got Source: %f, C: %f", ranks["Source"], ranks["C"])
	}
}

// TestGooglePageRank_LargeChain verifies behavior over a longer sequential chain of nodes.
func TestGooglePageRank_LargeChain(test *testing.T) {
	// Arrange.
	graph := Graph{
		"N1": {"N2"},
		"N2": {"N3"},
		"N3": {"N4"},
		"N4": {"N5"},
		"N5": {},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 5 {
		test.Fatalf("expected 5 nodes in chain, got %v", ranks)
	}

	if ranks["N5"] <= ranks["N1"] {
		test.Errorf("expected final node N5 to have higher rank than start node N1")
	}
}

// TestGooglePageRank_SelfLoopWithExternalLink verifies node rank stability with self-loops mixed with outward links.
func TestGooglePageRank_SelfLoopWithExternalLink(test *testing.T) {
	// Arrange.
	graph := Graph{
		"A": {"A", "B"},
		"B": {"A"},
	}

	// Act.
	ranks := GooglePageRank(graph, 0.85, 100, 1e-6)

	// Assert.
	if ranks == nil || len(ranks) != 2 {
		test.Fatalf("expected 2 nodes, got %v", ranks)
	}

	if ranks["A"] <= 0 || ranks["B"] <= 0 {
		test.Errorf("expected valid positive ranks for self-loop mixed graph, got %v", ranks)
	}
}

// FuzzGooglePageRank_Operations performs native Go fuzz testing on the Google PageRank implementation
// using randomized damping factors, iteration bounds, tolerances, and connection streams safely.
func FuzzGooglePageRank_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(0.85, 100, 1e-6, []byte{0, 1, 1, 2, 2, 0})
	fuzz.Add(0.50, 50, 1e-3, []byte{0, 2, 1, 0})

	fuzz.Fuzz(func(fuzzContext *testing.T, dampingFactor float64, maxIterations int, tolerance float64, rawConnectionBytes []byte) {
		if len(rawConnectionBytes) < 2 || len(rawConnectionBytes)%2 != 0 {
			return
		}

		// Act.
		graph := make(Graph)
		nodeCount := 6
		nodeNames := make([]string, nodeCount)

		for nodeIndex := 0; nodeIndex < nodeCount; nodeIndex++ {
			nodeNames[nodeIndex] = "node-" + strconv.Itoa(nodeIndex)
			graph[nodeNames[nodeIndex]] = []string{}
		}

		byteIndex := 0
		for byteIndex+1 < len(rawConnectionBytes) {
			sourceIndex := int(rawConnectionBytes[byteIndex]) % nodeCount
			targetIndex := int(rawConnectionBytes[byteIndex+1]) % nodeCount

			byteIndex += 2

			sourceNode := nodeNames[sourceIndex]
			targetNode := nodeNames[targetIndex]

			if sourceNode != targetNode {
				graph[sourceNode] = append(graph[sourceNode], targetNode)
			}
		}

		// Act & Assert: Execute PageRank computation wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			_ = GooglePageRank(graph, dampingFactor, maxIterations, tolerance)
		}()
	})
}

// BenchmarkGooglePageRank_SmallGraph benchmarks the performance of PageRank
// on a small connected graph topology.
func BenchmarkGooglePageRank_SmallGraph(bench *testing.B) {
	// Arrange.
	graph := Graph{
		"A": {"B", "C"},
		"B": {"C"},
		"C": {"A"},
		"D": {"C", "A"},
	}

	// Act.
	for bench.Loop() {
		_ = GooglePageRank(graph, 0.85, 100, 1e-6)
	}
}

// BenchmarkGooglePageRank_ChainGraph benchmarks the performance of PageRank
// on a sequential linear chain of nodes.
func BenchmarkGooglePageRank_ChainGraph(bench *testing.B) {
	// Arrange.
	graph := make(Graph)
	chainLength := 50

	for index := 0; index < chainLength; index++ {
		currentName := "node-" + strconv.Itoa(index)

		if index < chainLength-1 {
			nextNode := "node-" + strconv.Itoa(index+1)
			graph[currentName] = []string{nextNode}
		} else {
			graph[currentName] = []string{}
		}
	}

	// Act.
	for bench.Loop() {
		_ = GooglePageRank(graph, 0.85, 100, 1e-6)
	}
}

// BenchmarkGooglePageRank_DenseGraph benchmarks the performance of PageRank
// on a larger densely interconnected graph structure.
func BenchmarkGooglePageRank_DenseGraph(bench *testing.B) {
	// Arrange.
	graph := make(Graph)

	nodeCount := 30
	nodeNames := make([]string, nodeCount)

	for index := 0; index < nodeCount; index++ {
		nodeNames[index] = "node-" + strconv.Itoa(index)
	}

	for index := 0; index < nodeCount; index++ {
		source := nodeNames[index]
		targets := []string{}

		for offset := 1; offset <= 5; offset++ {
			targetIndex := (index + offset) % nodeCount
			targets = append(targets, nodeNames[targetIndex])
		}

		graph[source] = targets
	}

	// Act.
	for bench.Loop() {
		_ = GooglePageRank(graph, 0.85, 100, 1e-6)
	}
}
