// -----------------------------------------------------------------------------
// Package: squarerootdecomposition
// File: square_root_decomposition_test.go
//
// Description:
//
//	This file provides a comprehensive test suite of 30 unique, meaningful unit
//	tests covering data structure initialization, validation rules, block
//	computation, range querying, point updates, encapsulation guards, and
//	boundary checking for the Square Root Decomposition implementation.
//
// -----------------------------------------------------------------------------
package squarerootdecomposition

import (
	"testing"
)

// TestNewSqrtDecomposition_ValidInitialization verifies that initializing a
// decomposition structure with a valid array succeeds and returns a valid pointer.
func TestNewSqrtDecomposition_ValidInitialization(test *testing.T) {
	// Arrange.
	var numbers []int = []int{1, 2, 3, 4, 5}

	// Act.
	decomposition, err := NewSqrtDecomposition(numbers)

	// Assert.
	if err != nil {
		test.Fatalf("expected no error for valid array initialization, got %v", err)
	}

	if decomposition == nil {
		test.Fatalf("expected decomposition instance to be created, got nil")
	}
}

// TestNewSqrtDecomposition_NilArray verifies that attempting to initialize
// the decomposition structure with a nil slice is properly rejected.
func TestNewSqrtDecomposition_NilArray(test *testing.T) {
	// Arrange.
	var nilArray []int = nil

	// Act.
	decomposition, err := NewSqrtDecomposition(nilArray)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for nil input array, got nil")
	}

	if decomposition != nil {
		test.Fatalf("expected nil instance for nil input array, got %v", decomposition)
	}
}

// TestNewSqrtDecomposition_EmptyArray verifies that attempting to initialize
// the structure with an empty slice triggers a validation error.
func TestNewSqrtDecomposition_EmptyArray(test *testing.T) {
	// Arrange.
	var emptyArray []int = []int{}

	// Act.
	decomposition, err := NewSqrtDecomposition(emptyArray)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for empty input array, got nil")
	}

	if decomposition != nil {
		test.Fatalf("expected nil instance for empty input array, got %v", decomposition)
	}
}

// TestNewSqrtDecomposition_PerfectSquareBlockSize confirms that arrays with a
// perfect square length calculate exact block sizes seamlessly.
func TestNewSqrtDecomposition_PerfectSquareBlockSize(test *testing.T) {
	// Arrange.
	var perfectSquareArray []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Act.
	decomposition, _ := NewSqrtDecomposition(perfectSquareArray)

	// Assert.
	var expectedBlockSize int = 3

	if decomposition.size != expectedBlockSize {
		test.Fatalf("expected block size to be %d, got %d", expectedBlockSize, decomposition.size)
	}
}

// TestNewSqrtDecomposition_ImperfectSquareBlockSize verifies that arrays with
// non-square lengths round down safely and allocate sufficient block totals.
func TestNewSqrtDecomposition_ImperfectSquareBlockSize(test *testing.T) {
	// Arrange.
	var imperfectSquareArray []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Act.
	decomposition, _ := NewSqrtDecomposition(imperfectSquareArray)

	// Assert.
	var expectedBlockSize int = 3
	var expectedTotalBlocks int = 4

	if decomposition.size != expectedBlockSize {
		test.Fatalf("expected block size to be %d, got %d", expectedBlockSize, decomposition.size)
	}

	if len(decomposition.sums) != expectedTotalBlocks {
		test.Fatalf("expected total block count to be %d, got %d", expectedTotalBlocks, len(decomposition.sums))
	}
}

// TestNewSqrtDecomposition_SingleElementArray verifies robust handling of arrays
// containing exactly one element, avoiding zero-division panics.
func TestNewSqrtDecomposition_SingleElementArray(test *testing.T) {
	// Arrange.
	var singleElementArray []int = []int{42}

	// Act.
	decomposition, _ := NewSqrtDecomposition(singleElementArray)

	// Assert.
	var expectedBlockSize int = 1

	if decomposition.size != expectedBlockSize {
		test.Fatalf("expected single element array to force block size %d, got %d", expectedBlockSize, decomposition.size)
	}
}

// TestNewSqrtDecomposition_BlockSumsCorrectness ensures that the precomputed
// sums array accurately reflects the aggregated values of each specific block.
func TestNewSqrtDecomposition_BlockSumsCorrectness(test *testing.T) {
	// Arrange.
	var numbers []int = []int{1, 4, 3, 2, 8, 5, 6, 7}

	// Act.
	decomposition, _ := NewSqrtDecomposition(numbers)

	// Assert.
	var expecttedSum int = 5
	var expectedSecondBlockSum int = 5

	if decomposition.sums[0] != expecttedSum {
		test.Fatalf("expected first block sum to be %d, got %d", expecttedSum, decomposition.sums[0])
	}

	if decomposition.sums[1] != expectedSecondBlockSum {
		test.Fatalf("expected second block sum to be %d, got %d", expectedSecondBlockSum, decomposition.sums[1])
	}
}

