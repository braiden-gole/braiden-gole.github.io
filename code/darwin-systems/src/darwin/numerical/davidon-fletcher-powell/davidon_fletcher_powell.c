/**
 * @file        :   davidon_fletcher_powell.c
 * @author      :   Braiden Gole
 * @brief       :   Implementation of the Davidon-Fletcher-Powell (DFP) Quasi-Newton optimization algorithm.
 * @details     :   This source file provides an unconstrained nonlinear optimization routine using
 *                  the DFP formula to iteratively approximate the inverse Hessian matrix. It features 
 *                  defensive boundary checks, dynamic memory allocation safety, Euclidean norm and 
 *                  dot product utilities, an Armijo backtracking line search, and Rosenbrock test functions.
 * @version     :   1.0
 * @date        :   2026-11-11
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>

// Include header files.
#include "darwin/numerical/davidon-fletcher-powell/davidon_fletcher_powell.h"

/**
 * @brief       :   Calculates the Euclidean norm (length) of a numerical vector.
 * @details     :   Computes the square root of the sum of squared components across the specified vector space.
 * @param       :   vector: Pointer to the input vector array.
 * @param       :   space_dimension: The number of dimensions (elements) in the vector.
 * @return      :   double: The computed Euclidean norm, or 0.0 if inputs are invalid.
 */
double calculate_vector_normalization(const double* vector, size_t space_dimension)
{
    // Validate input pointers and dimensions to prevent segmentation faults.
    if (vector == NULL || space_dimension == 0)
    {
        return 0.0;
    }

    double sum_squared_components = 0.0;

    // Accumulate the squared value of each vector component.
    for (size_t index = 0; index < space_dimension; index++)
    {
        sum_squared_components += vector[index] * vector[index];
    }

    // Return the square root of the accumulated sum representing the Euclidean length.
    return sqrt(sum_squared_components);
}

/**
 * @brief       :   Calculates the dot product of two vectors of equal dimension.
 * @details     :   Iterates through each component, multiplying corresponding vector elements and accumulating the sum.
 * @param       :   vector: Pointer to the first vector array.
 * @param       :   operand: Pointer to the second vector array (operand).
 * @param       :   space_dimension: The dimension size shared by both vectors.
 * @return      :   double: The resulting scalar dot product value, or 0.0 on validation failure.
 */
double calculate_vector_dot_product(const double* vector, const double* operand, size_t space_dimension)
{
    // Validate both input vectors and space dimensions.
    if (vector == NULL || operand == NULL || space_dimension == 0)
    {
        return 0.0;
    }

    double accumulator_dot_product = 0.0;

    // Multiply corresponding elements across both vectors and accumulate the sum.
    for (size_t index = 0; index < space_dimension; index++)
    {
        accumulator_dot_product += vector[index] * operand[index];
    }

    return accumulator_dot_product;
}

/**
 * @brief       :   Evaluates Rosenbrock's test objective function at a given coordinate vector.
 * @details     :   Computes the standard 2D Rosenbrock banana function: f(x,y) = (1 - x)^2 + 100(y - x^2)^2.
 * @param       :   evaluation_vector: Pointer to the coordinate array.
 * @param       :   space_dimension: The dimension of the evaluation vector.
 * @return      :   double: The scalar function value at the evaluation point.
 */
double evaluate_rosenbrock_objective(const double* evaluation_vector, size_t space_dimension)
{
    // Ensure the vector is valid and contains at least 2 dimensions for Rosenbrock.
    if (evaluation_vector == NULL || space_dimension < 2)
    {
        return 0.0;
    }

    (void)space_dimension;

    // Extract x and y coordinates from the evaluation vector.
    double x_coordinate = evaluation_vector[0];
    double y_coordinate = evaluation_vector[1];
    
    // Compute and return the Rosenbrock function value.
    return (1.0 - x_coordinate) * (1.0 - x_coordinate) + 100.0 * (y_coordinate - x_coordinate * x_coordinate) * (y_coordinate - x_coordinate * x_coordinate);
}

/**
 * @brief       :   Evaluates Rosenbrock's analytical gradient vector at a given coordinate vector.
 * @details     :   Computes the partial derivatives of the Rosenbrock function with respect to x and y.
 * @param       :   evaluation_vector: Pointer to the current coordinate array.
 * @param       :   output_gradient_vector: Pointer to the pre-allocated array where gradient components are stored.
 */
