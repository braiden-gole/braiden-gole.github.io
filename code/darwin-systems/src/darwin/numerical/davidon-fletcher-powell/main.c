/**
 * @file        :   main.c
 * @author      :   Braiden Gole
 * @brief       :   Driver program for testing the Davidon-Fletcher-Powell (DFP) Quasi-Newton optimizer.
 * @details     :   This application file imports the DFP optimization header and executes the 
 *                  minimization of the 2D Rosenbrock banana function. It defines the objective and 
 *                  gradient evaluation functions, configures optimization hyperparameters, runs the 
 *                  solver, and robustly reports convergence results or execution errors.
 * @version     :   1.0
 * @date        :   2026-11-11
 * 
 * Compile:
 * clang ./src/darwin/numerical/davidon-fletcher-powell/main.c ./src/darwin/numerical/davidon-fletcher-powell/davidon_fletcher_powell.c -I./include -o ./bin/davidon_fletcher_powell_test
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

// Include header files.
#include "darwin/numerical/davidon-fletcher-powell/davidon_fletcher_powell.h"

int main(void)
{
    // Define problem space dimensions and optimization hyperparameters.
    double initial_guess_coordinates[2] = {-1.2, 1.0};

    const size_t space_dimension = 2;
    size_t maximum_optimization_iterations = 1000;
    
    double gradient_tolerance = 1e-6;
    double parameter_tolerance = 1e-6;

    // Print initial starting coordinates to standard output.
    printf("Starting optimization point: x = %f, y = %f\n", initial_guess_coordinates[0], initial_guess_coordinates[1]);

    // Execute the DFP optimization algorithm with defensive parameter passing.
    bool optimization_success = optimize_quasi_newton_dfp(
        evaluate_rosenbrock_objective,
        evaluate_rosenbrock_gradient,
        initial_guess_coordinates,
        space_dimension,
        maximum_optimization_iterations,
        gradient_tolerance,
        parameter_tolerance
    );

    // Evaluate convergence status and report results.
    if (optimization_success)
    {
        printf("Optimization converged successfully.\n");
        printf("Minimum found at: x = %f, y = %f\n", initial_guess_coordinates[0], initial_guess_coordinates[1]);
        printf("Objective function value at minimum: %e\n", evaluate_rosenbrock_objective(initial_guess_coordinates, space_dimension));
    }
    else
    {
        // Output failure notification to standard error stream.
        fprintf(stderr, "Error: Optimization failed to converge or encountered memory allocation failures.\n");

        return 1;
    }

    return 0;
}