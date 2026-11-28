"""
test_van_emde_boas.py

Automated unit test suite for verifying the recursive Van Emde Boas (vanEmdeBoas)
tree data structure implementation.

Overview:
    This module uses Python's built-in `unittest` framework to execute a matrix
    of 32 unique testing scenarios against the `VanEmdeBoas` set structure. It
    ensures mathematical and structural precision across structural cluster splitting,
    ordered predecessor/successor lookups, set transformations, and edge erasures.

Testing Scope & Architectural Coverage:
    - Initialization & Bounds: validates valid power-of-two universe allocations,
      asymmetric bit-slicing configurations, and base cases ($M = 2$).
    - Insertion Primitives & Sorting: Tracks proper min/max routing, nested sub-cluster
      element propagation, and automated summary synchronization.
    - Order Statistics: Audits exact boundaries for successor and predecessor query
      pipelines across internal and cross-cluster boundaries.
    - Erasure Mechanics & State Dissolution: Tests structural replacement logic for
      removing minimums, maximums, and cluster-emptying elements down to empty states.
    - Native Python Interfaces: Validates subscript syntax (`__getitem__`) and standard
      membership evaluations (`__contains__`).

Execution Instructions:
    Run the complete testing matrix natively from the shell or command line via:

        python -m unittest test_van_emde_boas.py
"""
import unittest

from van_emde_boas import VanEmdeBoas

