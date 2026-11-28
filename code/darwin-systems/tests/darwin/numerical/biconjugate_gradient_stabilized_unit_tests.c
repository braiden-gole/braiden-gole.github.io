/**
 * @file        :   biconjugate_gradient_stabilized_unit_tests.c
 * @author      :   Braiden Gole
 * @brief       :   Comprehensive unit test suite for the Biconjugate Gradient Stabilized (BiCGSTAB) linear solver module.
 * @details     :   Utilizes the Unity Test Framework to validate vector normalization, dot product 
 *                  computations, linear system solving accuracy, defensive pointer safety checks, 
 *                  and iteration limit error behaviors.
 * Compile:
 * clang -DUNITY_INCLUDE_DOUBLE ./tests/darwin/numerical/biconjugate_gradient_stabilized_unit_tests.c ./src/darwin/numerical/biconjugate-gradient-stabilized/biconjugate_gradient_stabilized.c ./src/unity.c -Iinclude -lm -o ./bin/biconjugate_gradient_stabilized_unit_test
 * @version     :   1.0
 * @date        :   2026-11-11
 */
#include <stdio.h>

#include "unity.h"
#include "darwin/numerical/biconjugate-gradient-stabilized/biconjugate_gradient_stabilized.h"

/**
 * @brief Sets up resources before each test execution.
 */
void setUp(void)
{
    ;;
}

/**
 * @brief Cleans up resources after each test execution.
 */
void tearDown(void)
{
    ;;
}

/**
 * @brief Helper callback function simulating identity matrix multiplication for testing.
 */
static void identity_matrix_multiply(int number_of_dimensions, const double* vector, double* product_vector, void* context)
{
    // Defend against null pointers.
    if (vector == NULL || product_vector == NULL)
    {
        return;
    }
    
    (void)context;

    // Multiply vector components by identity.
    for (int index = 0; index < number_of_dimensions; index++)
    {
        product_vector[index] = vector[index];
    }
}

/**
 * @brief Helper callback function simulating a diagonal matrix multiplication for testing.
 */
static void diagonal_matrix_multiply(int number_of_dimensions, const double* vector, double* product_vector, void* context)
{
    // Defend against null pointers.
    if (vector == NULL || product_vector == NULL)
    {
        return;
    }

    (void)context;

    // Multiply vector components by diagonal scaling factors.
    for (int index = 0; index < number_of_dimensions; index++)
    {
        product_vector[index] = (double)(index + 2) * vector[index];
    }
}

/**
 * @brief Helper callback function simulating a fractional diagonal matrix multiplication for testing.
 */
static void fractional_diagonal_matrix_multiply(int number_of_dimensions, const double* vector, double* product_vector, void* context)
{
    // Defend against null pointers.
    if (vector == NULL || product_vector == NULL)
    {
        return;
    }

    (void)context;

    // Multiply vector components by fractional diagonal scaling factors.
    for (int index = 0; index < number_of_dimensions; index++)
    {
        product_vector[index] = (0.5 * (double)(index + 1)) * vector[index];
    }
}

/**
 * @brief Helper callback function simulating a tridiagonal matrix multiplication for testing.
 */
static void tridiagonal_matrix_multiply(int number_of_dimensions, const double* vector, double* product_vector, void* context)
{
    // Defend against null pointers or invalid dimensions.
    if (vector == NULL || product_vector == NULL || number_of_dimensions != 3)
    {
        return;
    }

    (void)context;

    // Perform row operations for the 3x3 tridiagonal matrix multiplication.
    product_vector[0] = 4.0 * vector[0] + 1.0 * vector[1];
    product_vector[1] = 1.0 * vector[0] + 4.0 * vector[1] + 1.0 * vector[2];
    product_vector[2] = 1.0 * vector[1] + 4.0 * vector[2];
}

/**
 * @brief Helper callback function simulating a non-symmetric 4x4 convection-diffusion matrix multiplication.
 */
