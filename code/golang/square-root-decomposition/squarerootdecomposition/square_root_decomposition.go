// -----------------------------------------------------------------------------
// Package: squarerootdecomposition
// File: square_root_decomposition.go
//
// Description:
//
//	This file implements the Square Root Decomposition data structure in Go,
//	focusing on structured initialization, defensive validation bounds checking,
//	range sum queries, and dynamic pointwise updates.
//
//	The implementation prioritizes robust error handling against malformed inputs,
//	clean encapsulation by deep-copying slice payloads, and clear separation
//	between array lengths and block partitioning sizes.
//
// -----------------------------------------------------------------------------
package squarerootdecomposition

import (
	"errors"
	"fmt"
	"math"
)

// SquareRootDecomposition encapsulates the state and operations for square root decomposition,
// maintaining an underlying data array, precomputed block sums, array size capacity, and block sizing parameters.
type SquareRootDecomposition struct {
	array  []int
	length int
	sums   []int
	size   int
}

// validateIndex checks whether a target index is within the valid bounds of the array capacity,
// returning a descriptive error if the index is negative or exceeds the total length.
func (decomposition *SquareRootDecomposition) validateIndex(index int) error {
	// Index boundaries must fall within zero and length minus one.
	if index < 0 || index >= decomposition.length {
		return fmt.Errorf("Index %d is out of bounds for array size %d", index, decomposition.length)
	}

	return nil
}

// validateRange checks whether query boundaries are logically sound and within capacity limits,
// preventing out-of-bounds access or reversed start and end indices.
func (decomposition *SquareRootDecomposition) validateRange(left, right int) error {
	// Left boundary must not be negative, right boundary must not exceed length, and left must be less than or equal to right.
	if left < 0 || right >= decomposition.length || left > right {
		return fmt.Errorf("invalid query range: leftIndex %d, rightIndex %d for array size %d", left, right, decomposition.length)
	}

	return nil
}

// NewSqrtDecomposition initializes and returns a new SquareRootDecomposition instance.
// It validates input availability, performs defensive deep copying of the source data,
// calculates optimal block sizing, and precomputes initial block sum caches.
func NewSqrtDecomposition(array []int) (*SquareRootDecomposition, error) {
	// Input arrays must contain elements to initialize decomposition.
	if len(array) == 0 {
		return nil, errors.New("array size cannot be empty or nil")
	}

	// Create a defensive deep copy of the slice to protect against external source mutations.
	copied := make([]int, len(array))

	copy(copied, array)

	size := len(copied)

	// Determine optimal block size using the square root of total elements.
	blockSize := int(math.Sqrt(float64(size)))

	// Guarantee block size is at least one to prevent division-by-zero errors.
	if blockSize == 0 {
		blockSize = 1
	}

	// Calculate total number of blocks required and allocate the sums cache.
	numberOfBlocks := (size + blockSize - 1) / blockSize
	sums := make([]int, numberOfBlocks)

	// Precompute block sum aggregates for fast lookup during range operations.
	for index := 0; index < size; index++ {
		blockIndex := index / blockSize
		sums[blockIndex] += copied[index]
	}

	return &SquareRootDecomposition{
		array:  copied,
		length: size,
		sums:   sums,
		size:   blockSize,
	}, nil
}

// Size provides read access to the total number of elements managed by the structure.
func (decomposition *SquareRootDecomposition) Size() int {
	return decomposition.length
}

// RetrieveItem safely looks up and returns the value at a specified index with defensive bounds checking.
func (decomposition *SquareRootDecomposition) RetrieveItem(index int) (int, error) {
	// Verify index validity before accessing internal array elements.
	if err := decomposition.validateIndex(index); err != nil {
		return 0, err
	}

	return decomposition.array[index], nil
}

// QueryRange computes the sum of elements between the specified left and right indices inclusive.
// It optimizes performance by scanning partial boundary blocks element-by-element while adding fully
// contained blocks in O(1) time via precomputed block sums.
func (decomposition *SquareRootDecomposition) QueryRange(left, right int) (int, error) {
	// Ensure query range bounds are valid.
	if err := decomposition.validateRange(left, right); err != nil {
		return 0, err
	}

	sum := 0
	beginning := left / decomposition.size
	end := right / decomposition.size

	// If the entire query falls within a single block, perform a straightforward linear scan.
	if beginning == end {
		for index := left; index <= right; index++ {
			sum += decomposition.array[index]
		}

		return sum, nil
	}

	// Sum partial elements in the starting block.
	next := (beginning + 1) * decomposition.size

	for index := left; index < next; index++ {
		sum += decomposition.array[index]
	}

	// Accumulate fully contained middle blocks in O(1) time per block.
	for blockIndex := beginning + 1; blockIndex < end; blockIndex++ {
		sum += decomposition.sums[blockIndex]
	}

	// Sum partial elements in the ending block.
	endIndex := end * decomposition.size

	for index := endIndex; index <= right && index < decomposition.length; index++ {
		sum += decomposition.array[index]
	}

	return sum, nil
}

// UpdateItem modifies a single element at the specified index and dynamically updates
// the corresponding block sum cache to maintain state consistency efficiently.
func (decomposition *SquareRootDecomposition) UpdateItem(index, value int) error {
	// Verify index validity before applying updates.
	if err := decomposition.validateIndex(index); err != nil {
		return err
	}

	blockIndex := index / decomposition.size

	oldValue := decomposition.array[index]

	// Update the element value and adjust the block sum by the differential change.
	decomposition.array[index] = value
	decomposition.sums[blockIndex] += (value - oldValue)

	return nil
}
