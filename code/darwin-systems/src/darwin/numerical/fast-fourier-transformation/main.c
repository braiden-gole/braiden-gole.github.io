/**
 * @file        :   main.c
 * @author      :   Braiden Gole
 * @brief       :   Demonstration program for the Fast Fourier Transform (FFT) algorithm.
 * @details     :   Allocates a complex signal buffer, populates it with a sinusoidal wave, 
 *                  executes the forward FFT, prints the frequency spectrum, and performs 
 *                  an inverse FFT to verify signal reconstruction.
 * @version     :   1.0
 * @date        :   2026-11-10
 * 
 * Compile:
 * clang ./src/darwin/numerical/fast-fourier-transformation/main.c ./src/darwin/numerical/fast-fourier-transformation/fourier_transformation.c -I./include -lm -o ./bin/fourier_transformation_test
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>

// Include header files.
#include "darwin/numerical/fast-fourier-transformation/fourier_transformation.h"

/**
 * @brief           :   Entry point of the demonstration program.
 * @return          :   Returns 0 upon successful execution; otherwise 1 on failure.
 */
int main(void)
{
    // Define the total number of samples, which must be a positive power of two.
    size_t total_sample_count = 8;

    // Allocate memory for the complex signal buffer with defensive error handling.
    ComplexNumber* signal_buffer = (ComplexNumber*)malloc(total_sample_count * sizeof(ComplexNumber));

    if (signal_buffer == NULL)
    {
        fprintf(stderr, "%s", "Error: Memory allocation failed for the signal buffer.\n");

        return 1;
    }

    // Populate the signal buffer with a discrete sinusoidal test wave.
    for (size_t sample_index = 0; sample_index < total_sample_count; sample_index++)
    {
        signal_buffer[sample_index].real_component = sin(2.0 * M_PI * (double)sample_index / (double)total_sample_count);
        signal_buffer[sample_index].imaginary_component = 0.0;
    }

    // Print the original time-domain signal values.
    printf("%s", "Original Time-Domain Signal:\n");

    for (size_t sample_index = 0; sample_index < total_sample_count; sample_index++)
    {
        printf("\tSample %zu: Real = %.4f, Imaginary = %.4f\n", 
               sample_index, 
               signal_buffer[sample_index].real_component, 
               signal_buffer[sample_index].imaginary_component);
    }

    // Execute the forward Fast Fourier Transform (flag set to false).
    bool forward_transform_success = fast_fourier_transformation(signal_buffer, total_sample_count, false);

    if (!forward_transform_success)
    {
        fprintf(stderr, "%s", "Error: Fast Fourier Transform computation failed.\n");
        free(signal_buffer);

        return 1;
    }

    // Print the resulting frequency domain spectrum.
    printf("%s", "\nTransformed Frequency Spectrum (Forward FFT):\n");

    for (size_t frequency_index = 0; frequency_index < total_sample_count; frequency_index++)
    {
        printf("%c%s %zu: %s %.4lf %s %.4lf%c", '\t',
               "Bin",
               frequency_index,
               "Real =",
               signal_buffer[frequency_index].real_component,
               "Imaginary =",
               signal_buffer[frequency_index].imaginary_component,
               '\n');
    }

    // Execute the inverse Fast Fourier Transform to reconstruct the signal (flag set to true).
    bool inverse_transform_success = fast_fourier_transformation(signal_buffer, total_sample_count, true);

    if (!inverse_transform_success)
    {
        fprintf(stderr, "%s", "Error: Inverse Fast Fourier Transform computation failed.\n");
        free(signal_buffer);

        return 1;
    }

    // Print the reconstructed time-domain signal to verify accuracy.
    printf("%s", "\nReconstructed Signal (Inverse FFT):\n");

    for (size_t sample_index = 0; sample_index < total_sample_count; sample_index++)
    {
        printf("%c%s %zu: %s %.4lf %s %.4lf%c", '\t',
               "Sample",
               sample_index,
               "Real =",
               signal_buffer[sample_index].real_component,
               "Imaginary =",
               signal_buffer[sample_index].imaginary_component,
               '\n');
    }

    // Deallocate dynamically allocated memory to prevent leaks.
    free(signal_buffer);

    return 0;
}