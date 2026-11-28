// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file serves as an executable demonstration driver for the scaled dot-product
//	self-attention package. It initializes sample query, key, and value matrices,
//	invokes the core self-attention calculation function, handles any returned errors,
//	and prints out the resulting weighted context vectors row by row.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/braiden-gole/scaled-dot-product-self-attention/scaleddotproductselfattention"
)

func main() {
	// Initialize the query matrix representing token representations looking for information.
	queryMatrix := scaleddotproductselfattention.Matrix{
		{1.0, 0.0, 1.0},
		{0.0, 1.0, 1.0},
	}

	// Initialize the key matrix representing what information each token holds for matching.
	keyMatrix := scaleddotproductselfattention.Matrix{
		{1.0, 1.0, 0.0},
		{0.0, 1.0, 1.0},
	}

	// Initialize the value matrix containing the actual content vectors to be aggregated.
	valueMatrix := scaleddotproductselfattention.Matrix{
		{10.0, 20.0},
		{30.0, 40.0},
	}

	// Compute the scaled dot-product attention output using the initialized matrices.
	attentionOutput, calculationError := scaleddotproductselfattention.CalculateScaledDotProductAttention(
		queryMatrix,
		keyMatrix,
		valueMatrix,
	)

	// Terminate execution and display an error message if the attention calculation fails.
	if calculationError != nil {
		fmt.Printf("Error: %v\n", calculationError)

		return
	}

	// Print a header label for the calculation results.
	fmt.Println("Result:")

	// Iterate through each row of the final attention output and print its values.
	for rowIndex, rowValues := range attentionOutput {
		fmt.Printf("Row %d: %v\n", rowIndex, rowValues)
	}
}
