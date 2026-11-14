/**
 * @file        :   biconjugate_gradient_stabilized.c
 * @author      :   Braiden Gole
 * @brief       :   Implementation of the Biconjugate Gradient Stabilized (BiCGSTAB) linear solver algorithm.
 * @details     :   This source file provides an iterative routine for solving non-symmetric systems of linear equations 
 *                  using the BiCGSTAB method. It features defensive bounds and null-pointer checks, dynamic memory 
 *                  allocation safety, Euclidean vector normalization, and dot product utilities.
 * @version     :   1.0
 * @date        :   2026-11-11
 */
#include <stdio.h>
#include <stdlib.h>
#include <math.h>

// Include header files.
#include "darwin/numerical/biconjugate-gradient-stabilized/biconjugate_gradient_stabilized.h"

/**
 * @brief       :   Calculates the dot product of two numerical vectors.
 * @details     :   Iterates through each component, multiplying corresponding vector elements and accumulating the sum.
 * @param       :   number_of_dimensions: The dimension size shared by both vectors.
 * @param       :   vector: Pointer to the first vector array.
 * @param       :   operand: Pointer to the second vector array (operand).
 * @return      :   double: The resulting scalar dot product value.
 */
double compute_vector_dot_product(int number_of_dimensions, const double* vector, const double* operand)
{
    // Defend against null pointer arguments or invalid dimensions before dereferencing.
    if (vector == NULL || operand == NULL || number_of_dimensions <= 0)
    {
        return 0.0;
    }

    // Initialize cumulative sum tracker to zero.
    double cumulative_sum = 0;

    // Multiply corresponding elements across both vectors and accumulate the sum.
    for (int index = 0; index < number_of_dimensions; index++)
    {
        cumulative_sum += vector[index] * operand[index];
    }

    return cumulative_sum;
}

/**
 * @brief       :   Calculates the Euclidean normalization (norm) of a numerical vector.
 * @details     :   Computes the square root of the dot product of the vector with itself.
 * @param       :   number_of_dimensions: The number of dimensions (elements) in the vector.
 * @param       :   vector: Pointer to the input vector array.
 * @return      :   double: The computed Euclidean normalization value.
 */
double compute_vector_normalization(int number_of_dimensions, const double* vector)
{
    // Defend against null pointer arguments or invalid dimensions before dereferencing.
    if (vector == NULL || number_of_dimensions <= 0)
    {
        return 0.0;
    }

    // Compute normalization by taking the square root of the dot product of the vector with itself.
    return sqrt(compute_vector_dot_product(number_of_dimensions, vector, vector));
}

/**
 * @brief       :   Solves a linear system Ax = b using the Biconjugate Gradient Stabilized (BiCGSTAB) method.
 * @details     :   Iteratively refines an initial guess solution vector using matrix-vector multiplication callbacks,
 *                  adaptive stabilization parameters, and robust memory safety checks.
 * @param       :   callback: Function pointer to perform the matrix-vector multiplication.
 * @param       :   number_of_dimensions: Number of rows and columns in the square matrix.
 * @param       :   context: User-defined context passed to the matrix multiplication callback.
 * @param       :   operand: The right-hand side vector b of size n.
 * @param       :   initial_guess: Initial guess vector x on input.
 * @param       :   solution_vector: Output buffer where the computed solution is stored.
 * @param       :   maximum_iterations: Maximum number of iterations allowed before termination.
 * @param       :   convergence_tolerance: Target residual norm tolerance for convergence.
 * @return      :   Status: Status code indicating success or specific failure mode.
 */
