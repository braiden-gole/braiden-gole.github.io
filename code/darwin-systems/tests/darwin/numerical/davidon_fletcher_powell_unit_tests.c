/**
 * @file        :   davidon_fletcher_powell_unit_tests.c
 * @author      :   Braiden Gole
 * @brief       :   Comprehensive unit test suite for the Davidon-Fletcher-Powell (DFP) optimization module.
 * @details     :   Utilizes the Unity Test Framework to validate vector normalization, dot product 
 *                  computations, Rosenbrock objective and gradient evaluations, Armijo line search 
 *                  behavior, defensive pointer safety checks, and global minimum convergence accuracy.
 * Compile:
 * clang -DUNITY_INCLUDE_DOUBLE ./tests/darwin/numerical/davidon_fletcher_powell_unit_tests.c ./src/darwin/numerical/davidon-fletcher-powell/davidon_fletcher_powell.c ./src/unity.c -Iinclude -lm -o ./bin/davidon_fletcher_powell_unit_test
 * @version     :   1.0
 * @date        :   2026-11-11
 */
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#include "unity.h"
#include "darwin/numerical/davidon-fletcher-powell/davidon_fletcher_powell.h"

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
 * @brief Ensures passing a null vector pointer to normalization returns zero safely.
 */
void test_calculate_vector_normalization_null_pointer_returns_zero(void)
{
    // Arrange.
    const double* null_vector_pointer = NULL;
    size_t space_dimension = 3;

    // Act.
    double normalization_result = calculate_vector_normalization(null_vector_pointer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, normalization_result);
}

/**
 * @brief Ensures passing a zero space dimension to normalization returns zero safely.
 */
void test_calculate_vector_normalization_zero_dimension_returns_zero(void)
{
    // Arrange.
    double input_vector_buffer[3] = {1.0, 2.0, 3.0};
    size_t zero_dimension = 0;

    // Act.
    double normalization_result = calculate_vector_normalization(input_vector_buffer, zero_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, normalization_result);
}

/**
 * @brief Validates Euclidean norm computation for a unit vector.
 */
void test_calculate_vector_normalization_unit_vector_returns_one(void)
{
    // Arrange.
    double unit_vector_buffer[3] = {0.0, 1.0, 0.0};
    size_t space_dimension = 3;

    // Act.
    double normalization_result = calculate_vector_normalization(unit_vector_buffer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, normalization_result);
}

/**
 * @brief Validates Euclidean norm computation for a multi-dimensional vector.
 */
void test_calculate_vector_normalization_multidimensional_vector(void)
{
    // Arrange.
    double input_vector_buffer[3] = {3.0, 4.0, 0.0};
    size_t space_dimension = 3;

    // Act.
    double normalization_result = calculate_vector_normalization(input_vector_buffer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 5.0, normalization_result);
}

/**
 * @brief Ensures dot product returns zero if the left operand vector is null.
 */
