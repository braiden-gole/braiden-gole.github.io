/**
 * @file        :   fibonacci_heap.c
 * @author      :   Braiden Gole
 * @brief       :   Implementation of an advanced Fibonacci Heap priority queue structure.
 * @details     :   Provides $O(1)$ amortized time complexity for insertions, minimum node 
 * retrieval, and decrease-key operations. Uses circular doubly linked lists and 
 * lazy consolidation during node extractions.
 * @version     :   1.0
 * @date        :   2026-11-18
 */

#include <stdio.h>
#include <stdlib.h>

// Include header files.
#include "darwin/queue/fibonacci_heap.h"

/**
 * @brief Allocates and initializes an empty Fibonacci heap container instance.
 * @return Pointer to the newly created heap_t instance, or NULL on allocation failure.
 */
heap_t* initialize_heap(void)
{
    // Allocate memory for the main heap container structure.
    heap_t* heap = (heap_t*)malloc(sizeof(heap_t));

    // Validate defensively that memory allocation succeeded.
    if (heap == NULL)
    {
        return NULL;
    }

    // Initialize node counter and minimum pointer to default empty states.
    heap->number_of_nodes = 0;
    heap->minimum_node = NULL;

    // Return the successfully initialized heap pointer.
    return heap;
}

/**
 * @brief Inserts a new key-value pair into the Fibonacci heap in O(1) time.
 * @param heap Pointer to the target heap instance.
 * @param key Integer priority key value.
 * @param data Optional user-defined data payload pointer.
 * @return Pointer to the newly created node_t wrapper instance, or NULL on error.
 */
node_t* insert(heap_t* heap, int key, void* data)
{
    // Validate defensively that the target heap pointer is not null.
    if (heap == NULL)
    {
        return NULL;
    }

    // Allocate memory for the new heap node structure.
    node_t* node = (node_t*)malloc(sizeof(node_t));

    // Validate defensively that node memory allocation succeeded.
    if (node == NULL)
    {
        return NULL;
    }

    // Assign the input key and user data payload fields.
    node->key = key;
    node->data = data;

    // Initialize degree count and structural marker states.
    node->degree_count = 0;

    node->is_marked = false;
    
    // Set parent and child relationship pointers to null.
    node->parent = NULL;
    node->child = NULL;

    // Configure sibling pointers to create an isolated circular list of size one.
    node->left = node;
    node->right = node;

    // Check if the root list is currently empty.
    if (heap->minimum_node == NULL)
    {
        heap->minimum_node = node;
    }
    else
    {
        // Link the new node into the root list right beside the current minimum node.
        node->right = heap->minimum_node->right;
        node->left = heap->minimum_node;

        heap->minimum_node->right->left = node;
        heap->minimum_node->right = node;

        // Update the minimum pointer if the new node's key is smaller.
        if (node->key < heap->minimum_node->key)
        {
            heap->minimum_node = node;
        }
    }

    // Increment the total count of active nodes in the heap container.
    heap->number_of_nodes++;

    // Return a reference pointer to the newly created node.
    return node;
}

/**
 * @brief Returns the minimum node in O(1) time without removing it.
 * @param heap Read-only pointer to the target heap instance.
 * @return Pointer to the minimum node_t instance, or NULL if empty.
 */
node_t* heap_minimum(const heap_t* heap)
{
    // Validate defensively that the heap pointer is not null.
    if (heap == NULL)
    {
        return NULL;
    }

    // Return the cached minimum node pointer directly in O(1) time.
    return heap->minimum_node;
}

/**
 * @brief Internal helper function to link two trees of matching degrees together.
 * @param heap Pointer to the target heap instance.
 * @param child Node to be subordinated as a child.
 * @param parent Parent node receiving the child.
 */
static void link_trees(heap_t* heap, node_t* child, node_t* parent)
{
    // Validate defensively that all pointer arguments are valid.
    if (heap == NULL || child == NULL || parent == NULL)
    {
        return;
    }

    // Remove the child node from the root list by updating its siblings' pointers.
    child->left->right = child->right;
    child->right->left = child->left;

    // Establish parent pointer reference on the child node.
    child->parent = parent;

    // Check if the parent node currently has no children.
    if (parent->child == NULL)
    {
        parent->child = child;

        child->left = child;
        child->right = child;
    }
    else
    {
        // Insert the child into the parent's existing circular child list.
        child->right = parent->child->right;
        child->left = parent->child;

        parent->child->right->left = child;
        parent->child->right = child;
    }

    // Increment the degree count of the parent node.
    parent->degree_count++;

    // Reset the child node's marked flag as required by Fibonacci heap rules.
    child->is_marked = false;
}

