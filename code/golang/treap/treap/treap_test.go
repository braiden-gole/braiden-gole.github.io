// -----------------------------------------------------------------------------
// Package: treap
// File: treap_test.go
//
// Description:
//
//	This file contains unit tests for the generic Treap data structure.
//	The tests validate functional correctness, order statistics, tree balance,
//	and strict maintenance of Binary Search Tree (BST) and Heap invariants.
//
//	The test suite focuses on:
//
//	  - Correct insertion and deletion operations.
//	  - Exact order-statistic queries (FindKth and Rank).
//	  - Structural invariance (BST, Max-Heap, and Subtree Size constraints).
//	  - Handling of edge cases such as empty states and duplicate keys.
//	  - Deterministic random engine behavior via custom seeds.
//
// Notes:
//   - All tests use explicit Arrange / Act / Assert structure.
//   - The testing parameter is intentionally named "test" for clarity.
//   - Descriptive variable names are used throughout.
//   - Structural assertions run recursively using helper functions.
//
// -----------------------------------------------------------------------------
package treap

import (
	"cmp"
	"math/rand/v2"
	"sort"
	"testing"
)

// verifyTreapInvariants performs structural assertions checking BST ordering,
// Max-Heap priority layout, and accurate subtree sizing.
func verifyTreapInvariants[KeyType cmp.Ordered](test *testing.T, targetTreap *Treap[KeyType]) {
	test.Helper()

	var validateSubtree func(node *Node[KeyType]) int

	validateSubtree = func(node *Node[KeyType]) int {
		if node == nil {
			return 0
		}

		leftSubtreeSize := validateSubtree(node.Left)
		rightSubtreeSize := validateSubtree(node.Right)

		if node.Left != nil && node.Left.Key > node.Key {
			test.Fatalf("BST invariant violation: left child key (%v) > parent key (%v)",
				node.Left.Key,
				node.Key,
			)
		}

		if node.Right != nil && node.Right.Key < node.Key {
			test.Fatalf("BST invariant violation: right child key (%v) < parent key (%v)",
				node.Right.Key,
				node.Key,
			)
		}

		if node.Left != nil && node.Left.Priority > node.Priority {
			test.Fatalf("Max-Heap invariant violation: left child priority (%d) > parent priority (%d)",
				node.Left.Priority,
				node.Priority,
			)
		}

		if node.Right != nil && node.Right.Priority > node.Priority {
			test.Fatalf("Max-Heap invariant violation: right child priority (%d) > parent priority (%d)",
				node.Right.Priority,
				node.Priority,
			)
		}

		expectedSize := 1 + leftSubtreeSize + rightSubtreeSize

		if node.Size != expectedSize {
			test.Fatalf("Subtree size violation: node (%v) size is %d, expected %d",
				node.Key,
				node.Size,
				expectedSize,
			)
		}

		return expectedSize
	}

	validateSubtree(targetTreap.Root)
}

