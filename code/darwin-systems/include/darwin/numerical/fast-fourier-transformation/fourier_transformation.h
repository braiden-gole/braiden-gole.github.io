/**
 * @file        :   fourier_transformation.h
 * @author      :   Braiden Gole
 * @brief       :   Header file for the Fast Fourier Transform (FFT) algorithm.
 * @details     :   Defines the ComplexNumber structure and function prototype 
 *                  for computing the Fast Fourier Transform and its inverse.
 * @version     :   1.0
 * @date        :   2026-11-10
 */
#ifndef FOURIER_TRANSFORMATION_H
#define FOURIER_TRANSFORMATION_H

#include <stdbool.h>
#include <stddef.h>

/**
 * @brief Represents a complex number with explicit real and imaginary components.
 */
typedef struct
{
    double real_component;
    double imaginary_component;
} ComplexNumber;

/**
 * @brief           :   Computes the in-place Fast Fourier Transform (FFT) or Inverse FFT.
 * @param signal    :   Pointer to the array of complex numbers representing the signal buffer.
 * @param count     :   The total sample count, which must be a positive power of two.
 * @param flag      :   Boolean flag indicating direction (false for forward FFT, true for inverse FFT).
 * @return          :   True if the transformation succeeds; otherwise false.
 */
bool fast_fourier_transformation(ComplexNumber* signal, size_t count, bool flag);

#endif