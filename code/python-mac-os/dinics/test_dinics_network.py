"""
test_dinics_network.py

Unit testing suite for Dinic's maximum flow network algorithm implementation.

This module provides exhaustive test coverage validation leveraging Python's
built-in unittest framework. It validates parameter topological constraints,
state mutations across structural execution phases, residual graph capacities,
and behavioral flow benchmarks across asymmetric networks, back-edges, cycles,
and bottleneck distributions.
"""
import unittest

from dinics import DinicsNetwork, FlowEdge

class TestDinicsAlgorithm(unittest.TestCase):
    """
    Comprehensive unit testing suite enforcing structural integrity, mathematical
    correctness, and internal state validation for the DinicsNetwork pipeline.
    """

    def test_01_standard_graph(self):
        """Test a typical 6-node network with a known maximum flow."""
        # Arrange.
        network = DinicsNetwork(6)
        network.add_edge(0, 1, 16)
        network.add_edge(0, 2, 13)
        network.add_edge(1, 2, 10)
        network.add_edge(1, 3, 12)
        network.add_edge(2, 1, 4)
        network.add_edge(2, 4, 14)
        network.add_edge(3, 2, 9)
        network.add_edge(3, 5, 20)
        network.add_edge(4, 3, 7)
        network.add_edge(4, 5, 4)

        # Act.
        result = network.compute_maximum_flow(0, 5)

        # Assert.
        self.assertEqual(result, 23)

    def test_02_single_edge(self):
        """Test the simplest possible graph with only two nodes and one edge."""
        # Arrange.
        network = DinicsNetwork(2)
        network.add_edge(0, 1, 10)

        # Act.
        result = network.compute_maximum_flow(0, 1)

        # Assert.
        self.assertEqual(result, 10)

    def test_03_completely_disconnected(self):
        """Test a graph with no edges at all."""
        # Arrange.
        network = DinicsNetwork(4)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 0)

    def test_04_no_path_to_sink(self):
        """Test a graph where edges exist but do not connect the source to the sink."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 10)
        network.add_edge(2, 3, 10)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 0)

    def test_05_source_equals_sink(self):
        """Test the edge case where the source node is the same as the sink node."""
        # Arrange.
        network = DinicsNetwork(3)
        network.add_edge(0, 1, 10)
        network.add_edge(1, 2, 10)

        # Act.
        result = network.compute_maximum_flow(0, 0)

        # Assert.
        self.assertEqual(result, 0)

    def test_06_zero_capacity_edges(self):
        """Test that paths with an edge of 0 capacity do not transmit flow."""
        # Arrange.
        network = DinicsNetwork(3)
        network.add_edge(0, 1, 0)
        network.add_edge(1, 2, 10)

        # Act.
        result = network.compute_maximum_flow(0, 2)

        # Assert.
        self.assertEqual(result, 0)

    def test_07_bottleneck_at_source(self):
        """Test network constrained entirely by a tiny initial edge capacity."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 2)
        network.add_edge(1, 2, 50)
        network.add_edge(2, 3, 50)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 2)

    def test_08_bottleneck_at_sink(self):
        """Test network constrained entirely by a tiny final edge capacity."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 50)
        network.add_edge(1, 2, 50)
        network.add_edge(2, 3, 2)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 2)

    def test_09_bottleneck_in_middle(self):
        """Test network constrained by a constriction point in the middle layers."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 100)
        network.add_edge(1, 2, 5)
        network.add_edge(2, 3, 100)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 5)

    def test_10_multiple_parallel_paths(self):
        """Test simple parallel paths to ensure capacities aggregate cleanly."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 10)
        network.add_edge(1, 3, 10)
        network.add_edge(0, 2, 15)
        network.add_edge(2, 3, 15)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 25)

    def test_11_diamond_network(self):
        """Test a classic diamond layout where paths split and merge at a shared vertex."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 10)
        network.add_edge(0, 2, 10)
        network.add_edge(1, 2, 5)
        network.add_edge(1, 3, 10)
        network.add_edge(2, 3, 10)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 20)

    def test_12_backwards_edges(self):
        """Test a graph containing forward structural edges that point back toward the source."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 10)
        network.add_edge(1, 2, 10)
        network.add_edge(2, 1, 5)
        network.add_edge(2, 3, 10)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 10)

    def test_13_flow_canceling_required(self):
        """Test a graph layout that strictly requires the algorithm to push flow backward to optimize."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 2)
        network.add_edge(0, 2, 2)
        network.add_edge(1, 2, 1)
        network.add_edge(1, 3, 1)
        network.add_edge(2, 3, 2)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 3)

    def test_14_cyclic_network(self):
        """Test a network with a direct cycle to confirm it avoids infinite loops."""
        # Arrange.
        network = DinicsNetwork(4)
        network.add_edge(0, 1, 10)
        network.add_edge(1, 2, 10)
        network.add_edge(2, 1, 10)
        network.add_edge(2, 3, 10)

        # Act.
        result = network.compute_maximum_flow(0, 3)

        # Assert.
        self.assertEqual(result, 10)

    def test_15_residual_capacity_calculation(self):
        """Verify the calculation of residual capacity as flow updates."""
        # Arrange.
        edge = FlowEdge(0, 1, 15)

        # Act & Assert.
        self.assertEqual(edge.residual_capacity(), 15)

        # Act.
        edge.current_flow = 5

        # Assert.
        self.assertEqual(edge.residual_capacity(), 10)

    def test_16_reverse_edge_linking(self):
        """Verify that add_edge correctly creates and binds forward/backward objects."""
        # Arrange.
        network = DinicsNetwork(2)

        # Act.
        network.add_edge(0, 1, 5)

        # Assert.
        forward = network.adjacency_list[0][0]
        backward = network.adjacency_list[1][0]
        self.assertEqual(forward.reverse_edge, backward)
        self.assertEqual(backward.reverse_edge, forward)
        self.assertEqual(forward.capacity, 5)
        self.assertEqual(backward.capacity, 0)

    def test_17_large_capacities(self):
        """Test that large scales don't cause performance regressions or precision issues."""
        # Arrange.
        network = DinicsNetwork(3)
        network.add_edge(0, 1, 10 ** 9)
        network.add_edge(1, 2, 10 ** 9)

        # Act.
        result = network.compute_maximum_flow(0, 2)

        # Assert.
        self.assertEqual(result, 10 ** 9)

    def test_18_asymmetric_bipartite_matching_style(self):
        """Test a wider multi-layered structure mimicking a matching scenario."""
        # Arrange.
        network = DinicsNetwork(6)
        network.add_edge(0, 1, 1)
        network.add_edge(0, 2, 2)
        network.add_edge(1, 3, 1)
        network.add_edge(2, 3, 1)
        network.add_edge(2, 4, 1)
        network.add_edge(3, 5, 1)
        network.add_edge(4, 5, 2)

        # Act.
        result = network.compute_maximum_flow(0, 5)

        # Assert.
        self.assertEqual(result, 2)

    def test_19_multiple_consecutive_runs_on_same_instance(self):
        """Ensure running compute_maximum_flow multiple times returns consistent results."""
        # Arrange.
        network = DinicsNetwork(3)
        network.add_edge(0, 1, 10)
        network.add_edge(1, 2, 5)

        # Act.
        first_run = network.compute_maximum_flow(0, 2)
        second_run = network.compute_maximum_flow(0, 2)

        # Assert.
        self.assertEqual(first_run, 5)
        self.assertEqual(second_run, 0)

    def test_20_long_linear_chain(self):
        """Test a long series chain to evaluate Level Graph building across deep paths."""
        # Arrange.
        size = 20
        network = DinicsNetwork(size)
        for index in range(size - 1):
            network.add_edge(index, index + 1, 100 - index)

        # Act.
        result = network.compute_maximum_flow(0, size - 1)

        # Assert.
        self.assertEqual(result, 82)

if __name__ == "__main__":
    unittest.main()