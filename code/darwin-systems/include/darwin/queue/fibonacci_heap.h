/**
 * @file        :   fibonacci_heap.h
 * @author      :   Braiden Gole
 * @brief       :   Header definition for an advanced Fibonacci Heap priority queue structure.
 * @details     :   Declares node and heap structures along with public interface prototypes 
 * for managing $O(1)$ amortized priority queue operations, dynamic decrease-keys, 
 * and safe memory deallocation.
 * @version     :   1.0
 * @date        :   2026-11-18
 */

#ifndef FIBONACCI_HEAP_H
#define FIBONACCI_HEAP_H

#include <stdbool.h>
#include <stddef.h>

/**
 * @struct node
 * @brief Represents an individual node wrapper within the Fibonacci Heap structure.
 */
typedef struct node
{
    int key;
    void* data;
    size_t degree_count;
    bool is_marked;
    struct node* parent;
    struct node* child;
    struct node* left;
    struct node* right;
} node_t;

/**
 * @struct heap
 * @brief Represents the top-level container structure for the Fibonacci Heap.
 */
typedef struct heap
{
    size_t number_of_nodes;
    node_t* minimum_node;
} heap_t;

// Function prototypes.
heap_t* initialize_heap(void);
node_t* insert(heap_t* heap, int key, void* data);
node_t* heap_minimum(const heap_t* heap);
int extract_min(heap_t* heap);
void decrease_key(heap_t* heap, node_t* node, int new_key);
void free_heap(heap_t* heap);

#endif