void evaluate_rosenbrock_gradient(const double* evaluation_vector, double* output_gradient_vector)
{
    // Validate pointers to safeguard against memory errors.
    if (evaluation_vector == NULL || output_gradient_vector == NULL)
    {
        return;
    }

    // Extract coordinates.
    double x_coordinate = evaluation_vector[0];
    double y_coordinate = evaluation_vector[1];

    // Compute the partial derivative with respect to x.
    output_gradient_vector[0] = -2.0 * (1.0 - x_coordinate) - 400.0 * x_coordinate * (y_coordinate - x_coordinate * x_coordinate);
    
    // Compute the partial derivative with respect to y.
    output_gradient_vector[1] = 200.0 * (y_coordinate - x_coordinate * x_coordinate);
}

/**
 * @brief       :   Performs an Armijo backtracking line search to determine an optimal step size.
 * @details     :   Scales the search direction step size iteratively to satisfy the sufficient decrease condition,
 *                  ensuring numerical stability and steady descent toward a local minimum.
 * @param       :   objective: Function pointer to the objective cost function.
 * @param       :   gradient: Function pointer to the gradient evaluation function.
 * @param       :   current_position: Pointer to the current coordinate vector.
 * @param       :   search_direction: Pointer to the search direction vector.
 * @param       :   gradient_vector: Pointer to the current gradient vector.
 * @param       :   space_dimension: The dimension of the problem space.
 * @param       :   step_size: The initial step size scalar to test.
 * @param       :   contraction_factor: Multiplicative factor (< 1.0) to reduce step size upon failure.
 * @param       :   sufficient_decrease_constant: Constant governing the Armijo condition threshold.
 * @return      :   double: The accepted step size scalar, or 0.0 if line search fails.
 */
double perform_armijo_line_search(ObjectiveFunction objective, GradientFunction gradient,
    const double* current_position,
    const double* search_direction,
    const double* gradient_vector,
    size_t space_dimension,
    double step_size,
    double contraction_factor,
    double sufficient_decrease_constant)
{
    // Validate all function pointers, input arrays, dimensions, and hyperparameters.
    if (objective == NULL || gradient == NULL || current_position == NULL || search_direction == NULL ||
        gradient_vector == NULL || space_dimension == 0 || step_size <= 0.0 || contraction_factor <= 0.0 || contraction_factor >= 1.0)
    {
        return 0.0;
    }

    double step_size_scalar = step_size;
    double current_function = objective(current_position, space_dimension);
    double directional_derivative = calculate_vector_dot_product(gradient_vector, search_direction, space_dimension);

    // Abort line search if the computed search direction is not a valid descent direction.
    if (directional_derivative >= 0.0)
    {
        return 0.0;
    }

    // Allocate memory for the temporary trial position vector.
    double* trial_vector = (double*)malloc(sizeof(double) * space_dimension);

    if (trial_vector == NULL)
    {
        return 0.0;
    }

    const size_t maximum_iterations = 100;

    // Iteratively test step sizes using backtracking.
    for (size_t index = 0; index < maximum_iterations; index++)
    {
        // Compute the trial position coordinates using current step size scaling.
        for (size_t vector_index = 0; vector_index < space_dimension; vector_index++)
        {
            trial_vector[vector_index] = current_position[vector_index] + step_size_scalar * search_direction[vector_index];
        }

        // Evaluate objective function at the trial position.
        double trial_function = objective(trial_vector, space_dimension);

        // Check if trial step satisfies the Armijo sufficient decrease condition.
        if (trial_function <= current_function + sufficient_decrease_constant * step_size_scalar * directional_derivative)
        {
            free(trial_vector);

            return step_size_scalar;
        }

        // Reduce step size by multiplying with the contraction factor.
        step_size_scalar *= contraction_factor;
    }

    // Clean up trial vector memory if search fails to converge within iteration bounds.
    free(trial_vector);
    
    return step_size_scalar;
}

