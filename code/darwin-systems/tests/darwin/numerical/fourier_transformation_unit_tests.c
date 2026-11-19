/**
 * @file        :   fourier_transformation_unit_tests.c
 * @author      :   Braiden Gole
 * @brief       :   Comprehensive unit test suite for the Fast Fourier Transform (FFT) module.
 * @details     :   Utilizes the Unity Test Framework to validate the correctness, error handling, 
 *                  round-trip signal reconstruction, energy conservation (Parseval's theorem), 
 *                  and mathematical accuracy of the Cooley-Tukey radix-2 FFT and inverse FFT.
 * Compile:
 * clang -DUNITY_INCLUDE_DOUBLE ./tests/darwin/numerical/fourier_transformation_unit_tests.c ./src/darwin/numerical/fast-fourier-transformation/fourier_transformation.c ./src/unity.c -Iinclude -lm -o ./bin/fourier_transformation_unit_test
 * @version     :   1.0
 * @date        :   2026-11-10
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>

#include "unity.h"
#include "darwin/numerical/fast-fourier-transformation/fourier_transformation.h"

/**
 * @brief Sets up resources before each test.
 */
void setUp(void)
{
    ;;
}

/**
 * @brief Cleans up resources after each test.
 */
void tearDown(void)
{
    ;;
}

/**
 * @brief Ensures passing a null signal pointer returns false safely.
 */
void test_fft_null_signal_pointer_returns_false(void)
{
    // Arrange.
    ComplexNumber* null_signal = NULL;

    size_t count = 8;

    // Act.
    bool transformation_status = fast_fourier_transformation(null_signal, count, false);

    // Assert.
    TEST_ASSERT_FALSE(transformation_status);
}

/**
 * @brief Validates that a sample count of zero returns false.
 */
void test_fft_zero_sample_count_returns_false(void)
{
    // Arrange.
    ComplexNumber signal_buffer[1] = {{0.0, 0.0}};

    size_t count = 0;

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_FALSE(transformation_status);
}

/**
 * @brief Ensures non-power-of-two sample counts return false.
 */
void test_fft_non_power_of_two_count_three_returns_false(void)
{
    // Arrange.
    ComplexNumber signal_buffer[3] = {{1.0, 0.0}, {2.0, 0.0}, {3.0, 0.0}};

    size_t count = 3;

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_FALSE(transformation_status);
}

/**
 * @brief Verifies that sample counts not equal to a power of two like six return false.
 */
void test_fft_non_power_of_two_count_six_returns_false(void)
{
    // Arrange.
    ComplexNumber signal_buffer[6] = {
        {1.0, 0.0}, {2.0, 0.0}, {3.0, 0.0}, 
        {4.0, 0.0}, {5.0, 0.0}, {6.0, 0.0}
    };

    size_t count = 6;

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_FALSE(transformation_status);
}

/**
 * @brief Ensures sample counts exceeding safety limits return false.
 */
void test_fft_exceeding_max_sample_limit_returns_false(void)
{
    // Arrange.
    ComplexNumber* dummy_signal = (ComplexNumber*)malloc(sizeof(ComplexNumber));

    size_t excessive_count = ((size_t)1 << 30) + 2;

    // Act.
    bool transformation_status = fast_fourier_transformation(dummy_signal, excessive_count, false);

    // Assert.
    TEST_ASSERT_FALSE(transformation_status);

    free(dummy_signal);
}

/**
 * @brief Validates handling of a single sample input vector.
 */
void test_fft_single_sample_vector_succeeds(void)
{
    // Arrange.
    ComplexNumber signal_buffer[1] = {{5.0, 3.0}};

    size_t count = 1;

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(transformation_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 5.0, signal_buffer[0].real_component);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 3.0, signal_buffer[0].imaginary_component);
}

/**
 * @brief Tests transformation on a minimal valid power-of-two size of two.
 */
void test_fft_minimum_power_of_two_size_two_succeeds(void)
{
    // Arrange.
    ComplexNumber signal_buffer[2] = {{1.0, 0.0}, {1.0, 0.0}};

    size_t count = 2;

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(transformation_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 2.0, signal_buffer[0].real_component);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[0].imaginary_component);
}