static void nonsymmetric_matrix_multiply(int number_of_dimensions, const double* vector, double* product_vector, void* context)
{
    // Defend against null pointers or invalid dimensions.
    if (vector == NULL || product_vector == NULL || number_of_dimensions != 4)
    {
        return;
    }

    // Suppress unused context warning.
    (void)context;

    // Perform row operations for the non-symmetric 4x4 test matrix.
    product_vector[0] =  4.0 * vector[0] - 1.0 * vector[1];
    product_vector[1] = -2.0 * vector[0] + 4.0 * vector[1] - 1.0 * vector[2];
    product_vector[2] = -2.0 * vector[1] + 4.0 * vector[2] - 1.0 * vector[3];
    product_vector[3] = -2.0 * vector[2] + 4.0 * vector[3];
}

/**
 * @brief Ensures passing a null vector pointer to normalization returns zero safely.
 */
void test_compute_vector_normalization_null_pointer_returns_zero(void)
{
    // Arrange.
    const double* null_vector_pointer = NULL;

    int number_of_dimensions = 3;

    // Act.
    double normalization_result = compute_vector_normalization(number_of_dimensions, null_vector_pointer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, normalization_result);
}

/**
 * @brief Ensures passing a zero dimension to normalization returns zero safely.
 */
void test_compute_vector_normalization_zero_dimension_returns_zero(void)
{
    // Arrange.
    double vector_buffer[3] = {1.0, 2.0, 3.0};

    int zero_dimension = 0;

    // Act.
    double normalization_result = compute_vector_normalization(zero_dimension, vector_buffer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, normalization_result);
}

/**
 * @brief Ensures passing a negative dimension to normalization returns zero safely.
 */
void test_compute_vector_normalization_negative_dimension_returns_zero(void)
{
    // Arrange.
    double vector_buffer[3] = {1.0, 2.0, 3.0};

    int negative_dimension = -2;

    // Act.
    double normalization_result = compute_vector_normalization(negative_dimension, vector_buffer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, normalization_result);
}

/**
 * @brief Validates Euclidean norm computation for a unit vector.
 */
