/**
 * @file        :   interior_point_method.c
 * @author      :   Braiden Gole
 * @brief       :   Implementation of the Primal-Dual Path-Following Interior Point Method.
 * @details     :   Solves linear programming problems in standard form using interior point 
 *                  optimization with dynamic barrier parameter updates, full residual tracking, and Gaussian elimination.
 * @version     :   3.0
 * @date        :   2026-11-12
 */
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <stdint.h>
#include <stdbool.h>

// Include header files.
#include "darwin/numerical/primary-dual-path-following-interior-point-method/interior_point_method.h"

/**
 * @brief           :   Allocates and initializes a mathematical vector structure.
 * @param length    :   The number of elements contained in the vector.
 * @return          :   Pointer to the initialized vector, or NULL if allocation fails.
 */
Vector* initialize_vector(size_t length)
{
    // Validate that the requested length is non-zero and within safe memory bounds to prevent integer overflows.
    if (length == 0 || length > SIZE_MAX / sizeof(double))
    {
        return NULL;
    }

    // Allocate memory for the vector structure container.
    Vector* pointer = (Vector*)malloc(sizeof(Vector));

    if (pointer == NULL)
    {
        return NULL;
    }

    // Allocate zero-initialized memory for the vector values buffer.
    pointer->values = (double*)calloc(length, sizeof(double));

    if (pointer->values == NULL)
    {
        // Free the container structure if the values buffer allocation fails.
        free(pointer);

        return NULL;
    }

    // Assign the validated length to the vector structure.
    pointer->length = length;

    return pointer;
}

/**
 * @brief                   :   Allocates and initializes a mathematical matrix structure.
 * @param number_of_rows    :   The row count of the matrix.
 * @param number_of_columns :   The column count of the matrix.
 * @return                  :   Pointer to the initialized matrix, or NULL if allocation fails.
 */
Matrix* initialize_matrix(size_t number_of_rows, size_t number_of_columns)
{
    // Validate that matrix dimensions are positive non-zero integers.
    if (number_of_rows == 0 || number_of_columns == 0)
    {
        return NULL;
    }

    // Enforce strict overflow protection for row and column dimension multiplication.
    if (number_of_rows > SIZE_MAX / number_of_columns)
    {
        return NULL;
    }

    // Allocate memory for the matrix structure container.
    Matrix* matrix = (Matrix*)malloc(sizeof(Matrix));

    if (matrix == NULL)
    {
        return NULL;
    }

    // Allocate zero-initialized memory for the matrix values buffer.
    matrix->values = (double*)calloc(number_of_rows * number_of_columns, sizeof(double));

    if (matrix->values == NULL)
    {
        // Free the matrix structure container if the values buffer allocation fails.
        free(matrix);

        return NULL;
    }

    // Assign validated dimensions to the matrix structure.
    matrix->number_of_rows = number_of_rows;
    matrix->number_of_columns = number_of_columns;

    return matrix;
}

/**
 * @brief                       :   Solves a linear system using Gaussian elimination with partial pivoting.
 * @param coefficient_matrix    :   Pointer to the square coefficient matrix.
 * @param constant_vector       :   Pointer to the right-hand side constant vector.
 * @param solution_vector       :   Pointer to the output solution vector to be populated.
 * @return                      :   True if resolution succeeds; otherwise false.
 */