void test_calculate_vector_dot_product_left_null_returns_zero(void)
{
    // Arrange.
    const double* null_vector_pointer = NULL;
    double right_vector_buffer[2] = {1.0, 2.0};

    size_t space_dimension = 2;

    // Act.
    double dot_product_result = calculate_vector_dot_product(null_vector_pointer, right_vector_buffer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Ensures dot product returns zero if the right operand vector is null.
 */
void test_calculate_vector_dot_product_right_null_returns_zero(void)
{
    // Arrange.
    double left_vector_buffer[2] = {1.0, 2.0};
    const double* null_vector_pointer = NULL;

    size_t space_dimension = 2;

    // Act.
    double dot_product_result = calculate_vector_dot_product(left_vector_buffer, null_vector_pointer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Ensures dot product returns zero if space dimension is zero.
 */
void test_calculate_vector_dot_product_zero_dimension_returns_zero(void)
{
    // Arrange.
    double left_vector_buffer[2] = {1.0, 2.0};
    double right_vector_buffer[2] = {3.0, 4.0};

    size_t zero_dimension = 0;

    // Act.
    double dot_product_result = calculate_vector_dot_product(left_vector_buffer, right_vector_buffer, zero_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Validates dot product of orthogonal vectors evaluates to zero.
 */
void test_calculate_vector_dot_product_orthogonal_vectors_returns_zero(void)
{
    // Arrange.
    double left_vector_buffer[2] = {1.0, 0.0};
    double right_vector_buffer[2] = {0.0, 1.0};

    size_t space_dimension = 2;

    // Act.
    double dot_product_result = calculate_vector_dot_product(left_vector_buffer, right_vector_buffer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, dot_product_result);
}

/**
 * @brief Validates dot product calculation for general parallel vectors.
 */
void test_calculate_vector_dot_product_parallel_vectors(void)
{
    // Arrange.
    double left_vector_buffer[2] = {2.0, 3.0};
    double right_vector_buffer[2] = {4.0, 5.0};

    size_t space_dimension = 2;

    // Act.
    double dot_product_result = calculate_vector_dot_product(left_vector_buffer, right_vector_buffer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 23.0, dot_product_result);
}

/**
 * @brief Ensures Rosenbrock objective evaluation returns zero when pointer is null.
 */
void test_evaluate_rosenbrock_objective_null_pointer_returns_zero(void)
{
    // Arrange.
    const double* null_vector_pointer = NULL;

    size_t space_dimension = 2;

    // Act.
    double objective_result = evaluate_rosenbrock_objective(null_vector_pointer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, objective_result);
}

/**
 * @brief Ensures Rosenbrock objective evaluation returns zero when dimension is less than two.
 */
void test_evaluate_rosenbrock_objective_insufficient_dimension_returns_zero(void)
{
    // Arrange.
    double evaluation_vector_buffer[1] = {1.0};

    size_t insufficient_dimension = 1;

    // Act.
    double objective_result = evaluate_rosenbrock_objective(evaluation_vector_buffer, insufficient_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, objective_result);
}

/**
 * @brief Validates Rosenbrock objective function evaluates to zero at the global minimum (1.0, 1.0).
 */
void test_evaluate_rosenbrock_objective_at_global_minimum_returns_zero(void)
{
    // Arrange.
    double minimum_coordinates[2] = {1.0, 1.0};

    size_t space_dimension = 2;

    // Act.
    double objective_result = evaluate_rosenbrock_objective(minimum_coordinates, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, objective_result);
}

/**
 * @brief Validates Rosenbrock objective cost calculation at a known arbitrary coordinate point.
 */
void test_evaluate_rosenbrock_objective_at_known_point(void)
{
    // Arrange.
    double coordinate_buffer[2] = {0.0, 0.0};

    size_t space_dimension = 2;

    // Act.
    double objective_result = evaluate_rosenbrock_objective(coordinate_buffer, space_dimension);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, objective_result);
}

/**
 * @brief Ensures Rosenbrock gradient returns safely without crashing if input pointer is null.
 */
void test_evaluate_rosenbrock_gradient_null_evaluation_vector_returns_safely(void)
{
    // Arrange.
    const double* null_evaluation_vector = NULL;
    double output_gradient_buffer[2];

    // Act.
    evaluate_rosenbrock_gradient(null_evaluation_vector, output_gradient_buffer);

    // Assert.
    // Test passes if execution completes without segmentation fault.
    TEST_ASSERT_TRUE(true);
}

/**
 * @brief Ensures Rosenbrock gradient returns safely without crashing if output pointer is null.
 */
void test_evaluate_rosenbrock_gradient_null_output_vector_returns_safely(void)
{
    // Arrange.
    double evaluation_vector_buffer[2] = {1.0, 1.0};
    double* null_output_gradient = NULL;

    // Act.
    evaluate_rosenbrock_gradient(evaluation_vector_buffer, null_output_gradient);

    // Assert.
    TEST_ASSERT_TRUE(true);
}

/**
 * @brief Validates Rosenbrock gradient components evaluate to zero at the global minimum (1.0, 1.0).
 */
void test_evaluate_rosenbrock_gradient_at_global_minimum_is_zero(void)
{
    // Arrange.
    double minimum_coordinates[2] = {1.0, 1.0};
    double output_gradient_buffer[2] = {99.0, 99.0};

    // Act.
    evaluate_rosenbrock_gradient(minimum_coordinates, output_gradient_buffer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, output_gradient_buffer[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, output_gradient_buffer[1]);
}

/**
 * @brief Validates Rosenbrock gradient calculation at a known starting coordinate point.
 */
void test_evaluate_rosenbrock_gradient_at_known_point(void)
{
    // Arrange.
    double coordinate_buffer[2] = {0.0, 0.0};
    double output_gradient_buffer[2] = {0.0, 0.0};

    // Act.
    evaluate_rosenbrock_gradient(coordinate_buffer, output_gradient_buffer);

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, -2.0, output_gradient_buffer[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, output_gradient_buffer[1]);
}

/**
 * @brief Ensures Armijo line search returns zero when objective function pointer is null.
 */
void test_perform_armijo_line_search_null_objective_returns_zero(void)
{
    // Arrange.
    ObjectiveFunction null_objective = NULL;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double current_position[2] = {0.0, 0.0};
    double search_direction[2] = {-1.0, 0.0};
    double gradient_vector[2] = {-2.0, 0.0};

    size_t space_dimension = 2;

    // Act.
    double step_size_result = perform_armijo_line_search(
        null_objective, valid_gradient, current_position, search_direction,
        gradient_vector, space_dimension, 1.0, 0.5, 1e-4
    );

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, step_size_result);
}

/**
 * @brief Ensures Armijo line search returns zero when contraction factor is outside valid bounds.
 */
void test_perform_armijo_line_search_invalid_contraction_factor_returns_zero(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double current_position[2] = {0.0, 0.0};
    double search_direction[2] = {-1.0, 0.0};
    double gradient_vector[2] = {-2.0, 0.0};

    size_t space_dimension = 2;

    double invalid_contraction_factor = 1.5;

    // Act.
    double step_size_result = perform_armijo_line_search(
        valid_objective, valid_gradient, current_position, search_direction,
        gradient_vector, space_dimension, 1.0, invalid_contraction_factor, 1e-4
    );

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, step_size_result);
}

/**
 * @brief Ensures Armijo line search aborts and returns zero if search direction is not a descent direction.
 */
void test_perform_armijo_line_search_non_descent_direction_returns_zero(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double current_position[2] = {0.0, 0.0};
    double search_direction[2] = {-1.0, 0.0};
    double gradient_vector[2] = {-2.0, 0.0};
    
    size_t space_dimension = 2;

    // Act.
    double step_size_result = perform_armijo_line_search(
        valid_objective, valid_gradient, current_position, search_direction,
        gradient_vector, space_dimension, 1.0, 0.5, 1e-4
    );

    // Assert.
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 0.0, step_size_result);
}

/**
 * @brief Validates successful step size scaling determination under valid descent conditions.
 */
void test_perform_armijo_line_search_valid_descent_step_succeeds(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double current_position[2] = {0.0, 0.0};
    double search_direction[2] = {1.0, 2.0};
    double gradient_vector[2] = {-2.0, 0.0};

    size_t space_dimension = 2;

    // Act.
    double step_size_result = perform_armijo_line_search(
        valid_objective, valid_gradient, current_position, search_direction,
        gradient_vector, space_dimension, 1.0, 0.5, 1e-4
    );

    // Assert.
    TEST_ASSERT_TRUE(step_size_result > 0.0);
    TEST_ASSERT_TRUE(step_size_result <= 1.0);
}

/**
 * @brief Ensures DFP optimizer returns false if objective function pointer is null.
 */
void test_optimize_dfp_null_objective_returns_false(void)
{
    // Arrange.
    ObjectiveFunction null_objective = NULL;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {-1.2, 1.0};

    size_t space_dimension = 2;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        null_objective, valid_gradient, solution_vector, space_dimension, 1000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_FALSE(optimization_status);
}

/**
 * @brief Ensures DFP optimizer returns false if gradient function pointer is null.
 */
void test_optimize_dfp_null_gradient_returns_false(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction null_gradient = NULL;

    double solution_vector[2] = {-1.2, 1.0};

    size_t space_dimension = 2;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, null_gradient, solution_vector, space_dimension, 1000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_FALSE(optimization_status);
}

/**
 * @brief Ensures DFP optimizer returns false if solution position vector pointer is null.
 */
void test_optimize_dfp_null_solution_vector_returns_false(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double* null_solution_vector = NULL;

    size_t space_dimension = 2;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, null_solution_vector, space_dimension, 1000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_FALSE(optimization_status);
}

/**
 * @brief Ensures DFP optimizer returns false if space dimension is zero.
 */
void test_optimize_dfp_zero_space_dimension_returns_false(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {-1.2, 1.0};

    size_t zero_space_dimension = 0;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, solution_vector, zero_space_dimension, 1000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_FALSE(optimization_status);
}

/**
 * @brief Validates that starting directly at the global minimum converges immediately.
 */
void test_optimize_dfp_already_at_minimum_converges_immediately(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {1.0, 1.0};

    size_t space_dimension = 2;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, solution_vector, space_dimension, 1000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_TRUE(optimization_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-6, 1.0, solution_vector[1]);
}

/**
 * @brief Validates DFP optimization converges successfully on the standard Rosenbrock starting point (-1.2, 1.0).
 */
void test_optimize_dfp_rosenbrock_standard_starting_point_converges_successfully(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {-1.2, 1.0};

    size_t space_dimension = 2;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, solution_vector, space_dimension, 2000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_TRUE(optimization_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 1.0, solution_vector[1]);
}

/**
 * @brief Validates DFP optimization converges successfully from an alternative starting point (0.0, 0.0).
 */
void test_optimize_dfp_rosenbrock_alternative_starting_point_converges(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {0.0, 0.0};

    size_t space_dimension = 2;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, solution_vector, space_dimension, 2000, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_TRUE(optimization_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-4, 1.0, solution_vector[1]);
}

/**
 * @brief Validates behavior when maximum iteration limit is severely restricted.
 */
void test_optimize_dfp_exceeded_maximum_iterations_behavior(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {-1.2, 1.0};

    size_t space_dimension = 2;
    size_t restricted_max_iterations = 1;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, solution_vector, space_dimension, restricted_max_iterations, 1e-6, 1e-6
    );

    // Assert.
    TEST_ASSERT_FALSE(optimization_status);
}

/**
 * @brief Validates convergence accuracy under extremely tight gradient and parameter tolerances.
 */
void test_optimize_dfp_tight_tolerances_precision_validation(void)
{
    // Arrange.
    ObjectiveFunction valid_objective = evaluate_rosenbrock_objective;
    GradientFunction valid_gradient = evaluate_rosenbrock_gradient;

    double solution_vector[2] = {-1.0, 1.0};

    size_t space_dimension = 2;

    double tight_tolerance = 1e-8;

    // Act.
    bool optimization_status = optimize_quasi_newton_dfp(
        valid_objective, valid_gradient, solution_vector, space_dimension, 3000, tight_tolerance, tight_tolerance
    );

    // Assert.
    TEST_ASSERT_TRUE(optimization_status);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[0]);
    TEST_ASSERT_DOUBLE_WITHIN(1e-5, 1.0, solution_vector[1]);
}

int main(void)
{
    // Initialize the session.
    UNITY_BEGIN();

    // Register the tests to run.
    RUN_TEST(test_calculate_vector_normalization_null_pointer_returns_zero);
    RUN_TEST(test_calculate_vector_normalization_zero_dimension_returns_zero);
    RUN_TEST(test_calculate_vector_normalization_unit_vector_returns_one);
    RUN_TEST(test_calculate_vector_normalization_multidimensional_vector);
    RUN_TEST(test_calculate_vector_dot_product_left_null_returns_zero);
    RUN_TEST(test_calculate_vector_dot_product_right_null_returns_zero);
    RUN_TEST(test_calculate_vector_dot_product_zero_dimension_returns_zero);
    RUN_TEST(test_calculate_vector_dot_product_orthogonal_vectors_returns_zero);
    RUN_TEST(test_calculate_vector_dot_product_parallel_vectors);
    RUN_TEST(test_evaluate_rosenbrock_objective_null_pointer_returns_zero);
    RUN_TEST(test_evaluate_rosenbrock_objective_insufficient_dimension_returns_zero);
    RUN_TEST(test_evaluate_rosenbrock_objective_at_global_minimum_returns_zero);
    RUN_TEST(test_evaluate_rosenbrock_objective_at_known_point);
    RUN_TEST(test_evaluate_rosenbrock_gradient_null_evaluation_vector_returns_safely);
    RUN_TEST(test_evaluate_rosenbrock_gradient_null_output_vector_returns_safely);
    RUN_TEST(test_evaluate_rosenbrock_gradient_at_global_minimum_is_zero);
    RUN_TEST(test_evaluate_rosenbrock_gradient_at_known_point);
    RUN_TEST(test_perform_armijo_line_search_null_objective_returns_zero);
    RUN_TEST(test_perform_armijo_line_search_invalid_contraction_factor_returns_zero);
    RUN_TEST(test_perform_armijo_line_search_non_descent_direction_returns_zero);
    RUN_TEST(test_perform_armijo_line_search_valid_descent_step_succeeds);
    RUN_TEST(test_optimize_dfp_null_objective_returns_false);
    RUN_TEST(test_optimize_dfp_null_gradient_returns_false);
    RUN_TEST(test_optimize_dfp_null_solution_vector_returns_false);
    RUN_TEST(test_optimize_dfp_zero_space_dimension_returns_false);
    RUN_TEST(test_optimize_dfp_already_at_minimum_converges_immediately);
    RUN_TEST(test_optimize_dfp_rosenbrock_standard_starting_point_converges_successfully);
    RUN_TEST(test_optimize_dfp_rosenbrock_alternative_starting_point_converges);
    RUN_TEST(test_optimize_dfp_exceeded_maximum_iterations_behavior);
    RUN_TEST(test_optimize_dfp_tight_tolerances_precision_validation);

    // Destroy the session.
    return UNITY_END();
}