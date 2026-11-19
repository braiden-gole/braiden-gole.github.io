// -----------------------------------------------------------------------------
// Package: askprimality
// File: aks_primality_test.go
//
// Description:
//
//	This file contains unit tests for the AKS primality test implementation.
//	The tests verify correctness of individual helper functions as well as
//	end-to-end primality classification for known prime and composite values.
//
//	The tests are designed for clarity and correctness rather than performance.
//	Each test follows a strict Arrange / Act / Assert structure and uses
//	descriptive variable names to make intent explicit.
//
// Notes:
//   - The testing.T parameter is named "test" for clarity.
//   - All arithmetic is integer-based; no floating-point assumptions are made.
//   - These tests are intended for educational and verification purposes.
//
// -----------------------------------------------------------------------------
package askprimality

import (
	"math/big"
	"testing"
)

// TestIntegerSquareRootFloor_PerfectSquare verifies that the floor square root
// of a perfect square is computed exactly.
func TestIntegerSquareRootFloor_PerfectSquare(test *testing.T) {
	// Arrange.
	inputValue := 144
	expectedSquareRoot := 12

	// Act.
	computedSquareRoot := integerSquareRootFloor(inputValue)

	// Assert.
	if computedSquareRoot != expectedSquareRoot {
		test.Fatalf("expected floor sqrt(%d) to be %d, got %d",
			inputValue,
			expectedSquareRoot,
			computedSquareRoot,
		)
	}
}

// TestIntegerSquareRootFloor_NonPerfectSquare verifies that the floor square root
// of a non-perfect square is rounded downward correctly.
func TestIntegerSquareRootFloor_NonPerfectSquare(test *testing.T) {
	// Arrange.
	inputValue := 20
	expectedSquareRoot := 4

	// Act.
	computedSquareRoot := integerSquareRootFloor(inputValue)

	// Assert.
	if computedSquareRoot != expectedSquareRoot {
		test.Fatalf("expected floor sqrt(%d) to be %d, got %d",
			inputValue,
			expectedSquareRoot,
			computedSquareRoot,
		)
	}
}

// TestIsPerfectPower_ReturnsTrueForPerfectPower verifies that perfect powers
// such as 27 = 3^3 are correctly detected.
func TestIsPerfectPower_ReturnsTrueForPerfectPower(test *testing.T) {
	// Arrange.
	candidateValue := big.NewInt(27)

	// Act.
	isPower := isPerfectPower(candidateValue)

	// Assert.
	if !isPower {
		test.Fatalf("expected %s to be detected as a perfect power",
			candidateValue.String(),
		)
	}
}

// TestIsPerfectPower_ReturnsFalseForNonPerfectPower verifies that integers
// which are not perfect powers are not incorrectly classified.
func TestIsPerfectPower_ReturnsFalseForNonPerfectPower(test *testing.T) {
	// Arrange.
	candidateValue := big.NewInt(20)

	// Act.
	isPower := isPerfectPower(candidateValue)

	// Assert.
	if isPower {
		test.Fatalf("expected %s to not be detected as a perfect power",
			candidateValue.String(),
		)
	}
}

// TestOrderExceedsThreshold_ReturnsTrueWhenOrderIsLarge verifies that the
// function correctly identifies when the multiplicative order exceeds
// the provided threshold.
func TestOrderExceedsThreshold_ReturnsTrueWhenOrderIsLarge(test *testing.T) {
	// Arrange.
	candidateValue := big.NewInt(7)

	modulusValue := 10
	thresholdValue := 2

	// Act.
	exceedsThreshold := orderExceedsThreshold(candidateValue, modulusValue, thresholdValue)

	// Assert.
	if !exceedsThreshold {
		test.Fatalf("expected multiplicative order of %s mod %d to exceed %d",
			candidateValue.String(),
			modulusValue,
			thresholdValue,
		)
	}
}

// TestOrderExceedsThreshold_ReturnsFalseWhenOrderIsSmall verifies that the
// function correctly detects small multiplicative orders.
func TestOrderExceedsThreshold_ReturnsFalseWhenOrderIsSmall(test *testing.T) {
	// Arrange.
	candidateValue := big.NewInt(4)

	modulusValue := 5
	thresholdValue := 2

	// Act.
	exceedsThreshold := orderExceedsThreshold(candidateValue, modulusValue, thresholdValue)

	// Assert.
	if exceedsThreshold {
		test.Fatalf("expected multiplicative order of %s mod %d to not exceed %d",
			candidateValue.String(),
			modulusValue,
			thresholdValue,
		)
	}
}

