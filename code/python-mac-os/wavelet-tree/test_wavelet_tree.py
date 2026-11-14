"""
test_wavelet_tree.py

Automated unit test suite for verifying the custom class-based Wavelet Tree
data structure implementation.

Overview:
    This module uses Python's built-in `unittest` framework to execute a matrix
    of 20 unique testing scenarios against the `Node` and `WaveletTree` classes.
    It ensures algorithmic precision across standard structural setups, edge-case
    alphabets, empty structures, and varying rank boundary frontiers.

Testing Scope & Architectural Coverage:
    - Initialization & Topology: Validates empty sequences, single-element structures,
      homogeneous sequences, and negative or zero value configurations.
    - Standard Metrics: Verifies accurate bit-vector transitions, midpoint calculations,
      and local rank tracking tables.
    - Query Boundaries: Evaluates precise start, middle, and end index limitations,
      out-of-bounds symbols, and zero-index edge states.

Execution Instructions:
    Run the complete testing matrix natively from the shell or command line via:

        python -m unittest test_wavelet_tree.py
"""

import unittest

from wavelet_tree import WaveletTree, Node

class TestWaveletTree(unittest.TestCase):
    """
    Test suite containing 20 isolated test scenarios evaluating alphabet partitions,
    prefix sum counts, and recursive rank lookup paths.
    """

    def test_01_empty_sequence_initialization(self):
        """1. Verifies that initializing with an empty data sequence handles attributes safely."""
        # Arrange.
        sequence = []

        # Act.
        tree = WaveletTree(sequence)
        rank_result = tree.rank(0, 5)

        # Assert.
        self.assertIsNone(tree.root)
        self.assertEqual(rank_result, 0)

    def test_02_single_element_sequence(self):
        """2. Validates execution and queries on a sequence containing exactly one item."""
        # Arrange.
        sequence = [7]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(1, 7)
        r2 = tree.rank(1, 5)
        r3 = tree.rank(0, 7)

        # Assert.
        self.assertEqual(r1, 1)
        self.assertEqual(r2, 0)
        self.assertEqual(r3, 0)

    def test_03_homogeneous_sequence_all_same(self):
        """3. Tests rank resolution on a sequence where all elements are completely identical."""
        # Arrange.
        sequence = [4, 4, 4, 4, 4]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(5, 4)
        r2 = tree.rank(3, 4)
        r3 = tree.rank(5, 1)

        # Assert.
        self.assertEqual(r1, 5)
        self.assertEqual(r2, 3)
        self.assertEqual(r3, 0)

    def test_04_standard_sequence_middle_query(self):
        """4. Evaluates a standard lookup scenario tracking a value within a middle sequence slice."""
        # Arrange.
        sequence = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3]

        # Act.
        tree = WaveletTree(sequence)
        computed_rank = tree.rank(5, 1)

        # Assert.
        self.assertEqual(computed_rank, 2)

    def test_05_standard_sequence_full_range_query(self):
        """5. Checks total occurrences of a symbol across the entire initialized sequence bound."""
        # Arrange.
        sequence = [3, 1, 4, 1, 5, 9, 2, 6, 5, 3]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(10, 3)
        r2 = tree.rank(10, 5)
        r3 = tree.rank(10, 9)

        # Assert.
        self.assertEqual(r1, 2)
        self.assertEqual(r2, 2)
        self.assertEqual(r3, 1)

    def test_06_query_index_at_zero_boundary(self):
        """6. Assures that a query at index 0 immediately returns 0 for any valid symbol."""
        # Arrange.
        sequence = [1, 2, 3, 4, 5]

        # Act.
        tree = WaveletTree(sequence)
        computed_rank = tree.rank(0, 3)

        # Assert.
        self.assertEqual(computed_rank, 0)

    def test_07_query_value_below_minimum_alphabet(self):
        """7. Verifies that symbol queries strictly below the structural minimum gracefully fail fast."""
        # Arrange.
        sequence = [10, 20, 30, 40]

        # Act.
        tree = WaveletTree(sequence)
        computed_rank = tree.rank(3, 5)

        # Assert.
        self.assertEqual(computed_rank, 0)

    def test_08_query_value_above_maximum_alphabet(self):
        """8. Verifies that symbol queries strictly above the structural maximum gracefully fail fast."""
        # Arrange.
        sequence = [10, 20, 30, 40]

        # Act.
        tree = WaveletTree(sequence)
        computed_rank = tree.rank(4, 50)

        # Assert.
        self.assertEqual(computed_rank, 0)

    def test_09_query_value_unseen_but_within_range_limits(self):
        """9. Targets an unrepresented interior value that resides within the min/max alphabet span."""
        # Arrange.
        sequence = [1, 2, 4, 5]

        # Act.
        tree = WaveletTree(sequence)
        computed_rank = tree.rank(4, 3)

        # Assert.
        self.assertEqual(computed_rank, 0)

    def test_10_negative_values_alphabet_handling(self):
        """10. Validates tree partition mechanics when data contains negative integer ranges."""
        # Arrange.
        sequence = [-5, -1, -5, -2, -10, 0]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(4, -5)
        r2 = tree.rank(6, 0)

        # Assert.
        self.assertEqual(r1, 2)
        self.assertEqual(r2, 1)

    def test_11_sequential_increasing_array(self):
        """11. Evaluates tracking behavior when processing a perfectly ordered increasing sequence."""
        # Arrange.
        sequence = [1, 2, 3, 4, 5, 6, 7]

        tree = WaveletTree(sequence)

        # Act & Assert.
        for index, value in enumerate(sequence):
            self.assertEqual(tree.rank(index + 1, value), 1)

    def test_12_sequential_decreasing_array(self):
        """12. Evaluates tracking behavior when processing a perfectly ordered decreasing sequence."""
        # Arrange.
        sequence = [7, 6, 5, 4, 3, 2, 1]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(7, 4)
        r2 = tree.rank(3, 7)
        r3 = tree.rank(3, 4)

        # Assert.
        self.assertEqual(r1, 1)
        self.assertEqual(r2, 1)
        self.assertEqual(r3, 0)

    def test_13_alternating_binary_sequence(self):
        """13. Tests alphabet splits utilizing a compact two-symbol alternating sequence pattern."""
        # Arrange.
        sequence = [0, 1, 0, 1, 0, 1]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(6, 0)
        r2 = tree.rank(5, 1)

        # Assert.
        self.assertEqual(r1, 3)
        self.assertEqual(r2, 2)

    def test_14_large_alphabet_span_discrepancy(self):
        """14. Checks structural stability when two values form a vast interior gap."""
        # Arrange.
        sequence = [1, 1000000]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(2, 1000000)
        r2 = tree.rank(1, 1000000)

        # Assert.
        self.assertEqual(r1, 1)
        self.assertEqual(r2, 0)

    def test_15_query_index_exceeding_sequence_length(self):
        """15. Confirms execution behavior when the search index stretches past the valid collection boundary."""
        # Arrange.
        sequence = [5, 6, 7]
        tree = WaveletTree(sequence)

        # Act.
        computed_rank = tree.rank(5, 6)

        # Assert.
        self.assertEqual(computed_rank, 1)

    def test_16_internal_node_precomputed_counts_structure(self):
        """16. Directly checks the correctness of internal left-count tracking tables at the root node."""
        # Arrange.
        sequence = [3, 1, 4, 2]
        expected_counts = [0, 0, 1, 1, 2]

        # Act.
        tree = WaveletTree(sequence)

        # Assert.
        self.assertEqual(tree.root.left_counts, expected_counts)

    def test_17_all_zeros_sequence(self):
        """17. Assures structural precision when dealing entirely with absolute zero parameters."""
        # Arrange.
        sequence = [0, 0, 0, 0]

        # Act.
        tree = WaveletTree(sequence)
        computed_rank = tree.rank(3, 0)

        # Assert.
        self.assertEqual(computed_rank, 3)

    def test_18_reproducibility_invariance(self):
        """18. Guarantees that separate, identical rank queries yield the exact same response value."""
        # Arrange.
        sequence = [9, 3, 5, 9, 3, 5]

        tree = WaveletTree(sequence)

        # Act.
        first_query = tree.rank(4, 9)
        second_query = tree.rank(4, 9)

        # Assert.
        self.assertEqual(first_query, second_query)
        self.assertEqual(first_query, 2)

    def test_19_root_node_property_isolation(self):
        """19. Checks that root min and max properties conform with initialization limits."""
        # Arrange.
        sequence = [12, 99, 45, 67]

        # Act.
        tree = WaveletTree(sequence)

        # Assert.
        self.assertEqual(tree.minimum_value, 12)
        self.assertEqual(tree.maximum_value, 99)

    def test_20_dense_duplicate_distribution(self):
        """20. Tests complicated structural pathways with frequent repeating interspersed symbols."""
        # Arrange.
        sequence = [1, 2, 1, 2, 1, 2, 1, 2, 3, 3, 3]

        # Act.
        tree = WaveletTree(sequence)

        r1 = tree.rank(8, 1)
        r2 = tree.rank(8, 2)
        r3 = tree.rank(11, 3)

        # Assert.
        self.assertEqual(r1, 4)
        self.assertEqual(r2, 4)
        self.assertEqual(r3, 3)

if __name__ == "__main__":
    unittest.main()