/**
 * @file        :   davidon_fletcher_powell.h
 * @author      :   Braiden Gole
 * @brief       :   Header file for the Davidon-Fletcher-Powell (DFP) Quasi-Newton optimization algorithm.
 * @details     :   This header defines the function pointer signatures for objective and gradient 
 *                  evaluations, along with the public function prototype for executing the DFP 
 *                  unconstrained nonlinear optimization engine.
 * @version     :   1.0
 * @date        :   2026-11-11
 */
#ifndef DAVIDON_FLETCHER_POWELL_H
#define DAVIDON_FLETCHER_POWELL_H

#include <stddef.h>
#include <stdbool.h>

/**
 * @brief       :   Function pointer signature for objective cost function evaluation.
 * @param       :   evaluation_vector: Pointer to the coordinate array being evaluated.
 * @param       :   space_dimension: The dimension size of the evaluation vector.
 * @return      :   double: The computed scalar cost value.
 */
typedef double (*ObjectiveFunction)(const double* evaluation_vector, size_t space_dimension);

/**
 * @brief       :   Function pointer signature for analytical gradient evaluation.
 * @param       :   evaluation_vector: Pointer to the current coordinate array.
 * @param       :   output_gradient_vector: Pointer to the pre-allocated array where gradient components are stored.
 */
typedef void (*GradientFunction)(const double* evaluation_vector, double* output_gradient_vector);

// Function prototypes.
double calculate_vector_normalization(const double* vector, size_t space_dimension);
double evaluate_rosenbrock_objective(const double* evaluation_vector, size_t space_dimension);
void evaluate_rosenbrock_gradient(const double* evaluation_vector, double* output_gradient_vector);
double calculate_vector_dot_product(const double* vector, const double* operand, size_t space_dimension);

double perform_armijo_line_search(ObjectiveFunction objective, GradientFunction gradient,
    const double* current_position,
    const double* search_direction,
    const double* gradient_vector,
    size_t space_dimension,
    double step_size,
    double contraction_factor,
    double sufficient_decrease_constant);

bool optimize_quasi_newton_dfp(
    ObjectiveFunction objective_function,
    GradientFunction gradient_function,
    double* solution_position_vector,
    size_t space_dimension,
    size_t maximum_optimization_iterations,
    double gradient_tolerance,
    double parameter_tolerance
);

#endif