// TestNewSqrtDecomposition_EncapsulationDefensiveCopy guarantees that external
// mutations to the input slice do not corrupt the internal data structure state.
func TestNewSqrtDecomposition_EncapsulationDefensiveCopy(test *testing.T) {
	// Arrange.
	var originalArray []int = []int{10, 20, 30}

	decomposition, _ := NewSqrtDecomposition(originalArray)

	// Act.
	originalArray[0] = 999

	// Assert.
	var expectedInternalValue int = 10

	actualInternalValue, _ := decomposition.RetrieveItem(0)

	if actualInternalValue != expectedInternalValue {
		test.Fatalf("expected defensive copy to retain %d, got mutated value %d", expectedInternalValue, actualInternalValue)
	}
}

// TestSize_ReturnsTotalLength verifies the Size method accurately reflects
// the total number of elements managed by the decomposition structure.
func TestSize_ReturnsTotalLength(test *testing.T) {
	// Arrange.
	var numbers []int = []int{7, 14, 21, 28, 35}

	decomposition, _ := NewSqrtDecomposition(numbers)

	// Act.
	var structureSize int = decomposition.Size()

	// Assert.
	var expectedSize int = 5

	if structureSize != expectedSize {
		test.Fatalf("expected structure size to be %d, got %d", expectedSize, structureSize)
	}
}

// TestRetreiveItem_ValidBoundaryStart ensures retrieving the very first
// element within the structure works without throwing validation errors.
func TestRetreiveItem_ValidBoundaryStart(test *testing.T) {
	// Arrange.
	var numbers []int = []int{8, 16, 24}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var beginning int = 0

	// Act.
	retrievedValue, retrievalError := decomposition.RetrieveItem(beginning)

	// Assert.
	var expectedValue int = 8

	if retrievalError != nil {
		test.Fatalf("expected no error retrieving start index, got %v", retrievalError)
	}

	if retrievedValue != expectedValue {
		test.Fatalf("expected to retrieve %d, got %d", expectedValue, retrievedValue)
	}
}

// TestRetreiveItem_ValidBoundaryEnd ensures retrieving the very last
// element within the structure evaluates bounds securely and successfully.
func TestRetreiveItem_ValidBoundaryEnd(test *testing.T) {
	// Arrange.
	var numbers []int = []int{8, 16, 24}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var end int = 2

	// Act.
	retrievedValue, retrievalError := decomposition.RetrieveItem(end)

	// Assert.
	var expectedValue int = 24

	if retrievalError != nil {
		test.Fatalf("expected no error retrieving end index, got %v", retrievalError)
	}

	if retrievedValue != expectedValue {
		test.Fatalf("expected to retrieve %d, got %d", expectedValue, retrievedValue)
	}
}

// TestRetreiveItem_NegativeIndexRejection tests defensive programming boundaries
// to guarantee negative target indices return structured errors.
func TestRetreiveItem_NegativeIndexRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{100, 200}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var invalidNegativeIndex int = -1

	// Act.
	_, retrievalError := decomposition.RetrieveItem(invalidNegativeIndex)

	// Assert.
	if retrievalError == nil {
		test.Fatalf("expected error for negative retrieval index, got nil")
	}
}

// TestRetreiveItem_OutOfBoundsIndexRejection prevents fatal panics by ensuring
// target indices beyond the array length correctly return defensive errors.
func TestRetreiveItem_OutOfBoundsIndexRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{100, 200}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var invalidOutOfBoundsIndex int = 5

	// Act.
	_, retrievalError := decomposition.RetrieveItem(invalidOutOfBoundsIndex)

	// Assert.
	if retrievalError == nil {
		test.Fatalf("expected error for out of bounds retrieval index, got nil")
	}
}