/**
 * @brief Verifies constant DC signal concentrates all energy into bin zero.
 */
void test_fft_constant_dc_signal_concentrates_energy_at_bin_zero(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {2.0, 0.0}, {2.0, 0.0}, {2.0, 0.0}, {2.0, 0.0}
    };

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(transformation_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 8.0, signal_buffer[0].real_component);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[0].imaginary_component);

    for (size_t index = 1; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[index].imaginary_component);
    }
}

/**
 * @brief Validates impulse signal produces a flat frequency spectrum.
 */
void test_fft_impulse_signal_produces_flat_frequency_spectrum(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {1.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}
    };

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(transformation_status);

    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, signal_buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[index].imaginary_component);
    }
}

/**
 * @brief Verifies forward and inverse FFT round-trip reconstruction on a sine wave.
 */
void test_fft_round_trip_reconstruction_preserves_original_sine_wave(void)
{
    // Arrange.
    size_t count = 8;

    ComplexNumber original_signal[8];
    ComplexNumber working_buffer[8];

    for (size_t index = 0; index < count; index++)
    {
        double value = sin(2.0 * M_PI * (double)index / (double)count);

        original_signal[index].real_component = value;
        original_signal[index].imaginary_component = 0.0;
        working_buffer[index] = original_signal[index];
    }

    // Act.
    bool forward_success = fast_fourier_transformation(working_buffer, count, false);
    bool inverse_success = fast_fourier_transformation(working_buffer, count, true);

    // Assert.
    TEST_ASSERT_TRUE(forward_success);
    TEST_ASSERT_TRUE(inverse_success);

    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, original_signal[index].real_component, working_buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, original_signal[index].imaginary_component, working_buffer[index].imaginary_component);
    }
}

/**
 * @brief Validates round-trip reconstruction for a pulse signal.
 */
void test_fft_round_trip_reconstruction_pulse_signal(void)
{
    // Arrange.
    size_t count = 8;

    ComplexNumber working_buffer[8] = {
        {0.0, 0.0}, {1.0, 0.0}, {0.0, 0.0}, {0.0, 0.0},
        {0.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}
    };

    ComplexNumber expected_signal[8];

    for (size_t index = 0; index < count; index++)
    {
        expected_signal[index] = working_buffer[index];
    }

    // Act.
    fast_fourier_transformation(working_buffer, count, false);
    fast_fourier_transformation(working_buffer, count, true);

    // Assert.
    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, expected_signal[index].real_component, working_buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, expected_signal[index].imaginary_component, working_buffer[index].imaginary_component);
    }
}

/**
 * @brief Ensures an all-zero signal buffer remains zero after transformation.
 */
void test_fft_zero_signal_buffer_remains_zero(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {0.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}
    };

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(transformation_status);

    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, signal_buffer[index].imaginary_component);
    }
}

/**
 * @brief Tests transformation on a pure imaginary input buffer.
 */
void test_fft_pure_imaginary_input_succeeds(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {0.0, 1.0}, {0.0, 2.0}, {0.0, 3.0}, {0.0, 4.0}
    };

    // Act.
    bool transformation_status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(transformation_status);
}

/**
 * @brief Validates 16-sample size power of two forward and inverse transform.
 */
void test_fft_sixteen_sample_vector_round_trip(void)
{
    // Arrange.
    size_t count = 16;

    ComplexNumber signal_buffer[16];
    ComplexNumber backup_buffer[16];

    for (size_t index = 0; index < count; index++)
    {
        double value = cos(2.0 * M_PI * (double)index / (double)count);

        signal_buffer[index].real_component = value;
        signal_buffer[index].imaginary_component = 0.5 * value;

        backup_buffer[index] = signal_buffer[index];
    }

    // Act.
    fast_fourier_transformation(signal_buffer, count, false);
    fast_fourier_transformation(signal_buffer, count, true);

    // Assert.
    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, backup_buffer[index].real_component, signal_buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, backup_buffer[index].imaginary_component, signal_buffer[index].imaginary_component);
    }
}

/**
 * @brief Validates 32-sample size power of two transform success.
 */
