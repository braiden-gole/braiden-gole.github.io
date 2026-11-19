/**
 * @file        :   suffix_tree.c
 * @author      :   Braiden Gole
 * @brief       :   Compressed suffix tree implementation for substring search and occurrence indexing.
 * @details     :   Constructs a space-efficient suffix tree using start/end edge indices and a
 *                  unique sentinel character so every suffix terminates at a distinct leaf node.
 *                  Supports existence checks and recursive collection of starting indices.
 * @version     :   1.0
 * @date        :   2026-11-19
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Include header files.
#include "darwin/pattern-matching/suffix_tree.h"

/**
 * @brief       :   Allocates and initializes a single suffix tree node.
 * @details     :   Zero-initializes all child pointers via calloc and populates edge bounds,
 *                  leaf status, and original suffix starting position.
 * @param       :   start: Inclusive starting index of the node's incoming edge in the corpus.
 * @param       :   end: Inclusive ending index of the node's incoming edge in the corpus.
 * @param       :   position: Starting suffix index stored on leaf nodes, or -1 for internal nodes.
 * @param       :   is_leaf: Boolean flag indicating whether the node represents a terminal suffix leaf.
 * @return      :   node_t*: Pointer to the initialized node, or NULL on allocation failure.
 */
static inline node_t* initialize_node(int32_t start, int32_t end, int32_t position, bool is_leaf)
{
    // Allocate a zeroed node structure so all child pointers start as NULL.
    node_t* node = (node_t*)calloc(1, sizeof(node_t));

    if (node == NULL)
    {
        return NULL;
    }

    // Populate edge index bounds and node metadata fields.
    node->start = start;
    node->end = end;
    node->position = position;
    node->is_leaf = is_leaf;

    return node;
}

/**
 * @brief       :   Computes the character length of an edge represented by a node.
 * @details     :   Validates non-negative bounds and guards against signed 32-bit overflow
 *                  before returning the inclusive edge span length.
 * @param       :   node: Constant pointer to the node whose incoming edge span is measured.
 * @return      :   int32_t: Edge character span length, or 0 if the node or span is invalid.
 */
static inline int32_t retrieve_edge_span_length(const node_t* node)
{
    // Reject null nodes or malformed edge index ranges.
    if (node == NULL)
    {
        return 0;
    }

    if (node->start < 0 || node->end < node->start)
    {
        return 0;
    }

    // Guard against signed integer overflow when adding one to the span difference.
    if (node->end - node->start > INT32_MAX - 1)
    {
        return 0;
    }

    return node->end - node->start + 1;
}

/**
 * @brief       :   Inserts the suffix starting at a given index into the suffix tree.
 * @details     :   Walks existing edges from the root, creates a new leaf when a branch is open,
 *                  and splits an existing edge into an internal node plus a new leaf on mismatch.
 * @param       :   tree: Pointer to the suffix tree receiving the new suffix.
 * @param       :   index: Starting offset of the suffix within the sentinel-terminated corpus.
 * @return      :   bool: true if the suffix was inserted successfully, otherwise false.
 */
static inline bool insert(tree_t* tree, int32_t index)
{
    // Validate tree state and suffix index bounds before traversing.
    if (tree == NULL || tree->text == NULL || tree->root == NULL)
    {
        return false;
    }

    if (tree->length <= 0 || index < 0 || index >= tree->length)
    {
        return false;
    }

    node_t* root = tree->root;

    int32_t offset = index;

    while (offset < tree->length)
    {
        // Look up the child edge beginning with the next unmatched suffix character.
        unsigned char dispatch = (unsigned char)tree->text[offset];

        node_t* node = root->children[dispatch];

        // Attach a new terminal leaf directly when no outgoing edge starts with this character.
        if (node == NULL)
        {
            node_t* new_leaf = initialize_node(offset, tree->length - 1, index, true);

            if (new_leaf == NULL)
            {
                return false;
            }

            root->children[dispatch] = new_leaf;

            return true;
        }

        int32_t edge_count = retrieve_edge_span_length(node);

        if (edge_count <= 0)
        {
            return false;
        }

        int32_t cursor = 0;

        // Compare characters along the matched edge until the edge ends or a mismatch occurs.
        while (cursor < edge_count && offset + cursor < tree->length)
        {
            char existing_edge = tree->text[node->start + cursor];
            char incoming_suffix = tree->text[offset + cursor];

            if (existing_edge != incoming_suffix)
            {
                // Split the existing edge at the mismatch point and attach both branches.
                int32_t split_prefix_end = node->start + cursor - 1;

                node_t* split_internal_node = initialize_node(
                    node->start,
                    split_prefix_end,
                    -1,
                    false
                );

                if (split_internal_node == NULL)
                {
                    return false;
                }

                unsigned char retained_edge_key = (unsigned char)tree->text[node->start + cursor];
                unsigned char diverged_suffix_key = (unsigned char)tree->text[offset + cursor];

                node->start = node->start + cursor;

                node_t* diverged_suffix_leaf = initialize_node(
                    offset + cursor,
                    tree->length - 1,
                    index,
                    true
                );

                if (diverged_suffix_leaf == NULL)
                {
                    free(split_internal_node);

                    return false;
                }

                split_internal_node->children[retained_edge_key] = node;
                split_internal_node->children[diverged_suffix_key] = diverged_suffix_leaf;
                root->children[dispatch] = split_internal_node;

                return true;
            }

            cursor++;
        }

        // Advance down the tree if the entire edge matched, or fail if the suffix ended early.
        if (cursor == edge_count)
        {
            offset += edge_count;
            root = node;
        }
        else
        {
            // Reaching the end of the suffix before exhausting the edge should not happen with a unique sentinel.
            return false;
        }
    }

    // Every suffix terminates on a distinct leaf once a unique sentinel is appended.
    return false;
}

