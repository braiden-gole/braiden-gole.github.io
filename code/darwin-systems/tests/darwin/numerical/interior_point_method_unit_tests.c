/**
 * @file        :   interior_point_method_unit_tests.c
 * @author      :   Braiden Gole
 * @brief       :   Comprehensive unit test suite for the Primal-Dual Path-Following Interior Point Method module.
 * @details     :   Utilizes the Unity Test Framework to validate the correctness, error handling, 
 *                  memory allocation boundary checks, Gaussian elimination solver accuracy, 
 *                  and mathematical optimization convergence of the interior point linear programming solver.
 * Compile:
 * clang -DUNITY_INCLUDE_DOUBLE ./tests/darwin/numerical/interior_point_method_unit_tests.c ./src/darwin/numerical/primary-dual-path-following-interior-point-method/interior_point_method.h.c ./src/unity.c -Iinclude -lm -o ./bin/interior_point_method_unit_test
 * @version     :   1.0
 * @date        :   2026-11-12
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>

#include "unity.h"
#include "darwin/numerical/primary-dual-path-following-interior-point-method/interior_point_method.h"

/**
 * @brief Sets up resources before each test.
 */
void setUp(void)
{
    ;;
}

/**
 * @brief Cleans up resources after each test.
 */
void tearDown(void)
{
    ;;
}

/**
 * @brief Validates that initializing a vector with zero length returns a null pointer safely.
 */
void test_vector_initialization_zero_length_returns_null(void)
{
    // Arrange.
    size_t length = 0;

    // Act.
    Vector* vector = initialize_vector(length);

    // Assert.
    TEST_ASSERT_NULL(vector);
}

/**
 * @brief Validates that initializing a vector with excessive length exceeding memory bounds returns null.
 */
void test_vector_initialization_excessive_length_returns_null(void)
{
    // Arrange.
    size_t excessive_length = ((size_t)1 << 60);

    // Act.
    Vector* vector = initialize_vector(excessive_length);

    // Assert.
    TEST_ASSERT_NULL(vector);
}

/**
 * @brief Verifies that initializing a vector with a valid length succeeds and allocates memory.
 */
void test_vector_initialization_valid_length_succeeds(void)
{
    // Arrange.
    size_t length = 5;

    // Act.
    Vector* vector = initialize_vector(length);

    // Assert.
    TEST_ASSERT_NOT_NULL(vector);
    TEST_ASSERT_EQUAL_INT(length, vector->length);
    TEST_ASSERT_NOT_NULL(vector->values);

    erase_vector(vector);
}

/**
 * @brief Validates that initializing a matrix with zero rows returns a null pointer safely.
 */
void test_matrix_initialization_zero_rows_returns_null(void)
{
    // Arrange.
    size_t rows = 0;
    size_t columns = 3;

    // Act.
    Matrix* matrix = initialize_matrix(rows, columns);

    // Assert.
    TEST_ASSERT_NULL(matrix);
}

/**
 * @brief Validates that initializing a matrix with zero columns returns a null pointer safely.
 */
void test_matrix_initialization_zero_columns_returns_null(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 0;

    // Act.
    Matrix* matrix = initialize_matrix(rows, columns);

    // Assert.
    TEST_ASSERT_NULL(matrix);
}

/**
 * @brief Verifies that initializing a matrix with valid dimensions succeeds.
 */
void test_matrix_initialization_valid_dimensions_succeeds(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    // Act.
    Matrix* matrix = initialize_matrix(rows, columns);

    // Assert.
    TEST_ASSERT_NOT_NULL(matrix);
    TEST_ASSERT_EQUAL_INT(rows, matrix->number_of_rows);
    TEST_ASSERT_EQUAL_INT(columns, matrix->number_of_columns);
    TEST_ASSERT_NOT_NULL(matrix->values);

    erase_matrix(matrix);
}

/**
 * @brief Ensures erasing a null vector pointer executes safely without crashing.
 */
void test_erase_vector_null_pointer_safe(void)
{
    // Arrange.
    Vector* null_vector = NULL;

    // Act & Assert.
    erase_vector(null_vector);

    TEST_ASSERT_NULL(null_vector);
}

/**
 * @brief Ensures erasing a null matrix pointer executes safely without crashing.
 */
void test_erase_matrix_null_pointer_safe(void)
{
    // Arrange.
    Matrix* null_matrix = NULL;

    // Act & Assert.
    erase_matrix(null_matrix);
    
    TEST_ASSERT_NULL(null_matrix);
}

