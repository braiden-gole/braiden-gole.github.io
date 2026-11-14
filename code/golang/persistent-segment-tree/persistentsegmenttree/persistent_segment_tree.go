// -----------------------------------------------------------------------------
// Package: persistentsegmenttree
// File: persistent_segment_tree.go
//
// Description:
//
//	This file implements a Persistent Segment Tree using path copying.
//	It manages multiple historical versions of an integer array efficiently
//	via a slice of root Node pointers, supporting immutable range queries
//	and versioned point updates.
//
// Intended Coverage:
//   - Node and SegmentTree structural definitions.
//   - Recursive segment tree construction with bounds checking.
//   - Path-copied recursive updates generating new historical versions.
//   - Logarithmic range sum queries across any version.
//   - Public API methods: NewSegmentTree, InsertVersion, and QueryRange.
//
// -----------------------------------------------------------------------------
package persistentsegmenttree

import "errors"

// Node represents a single vertex in the persistent segment tree.
// It stores an accumulated value and pointers to its left and right child nodes.
type Node struct {
	value int
	left  *Node
	right *Node
}

// SegmentTree manages the historical versions of the persistent segment tree
// using a slice of root pointers and records the total length of the base array.
type SegmentTree struct {
	root   []*Node
	length int
}

// buildSegmentTree recursively constructs a balanced segment tree from an initial slice of values.
// It returns the newly allocated root node or an error if the range boundaries are invalid.
func buildSegmentTree(values []int, left int, right int) (*Node, error) {
	// Validate that the current construction boundaries are within valid array indices.
	if left < 0 || right >= len(values) || left > right {
		return nil, errors.New("invalid range boundaries during tree construction")
	}

	// Allocate a new tree node for the current segment range.
	currentNode := &Node{}

	// If this is a leaf node, assign the specific value from the initial values slice.
	if left == right {
		currentNode.value = values[left]

		return currentNode, nil
	}

	// Compute the midpoint index to divide the current segment in half.
	middleIndex := left + (right-left)/2

	// Recursively construct the left subtree for the lower half of the range.
	leftSubtreeNode, buildError := buildSegmentTree(values, left, middleIndex)

	if buildError != nil {
		return nil, buildError
	}

	// Recursively construct the right subtree for the upper half of the range.
	rightSubtreeNode, buildError := buildSegmentTree(values, middleIndex+1, right)

	if buildError != nil {
		return nil, buildError
	}

	// Attach the child subtrees and aggregate their values into the current node.
	currentNode.left = leftSubtreeNode
	currentNode.right = rightSubtreeNode
	currentNode.value = leftSubtreeNode.value + rightSubtreeNode.value

	return currentNode, nil
}

// recursiveUpdate performs a path-copied update on the segment tree, creating new nodes
// along the update path while sharing unmodified subtrees with the previous version.
func recursiveUpdate(current *Node, left int, right int, targetIndex int, newValue int) (*Node, error) {
	// Verify that the target index falls within the active segment bounds.
	if targetIndex < left || targetIndex > right {
		return nil, errors.New("target index falls outside current recursion range")
	}

	// Create a new token node representing this path-copied level in the tree.
	newToken := &Node{}

	// When reaching the leaf node corresponding to the target index, set the new value.
	if left == right {
		newToken.value = newValue

		return newToken, nil
	}

	// Calculate the middle index to determine which child subtree contains the target index.
	middleIndex := left + (right-left)/2

	// Branch left if the target index is less than or equal to the midpoint.
	if targetIndex <= middleIndex {
		var updatedLeft *Node
		var updateError error

		// If the previous version's left child exists, recurse down it.
		if current != nil && current.left != nil {
			updatedLeft, updateError = recursiveUpdate(current.left, left, middleIndex, targetIndex, newValue)
		} else {
			// If the left subtree is missing, dynamically build a zero-initialized subtree.
			segmentLength := (middleIndex - left) + 1
			zeroValues := make([]int, segmentLength)

			builtSubtree, buildError := buildSegmentTree(zeroValues, 0, segmentLength-1)

			if buildError != nil {
				return nil, buildError
			}

			// Perform the recursive update on the newly generated zero-initialized subtree.
			updatedLeft, updateError = recursiveUpdate(builtSubtree, 0, segmentLength-1, targetIndex-left, newValue)
		}

		if updateError != nil {
			return nil, updateError
		}

		// Assign the updated left child and preserve a pointer to the unchanged right child.
		newToken.left = updatedLeft

		if current != nil {
			newToken.right = current.right
		}
	} else {
		// Branch right if the target index is greater than the midpoint.
		var updatedRight *Node
		var updateError error

		// If the previous version's right child exists, recurse down it.
		if current != nil && current.right != nil {
			updatedRight, updateError = recursiveUpdate(current.right, middleIndex+1, right, targetIndex, newValue)
		} else {
			// If the right subtree is missing, dynamically build a zero-initialized subtree.
			segmentLength := right - middleIndex
			zeroValues := make([]int, segmentLength)

			builtSubtree, buildError := buildSegmentTree(zeroValues, 0, segmentLength-1)

			if buildError != nil {
				return nil, buildError
			}

			// Perform the recursive update on the newly generated zero-initialized subtree.
			updatedRight, updateError = recursiveUpdate(builtSubtree, 0, segmentLength-1, targetIndex-(middleIndex+1), newValue)
		}

		if updateError != nil {
			return nil, updateError
		}

		// Assign the updated right child and preserve a pointer to the unchanged left child.
		newToken.right = updatedRight

		if current != nil {
			newToken.left = current.left
		}
	}

	// Retrieve values from children safely to compute the updated aggregated node sum.
	leftValue := 0

	if newToken.left != nil {
		leftValue = newToken.left.value
	}

	rightValue := 0

	if newToken.right != nil {
		rightValue = newToken.right.value
	}

	// Update the current node's value with the sum of its children.
	newToken.value = leftValue + rightValue

	return newToken, nil
}

