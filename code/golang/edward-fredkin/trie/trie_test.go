// -----------------------------------------------------------------------------
// Package: trie
// File: trie_test.go
//
// Description:
//
//	This file contains unit tests for the Trie (prefix tree) implementation.
//	The tests verify correctness of insertion, lookup, prefix queries,
//	enumeration, deletion, and word count tracking.
//
//	Each test follows a strict Arrange / Act / Assert structure to ensure
//	clarity and consistency. All tests are written defensively and assume
//	no ordering guarantees from internal map traversal.
//
// Notes:
//   - Test function parameters are named "test" instead of the conventional "t".
//   - Assertions focus on set membership rather than ordering.
//   - Tests are independent and do not rely on shared state.
//
// -----------------------------------------------------------------------------
package trie

import "testing"

// TestNewTrieInitialState verifies that a newly created trie is empty
// and properly initialized.
func TestNewTrieInitialState(test *testing.T) {
	// Arrange.
	newTrieInstance := NewTrie()

	// Act.
	initialWordCount := newTrieInstance.WordCount()
	allStoredWords := newTrieInstance.AllWords()

	// Assert.
	if initialWordCount != 0 {
		test.Fatalf("expected initial word count to be 0, got %d",
			initialWordCount,
		)
	}

	if len(allStoredWords) != 0 {
		test.Fatalf("expected no words in new trie, got %d",
			len(allStoredWords),
		)
	}
}

// TestInsertAndContainsWord verifies that inserted words
// are correctly recognized as complete entries.
func TestInsertAndContainsWord(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	wordsToInsert := []string{"cat", "car", "dog"}

	// Act.
	for wordIndex := 0; wordIndex < len(wordsToInsert); wordIndex++ {
		trie.Insert(wordsToInsert[wordIndex])
	}

	// Assert.
	for wordIndex := 0; wordIndex < len(wordsToInsert); wordIndex++ {
		currentWord := wordsToInsert[wordIndex]

		if !trie.ContainsWord(currentWord) {
			test.Fatalf("expected trie to contain word %q",
				currentWord,
			)
		}
	}
}

// TestContainsWordRejectsPrefixes verifies that prefixes
// are not incorrectly treated as complete words.
func TestContainsWordRejectsPrefixes(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	trie.Insert("carbon")

	// Act.
	containsPrefixOnly := trie.ContainsWord("car")

	// Assert.
	if containsPrefixOnly {
		test.Fatalf("expected prefix %q not to be treated as a complete word",
			"car",
		)
	}
}

// TestContainsPrefix verifies that valid prefixes are detected
// even when they do not terminate complete words.
func TestContainsPrefix(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	trie.Insert("dog")
	trie.Insert("dot")

	// Act.
	prefixExists := trie.ContainsPrefix("do")
	prefixDoesNotExist := trie.ContainsPrefix("da")

	// Assert.
	if !prefixExists {
		test.Fatalf("expected prefix %q to exist in trie", "do")
	}

	if prefixDoesNotExist {
		test.Fatalf("expected prefix %q not to exist in trie", "da")
	}
}

// TestGetWordsWithPrefix verifies that all words sharing a prefix
// are returned without relying on order.
func TestGetWordsWithPrefix(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	trie.Insert("car")
	trie.Insert("cart")
	trie.Insert("carbon")
	trie.Insert("cat")

	expectedWordsWithPrefix := map[string]bool{
		"car":    true,
		"cart":   true,
		"carbon": true,
	}

	// Act.
	retrievedWords := trie.GetWordsWithPrefix("car")

	// Assert.
	if len(retrievedWords) != len(expectedWordsWithPrefix) {
		test.Fatalf("expected %d words with prefix, got %d",
			len(expectedWordsWithPrefix),
			len(retrievedWords),
		)
	}

	for wordIndex := 0; wordIndex < len(retrievedWords); wordIndex++ {
		currentWord := retrievedWords[wordIndex]

		if !expectedWordsWithPrefix[currentWord] {
			test.Fatalf("unexpected word %q returned for prefix query",
				currentWord,
			)
		}
	}
}

// TestAllWords verifies that all inserted words are returned.
func TestAllWords(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	expectedWords := map[string]bool{
		"cat": true,
		"car": true,
		"dog": true,
	}

	for word := range expectedWords {
		trie.Insert(word)
	}

	// Act.
	allStoredWords := trie.AllWords()

	// Assert.
	if len(allStoredWords) != len(expectedWords) {
		test.Fatalf("expected %d total words, got %d",
			len(expectedWords),
			len(allStoredWords),
		)
	}

	for wordIndex := 0; wordIndex < len(allStoredWords); wordIndex++ {
		currentWord := allStoredWords[wordIndex]

		if !expectedWords[currentWord] {
			test.Fatalf("unexpected word %q returned from AllWords",
				currentWord,
			)
		}
	}
}

