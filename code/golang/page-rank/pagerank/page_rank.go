// -----------------------------------------------------------------------------
// Package: pagerank
// File: pagerank.go
//
// Description:
//
//	This file implements the Google PageRank algorithm using the power iteration
//	method. It computes the relative importance of nodes in a directed graph
//	by treating links as votes, incorporating a damping factor to simulate
//	random web surfer jumps, and handling dangling nodes (pages with no outgoing
//	links) by redistributing their rank evenly across the entire graph.
//
//	The implementation includes defensive checks for parameters, automatic
//	discovery of target nodes omitted from graph keys, and an early termination
//	check based on convergence tolerance.
//
// Notes:
//   - The graph is represented as an adjacency list mapping node IDs to outgoing links.
//   - Convergence is measured using the absolute difference sum across all ranks.
//
// -----------------------------------------------------------------------------
package pagerank

import "math"

// Graph represents a web graph as an adjacency list where keys are node IDs
// and values are slices representing outgoing links.
type Graph map[string][]string

// GooglePageRank computes the PageRank scores for nodes in a directed graph
// using the power iteration method. It returns a map of node IDs to their
// calculated PageRank values, or nil if validation checks fail.
func GooglePageRank(graph Graph, dampingFactor float64, maximumIterations int, tolerance float64) map[string]float64 {
	// Validate that the damping factor falls within the valid probability range of zero to one.
	if dampingFactor < 0.0 || dampingFactor > 1.0 {
		return nil
	}

	// Ensure maximum iterations are positive and tolerance is non-negative.
	if maximumIterations <= 0 || tolerance < 0.0 {
		return nil
	}

	// Verify that the input graph is not empty.
	if len(graph) == 0 {
		return nil
	}

	// Initialize a set to collect all unique nodes, including dangling targets not defined as keys.
	nodeSet := make(map[string]bool)

	// Scan all graph keys and their target links to discover every unique node.
	for node, targets := range graph {
		nodeSet[node] = true

		for _, target := range targets {
			nodeSet[target] = true
		}
	}

	// Flatten the unique node set into a slice for indexed iteration.
	nodes := make([]string, 0, len(nodeSet))

	for node := range nodeSet {
		nodes = append(nodes, node)
	}

	// Initialize lookup maps to track incoming links and out-degrees for every node.
	incoming := make(map[string][]string)
	outDegree := make(map[string]int)

	// Ensure every discovered node has an entry in the out-degree and incoming structures.
	for _, node := range nodes {
		incoming[node] = []string{}
		outDegree[node] = 0
	}

	// Populate out-degrees and incoming link lists based on the graph edges.
	for node, targets := range graph {
		outDegree[node] = len(targets)

		for _, target := range targets {
			incoming[target] = append(incoming[target], node)
		}
	}

	// Convert the total node count to a floating-point value for rank calculations.
	numberOfNodes := float64(len(nodes))

	// Guard against a zero node count to prevent division-by-zero errors.
	if numberOfNodes == 0 {
		return nil
	}

	// Initialize all node ranks uniformly to an equal probability share.
	ranks := make(map[string]float64)

	for _, node := range nodes {
		ranks[node] = 1.0 / numberOfNodes
	}

	// Compute the base rank contributed to every node by random surfer teleports.
	baseRank := (1.0 - dampingFactor) / numberOfNodes

	// Iterate up to the maximum limit or until convergence criteria are met.
	for iteration := 0; iteration < maximumIterations; iteration++ {
		// Allocate a map to store updated rank values for the current iteration.
		newRanks := make(map[string]float64)

		// Track the cumulative rank of all dangling nodes with zero outgoing links.
		danglingSum := 0.0

		// Accumulate the rank scores of all dangling nodes to redistribute them evenly.
		for _, node := range nodes {
			if outDegree[node] == 0 {
				danglingSum += ranks[node]
			}
		}

		// Initialize the convergence difference metric for this iteration step.
		difference := 0.0

		// Compute the new rank for each node using incoming link contributions and redistributed dangling ranks.
		for _, node := range nodes {
			sum := 0.0

			// Sum the normalized ranks contributed by all incoming neighbor nodes.
			for _, incomingNode := range incoming[node] {
				// Only consider incoming nodes that have valid out-degrees.
				if outDegree[incomingNode] > 0 {
					sum += ranks[incomingNode] / float64(outDegree[incomingNode])
				}
			}

			// Combine the base teleport probability, incoming link votes, and redistributed dangling rank.
			newRanks[node] = baseRank + dampingFactor*(sum+danglingSum/numberOfNodes)
		}

		// Update the global rank values and accumulate the absolute difference for convergence tracking.
		for _, node := range nodes {
			difference += math.Abs(newRanks[node] - ranks[node])
			ranks[node] = newRanks[node]
		}

		// Terminate early if the rank changes have converged below the tolerance threshold.
		if difference < tolerance {
			break
		}
	}

	// Return the final computed PageRank scores.
	return ranks
}