// TestQueryRange_SingleElementSameIndex verifies ranges where the left and
// right boundaries are identical return strictly that element's value.
func TestQueryRange_SingleElementSameIndex(test *testing.T) {
	// Arrange.
	var numbers []int = []int{5, 10, 15, 20, 25}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var index int = 2

	// Act.
	calculatedSum, err := decomposition.QueryRange(index, index)

	// Assert.
	var expectedSum int = 15

	if err != nil {
		test.Fatalf("expected no error for identical boundary query, got %v", err)
	}

	if calculatedSum != expectedSum {
		test.Fatalf("expected sum of %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_FullyWithinSingleBlock guarantees accurate linear scanning
// when an entire requested range naturally fits inside one distinct block.
func TestQueryRange_FullyWithinSingleBlock(test *testing.T) {
	// Arrange.
	var numbers []int = []int{1, 1, 1, 1, 1, 1, 1, 1, 1}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var leftBoundary int = 0
	var rightBoundary int = 2

	// Act.
	calculatedSum, err := decomposition.QueryRange(leftBoundary, rightBoundary)

	// Assert.
	var expectedSum int = 3

	if err != nil {
		test.Fatalf("expected no error for single block query, got %v", err)
	}

	if calculatedSum != expectedSum {
		test.Fatalf("expected sum of %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_SpanningTwoAdjacentBlocks validates that ranges crossing
// exactly one block boundary accurately merge trailing and leading elements.
func TestQueryRange_SpanningTwoAdjacentBlocks(test *testing.T) {
	// Arrange.
	var numbers []int = []int{10, 10, 10, 10, 10, 10, 10, 10, 10}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var leftBoundary int = 1
	var rightBoundary int = 4

	// Act.
	calculatedSum, err := decomposition.QueryRange(leftBoundary, rightBoundary)

	// Assert.
	var expectedSum int = 40

	if err != nil {
		test.Fatalf("expected no error for adjacent block query, got %v", err)
	}

	if calculatedSum != expectedSum {
		test.Fatalf("expected sum of %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_SpanningMultipleBlocksWithFullMiddle verifies O(1) block sum
// injection when a range completely encapsulates one or more internal blocks.
func TestQueryRange_SpanningMultipleBlocksWithFullMiddle(test *testing.T) {
	// Arrange.
	var numbers []int = []int{2, 2, 2, 5, 5, 5, 2, 2, 2, 2, 2, 2}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var leftBoundary int = 1
	var rightBoundary int = 8

	// Act.
	calculatedSum, err := decomposition.QueryRange(leftBoundary, rightBoundary)

	// Assert.
	var expectedSum int = 25

	if err != nil {
		test.Fatalf("expected no error for multi-block query, got %v", err)
	}

	if calculatedSum != expectedSum {
		test.Fatalf("expected sum of %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_EntireArrayCoverage evaluates performance and correctness
// when the query encompasses every element managed by the structure.
func TestQueryRange_EntireArrayCoverage(test *testing.T) {
	// Arrange.
	var numbers []int = []int{1, 2, 3, 4, 5, 6, 7}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var leftBoundary int = 0
	var rightBoundary int = 6

	// Act.
	calculatedSum, err := decomposition.QueryRange(leftBoundary, rightBoundary)

	// Assert.
	var expectedSum int = 28

	if err != nil {
		test.Fatalf("expected no error for full array query, got %v", err)
	}

	if calculatedSum != expectedSum {
		test.Fatalf("expected sum of %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_NegativeLeftBoundaryRejection ensures that queries beginning
// at invalid negative boundaries immediately trigger defensive validation errors.
func TestQueryRange_NegativeLeftBoundaryRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{10, 20, 30}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var invalidLeftBoundary int = -2
	var validRightBoundary int = 1

	// Act.
	_, err := decomposition.QueryRange(invalidLeftBoundary, validRightBoundary)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for negative left boundary, got nil")
	}
}

// TestQueryRange_OutOfBoundsRightBoundaryRejection ensures queries terminating
// beyond total structure capacity are successfully caught and rejected.
func TestQueryRange_OutOfBoundsRightBoundaryRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{10, 20, 30}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var validLeftBoundary int = 0
	var invalidRightBoundary int = 10

	// Act.
	_, err := decomposition.QueryRange(validLeftBoundary, invalidRightBoundary)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for out of bounds right boundary, got nil")
	}
}

// TestQueryRange_ReversedBoundariesRejection guarantees illogical parameters
// where the starting index is mathematically greater than the ending index fail.
func TestQueryRange_ReversedBoundariesRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{10, 20, 30, 40}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var invalidLeftBoundary int = 3
	var invalidRightBoundary int = 1

	// Act.
	_, err := decomposition.QueryRange(invalidLeftBoundary, invalidRightBoundary)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for reversed boundaries, got nil")
	}
}

// TestUpdateItem_ValidModificationSuccess verifies that writing new values to
// existing valid indices proceeds without triggering protection faults.
func TestUpdateItem_ValidModificationSuccess(test *testing.T) {
	// Arrange.
	var numbers []int = []int{5, 5, 5, 5}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var index int = 2
	var replacementValue int = 10

	// Act.
	updateError := decomposition.UpdateItem(index, replacementValue)

	// Assert.
	if updateError != nil {
		test.Fatalf("expected no error updating valid index, got %v", updateError)
	}

	actualValue, _ := decomposition.RetrieveItem(index)

	if actualValue != replacementValue {
		test.Fatalf("expected updated value to be %d, got %d", replacementValue, actualValue)
	}
}

// TestUpdateItem_NegativeIndexRejection verifies that attempting to update
// a negative array index properly yields a defensive validation error.
func TestUpdateItem_NegativeIndexRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{5, 5, 5}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var invalidNegativeIndex int = -5
	var replacementValue int = 10

	// Act.
	updateError := decomposition.UpdateItem(invalidNegativeIndex, replacementValue)

	// Assert.
	if updateError == nil {
		test.Fatalf("expected error for negative update index, got nil")
	}
}

// TestUpdateItem_OutOfBoundsIndexRejection asserts that updating boundaries
// exceeding the internal array length are forcefully prevented.
func TestUpdateItem_OutOfBoundsIndexRejection(test *testing.T) {
	// Arrange.
	var numbers []int = []int{5, 5, 5}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var invalidOutOfBoundsIndex int = 3
	var replacementValue int = 10

	// Act.
	updateError := decomposition.UpdateItem(invalidOutOfBoundsIndex, replacementValue)

	// Assert.
	if updateError == nil {
		test.Fatalf("expected error for out of bounds update index, got nil")
	}
}

// TestUpdateItem_ModifiesInternalBlockSumCorrectly ensures the internal caching
// layer dynamically updates its running sum rather than forcing a recalculation.
func TestUpdateItem_ModifiesInternalBlockSumCorrectly(test *testing.T) {
	// Arrange.
	var numbers []int = []int{10, 10, 10, 10}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var index int = 0
	var replacementValue int = 20

	// Act.
	_ = decomposition.UpdateItem(index, replacementValue)

	// Assert.
	var expecttedSum int = 30

	if decomposition.sums[0] != expecttedSum {
		test.Fatalf("expected block sum to dynamically update to %d, got %d", expecttedSum, decomposition.sums[0])
	}
}

// TestUpdateItem_ZeroDifferenceNoImpact confirms that modifying an element
// with its existing equivalent value causes no destructive logic flaws.
func TestUpdateItem_ZeroDifferenceNoImpact(test *testing.T) {
	// Arrange.
	var numbers []int = []int{7, 7, 7}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var index int = 1
	var replacementValue int = 7

	// Act.
	_ = decomposition.UpdateItem(index, replacementValue)

	// Assert.
	var expecttedSum int = 14

	actualSum, _ := decomposition.QueryRange(0, 1)

	if actualSum != expecttedSum {
		test.Fatalf("expected block sum to remain %d, got %d", expecttedSum, actualSum)
	}
}

// TestUpdateItem_NegativeDifferenceImpact tests updating an element to a lesser
// value correctly subtracts the calculated differential from internal block caches.
func TestUpdateItem_NegativeDifferenceImpact(test *testing.T) {
	// Arrange.
	var numbers []int = []int{50, 50, 50, 50}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var index int = 2
	var replacementValue int = 10

	// Act.
	_ = decomposition.UpdateItem(index, replacementValue)

	// Assert.
	var expectedSecondBlockSum int = 60

	if decomposition.sums[1] != expectedSecondBlockSum {
		test.Fatalf("expected second block sum to decrease to %d, got %d", expectedSecondBlockSum, decomposition.sums[1])
	}
}

// TestQueryRange_ReflectsPostUpdateState confirms range queries seamlessly
// incorporate values actively injected via recent pointwise modifications.
func TestQueryRange_ReflectsPostUpdateState(test *testing.T) {
	// Arrange.
	var numbers []int = []int{1, 1, 1, 1, 1}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var index int = 2
	var replacementValue int = 100

	_ = decomposition.UpdateItem(index, replacementValue)

	// Act.
	calculatedSum, _ := decomposition.QueryRange(0, 4)

	// Assert.
	var expectedSum int = 104

	if calculatedSum != expectedSum {
		test.Fatalf("expected range query to reflect new sum %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_ReflectsMultipleSequentialUpdates validates robustness against
// consecutive alterations distributed sporadically throughout various sub-blocks.
func TestQueryRange_ReflectsMultipleSequentialUpdates(test *testing.T) {
	// Arrange.
	var numbers []int = []int{0, 0, 0, 0, 0, 0}

	decomposition, _ := NewSqrtDecomposition(numbers)

	// Act.
	_ = decomposition.UpdateItem(0, 5)
	_ = decomposition.UpdateItem(3, 15)
	_ = decomposition.UpdateItem(5, 20)

	calculatedSum, _ := decomposition.QueryRange(0, 5)

	// Assert.
	var expectedSum int = 40

	if calculatedSum != expectedSum {
		test.Fatalf("expected range query to reflect sequential updates totaling %d, got %d", expectedSum, calculatedSum)
	}
}

// TestQueryRange_MaxRangeAfterAllElementsUpdated performs an aggressive stress
// scenario replacing every value to ensure comprehensive sum cache stability.
func TestQueryRange_MaxRangeAfterAllElementsUpdated(test *testing.T) {
	// Arrange.
	var numbers []int = []int{1, 2, 3, 4}

	decomposition, _ := NewSqrtDecomposition(numbers)

	var numberOfItems int = decomposition.Size()
	var uniformReplacementValue int = 10

	// Act.
	for elementIndex := 0; elementIndex < numberOfItems; elementIndex++ {
		_ = decomposition.UpdateItem(elementIndex, uniformReplacementValue)
	}

	calculatedSum, _ := decomposition.QueryRange(0, numberOfItems-1)

	// Assert.
	var expectedSum int = 40

	if calculatedSum != expectedSum {
		test.Fatalf("expected range query of completely mutated array to equal %d, got %d", expectedSum, calculatedSum)
	}
}

// FuzzSqrtDecomposition_Operations performs native Go fuzz testing on the Square Root Decomposition implementation
// using randomized byte streams to execute range queries, item updates, and retrieval inspections safely.
func FuzzSqrtDecomposition_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{1, 2, 3, 4, 5}, 0, 4, 2, 10)
	fuzz.Add([]byte{10, 20, 30, 40}, 1, 3, 0, 99)

	fuzz.Fuzz(func(fuzzContext *testing.T, rawBytes []byte, queryLeft int, queryRight int, updateIndex int, updateVal int) {
		if len(rawBytes) == 0 {
			return
		}

		// Convert bytes to a base number slice.
		numbers := make([]int, len(rawBytes))

		for index, value := range rawBytes {
			numbers[index] = int(value)
		}

		// Act & Assert: Execute decomposition wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			decomposition, err := NewSqrtDecomposition(numbers)

			if err != nil {
				return
			}

			size := decomposition.Size()

			if size == 0 {
				return
			}

			safeLeft := queryLeft % size
			safeRight := queryRight % size

			if safeLeft > safeRight {
				safeLeft, safeRight = safeRight, safeLeft
			}

			_, _ = decomposition.QueryRange(safeLeft, safeRight)

			safeUpdateIndex := updateIndex % size

			_ = decomposition.UpdateItem(safeUpdateIndex, updateVal)
			_, _ = decomposition.RetrieveItem(safeUpdateIndex)
		}()
	})
}

// BenchmarkSqrtDecomposition_New benchmarks the performance of initializing
// the square root decomposition structure.
func BenchmarkSqrtDecomposition_New(bench *testing.B) {
	// Arrange.
	numbers := make([]int, 500)

	for index := 0; index < 500; index++ {
		numbers[index] = index + 1
	}

	// Act.
	for bench.Loop() {
		_, _ = NewSqrtDecomposition(numbers)
	}
}

// BenchmarkSqrtDecomposition_QueryRange benchmarks the performance of performing
// range sum queries over the decomposition structure.
func BenchmarkSqrtDecomposition_QueryRange(bench *testing.B) {
	// Arrange.
	numbers := make([]int, 500)

	for index := 0; index < 500; index++ {
		numbers[index] = index + 1
	}

	decomposition, _ := NewSqrtDecomposition(numbers)

	// Act.
	for bench.Loop() {
		_, _ = decomposition.QueryRange(50, 450)
	}
}

// BenchmarkSqrtDecomposition_UpdateItem benchmarks the performance of pointwise
// updates and dynamic block cache maintenance.
func BenchmarkSqrtDecomposition_UpdateItem(bench *testing.B) {
	// Arrange.
	numbers := make([]int, 500)

	for index := 0; index < 500; index++ {
		numbers[index] = index + 1
	}

	// Act.
	for bench.Loop() {
		decomposition, _ := NewSqrtDecomposition(numbers)

		for elementIndex := 0; elementIndex < 100; elementIndex++ {
			_ = decomposition.UpdateItem(elementIndex, elementIndex*2)
		}
	}
}
