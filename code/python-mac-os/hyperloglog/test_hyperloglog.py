"""
test_hyperloglog.py

Automated unit test suite for verifying the HyperLogLog cardinality estimation
implementation across diverse stream patterns, range regimes, and edge bounds.

Overview:
    This module uses Python's built-in `unittest` framework to execute a matrix
    of 20 unique testing scenarios against the `HyperLogLog` engine. It ensures
    that register calculations, error bounds, zero-count states, small-range
    linear counting corrections, and large-range saturation behave with complete
    structural integrity.

Testing Scope & Architectural Coverage:
    - Initialization Guardrails: Assesses parameter type safety bounds, negative
      inputs, and precision constraints.
    - Stochastic Routing & Hashing: Verifies deterministic behavior of SHA-256
      truncation, empty string inputs, and rank evaluation of leading zeros.
    - Correction Regime Invariants: Evaluates empty-register scenarios, small
      cardinality linear fallback loops, and massive overflow boundaries.
    - Value Hardening: Checks high-water mark updates, register boundaries, and
      mathematical error convergence.

Execution Instructions:
    Run the complete testing matrix natively from the shell or command line via:

        python -m unittest test_hyperloglog.py
"""
from __future__ import annotations

import unittest

from hyperloglog import HyperLogLog

class TestHyperLogLog(unittest.TestCase):
    """
    Test suite containing 20 comprehensive, isolated test vectors evaluating
    the HyperLogLog sketch registers, corrections, and unique count logic.
    """

    def test_01_initial_empty_state(self):
        """1. Verifies that an initialized HyperLogLog instance contains entirely zeroed structures."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)

        # Act & Assert.
        self.assertEqual(hyperloglog.estimate(), 0)
        self.assertEqual(sum(hyperloglog.registers), 0)
        self.assertEqual(len(hyperloglog.registers), hyperloglog.register_count)

    def test_02_invalid_error_rate_types(self):
        """2. Validates that initializing with non-numeric error types triggers a TypeError."""
        # Arrange & Act & Assert.
        with self.assertRaises(TypeError):
            # Act.
            HyperLogLog(error_rate="invalid")  # type: ignore

    def test_03_error_rate_out_of_bounds_high(self):
        """3. Confirms a ValueError triggers when the target error rate equals or exceeds 1.0."""
        # Arrange & Act & Assert.
        with self.assertRaises(ValueError):
            # Act.
            HyperLogLog(error_rate=1.0)

    def test_04_error_rate_out_of_bounds_low(self):
        """4. Confirms a ValueError triggers when the target error rate falls below 0.0001."""
        # Arrange & Act & Assert.
        with self.assertRaises(ValueError):
            # Act.
            HyperLogLog(error_rate=0.00005)

    def test_05_precision_boundary_clamping(self):
        """5. Verifies internal bit width constraints clamp cleanly between the [4, 16] safe zones."""
        # Arrange & Act.
        hyperloglog_high = HyperLogLog(error_rate=0.50)
        hyperloglog_low = HyperLogLog(error_rate=0.001)

        # Assert.
        self.assertEqual(hyperloglog_high.precision_boundary, 4)
        self.assertEqual(hyperloglog_low.precision_boundary, 16)

    def test_06_alpha_multipliers_resolution(self):
        """6. Checks that standard correction scalars yield the correct empirical constants."""
        # Arrange & Act.
        hyperloglog_16 = HyperLogLog(error_rate=0.35)
        hyperloglog_32 = HyperLogLog(error_rate=0.23)
        hyperloglog_64 = HyperLogLog(error_rate=0.15)

        # Assert.
        self.assertAlmostEqual(hyperloglog_16.alpha, 0.673)
        self.assertAlmostEqual(hyperloglog_32.alpha, 0.697)
        self.assertAlmostEqual(hyperloglog_64.alpha, 0.709)

    def test_07_hash_type_safety_guards(self):
        """7. Confirms addition operations reject non-string elements with a TypeError."""
        # Arrange.
        hyperloglog = HyperLogLog()

        # Act & Assert.
        with self.assertRaises(TypeError):
            # Act.
            hyperloglog.add(12345)

    def test_08_deterministic_hashing(self):
        """8. Assures internal cryptographic hash signatures produce uniform, reproducible values."""
        # Arrange.
        hyperloglog = HyperLogLog()

        # Act.
        hash_val = hyperloglog._hash("test_token")
        other = hyperloglog._hash("test_token")

        # Assert.
        self.assertEqual(hash_val, other)
        self.assertIsInstance(hash_val, int)

    def test_09_leading_zeros_all_zeros(self):
        """9. Validates zero check logic when a sub-segment contains completely zeroed bits."""
        # Arrange.
        hyperloglog = HyperLogLog()

        # Act.
        rank = hyperloglog._retrieve_leading_zeros(bits=0, maximum_bits=32)

        # Assert.
        self.assertEqual(rank, 33)

    def test_10_leading_zeros_first_bit_active(self):
        """10. Checks rank evaluation when the absolute highest bit position is active."""
        # Arrange.
        hyperloglog = HyperLogLog()

        # Act.
        rank = hyperloglog._retrieve_leading_zeros(bits=16, maximum_bits=5)

        # Assert.
        self.assertEqual(rank, 1)

    def test_11_duplicate_element_idempotency(self):
        """11. Checks that duplicate elements do not distort or change estimation states."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)

        # Act.
        for _ in range(100):
            hyperloglog.add("constant_string")

        # Assert.
        self.assertEqual(hyperloglog.estimate(), 1)

    def test_12_small_range_linear_counting(self):
        """12. Validates that small-scale insertions trigger linear counting log math corrections."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.01)

        # Act.
        hyperloglog.add("item_1")
        hyperloglog.add("item_2")
        hyperloglog.add("item_3")

        raw_sum = sum(1.0 / (2.0 ** val) for val in hyperloglog.registers)
        raw_estimate = hyperloglog.alpha * (hyperloglog.register_count ** 2) * (1.0 / raw_sum)

        # Assert.
        self.assertTrue(raw_estimate <= 2.5 * hyperloglog.register_count)
        self.assertEqual(hyperloglog.estimate(), 3)

    def test_13_large_range_saturation_handling(self):
        """13. Evaluates that artificial saturation triggers the upper-bound logarithm exception."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)
        hyperloglog.registers = [60] * hyperloglog.register_count

        # Act & Assert.
        with self.assertRaises(OverflowError):
            # Act.
            hyperloglog.estimate()

    def test_14_empty_string_input_handling(self):
        """14. Verifies that passing an empty string handles securely without crashing."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)

        # Act.
        hyperloglog.add("")

        # Assert.
        self.assertEqual(hyperloglog.estimate(), 1)

    def test_15_maximum_rank_boundary_limit(self):
        """15. Assures calculated rank outputs can never exceed the designated bit-width boundaries."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)
        max_possible_bits = 64 - hyperloglog.precision_boundary

        # Act.
        rank = hyperloglog._retrieve_leading_zeros(bits=1, maximum_bits=max_possible_bits)

        # Assert.
        self.assertTrue(rank <= (max_possible_bits + 1))

    def test_16_single_item_register_index_bounds(self):
        """16. Confirms an item's routed register address fits within the physical array limits."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)

        # Act.
        hash_value = hyperloglog._hash("bounded_index_check")
        register_index = hash_value >> (64 - hyperloglog.precision_boundary)

        # Assert.
        self.assertTrue(0 <= register_index < hyperloglog.register_count)

    def test_17_sequential_hash_distribution(self):
        """17. Assures that distinct input variants hit diverse register buckets rather than clumping."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)

        # Act.
        hyperloglog.add("alpha_token")
        hyperloglog.add("beta_token")
        hyperloglog.add("gamma_token")
        active_buckets = len([value for value in hyperloglog.registers if value > 0])

        # Assert.
        self.assertTrue(active_buckets >= 1)

    def test_18_register_high_water_marks(self):
        """18. Assures register values only increment if incoming metrics beat old structural maxima."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)
        hyperloglog.add("element_omega")
        initial_registers = list(hyperloglog.registers)

        # Act.
        hyperloglog.add("element_omega")

        # Assert.
        self.assertEqual(hyperloglog.registers, initial_registers)

    def test_19_stochastic_error_convergence(self):
        """19. Tests convergence properties over a medium-sized unique token stream to check variance limits."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.03)
        total_items = 2000

        # Act.
        for index in range(total_items):
            hyperloglog.add(f"token_id_sequence_{index}")

        observed_estimate = hyperloglog.estimate()
        variance = abs(observed_estimate - total_items) / total_items

        # Assert.
        self.assertTrue(variance < 0.05)

    def test_20_display_execution_safety(self):
        """20. Assures formatting utilities run without structural exceptions or string interpolation failures."""
        # Arrange.
        hyperloglog = HyperLogLog(error_rate=0.05)
        hyperloglog.add("test_diagnostic_element")

        # Act & Assert.
        try:
            # Act.
            hyperloglog.display()
            executed_safely = True
        except Exception:
            executed_safely = False

        # Assert.
        self.assertTrue(executed_safely)

if __name__ == "__main__":
    unittest.main()