bool solve_linear_system_gaussian_elimination(Matrix* coefficient_matrix, Vector* constant_vector, Vector* solution_vector)
{
    // Validate input pointers to prevent segmentation faults.
    if (coefficient_matrix == NULL || constant_vector == NULL || solution_vector == NULL)
    {
        return false;
    }

    // Verify that internal values buffers are allocated and valid.
    if (coefficient_matrix->values == NULL || constant_vector->values == NULL || solution_vector->values == NULL)
    {
        return false;
    }

    size_t dimension = coefficient_matrix->number_of_rows;

    // Validate that matrix dimensions and vector lengths match.
    if (dimension != coefficient_matrix->number_of_columns || dimension != constant_vector->length || dimension != solution_vector->length)
    {
        return false;
    }

    size_t augmented_column_count = dimension + 1;

    // Allocate memory for the augmented matrix workspace buffer with overflow protection.
    if (dimension > SIZE_MAX / augmented_column_count || (dimension * augmented_column_count) > SIZE_MAX / sizeof(double))
    {
        return false;
    }

    double* augmented_matrix = (double*)malloc(dimension * augmented_column_count * sizeof(double));

    if (augmented_matrix == NULL)
    {
        return false;
    }

    // Construct the augmented matrix workspace buffer by combining coefficient matrix and constant vector.
    for (size_t row_index = 0; row_index < dimension; row_index++)
    {
        for (size_t column_index = 0; column_index < dimension; column_index++)
        {
            augmented_matrix[row_index * augmented_column_count + column_index] = 
                coefficient_matrix->values[row_index * dimension + column_index];
        }

        augmented_matrix[row_index * augmented_column_count + dimension] = constant_vector->values[row_index];
    }

    const double singularity_tolerance = 1e-12;

    // Execute forward elimination with partial pivoting for numerical stability.
    for (size_t pivot_row = 0; pivot_row < dimension; pivot_row++)
    {
        size_t maximum_row = pivot_row;

        double maximum_pivot_value = fabs(augmented_matrix[pivot_row * augmented_column_count + pivot_row]);

        // Search for the row with the largest pivot candidate in the current column.
        for (size_t candidate_row = pivot_row + 1; candidate_row < dimension; candidate_row++)
        {
            double candidate_value = fabs(augmented_matrix[candidate_row * augmented_column_count + pivot_row]);

            if (candidate_value > maximum_pivot_value)
            {
                maximum_pivot_value = candidate_value;
                maximum_row = candidate_row;
            }
        }

        // Verify that the matrix is non-singular and the pivot exceeds numerical tolerance limits.
        if (maximum_pivot_value < singularity_tolerance || isnan(maximum_pivot_value) || isinf(maximum_pivot_value))
        {
            free(augmented_matrix);

            return false;
        }

        // Swap rows if a better pivot candidate was found to maintain stability.
        if (maximum_row != pivot_row)
        {
            for (size_t column_index = pivot_row; column_index < augmented_column_count; column_index++)
            {
                double temporary_value = augmented_matrix[pivot_row * augmented_column_count + column_index];

                augmented_matrix[pivot_row * augmented_column_count + column_index] = augmented_matrix[maximum_row * augmented_column_count + column_index];
                augmented_matrix[maximum_row * augmented_column_count + column_index] = temporary_value;
            }
        }

        // Eliminate column entries below the pivot row.
        for (size_t target_row = pivot_row + 1; target_row < dimension; target_row++)
        {
            double scaling_factor = augmented_matrix[target_row * augmented_column_count + pivot_row] / augmented_matrix[pivot_row * augmented_column_count + pivot_row];

            for (size_t column_index = pivot_row; column_index < augmented_column_count; column_index++)
            {
                augmented_matrix[target_row * augmented_column_count + column_index] -= scaling_factor * augmented_matrix[pivot_row * augmented_column_count + column_index];
            }
        }
    }

    // Execute back substitution to resolve the solution vector from the upper triangular form.
    for (long long row_index = (long long)dimension - 1; row_index >= 0; row_index--)
    {
        double current_sum = augmented_matrix[row_index * augmented_column_count + dimension];

        // Subtract known variable contributions from the running sum.
        for (size_t column_index = row_index + 1; column_index < dimension; column_index++)
        {
            current_sum -= augmented_matrix[row_index * augmented_column_count + column_index] * solution_vector->values[column_index];
        }

        double diagonal_coefficient = augmented_matrix[row_index * augmented_column_count + row_index];

        // Verify that the diagonal coefficient is non-zero before performing division.
        if (fabs(diagonal_coefficient) < singularity_tolerance || isnan(diagonal_coefficient) || isinf(diagonal_coefficient))
        {
            free(augmented_matrix);

            return false;
        }

        // Compute the final value for the current solution vector index.
        solution_vector->values[row_index] = current_sum / diagonal_coefficient;
    }

    // Free the temporary augmented matrix workspace buffer.
    free(augmented_matrix);

    return true;
}

/**
 * @brief                           :   Executes the primal-dual path-following interior point optimization algorithm.
 * @param matrix                    :   Pointer to the linear constraint coefficient matrix.
 * @param vector                    :   Pointer to the constraint bound vector.
 * @param objective                 :   Pointer to the objective function coefficient vector.
 * @param initial_solution          :   Pointer to the initial primal guess vector (modified in place with optimal values).
 * @param maximum_iterations        :   The maximum iteration limit allowed before termination.
 * @param convergence_tolerance     :   The target duality gap tolerance for convergence.
 * @return                          :   True if optimization converges successfully; otherwise false.
 */