void test_compute_vector_normalization_unit_vector_returns_one(void)
{
    // Arrange.
    double unit_vector_buffer[3] = {0.0, 1.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    double normalization_result = compute_vector_normalization(number_of_dimensions, unit_vector_buffer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, normalization_result);
}

/**
 * @brief Validates Euclidean norm computation for a multi-dimensional vector.
 */
void test_compute_vector_normalization_multidimensional_vector(void)
{
    // Arrange.
    double input_vector_buffer[3] = {3.0, 4.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    double normalization_result = compute_vector_normalization(number_of_dimensions, input_vector_buffer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 5.0, normalization_result);
}

/**
 * @brief Ensures dot product returns zero if the left operand vector is null.
 */
void test_compute_vector_dot_product_left_null_returns_zero(void)
{
    // Arrange.
    const double* null_vector_pointer = NULL;

    double operand[2] = {1.0, 2.0};

    int number_of_dimensions = 2;

    // Act.
    double dot_product_result = compute_vector_dot_product(number_of_dimensions, null_vector_pointer, operand);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Ensures dot product returns zero if the right operand vector is null.
 */
void test_compute_vector_dot_product_right_null_returns_zero(void)
{
    // Arrange.
    double vector[2] = {1.0, 2.0};

    const double* null_vector_pointer = NULL;

    int number_of_dimensions = 2;

    // Act.
    double dot_product_result = compute_vector_dot_product(number_of_dimensions, vector, null_vector_pointer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Ensures dot product returns zero if dimension is zero.
 */
void test_compute_vector_dot_product_zero_dimension_returns_zero(void)
{
    // Arrange.
    double vector[2] = {1.0, 2.0};
    double operand[2] = {3.0, 4.0};

    int zero_dimension = 0;

    // Act.
    double dot_product_result = compute_vector_dot_product(zero_dimension, vector, operand);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Ensures dot product returns zero if dimension is negative.
 */
void test_compute_vector_dot_product_negative_dimension_returns_zero(void)
{
    // Arrange.
    double vector[2] = {1.0, 2.0};
    double operand[2] = {3.0, 4.0};

    int negative_dimension = -4;

    // Act.
    double dot_product_result = compute_vector_dot_product(negative_dimension, vector, operand);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Validates dot product of orthogonal vectors evaluates to zero.
 */
void test_compute_vector_dot_product_orthogonal_vectors_returns_zero(void)
{
    // Arrange.
    double vector[2] = {1.0, 0.0};
    double operand[2] = {0.0, 1.0};

    int number_of_dimensions = 2;

    // Act.
    double dot_product_result = compute_vector_dot_product(number_of_dimensions, vector, operand);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Validates dot product calculation for general parallel vectors.
 */
void test_compute_vector_dot_product_parallel_vectors(void)
{
    // Arrange.
    double vector[2] = {2.0, 3.0};
    double operand[2] = {4.0, 5.0};

    int number_of_dimensions = 2;

    // Act.
    double dot_product_result = compute_vector_dot_product(number_of_dimensions, vector, operand);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 23.0, dot_product_result);
}

/**
 * @brief Ensures solver returns error null pointer when callback function is null.
 */
void test_solve_null_callback_returns_error(void)
{
    // Arrange.
    MultiplyFunction null_callback = NULL;

    double operand[2] = {1.0, 1.0};
    double initial_guess[2] = {0.0, 0.0};
    double solution_vector[2] = {0.0, 0.0};

    int number_of_dimensions = 2;

    // Act.
    Status status = solve(null_callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_NULL_POINTER, status);
}

/**
 * @brief Ensures solver returns error null pointer when initial guess is null.
 */
void test_solve_null_initial_guess_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[2] = {1.0, 1.0};

    double* null_initial_guess = NULL;

    double solution_vector[2] = {0.0, 0.0};

    int number_of_dimensions = 2;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, null_initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_NULL_POINTER, status);
}

/**
 * @brief Ensures solver returns error null pointer when operand vector is null.
 */
void test_solve_null_operand_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    const double* null_operand = NULL;

    double initial_guess[2] = {0.0, 0.0};
    double solution_vector[2] = {0.0, 0.0};

    int number_of_dimensions = 2;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, null_operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_NULL_POINTER, status);
}

/**
 * @brief Ensures solver returns error null pointer when solution vector buffer is null.
 */
void test_solve_null_solution_vector_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[2] = {1.0, 1.0};
    double initial_guess[2] = {0.0, 0.0};

    double* null_solution_vector = NULL;

    int number_of_dimensions = 2;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, null_solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_NULL_POINTER, status);
}

/**
 * @brief Ensures solver returns invalid dimension error when dimension count is zero.
 */
void test_solve_zero_dimension_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[1] = {1.0};
    double initial_guess[1] = {0.0};
    double solution_vector[1] = {0.0};

    int zero_dimension = 0;

    // Act.
    Status status = solve(callback, zero_dimension, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_INVALID_DIMENSION, status);
}

/**
 * @brief Ensures solver returns invalid dimension error when dimension count is negative.
 */
void test_solve_negative_dimension_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[1] = {1.0};
    double initial_guess[1] = {0.0};
    double solution_vector[1] = {0.0};

    int negative_dimension = -3;

    // Act.
    Status status = solve(callback, negative_dimension, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_INVALID_DIMENSION, status);
}

/**
 * @brief Ensures solver returns invalid dimension error when maximum iterations is zero.
 */
void test_solve_zero_max_iterations_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[2] = {1.0, 1.0};
    double initial_guess[2] = {0.0, 0.0};
    double solution_vector[2] = {0.0, 0.0};

    int number_of_dimensions = 2;
    int zero_iterations = 0;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, zero_iterations, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_INVALID_DIMENSION, status);
}

/**
 * @brief Ensures solver returns invalid dimension error when convergence tolerance is negative.
 */
void test_solve_negative_tolerance_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[2] = {1.0, 1.0};
    double initial_guess[2] = {0.0, 0.0};
    double solution_vector[2] = {0.0, 0.0};

    int number_of_dimensions = 2;
    double negative_tolerance = -1e-4;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, negative_tolerance);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_INVALID_DIMENSION, status);
}

/**
 * @brief Validates solver returns success immediately if initial guess is already the exact solution.
 */
void test_solve_initial_guess_already_correct_succeeds_immediately(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[2] = {3.0, 5.0};
    double initial_guess[2] = {3.0, 5.0};
    double solution_vector[2] = {0.0, 0.0};

    int number_of_dimensions = 2;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 3.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 5.0, solution_vector[1]);
}

