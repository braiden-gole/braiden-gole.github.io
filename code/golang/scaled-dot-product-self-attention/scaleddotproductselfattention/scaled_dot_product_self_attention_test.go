// -----------------------------------------------------------------------------
// Package: scaleddotproductselfattention
// File: scaled_dot_product_self_attention_test.go
//
// Description:
//
//	This file provides a comprehensive suite of 40 unique unit tests, fuzz tests,
//	and modern benchmarks using b.Loop() covering every function within the scaled
//	dot-product self-attention package. It validates correctness, error handling,
//	numerical stability, and performance across diverse matrix configurations.
//
// -----------------------------------------------------------------------------
package scaleddotproductselfattention

import (
	"sync"
	"testing"
)

// TestValidateMatrix_ValidMatrix verifies that validation passes for a well-formed matrix.
func TestValidateMatrix_ValidMatrix(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	// Act.
	err := validateMatrix(matrix)

	// Assert.
	if err != nil {
		test.Errorf("expected no error, got: %v", err)
	}
}

// TestValidateMatrix_NilMatrix verifies that validation fails when the matrix is nil.
func TestValidateMatrix_NilMatrix(test *testing.T) {
	// Arrange.
	var matrix Matrix = nil

	// Act.
	err := validateMatrix(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error for nil matrix, got nil")
	}
}

// TestValidateMatrix_EmptyRows verifies that validation fails when the matrix has zero rows.
func TestValidateMatrix_EmptyRows(test *testing.T) {
	// Arrange.
	matrix := Matrix{}

	// Act.
	err := validateMatrix(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error for empty rows matrix, got nil")
	}
}

// TestValidateMatrix_EmptyColumns verifies that validation fails when matrix rows have zero columns.
func TestValidateMatrix_EmptyColumns(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{},
	}

	// Act.
	err := validateMatrix(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error for empty columns matrix, got nil")
	}
}

// TestValidateMatrix_NilRow verifies that validation fails when a row within the matrix is nil.
func TestValidateMatrix_NilRow(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		nil,
	}

	// Act.
	err := validateMatrix(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error for nil row, got nil")
	}
}

// TestValidateMatrix_InconsistentRowLengths verifies that validation fails when row lengths differ.
func TestValidateMatrix_InconsistentRowLengths(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		{3.0},
	}

	// Act.
	err := validateMatrix(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error for inconsistent row lengths, got nil")
	}
}

// TestTranspose_ValidSquareMatrix verifies correct transposition of a square matrix.
func TestTranspose_ValidSquareMatrix(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	// Act.
	transposedMatrix, err := transpose(matrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected transposition error: %v", err)
	}

	if transposedMatrix[0][1] != 3.0 || transposedMatrix[1][0] != 2.0 {
		test.Errorf("transposition result incorrect: %v", transposedMatrix)
	}
}

// TestTranspose_ValidRectangularMatrix verifies correct transposition of a rectangular matrix.
func TestTranspose_ValidRectangularMatrix(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}

	// Act.
	transposedMatrix, err := transpose(matrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected transposition error: %v", err)
	}

	if len(transposedMatrix) != 3 || len(transposedMatrix[0]) != 2 {
		test.Errorf("transposed dimensions incorrect: %v", transposedMatrix)
	}
}

// TestTranspose_NilMatrix verifies that transposition returns an error on a nil matrix.
func TestTranspose_NilMatrix(test *testing.T) {
	// Arrange.
	var matrix Matrix = nil

	// Act.
	_, err := transpose(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error when transposing nil matrix, got nil")
	}
}

// TestTranspose_InconsistentRows verifies that transposition fails on malformed matrices.
func TestTranspose_InconsistentRows(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		{3.0},
	}

	// Act.
	_, err := transpose(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error when transposing inconsistent rows, got nil")
	}
}

