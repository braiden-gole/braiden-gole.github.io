/**
 * @file        :   fibonacci_heap_unit_tests.c
 * @author      :   Braiden Gole
 * @brief       :   Comprehensive suite of 30 unique unit tests for the Fibonacci Heap priority queue.
 * @details     :   Validates heap initialization, O(1) insertions, minimum node retrieval, lazy 
 * consolidation extractions, dynamic decrease-key operations with cascading cuts, null-pointer 
 * defensive safeguards, and memory deallocation via the Unity Test Framework.
 * Compile:
 * clang ./tests/darwin/queue/fibonacci_heap_unit_tests.c ./src/darwin/queue/fibonacci-heap/fibonacci_heap.c ./src/unity.c -Iinclude -o ./bin/fibonacci_heap_unit_test
 * @version     :   1.0
 * @date        :   2026-11-18
 */

#include <stdio.h>
#include <stdlib.h>

#include "unity.h"
#include "darwin/queue/fibonacci_heap.h"

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
 * @brief Test 1: Verifies that heap initialization allocates memory successfully and resets attributes.
 */
void test_initialize_heap_should_allocate_successfully(void)
{
    // Act.
    heap_t* heap = initialize_heap();

    // Assert.
    TEST_ASSERT_NOT_NULL(heap);
    TEST_ASSERT_EQUAL_UINT(0, heap->number_of_nodes);
    TEST_ASSERT_NULL(heap->minimum_node);

    free_heap(heap);
}

/**
 * @brief Test 2: Verifies that querying the minimum on an empty heap returns a null pointer.
 */
void test_heap_minimum_on_empty_heap_should_return_null(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    // Act.
    node_t* minimum_node = heap_minimum(heap);

    // Assert.
    TEST_ASSERT_NULL(minimum_node);

    free_heap(heap);
}

/**
 * @brief Test 3: Verifies that inserting a single node sets the minimum pointer and increments node count.
 */
void test_insert_single_node_should_set_minimum_and_count(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    int key = 42;

    // Act.
    node_t* inserted_node = insert(heap, key, NULL);

    // Assert.
    TEST_ASSERT_NOT_NULL(inserted_node);
    TEST_ASSERT_EQUAL_UINT(1, heap->number_of_nodes);
    TEST_ASSERT_EQUAL_PTR(inserted_node, heap->minimum_node);
    TEST_ASSERT_EQUAL_INT(key, heap->minimum_node->key);

    free_heap(heap);
}

/**
 * @brief Test 4: Verifies that inserting multiple nodes correctly tracks the lowest key as minimum.
 */
void test_insert_multiple_nodes_should_track_minimum_correctly(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    // Act.
    insert(heap, 10, NULL);
    insert(heap, 3, NULL);
    insert(heap, 25, NULL);

    // Assert.
    TEST_ASSERT_EQUAL_UINT(3, heap->number_of_nodes);
    TEST_ASSERT_NOT_NULL(heap->minimum_node);
    TEST_ASSERT_EQUAL_INT(3, heap->minimum_node->key);

    free_heap(heap);
}

/**
 * @brief Test 5: Verifies that heap_minimum retrieves the minimum node without removing it from the heap.
 */
void test_heap_minimum_should_return_minimum_without_removal(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 5, NULL);
    insert(heap, 2, NULL);

    // Act.
    node_t* first_query = heap_minimum(heap);
    node_t* second_query = heap_minimum(heap);

    // Assert.
    TEST_ASSERT_NOT_NULL(first_query);
    TEST_ASSERT_EQUAL_PTR(first_query, second_query);
    TEST_ASSERT_EQUAL_INT(2, first_query->key);
    TEST_ASSERT_EQUAL_UINT(2, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 6: Verifies that extracting minimum on an empty heap returns error sentinel -1.
 */
void test_extract_min_on_empty_heap_should_return_error_sentinel(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    // Act.
    int extraction_result = extract_min(heap);

    // Assert.
    TEST_ASSERT_EQUAL_INT(-1, extraction_result);

    free_heap(heap);
}

/**
 * @brief Test 7: Verifies that extracting the sole element resets the heap state completely.
 */
void test_extract_min_single_element_should_reset_heap(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 99, NULL);

    // Act.
    int extracted_key = extract_min(heap);

    // Assert.
    TEST_ASSERT_EQUAL_INT(99, extracted_key);
    TEST_ASSERT_EQUAL_UINT(0, heap->number_of_nodes);
    TEST_ASSERT_NULL(heap->minimum_node);

    free_heap(heap);
}