class TestVanEmdeBoas(unittest.TestCase):
    """
    Test suite containing 32 isolated test vectors evaluating cluster routing,
    order statistics, empty-tree short-circuits, and erasures inside a vanEmdeBoas tree.
    """

    def test_01_base_case_initialization_state(self):
        """1. Verifies that an M=2 base-case tree initializes with null clusters and summary."""
        # Arrange & Act.
        vanEmdeBoas = VanEmdeBoas(2)

        # Assert.
        self.assertEqual(vanEmdeBoas.universe_size, 2)
        self.assertIsNone(vanEmdeBoas.summary)
        self.assertIsNone(vanEmdeBoas.cluster)

    def test_02_base_case_square_root_fallbacks(self):
        """2. Checks that an M=2 base-case tree safely sets its lower and upper square roots to 1."""
        # Arrange & Act.
        vanEmdeBoas = VanEmdeBoas(2)

        # Assert.
        self.assertEqual(vanEmdeBoas.lower_square_root, 1)
        self.assertEqual(vanEmdeBoas.upper_square_root, 1)

    def test_03_standard_universe_asymmetric_slicing(self):
        """3. Validates correct bit-slicing square root splits for non-perfect square capacities (M=8)."""
        # Arrange & Act.
        vanEmdeBoas = VanEmdeBoas(8)

        # Assert.
        self.assertEqual(vanEmdeBoas.lower_square_root, 2)
        self.assertEqual(vanEmdeBoas.upper_square_root, 4)

    def test_04_perfect_square_universe_slicing(self):
        """4. Validates correct symmetric bit-slicing square root splits for perfect square capacities (M=16)."""
        # Arrange & Act.
        vanEmdeBoas = VanEmdeBoas(16)

        # Assert.
        self.assertEqual(vanEmdeBoas.lower_square_root, 4)
        self.assertEqual(vanEmdeBoas.upper_square_root, 4)

    def test_05_nested_structural_allocation_depth(self):
        """5. Verifies nested sub-cluster generation matches upper square root capacities."""
        # Arrange & Act.
        vanEmdeBoas = VanEmdeBoas(16)

        # Assert.
        self.assertEqual(len(vanEmdeBoas.cluster), vanEmdeBoas.upper_square_root)
        self.assertEqual(vanEmdeBoas.summary.universe_size, vanEmdeBoas.upper_square_root)

    def test_06_empty_tree_boundary_defaults(self):
        """6. Confirms that an unmutated tree defaults its tracking boundaries strictly to None."""
        # Arrange & Act.
        vanEmdeBoas = VanEmdeBoas(16)

        # Assert.
        self.assertIsNone(vanEmdeBoas.minimum)
        self.assertIsNone(vanEmdeBoas.maximum)

    def test_07_single_element_insertion_boundaries(self):
        """7. Verifies that the first inserted key initializes both minimum and maximum boundaries."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        # Act.
        vanEmdeBoas.insert(5)

        # Assert.
        self.assertEqual(vanEmdeBoas.minimum, 5)
        self.assertEqual(vanEmdeBoas.maximum, 5)

    def test_08_two_ordered_elements_insertion(self):
        """8. Confirms correct min/max splitting when adding two perfectly sorted items."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        # Act.
        vanEmdeBoas.insert(4)
        vanEmdeBoas.insert(12)

        # Assert.
        self.assertEqual(vanEmdeBoas.minimum, 4)
        self.assertEqual(vanEmdeBoas.maximum, 12)

    def test_09_two_reverse_elements_insertion_swap(self):
        """9. Validates the structural value swapping mechanics when a smaller item follows a larger one."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        # Act.
        vanEmdeBoas.insert(12)
        vanEmdeBoas.insert(4)

        # Assert.
        self.assertEqual(vanEmdeBoas.minimum, 4)
        self.assertEqual(vanEmdeBoas.maximum, 12)

    def test_10_nested_cluster_propagation_on_insert(self):
        """10. Ensures that high/low coordinate splitting accurately routes items to sub-clusters."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)
        vanEmdeBoas.insert(2)

        # Act.
        vanEmdeBoas.insert(6)
        cluster_index = vanEmdeBoas._high(6)
        index = vanEmdeBoas._low(6)

        # Assert.
        self.assertEqual(vanEmdeBoas.cluster[cluster_index].minimum, index)

    def test_11_summary_synchronization_on_new_cluster(self):
        """11. Checks that summary structures accurately log cluster occupation changes."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)
        vanEmdeBoas.insert(1)

        # Act.
        vanEmdeBoas.insert(9)

        # Assert.
        self.assertTrue(vanEmdeBoas.summary.__contains__(2))

    def test_12_contains_short_circuit_on_min_max(self):
        """12. Assures membership checks short-circuit positively when evaluating top boundaries."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(3)
        vanEmdeBoas.insert(14)

        # Act & Assert.
        self.assertIn(3, vanEmdeBoas)
        self.assertIn(14, vanEmdeBoas)

    def test_13_contains_negative_evaluation(self):
        """13. Assures membership lookups evaluate to False for missing keys."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)
        vanEmdeBoas.insert(5)

        # Act & Assert.
        self.assertNotIn(8, vanEmdeBoas)

    def test_14_subscript_retrieval_success(self):
        """14. Verifies standard dict-like bracket lookup syntax over existing items."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)
        vanEmdeBoas.insert(7)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas[7], 7)

    def test_15_subscript_missing_key_exception(self):
        """15. Validates that subscripting an absent key throws a KeyError."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        # Act & Assert.
        with self.assertRaises(KeyError):
            # Act.
            _ = vanEmdeBoas[13]

    def test_16_base_case_successor_evaluation(self):
        """16. Tests successor lookup mechanics within an M=2 base layout."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(2)

        vanEmdeBoas.insert(0)
        vanEmdeBoas.insert(1)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.successor(0), 1)
        self.assertIsNone(vanEmdeBoas.successor(1))

    def test_17_base_case_predecessor_evaluation(self):
        """17. Tests predecessor lookup mechanics within an M=2 base layout."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(2)

        vanEmdeBoas.insert(0)
        vanEmdeBoas.insert(1)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.predecessor(1), 0)
        self.assertIsNone(vanEmdeBoas.predecessor(0))

    def test_18_successor_short_circuit_via_global_minimum(self):
        """18. Checks if lookups below the global minimum quickly return the minimum."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(4)
        vanEmdeBoas.insert(8)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.successor(2), 4)

    def test_19_successor_internal_same_cluster(self):
        """19. Validates searching for successors located within the identical sub-cluster."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(1)
        vanEmdeBoas.insert(4)
        vanEmdeBoas.insert(6)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.successor(4), 6)

    def test_20_successor_cross_cluster_traversal(self):
        """20. Evaluates successor lookup routing when skipping across adjacent cluster gaps."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(0)
        vanEmdeBoas.insert(5)
        vanEmdeBoas.insert(10)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.successor(5), 10)

    def test_21_predecessor_short_circuit_via_global_maximum(self):
        """21. Checks if lookups above the global maximum quickly return the maximum."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(2)
        vanEmdeBoas.insert(11)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.predecessor(15), 11)

    def test_22_predecessor_internal_same_cluster(self):
        """22. Validates searching for predecessors located within the identical sub-cluster."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(0)
        vanEmdeBoas.insert(9)
        vanEmdeBoas.insert(11)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.predecessor(11), 9)

    def test_23_predecessor_cross_cluster_traversal(self):
        """23. Evaluates predecessor lookup routing when skipping backward across cluster gaps."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(0)
        vanEmdeBoas.insert(4)
        vanEmdeBoas.insert(12)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.predecessor(12), 4)

    def test_24_erase_absent_element_no_op(self):
        """24. Confirms that attempting to erase a missing element exits silently with no mutations."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)
        vanEmdeBoas.insert(3)

        # Act.
        vanEmdeBoas.erase(9)

        # Assert.
        self.assertEqual(vanEmdeBoas.minimum, 3)

    def test_25_erase_single_element_depletion(self):
        """25. Assures erasing the last remaining element returns the tree to a clean empty state."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)
        vanEmdeBoas.insert(10)

        # Act.
        vanEmdeBoas.erase(10)

        # Assert.
        self.assertIsNone(vanEmdeBoas.minimum)
        self.assertIsNone(vanEmdeBoas.maximum)

    def test_26_erase_base_case_elements(self):
        """26. Validates mutation behaviors when clearing components from an M=2 node structure."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(2)

        vanEmdeBoas.insert(0)
        vanEmdeBoas.insert(1)

        # Act.
        vanEmdeBoas.erase(0)

        # Assert.
        self.assertEqual(vanEmdeBoas.minimum, 1)
        self.assertEqual(vanEmdeBoas.maximum, 1)

    def test_27_erase_global_minimum_promotion(self):
        """27. Tests deleting the global minimum, ensuring its successor gets cleanly promoted."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(2)
        vanEmdeBoas.insert(5)
        vanEmdeBoas.insert(9)

        # Act.
        vanEmdeBoas.erase(2)

        # Assert.
        self.assertEqual(vanEmdeBoas.minimum, 5)

    def test_28_erase_cluster_emptying_summary_cleanup(self):
        """28. Verifies summary pointers drop vacant cluster indexes when the final item inside them is deleted."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(1)
        vanEmdeBoas.insert(6)

        self.assertTrue(vanEmdeBoas.summary.__contains__(1))

        # Act.
        vanEmdeBoas.erase(6)

        # Assert.
        self.assertFalse(vanEmdeBoas.summary.__contains__(1))

    def test_29_erase_global_maximum_demotion(self):
        """29. Tests deleting the global maximum, ensuring the previous element becomes the new maximum."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        vanEmdeBoas.insert(1)
        vanEmdeBoas.insert(4)
        vanEmdeBoas.insert(13)

        # Act.
        vanEmdeBoas.erase(13)

        # Assert.
        self.assertEqual(vanEmdeBoas.maximum, 4)

    def test_30_to_list_empty_representation(self):
        """30. Validates that calling to_list over empty structures yields an empty list."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        # Act & Assert.
        self.assertEqual(vanEmdeBoas.to_list(), [])

    def test_31_bulk_sorted_list_generation(self):
        """31. Verifies that to_list systematically outputs full sorted contents via successor iterations."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        items = [14, 2, 7, 5, 11]

        # Act.
        for item in items:
            vanEmdeBoas.insert(item)

        # Assert.
        self.assertEqual(vanEmdeBoas.to_list(), [2, 5, 7, 11, 14])

    def test_32_comprehensive_interleaved_mutations(self):
        """32. Audits structure stability across long pipelines of mixed insertions, searches, and erasures."""
        # Arrange.
        vanEmdeBoas = VanEmdeBoas(16)

        items = [9, 2, 15, 0, 6, 4]

        # Act.
        for item in items:
            vanEmdeBoas.insert(item)

        vanEmdeBoas.erase(0)
        vanEmdeBoas.erase(9)

        # Assert.
        self.assertEqual(vanEmdeBoas.to_list(), [2, 4, 6, 15])
        self.assertEqual(vanEmdeBoas.minimum, 2)
        self.assertEqual(vanEmdeBoas.maximum, 15)

if __name__ == '__main__':
    unittest.main()