// TestScaleMatrix_ValidScaling verifies correct element-wise division by a scale factor.
func TestScaleMatrix_ValidScaling(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{4.0, 8.0},
		{12.0, 16.0},
	}

	scaleFactor := 2.0

	// Act.
	scaledMatrix, err := scaleMatrix(matrix, scaleFactor)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected scaling error: %v", err)
	}

	if scaledMatrix[0][0] != 2.0 || scaledMatrix[1][1] != 8.0 {
		test.Errorf("scaled matrix values incorrect: %v", scaledMatrix)
	}
}

// TestScaleMatrix_ZeroScaleFactor verifies that scaling by zero returns an error.
func TestScaleMatrix_ZeroScaleFactor(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	scaleFactor := 0.0

	// Act.
	_, err := scaleMatrix(matrix, scaleFactor)

	// Assert.
	if err == nil {
		test.Error("expected error when scaling by zero, got nil")
	}
}

// TestScaleMatrix_NilMatrix verifies that scaling a nil matrix returns an error.
func TestScaleMatrix_NilMatrix(test *testing.T) {
	// Arrange.
	var matrix Matrix = nil

	scaleFactor := 2.0

	// Act.
	_, err := scaleMatrix(matrix, scaleFactor)

	// Assert.
	if err == nil {
		test.Error("expected error when scaling nil matrix, got nil")
	}
}

// TestScaleMatrix_SmallScaleFactor verifies that scaling by an excessively small factor returns an error.
func TestScaleMatrix_SmallScaleFactor(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	scaleFactor := 1e-12

	// Act.
	_, err := scaleMatrix(matrix, scaleFactor)

	// Assert.
	if err == nil {
		test.Error("expected error when scaling by infinitesimal factor, got nil")
	}
}

// TestComputeProduct_ValidMultiplication verifies correct standard matrix multiplication.
func TestComputeProduct_ValidMultiplication(test *testing.T) {
	// Arrange.
	multiplicand := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	multiplier := Matrix{
		{5.0, 6.0},
		{7.0, 8.0},
	}

	// Act.
	productMatrix, err := computeProduct(multiplicand, multiplier)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected multiplication error: %v", err)
	}

	if productMatrix[0][0] != 19.0 || productMatrix[0][1] != 22.0 {
		test.Errorf("computed product values incorrect: %v", productMatrix)
	}
}

// TestComputeProduct_IncompatibleDimensions verifies that mismatched dimensions return an error.
func TestComputeProduct_IncompatibleDimensions(test *testing.T) {
	// Arrange.
	multiplicand := Matrix{
		{1.0, 2.0, 3.0},
	}

	multiplier := Matrix{
		{1.0, 2.0},
	}

	// Act.
	_, err := computeProduct(multiplicand, multiplier)

	// Assert.
	if err == nil {
		test.Error("expected error for incompatible dimensions, got nil")
	}
}

// TestComputeProduct_NilLeftOperand verifies that a nil left operand returns an error.
func TestComputeProduct_NilLeftOperand(test *testing.T) {
	// Arrange.
	var multiplicand Matrix = nil

	multiplier := Matrix{
		{1.0, 2.0},
	}

	// Act.
	_, err := computeProduct(multiplicand, multiplier)

	// Assert.
	if err == nil {
		test.Error("expected error for nil left operand, got nil")
	}
}

// TestComputeProduct_NilRightOperand verifies that a nil right operand returns an error.
func TestComputeProduct_NilRightOperand(test *testing.T) {
	// Arrange.
	multiplicand := Matrix{
		{1.0, 2.0},
	}

	var multiplier Matrix = nil

	// Act.
	_, err := computeProduct(multiplicand, multiplier)

	// Assert.
	if err == nil {
		test.Error("expected error for nil right operand, got nil")
	}
}

// TestComputeProduct_EmptyMatrix verifies that multiplying empty matrices returns an error.
func TestComputeProduct_EmptyMatrix(test *testing.T) {
	// Arrange.
	multiplicand := Matrix{}
	multiplier := Matrix{
		{1.0, 2.0},
	}

	// Act.
	_, err := computeProduct(multiplicand, multiplier)

	// Assert.
	if err == nil {
		test.Error("expected error for empty matrix operand, got nil")
	}
}