void test_fft_thirty_two_sample_vector_succeeds(void)
{
    // Arrange.
    size_t count = 32;

    ComplexNumber* signal_buffer = (ComplexNumber*)malloc(count * sizeof(ComplexNumber));

    for (size_t index = 0; index < count; index++)
    {
        signal_buffer[index].real_component = (double)(index % 5);
        signal_buffer[index].imaginary_component = 0.0;
    }

    // Act.
    bool status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(status);

    free(signal_buffer);
}

/**
 * @brief Verifies Parseval's theorem (energy conservation between time and frequency domains).
 */
void test_fft_parseval_theorem_energy_conservation(void)
{
    // Arrange.
    size_t count = 8;

    ComplexNumber signal_buffer[8];

    double time_domain_energy = 0.0;

    for (size_t index = 0; index < count; index++)
    {
        signal_buffer[index].real_component = (double)(index + 1);
        signal_buffer[index].imaginary_component = 0.0;

        time_domain_energy += (signal_buffer[index].real_component * signal_buffer[index].real_component);
    }

    // Act.
    fast_fourier_transformation(signal_buffer, count, false);

    double frequency_domain_energy = 0.0;

    for (size_t index = 0; index < count; index++)
    {
        double real = signal_buffer[index].real_component;
        double imaginary = signal_buffer[index].imaginary_component;

        frequency_domain_energy += (real * real + imaginary * imaginary);
    }

    frequency_domain_energy /= (double)count;

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, time_domain_energy, frequency_domain_energy);
}

/**
 * @brief Tests alternating sign signal producing expected Nyquist frequency response.
 */
void test_fft_alternating_sign_signal_nyquist_frequency(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {1.0, 0.0}, {-1.0, 0.0}, {1.0, 0.0}, {-1.0, 0.0}
    };

    // Act.
    bool status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 4.0, signal_buffer[2].real_component);
}

/**
 * @brief Verifies in-place mutation of the input buffer.
 */
void test_fft_modifies_buffer_in_place(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {1.0, 0.0}, {2.0, 0.0}, {3.0, 0.0}, {4.0, 0.0}
    };

    ComplexNumber snapshot[4];

    for (size_t index = 0; index < count; index++)
    {
        snapshot[index] = signal_buffer[index];
    }

    // Act.
    fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    bool content_changed = false;

    for (size_t index = 0; index < count; index++)
    {
        if (signal_buffer[index].real_component != snapshot[index].real_component)
        {
            content_changed = true;

            break;
        }
    }

    TEST_ASSERT_TRUE(content_changed);
}

/**
 * @brief Validates inverse transform scaling factor divides correctly by N.
 */
void test_fft_inverse_scaling_factor_correctness(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {12.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}, {0.0, 0.0}
    };

    // Act.
    fast_fourier_transformation(signal_buffer, count, true);

    // Assert.
    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, 3.0, signal_buffer[index].real_component);
    }
}

/**
 * @brief Tests sequential forward transformations without memory corruption.
 */
void test_fft_multiple_sequential_transforms_stability(void)
{
    // Arrange.
    size_t count = 4;

    ComplexNumber signal_buffer[4] = {
        {1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}, {4.0, 4.0}
    };

    // Act & Assert.
    TEST_ASSERT_TRUE(fast_fourier_transformation(signal_buffer, count, false));
    TEST_ASSERT_TRUE(fast_fourier_transformation(signal_buffer, count, false));
}

/**
 * @brief Validates round-trip accuracy with random phase components.
 */
void test_fft_random_phase_complex_signal_round_trip(void)
{
    // Arrange.
    size_t count = 8;

    ComplexNumber buffer[8] = {
        {1.5, -2.3}, {0.5, 1.1}, {-1.2, 3.4}, {2.1, 0.0},
        {-0.5, -0.5}, {3.3, 2.2}, {-2.2, -1.1}, {1.0, 1.0}
    };

    ComplexNumber reference[8];

    for (size_t index = 0; index < count; index++)
    {
        reference[index] = buffer[index];
    }

    // Act.
    fast_fourier_transformation(buffer, count, false);
    fast_fourier_transformation(buffer, count, true);

    // Assert.
    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-4, reference[index].real_component, buffer[index].real_component);
        TEST_ASSERT_DOUBLE_WITHIN(1e-4, reference[index].imaginary_component, buffer[index].imaginary_component);
    }
}

