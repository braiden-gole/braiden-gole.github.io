/**
 * @file        :   main.c
 * @author      :   Braiden Gole
 * @brief       :   Demonstrates substring search and occurrence indexing using a Suffix Tree.
 * @details     :   Constructs a compressed suffix tree over the sample corpus
 * "banana_bandana" and evaluates multiple query patterns—including repeated
 * substrings, missing substrings, and an empty pattern—to verify substring
 * existence checks and starting index collection.
 * Compile:
 * clang ./src/darwin/pattern-matching/suffix-tree/main.c ./src/darwin/pattern-matching/suffix-tree/suffix_tree.c -I./include -o ./bin/suffix_tree_test
 * @version     :   1.0
 * @date        :   2026-11-19
 */
#include <stdio.h>
#include <stdlib.h>

// Include header files.
#include "darwin/pattern-matching/suffix_tree.h"

/**
 * @brief   Prints an array of occurrence indices formatted as a bracketed list.
 * @param   occurrence_indices  Pointer to the array of starting indices to display.
 * @param   occurrence_count    Number of valid entries stored in occurrence_indices.
 * @return  None.
 */
static inline void print_occurrence_indices(const int32_t* occurrence_indices, size_t occurrence_count)
{
    // Guard against null or empty occurrence arrays before printing.
    if (occurrence_indices == NULL || occurrence_count == 0)
    {
        printf("%s", "[]");

        return;
    }

    // Open the bracketed list of indices.
    printf("%c", '[');

    for (size_t index_cursor = 0; index_cursor < occurrence_count; index_cursor++)
    {
        // Append a comma separator between consecutive indices.
        if (index_cursor < occurrence_count - 1)
        {
            printf("%i, ", occurrence_indices[index_cursor]);
        }
        else
        {
            printf("%i", occurrence_indices[index_cursor]);
        }
    }

    // Close the bracketed list of indices.
    printf("%c", ']');
}

int main(void)
{
    const char* sample_corpus = "banana_bandana";

    // Build the suffix tree from the sample corpus string.
    tree_t* suffix_tree = initialize_tree(sample_corpus);

    if (suffix_tree == NULL)
    {
        fprintf(stderr, "%s", "Failed to build suffix tree for sample.\n");

        return EXIT_FAILURE;
    }

    printf("Corpus: \"%s\"\n\n", sample_corpus);

    const char* query_patterns[] = {
        "ana",
        "band",
        "nana",
        "ban",
        "xyz",
        ""
    };

    size_t query_pattern_count = sizeof(query_patterns) / sizeof(query_patterns[0]);

    for (size_t pattern_cursor = 0; pattern_cursor < query_pattern_count; pattern_cursor++)
    {
        const char* active_pattern = query_patterns[pattern_cursor];

        // Skip null entries defensively in case the pattern table changes.
        if (active_pattern == NULL)
        {
            continue;
        }

        // Check substring presence and collect up to 32 starting positions.
        bool is_substring_present = contains(suffix_tree, active_pattern);

        int32_t matched_positions[32] = {0};

        size_t match_total = find(
            suffix_tree,
            active_pattern,
            matched_positions,
            sizeof(matched_positions) / sizeof(matched_positions[0])
        );

        printf("Pattern \"%-8s\" | contains: %-5s | occurrences (%zu): ",
               active_pattern,
               is_substring_present ? "true" : "false",
               match_total);

        print_occurrence_indices(matched_positions, match_total);

        printf("\n");
    }

    // Release all heap memory allocated for the suffix tree.
    destroy(suffix_tree);

    suffix_tree = NULL;

    return EXIT_SUCCESS;
}