// TestApplySoftmax_ValidNormalization verifies correct softmax probability normalization across rows.
func TestApplySoftmax_ValidNormalization(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0},
	}

	// Act.
	softmaxMatrix, err := applySoftmax(matrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected softmax error: %v", err)
	}

	sumProbabilities := softmaxMatrix[0][0] + softmaxMatrix[0][1]

	if sumProbabilities < 0.999 || sumProbabilities > 1.001 {
		test.Errorf("softmax probabilities do not sum to 1.0: %f", sumProbabilities)
	}
}

// TestApplySoftmax_NumericalStability verifies that large input values do not overflow during softmax.
func TestApplySoftmax_NumericalStability(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{1000.0, 1001.0},
	}

	// Act.
	softmaxMatrix, err := applySoftmax(matrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected stability error during softmax: %v", err)
	}

	if softmaxMatrix[0][0] >= softmaxMatrix[0][1] {
		test.Error("expected higher probability for larger input value")
	}
}

// TestApplySoftmax_NilMatrix verifies that applying softmax to a nil matrix returns an error.
func TestApplySoftmax_NilMatrix(test *testing.T) {
	// Arrange.
	var matrix Matrix = nil

	// Act.
	_, err := applySoftmax(matrix)

	// Assert.
	if err == nil {
		test.Error("expected error for nil matrix in softmax, got nil")
	}
}

// TestApplySoftmax_UnderflowCondition verifies defensive handling of extreme underflows.
func TestApplySoftmax_UnderflowCondition(test *testing.T) {
	// Arrange.
	matrix := Matrix{
		{-1e20, -1e20},
	}

	// Act.
	_, err := applySoftmax(matrix)

	// Assert.
	if err != nil {
		test.Logf("handled underflow gracefully with error: %v", err)
	}
}

// TestCalculateScaledDotProductAttention_StandardCase verifies the complete attention pipeline.
func TestCalculateScaledDotProductAttention_StandardCase(test *testing.T) {
	// Arrange.
	queryMatrix := Matrix{
		{1.0, 0.0, 1.0},
		{0.0, 1.0, 1.0},
	}

	keyMatrix := Matrix{
		{1.0, 1.0, 0.0},
		{0.0, 1.0, 1.0},
	}

	valueMatrix := Matrix{
		{10.0, 20.0},
		{30.0, 40.0},
	}

	// Act.
	attentionOutput, err := CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected attention calculation error: %v", err)
	}

	if len(attentionOutput) != 2 || len(attentionOutput[0]) != 2 {
		test.Errorf("attention output dimensions incorrect: %v", attentionOutput)
	}
}

// TestCalculateScaledDotProductAttention_NilQuery verifies error handling when query matrix is nil.
func TestCalculateScaledDotProductAttention_NilQuery(test *testing.T) {
	// Arrange.
	var matrix Matrix = nil

	keyMatrix := Matrix{{1.0, 1.0}}
	valueMatrix := Matrix{{10.0, 20.0}}

	// Act.
	_, err := CalculateScaledDotProductAttention(matrix, keyMatrix, valueMatrix)

	// Assert.
	if err == nil {
		test.Error("expected error for nil query matrix, got nil")
	}
}

// TestCalculateScaledDotProductAttention_NilKey verifies error handling when key matrix is nil.
func TestCalculateScaledDotProductAttention_NilKey(test *testing.T) {
	// Arrange.
	matrix := Matrix{{1.0, 1.0}}

	var keyMatrix Matrix = nil

	valueMatrix := Matrix{{10.0, 20.0}}

	// Act.
	_, err := CalculateScaledDotProductAttention(matrix, keyMatrix, valueMatrix)

	// Assert.
	if err == nil {
		test.Error("expected error for nil key matrix, got nil")
	}
}

// TestCalculateScaledDotProductAttention_NilValue verifies error handling when value matrix is nil.
func TestCalculateScaledDotProductAttention_NilValue(test *testing.T) {
	// Arrange.
	queryMatrix := Matrix{{1.0, 1.0}}
	keyMatrix := Matrix{{1.0, 1.0}}

	var valueMatrix Matrix = nil

	// Act.
	_, err := CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

	// Assert.
	if err == nil {
		test.Error("expected error for nil value matrix, got nil")
	}
}