/**
 * @brief Consolidates the root list of the heap, combining trees with matching degrees.
 * @param heap Pointer to the target heap instance.
 */
static void consolidate_heap(heap_t* heap)
{
    // Validate defensively that the heap and its minimum pointer are valid.
    if (heap == NULL || heap->minimum_node == NULL)
    {
        return;
    }

    // Define a safe maximum degree bound for the degree lookup table.
    size_t maximum_degree_bound = 90;

    // Allocate a lookup table to track root trees by their degree counts.
    node_t** degree_table = (node_t**)calloc(maximum_degree_bound, sizeof(node_t*));

    // Validate defensively that table memory allocation succeeded.
    if (degree_table == NULL)
    {
        return;
    }

    // Count the total number of root nodes currently present in the root list.
    size_t count = 0;

    node_t* current = heap->minimum_node;

    if (current != NULL)
    {
        node_t* node = current;

        do {
            count++;

            current = current->right;
        } while (current != node);
    }

    // Allocate a temporary array to store pointers to all root nodes safely.
    node_t** array = (node_t**)malloc(count * sizeof(node_t*));

    // Validate defensively that the temporary array allocation succeeded.
    if (array == NULL)
    {
        free(degree_table);

        return;
    }

    // Populate the temporary array by iterating through the circular root list.
    current = heap->minimum_node;

    for (size_t index = 0; index < count; index++)
    {
        array[index] = current;

        current = current->right;
    }

    // Process each root node and merge trees that share the exact same degree.
    for (size_t index = 0; index < count; index++)
    {
        node_t* process_node = array[index];

        size_t node_degree = process_node->degree_count;

        // Keep linking trees of matching degree until no conflicts remain in the table slot.
        while (node_degree < maximum_degree_bound && degree_table[node_degree])
        {
            node_t* conflicting_node = degree_table[node_degree];

            // Ensure the node with the smaller key acts as the parent tree root.
            if (process_node->key > conflicting_node->key)
            {
                node_t* temporary = process_node;
                process_node = conflicting_node;
                conflicting_node = temporary;
            }

            // Link the conflicting tree under the current process node.
            link_trees(heap, conflicting_node, process_node);

            // Clear the table entry and increment the degree index for the next check.
            degree_table[node_degree] = NULL;
            node_degree++;
        }

        // Store the final resulting tree root into its corresponding degree table slot.
        if (node_degree < maximum_degree_bound)
        {
            degree_table[node_degree] = process_node;
        }
    }

    // Temporarily reset the minimum node pointer to rebuild the root list.
    heap->minimum_node = NULL;

    // Reconstruct the root list from the non-empty entries in the degree table.
    for (size_t index = 0; index < maximum_degree_bound; index++)
    {
        if (degree_table[index] != NULL)
        {
            if (heap->minimum_node == NULL)
            {
                heap->minimum_node = degree_table[index];
                heap->minimum_node->left = heap->minimum_node;
                heap->minimum_node->right = heap->minimum_node;
            }
            else
            {
                // Insert the table root tree back into the circular root list.
                degree_table[index]->right = heap->minimum_node->right;
                degree_table[index]->left = heap->minimum_node;

                heap->minimum_node->right->left = degree_table[index];
                heap->minimum_node->right = degree_table[index];

                // Update the minimum pointer if this tree root holds a smaller key.
                if (degree_table[index]->key < heap->minimum_node->key)
                {
                    heap->minimum_node = degree_table[index];
                }
            }
        }
    }

    // Free the temporary array and degree table memory buffers.
    free(array);
    free(degree_table);
}

/**
 * @brief Extracts and removes the minimum element from the Fibonacci heap.
 * @param heap Pointer to the target heap instance.
 * @return The integer key of the extracted minimum node, or -1 on error.
 */
int extract_min(heap_t* heap)
{
    // Validate defensively that the heap is not null and not empty.
    if (heap == NULL || heap->minimum_node == NULL)
    {
        return -1;
    }

    // Cache the minimum node pointer and its priority key value.
    node_t* minimum_node = heap->minimum_node;

    int minimum_key = minimum_node->key;

    // Promote all children of the minimum node directly into the root list.
    if (minimum_node->child != NULL)
    {
        node_t* child_current = minimum_node->child;
        node_t* first_child = child_current;

        do
        {
            node_t* next_child = child_current->right;

            // Clear parent pointers since these nodes are now root trees.
            child_current->parent = NULL;

            // Splice each child node into the root list next to the minimum node.
            child_current->right = heap->minimum_node->right;
            child_current->left = heap->minimum_node;

            heap->minimum_node->right->left = child_current;
            heap->minimum_node->right = child_current;

            child_current = next_child;
        } while (child_current != first_child);
    }

    // Remove the minimum node from the circular root list.
    minimum_node->left->right = minimum_node->right;
    minimum_node->right->left = minimum_node->left;

    // Check if the extracted node was the only node left in the root list.
    if (minimum_node == minimum_node->right)
    {
        heap->minimum_node = NULL;
    }
    else
    {
        // Temporarily assign a valid root node and trigger heap consolidation.
        heap->minimum_node = minimum_node->right;

        consolidate_heap(heap);
    }

    // Decrement the total count of nodes and free the extracted node memory.
    heap->number_of_nodes--;

    free(minimum_node);

    // Return the key value of the extracted minimum element.
    return minimum_key;
}

