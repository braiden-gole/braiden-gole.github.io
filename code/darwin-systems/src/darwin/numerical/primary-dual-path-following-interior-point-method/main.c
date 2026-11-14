/**
 * @file        :   main.c
 * @author      :   Braiden Gole
 * @brief       :   Executable driver for the Primal-Dual Path-Following Interior Point Method.
 * @details     :   Sets up a linear programming problem, executes the interior point solver, 
 *                  and outputs the optimal primal variables.
 * @version     :   1.0
 * @date        :   2026-11-12
 * 
 * Compile:
 * clang ./src/darwin/numerical/primary-dual-path-following-interior-point-method/main.c ./src/darwin/numerical/primary-dual-path-following-interior-point-method/interior_point_method.h.c -I./include -lm -o ./bin/interior_point_method_test
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

// Include header files.
#include "darwin/numerical/primary-dual-path-following-interior-point-method/interior_point_method.h"

int main(void)
{
    // Define problem dimensions and variable counts.
    size_t length = 2;
    size_t count = 3;

    // Allocate constraint matrix and working vectors.
    Matrix* matrix = initialize_matrix(length, count);
    Vector* resource_limit = initialize_vector(length);
    Vector* objective = initialize_vector(count);
    Vector* initial_solution = initialize_vector(count);

    // Validate that all structures were successfully allocated.
    if (matrix == NULL || resource_limit == NULL || objective == NULL || initial_solution == NULL)
    {
        fprintf(stderr, "%s", "Error: Memory allocation failed for optimization problem structures.\n");

        erase_matrix(matrix);
        erase_vector(resource_limit);
        erase_vector(objective);
        erase_vector(initial_solution);

        return 1;
    }

    // Populate constraint coefficient matrix values.
    matrix->values[0 * count + 0] = 1.0;
    matrix->values[0 * count + 1] = 1.0;
    matrix->values[0 * count + 2] = 1.0;
    
    matrix->values[1 * count + 0] = 2.0;
    matrix->values[1 * count + 1] = 1.0;
    matrix->values[1 * count + 2] = 0.0;

    // Populate resource limit vector values.
    resource_limit->values[0] = 10.0;
    resource_limit->values[1] = 14.0;

    // Populate objective function coefficient values.
    objective->values[0] = -1.0;
    objective->values[1] = -2.0;
    objective->values[2] = 0.0;

    // Populate initial feasible primal guess vector values.
    initial_solution->values[0] = 1.0;
    initial_solution->values[1] = 1.0;
    initial_solution->values[2] = 1.0;

    // Define algorithm configuration parameters.
    size_t maximum_iterations = 100;

    double convergence_tolerance = 1e-8;

    // Execute the primal-dual path-following interior point optimization algorithm.
    bool status = execute_primal_dual_interior_point_method(
        matrix,
        resource_limit,
        objective,
        initial_solution,
        maximum_iterations,
        convergence_tolerance
    );

    // Evaluate optimization status and print results or errors accordingly.
    if (status)
    {
        printf("%s", "Optimization completed successfully.\n");
        printf("%s", "Optimal primal solution variables:\n");

        for (size_t index = 0; index < count; index++)
        {
            printf("\tx[%zu] = %.6f\n", index, initial_solution->values[index]);
        }
    }
    else
    {
        fprintf(stderr, "Error: Interior point optimization failed to converge.\n");
    }

    // Clean up all dynamically allocated memory resources safely.
    erase_matrix(matrix);
    erase_vector(resource_limit);
    erase_vector(objective);
    erase_vector(initial_solution);

    return status ? 0 : 1;
}