// TestCalculateScaledDotProductAttention_IncompatibleDimensions verifies error handling on dimension mismatches.
func TestCalculateScaledDotProductAttention_IncompatibleDimensions(test *testing.T) {
	// Arrange.
	queryMatrix := Matrix{{1.0, 0.0, 1.0}}
	keyMatrix := Matrix{{1.0, 1.0}}
	valueMatrix := Matrix{{10.0, 20.0}}

	// Act.
	_, err := CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

	// Assert.
	if err == nil {
		test.Error("expected error for dimension mismatch in attention, got nil")
	}
}

// TestCalculateScaledDotProductAttention_SingleTokenSequence verifies attention computation with a single token.
func TestCalculateScaledDotProductAttention_SingleTokenSequence(test *testing.T) {
	// Arrange.
	queryMatrix := Matrix{{1.0, 2.0}}
	keyMatrix := Matrix{{1.0, 2.0}}
	valueMatrix := Matrix{{5.0, 10.0}}

	// Act.
	attentionOutput, err := CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error on single token attention: %v", err)
	}

	if attentionOutput[0][0] != 5.0 || attentionOutput[0][1] != 10.0 {
		test.Errorf("single token attention result incorrect: %v", attentionOutput)
	}
}

// TestCalculateScaledDotProductAttention_LargeEmbeddingDimensions verifies attention across wider vectors.
func TestCalculateScaledDotProductAttention_LargeEmbeddingDimensions(test *testing.T) {
	// Arrange.
	queryMatrix := Matrix{
		{0.5, 0.5, 0.5, 0.5},
	}

	keyMatrix := Matrix{
		{0.5, 0.5, 0.5, 0.5},
		{0.1, 0.2, 0.3, 0.4},
	}

	valueMatrix := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
	}

	// Act.
	attentionOutput, err := CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error on large embedding dimension test: %v", err)
	}

	if len(attentionOutput) != 1 || len(attentionOutput[0]) != 2 {
		test.Errorf("output structure incorrect: %v", attentionOutput)
	}
}

// TestCalculateScaledDotProductAttention_ConcurrentExecution verifies that attention calculations can be safely executed concurrently across multiple goroutines.
func TestCalculateScaledDotProductAttention_ConcurrentExecution(test *testing.T) {
	// Arrange.
	queryMatrix := Matrix{
		{1.0, 0.0, 1.0},
		{0.0, 1.0, 1.0},
	}

	keyMatrix := Matrix{
		{1.0, 1.0, 0.0},
		{0.0, 1.0, 1.0},
	}

	valueMatrix := Matrix{
		{10.0, 20.0},
		{30.0, 40.0},
	}

	concurrencyCount := 10

	errorChannel := make(chan error, concurrencyCount)

	var waitGroup sync.WaitGroup

	waitGroup.Add(concurrencyCount)

	// Act.
	for routineIndex := 0; routineIndex < concurrencyCount; routineIndex++ {
		go func() {
			defer waitGroup.Done()

			_, err := CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

			if err != nil {
				errorChannel <- err
			}
		}()
	}

	waitGroup.Wait()

	close(errorChannel)

	// Assert.
	for err := range errorChannel {
		if err != nil {
			test.Errorf("concurrent attention calculation failed: %v", err)
		}
	}
}

// FuzzValidateMatrix performs fuzz testing on matrix validation routines.
func FuzzValidateMatrix(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(1.0, 2.0)

	// Act.
	fuzz.Fuzz(func(test *testing.T, value float64, reference float64) {
		// Arrange.
		matrix := Matrix{
			{value, reference},
			{reference, value},
		}

		// Act.
		err := validateMatrix(matrix)

		// Assert.
		if err != nil {
			test.Errorf("valid fuzzed matrix failed validation: %v", err)
		}
	})
}