/**
 * @brief Cuts a node from its parent and re-inserts it into the root list.
 * @param heap Pointer to the target heap instance.
 * @param child Pointer to the node being cut.
 * @param parent Pointer to the parent node.
 */
static void cut_node(heap_t* heap, node_t* child, node_t* parent)
{
    // Validate defensively that all pointer arguments are valid.
    if (heap == NULL || child == NULL || parent == NULL)
    {
        return;
    }

    // Update parent's child pointer if the cut node was directly referenced.
    if (parent->child == child)
    {
        if (child->right == child)
        {
            parent->child = NULL;
        }
        else
        {
            parent->child = child->right;
        }
    }

    // Remove the child node from its sibling list under the parent.
    child->left->right = child->right;
    child->right->left = child->left;

    // Decrement the parent node's degree count.
    parent->degree_count--;

    // Insert the cut node back into the root list of the heap.
    child->right = heap->minimum_node->right;
    child->left = heap->minimum_node;

    heap->minimum_node->right->left = child;
    heap->minimum_node->right = child;

    // Clear parent reference and reset the marked flag state.
    child->parent = NULL;
    child->is_marked = false;
}

/**
 * @brief Performs cascading cuts recursively to maintain structural heap invariants.
 * @param heap Pointer to the target heap instance.
 * @param node Pointer to the target node that was cut.
 */
static void cascading_cut(heap_t* heap, node_t* node)
{
    // Validate defensively that heap and node pointers are valid.
    if (heap == NULL || node == NULL)
    {
        return;
    }

    // Retrieve the parent pointer of the target node.
    node_t* parent = node->parent;

    // Check if a parent exists to evaluate for cascading cuts.
    if (parent != NULL)
    {
        // If the node is not marked, mark it; otherwise, cut it and recurse upward.
        if (!node->is_marked)
        {
            node->is_marked = true;
        }
        else
        {
            cut_node(heap, node, parent);
            cascading_cut(heap, parent);
        }
    }
}

/**
 * @brief Decreases the key value of a specific node in the heap.
 * @param heap Pointer to the target heap instance.
 * @param node Pointer to the target node.
 * @param key New lower integer key value.
 */
void decrease_key(heap_t* heap, node_t* node, int key)
{
    // Validate defensively that heap and node pointers are valid.
    if (heap == NULL || node == NULL)
    {
        return;
    }

    // Validate defensively that the new key value is strictly lower than the current key.
    if (key > node->key)
    {
        return;
    }

    // Update the node's key value.
    node->key = key;

    // Retrieve the parent pointer of the updated node.
    node_t* parent = node->parent;

    // Check if heap order property is violated (child key smaller than parent key).
    if (parent != NULL && node->key < parent->key)
    {
        cut_node(heap, node, parent);
        cascading_cut(heap, parent);
    }

    // Update the global minimum node pointer if the new key is smaller than the current minimum.
    if (heap->minimum_node == NULL || node->key < heap->minimum_node->key)
    {
        heap->minimum_node = node;
    }
}

/**
 * @brief Recursively frees memory allocated for node trees.
 * @param node Pointer to the root of the node tree.
 */
static void free_node_recursive(node_t* node)
{
    // Base case: return immediately if the node pointer is null.
    if (node == NULL)
    {
        return;
    }

    // Walk through the circular sibling list and free each node and its subtrees.
    node_t* current = node;

    do
    {
        node_t* next = current->right;

        // Recursively free any child subtrees attached to this node.
        if (current->child != NULL)
        {
            free_node_recursive(current->child);
        }

        // Free the current node memory allocation.
        free(current);

        current = next;
    } while (current != node);
}

/**
 * @brief Destroys and cleans up all memory allocated for the heap and its nodes.
 * @param heap Pointer to the target heap instance.
 */
void free_heap(heap_t* heap)
{
    // Validate defensively that the heap pointer is not null.
    if (heap == NULL)
    {
        return;
    }

    // Recursively free all node trees attached to the heap if any exist.
    if (heap->minimum_node != NULL)
    {
        free_node_recursive(heap->minimum_node);
    }

    // Free the top-level heap container structure memory.
    free(heap);
}