/**
 * @brief Ensures boundary bit-reversal works correctly for 8-element arrays.
 */
void test_fft_bit_reversal_permutation_integrity(void)
{
    // Arrange.
    size_t count = 8;

    ComplexNumber signal_buffer[8];

    for (size_t index = 0; index < count; index++)
    {
        signal_buffer[index].real_component = (double)index;
        signal_buffer[index].imaginary_component = 0.0;
    }

    // Act.
    bool forward = fast_fourier_transformation(signal_buffer, count, false);
    bool inverse = fast_fourier_transformation(signal_buffer, count, true);

    // Assert.
    TEST_ASSERT_TRUE(forward);
    TEST_ASSERT_TRUE(inverse);

    for (size_t index = 0; index < count; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-5, (double)index, signal_buffer[index].real_component);
    }
}

/**
 * @brief Validates failure handling for size 5 (non power of two).
 */
void test_fft_non_power_of_two_count_five_returns_false(void)
{
    // Arrange.
    ComplexNumber signal_buffer[5] = {
        {1.0, 0.0}, {2.0, 0.0}, {3.0, 0.0}, {4.0, 0.0}, {5.0, 0.0}
    };
    
    size_t count = 5;

    // Act.
    bool status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_FALSE(status);
}

/**
 * @brief Validates failure handling for size 10 (non power of two).
 */
void test_fft_non_power_of_two_count_ten_returns_false(void)
{
    // Arrange.
    ComplexNumber signal_buffer[10];

    size_t count = 10;

    // Act.
    bool status = fast_fourier_transformation(signal_buffer, count, false);

    // Assert.
    TEST_ASSERT_FALSE(status);
}

int main(void)
{
    // Initialize the session.
    UNITY_BEGIN();

    // Register the tests to run.
    RUN_TEST(test_fft_null_signal_pointer_returns_false);
    RUN_TEST(test_fft_zero_sample_count_returns_false);
    RUN_TEST(test_fft_non_power_of_two_count_three_returns_false);
    RUN_TEST(test_fft_non_power_of_two_count_six_returns_false);
    RUN_TEST(test_fft_exceeding_max_sample_limit_returns_false);
    RUN_TEST(test_fft_single_sample_vector_succeeds);
    RUN_TEST(test_fft_minimum_power_of_two_size_two_succeeds);
    RUN_TEST(test_fft_constant_dc_signal_concentrates_energy_at_bin_zero);
    RUN_TEST(test_fft_impulse_signal_produces_flat_frequency_spectrum);
    RUN_TEST(test_fft_round_trip_reconstruction_preserves_original_sine_wave);
    RUN_TEST(test_fft_round_trip_reconstruction_pulse_signal);
    RUN_TEST(test_fft_zero_signal_buffer_remains_zero);
    RUN_TEST(test_fft_pure_imaginary_input_succeeds);
    RUN_TEST(test_fft_sixteen_sample_vector_round_trip);
    RUN_TEST(test_fft_thirty_two_sample_vector_succeeds);
    RUN_TEST(test_fft_parseval_theorem_energy_conservation);
    RUN_TEST(test_fft_alternating_sign_signal_nyquist_frequency);
    RUN_TEST(test_fft_modifies_buffer_in_place);
    RUN_TEST(test_fft_inverse_scaling_factor_correctness);
    RUN_TEST(test_fft_multiple_sequential_transforms_stability);
    RUN_TEST(test_fft_random_phase_complex_signal_round_trip);
    RUN_TEST(test_fft_bit_reversal_permutation_integrity);
    RUN_TEST(test_fft_non_power_of_two_count_five_returns_false);
    RUN_TEST(test_fft_non_power_of_two_count_ten_returns_false);
    RUN_TEST(test_fft_round_trip_reconstruction_preserves_original_sine_wave);
    
    // Destroy the test session.
    return UNITY_END();
}
