/**
 * @file        :   fourier_transformation.c
 * @author      :   Braiden Gole
 * @brief       :   Implementation of the Fast Fourier Transform (FFT) algorithm.
 * @details     :   Computes the Discrete Fourier Transform (DFT) or its inverse 
 * using the Cooley-Tukey iterative radix-2 algorithm in O(N log N) time complexity.
 * @version     :   1.0
 * @date        :   2026-11-10
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>

// Include header files.
#include "darwin/numerical/fast-fourier-transformation/fourier_transformation.h"

/**
 * @brief           :   Validates whether a given integer is a positive power of two.
 * @param candidate :   The integer candidate to test.
 * @return          :   True if the candidate is a positive power of two; otherwise false.
 */
static inline bool is_power_of_two(int candidate)
{
    if (candidate <= 0)
    {
        return false;
    }

    return (candidate & (candidate - 1)) == 0;
}

/**
 * @brief           :   Computes the in-place Fast Fourier Transform (FFT) or Inverse FFT.
 * @param signal    :   Pointer to the array of complex numbers representing the signal buffer.
 * @param count     :   The total sample count, which must be a positive power of two.
 * @param flag      :   Boolean flag indicating direction (false for forward FFT, true for inverse FFT).
 * @return          :   True if the transformation succeeds; otherwise false.
 * @details         :   Performs bit-reversal permutation followed by iterative Cooley-Tukey 
 *                      butterfly operations using dynamic trigonometric twiddle factors.
 */
bool fast_fourier_transformation(ComplexNumber* signal, size_t count, bool flag)
{
    // Validate input pointer to prevent segmentation faults.
    if (signal == NULL)
    {
        return false;
    }

    // Validate that the total sample count is a positive power of two.
    if (!is_power_of_two((int)count))
    {
        return false;
    }

    // Enforce safety limits against excessive memory allocation or overflow.
    if (count > ((size_t)1 << 30))
    {
        return false;
    }

    size_t reversal_index = 0;

    // Perform bit-reversal permutation to reorder elements for butterfly stages.
    for (size_t index = 0; index < count; index++)
    {
        if (index < reversal_index)
        {
            ComplexNumber temporary = signal[index];
            signal[index] = signal[reversal_index];
            signal[reversal_index] = temporary;
        }

        size_t mask = count >> 1;

        while ((reversal_index & mask) != 0)
        {
            reversal_index &= ~mask;
            mask >>= 1;
        }

        reversal_index |= mask;
    }

    // Execute the iterative Cooley-Tukey butterfly stages.
    for (size_t step = 2; step <= count; step <<= 1)
    {
        double angle_multiplier = -2.0 * M_PI / (double)step;

        // Invert angle multiplier if computing the inverse transform.
        if (flag)
        {
            angle_multiplier = -angle_multiplier;
        }

        ComplexNumber twiddle_factor_step =
        {
            cos(angle_multiplier),
            sin(angle_multiplier)
        };

        for (size_t block_index = 0; block_index < count; block_index += step)
        {
            ComplexNumber current = { 1.0, 0.0 };

            size_t half_step = step >> 1;

            for (size_t butterfly_index = 0; butterfly_index < half_step; butterfly_index++)
            {
                size_t lower_index = block_index + butterfly_index;
                size_t upper_index = lower_index + half_step;

                ComplexNumber upper_target = signal[upper_index];
                
                // Scale the upper branch by the current twiddle factor.
                ComplexNumber scaled_upper_branch =
                {
                    upper_target.real_component * current.real_component - upper_target.imaginary_component * current.imaginary_component,
                    upper_target.real_component * current.imaginary_component + upper_target.imaginary_component * current.real_component
                };

                ComplexNumber lower_target = signal[lower_index];

                // Compute butterfly additions and subtractions.
                signal[lower_index].real_component = lower_target.real_component + scaled_upper_branch.real_component;
                signal[lower_index].imaginary_component = lower_target.imaginary_component + scaled_upper_branch.imaginary_component;

                signal[upper_index].real_component = lower_target.real_component - scaled_upper_branch.real_component;
                signal[upper_index].imaginary_component = lower_target.imaginary_component - scaled_upper_branch.imaginary_component;

                // Update the current twiddle factor for the next butterfly iteration.
                double next_real_current = current.real_component * twiddle_factor_step.real_component - current.imaginary_component * twiddle_factor_step.imaginary_component;
                double next_imag_current = current.real_component * twiddle_factor_step.imaginary_component + current.imaginary_component * twiddle_factor_step.real_component;

                current.real_component = next_real_current;
                current.imaginary_component = next_imag_current;
            }
        }
    }

    // Apply normalization scaling factor if computing the inverse transform.
    if (flag)
    {
        double scaling_factor = (double)count;
        
        for (size_t index = 0; index < count; index++)
        {
            signal[index].real_component /= scaling_factor;
            signal[index].imaginary_component /= scaling_factor;
        }
    }

    return true;
}