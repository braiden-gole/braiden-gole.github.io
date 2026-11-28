/**
 * @file        :   interior_point_method.h
 * @author      :   Braiden Gole
 * @brief       :   Header file for the Primal-Dual Path-Following Interior Point Method.
 * @details     :   Defines vector and matrix data structures, memory management utilities, 
 *                  Gaussian elimination linear solvers, and the interior point optimization routine.
 * @version     :   1.0
 * @date        :   2026-11-12
 */
#ifndef INTERIOR_POINT_METHOD_H
#define INTERIOR_POINT_METHOD_H

#include <stddef.h>
#include <stdbool.h>

/**
 * @brief   Represents a dynamically allocated mathematical vector structure.
 */
typedef struct
{
    size_t length;
    double* values;
} Vector;

/**
 * @brief   Represents a dynamically allocated mathematical matrix structure.
 */
typedef struct
{
    size_t number_of_rows;
    size_t number_of_columns;
    double* values;
} Matrix;

// Function prototypes.
Vector* initialize_vector(size_t length);
Matrix* initialize_matrix(size_t number_of_rows, size_t number_of_columns);

bool solve_linear_system_gaussian_elimination(Matrix* coefficient_matrix, Vector* right_hand_side_vector, Vector* solution_vector);
bool execute_primal_dual_interior_point_method(Matrix* matrix, Vector* vector, Vector* objective, Vector* initial_solution, size_t maximum_iterations, double convergence_tolerance);

void erase_vector(Vector* pointer);
void erase_matrix(Matrix* matrix);

#endif