/**
 * @brief       :   Executes the Davidon-Fletcher-Powell (DFP) Quasi-Newton optimization algorithm.
 * @details     :   Iteratively minimizes an objective function by approximating the inverse Hessian matrix
 *                  using gradient and position displacements. Includes robust memory safety checks, convergence
 *                  tests based on gradient norm and parameter displacement tolerances, and DFP matrix updates.
 * @param       :   objective_function: Function pointer to the objective evaluation method.
 * @param       :   gradient_function: Function pointer to the gradient evaluation method.
 * @param       :   solution_position_vector: Pointer to initial guess coordinates (updated in-place to final solution).
 * @param       :   space_dimension: The dimension size of the optimization vector space.
 * @param       :   maximum_optimization_iterations: The upper limit of iteration loops allowed.
 * @param       :   gradient_tolerance: Threshold norm for gradient convergence stopping condition.
 * @param       :   parameter_tolerance: Threshold norm for position displacement convergence stopping condition.
 * @return      :   bool: true if optimization successfully converged within tolerances, false otherwise.
 */
bool optimize_quasi_newton_dfp(
    ObjectiveFunction objective_function,
    GradientFunction gradient_function,
    double* solution_position_vector,
    size_t space_dimension,
    size_t maximum_optimization_iterations,
    double gradient_tolerance,
    double parameter_tolerance)
{
    // Defensive check for invalid pointers or zero dimensions.
    if (objective_function == NULL || gradient_function == NULL || solution_position_vector == NULL || space_dimension == 0)
    {
        return false;
    }

    // Allocate memory blocks for vector quantities and matrix approximations.
    double* current_gradient_vector = (double*)malloc(sizeof(double) * space_dimension);
    double* next_gradient_vector = (double*)malloc(sizeof(double) * space_dimension);
    double* search_direction_vector = (double*)malloc(sizeof(double) * space_dimension);
    double* position_displacement_vector = (double*)malloc(sizeof(double) * space_dimension);
    double* gradient_displacement_vector = (double*)malloc(sizeof(double) * space_dimension);
    double* hessian_gradient_product_vector = (double*)malloc(sizeof(double) * space_dimension);
    double* inverse_hessian_matrix = (double*)malloc(sizeof(double) * space_dimension * space_dimension);

    // Defensive check ensuring all dynamic memory allocations succeeded.
    if (current_gradient_vector == NULL || next_gradient_vector == NULL ||
        search_direction_vector == NULL || position_displacement_vector == NULL ||
        gradient_displacement_vector == NULL || hessian_gradient_product_vector == NULL ||
        inverse_hessian_matrix == NULL)
    {
        // Free any allocated memory blocks before returning failure.
        free(current_gradient_vector);
        free(next_gradient_vector);
        free(search_direction_vector);
        free(position_displacement_vector);
        free(gradient_displacement_vector);
        free(hessian_gradient_product_vector);
        free(inverse_hessian_matrix);
        
        return false;
    }

    // Initialize inverse Hessian approximation matrix to the Identity Matrix.
    for (size_t row_index = 0; row_index < space_dimension; row_index++)
    {
        for (size_t column_index = 0; column_index < space_dimension; column_index++)
        {
            inverse_hessian_matrix[row_index * space_dimension + column_index] = (row_index == column_index) ? 1.0 : 0.0;
        }
    }

    // Compute initial gradient vector at starting position.
    gradient_function(solution_position_vector, current_gradient_vector);

    bool optimization_converged = false;

    // Main optimization iteration loop.
    for (size_t iteration_counter = 0; iteration_counter < maximum_optimization_iterations; iteration_counter++)
    {
        // Calculate the normalization (Euclidean norm) of the current gradient vector.
        double current_gradient_normalization = calculate_vector_normalization(current_gradient_vector, space_dimension);

        // Check if gradient norm has fallen below tolerance threshold (successful convergence).
        if (current_gradient_normalization < gradient_tolerance)
        {
            optimization_converged = true;
            break;
        }

        // Compute search direction vector: p_k = - H_k * g_k using matrix-vector multiplication.
        for (size_t row_index = 0; row_index < space_dimension; row_index++)
        {
            search_direction_vector[row_index] = 0.0;

            for (size_t column_index = 0; column_index < space_dimension; column_index++)
            {
                search_direction_vector[row_index] -= inverse_hessian_matrix[row_index * space_dimension + column_index] * current_gradient_vector[column_index];
            }
        }

        // Determine step size scaling factor using Armijo line search.
        double step_size_scalar = perform_armijo_line_search(
            objective_function,
            gradient_function,
            solution_position_vector,
            search_direction_vector,
            current_gradient_vector,
            space_dimension,
            1.0,
            0.5,
            1e-4
        );

        // Terminate if line search fails to find a valid step size.
        if (step_size_scalar <= 0.0)
        {
            break;
        }

        // Update position vector in-place and compute position displacement vector (s_k).
        for (size_t vector_index = 0; vector_index < space_dimension; vector_index++)
        {
            position_displacement_vector[vector_index] = step_size_scalar * search_direction_vector[vector_index];

            solution_position_vector[vector_index] += position_displacement_vector[vector_index];
        }

        // Check if parameter step size meets convergence criteria.
        if (calculate_vector_normalization(position_displacement_vector, space_dimension) < parameter_tolerance)
        {
            optimization_converged = true;
            break;
        }

        // Compute new gradient vector at updated position and calculate gradient displacement vector (y_k).
        gradient_function(solution_position_vector, next_gradient_vector);

        for (size_t vector_index = 0; vector_index < space_dimension; vector_index++)
        {
            gradient_displacement_vector[vector_index] = next_gradient_vector[vector_index] - current_gradient_vector[vector_index];
        }

        // Calculate inner product of position displacement and gradient displacement (s_k^T * y_k).
        double position_gradient_inner_product = calculate_vector_dot_product(position_displacement_vector, gradient_displacement_vector, space_dimension);

        // Compute H_k * y_k product term.
        for (size_t row_index = 0; row_index < space_dimension; row_index++)
        {
            hessian_gradient_product_vector[row_index] = 0.0;

            for (size_t column_index = 0; column_index < space_dimension; column_index++)
            {
                hessian_gradient_product_vector[row_index] += inverse_hessian_matrix[row_index * space_dimension + column_index] * gradient_displacement_vector[column_index];
            }
        }

        // Calculate quadratic form denominator: y_k^T * H_k * y_k.
        double denominator_quadratic_form = calculate_vector_dot_product(gradient_displacement_vector, hessian_gradient_product_vector, space_dimension);

        // Verify curvature conditions and numerical stability before updating inverse Hessian matrix.
        if (position_gradient_inner_product > 1e-10 && denominator_quadratic_form > 1e-10)
        {
            double* updated_inverse_hessian_matrix = (double*)malloc(sizeof(double) * space_dimension * space_dimension);

            if (updated_inverse_hessian_matrix != NULL)
            {
                // Apply Davidon-Fletcher-Powell (DFP) update formula element by element.
                for (size_t row_index = 0; row_index < space_dimension; row_index++)
                {
                    for (size_t column_index = 0; column_index < space_dimension; column_index++)
                    {
                        double baseline_hessian_element = inverse_hessian_matrix[row_index * space_dimension + column_index];
                        double position_outer_product_element = (position_displacement_vector[row_index] * position_displacement_vector[column_index]) / position_gradient_inner_product;
                        double curvature_product_element = (hessian_gradient_product_vector[row_index] * hessian_gradient_product_vector[column_index]) / denominator_quadratic_form;

                        updated_inverse_hessian_matrix[row_index * space_dimension + column_index] = baseline_hessian_element + position_outer_product_element - curvature_product_element;
                    }
                }

                // Copy updated matrix back into primary storage container.
                for (size_t matrix_element_index = 0; matrix_element_index < space_dimension * space_dimension; matrix_element_index++)
                {
                    inverse_hessian_matrix[matrix_element_index] = updated_inverse_hessian_matrix[matrix_element_index];
                }

                // Free temporary updated matrix block.
                free(updated_inverse_hessian_matrix);
            }
        }

        // Advance current gradient vector state to prepare for the next iteration.
        for (size_t vector_index = 0; vector_index < space_dimension; vector_index++)
        {
            current_gradient_vector[vector_index] = next_gradient_vector[vector_index];
        }
    }

    // Release all dynamically allocated memory blocks to prevent memory leaks.
    free(current_gradient_vector);
    free(next_gradient_vector);
    free(search_direction_vector);
    free(position_displacement_vector);
    free(gradient_displacement_vector);
    free(hessian_gradient_product_vector);
    free(inverse_hessian_matrix);

    return optimization_converged;
}