/**
 * @brief Validates that Gaussian elimination returns false when given null pointers.
 */
void test_gaussian_elimination_null_pointers_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 2);
    Vector* vector = initialize_vector(2);
    Vector* solution = initialize_vector(2);

    // Act & Assert.
    TEST_ASSERT_FALSE(solve_linear_system_gaussian_elimination(NULL, vector, solution));
    TEST_ASSERT_FALSE(solve_linear_system_gaussian_elimination(matrix, NULL, solution));
    TEST_ASSERT_FALSE(solve_linear_system_gaussian_elimination(matrix, vector, NULL));

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(solution);
}

/**
 * @brief Validates that dimension mismatches in Gaussian elimination return false.
 */
void test_gaussian_elimination_dimension_mismatch_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 2);
    Vector* vector = initialize_vector(3);
    Vector* solution = initialize_vector(2);

    // Act.
    bool status = solve_linear_system_gaussian_elimination(matrix, vector, solution);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(solution);
}

/**
 * @brief Verifies that solving a singular matrix system via Gaussian elimination returns false.
 */
void test_gaussian_elimination_singular_matrix_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 2);

    matrix->values[0] = 1.0; matrix->values[1] = 2.0;
    matrix->values[2] = 2.0; matrix->values[3] = 4.0;

    Vector* vector = initialize_vector(2);

    vector->values[0] = 1.0;
    vector->values[1] = 2.0;

    Vector* solution = initialize_vector(2);

    // Act.
    bool status = solve_linear_system_gaussian_elimination(matrix, vector, solution);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(solution);
}

/**
 * @brief Validates successful resolution of a valid 2x2 linear system using Gaussian elimination.
 */
void test_gaussian_elimination_valid_2x2_system_succeeds(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 2);

    matrix->values[0] = 2.0; matrix->values[1] = 1.0;
    matrix->values[2] = 1.0; matrix->values[3] = 3.0;

    Vector* vector = initialize_vector(2);

    vector->values[0] = 5.0;
    vector->values[1] = 5.0;

    Vector* solution = initialize_vector(2);

    // Act.
    bool status = solve_linear_system_gaussian_elimination(matrix, vector, solution);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 2.0, solution->values[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution->values[1]);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(solution);
}

/**
 * @brief Validates successful resolution of a valid 3x3 linear system using Gaussian elimination.
 */
void test_gaussian_elimination_valid_3x3_system_succeeds(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(3, 3);
    
    matrix->values[0] = 4.0; matrix->values[1] = 1.0; matrix->values[2] = 0.0;
    matrix->values[3] = 1.0; matrix->values[4] = 4.0; matrix->values[5] = 1.0;
    matrix->values[6] = 0.0; matrix->values[7] = 1.0; matrix->values[8] = 3.0;

    Vector* vector = initialize_vector(3);

    vector->values[0] = 5.0;
    vector->values[1] = 6.0;
    vector->values[2] = 4.0;

    Vector* solution = initialize_vector(3);

    // Act.
    bool status = solve_linear_system_gaussian_elimination(matrix, vector, solution);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution->values[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution->values[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution->values[2]);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(solution);
}

/**
 * @brief Ensures passing a null constraint matrix pointer to the interior point method returns false.
 */
void test_interior_point_null_matrix_pointer_returns_false(void)
{
    // Arrange.
    Vector* vector = initialize_vector(2);
    Vector* objective = initialize_vector(3);
    Vector* solution = initialize_vector(3);

    // Act.
    bool status = execute_primal_dual_interior_point_method(NULL, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Ensures passing a null constraint bound vector to the interior point method returns false.
 */
void test_interior_point_null_vector_pointer_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 3);
    Vector* objective = initialize_vector(3);
    Vector* solution = initialize_vector(3);

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, NULL, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Ensures passing a null objective vector to the interior point method returns false.
 */
void test_interior_point_null_objective_pointer_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 3);
    Vector* vector = initialize_vector(2);
    Vector* solution = initialize_vector(3);

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, NULL, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(solution);
}

/**
 * @brief Ensures passing a null initial solution vector to the interior point method returns false.
 */
void test_interior_point_null_initial_solution_pointer_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 3);
    Vector* vector = initialize_vector(2);
    Vector* objective = initialize_vector(3);

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, NULL, 100, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
}

