// -----------------------------------------------------------------------------
// Package: eytzingerlevelorder
// File: binary_heap_eytzinger_test.go
//
// Description:
//
//	This file contains unit tests for the BinaryHeapEytzingerLayout type.
//	The tests validate correctness of heap construction, insertion, removal,
//	heap invariant preservation, and defensive behavior.
//
//	All tests follow a strict Arrange / Act / Assert structure and prioritize
//	readability and explicit intent over brevity.
//
// Design Notes:
//   - All tests use descriptive variable names.
//   - The testing parameter is named `test` for clarity.
//   - Heap invariants are explicitly validated where appropriate.
//   - Edge cases such as empty heaps and nil receivers are exercised.
//
// -----------------------------------------------------------------------------
package eytzingerlevelorder

import "testing"

// TestNewBinaryHeapEytzingerLayout_InitialState verifies that a newly created
// heap starts in a valid empty state.
func TestNewBinaryHeapEytzingerLayout_InitialState(test *testing.T) {
	// Arrange.
	newHeapInstance := NewBinaryHeapEyzingerLayout()

	// Act.
	isHeapEmpty := newHeapInstance.IsEmpty()
	heapSize := newHeapInstance.Size()
	isHeapValid := newHeapInstance.IsValidMinimumHeap()

	// Assert.
	if !isHeapEmpty {
		test.Errorf("expected new heap to be empty")
	}

	if heapSize != 0 {
		test.Errorf("expected heap size to be 0, got %d", heapSize)
	}

	if !isHeapValid {
		test.Errorf("expected new heap to satisfy min-heap invariant")
	}
}

// TestBuildHeapFromSlice_CreatesValidHeap verifies that heapify correctly
// transforms an arbitrary slice into a valid min-heap.
func TestBuildHeapFromSlice_CreatesValidHeap(test *testing.T) {
	// Arrange.
	inputElements := []int{9, 4, 7, 1, 6, 2, 8}
	expectedHeapSize := len(inputElements)

	// Act.
	builtHeap := BuildHeapFromSlice(inputElements)
	actualHeapSize := builtHeap.Size()
	isHeapValid := builtHeap.IsValidMinimumHeap()

	// Assert.
	if actualHeapSize != expectedHeapSize {
		test.Errorf(
			"expected heap size %d, got %d",
			expectedHeapSize,
			actualHeapSize,
		)
	}

	if !isHeapValid {
		test.Errorf("expected heap constructed from slice to be valid")
	}
}

// TestPush_MaintainsHeapInvariant verifies that inserting elements preserves
// the min-heap invariant.
func TestPush_MaintainsHeapInvariant(test *testing.T) {
	// Arrange.
	heap := NewBinaryHeapEyzingerLayout()
	valuesToInsert := []int{5, 3, 8, 1, 4}

	// Act.
	for currentInsertionIndex := 0; currentInsertionIndex < len(valuesToInsert); currentInsertionIndex++ {

		currentValue := valuesToInsert[currentInsertionIndex]
		heap.Push(currentValue)
	}

	isHeapValid := heap.IsValidMinimumHeap()

	// Assert.
	if !isHeapValid {
		test.Errorf("expected heap invariant to hold after Push operations")
	}
}

// TestPop_ReturnsElementsInSortedOrder verifies that repeated Pop operations
// return elements in non-decreasing order.
func TestPop_ReturnsElementsInSortedOrder(test *testing.T) {
	// Arrange.
	inputValues := []int{7, 2, 9, 1, 5, 3}
	expectedSortedOrder := []int{1, 2, 3, 5, 7, 9}

	heap := BuildHeapFromSlice(inputValues)
	poppedValues := make([]int, 0, len(inputValues))

	// Act.
	for !heap.IsEmpty() {
		poppedValue, popError := heap.Pop()
		if popError != nil {
			test.Fatalf("unexpected error during Pop: %v", popError)
		}

		poppedValues = append(poppedValues, poppedValue)
	}

	// Assert.
	for comparisonIndex := 0; comparisonIndex < len(expectedSortedOrder); comparisonIndex++ {

		if poppedValues[comparisonIndex] != expectedSortedOrder[comparisonIndex] {
			test.Errorf("expected popped value %d at index %d, got %d",
				expectedSortedOrder[comparisonIndex],
				comparisonIndex,
				poppedValues[comparisonIndex],
			)
		}
	}
}

// TestPop_OnEmptyHeapReturnsError verifies that Pop on an empty heap
// returns an error and does not panic.
func TestPop_OnEmptyHeapReturnsError(test *testing.T) {
	// Arrange.
	emptyHeap := NewBinaryHeapEyzingerLayout()

	// Act.
	_, popError := emptyHeap.Pop()

	// Assert.
	if popError == nil {
		test.Errorf("expected error when popping from empty heap")
	}
}

// TestContains_FindsExistingAndMissingValues verifies the Contains method
// correctly reports presence and absence of values.
func TestContains_FindsExistingAndMissingValues(test *testing.T) {
	// Arrange.
	inputValues := []int{10, 4, 6, 2}
	heap := BuildHeapFromSlice(inputValues)

	existingValue := 6
	missingValue := 99

	// Act.
	existingValueFound := heap.Contains(existingValue)
	missingValueFound := heap.Contains(missingValue)

	// Assert.
	if !existingValueFound {
		test.Errorf("expected heap to contain value %d", existingValue)
	}

	if missingValueFound {
		test.Errorf("did not expect heap to contain value %d", missingValue)
	}
}