/**
 * @brief       :   Locates the edge or subtree corresponding to a non-empty search pattern.
 * @details     :   Matches the pattern character-by-character across compressed edges and reports
 *                  any remaining characters inside the final edge if the match ends mid-edge.
 * @param       :   tree: Constant pointer to the constructed suffix tree.
 * @param       :   pattern: Null-terminated pattern string to search for.
 * @param       :   edge_remainder: Optional output pointer receiving unread characters on the matched edge.
 * @return      :   const node_t*: Pointer to the node whose edge contains the pattern end, or NULL if absent.
 */
static inline const node_t* locate_pattern_subtree(const tree_t* tree, const char* pattern, int32_t* edge_remainder)
{
    // Reject uninitialized trees, empty corpus lengths, or null/empty patterns.
    if (tree == NULL || tree->text == NULL || tree->root == NULL)
    {
        return NULL;
    }

    if (tree->length <= 0 || pattern == NULL || pattern[0] == '\0')
    {
        return NULL;
    }

    size_t length = strlen(pattern);

    if (length > (size_t)tree->length || length >= (size_t)INT32_MAX)
    {
        return NULL;
    }

    int32_t pattern_length = (int32_t)length;

    const node_t* active_search_node = tree->root;

    int32_t cursor = 0;

    while (cursor < pattern_length)
    {
        // Step to the child edge matching the current pattern character.
        unsigned char query = (unsigned char)pattern[cursor];

        const node_t* subtree_edge = active_search_node->children[query];

        if (subtree_edge == NULL)
        {
            return NULL;
        }

        int32_t span_length = retrieve_edge_span_length(subtree_edge);

        if (span_length <= 0)
        {
            return NULL;
        }

        int32_t span_offset_cursor = 0;

        // Verify every character of the pattern against the current edge span.
        while (span_offset_cursor < span_length && cursor < pattern_length)
        {
            char corpus_edge_character = tree->text[subtree_edge->start + span_offset_cursor];

            if (corpus_edge_character != pattern[cursor])
            {
                return NULL;
            }

            span_offset_cursor++;
            cursor++;
        }

        // Return the matched edge as soon as the entire pattern has been matched.
        if (cursor == pattern_length)
        {
            if (edge_remainder != NULL)
            {
                *edge_remainder = (span_offset_cursor < span_length) ? (span_length - span_offset_cursor) : 0;
            }

            return subtree_edge;
        }

        active_search_node = subtree_edge;
    }

    return NULL;
}

/**
 * @brief       :   Recursively collects starting suffix positions from leaf nodes in a subtree.
 * @details     :   Traverses child slots in ascending ASCII order and stops once the output
 *                  buffer reaches maximum capacity.
 * @param       :   node: Constant pointer to the current subtree node being inspected.
 * @param       :   occurrences: Output array storing matched starting indices.
 * @param       :   capacity: Maximum number of indices that may be written to occurrences.
 * @param       :   count: Pointer to the running counter of collected match positions.
 * @return      :   void
 */
static inline void collect_leaf_positions_recursive(const node_t* node, int32_t* occurrences, size_t capacity, size_t* count)
{
    // Validate pointers and stop early if the caller's output buffer is full.
    if (node == NULL || occurrences == NULL || count == NULL)
    {
        return;
    }

    if (*count >= capacity)
    {
        return;
    }

    // Record the suffix start position once recursion reaches a leaf node.
    if (node->is_leaf)
    {
        if (node->position >= 0)
        {
            occurrences[*count] = node->position;
            (*count)++;
        }

        return;
    }

    // Walk non-null child subtrees in ascending character-code order.
    for (int32_t alphabet_slot = 0; alphabet_slot < ASCII_SIZE; alphabet_slot++)
    {
        if (*count >= capacity)
        {
            return;
        }

        if (node->children[alphabet_slot] != NULL)
        {
            collect_leaf_positions_recursive(
                node->children[alphabet_slot],
                occurrences,
                capacity,
                count
            );
        }
    }
}

/**
 * @brief       :   Recursively frees a suffix tree node and all of its descendants.
 * @details     :   Performs post-order traversal across all ASCII child slots and clears
 *                  pointers before freeing each node.
 * @param       :   node: Pointer to the root of the node subtree to destroy.
 * @return      :   void
 */