/**
 * @brief Validates BiCGSTAB correctly solves an identity matrix system Ax = b.
 */
void test_solve_identity_matrix_system(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[3] = {2.0, -4.0, 7.0};
    double initial_guess[3] = {0.0, 0.0, 0.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 2.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, -4.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 7.0, solution_vector[2]);
}

/**
 * @brief Validates BiCGSTAB correctly solves a 5-dimensional identity matrix system.
 */
void test_solve_five_dimensional_identity_system(void)
{
    // Arrange.
    MultiplyFunction callback = identity_matrix_multiply;

    double operand[5] = {1.0, 2.0, 3.0, 4.0, 5.0};
    double initial_guess[5] = {0.0, 0.0, 0.0, 0.0, 0.0};
    double solution_vector[5] = {0.0, 0.0, 0.0, 0.0, 0.0};

    int number_of_dimensions = 5;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);

    for (int index = 0; index < 5; index++)
    {
        TEST_ASSERT_DOUBLE_WITHIN(1e-6, operand[index], solution_vector[index]);
    }
}

/**
 * @brief Validates BiCGSTAB correctly solves a diagonal matrix system.
 */
void test_solve_diagonal_matrix_system(void)
{
    // Arrange.
    MultiplyFunction callback = diagonal_matrix_multiply;

    double operand[3] = {4.0, 9.0, 16.0};
    double initial_guess[3] = {0.0, 0.0, 0.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 2.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 3.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 4.0, solution_vector[2]);
}

/**
 * @brief Validates BiCGSTAB correctly solves a fractional diagonal matrix system.
 */
void test_solve_fractional_diagonal_matrix_system(void)
{
    // Arrange.
    MultiplyFunction callback = fractional_diagonal_matrix_multiply;

    double operand[3] = {0.5, 2.0, 4.5};
    double initial_guess[3] = {0.0, 0.0, 0.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 2.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 3.0, solution_vector[2]);
}

/**
 * @brief Validates BiCGSTAB correctly solves a standard tridiagonal matrix system.
 */
void test_solve_tridiagonal_matrix_system_succeeds(void)
{
    // Arrange.
    MultiplyFunction callback = tridiagonal_matrix_multiply;

    double operand[3] = {5.0, 6.0, 5.0};
    double initial_guess[3] = {0.0, 0.0, 0.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 100, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[2]);
}

/**
 * @brief Validates BiCGSTAB converges starting from a non-zero random initial guess vector.
 */
void test_solve_random_initial_guess_tridiagonal_system(void)
{
    // Arrange.
    MultiplyFunction callback = tridiagonal_matrix_multiply;

    double operand[3] = {5.0, 6.0, 5.0};
    double initial_guess[3] = {0.5, -0.5, 2.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 200, 1e-6);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[2]);
}

/**
 * @brief Validates behavior when maximum iteration limit is severely restricted on a non-trivial system.
 */
void test_solve_exceeded_maximum_iterations_returns_error(void)
{
    // Arrange.
    MultiplyFunction callback = tridiagonal_matrix_multiply;

    double operand[3] = {5.0, 6.0, 5.0};
    double initial_guess[3] = {0.0, 0.0, 0.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;
    int restricted_iterations = 1;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, restricted_iterations, 1e-12);

    // Assert.
    TEST_ASSERT_EQUAL(ERROR_MAXIMUM_ITERATIONS_EXCEEDED, status);
}

/**
 * @brief Validates solver precision under extremely tight convergence tolerances.
 */
void test_solve_tight_tolerance_precision_validation(void)
{
    // Arrange.
    MultiplyFunction callback = tridiagonal_matrix_multiply;

    double operand[3] = {5.0, 6.0, 5.0};
    double initial_guess[3] = {0.5, 0.5, 0.5};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    double tight_tolerance = 1e-10;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 500, tight_tolerance);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-8, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-8, 1.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-8, 1.0, solution_vector[2]);
}

/**
 * @brief Validates rapid termination under a loose convergence tolerance threshold.
 */
