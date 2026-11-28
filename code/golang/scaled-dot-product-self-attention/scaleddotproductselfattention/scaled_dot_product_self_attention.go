// -----------------------------------------------------------------------------
// Package: scaleddotproductselfattention
// File: scaled_dot_product_self_attention.go
//
// Description:
//
//	This file implements the scaled dot-product self-attention mechanism
//	used by transformer neural network architectures in pure Go. It provides
//	robust matrix operations including dimension validation, multiplication,
//	transposition, scaling, and numerical stability-protected softmax normalization.
//
//	The implementation prioritizes defensive programming against malformed or
//	incompatible matrix dimensions, floating-point safety, and clean separation
//	of concerns across each step of the self-attention formula:
//	Attention(Q, K, V) = softmax((Q * K^T) / sqrt(d_k)) * V.
//
// -----------------------------------------------------------------------------
package scaleddotproductselfattention

import (
	"errors"
	"math"
)

// Matrix represents a two-dimensional slice of floating-point numbers used for
// all matrix arithmetic and tensor operations within the self-attention pipeline.
type Matrix [][]float64

// validateMatrix ensures that the provided matrix is non-nil, non-empty,
// and possesses consistent row lengths to prevent runtime out-of-bounds panics
// during matrix calculations.
func validateMatrix(inputMatrix Matrix) error {
	// Guard against nil matrix references.
	if inputMatrix == nil {
		return errors.New("input matrix cannot be nil.")
	}

	rowCount := len(inputMatrix)

	// Ensure the matrix contains at least one row.
	if rowCount == 0 {
		return errors.New("input matrix cannot be empty.")
	}

	columnCount := len(inputMatrix[0])

	// Ensure the matrix rows contain at least one column.
	if columnCount == 0 {
		return errors.New("input matrix rows cannot be empty.")
	}

	// Validate every row for nil pointers and uniform length.
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		if inputMatrix[rowIndex] == nil {
			return errors.New("input matrix row cannot be nil.")
		}

		if len(inputMatrix[rowIndex]) != columnCount {
			return errors.New("input matrix contains inconsistent row lengths.")
		}
	}

	return nil
}

// computeProduct performs standard matrix multiplication between a left operand matrix
// and a right operand matrix with rigorous dimensional validation and compatibility checks.
func computeProduct(multiplicandMatrix Matrix, multiplierMatrix Matrix) (Matrix, error) {
	// Validate the left-hand operand matrix structure.
	if validationError := validateMatrix(multiplicandMatrix); validationError != nil {
		return nil, validationError
	}

	// Validate the right-hand operand matrix structure.
	if validationError := validateMatrix(multiplierMatrix); validationError != nil {
		return nil, validationError
	}

	rowCountLeft := len(multiplicandMatrix)
	columnCountLeft := len(multiplicandMatrix[0])

	rowCountRight := len(multiplierMatrix)
	columnCountRight := len(multiplierMatrix[0])

	// Ensure inner dimensions match for valid matrix multiplication.
	if columnCountLeft != rowCountRight {
		return nil, errors.New("matrix dimensions are incompatible for multiplication.")
	}

	resultMatrix := make(Matrix, rowCountLeft)

	// Compute dot products across rows and columns.
	for rowIndex := 0; rowIndex < rowCountLeft; rowIndex++ {
		resultMatrix[rowIndex] = make([]float64, columnCountRight)

		for columnIndex := 0; columnIndex < columnCountRight; columnIndex++ {
			var dotProductSum float64 = 0.0

			for innerIndex := 0; innerIndex < columnCountLeft; innerIndex++ {
				dotProductSum += multiplicandMatrix[rowIndex][innerIndex] * multiplierMatrix[innerIndex][columnIndex]
			}

			resultMatrix[rowIndex][columnIndex] = dotProductSum
		}
	}

	return resultMatrix, nil
}

// transpose safely swaps the rows and columns of a given two-dimensional matrix,
// aligning dimensions so that the Key matrix can be multiplied with the Query matrix.
func transpose(inputMatrix Matrix) (Matrix, error) {
	// Validate the input matrix structure prior to transposition.
	if validationError := validateMatrix(inputMatrix); validationError != nil {
		return nil, validationError
	}

	rowCount := len(inputMatrix)
	columnCount := len(inputMatrix[0])

	transposedMatrix := make(Matrix, columnCount)

	// Swap indices to construct the transposed matrix layout.
	for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
		transposedMatrix[columnIndex] = make([]float64, rowCount)

		for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
			transposedMatrix[columnIndex][rowIndex] = inputMatrix[rowIndex][columnIndex]
		}
	}

	return transposedMatrix, nil
}