/**
 * @brief Validates that dimension mismatches between matrix rows and bound vector return false.
 */
void test_interior_point_dimension_mismatch_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 3);
    Vector* vector = initialize_vector(3);
    Vector* objective = initialize_vector(3);
    Vector* solution = initialize_vector(3);

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Validates that having fewer columns than rows in the constraint matrix returns false.
 */
void test_interior_point_fewer_columns_than_rows_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(3, 2);
    Vector* vector = initialize_vector(3);
    Vector* objective = initialize_vector(2);
    Vector* solution = initialize_vector(2);

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Validates that zero maximum iterations limit returns false.
 */
void test_interior_point_zero_iterations_limit_returns_false(void)
{
    // Arrange.
    Matrix* matrix = initialize_matrix(2, 3);
    Vector* vector = initialize_vector(2);
    Vector* objective = initialize_vector(3);
    Vector* solution = initialize_vector(3);

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 0, 1e-8);

    // Assert.
    TEST_ASSERT_FALSE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Verifies successful optimization convergence on standard 2-constraint 3-variable linear program.
 */
void test_interior_point_standard_2_constraint_3_variable_lp_succeeds(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0; matrix->values[2] = 1.0;
    matrix->values[3] = 2.0; matrix->values[4] = 1.0; matrix->values[5] = 0.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 10.0;
    vector->values[1] = 14.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -1.0;
    objective->values[1] = -2.0;
    objective->values[2] = 0.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 1.0;
    solution->values[1] = 1.0;
    solution->values[2] = 1.0;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 4.0, solution->values[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 6.0, solution->values[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 0.0, solution->values[2]);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Verifies interior point optimization on a single constraint multiple variable system.
 */
void test_interior_point_single_constraint_multiple_variables_succeeds(void)
{
    // Arrange.
    size_t rows = 1;
    size_t columns = 2;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 5.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -3.0;
    objective->values[1] = -1.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 1.0;
    solution->values[1] = 1.0;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 5.0, solution->values[0]);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Validates successful optimization under strict tight convergence tolerance.
 */
void test_interior_point_tight_convergence_tolerance_succeeds(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0; matrix->values[2] = 1.0;
    matrix->values[3] = 2.0; matrix->values[4] = 1.0; matrix->values[5] = 0.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 10.0;
    vector->values[1] = 14.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -1.0;
    objective->values[1] = -2.0;
    objective->values[2] = 0.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 1.0;
    solution->values[1] = 1.0;
    solution->values[2] = 1.0;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 200, 1e-10);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 4.0, solution->values[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 6.0, solution->values[1]);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Validates behavior when initial solution is already at the optimal point.
 */
void test_interior_point_already_optimal_initial_solution_succeeds(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0; matrix->values[2] = 1.0;
    matrix->values[3] = 2.0; matrix->values[4] = 1.0; matrix->values[5] = 0.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 10.0;
    vector->values[1] = 14.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -1.0;
    objective->values[1] = -2.0;
    objective->values[2] = 0.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 4.0;
    solution->values[1] = 6.0;
    solution->values[2] = 0.5;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_TRUE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Validates that negative initial guesses are handled safely and corrected by the interior initialization.
 */
void test_interior_point_negative_initial_guess_handled_safely(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0; matrix->values[2] = 1.0;
    matrix->values[3] = 2.0; matrix->values[4] = 1.0; matrix->values[5] = 0.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 10.0;
    vector->values[1] = 14.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -1.0;
    objective->values[1] = -2.0;
    objective->values[2] = 0.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = -5.0;
    solution->values[1] = -10.0;
    solution->values[2] = -2.0;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_TRUE(status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 4.0, solution->values[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 6.0, solution->values[1]);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Verifies optimization succeeds with zero right-hand side constant vector.
 */
void test_interior_point_zero_right_hand_side_vector_succeeds(void)
{
    // Arrange.
    size_t rows = 1;
    size_t columns = 2;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 0.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -1.0;
    objective->values[1] = -1.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 1.0;
    solution->values[1] = 1.0;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

    // Assert.
    TEST_ASSERT_TRUE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Verifies optimization succeeds when objective function coefficients are all zero.
 */
void test_interior_point_zero_objective_coefficients_succeeds(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0; matrix->values[2] = 1.0;
    matrix->values[3] = 2.0; matrix->values[4] = 1.0; matrix->values[5] = 0.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 10.0;
    vector->values[1] = 14.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = 0.0;
    objective->values[1] = 0.0;
    objective->values[2] = 0.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 1.0;
    solution->values[1] = 1.0;
    solution->values[2] = 1.0;

    // Act.
    bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 50, 1e-8);

    // Assert.
    TEST_ASSERT_TRUE(status);

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

/**
 * @brief Validates memory allocation and cleanup stability across repeated test executions.
 */
void test_interior_point_repeated_execution_stability(void)
{
    // Arrange.
    size_t rows = 2;
    size_t columns = 3;

    Matrix* matrix = initialize_matrix(rows, columns);

    matrix->values[0] = 1.0; matrix->values[1] = 1.0; matrix->values[2] = 1.0;
    matrix->values[3] = 2.0; matrix->values[4] = 1.0; matrix->values[5] = 0.0;

    Vector* vector = initialize_vector(rows);

    vector->values[0] = 10.0;
    vector->values[1] = 14.0;

    Vector* objective = initialize_vector(columns);

    objective->values[0] = -1.0;
    objective->values[1] = -2.0;
    objective->values[2] = 0.0;

    Vector* solution = initialize_vector(columns);

    solution->values[0] = 1.0;
    solution->values[1] = 1.0;
    solution->values[2] = 1.0;

    // Act & Assert.
    for (int run = 0; run < 5; run++)
    {
        solution->values[0] = 1.0;
        solution->values[1] = 1.0;
        solution->values[2] = 1.0;

        bool status = execute_primal_dual_interior_point_method(matrix, vector, objective, solution, 100, 1e-8);

        TEST_ASSERT_TRUE(status);
    }

    erase_matrix(matrix);
    erase_vector(vector);
    erase_vector(objective);
    erase_vector(solution);
}

int main(void)
{
    // Initialize the session.
    UNITY_BEGIN();

    // Register the tests to run.
    RUN_TEST(test_vector_initialization_zero_length_returns_null);
    RUN_TEST(test_vector_initialization_excessive_length_returns_null);
    RUN_TEST(test_vector_initialization_valid_length_succeeds);
    RUN_TEST(test_matrix_initialization_zero_rows_returns_null);
    RUN_TEST(test_matrix_initialization_zero_columns_returns_null);
    RUN_TEST(test_matrix_initialization_valid_dimensions_succeeds);
    RUN_TEST(test_erase_vector_null_pointer_safe);
    RUN_TEST(test_erase_matrix_null_pointer_safe);
    RUN_TEST(test_gaussian_elimination_null_pointers_returns_false);
    RUN_TEST(test_gaussian_elimination_dimension_mismatch_returns_false);
    RUN_TEST(test_gaussian_elimination_singular_matrix_returns_false);
    RUN_TEST(test_gaussian_elimination_valid_2x2_system_succeeds);
    RUN_TEST(test_gaussian_elimination_valid_3x3_system_succeeds);
    RUN_TEST(test_interior_point_null_matrix_pointer_returns_false);
    RUN_TEST(test_interior_point_null_vector_pointer_returns_false);
    RUN_TEST(test_interior_point_null_objective_pointer_returns_false);
    RUN_TEST(test_interior_point_null_initial_solution_pointer_returns_false);
    RUN_TEST(test_interior_point_dimension_mismatch_returns_false);
    RUN_TEST(test_interior_point_fewer_columns_than_rows_returns_false);
    RUN_TEST(test_interior_point_zero_iterations_limit_returns_false);
    RUN_TEST(test_interior_point_standard_2_constraint_3_variable_lp_succeeds);
    RUN_TEST(test_interior_point_single_constraint_multiple_variables_succeeds);
    RUN_TEST(test_interior_point_tight_convergence_tolerance_succeeds);
    RUN_TEST(test_interior_point_already_optimal_initial_solution_succeeds);
    RUN_TEST(test_interior_point_negative_initial_guess_handled_safely);
    RUN_TEST(test_interior_point_zero_right_hand_side_vector_succeeds);
    RUN_TEST(test_interior_point_zero_objective_coefficients_succeeds);
    RUN_TEST(test_vector_initialization_valid_length_succeeds);
    RUN_TEST(test_matrix_initialization_valid_dimensions_succeeds);
    RUN_TEST(test_interior_point_repeated_execution_stability);

    // Destroy the session
    return UNITY_END();
}