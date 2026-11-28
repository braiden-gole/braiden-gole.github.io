/**
 * @file        :   suffix_tree.h
 * @author      :   Braiden Gole
 * @brief       :   Public interface and data structures for a compressed Suffix Tree.
 * @details     :   Defines the node and tree structures using edge index spans and a
 *                  unique sentinel character to support substring existence checks and
 *                  occurrence index queries in C.
 * @version     :   1.0
 * @date        :   2026-11-19
 */

#ifndef SUFFIX_TREE_H
#define SUFFIX_TREE_H

#include <stdbool.h>
#include <stdint.h>

#define ASCII_SIZE 256
#define SUFFIX_TREE_SENTINEL '$'

/**
 * @brief       :   Represents a compressed edge and node inside the suffix tree.
 * @details     :   Stores inclusive [start, end] character offsets into the tree corpus,
 *                  a leaf indicator, the starting suffix index when terminal, and child
 *                  pointers indexed by ASCII character value.
 */
typedef struct node
{
    int32_t start;
    int32_t end;
    int32_t position;
    bool is_leaf;
    struct node* children[ASCII_SIZE];
} node_t;

/**
 * @brief       :   Top-level handle for a constructed suffix tree instance.
 * @details     :   Owns the sentinel-terminated corpus buffer, tracks its character
 *                  length, and points to the root node of the suffix tree.
 */
typedef struct tree
{
    char* text;
    int32_t length;
    node_t* root;
} tree_t;


// Function prototypes.
tree_t* initialize_tree(const char* input);
bool contains(const tree_t* tree, const char* pattern);
size_t find(const tree_t* tree, const char* pattern, int32_t* output_indices, size_t maximum_output_capacity);
void destroy(tree_t* tree);

#endif