void test_solve_loose_tolerance_fast_convergence(void)
{
    // Arrange.
    MultiplyFunction callback = tridiagonal_matrix_multiply;

    double operand[3] = {5.0, 6.0, 5.0};
    double initial_guess[3] = {0.0, 0.0, 0.0};
    double solution_vector[3] = {0.0, 0.0, 0.0};

    int number_of_dimensions = 3;

    double loose_tolerance = 1e-2;

    // Act.
    Status status = solve(callback, number_of_dimensions, NULL, operand, initial_guess, solution_vector, 50, loose_tolerance);

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-2, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-2, 1.0, solution_vector[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-2, 1.0, solution_vector[2]);
}

/**
 * @brief Comprehensive stress test validating BiCGSTAB convergence on a challenging 4x4 non-symmetric linear system.
 */
void test_solve_nonsymmetric_convection_diffusion_system_converges_precisely(void)
{
    // Arrange.
    MultiplyFunction callback = nonsymmetric_matrix_multiply;

    int number_of_dimensions = 4;
    
    double operand[4] = {2.0, 3.0, 4.0, 10.0};
    double initial_guess[4] = {0.0, 0.0, 0.0, 0.0};
    double computed_solution[4] = {0.0, 0.0, 0.0, 0.0};
    
    int maximum_allowed_iterations = 200;

    double target_convergence_tolerance = 1e-12;

    // Act.
    Status execution_status = solve(
        callback,
        number_of_dimensions,
        NULL,
        operand,
        initial_guess,
        computed_solution,
        maximum_allowed_iterations,
        target_convergence_tolerance
    );

    // Assert.
    TEST_ASSERT_EQUAL(SUCCESS, execution_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-9, 1.0, computed_solution[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-9, 2.0, computed_solution[1]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-9, 3.0, computed_solution[2]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-9, 4.0, computed_solution[3]);
}

int main(void)
{
    // Initialize the session.
    UNITY_BEGIN();

    // Register the tests to run.
    RUN_TEST(test_compute_vector_normalization_null_pointer_returns_zero);
    RUN_TEST(test_compute_vector_normalization_zero_dimension_returns_zero);
    RUN_TEST(test_compute_vector_normalization_negative_dimension_returns_zero);
    RUN_TEST(test_compute_vector_normalization_unit_vector_returns_one);
    RUN_TEST(test_compute_vector_normalization_multidimensional_vector);
    RUN_TEST(test_compute_vector_dot_product_left_null_returns_zero);
    RUN_TEST(test_compute_vector_dot_product_right_null_returns_zero);
    RUN_TEST(test_compute_vector_dot_product_zero_dimension_returns_zero);
    RUN_TEST(test_compute_vector_dot_product_negative_dimension_returns_zero);
    RUN_TEST(test_compute_vector_dot_product_orthogonal_vectors_returns_zero);
    RUN_TEST(test_compute_vector_dot_product_parallel_vectors);
    RUN_TEST(test_solve_null_callback_returns_error);
    RUN_TEST(test_solve_null_initial_guess_returns_error);
    RUN_TEST(test_solve_null_operand_returns_error);
    RUN_TEST(test_solve_null_solution_vector_returns_error);
    RUN_TEST(test_solve_zero_dimension_returns_error);
    RUN_TEST(test_solve_negative_dimension_returns_error);
    RUN_TEST(test_solve_zero_max_iterations_returns_error);
    RUN_TEST(test_solve_negative_tolerance_returns_error);
    RUN_TEST(test_solve_initial_guess_already_correct_succeeds_immediately);
    RUN_TEST(test_solve_identity_matrix_system);
    RUN_TEST(test_solve_five_dimensional_identity_system);
    RUN_TEST(test_solve_diagonal_matrix_system);
    RUN_TEST(test_solve_fractional_diagonal_matrix_system);
    RUN_TEST(test_solve_tridiagonal_matrix_system_succeeds);
    RUN_TEST(test_solve_random_initial_guess_tridiagonal_system);
    RUN_TEST(test_solve_exceeded_maximum_iterations_returns_error);
    RUN_TEST(test_solve_tight_tolerance_precision_validation);
    RUN_TEST(test_solve_loose_tolerance_fast_convergence);
    RUN_TEST(test_solve_nonsymmetric_convection_diffusion_system_converges_precisely);

    // Destroy the session.
    return UNITY_END();
}