Status solve(MultiplyFunction callback, int number_of_dimensions, void* context, 
    const double* operand, const double* initial_guess, double* solution_vector, int maximum_iterations, double convergence_tolerance)
{
    // Defend against null pointer arguments.
    if (callback == NULL || initial_guess == NULL || operand == NULL || solution_vector == NULL)
    {
        return ERROR_NULL_POINTER;
    }

    // Defend against invalid matrix dimensions.
    if (number_of_dimensions <= 0)
    {
        return ERROR_INVALID_DIMENSION;
    }

    // Defend against invalid iteration limits or negative tolerance.
    if (maximum_iterations <= 0 || convergence_tolerance < 0.0)
    {
        return ERROR_INVALID_DIMENSION;
    }

    Status execute_status = SUCCESS;

    int index = 0;

    // Allocate dynamic working memory buffers with defensive checks.
    double* residual = (double*)malloc(sizeof(double) * number_of_dimensions);
    double* shadow_residual = (double*)malloc(sizeof(double) * number_of_dimensions);
    double* search_direction = (double*)malloc(sizeof(double) * number_of_dimensions);
    double* product = (double*)malloc(sizeof(double) * number_of_dimensions);
    double* intermediate = (double*)malloc(sizeof(double) * number_of_dimensions);
    double* multiply = (double*)malloc(sizeof(double) * number_of_dimensions);

    // Verify that all memory allocations succeeded successfully.
    if (residual == NULL || shadow_residual == NULL || 
        search_direction == NULL || product == NULL || 
        intermediate == NULL || multiply == NULL)
    {
        // Safely free any successfully allocated buffers before returning allocation failure.
        (residual != NULL) ? free(residual) : (void)0;
        (shadow_residual != NULL) ? free(shadow_residual) : (void)0;
        (search_direction != NULL) ? free(search_direction) : (void)0;
        (product != NULL) ? free(product) : (void)0;
        (intermediate != NULL) ? free(intermediate) : (void)0;
        (multiply != NULL) ? free(multiply) : (void)0;

        return ERROR_ALLOCATION_FAILED;
    }

    double current_inner_product = 1.0;
    double previous_inner_product = 1.0;
    double scalar_alpha = 1.0;
    double scalar_omega_current = 1.0;
    double scalar_beta = 0.0;

    // Initialize solution vector and working vectors to zero.
    for (index = 0; index < number_of_dimensions; index++)
    {
        solution_vector[index] = initial_guess[index];
        search_direction[index] = 0.0;
        product[index] = 0.0;
    }

    // Compute initial matrix-vector product using the initial guess solution vector.
    callback(number_of_dimensions, solution_vector, residual, context);

    // Compute initial residual vector by subtracting the product from the operand vector.
    for (index = 0; index < number_of_dimensions; index++)
    {
        residual[index] = operand[index] - residual[index];
    }

    // Check if initial guess already satisfies convergence tolerance.
    if (compute_vector_normalization(number_of_dimensions, residual) <= convergence_tolerance)
    {
        execute_status = SUCCESS;
        
        // Free allocated memory buffers before returning successful early exit.
        free(residual);
        free(shadow_residual);
        free(search_direction);
        free(product);
        free(intermediate);
        free(multiply);

        return execute_status;
    }

    // Populate shadow residual vector with initial residual values.
    for (index = 0; index < number_of_dimensions; index++)
    {
        shadow_residual[index] = residual[index];
    }

    int iteration_counter;

    // Execute the main iterative loop of the BiCGSTAB algorithm.
    for (iteration_counter = 1; iteration_counter <= maximum_iterations; iteration_counter++)
    {
        // Compute current inner product rho between shadow residual and residual vectors.
        current_inner_product = compute_vector_dot_product(number_of_dimensions, shadow_residual, residual);

        // Defend against mathematical breakdown where rho inner product is zero.
        if (fabs(current_inner_product) < 1e-16)
        {
            execute_status = ERROR_BREAKDOWN_RHO_ZERO;

            break;
        }

        // Handle first iteration search direction setup separately from subsequent iterations.
        if (iteration_counter == 1)
        {
            for (index = 0; index < number_of_dimensions; index++)
            {
                search_direction[index] = residual[index];
            }
        }
        else
        {
            // Defend against division by zero when omega parameter is zero.
            if (fabs(scalar_omega_current) < 1e-16)
            {
                execute_status = ERROR_BREAKDOWN_OMEGA_ZERO;

                break;
            }

            // Compute scaling factor beta for the search direction update.
            scalar_beta = (current_inner_product / previous_inner_product) * (scalar_alpha / scalar_omega_current);

            // Update search direction vector using residual and previous product vectors.
            for (index = 0; index < number_of_dimensions; index++)
            {
                search_direction[index] = residual[index] + 
                    scalar_beta * (search_direction[index] - scalar_omega_current * product[index]);
            }
        }

        // Compute matrix-vector product using the callback function with search direction.
        callback(number_of_dimensions, search_direction, product, context);

        // Compute denominator dot product between shadow residual and product vector.
        double denominator_shadow_product = compute_vector_dot_product(number_of_dimensions, shadow_residual, product);

        // Defend against division by zero when shadow product denominator is zero.
        if (fabs(denominator_shadow_product) < 1e-16)
        {
            execute_status = ERROR_BREAKDOWN_ALPHA_ZERO;
            
            break;
        }

        // Compute scaling factor alpha.
        scalar_alpha = current_inner_product / denominator_shadow_product;

        // Compute intermediate residual vector.
        for (index = 0; index < number_of_dimensions; index++)
        {
            intermediate[index] = residual[index] - scalar_alpha * product[index];
        }

        // Check for early convergence on the intermediate residual norm.
        if (compute_vector_normalization(number_of_dimensions, intermediate) <= convergence_tolerance)
        {
            for (index = 0; index < number_of_dimensions; index++)
            {
                solution_vector[index] += scalar_alpha * search_direction[index];
            }

            execute_status = SUCCESS;

            break;
        }

        // Compute second matrix-vector product using the callback function with intermediate vector.
        callback(number_of_dimensions, intermediate, multiply, context);

        // Compute numerator and denominator dot products for stabilization parameter omega.
        double numerator_multiply_intermediate = compute_vector_dot_product(number_of_dimensions, multiply, intermediate);
        double denominator_multiply_multiply = compute_vector_dot_product(number_of_dimensions, multiply, multiply);

        // Defend against division by zero when stabilization denominator is zero.
        if (fabs(denominator_multiply_multiply) < 1e-16)
        {
            execute_status = ERROR_BREAKDOWN_OMEGA_ZERO;

            break;
        }

        // Compute stabilization parameter omega.
        scalar_omega_current = numerator_multiply_intermediate / denominator_multiply_multiply;

        // Defend against omega being zero after computation.
        if (fabs(scalar_omega_current) < 1e-16)
        {
            execute_status = ERROR_BREAKDOWN_OMEGA_ZERO;

            break;
        }

        // Update solution vector and residual vector for the current iteration step.
        for (index = 0; index < number_of_dimensions; index++)
        {
            solution_vector[index] += scalar_alpha * search_direction[index] + scalar_omega_current * intermediate[index];
            residual[index] = intermediate[index] - scalar_omega_current * multiply[index];
        }

        // Check convergence criteria against current residual normalization norm.
        if (compute_vector_normalization(number_of_dimensions, residual) <= convergence_tolerance)
        {
            execute_status = SUCCESS;

            break;
        }

        // Store current inner product rho as previous inner product for the next iteration step.
        previous_inner_product = current_inner_product;
    }

    // Check if the loop terminated because maximum allowed iterations were exceeded.
    if (iteration_counter > maximum_iterations && execute_status == SUCCESS)
    {
        execute_status = ERROR_MAXIMUM_ITERATIONS_EXCEEDED;
    }

    // Safely deallocate all dynamic working memory buffers to prevent memory leaks.
    free(residual);
    free(shadow_residual);
    free(search_direction);
    free(product);
    free(intermediate);
    free(multiply);

    return execute_status;
}