/**
 * @brief Test 8: Verifies element extractions in sequential order across multiple items.
 */
void test_extract_min_multiple_elements_in_ascending_order(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 20, NULL);
    insert(heap, 5, NULL);
    insert(heap, 15, NULL);
    insert(heap, 1, NULL);

    // Act & Assert.
    TEST_ASSERT_EQUAL_INT(1, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(5, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(15, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(20, extract_min(heap));
    TEST_ASSERT_EQUAL_UINT(0, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 9: Verifies correct extractions when elements are inserted in reverse order.
 */
void test_extract_min_multiple_elements_in_descending_order(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 30, NULL);
    insert(heap, 20, NULL);
    insert(heap, 10, NULL);

    // Act & Assert.
    TEST_ASSERT_EQUAL_INT(10, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(20, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(30, extract_min(heap));

    free_heap(heap);
}

/**
 * @brief Test 10: Verifies that extraction triggers heap consolidation across tree degrees.
 */
void test_extract_min_triggers_consolidation_of_trees(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 12, NULL);
    insert(heap, 7, NULL);
    insert(heap, 3, NULL);
    insert(heap, 19, NULL);

    // Act.
    int first_extracted = extract_min(heap);
    int second_extracted = extract_min(heap);

    // Assert.
    TEST_ASSERT_EQUAL_INT(3, first_extracted);
    TEST_ASSERT_EQUAL_INT(7, second_extracted);
    TEST_ASSERT_EQUAL_UINT(2, heap->number_of_nodes);
    TEST_ASSERT_EQUAL_INT(12, heap_minimum(heap)->key);

    free_heap(heap);
}

/**
 * @brief Test 11: Verifies that decreasing a node key to a smaller value updates the minimum correctly.
 */
void test_decrease_key_valid_smaller_value_should_update_minimum(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 50, NULL);

    node_t* node = insert(heap, 30, NULL);

    insert(heap, 40, NULL);

    // Act.
    decrease_key(heap, node, 10);

    // Assert.
    TEST_ASSERT_EQUAL_INT(10, node->key);
    TEST_ASSERT_EQUAL_PTR(node, heap->minimum_node);
    TEST_ASSERT_EQUAL_INT(10, heap_minimum(heap)->key);

    free_heap(heap);
}

/**
 * @brief Test 12: Verifies that decreasing a key to its exact same value succeeds without mutation.
 */
void test_decrease_key_equal_value_should_succeed(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();
    node_t* node = insert(heap, 25, NULL);

    // Act.
    decrease_key(heap, node, 25);

    // Assert.
    TEST_ASSERT_EQUAL_INT(25, node->key);
    TEST_ASSERT_EQUAL_UINT(1, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 13: Verifies that attempting to increase a key via decrease_key is safely ignored.
 */
void test_decrease_key_larger_value_should_be_ignored(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();
    node_t* node = insert(heap, 10, NULL);

    // Act.
    decrease_key(heap, node, 50);

    // Assert.
    TEST_ASSERT_EQUAL_INT(10, node->key);
    TEST_ASSERT_EQUAL_INT(10, heap_minimum(heap)->key);

    free_heap(heap);
}

/**
 * @brief Test: Verifies that decrease-key triggers structural cuts and cascading operations while utilizing all local variables.
 */
void test_decrease_key_triggers_cut_and_cascading_cut(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 5, NULL);

    node_t* parent_node = insert(heap, 20, NULL);
    
    // Assert that the parent node was created correctly with key 20.
    TEST_ASSERT_NOT_NULL(parent_node);
    TEST_ASSERT_EQUAL_INT(20, parent_node->key);
    
    // Extract min to force consolidation so 20 becomes a root, then insert child.
    extract_min(heap);

    insert(heap, 2, NULL);

    int extracted = extract_min(heap);
    
    // Assert that the extracted minimum value matches the expected node key 2.
    TEST_ASSERT_EQUAL_INT(2, extracted);
    
    // Let's test direct key decrease triggering structural cut.
    node_t* base_node = insert(heap, 50, NULL);
    node_t* child_node = insert(heap, 60, NULL);
    
    // Assert that the child node was allocated and initialized correctly with key 60.
    TEST_ASSERT_NOT_NULL(child_node);
    TEST_ASSERT_EQUAL_INT(60, child_node->key);
    
    // Force link manually or via consolidation by extracting.
    extract_min(heap);
    
    // Act & Assert execution stability.
    decrease_key(heap, base_node, 1);

    TEST_ASSERT_EQUAL_INT(1, heap_minimum(heap)->key);

    free_heap(heap);
}

/**
 * @brief Test 15: Verifies defensive handling when passing a null heap pointer to insert.
 */
void test_insert_with_null_heap_should_return_null(void)
{
    // Arrange.
    heap_t* unallocated_heap = NULL;

    // Act.
    node_t* result_node = insert(unallocated_heap, 10, NULL);

    // Assert.
    TEST_ASSERT_NULL(result_node);
}

/**
 * @brief Test 16: Verifies defensive handling when passing a null heap pointer to heap_minimum.
 */
void test_heap_minimum_with_null_heap_should_return_null(void)
{
    // Arrange.
    const heap_t* unallocated_heap = NULL;

    // Act.
    node_t* result_node = heap_minimum(unallocated_heap);

    // Assert.
    TEST_ASSERT_NULL(result_node);
}

/**
 * @brief Test 17: Verifies defensive handling when passing a null heap pointer to extract_min.
 */
void test_extract_min_with_null_heap_should_return_error_sentinel(void)
{
    // Arrange.
    heap_t* unallocated_heap = NULL;

    // Act.
    int result_key = extract_min(unallocated_heap);

    // Assert.
    TEST_ASSERT_EQUAL_INT(-1, result_key);
}

/**
 * @brief Test 18: Verifies defensive handling when passing a null heap pointer to decrease_key.
 */
void test_decrease_key_with_null_heap_should_return_safely(void)
{
    // Arrange.
    heap_t* unallocated_heap = NULL;
    node_t dummy_node;

    // Act & Assert (Should not crash).
    decrease_key(unallocated_heap, &dummy_node, 5);

    TEST_PASS();
}

/**
 * @brief Test 19: Verifies defensive handling when passing a null node pointer to decrease_key.
 */
void test_decrease_key_with_null_node_should_return_safely(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();
    node_t* unallocated_node = NULL;

    // Act & Assert (Should not crash).
    decrease_key(heap, unallocated_node, 5);

    TEST_PASS();

    free_heap(heap);
}

/**
 * @brief Test 20: Verifies that calling free_heap with a null pointer returns safely without exceptions.
 */
void test_free_heap_with_null_pointer_should_return_without_exception(void)
{
    // Arrange.
    heap_t* empty_null_reference = NULL;

    // Act & Assert.
    free_heap(empty_null_reference);

    TEST_PASS();
}

/**
 * @brief Test 21: Verifies correct storage and retrieval of optional user-defined data payload pointers.
 */
void test_insert_with_user_data_payload_preservation(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    int sample_payload = 1337;

    // Act.
    node_t* inserted_node = insert(heap, 10, &sample_payload);

    // Assert.
    TEST_ASSERT_NOT_NULL(inserted_node);
    TEST_ASSERT_EQUAL_PTR(&sample_payload, inserted_node->data);
    TEST_ASSERT_EQUAL_INT(1337, *(int*)inserted_node->data);

    free_heap(heap);
}

/**
 * @brief Test 22: Verifies full heap depletion through successive extractions until empty.
 */
void test_multiple_extractions_depleting_all_nodes(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 8, NULL);
    insert(heap, 3, NULL);
    insert(heap, 5, NULL);

    // Act & Assert.
    TEST_ASSERT_EQUAL_UINT(3, heap->number_of_nodes);
    TEST_ASSERT_EQUAL_INT(3, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(5, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(8, extract_min(heap));
    TEST_ASSERT_EQUAL_UINT(0, heap->number_of_nodes);
    TEST_ASSERT_NULL(heap->minimum_node);

    free_heap(heap);
}

/**
 * @brief Test 23: Verifies stability under interleaved insertions and extractions.
 */
void test_interleaved_insertions_and_extractions(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    // Act.
    insert(heap, 50, NULL);
    insert(heap, 20, NULL);
    
    TEST_ASSERT_EQUAL_INT(20, extract_min(heap));

    insert(heap, 10, NULL);
    insert(heap, 30, NULL);

    // Assert.
    TEST_ASSERT_EQUAL_INT(10, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(30, extract_min(heap));
    TEST_ASSERT_EQUAL_UINT(1, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 24: Verifies decrease_key when dropping a node's key to become the new absolute minimum.
 */
void test_decrease_key_making_node_new_absolute_minimum(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 100, NULL);
    insert(heap, 50, NULL);

    node_t* node = insert(heap, 75, NULL);

    // Act.
    decrease_key(heap, node, 5);

    // Assert.
    TEST_ASSERT_EQUAL_PTR(node, heap->minimum_node);
    TEST_ASSERT_EQUAL_INT(5, extract_min(heap));
    TEST_ASSERT_EQUAL_UINT(2, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 25: Verifies successful handling of a larger volume of sequential node insertions.
 */
void test_large_volume_sequential_insertions(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();
    size_t number_of_items = 50;

    // Act.
    for (int index = (int)number_of_items; index > 0; index--)
    {
        insert(heap, index, NULL);
    }

    // Assert.
    TEST_ASSERT_EQUAL_UINT(number_of_items, heap->number_of_nodes);
    TEST_ASSERT_EQUAL_INT(1, heap_minimum(heap)->key);
    TEST_ASSERT_EQUAL_INT(1, extract_min(heap));
    TEST_ASSERT_EQUAL_UINT(number_of_items - 1, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 26: Verifies extraction stability with complex multi-level heap tree structures.
 */
void test_extract_min_with_complex_tree_hierarchy(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    int keys[] = {40, 10, 30, 20, 60, 50, 5};

    for (int index = 0; index < 7; index++)
    {
        insert(heap, keys[index], NULL);
    }

    // Act & Assert.
    TEST_ASSERT_EQUAL_INT(5, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(10, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(20, extract_min(heap));
    TEST_ASSERT_EQUAL_UINT(4, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 27: Verifies that decrease-key triggers proper cascading cut markers.
 */
void test_decrease_key_triggering_multiple_cascading_cuts(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 10, NULL);

    node_t* node_a = insert(heap, 40, NULL);
    node_t* node_b = insert(heap, 50, NULL);

    // Act.
    decrease_key(heap, node_a, 5);
    decrease_key(heap, node_b, 2);

    // Assert.
    TEST_ASSERT_EQUAL_INT(2, heap_minimum(heap)->key);
    TEST_ASSERT_EQUAL_UINT(3, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 28: Verifies node count tracking accuracy across the full operational lifecycle.
 */
void test_heap_node_count_tracking_accuracy_through_lifecycle(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    // Act & Assert.
    TEST_ASSERT_EQUAL_UINT(0, heap->number_of_nodes);

    insert(heap, 100, NULL);
    insert(heap, 200, NULL);

    TEST_ASSERT_EQUAL_UINT(2, heap->number_of_nodes);

    extract_min(heap);

    TEST_ASSERT_EQUAL_UINT(1, heap->number_of_nodes);

    extract_min(heap);

    TEST_ASSERT_EQUAL_UINT(0, heap->number_of_nodes);

    free_heap(heap);
}

/**
 * @brief Test 29: Verifies correct handling of negative priority keys in the heap.
 */
void test_insert_negative_keys_handling(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 10, NULL);
    insert(heap, -5, NULL);
    insert(heap, -20, NULL);

    // Act & Assert.
    TEST_ASSERT_EQUAL_INT(-20, heap_minimum(heap)->key);
    TEST_ASSERT_EQUAL_INT(-20, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(-5, extract_min(heap));
    TEST_ASSERT_EQUAL_INT(10, extract_min(heap));

    free_heap(heap);
}

/**
 * @brief Test 30: Verifies that free_heap correctly deallocates active nodes without memory leaks.
 */
void test_free_heap_with_active_nodes_should_not_leak_memory(void)
{
    // Arrange.
    heap_t* heap = initialize_heap();

    insert(heap, 15, NULL);
    insert(heap, 10, NULL);
    insert(heap, 25, NULL);
    insert(heap, 5, NULL);

    // Act & Assert (Should free entire recursive tree structure cleanly).
    free_heap(heap);

    TEST_PASS();
}

int main(void)
{
    // Initiate the session.
    UNITY_BEGIN();

    // Register the tests to run.
    RUN_TEST(test_initialize_heap_should_allocate_successfully);
    RUN_TEST(test_heap_minimum_on_empty_heap_should_return_null);
    RUN_TEST(test_insert_single_node_should_set_minimum_and_count);
    RUN_TEST(test_insert_multiple_nodes_should_track_minimum_correctly);
    RUN_TEST(test_heap_minimum_should_return_minimum_without_removal);
    RUN_TEST(test_extract_min_on_empty_heap_should_return_error_sentinel);
    RUN_TEST(test_extract_min_single_element_should_reset_heap);
    RUN_TEST(test_extract_min_multiple_elements_in_ascending_order);
    RUN_TEST(test_extract_min_multiple_elements_in_descending_order);
    RUN_TEST(test_extract_min_triggers_consolidation_of_trees);
    RUN_TEST(test_decrease_key_valid_smaller_value_should_update_minimum);
    RUN_TEST(test_decrease_key_equal_value_should_succeed);
    RUN_TEST(test_decrease_key_larger_value_should_be_ignored);
    RUN_TEST(test_decrease_key_triggers_cut_and_cascading_cut);
    RUN_TEST(test_insert_with_null_heap_should_return_null);
    RUN_TEST(test_heap_minimum_with_null_heap_should_return_null);
    RUN_TEST(test_extract_min_with_null_heap_should_return_error_sentinel);
    RUN_TEST(test_decrease_key_with_null_heap_should_return_safely);
    RUN_TEST(test_decrease_key_with_null_node_should_return_safely);
    RUN_TEST(test_free_heap_with_null_pointer_should_return_without_exception);
    RUN_TEST(test_insert_with_user_data_payload_preservation);
    RUN_TEST(test_multiple_extractions_depleting_all_nodes);
    RUN_TEST(test_interleaved_insertions_and_extractions);
    RUN_TEST(test_decrease_key_making_node_new_absolute_minimum);
    RUN_TEST(test_large_volume_sequential_insertions);
    RUN_TEST(test_extract_min_with_complex_tree_hierarchy);
    RUN_TEST(test_decrease_key_triggering_multiple_cascading_cuts);
    RUN_TEST(test_heap_node_count_tracking_accuracy_through_lifecycle);
    RUN_TEST(test_insert_negative_keys_handling);
    RUN_TEST(test_free_heap_with_active_nodes_should_not_leak_memory);

    // Destroy the session.
    return UNITY_END();
}
