"""
test_skip_list.py

A comprehensive, rigorous unit test suite establishing 20 distinct functional,
boundary, and structural invariant verification assertions for a probabilistic,
multi-layered SkipList class implementation.

Testing Domains Covered:
    1. Lifecycle & Structural Integrity: Empty initialization states, node link
       array sizing restrictions, and automatic high-level lane contraction sweeps.
    2. Mutator Mechanics: Layer-linked pointer insertions, duplicate element updates,
       non-existent payload updates, and multi-lane target erasing routines.
    3. Traversal Routing: Early-exit boundary protection against None queries, missing
       element lane drops, and exact match pointer hits.
    4. Sorting Invariance & Type Uniformity: Enforcing strict spatial sorting across
       mixed data representations including alphanumeric strings, floating points, and
       negative scale numeric dimensions.
    5. Monotonic & Stochastic Stress Testing: Edge configuration stability tracking
       under continuous ascending, descending, and chaotic high-volume fuzzed input arrays.

Execution Matrix:
    Run explicitly via command line interface using standard unittest execution flags:
        $ python3 -m unittest test_skip_list.py
"""
import unittest
import random

from skip_list import SkipList, Node

class TestSkipList(unittest.TestCase):
    def setUp(self):
        """Set up an empty skip list instance before executing each test case."""
        self.skip_list = SkipList()

    def test_01_initialization_empty_list(self):
        """Ensure an initialized skip list engine starts cleanly with empty structures."""
        # Arrange & Act.
        # (SkipList initialized in setUp)

        # Assert.
        self.assertEqual(self.skip_list.level, 0)
        self.assertEqual(len(self.skip_list.head.forward), self.skip_list.maximum_level + 1)
        self.assertIsNone(self.skip_list.head.forward[0])

    def test_02_single_key_insertion(self):
        """Verify inserting a baseline item establishes references correctly on Level 0."""
        # Arrange.
        # (SkipList initialized in setUp)

        # Act.
        self.skip_list.insert(10, "Apple")

        # Assert.
        self.assertEqual(self.skip_list.search(10), "Apple")

        base_first_node = self.skip_list.head.forward[0]

        self.assertIsNotNone(base_first_node)
        self.assertEqual(base_first_node.key, 10)
        self.assertEqual(base_first_node.value, "Apple")

    def test_03_duplicate_insertion_redirects_to_update(self):
        """Ensure inserting an identical key maps to update instead of corrupting paths."""
        # Arrange.
        self.skip_list.insert(42, "Original")

        # Act.
        self.skip_list.insert(42, "Mutated")

        # Assert.
        self.assertEqual(self.skip_list.search(42), "Mutated")
        self.assertIsNone(self.skip_list.head.forward[0].forward[0])

    def test_04_none_insertion_prohibition(self):
        """Ensure insertion engine rejects uninitialized or None type keys outright."""
        # Arrange.
        # (SkipList initialized in setUp)

        # Act & Assert.
        with self.assertRaises(ValueError):
            # Act.
            self.skip_list.insert(None, "Invalid")

    def test_05_update_existing_key(self):
        """Verify that updating a pre-existing key alters values and returns True."""
        # Arrange.
        self.skip_list.insert(15, "Banana")

        # Act.
        update_status = self.skip_list.update(15, "Cranberry")

        # Assert.
        self.assertTrue(update_status)
        self.assertEqual(self.skip_list.search(15), "Cranberry")

    def test_06_update_non_existent_key(self):
        """Verify that attempting to update an unallocated key safely returns False."""
        # Arrange.
        self.skip_list.insert(5, "Elderberry")

        # Act.
        update_status = self.skip_list.update(99, "Ghost")

        # Assert.
        self.assertFalse(update_status)
        self.assertIsNone(self.skip_list.search(99))

    def test_07_search_empty_list(self):
        """Verify that looking up elements in an unpopulated matrix safely returns None."""
        # Arrange.
        # (SkipList initialized in setUp)

        # Act & Assert.
        self.assertIsNone(self.skip_list.search(100))

    def test_08_search_non_existent_key(self):
        """Verify targeted lane traversal returns None when a value is absent from a populated list."""
        # Arrange.
        for key, value in [(10, "A"), (20, "B"), (30, "C")]:
            self.skip_list.insert(key, value)

        # Act & Assert.
        self.assertIsNone(self.skip_list.search(15))

    def test_09_search_valid_keys(self):
        """Verify search tracking successfully locates all inserted elements across data nodes."""
        # Arrange.
        items = [(25, "X"), (12, "Y"), (38, "Z"), (6, "W")]

        for key, value in items:
            self.skip_list.insert(key, value)

        # Act & Assert.
        for key, value in items:
            self.assertEqual(self.skip_list.search(key), value)

    def test_10_erase_valid_key(self):
        """Verify that erasing a valid node returns True and restructures pointers cleanly."""
        # Arrange.
        self.skip_list.insert(7, "Cherry")
        self.skip_list.insert(14, "Date")

        # Act.
        erase_status = self.skip_list.erase(7)

        # Assert.
        self.assertTrue(erase_status)
        self.assertIsNone(self.skip_list.search(7))
        self.assertEqual(self.skip_list.search(14), "Date")

    def test_11_erase_non_existent_key(self):
        """Verify that attempting to remove an unallocated key safely returns False."""
        # Arrange.
        self.skip_list.insert(50, "Fig")

        # Act.
        erase_status = self.skip_list.erase(999)

        # Assert.
        self.assertFalse(erase_status)

    def test_12_level_zero_ordering_invariance(self):
        """Verify that regardless of insertion order, the base layer always maintains strict sorted sorting."""
        # Arrange.
        unordered_keys = [50, 10, 40, 20, 30]

        for key in unordered_keys:
            self.skip_list.insert(key, f"Val_{key}")

        # Act.
        collected_keys = []
        current_node = self.skip_list.head.forward[0]

        while current_node is not None:
            collected_keys.append(current_node.key)
            current_node = current_node.forward[0]

        # Assert.
        self.assertEqual(collected_keys, [10, 20, 30, 40, 50])

    def test_13_global_level_trim_on_erasure(self):
        """Test that upper vertical lane layers contract when their highest elements are erased."""
        # Arrange.
        target_node = Node(100, "Century", 5)

        self.skip_list.level = 5

        for index in range(6):
            self.skip_list.head.forward[index] = target_node

        # Act.
        self.skip_list.erase(100)

        # Assert.
        self.assertEqual(self.skip_list.level, 0)

    def test_14_node_pointer_array_sizing(self):
        """Verify structural Node instantiation allocates a forward pointer list tracking level offsets."""
        # Arrange & Act.
        level_two_node = Node(key=5, value="Test", level=2)

        # Assert.
        self.assertEqual(len(level_two_node.forward), 3)
        self.assertIsNone(level_two_node.forward[0])

    def test_15_string_keys_lexicographical_sorting(self):
        """Verify that the matrix lanes accurately index alternative comparable scalar sequences like text strings."""
        # Arrange.
        text_items = [("delta", 4), ("alpha", 1), ("omega", 5), ("beta", 2)]

        for key, value in text_items:
            self.skip_list.insert(key, value)

        # Act.
        current_node = self.skip_list.head.forward[0]
        sorted_keys_output = []

        while current_node:
            sorted_keys_output.append(current_node.key)
            current_node = current_node.forward[0]

        # Assert.
        self.assertEqual(sorted_keys_output, ["alpha", "beta", "delta", "omega"])

    def test_16_negative_and_floating_point_keys(self):
        """Verify multi-level skip routers handle a combination of signed integers and decimal floats."""
        # Arrange.
        mixed_keys = [-15, 0, 3.1415, -2.5, 100.7]

        for key in mixed_keys:
            self.skip_list.insert(key, True)

        # Act & Assert.
        for key in mixed_keys:
            self.assertTrue(self.skip_list.search(key))

    def test_17_search_none_handling(self):
        """Ensure search retrieval traversals handle None parameter boundaries gracefully by returning None."""
        # Arrange.
        self.skip_list.insert(10, "Valid")

        # Act & Assert.
        self.assertIsNone(self.skip_list.search(None))

    def test_18_sequential_ascending_insertion_stability(self):
        """Verify inserting values in monotonically increasing order keeps all records accessible."""
        # Arrange.
        ascending_sequence = list(range(1, 51))

        for key in ascending_sequence:
            self.skip_list.insert(key, f"Num_{key}")

        # Act & Assert.
        for key in ascending_sequence:
            self.assertEqual(self.skip_list.search(key), f"Num_{key}")

    def test_19_sequential_descending_insertion_stability(self):
        """Verify inserting values in monotonically decreasing order maps lane entries correctly."""
        # Arrange.
        descending_sequence = list(range(50, 0, -1))

        for key in descending_sequence:
            self.skip_list.insert(key, f"Num_{key}")

        # Act & Assert.
        for key in descending_sequence:
            self.assertEqual(self.skip_list.search(key), f"Num_{key}")

    def test_20_massive_randomized_stride(self):
        """Execute a chaotic randomized insertion, lookup, and deletion fuzz sequence under load conditions."""
        # Arrange.
        random.seed(2026)

        large_dataset = list(set(random.randint(-10000, 10000) for _ in range(300)))

        # Act.
        for key in large_dataset:
            self.skip_list.insert(key, f"Value_{key}")

        for key in large_dataset:
            self.assertEqual(self.skip_list.search(key), f"Value_{key}")

        erasure_targets = large_dataset[:150]

        for key in erasure_targets:
            self.assertTrue(self.skip_list.erase(key))

        # Assert.
        for key in large_dataset:
            if key in erasure_targets:
                self.assertIsNone(self.skip_list.search(key))
            else:
                self.assertEqual(self.skip_list.search(key), f"Value_{key}")

if __name__ == "__main__":
    unittest.main()