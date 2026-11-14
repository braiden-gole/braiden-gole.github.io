/**
 * @file        :   biconjugate_gradient_stabilized_example.c
 * @author      :   Braiden Gole
 * @brief       :   Example usage of the Biconjugate Gradient Stabilized (BiCGSTAB) linear solver algorithm.
 * @details     :   This source file demonstrates how to set up a linear system Ax = b, define a custom 
 *                  matrix-vector multiplication callback function, execute the BiCGSTAB solver routine, 
 *                  and validate the resulting solution vector.
 * @version     :   1.0
 * @date        :   2026-11-11
 * 
 * Compile:
 * clang ./src/darwin/numerical/biconjugate-gradient-stabilized/main.c ./src/darwin/numerical/biconjugate-gradient-stabilized/biconjugate_gradient_stabilized.c -I./include -lm -o ./bin/biconjugate_gradient_stabilized_test
 */
#include <stdio.h>
#include <stdlib.h>

// Include header files.
#include "darwin/numerical/biconjugate-gradient-stabilized/biconjugate_gradient_stabilized.h"

/**
 * @brief       :   Custom matrix-vector multiplication callback for a sample 3x3 tridiagonal matrix.
 * @details     :   Computes product_vector = matrix_a * input_vector where matrix_a is defined as:
 *                  [ 4.0, 1.0, 0.0 ]
 *                  [ 1.0, 4.0, 1.0 ]
 *                  [ 0.0, 1.0, 4.0 ]
 * @param       :   number_of_dimensions: The dimension size of the vectors.
 * @param       :   input_vector: Pointer to the input vector.
 * @param       :   product_vector: Pointer to store the resulting matrix-vector product.
 * @param       :   context: User-defined context pointer passed through the solver.
 */
static void sample_matrix_vector_multiply(int number_of_dimensions, const double* vector, double* operand, void* context)
{
    // Validate input and output pointers to prevent memory access faults.
    if (vector == NULL || operand == NULL || number_of_dimensions != 3)
    {
        return;
    }

    // Suppress unused context parameter warning.
    (void)context;

    // Perform row operations for the 3x3 tridiagonal matrix multiplication.
    operand[0] = 4.0 * vector[0] + 1.0 * vector[1];
    operand[1] = 1.0 * vector[0] + 4.0 * vector[1] + 1.0 * vector[2];
    operand[2] = 1.0 * vector[1] + 4.0 * vector[2];
}

int main(void)
{
    // Define the matrix dimension size for our sample linear system.
    int number_of_dimensions = 3;

    // Define the right-hand side vector b corresponding to exact solution [1.0, 1.0, 1.0].
    double operand[3] = {5.0, 6.0, 5.0};

    // Define an initial guess vector initialized to zero.
    double initial_guess[3] = {0.0, 0.0, 0.0};

    // Allocate memory buffer to store the computed solution vector.
    double* computed_solution_vector = (double*)malloc(sizeof(double) * number_of_dimensions);

    // Defend against memory allocation failure.
    if (computed_solution_vector == NULL)
    {
        fprintf(stderr, "%s", "Error: Failed to allocate memory for computed solution vector.\n");

        return 1;
    }

    // Configure solver hyperparameters and iteration limits.
    int maximum_allowed_iterations = 100;

    double target_convergence_tolerance = 1e-6;

    // Execute the BiCGSTAB solver routine.
    Status execution_status = solve(
        sample_matrix_vector_multiply,
        number_of_dimensions,
        NULL,
        operand,
        initial_guess,
        computed_solution_vector,
        maximum_allowed_iterations,
        target_convergence_tolerance
    );

    // Evaluate the returned status code and handle success or failure modes.
    if (execution_status == SUCCESS)
    {
        printf("%s", "BiCGSTAB solver completed successfully.\n");
        printf("%s", "Computed Solution Vector x:\n");
        
        for (int index_iterator = 0; index_iterator < number_of_dimensions; index_iterator++)
        {
            printf("%cx[%d] = %.6f\n", '\t', index_iterator, computed_solution_vector[index_iterator]);
        }
    }
    else
    {
        fprintf(stderr, "%s %d%c", "BiCGSTAB solver failed with status code:", execution_status, '\n');
    }

    // Safely deallocate the computed solution vector memory buffer.
    free(computed_solution_vector);

    return EXIT_SUCCESS;
}