// TestNewTreap verifies that initializing a Treap creates an empty structure.
func TestNewTreap(test *testing.T) {
	// Arrange.
	seedState := uint64(1)
	seedStream := uint64(2)

	// Act.
	treap := NewTreap[int](seedState, seedStream)

	// Assert.
	if treap.Root != nil {
		test.Fatalf("expected root node to be nil for new Treap")
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_InsertSingle verifies that inserting a single key creates a valid root.
func TestTreap_InsertSingle(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](10, 20)

	targetKey := 42

	// Act.
	treap.Insert(targetKey)

	// Assert.
	retrievedKey, isFound := treap.FindKth(0)

	if !isFound || retrievedKey != targetKey {
		test.Fatalf("expected key %d at rank 0, got %d (found: %v)",
			targetKey,
			retrievedKey,
			isFound,
		)
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_InsertSequential verifies that keys inserted in ascending order
// maintain sorted order statistics and structural invariants.
func TestTreap_InsertSequential(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](1, 1)

	numberOfKeys := 10

	// Act.
	for currentKey := 1; currentKey <= numberOfKeys; currentKey++ {
		treap.Insert(currentKey)
	}

	// Assert.
	for indexRank := 0; indexRank < numberOfKeys; indexRank++ {
		expectedValue := indexRank + 1
		retrievedKey, _ := treap.FindKth(indexRank)

		if retrievedKey != expectedValue {
			test.Fatalf("expected key %d at rank %d, got %d",
				expectedValue,
				indexRank,
				retrievedKey,
			)
		}
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_InsertReverseSequential verifies that keys inserted in descending order
// maintain proper internal ordering.
func TestTreap_InsertReverseSequential(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](5, 5)

	startingKey := 10

	// Act.
	for currentKey := startingKey; currentKey >= 1; currentKey-- {
		treap.Insert(currentKey)
	}

	// Assert.
	for indexRank := 0; indexRank < startingKey; indexRank++ {
		expectedValue := indexRank + 1
		retrievedKey, _ := treap.FindKth(indexRank)

		if retrievedKey != expectedValue {
			test.Fatalf("expected key %d at rank %d, got %d",
				expectedValue,
				indexRank,
				retrievedKey,
			)
		}
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_FindKthOutOfBounds verifies defensive return values for invalid ranks.
func TestTreap_FindKthOutOfBounds(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](1, 2)

	treap.Insert(10)
	treap.Insert(20)

	// Act.
	_, isNegativeFound := treap.FindKth(-1)
	_, isExceedingFound := treap.FindKth(2)

	// Assert.
	if isNegativeFound {
		test.Fatalf("expected isFound to be false for negative rank query")
	}

	if isExceedingFound {
		test.Fatalf("expected isFound to be false for rank query beyond total count")
	}
}

// TestTreap_RankExistingKeys verifies position reporting for keys present in the tree.
func TestTreap_RankExistingKeys(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](100, 200)

	keysToInsert := []int{50, 20, 80, 10, 30}

	for _, key := range keysToInsert {
		treap.Insert(key)
	}

	expectedRanks := map[int]int{
		10: 0,
		20: 1,
		30: 2,
		50: 3,
		80: 4,
	}

	// Act & Assert.
	for key, expectedRank := range expectedRanks {
		calculatedRank := treap.Rank(key)

		if calculatedRank != expectedRank {
			test.Fatalf("expected rank for key %d to be %d, got %d",
				key,
				expectedRank,
				calculatedRank,
			)
		}
	}
}

// TestTreap_RankNonExistingKeys verifies relative rank queries for keys not in the tree.
func TestTreap_RankNonExistingKeys(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](42, 84)

	treap.Insert(10)
	treap.Insert(30)
	treap.Insert(50)

	// Act & Assert.
	if rank := treap.Rank(5); rank != 0 {
		test.Fatalf("expected rank 0 for key below minimum, got %d", rank)
	}

	if rank := treap.Rank(25); rank != 1 {
		test.Fatalf("expected rank 1 for key between 10 and 30, got %d", rank)
	}

	if rank := treap.Rank(60); rank != 3 {
		test.Fatalf("expected rank 3 for key above maximum, got %d", rank)
	}
}

// TestTreap_DeleteLeaf verifies removing a leaf node correctly updates the tree.
func TestTreap_DeleteLeaf(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](1, 2)

	treap.Insert(10)
	treap.Insert(20)

	// Act.
	treap.Delete(20)

	// Assert.
	if _, isFound := treap.FindKth(1); isFound {
		test.Fatalf("expected second node to be removed from tree")
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_DeleteRoot verifies deleting the current root restructures correctly.
func TestTreap_DeleteRoot(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](7, 14)

	treap.Insert(20)
	treap.Insert(10)
	treap.Insert(30)

	rootKeyBeforeDeletion := treap.Root.Key

	// Act.
	treap.Delete(rootKeyBeforeDeletion)

	// Assert.
	if treap.Rank(10) != 0 || treap.Rank(30) != 1 {
		test.Fatalf("expected remaining keys 10 and 30 to hold ranks 0 and 1 after root deletion")
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_DeleteNonExistentKey verifies that removing a missing key leaves the tree intact.
func TestTreap_DeleteNonExistentKey(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](9, 9)

	treap.Insert(10)
	treap.Insert(20)

	nonExistentKey := 999

	// Act.
	treap.Delete(nonExistentKey)

	// Assert.
	retrievedKey, _ := treap.FindKth(0)

	if retrievedKey != 10 {
		test.Fatalf("expected first element to remain 10, got %d", retrievedKey)
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_DeleteFromEmptyTreap verifies safety when attempting deletion on an empty tree.
func TestTreap_DeleteFromEmptyTreap(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](0, 0)

	// Act.
	treap.Delete(5)

	// Assert.
	if treap.Root != nil {
		test.Fatalf("expected root to remain nil after deleting from empty tree")
	}
}

// TestTreap_DeleteAll verifies sequential deletion until the tree is empty.
func TestTreap_DeleteAll(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](123, 456)

	elementsToInsert := []int{15, 8, 23, 42, 4}

	for _, element := range elementsToInsert {
		treap.Insert(element)
	}

	// Act.
	for _, element := range elementsToInsert {
		treap.Delete(element)
		verifyTreapInvariants(test, treap)
	}

	// Assert.
	if treap.Root != nil {
		test.Fatalf("expected treap to be completely empty after removing all elements")
	}
}

// TestTreap_DuplicateInsertions verifies handling duplicate key values.
func TestTreap_DuplicateInsertions(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](11, 22)

	duplicateKey := 10

	// Act.
	treap.Insert(duplicateKey)
	treap.Insert(duplicateKey)
	treap.Insert(duplicateKey)

	// Assert.
	if treap.Root.Size != 3 {
		test.Fatalf("expected total size 3 for duplicate keys, got %d",
			treap.Root.Size,
		)
	}

	treap.Delete(duplicateKey)

	if treap.Root.Size != 2 {
		test.Fatalf("expected total size 2 after single instance deletion, got %d",
			treap.Root.Size,
		)
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_StringKeys verifies generic support for ordered string types.
func TestTreap_StringKeys(test *testing.T) {
	// Arrange.
	treap := NewTreap[string](55, 66)

	wordList := []string{"orange", "apple", "banana"}

	// Act.
	for _, word := range wordList {
		treap.Insert(word)
	}

	// Assert.
	retrievedWord, _ := treap.FindKth(0)

	if retrievedWord != "apple" {
		test.Fatalf("expected 'apple' at rank 0, got '%s'", retrievedWord)
	}

	calculatedRank := treap.Rank("banana")

	if calculatedRank != 1 {
		test.Fatalf("expected rank 1 for 'banana', got %d", calculatedRank)
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_Float64Keys verifies generic support for floating point types.
func TestTreap_Float64Keys(test *testing.T) {
	// Arrange.
	treap := NewTreap[float64](3, 9)

	floatList := []float64{3.14, 1.41, 2.71, 0.57}

	// Act.
	for _, floatValue := range floatList {
		treap.Insert(floatValue)
	}

	// Assert.
	retrievedValue, _ := treap.FindKth(0)

	if retrievedValue != 0.57 {
		test.Fatalf("expected 0.57 at rank 0, got %f", retrievedValue)
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_Split verifies the Split primitive divides keys correctly relative to target threshold.
func TestTreap_Split(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](1, 2)

	for value := 10; value <= 50; value += 10 {
		treap.Insert(value)
	}

	splitThreshold := 30

	// Act.
	leftSubtree, rightSubtree := Split(treap.Root, splitThreshold)

	// Assert.
	if leftSubtree == nil || rightSubtree == nil {
		test.Fatalf("expected non-nil subtrees after Split operation")
	}

	if leftSubtree.Size != 3 || rightSubtree.Size != 2 {
		test.Fatalf("expected left size 3 and right size 2, got left=%d, right=%d",
			leftSubtree.Size,
			rightSubtree.Size,
		)
	}
}

// TestTreap_Merge verifies combining two disjoint treap structures.
func TestTreap_Merge(test *testing.T) {
	// Arrange.
	firstTreap := NewTreap[int](1, 2)

	firstTreap.Insert(10)
	firstTreap.Insert(20)

	secondTreap := NewTreap[int](3, 4)

	secondTreap.Insert(30)
	secondTreap.Insert(40)

	// Act.
	mergedRoot := Merge(firstTreap.Root, secondTreap.Root)

	// Assert.
	if mergedRoot.Size != 4 {
		test.Fatalf("expected merged root size 4, got %d", mergedRoot.Size)
	}
}

// TestTreap_StressRandomDataset verifies correctness against standard slice sorting across 1000 items.
func TestTreap_StressRandomDataset(test *testing.T) {
	// Arrange.
	treap := NewTreap[int](999, 888)

	elementCount := 1000
	insertedElements := make([]int, elementCount)

	for index := 0; index < elementCount; index++ {
		randomValue := rand.IntN(10000)
		insertedElements[index] = randomValue
		treap.Insert(randomValue)
	}

	verifyTreapInvariants(test, treap)

	// Act.
	sort.Ints(insertedElements)

	// Assert.
	for indexRank := 0; indexRank < elementCount; indexRank++ {
		retrievedKey, _ := treap.FindKth(indexRank)

		if retrievedKey != insertedElements[indexRank] {
			test.Fatalf("mismatch at rank %d: expected %d, got %d",
				indexRank,
				insertedElements[indexRank],
				retrievedKey,
			)
		}
	}

	for index := 0; index < elementCount/2; index++ {
		treap.Delete(insertedElements[index])
	}

	verifyTreapInvariants(test, treap)
}

// TestTreap_DeterministicSeed verifies that two instances with matching seed configurations produce identical structures.
func TestTreap_DeterministicSeed(test *testing.T) {
	// Arrange.
	stateSeed := uint64(12345)
	streamSeed := uint64(67890)

	firstTreap := NewTreap[int](stateSeed, streamSeed)
	secondTreap := NewTreap[int](stateSeed, streamSeed)

	insertionValues := []int{45, 12, 89, 33, 27, 91, 5, 62}

	// Act.
	for _, value := range insertionValues {
		firstTreap.Insert(value)
		secondTreap.Insert(value)
	}

	// Assert.
	var compareNodes func(source, target *Node[int]) bool

	compareNodes = func(source, target *Node[int]) bool {
		if source == nil && target == nil {
			return true
		}

		if source == nil || target == nil {
			return false
		}

		return source.Key == target.Key &&
			source.Priority == target.Priority &&
			compareNodes(source.Left, target.Left) &&
			compareNodes(source.Right, target.Right)
	}

	if !compareNodes(firstTreap.Root, secondTreap.Root) {
		test.Fatalf("expected identical internal structures for identical seeds")
	}
}

// FuzzTreap_Operations performs native Go fuzz testing on the Treap implementation
// using randomized byte streams to execute insertions, deletions, rank searches, and invariant checks safely.
func FuzzTreap_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 10, 1, 20, 2, 10, 0, 30})
	fuzz.Add([]byte{0, 5, 0, 15, 1, 5, 2, 15})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawOperationBytes []byte) {
		if len(rawOperationBytes) < 2 {
			return
		}

		// Act.
		treap := NewTreap[int](42, 84)
		byteIndex := 0

		for byteIndex+1 < len(rawOperationBytes) {
			operationType := rawOperationBytes[byteIndex] % 3
			byteIndex++

			payloadValue := int(rawOperationBytes[byteIndex])
			byteIndex++

			// Execute treap operations wrapped in recovery to guarantee fuzzer stability.
			func() {
				defer func() {
					_ = recover()
				}()

				switch operationType {
				case 0:
					treap.Insert(payloadValue)
				case 1:
					treap.Delete(payloadValue)
				case 2:
					_, _ = treap.FindKth(payloadValue % 10)
				}
			}()
		}
	})
}

// BenchmarkTreap_Insert benchmarks the performance of inserting keys into the treap.
func BenchmarkTreap_Insert(bench *testing.B) {
	// Act.
	for bench.Loop() {
		treap := NewTreap[int](1, 2)

		for elementIndex := 0; elementIndex < 300; elementIndex++ {
			treap.Insert(elementIndex)
		}
	}
}

// BenchmarkTreap_FindKth benchmarks the performance of order-statistic queries
// across a pre-populated treap structure.
func BenchmarkTreap_FindKth(bench *testing.B) {
	// Arrange.
	treap := NewTreap[int](5, 10)

	for elementIndex := 0; elementIndex < 300; elementIndex++ {
		treap.Insert(elementIndex)
	}

	// Act.
	for bench.Loop() {
		for rankIndex := 0; rankIndex < 300; rankIndex += 10 {
			_, _ = treap.FindKth(rankIndex % 300)
		}
	}
}

// BenchmarkTreap_MixedWorkload benchmarks a mixed workload of insertions,
// deletions, and rank queries executed iteratively on the treap using a clean switch statement.
func BenchmarkTreap_MixedWorkload(bench *testing.B) {
	// Act.
	for bench.Loop() {
		treap := NewTreap[int](99, 101)

		for operationIndex := 0; operationIndex < 200; operationIndex++ {
			switch operationIndex % 3 {
			case 0:
				treap.Insert(operationIndex)
			case 1:
				treap.Delete(operationIndex % 50)
			case 2:
				_ = treap.Rank(operationIndex % 100)
			}
		}
	}
}