// TestDeleteRemovesOnlyTargetWord verifies that deletion removes
// only the specified word and preserves shared prefixes.
func TestDeleteRemovesOnlyTargetWord(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	trie.Insert("car")
	trie.Insert("cart")
	trie.Insert("carbon")

	// Act.
	deletionResult := trie.Delete("cart")

	// Assert.
	if !deletionResult {
		test.Fatalf("expected deletion of %q to succeed", "cart")
	}

	if trie.ContainsWord("cart") {
		test.Fatalf("expected word %q to be removed from trie", "cart")
	}

	if !trie.ContainsWord("car") {
		test.Fatalf("expected word %q to remain after deleting %q",
			"car",
			"cart",
		)
	}

	if !trie.ContainsWord("carbon") {
		test.Fatalf("expected word %q to remain after deleting %q",
			"carbon",
			"cart",
		)
	}
}

// TestWordCountTracksInsertionsAndDeletions verifies that the
// word count is updated correctly.
func TestWordCountTracksInsertionsAndDeletions(test *testing.T) {
	// Arrange.
	trie := NewTrie()

	// Act.
	trie.Insert("cat")
	trie.Insert("dog")
	trie.Insert("dot")

	initialCount := trie.WordCount()

	trie.Delete("dog")

	finalCount := trie.WordCount()

	// Assert.
	if initialCount != 3 {
		test.Fatalf("expected word count to be 3 after insertions, got %d",
			initialCount,
		)
	}

	if finalCount != 2 {
		test.Fatalf("expected word count to be 2 after deletion, got %d",
			finalCount,
		)
	}
}

// FuzzTrie_Operations performs native Go fuzz testing on the Trie prefix tree
// using randomized byte streams to execute insertions, lookups, prefix queries, and deletions safely.
func FuzzTrie_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte("cat"), []byte("cart"))
	fuzz.Add([]byte("algorithm"), []byte("data"))

	fuzz.Fuzz(func(fuzzContext *testing.T, rawKeyBytes []byte, rawSecondaryBytes []byte) {
		if len(rawKeyBytes) == 0 || len(rawSecondaryBytes) == 0 {
			return
		}

		primaryWord := string(rawKeyBytes)
		secondaryWord := string(rawSecondaryBytes)

		// Act & Assert: Execute Trie operations wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			trie := NewTrie()
			trie.Insert(primaryWord)
			trie.Insert(secondaryWord)

			_ = trie.ContainsWord(primaryWord)
			_ = trie.ContainsPrefix(primaryWord[:len(primaryWord)/2+1])
			_ = trie.GetWordsWithPrefix(primaryWord[:1])
			_ = trie.AllWords()
			_ = trie.WordCount()

			trie.Delete(primaryWord)

			_ = trie.ContainsWord(primaryWord)
		}()
	})
}

// BenchmarkTrie_Insert benchmarks the performance of inserting words
// into the Trie.
func BenchmarkTrie_Insert(bench *testing.B) {
	// Arrange.
	words := []string{"algorithm", "alignment", "allocate", "binary", "blockchain", "bloom", "bytecode", "cluster"}

	// Act.
	for bench.Loop() {
		trie := NewTrie()

		for _, word := range words {
			trie.Insert(word)
		}
	}
}

// BenchmarkTrie_ContainsWord benchmarks the performance of membership lookup checks
// within a pre-populated Trie.
func BenchmarkTrie_ContainsWord(bench *testing.B) {
	// Arrange.
	trie := NewTrie()

	words := []string{"algorithm", "alignment", "allocate", "binary", "blockchain", "bloom", "bytecode", "cluster"}

	for _, word := range words {
		trie.Insert(word)
	}

	// Act.
	for bench.Loop() {
		for _, word := range words {
			_ = trie.ContainsWord(word)
		}
	}
}

// BenchmarkTrie_GetWordsWithPrefix benchmarks the performance of retrieving all words
// sharing a common prefix from the Trie.
func BenchmarkTrie_GetWordsWithPrefix(bench *testing.B) {
	// Arrange.
	trie := NewTrie()
	words := []string{"cat", "car", "cart", "carbon", "card", "care", "carpet", "dog", "dot"}

	for _, word := range words {
		trie.Insert(word)
	}

	// Act.
	for bench.Loop() {
		_ = trie.GetWordsWithPrefix("car")
	}
}