static inline void destroy_suffix_tree_node_recursive(node_t* node)
{
    if (node == NULL)
    {
        return;
    }

    // Destroy all child subtrees before releasing the parent node.
    for (int32_t ascii_slot = 0; ascii_slot < ASCII_SIZE; ascii_slot++)
    {
        if (node->children[ascii_slot] != NULL)
        {
            destroy_suffix_tree_node_recursive(node->children[ascii_slot]);

            node->children[ascii_slot] = NULL;
        }
    }

    free(node);
}

/**
 * @brief       :   Constructs and populates a suffix tree for a given input corpus.
 * @details     :   Copies the corpus with an appended unique sentinel character, allocates
 *                  the root node, and inserts every suffix position into the tree.
 * @param       :   corpus: Null-terminated input string to index.
 * @return      :   tree_t*: Pointer to the constructed suffix tree, or NULL on failure.
 */
tree_t* initialize_tree(const char* corpus)
{
    // Reject null or empty inputs as well as lengths that would overflow index bounds.
    if (corpus == NULL)
    {
        return NULL;
    }

    size_t length = strlen(corpus);

    if (length == 0 || length >= (size_t)(INT32_MAX - 1))
    {
        return NULL;
    }

    // Allocate the top-level tree structure.
    tree_t* tree = (tree_t*)calloc(1, sizeof(tree_t));

    if (tree == NULL)
    {
        return NULL;
    }

    // Append a unique sentinel character so every suffix ends on its own leaf.
    size_t corpus_with_sentinel_length = length + 1;

    tree->length = (int32_t)corpus_with_sentinel_length;
    tree->text = (char*)malloc(corpus_with_sentinel_length + 1);

    if (tree->text == NULL)
    {
        free(tree);

        return NULL;
    }

    memcpy(tree->text, corpus, length);

    tree->text[length] = SUFFIX_TREE_SENTINEL;
    tree->text[corpus_with_sentinel_length] = '\0';

    // Initialize the empty root node before inserting suffixes.
    tree->root = initialize_node(-1, -1, -1, false);

    if (tree->root == NULL)
    {
        free(tree->text);

        tree->text = NULL;

        free(tree);

        return NULL;
    }

    // Insert every suffix starting index into the tree.
    for (int32_t suffix_index = 0; suffix_index < tree->length; suffix_index++)
    {
        bool is_inserted = insert(tree, suffix_index);

        if (!is_inserted)
        {
            destroy(tree);

            return NULL;
        }
    }

    return tree;
}

/**
 * @brief       :   Checks whether a pattern occurs anywhere inside the suffix tree corpus.
 * @details     :   Locates the pattern path in the tree and returns true if a matching
 *                  edge or subtree exists.
 * @param       :   tree: Constant pointer to the suffix tree to query.
 * @param       :   pattern: Null-terminated substring pattern to check.
 * @return      :   bool: true if the pattern is found, otherwise false.
 */
bool contains(const tree_t* tree, const char* pattern)
{
    // Ignore edge remainder details since only existence is tested here.
    int32_t edge_characters = 0;

    const node_t* node = locate_pattern_subtree(tree, pattern, &edge_characters);

    return node != NULL;
}

/**
 * @brief       :   Finds starting indices of all occurrences of a pattern in the corpus.
 * @details     :   Locates the subtree matching the pattern and gathers leaf suffix positions
 *                  up to the caller-specified capacity.
 * @param       :   tree: Constant pointer to the suffix tree to search.
 * @param       :   pattern: Null-terminated substring pattern to locate.
 * @param       :   indices: Output buffer receiving starting occurrence indices.
 * @param       :   maximum_capacity: Maximum number of indices that may be written to indices.
 * @return      :   size_t: Total number of occurrence indices written to indices.
 */
size_t find(const tree_t* tree, const char* pattern, int32_t* indices, size_t maximum_capacity)
{
    // Validate output buffer availability before searching.
    if (indices == NULL || maximum_capacity == 0)
    {
        return 0;
    }

    int32_t edge_characters = 0;

    const node_t* root = locate_pattern_subtree(tree, pattern, &edge_characters);

    if (root == NULL)
    {
        return 0;
    }

    // Collect matching starting positions across all leaves under the matched subtree.
    size_t ocurrences = 0;

    collect_leaf_positions_recursive(root, indices, maximum_capacity, &ocurrences);

    return ocurrences;
}

/**
 * @brief       :   Frees all memory associated with a suffix tree instance.
 * @details     :   Recursively destroys all nodes starting from the root, frees the stored
 *                  corpus buffer, and releases the tree handle.
 * @param       :   tree: Pointer to the suffix tree to destroy.
 * @return      :   void
 */
void destroy(tree_t* tree)
{
    if (tree == NULL)
    {
        return;
    }

    // Tear down all tree nodes before releasing the corpus text and tree struct.
    destroy_suffix_tree_node_recursive(tree->root);

    tree->root = NULL;

    free(tree->text);

    tree->text = NULL;

    free(tree);
}