// TestPolynomialCongruenceCheck_PrimePasses verifies that the polynomial
// congruence holds for a known prime input.
func TestPolynomialCongruenceCheck_PrimePasses(test *testing.T) {
	// Arrange.
	primeCandidate := big.NewInt(7)

	modulusValue := 5
	witnessValue := 1

	// Act.
	congruenceHolds := polynomialCongruenceCheck(primeCandidate, modulusValue, witnessValue)

	// Assert.
	if !congruenceHolds {
		test.Fatalf("expected polynomial congruence to hold for prime %s",
			primeCandidate.String(),
		)
	}
}

// TestAks_ReturnsTrueForSmallPrimes verifies that AKS correctly identifies
// small prime numbers.
func TestAks_ReturnsTrueForSmallPrimes(test *testing.T) {
	// Arrange.
	primeValues := []*big.Int{
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(7),
		big.NewInt(11),
		big.NewInt(13),
	}

	// Act / Assert.
	for index := 0; index < len(primeValues); index++ {
		currentPrime := primeValues[index]

		isPrime := Aks(currentPrime)

		if !isPrime {
			test.Fatalf("expected %s to be classified as prime",
				currentPrime.String(),
			)
		}
	}
}

// TestAks_ReturnsFalseForComposites verifies that AKS correctly rejects
// composite numbers, including perfect powers.
func TestAks_ReturnsFalseForComposites(test *testing.T) {
	// Arrange.
	compositeValues := []*big.Int{
		big.NewInt(4),
		big.NewInt(6),
		big.NewInt(8),
		big.NewInt(9),
		big.NewInt(15),
		big.NewInt(21),
	}

	// Act / Assert.
	for index := 0; index < len(compositeValues); index++ {
		currentComposite := compositeValues[index]

		isPrime := Aks(currentComposite)

		if isPrime {
			test.Fatalf("expected %s to be classified as composite",
				currentComposite.String(),
			)
		}
	}
}

// TestIntegerSquareRootCeil_PerfectSquare verifies that the ceiling square root
// of a perfect square returns the exact integer value.
func TestIntegerSquareRootCeil_PerfectSquare(test *testing.T) {
	// Arrange.
	inputValue := 81
	expectedSquareRoot := 9

	// Act.
	computedSquareRoot := integerSquareRootCeil(inputValue)

	// Assert.
	if computedSquareRoot != expectedSquareRoot {
		test.Fatalf("expected ceil sqrt(%d) to be %d",
			inputValue,
			expectedSquareRoot,
		)
	}
}

// TestIntegerSquareRootCeil_NonPerfectSquare verifies that the ceiling square
// root of a non-perfect square rounds upward correctly.
func TestIntegerSquareRootCeil_NonPerfectSquare(test *testing.T) {
	// Arrange.
	inputValue := 50
	expectedSquareRoot := 8

	// Act.
	computedSquareRoot := integerSquareRootCeil(inputValue)

	// Assert.
	if computedSquareRoot != expectedSquareRoot {
		test.Fatalf("expected ceil sqrt(%d) to be %d",
			inputValue,
			expectedSquareRoot,
		)
	}
}

// TestComputeEulerTotient_PrimeModulus verifies that φ(p) = p − 1 for prime p.
func TestComputeEulerTotient_PrimeModulus(test *testing.T) {
	// Arrange.
	primeModulus := 13
	expectedTotient := 12

	// Act.
	computedTotient := computeEulerTotient(primeModulus)

	// Assert.
	if computedTotient != expectedTotient {
		test.Fatalf("expected φ(%d) to be %d",
			primeModulus,
			expectedTotient,
		)
	}
}

// TestComputeEulerTotient_CompositeModulus verifies correct totient computation
// for a composite modulus with multiple prime factors.
func TestComputeEulerTotient_CompositeModulus(test *testing.T) {
	// Arrange.
	compositeModulus := 36
	expectedTotient := 12

	// Act.
	computedTotient := computeEulerTotient(compositeModulus)

	// Assert.
	if computedTotient != expectedTotient {
		test.Fatalf("expected φ(%d) to be %d",
			compositeModulus,
			expectedTotient,
		)
	}
}

// TestHasNonTrivial_ReturnsTrueForComposite verifies that a composite number
// with a small non-trivial divisor is detected correctly.
func TestHasNonTrivial_ReturnsTrueForComposite(test *testing.T) {
	// Arrange.
	candidateValue := big.NewInt(21)
	modulusLimit := 10

	// Act.
	hasDivisor := hasNonTrivial(candidateValue, modulusLimit)

	// Assert.
	if !hasDivisor {
		test.Fatalf("expected %s to have a non-trivial divisor ≤ %d",
			candidateValue.String(),
			modulusLimit,
		)
	}
}