// recursiveQuery traverses the segment tree to compute the aggregate sum
// over a specified target range [leftIndex, rightIndex].
func recursiveQuery(current *Node, left int, right int, leftIndex int, rightIndex int) (int, error) {
	// If the node is nil, it contributes zero to the sum.
	if current == nil {
		return 0, nil
	}

	// If the current range is completely outside the query range, return zero.
	if leftIndex > right || rightIndex < left {
		return 0, nil
	}

	// If the current range is fully contained within the query range, return its accumulated value.
	if leftIndex <= left && right <= rightIndex {
		return current.value, nil
	}

	// Calculate the middle index to split the query across child subtrees.
	middleIndex := left + (right-left)/2

	// Recursively query the left subtree for overlapping segments.
	leftSum, queryError := recursiveQuery(current.left, left, middleIndex, leftIndex, rightIndex)

	if queryError != nil {
		return 0, queryError
	}

	// Recursively query the right subtree for overlapping segments.
	rightSum, queryError := recursiveQuery(current.right, middleIndex+1, right, leftIndex, rightIndex)

	if queryError != nil {
		return 0, queryError
	}

	// Combine and return the sums from both subtrees.
	return leftSum + rightSum, nil
}

// NewSegmentTree initializes a new persistent segment tree from a slice of values
// and returns the constructed segment tree container holding Version 0.
func NewSegmentTree(values []int) (*SegmentTree, error) {
	// Ensure the input values slice is not empty before building.
	if len(values) == 0 {
		return nil, errors.New("values slice cannot be empty")
	}

	length := len(values)

	// Build the base tree representing Version 0.
	initialRoot, buildError := buildSegmentTree(values, 0, length-1)

	if buildError != nil {
		return nil, buildError
	}

	// Package the root node into a SegmentTree structure containing Version 0.
	return &SegmentTree{
		root:   []*Node{initialRoot},
		length: length,
	}, nil
}

// InsertVersion creates a new historical version of the segment tree by updating
// a single value at the specified target index, returning the new version index.
func (tree *SegmentTree) InsertVersion(versionIndex int, targetIndex int, newValue int) (int, error) {
	// Verify that the requested historical version index exists.
	if versionIndex < 0 || versionIndex >= len(tree.root) {
		return -1, errors.New("provided version index is out of bounds")
	}

	// Verify that the target update index is within array bounds.
	if targetIndex < 0 || targetIndex >= tree.length {
		return -1, errors.New("target index for update is out of bounds")
	}

	// Fetch the root node of the selected base version.
	latestRoot := tree.root[versionIndex]

	// Generate a new path-copied root node via recursive update.
	newRoot, updateError := recursiveUpdate(latestRoot, 0, tree.length-1, targetIndex, newValue)

	if updateError != nil {
		return -1, updateError
	}

	// Append the new root node to the version history slice.
	tree.root = append(tree.root, newRoot)
	newVersionIndex := len(tree.root) - 1

	return newVersionIndex, nil
}

// QueryRange computes the aggregate sum over a closed range [left, right]
// against a specific historical tree version.
func (tree *SegmentTree) QueryRange(versionIndex int, left int, right int) (int, error) {
	// Verify that the requested historical version index exists.
	if versionIndex < 0 || versionIndex >= len(tree.root) {
		return 0, errors.New("provided version index is out of bounds")
	}

	// Verify that the query range boundaries are valid.
	if left < 0 || right >= tree.length || left > right {
		return 0, errors.New("query index boundaries are invalid")
	}

	// Locate the root node for the requested version.
	node := tree.root[versionIndex]

	// Execute the recursive query over the targeted tree version.
	return recursiveQuery(node, 0, tree.length-1, left, right)
}
