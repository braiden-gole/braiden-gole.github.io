/**
 * @file        :   main.c
 * @author      :   Braiden Gole
 * @brief       :   Driver program to test a Fibonacci Heap implementation.
 * @details     :   This file demonstrates the basic functionality of the Fibonacci Heap data structure.
 *                  It tests heap initialization, node insertions, minimum value retrieval, element extractions,
 *                  dynamic decrease-key operations, and proper memory cleanup.
 * Compile:
 * clang ./src/darwin/queue/fibonacci-heap/main.c ./src/darwin/queue/fibonacci-heap/fibonacci_heap.c -Iinclude -o ./bin/fibonacci_heap_test
 * @version     :   1.0
 * @date        :   2026-11-18
 */

#include <stdio.h>
#include <stdlib.h>

// Include header files.
#include "darwin/queue/fibonacci_heap.h"

int main(void)
{
    // Initialize an empty Fibonacci heap container.
    heap_t* heap = initialize_heap();

    // Validate defensively that heap memory allocation succeeded.
    if (heap == NULL)
    {
        fprintf(stderr, "%s", "Failed to initialize heap.\n");

        return EXIT_FAILURE;
    }

    // Insert multiple test nodes into the heap, preserving a reference to node ten.
    node_t* ten = insert(heap, 10, NULL);

    insert(heap, 3, NULL);
    insert(heap, 15, NULL);
    insert(heap, 1, NULL);
    insert(heap, 7, NULL);

    // Retrieve the current minimum node from the heap structure.
    node_t* current_minimum = heap_minimum(heap);

    // Print the current minimum key value if the node pointer is valid.
    if (current_minimum != NULL)
    {
        printf("%s %i%c", "Current minimum key value:", current_minimum->key, '\n');
    }

    // Extract the minimum elements sequentially from the heap.
    printf("Extracting minimum element: %d\n", extract_min(heap));
    printf("Extracting minimum element: %d\n", extract_min(heap));

    // Retrieve the updated minimum node reference after extractions.
    current_minimum = heap_minimum(heap);

    // Print the new minimum key value if the pointer is valid.
    if (current_minimum != NULL)
    {
        printf("New minimum key value after extractions: %d\n", current_minimum->key);
    }

    // Perform a decrease-key operation to lower node ten's key value down to zero.
    printf("Decreasing key value of node 10 to 0...\n");

    decrease_key(heap, ten, 0);

    // Retrieve the minimum node after the decrease-key operation triggers structural updates.
    current_minimum = heap_minimum(heap);

    // Print the resulting minimum key value if valid.
    if (current_minimum != NULL)
    {
        printf("Minimum key value after decrease-key operation: %d\n", current_minimum->key);
    }

    // Clean up and release all dynamically allocated heap memory recursively.
    printf("Cleaning up and freeing heap memory...\n");

    free_heap(heap);

    // Indicate successful test execution and exit program.
    printf("Fibonacci heap test execution completed successfully.\n");

    return EXIT_SUCCESS;
}