bool execute_primal_dual_interior_point_method(Matrix* matrix, Vector* vector, Vector* objective, Vector* initial_solution, size_t maximum_iterations, double convergence_tolerance)
{
    // Validate input pointers to prevent null pointer dereferences.
    if (matrix == NULL || vector == NULL || objective == NULL || initial_solution == NULL)
    {
        return false;
    }

    // Validate internal memory buffers of all input structures.
    if (matrix->values == NULL || vector->values == NULL || objective->values == NULL || initial_solution->values == NULL)
    {
        return false;
    }

    size_t number_of_rows = matrix->number_of_rows;
    size_t number_of_columns = matrix->number_of_columns;

    // Verify dimension compatibility between the constraint matrix, bounds vector, objective vector, and initial solution.
    if (number_of_rows != vector->length || number_of_columns != objective->length || initial_solution->length != number_of_columns)
    {
        return false;
    }

    // Verify that the problem system dimensions are valid for interior point execution.
    if (number_of_columns < number_of_rows)
    {
        return false;
    }

    // Allocate workspace tracking vectors for primal variables, dual variables, and slack variables.
    Vector* primal = initialize_vector(number_of_columns);
    Vector* dual = initialize_vector(number_of_rows);
    Vector* stack = initialize_vector(number_of_columns);

    // Verify successful allocation of all working vectors.
    if (primal == NULL || dual == NULL || stack == NULL)
    {
        erase_vector(primal);
        erase_vector(dual);
        erase_vector(stack);

        return false;
    }

    // Initialize primal variables and slack variables to strictly positive values to satisfy interior constraints.
    for (size_t index = 0; index < number_of_columns; index++)
    {
        primal->values[index] = initial_solution->values[index] > 0.0 ? initial_solution->values[index] : 1.0;
        stack->values[index] = 1.0;
    }

    // Run the iterative path-following optimization loop up to the maximum iteration limit.
    for (size_t iteration = 0; iteration < maximum_iterations; iteration++)
    {
        // Allocate residual vectors for primal, dual, and complementarity conditions.
        Vector* primal_residual = initialize_vector(number_of_rows);
        Vector* dual_residual = initialize_vector(number_of_columns);
        Vector* complementarity_residual = initialize_vector(number_of_columns);

        if (primal_residual == NULL || dual_residual == NULL || complementarity_residual == NULL)
        {
            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return false;
        }

        // Compute primal residual: r_p = b - A * x.
        for (size_t row_index = 0; row_index < number_of_rows; row_index++)
        {
            double constraint_sum = 0.0;

            for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
            {
                constraint_sum += matrix->values[row_index * number_of_columns + variable_index] * primal->values[variable_index];
            }

            primal_residual->values[row_index] = vector->values[row_index] - constraint_sum;
        }

        // Compute dual residual: r_d = c - A^T * y - s.
        for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
        {
            double transpose_product_sum = 0.0;

            for (size_t row_index = 0; row_index < number_of_rows; row_index++)
            {
                transpose_product_sum += matrix->values[row_index * number_of_columns + variable_index] * dual->values[row_index];
            }

            dual_residual->values[variable_index] = objective->values[variable_index] - transpose_product_sum - stack->values[variable_index];
        }

        // Compute current duality gap across all variables.
        double duality_gap = 0.0;

        for (size_t index = 0; index < number_of_columns; index++)
        {
            duality_gap += primal->values[index] * stack->values[index];
        }

        duality_gap /= (double)number_of_columns;

        // Check if duality gap evaluation resulted in numerical instability.
        if (isnan(duality_gap) || isinf(duality_gap))
        {
            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return false;
        }

        // Check if convergence criteria (duality gap and residual magnitudes) have been satisfied.
        bool is_primal_feasible = true;
        bool is_dual_feasible = true;

        for (size_t row_index = 0; row_index < number_of_rows; row_index++)
        {
            if (fabs(primal_residual->values[row_index]) > convergence_tolerance * 10.0)
            {
                is_primal_feasible = false;

                break;
            }
        }

        for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
        {
            if (fabs(dual_residual->values[variable_index]) > convergence_tolerance * 10.0)
            {
                is_dual_feasible = false;

                break;
            }
        }

        if (duality_gap < convergence_tolerance && is_primal_feasible && is_dual_feasible)
        {
            for (size_t index = 0; index < number_of_columns; index++)
            {
                initial_solution->values[index] = primal->values[index];
            }

            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return true;
        }

        double centering_parameter = 0.1;
        double barrier_parameter = centering_parameter * duality_gap;

        // Compute complementarity residual: r_c = mu * e - X * s.
        for (size_t index = 0; index < number_of_columns; index++)
        {
            complementarity_residual->values[index] = barrier_parameter - (primal->values[index] * stack->values[index]);
        }

        // Allocate search direction vectors for Newton step calculations.
        Vector* delta_primal = initialize_vector(number_of_columns);
        Vector* delta_dual = initialize_vector(number_of_rows);
        Vector* delta_slack = initialize_vector(number_of_columns);

        if (delta_primal == NULL || delta_dual == NULL || delta_slack == NULL)
        {
            erase_vector(delta_primal);
            erase_vector(delta_dual);
            erase_vector(delta_slack);
            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return false;
        }

        // Allocate Schur complement matrix and right-hand side vector for normal equations.
        Matrix* schur_complement_matrix = initialize_matrix(number_of_rows, number_of_rows);
        Vector* schur_right_hand_side_vector = initialize_vector(number_of_rows);

        if (schur_complement_matrix == NULL || schur_right_hand_side_vector == NULL)
        {
            erase_matrix(schur_complement_matrix);
            erase_vector(schur_right_hand_side_vector);
            erase_vector(delta_primal);
            erase_vector(delta_dual);
            erase_vector(delta_slack);
            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return false;
        }

        Vector* scaling_weight = initialize_vector(number_of_columns);

        if (scaling_weight == NULL)
        {
            erase_matrix(schur_complement_matrix);
            erase_vector(schur_right_hand_side_vector);
            erase_vector(delta_primal);
            erase_vector(delta_dual);
            erase_vector(delta_slack);
            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return false;
        }

        // Assemble diagonal scaling weights where each weight D = x / s.
        for (size_t index = 0; index < number_of_columns; index++)
        {
            if (stack->values[index] <= 0.0)
            {
                erase_vector(scaling_weight);
                erase_matrix(schur_complement_matrix);
                erase_vector(schur_right_hand_side_vector);
                erase_vector(delta_primal);
                erase_vector(delta_dual);
                erase_vector(delta_slack);
                erase_vector(primal_residual);
                erase_vector(dual_residual);
                erase_vector(complementarity_residual);
                erase_vector(primal);
                erase_vector(dual);
                erase_vector(stack);

                return false;
            }

            scaling_weight->values[index] = primal->values[index] / stack->values[index];
        }

        // Construct the Schur complement matrix expression: A * D * A^T.
        for (size_t row_index = 0; row_index < number_of_rows; row_index++)
        {
            for (size_t column_index = 0; column_index < number_of_rows; column_index++)
            {
                double cumulative_sum = 0.0;

                for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
                {
                    double constraint_element_row = matrix->values[row_index * number_of_columns + variable_index];
                    double constraint_element_column = matrix->values[column_index * number_of_columns + variable_index];

                    cumulative_sum += constraint_element_row * scaling_weight->values[variable_index] * constraint_element_column;
                }

                schur_complement_matrix->values[row_index * number_of_rows + column_index] = cumulative_sum;
            }
        }

        // Assemble Schur right-hand side: r_p - A * (S^(-1) * r_c - D * r_d).
        for (size_t row_index = 0; row_index < number_of_rows; row_index++)
        {
            double constraint_term_sum = 0.0;

            for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
            {
                double constraint_value = matrix->values[row_index * number_of_columns + variable_index];
                double inverse_slack_term = complementarity_residual->values[variable_index] / stack->values[variable_index];
                double scaling_residual_term = scaling_weight->values[variable_index] * dual_residual->values[variable_index];

                constraint_term_sum += constraint_value * (inverse_slack_term - scaling_residual_term);
            }

            schur_right_hand_side_vector->values[row_index] = primal_residual->values[row_index] - constraint_term_sum;
        }

        // Solve for the dual search direction vector using Gaussian elimination.
        bool status = solve_linear_system_gaussian_elimination(schur_complement_matrix, schur_right_hand_side_vector, delta_dual);

        if (!status)
        {
            erase_vector(scaling_weight);
            erase_matrix(schur_complement_matrix);
            erase_vector(schur_right_hand_side_vector);
            erase_vector(delta_primal);
            erase_vector(delta_dual);
            erase_vector(delta_slack);
            erase_vector(primal_residual);
            erase_vector(dual_residual);
            erase_vector(complementarity_residual);
            erase_vector(primal);
            erase_vector(dual);
            erase_vector(stack);

            return false;
        }

        // Back-calculate slack and primal search directions with correct mathematical signs.
        for (size_t index = 0; index < number_of_columns; index++)
        {
            double constraint_product_sum = 0.0;

            for (size_t constraint_index = 0; constraint_index < number_of_rows; constraint_index++)
            {
                constraint_product_sum += matrix->values[constraint_index * number_of_columns + index] * delta_dual->values[constraint_index];
            }

            // Correct sign for slack search direction: delta_s = r_d - A^T * delta_y.
            delta_slack->values[index] = dual_residual->values[index] - constraint_product_sum;
            
            // Correct sign for primal search direction: delta_x = S^(-1) * r_c - D * delta_s.
            double inverse_slack_term = complementarity_residual->values[index] / stack->values[index];
            double scaling_slack_term = scaling_weight->values[index] * delta_slack->values[index];
            
            delta_primal->values[index] = inverse_slack_term - scaling_slack_term;
        }

        double step_length_primal = 1.0;
        double step_length_dual = 1.0;

        const double boundary_safety_factor = 0.99;

        // Determine maximum allowable step lengths maintaining strict variable positivity via fraction-to-boundary rule.
        for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
        {
            if (delta_primal->values[variable_index] < 0.0)
            {
                double candidate_step = -boundary_safety_factor * primal->values[variable_index] / delta_primal->values[variable_index];

                if (candidate_step < step_length_primal)
                {
                    step_length_primal = candidate_step;
                }
            }

            if (delta_slack->values[variable_index] < 0.0)
            {
                double candidate_step = -boundary_safety_factor * stack->values[variable_index] / delta_slack->values[variable_index];

                if (candidate_step < step_length_dual)
                {
                    step_length_dual = candidate_step;
                }
            }
        }

        // Update primal and slack variables using calculated step lengths.
        for (size_t variable_index = 0; variable_index < number_of_columns; variable_index++)
        {
            primal->values[variable_index] += step_length_primal * delta_primal->values[variable_index];
            stack->values[variable_index] += step_length_dual * delta_slack->values[variable_index];
        }

        // Update dual variables using calculated step lengths.
        for (size_t constraint_index = 0; constraint_index < number_of_rows; constraint_index++)
        {
            dual->values[constraint_index] += step_length_dual * delta_dual->values[constraint_index];
        }

        // Clean up temporary memory allocated specifically for the current iteration.
        erase_vector(scaling_weight);
        erase_matrix(schur_complement_matrix);
        erase_vector(schur_right_hand_side_vector);
        erase_vector(delta_primal);
        erase_vector(delta_dual);
        erase_vector(delta_slack);
        erase_vector(primal_residual);
        erase_vector(dual_residual);
        erase_vector(complementarity_residual);
    }

    // Clean up persistent working vectors before returning failure upon reaching maximum iterations.
    erase_vector(primal);
    erase_vector(dual);
    erase_vector(stack);

    return false;
}

/**
 * @brief           :   Safely deallocates a vector structure and prevents dangling pointers.
 * @param pointer   :   Pointer to the vector structure to erase.
 */
void erase_vector(Vector* pointer)
{
    // Verify that the vector pointer is non-null before attempting deallocation.
    if (pointer != NULL)
    {
        // Free the internal values buffer if it has been allocated.
        if (pointer->values != NULL)
        {
            free(pointer->values);

            // Nullify the pointer immediately after freeing to prevent use-after-free bugs.
            pointer->values = NULL;
        }

        // Free the vector container structure itself.
        free(pointer);
    }
}

/**
 * @brief           :   Safely deallocates a matrix structure and prevents dangling pointers.
 * @param matrix    :   Pointer to the matrix structure to erase.
 */
void erase_matrix(Matrix* matrix)
{
    // Verify that the matrix pointer is non-null before attempting deallocation.
    if (matrix != NULL)
    {
        // Free the internal values buffer if it has been allocated.
        if (matrix->values != NULL)
        {
            free(matrix->values);

            // Nullify the pointer immediately after freeing to prevent use-after-free bugs.
            matrix->values = NULL;
        }

        // Free the matrix container structure itself.
        free(matrix);
    }
}