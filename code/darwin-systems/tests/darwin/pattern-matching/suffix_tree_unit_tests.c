/**
 * @file        :   suffix_tree_unit_tests.c
 * @author      :   Braiden Gole
 * @brief       :   Unit tests for the compressed Suffix Tree pattern matching engine.
 * @details     :   Validates tree construction, exact substring matching, multiple and
 *                  overlapping occurrence indexing, edge-splitting behavior, empty or
 *                  missing patterns, output buffer capacity bounds, and defensive NULL
 *                  handling using the Unity Test Framework.
 * Compile:
 * clang ./tests/darwin/pattern-matching/suffix_tree_unit_tests.c ./src/darwin/pattern-matching/suffix-tree/suffix_tree.c ./src/unity.c -I./include -o ./bin/suffix_tree_unit_test
 * @version     :   1.0
 * @date        :   2026-11-19
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "unity.h"
#include "darwin/pattern-matching/suffix_tree.h"

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
 * @brief Verifies that a valid non-empty corpus constructs a non-null suffix tree handle.
 */
void test_suffix_tree_initializes_valid_corpus_successfully(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";

    // Act.
    tree_t* tree = initialize_tree(corpus);

    // Assert.
    TEST_ASSERT_NOT_NULL(tree);
    TEST_ASSERT_NOT_NULL(tree->root);
    TEST_ASSERT_NOT_NULL(tree->text);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that passing a NULL corpus pointer to initialize_tree returns NULL safely.
 */
void test_suffix_tree_returns_null_when_corpus_is_null(void)
{
    // Arrange.
    const char* null_corpus = NULL;

    // Act.
    tree_t* tree = initialize_tree(null_corpus);

    // Assert.
    TEST_ASSERT_NULL(tree);
}

/**
 * @brief Verifies that passing an empty string corpus to initialize_tree returns NULL safely.
 */
void test_suffix_tree_returns_null_when_corpus_is_empty(void)
{
    // Arrange.
    const char* empty_corpus = "";

    // Act.
    tree_t* tree = initialize_tree(empty_corpus);

    // Assert.
    TEST_ASSERT_NULL(tree);
}

/**
 * @brief Verifies that a single pattern match located in the middle of a string is found cleanly.
 */
void test_suffix_tree_finds_single_pattern_match_in_middle(void)
{
    // Arrange.
    const char* corpus = "abcdefgh";
    const char* pattern = "cde";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[8] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 8);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(1, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(2, indices[0]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a pattern starting at index 0 of the corpus matches accurately.
 */
void test_suffix_tree_matches_pattern_at_start_of_corpus(void)
{
    // Arrange.
    const char* corpus = "darwin_pattern_engine";
    const char* pattern = "darwin";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[21] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 21);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(1, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(0, indices[0]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a pattern terminating at the final character of the corpus matches accurately.
 */
void test_suffix_tree_matches_pattern_at_end_of_corpus(void)
{
    // Arrange.
    const char* corpus = "suffix_tree_structure";
    const char* pattern = "structure";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[21] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 21);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(1, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(12, indices[0]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that multiple occurrences of a repeated pattern are all collected.
 */
void test_suffix_tree_finds_multiple_occurrences_in_banana_bandana(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";
    const char* pattern = "ana";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[14] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 14);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(3, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(11, indices[0]);
    TEST_ASSERT_EQUAL_INT32(3, indices[1]);
    TEST_ASSERT_EQUAL_INT32(1, indices[2]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that overlapping repeating characters produce all valid starting positions.
 */
void test_suffix_tree_handles_overlapping_repeated_substrings(void)
{
    // Arrange.
    const char* corpus = "aaaaa";
    const char* pattern = "aaa";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[5] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 5);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(3, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(2, indices[0]);
    TEST_ASSERT_EQUAL_INT32(1, indices[1]);
    TEST_ASSERT_EQUAL_INT32(0, indices[2]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a pattern absent from the corpus returns false and zero occurrences.
 */
void test_suffix_tree_returns_false_for_missing_pattern(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";
    const char* pattern = "xyz";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[14] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 14);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that searching for the entire corpus string matches at index 0.
 */
void test_suffix_tree_matches_entire_corpus_as_pattern(void)
{
    // Arrange.
    const char* corpus = "exact_corpus_match";
    const char* pattern = "exact_corpus_match";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[18] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 18);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(1, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(0, indices[0]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a pattern longer than the corpus returns false and zero occurrences.
 */
void test_suffix_tree_rejects_pattern_longer_than_corpus(void)
{
    // Arrange.
    const char* corpus = "short";
    const char* pattern = "longer_than_corpus";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[5] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 5);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that pattern matching distinguishes uppercase and lowercase letters accurately.
 */
void test_suffix_tree_enforces_case_sensitivity_during_search(void)
{
    // Arrange.
    const char* corpus = "SuffixTreeEngine";
    const char* pattern = "suffixtree";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[16] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 16);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a single-character pattern finds all matching positions in ASCII traversal order.
 */
void test_suffix_tree_finds_all_single_character_occurrences(void)
{
    // Arrange.
    const char* corpus = "abaca";
    const char* pattern = "a";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[5] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 5);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(3, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(4, indices[0]);
    TEST_ASSERT_EQUAL_INT32(0, indices[1]);
    TEST_ASSERT_EQUAL_INT32(2, indices[2]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a single-character corpus handles queries for that character and others.
 */
void test_suffix_tree_handles_single_character_corpus(void)
{
    // Arrange.
    const char* corpus = "z";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[2] = { 0 };

    // Act.
    bool exists = contains(tree, "z");

    size_t count = find(tree, "z", indices, 2);

    bool missing = contains(tree, "a");

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(1, count);
    TEST_ASSERT_EQUAL_INT32(0, indices[0]);
    TEST_ASSERT_FALSE(missing);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that numeric digits inside the corpus are indexed and retrieved accurately.
 */
void test_suffix_tree_matches_numeric_substrings(void)
{
    // Arrange.
    const char* corpus = "build_2026_v1_2026_release";
    const char* pattern = "2026";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[27] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 27);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(2, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(14, indices[0]);
    TEST_ASSERT_EQUAL_INT32(6, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that whitespace and punctuation characters are preserved during edge matching.
 */
void test_suffix_tree_matches_whitespace_and_punctuation(void)
{
    // Arrange.
    const char* corpus = "item: [ok], item: [ok]";
    const char* pattern = ": [ok]";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[22] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 22);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(2, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(16, indices[0]);
    TEST_ASSERT_EQUAL_INT32(4, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a mismatch on the very last character of a pattern returns false.
 */
void test_suffix_tree_returns_false_when_final_character_mismatches(void)
{
    // Arrange.
    const char* corpus = "aaaaabaaaaab";
    const char* pattern = "aaaaac";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[12] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 12);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a pattern matching only a prefix at the end of the corpus returns false.
 */
void test_suffix_tree_returns_false_for_incomplete_suffix_at_end(void)
{
    // Arrange.
    const char* corpus = "prefix_ending_in_suf";
    const char* pattern = "suffix";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[22] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 22);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that control characters such as newline and tab are matched accurately.
 */
void test_suffix_tree_handles_newline_and_tab_sequences(void)
{
    // Arrange.
    const char* corpus = "line1\n\tval\nline2\n\tval";
    const char* pattern = "\n\tval";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[22] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 22);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(2, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(5, indices[0]);
    TEST_ASSERT_EQUAL_INT32(16, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that non-sentinel symbols such as '#' inside the corpus match accurately.
 */
void test_suffix_tree_handles_hash_symbol_inside_corpus(void)
{
    // Arrange.
    const char* corpus = "tag_#50_and_#50_total";
    const char* pattern = "#50";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[23] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 23);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(2, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(4, indices[0]);
    TEST_ASSERT_EQUAL_INT32(12, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that an empty pattern returns false for contains and zero for find.
 */
void test_suffix_tree_returns_false_for_empty_pattern(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";
    const char* empty_pattern = "";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[14] = { 0 };

    // Act.
    bool exists = contains(tree, empty_pattern);
    size_t number_of_occurrences = find(tree, empty_pattern, indices, 14);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that passing a NULL pattern to contains and find returns false and zero safely.
 */
void test_suffix_tree_returns_false_for_null_pattern(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";
    const char* null_pattern = NULL;

    tree_t* tree = initialize_tree(corpus);
    
    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[14] = { 0 };

    // Act.
    bool exists = contains(tree, null_pattern);
    size_t number_of_occurrences = find(tree, null_pattern, indices, 14);

    // Assert.
    TEST_ASSERT_FALSE(exists);
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that passing a NULL tree handle to contains and find fails safely.
 */
void test_suffix_tree_returns_false_when_tree_is_null(void)
{
    // Arrange.
    const tree_t* null_tree = NULL;

    const char* pattern = "ana";

    int32_t indices[4] = { 0 };

    // Act.
    bool exists = contains(null_tree, pattern);
    size_t number_of_occurrences = find(null_tree, pattern, indices, 4);

    // Assert.
    TEST_ASSERT_FALSE(exists);

    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);
}

/**
 * @brief Verifies that passing a NULL output array to find returns zero safely.
 */
void test_suffix_tree_find_returns_zero_when_output_buffer_is_null(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    // Act.
    size_t number_of_occurrences = find(tree, "ana", NULL, 10);

    // Assert.
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that passing a zero output capacity to find returns zero without writing indices.
 */
void test_suffix_tree_find_returns_zero_when_capacity_is_zero(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[2] = { -1, -1 };

    // Act.
    size_t number_of_occurrences = find(tree, "ana", indices, 0);

    // Assert.
    TEST_ASSERT_EQUAL_UINT64(0, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(-1, indices[0]);
    TEST_ASSERT_EQUAL_INT32(-1, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that find caps collected indices at the caller-provided buffer capacity.
 */
void test_suffix_tree_find_respects_smaller_output_capacity(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[2] = { -1, -1 };

    // Act.
    size_t number_of_occurrences = find(tree, "ana", indices, 2);

    // Assert.
    TEST_ASSERT_EQUAL_UINT64(2, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(11, indices[0]);
    TEST_ASSERT_EQUAL_INT32(3, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that an edge split followed by another split on the same prefix works cleanly.
 */
void test_suffix_tree_handles_cascading_edge_splits(void)
{
    // Arrange.
    const char* corpus = "pre_alpha_pre_beta_pre_gamma";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[8] = { 0 };

    // Act.
    bool exists = contains(tree, "pre_");
    size_t number_of_prefixes = find(tree, "pre_", indices, 8);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(3, number_of_prefixes);
    TEST_ASSERT_EQUAL_INT32(0, indices[0]);
    TEST_ASSERT_EQUAL_INT32(10, indices[1]);
    TEST_ASSERT_EQUAL_INT32(19, indices[2]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that distinct occurrences of "ban" in "banana_bandana" are located accurately.
 */
void test_suffix_tree_finds_ban_occurrences_in_banana_bandana(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";
    const char* pattern = "ban";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[14] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);

    size_t number_of_occurrences = find(tree, pattern, indices, 14);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(2, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(0, indices[0]);
    TEST_ASSERT_EQUAL_INT32(7, indices[1]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that a pattern appearing only once near the middle of a repeated string is isolated.
 */
void test_suffix_tree_finds_unique_substring_inside_repetitive_text(void)
{
    // Arrange.
    const char* corpus = "banana_bandana";
    const char* pattern = "band";

    tree_t* tree = initialize_tree(corpus);

    TEST_ASSERT_NOT_NULL(tree);

    int32_t indices[14] = { 0 };

    // Act.
    bool exists = contains(tree, pattern);
    size_t number_of_occurrences = find(tree, pattern, indices, 14);

    // Assert.
    TEST_ASSERT_TRUE(exists);
    TEST_ASSERT_EQUAL_UINT64(1, number_of_occurrences);
    TEST_ASSERT_EQUAL_INT32(7, indices[0]);

    // Cleanup.
    destroy(tree);
}

/**
 * @brief Verifies that calling destroy on a NULL tree pointer is a safe no-op.
 */
void test_suffix_tree_destroy_handles_null_pointer_safely(void)
{
    // Arrange.
    tree_t* null_tree = NULL;

    // Act.
    destroy(null_tree);

    // Assert.
    TEST_PASS();
}

int main(void)
{
    // Initiate the session.
    UNITY_BEGIN();

    // Register the 30 tests to run.
    RUN_TEST(test_suffix_tree_initializes_valid_corpus_successfully);
    RUN_TEST(test_suffix_tree_returns_null_when_corpus_is_null);
    RUN_TEST(test_suffix_tree_returns_null_when_corpus_is_empty);
    RUN_TEST(test_suffix_tree_finds_single_pattern_match_in_middle);
    RUN_TEST(test_suffix_tree_matches_pattern_at_start_of_corpus);
    RUN_TEST(test_suffix_tree_matches_pattern_at_end_of_corpus);
    RUN_TEST(test_suffix_tree_finds_multiple_occurrences_in_banana_bandana);
    RUN_TEST(test_suffix_tree_handles_overlapping_repeated_substrings);
    RUN_TEST(test_suffix_tree_returns_false_for_missing_pattern);
    RUN_TEST(test_suffix_tree_matches_entire_corpus_as_pattern);
    RUN_TEST(test_suffix_tree_rejects_pattern_longer_than_corpus);
    RUN_TEST(test_suffix_tree_enforces_case_sensitivity_during_search);
    RUN_TEST(test_suffix_tree_finds_all_single_character_occurrences);
    RUN_TEST(test_suffix_tree_handles_single_character_corpus);
    RUN_TEST(test_suffix_tree_matches_numeric_substrings);
    RUN_TEST(test_suffix_tree_matches_whitespace_and_punctuation);
    RUN_TEST(test_suffix_tree_returns_false_when_final_character_mismatches);
    RUN_TEST(test_suffix_tree_returns_false_for_incomplete_suffix_at_end);
    RUN_TEST(test_suffix_tree_handles_newline_and_tab_sequences);
    RUN_TEST(test_suffix_tree_handles_hash_symbol_inside_corpus);
    RUN_TEST(test_suffix_tree_returns_false_for_empty_pattern);
    RUN_TEST(test_suffix_tree_returns_false_for_null_pattern);
    RUN_TEST(test_suffix_tree_returns_false_when_tree_is_null);
    RUN_TEST(test_suffix_tree_find_returns_zero_when_output_buffer_is_null);
    RUN_TEST(test_suffix_tree_find_returns_zero_when_capacity_is_zero);
    RUN_TEST(test_suffix_tree_find_respects_smaller_output_capacity);
    RUN_TEST(test_suffix_tree_handles_cascading_edge_splits);
    RUN_TEST(test_suffix_tree_finds_ban_occurrences_in_banana_bandana);
    RUN_TEST(test_suffix_tree_finds_unique_substring_inside_repetitive_text);
    RUN_TEST(test_suffix_tree_destroy_handles_null_pointer_safely);

    // Destroy the session.
    return UNITY_END();
}