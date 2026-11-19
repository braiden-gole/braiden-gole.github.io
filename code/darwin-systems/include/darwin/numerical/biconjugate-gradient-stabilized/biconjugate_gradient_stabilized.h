/**
 * @file        :   biconjugate_gradient_stabilized.h
 * @author      :   Braiden Gole
 * @brief       :   Header definition for the Biconjugate Gradient Stabilized (BiCGSTAB) linear solver algorithm.
 * @details     :   This header file provides type definitions, status error codes, callback function pointers, 
 *                  and public function prototypes required to execute the iterative BiCGSTAB linear solver.
 * @version     :   1.0
 * @date        :   2026-11-11
 */
#ifndef BICONJUGATE_GRADIENT_STABILIZED_H
#define BICONJUGATE_GRADIENT_STABILIZED_H

/**
 * @brief       :   Status and error code enumeration for the BiCGSTAB linear solver.
 */
typedef enum
{
    SUCCESS = 0,
    ERROR_NULL_POINTER = -1,
    ERROR_INVALID_DIMENSION = -2,
    ERROR_ALLOCATION_FAILED = -3,
    ERROR_BREAKDOWN_RHO_ZERO = -4,
    ERROR_BREAKDOWN_ALPHA_ZERO = -5,
    ERROR_BREAKDOWN_OMEGA_ZERO = -6,
    ERROR_MAXIMUM_ITERATIONS_EXCEEDED = -1
} Status;

/**
 * @brief       :   Function pointer type for matrix-vector multiplication.
 * @details     :   Computes the matrix-vector product vector_a * input_vector.
 * @param       :   number_of_dimensions: Dimension size of the vectors.
 * @param       :   vector: Pointer to the input vector.
 * @param       :   product_vector: Pointer to store the resulting product vector.
 * @param       :   context: User-defined context pointer passed to the callback.
 */
typedef void (*MultiplyFunction)(int number_of_dimensions, const double* vector, double* product_vector, void* context);

// Function prototypes.
double compute_vector_dot_product(int number_of_dimensions, const double* vector, const double* operand);
double compute_vector_normalization(int number_of_dimensions, const double* vector);

Status solve(MultiplyFunction callback, int number_of_dimensions, void* context, 
    const double* operand, const double* initial_guess, double* solution_vector, int maximum_iterations, double convergence_tolerance);

#endif