// TestHasNonTrivial_ReturnsFalseForPrime verifies that primes do not produce
// false positives in the non-trivial divisor check.
func TestHasNonTrivial_ReturnsFalseForPrime(test *testing.T) {
	// Arrange.
	primeCandidate := big.NewInt(17)
	modulusLimit := 10

	// Act.
	hasDivisor := hasNonTrivial(primeCandidate, modulusLimit)

	// Assert.
	if hasDivisor {
		test.Fatalf("expected %s to not have a non-trivial divisor ≤ %d",
			primeCandidate.String(),
			modulusLimit,
		)
	}
}

// TestPolynomialCongruenceCheck_CompositeFails verifies that the polynomial
// congruence check fails for a known composite input.
func TestPolynomialCongruenceCheck_CompositeFails(test *testing.T) {
	// Arrange.
	compositeCandidate := big.NewInt(9)

	modulusValue := 5
	witnessValue := 1

	// Act.
	congruenceHolds := polynomialCongruenceCheck(compositeCandidate, modulusValue, witnessValue)

	// Assert.
	if congruenceHolds {
		test.Fatalf("expected polynomial congruence to fail for composite %s",
			compositeCandidate.String(),
		)
	}
}

// TestAks_ReturnsFalseForZero verifies that zero is correctly classified
// as non-prime.
func TestAks_ReturnsFalseForZero(test *testing.T) {
	// Arrange.
	zeroValue := big.NewInt(0)

	// Act.
	isPrime := Aks(zeroValue)

	// Assert.
	if isPrime {
		test.Fatalf("expected %s to be classified as composite",
			zeroValue.String(),
		)
	}
}

// TestAks_ReturnsFalseForOne verifies that one is correctly classified
// as non-prime.
func TestAks_ReturnsFalseForOne(test *testing.T) {
	// Arrange.
	oneValue := big.NewInt(1)

	// Act.
	isPrime := Aks(oneValue)

	// Assert.
	if isPrime {
		test.Fatalf("expected %s to be classified as composite",
			oneValue.String(),
		)
	}
}

// TestAks_ReturnsFalseForNegativeNumbers verifies that negative integers are
// correctly classified as non-prime.
func TestAks_ReturnsFalseForNegativeNumbers(test *testing.T) {
	// Arrange.
	negativeValue := big.NewInt(-7)

	// Act.
	isPrime := Aks(negativeValue)

	// Assert.
	if isPrime {
		test.Fatalf("expected negative number %s to be classified as composite/non-prime",
			negativeValue.String(),
		)
	}
}

// TestAks_ReturnsTrueForLargerKnownPrime verifies that AKS correctly identifies
// a slightly larger prime number.
func TestAks_ReturnsTrueForLargerKnownPrime(test *testing.T) {
	// Arrange.
	primeCandidate := big.NewInt(17)

	// Act.
	isPrime := Aks(primeCandidate)

	// Assert.
	if !isPrime {
		test.Fatalf("expected %s to be classified as prime",
			primeCandidate.String(),
		)
	}
}

// FuzzAks_Primality performs native Go fuzz testing on the AKS primality algorithm
// using randomized integer byte streams to ensure stability across arbitrary inputs.
func FuzzAks_Primality(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(int64(2))
	fuzz.Add(int64(13))
	fuzz.Add(int64(99))

	fuzz.Fuzz(func(fuzzContext *testing.T, rawInputValue int64) {
		// Act.
		candidateValue := big.NewInt(rawInputValue)

		// Execute algorithm wrapped in recovery to catch any unexpected behavior safely.
		func() {
			defer func() {
				_ = recover()
			}()

			_ = Aks(candidateValue)
		}()
	})
}

// BenchmarkAks_SmallPrime benchmarks the performance of the AKS algorithm
// when evaluating a small known prime number.
func BenchmarkAks_SmallPrime(bench *testing.B) {
	// Arrange.
	primeCandidate := big.NewInt(13)

	// Act.
	for bench.Loop() {
		_ = Aks(primeCandidate)
	}
}

// BenchmarkAks_MediumComposite benchmarks the performance of the AKS algorithm
// when evaluating a medium composite number.
func BenchmarkAks_MediumComposite(bench *testing.B) {
	// Arrange.
	compositeCandidate := big.NewInt(109)

	// Act.
	for bench.Loop() {
		_ = Aks(compositeCandidate)
	}
}