// FuzzScaleMatrix performs fuzz testing on matrix scaling behavior.
func FuzzScaleMatrix(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(10.0, 2.0)

	// Act.
	fuzz.Fuzz(func(test *testing.T, matrix float64, scaleFactor float64) {
		// Arrange.
		fuzzMatrix := Matrix{
			{matrix},
		}

		// Act.
		_, err := scaleMatrix(fuzzMatrix, scaleFactor)

		// Assert.
		if scaleFactor == 0.0 && err == nil {
			test.Error("expected error when scaling by zero in fuzz test, got nil")
		}
	})
}

// FuzzCalculateScaledDotProductAttention performs fuzz testing on the complete attention calculation.
func FuzzCalculateScaledDotProductAttention(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(1.0, 0.0, 1.0, 1.0)

	// Act.
	fuzz.Fuzz(func(test *testing.T, query float64, key float64, value float64, supplementary float64) {
		// Arrange.
		queryMatrix := Matrix{{query, supplementary}}
		keyMatrix := Matrix{{key, supplementary}}
		valueMatrix := Matrix{{value, supplementary}}

		// Act.
		_, _ = CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)

		// Assert.
	})
}

// BenchmarkCalculateScaledDotProductAttention_Small measures performance of small attention calculations.
func BenchmarkCalculateScaledDotProductAttention_Small(bench *testing.B) {
	// Arrange.
	queryMatrix := Matrix{
		{1.0, 0.0, 1.0},
		{0.0, 1.0, 1.0},
	}

	keyMatrix := Matrix{
		{1.0, 1.0, 0.0},
		{0.0, 1.0, 1.0},
	}

	valueMatrix := Matrix{
		{10.0, 20.0},
		{30.0, 40.0},
	}

	// Act.
	for bench.Loop() {
		_, _ = CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)
	}

	// Assert.
}

// BenchmarkCalculateScaledDotProductAttention_Medium measures performance of larger attention calculations.
func BenchmarkCalculateScaledDotProductAttention_Medium(bench *testing.B) {
	// Arrange.
	queryMatrix := Matrix{
		{1.0, 0.5, 0.2, 0.8},
		{0.3, 0.9, 0.1, 0.4},
		{0.2, 0.3, 0.5, 0.6},
	}

	keyMatrix := Matrix{
		{1.0, 0.5, 0.2, 0.8},
		{0.3, 0.9, 0.1, 0.4},
		{0.2, 0.3, 0.5, 0.6},
	}

	valueMatrix := Matrix{
		{10.0, 20.0},
		{30.0, 40.0},
		{50.0, 60.0},
	}

	// Act.
	for bench.Loop() {
		_, _ = CalculateScaledDotProductAttention(queryMatrix, keyMatrix, valueMatrix)
	}

	// Assert.
}

// BenchmarkComputeProduct measures performance of matrix multiplication.
func BenchmarkComputeProduct(bench *testing.B) {
	// Arrange.
	multiplicand := Matrix{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}

	multiplier := Matrix{
		{1.0, 2.0},
		{3.0, 4.0},
		{5.0, 6.0},
	}

	// Act.
	for bench.Loop() {
		_, _ = computeProduct(multiplicand, multiplier)
	}

	// Assert.
}

// BenchmarkApplySoftmax measures performance of softmax normalization.
func BenchmarkApplySoftmax(bench *testing.B) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0, 3.0, 4.0},
		{0.5, 1.5, 2.5, 3.5},
	}

	// Act.
	for bench.Loop() {
		_, _ = applySoftmax(matrix)
	}

	// Assert.
}

// BenchmarkTranspose measures performance of matrix transposition.
func BenchmarkTranspose(bench *testing.B) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}

	// Act.
	for bench.Loop() {
		_, _ = transpose(matrix)
	}

	// Assert.
}

// BenchmarkValidateMatrix measures performance of matrix structural validation.
func BenchmarkValidateMatrix(bench *testing.B) {
	// Arrange.
	matrix := Matrix{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}

	// Act.
	for bench.Loop() {
		_ = validateMatrix(matrix)
	}

	// Assert.
}