// scaleMatrix divides every element in the matrix by a scale factor derived from
// the square root of the key dimension, preventing gradient vanishing issues during training.
func scaleMatrix(matrix Matrix, scaleFactor float64) (Matrix, error) {
	// Validate the matrix structure before applying scaling.
	if validationError := validateMatrix(matrix); validationError != nil {
		return nil, validationError
	}

	// Guard against division by zero or extremely small scale factors.
	if math.Abs(scaleFactor) < 1e-9 {
		return nil, errors.New("scale factor cannot be zero or excessively close to zero.")
	}

	rowCount := len(matrix)
	columnCount := len(matrix[0])

	scaledMatrix := make(Matrix, rowCount)

	// Scale each element by the provided factor.
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		scaledMatrix[rowIndex] = make([]float64, columnCount)

		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			scaledMatrix[rowIndex][columnIndex] = matrix[rowIndex][columnIndex] / scaleFactor
		}
	}

	return scaledMatrix, nil
}

// applySoftmax normalizes the rows of a matrix so they form a valid probability distribution
// summing to 1.0. It incorporates max-subtraction for numerical stability and guards
// against exponential sum underflows.
func applySoftmax(inputMatrix Matrix) (Matrix, error) {
	// Validate the input matrix before applying softmax normalization.
	if validationError := validateMatrix(inputMatrix); validationError != nil {
		return nil, validationError
	}

	rowCount := len(inputMatrix)
	columnCount := len(inputMatrix[0])

	softmaxMatrix := make(Matrix, rowCount)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		softmaxMatrix[rowIndex] = make([]float64, columnCount)

		// Find the maximum value in the current row for numerical stability.
		maximumValue := inputMatrix[rowIndex][0]

		for columnIndex := 1; columnIndex < columnCount; columnIndex++ {
			if inputMatrix[rowIndex][columnIndex] > maximumValue {
				maximumValue = inputMatrix[rowIndex][columnIndex]
			}
		}

		// Calculate exponential values relative to the row maximum and accumulate their sum.
		var exponentialSum float64 = 0.0

		rowExponentials := make([]float64, columnCount)

		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			rowExponentials[columnIndex] = math.Exp(inputMatrix[rowIndex][columnIndex] - maximumValue)
			exponentialSum += rowExponentials[columnIndex]
		}

		// Guard against exponential sum underflow to zero.
		if math.Abs(exponentialSum) < 1e-9 {
			return nil, errors.New("exponential sum underflowed to zero during softmax calculation.")
		}

		// Normalize each element in the row by the accumulated exponential sum.
		for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
			softmaxMatrix[rowIndex][columnIndex] = rowExponentials[columnIndex] / exponentialSum
		}
	}
	return softmaxMatrix, nil
}

// CalculateScaledDotProductAttention computes the core self-attention formula:
// Attention(Query, Key, Value) = softmax(Query * Key^T / sqrt(d_k)) * Value.
// It validates all input matrices, transposes the keys, computes raw scores, scales them,
// applies softmax to get attention weights, and multiplies by the value matrix.
func CalculateScaledDotProductAttention(queryMatrix Matrix, keyMatrix Matrix, valueMatrix Matrix) (Matrix, error) {
	// Validate all incoming matrices for structural integrity.
	if validationError := validateMatrix(queryMatrix); validationError != nil {
		return nil, validationError
	}

	if validationError := validateMatrix(keyMatrix); validationError != nil {
		return nil, validationError
	}

	if validationError := validateMatrix(valueMatrix); validationError != nil {
		return nil, validationError
	}

	// Step 1: Transpose the Key matrix so dimensions align for multiplication with the Query matrix.
	transposedKeyMatrix, transposeError := transpose(keyMatrix)

	if transposeError != nil {
		return nil, transposeError
	}

	// Step 2: Multiply Query matrix by Transposed Key matrix to get raw attention scores.
	rawAttentionScores, multiplicationError := computeProduct(queryMatrix, transposedKeyMatrix)

	if multiplicationError != nil {
		return nil, multiplicationError
	}

	// Step 3: Determine the scaling factor based on the dimensionality of the key vectors.
	keyDimension := float64(len(keyMatrix[0]))
	scalingFactor := math.Sqrt(keyDimension)

	// Step 4: Scale the raw attention scores to prevent gradients from vanishing during training.
	scaledAttentionScores, scaleError := scaleMatrix(rawAttentionScores, scalingFactor)

	if scaleError != nil {
		return nil, scaleError
	}

	// Step 5: Apply Softmax function to turn scores into attention weights.
	attentionWeights, softmaxError := applySoftmax(scaledAttentionScores)

	if softmaxError != nil {
		return nil, softmaxError
	}

	// Step 6: Multiply attention weights by the Value matrix to get the final context vectors.
	finalAttentionOutput, finalMultiplicationError := computeProduct(attentionWeights, valueMatrix)

	if finalMultiplicationError != nil {
		return nil, finalMultiplicationError
	}

	return finalAttentionOutput, nil
}