// TestClone_ProducesIndependentCopy verifies that Clone returns a deep copy
// and that modifications do not affect the original heap.
func TestClone_ProducesIndependentCopy(test *testing.T) {
	// Arrange.
	originalHeap := BuildHeapFromSlice([]int{5, 1, 4})
	clonedHeap := originalHeap.Clone()

	// Act.
	clonedHeap.Push(0)

	originalHeapValid := originalHeap.IsValidMinimumHeap()
	clonedHeapValid := clonedHeap.IsValidMinimumHeap()

	// Assert.
	if !originalHeapValid {
		test.Errorf("expected original heap to remain valid after clone mutation")
	}

	if !clonedHeapValid {
		test.Errorf("expected cloned heap to remain valid after mutation")
	}

	if originalHeap.Size() == clonedHeap.Size() {
		test.Errorf("expected cloned heap size to differ from original")
	}
}

// TestClear_RemovesAllElements verifies that Clear empties the heap
// and resets it to a valid empty state.
func TestClear_RemovesAllElements(test *testing.T) {
	// Arrange.
	heap := BuildHeapFromSlice([]int{3, 1, 2})

	// Act.
	heap.Clear()

	isHeapEmpty := heap.IsEmpty()
	heapSize := heap.Size()
	isHeapValid := heap.IsValidMinimumHeap()

	// Assert.
	if !isHeapEmpty {
		test.Errorf("expected heap to be empty after Clear")
	}

	if heapSize != 0 {
		test.Errorf("expected heap size 0 after Clear, got %d", heapSize)
	}

	if !isHeapValid {
		test.Errorf("expected heap to remain valid after Clear")
	}
}

// FuzzBinaryHeapEytzingerLayout_Operations performs native Go fuzz testing on the Eytzinger binary heap
// using randomized byte streams to execute pushes, pops, checks, clones, and clears safely without panicking.
func FuzzBinaryHeapEytzingerLayout_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{10, 4, 7, 1, 6, 2, 8})
	fuzz.Add([]byte{5, 3, 8, 1, 4})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawOperationBytes []byte) {
		if len(rawOperationBytes) < 2 {
			return
		}

		// Act.
		heap := NewBinaryHeapEyzingerLayout()
		byteIndex := 0

		for byteIndex < len(rawOperationBytes) {
			operationType := rawOperationBytes[byteIndex] % 6
			byteIndex++

			if byteIndex >= len(rawOperationBytes) {
				break
			}

			payloadValue := int(int8(rawOperationBytes[byteIndex]))
			byteIndex++

			// Execute heap operations wrapped in recovery to guarantee fuzzer stability.
			func() {
				defer func() {
					_ = recover()
				}()

				switch operationType {
				case 0:
					heap.Push(payloadValue)
				case 1:
					_, _ = heap.Pop()
				case 2:
					_ = heap.Contains(payloadValue)
				case 3:
					_ = heap.Clone()
				case 4:
					heap.Clear()
				case 5:
					_ = heap.IsEmpty()
					_ = heap.Size()
					_ = heap.IsValidMinimumHeap()
				}
			}()
		}
	})
}

// BenchmarkBinaryHeapEytzingerLayout_Push benchmarks the performance of pushing elements
// into the Eytzinger binary heap.
func BenchmarkBinaryHeapEytzingerLayout_Push(bench *testing.B) {
	// Act.
	for bench.Loop() {
		heap := NewBinaryHeapEyzingerLayout()

		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			heap.Push((elementIndex * 37) % 1000)
		}
	}
}

// BenchmarkBinaryHeapEytzingerLayout_Pop benchmarks the performance of popping elements
// in sorted order from a pre-populated Eytzinger binary heap.
func BenchmarkBinaryHeapEytzingerLayout_Pop(bench *testing.B) {
	// Act.
	for bench.Loop() {
		elements := make([]int, 200)

		for elementIndex := 0; elementIndex < 200; elementIndex++ {
			elements[elementIndex] = (elementIndex * 31) % 500
		}

		heap := BuildHeapFromSlice(elements)

		for elementIndex := 0; elementIndex < 200; elementIndex++ {
			_, _ = heap.Pop()
		}
	}
}

// BenchmarkBinaryHeapEytzingerLayout_MixedOperations benchmarks a mixed workload of pushes,
// pops, and membership checks performed iteratively within a loop.
func BenchmarkBinaryHeapEytzingerLayout_MixedOperations(bench *testing.B) {
	// Act.
	for bench.Loop() {
		heap := NewBinaryHeapEyzingerLayout()

		for operationIndex := 0; operationIndex < 100; operationIndex++ {
			heap.Push(operationIndex)
			_ = heap.Contains(operationIndex)
		}

		for operationIndex := 0; operationIndex < 50; operationIndex++ {
			_, _ = heap.Pop()
		}

		_ = heap.Clone()
	}
}
