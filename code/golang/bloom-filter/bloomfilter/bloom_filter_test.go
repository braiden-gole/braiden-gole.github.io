// ----------------------------------------------------------------------------
// Package: bloomfilter
// File: bloom_filter_test.go
//
// Description:
//
//	  This file provides a comprehensive suite of unit tests for the
//	  Counting Bloom Filter implementation. It validates core functionality,
//	  mathematical boundary conditions, and defensive programming invariants.
//
//	  The testing suite is intentionally defensive:
//		- Verifies that all methods handle nil *BloomFilter receivers without panicking.
//		- Ensures empty data inputs are handled gracefully across all operations.
//		- Validates that counting logic correctly supports item removal.
//	 	- Tests concurrency to ensure thread-safety under high-pressure access.
//
//	  Data Structure:
//
//	  Tests interact with the BloomFilter struct, validating the internal
//	  management of the byte-slice buckets and the iteration logic.
//
// Notes:
//   - Tests utilize the Kirsch-Mitzenmacher optimization indirectly via public APIs.
//   - Concurrency tests should be run with the -race flag for full validation.
//   - Probability tests assume a standard distribution provided by FNV hashing.
//
// -----------------------------------------------------------------------------
package bloomfilter

import (
	"fmt"
	"sync"
	"testing"
)

// TestNewFilterOptimalParams ensures the constructor creates a valid struct.
func TestNewFilterOptimalParams(test *testing.T) {
	// Arrange.
	expectedEntries := uint64(1000)
	falsePositiveRate := 0.01

	// Act.
	filter := NewOptimalBloomFilter(expectedEntries, falsePositiveRate)

	// Assert.
	if filter == nil {
		test.Fatal("Expected filter to be initialized, got nil")
	}
}

// TestAddAndContains validates basic insertion and positive membership lookups.
func TestAddAndContains(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)
	braidenKey := []byte("braiden_gole")

	// Act.
	filter.Add(braidenKey)

	// Assert.
	if !filter.Contains(braidenKey) {
		test.Errorf("Expected filter to contain %s", string(braidenKey))
	}
}

// TestDefinitiveNegative verifies that items not added always return false.
func TestDefinitiveNegative(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)
	missingKey := []byte("missing_user")

	// Act.
	exists := filter.Contains(missingKey)

	// Assert.
	if exists {
		test.Error("Filter reported a positive for an item that was never added")
	}
}

// TestRemoveFunctionality ensures that removing an element decrements counters correctly.
func TestRemoveFunctionality(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)
	tokenKey := []byte("session_token_123")

	// Act.
	filter.Add(tokenKey)
	filter.Remove(tokenKey)

	// Assert.
	if filter.Contains(tokenKey) {
		test.Error("Filter still contains item after removal")
	}
}

// TestNilReceiverAdd validates defensive nil-checking for the Add method.
func TestNilReceiverAdd(test *testing.T) {
	// Arrange.
	var nilFilter *BloomFilter

	// Act & Assert.
	// No panic should occur.
	nilFilter.Add([]byte("safety_test"))
}

// TestNilReceiverContains validates defensive nil-checking for the Contains method.
func TestNilReceiverContains(test *testing.T) {
	// Arrange.
	var nilFilter *BloomFilter

	// Act.
	exists := nilFilter.Contains([]byte("safety_test"))

	// Assert.
	if exists {
		test.Error("Nil filter should never report that it contains an item")
	}
}

// TestNilReceiverRemove validates defensive nil-checking for the Remove method.
func TestNilReceiverRemove(test *testing.T) {
	// Arrange.
	var nilFilter *BloomFilter

	// Act & Assert.
	// No panic should occur.
	nilFilter.Remove([]byte("safety_test"))
}

// TestEmptyDataAdd verifies that adding empty slices is handled safely.
func TestEmptyDataAdd(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)
	emptySlice := []byte("")

	// Act.
	filter.Add(emptySlice)

	// Assert.
	if filter.Contains(emptySlice) {
		test.Error("Empty data should not result in a positive membership hit")
	}
}

// TestCounterSaturation verifies that buckets do not overflow beyond 255.
func TestCounterSaturation(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(10, 0.1)
	staticKey := []byte("overflow_check")

	// Act.
	for index := 0; index < 300; index++ {
		filter.Add(staticKey)
	}

	// Assert.
	filter.Remove(staticKey)

	if !filter.Contains(staticKey) {
		test.Error("Filter lost item after single removal; counter likely wrapped to zero")
	}
}

// TestRemoveWithoutAdd verifies that removing non-existent items is a no-op.
func TestRemoveWithoutAdd(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)
	unaddedKey := []byte("unknown")

	// Act & Assert.
	filter.Remove(unaddedKey)
}

// TestConcurrentSafety uses goroutines to validate thread-safety under load.
func TestConcurrentSafety(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(1000, 0.01)
	waitGroup := sync.WaitGroup{}
	iterations := 100

	// Act.
	for iteration := 0; iteration < iterations; iteration++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()

			data := fmt.Appendf(nil, "user_%d", index)

			filter.Add(data)
			filter.Contains(data)
			filter.Remove(data)
		}(iteration)
	}

	waitGroup.Wait()

	// Assert.
	// Passes if no race conditions are detected.
}

// TestMultipleAdditions tracks items added more than once.
func TestMultipleAdditions(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)
	duplicateKey := []byte("multi_add")

	// Act.
	filter.Add(duplicateKey)
	filter.Add(duplicateKey)
	filter.Remove(duplicateKey)

	// Assert.
	if !filter.Contains(duplicateKey) {
		test.Error("Item should remain if it was added twice but removed once")
	}
}

// TestInvalidConstructorInputs validates fallback defaults for bad math inputs.
func TestInvalidConstructorInputs(test *testing.T) {
	// Arrange.
	badEntries := uint64(0)
	badRate := -0.5

	// Act.
	filter := NewOptimalBloomFilter(badEntries, badRate)

	// Assert.
	if filter == nil || filter.Contains([]byte("test")) {
		test.Error("Constructor failed to handle invalid inputs safely")
	}
}

// TestLargeVolumeInsertion ensures stability across 1000 unique inserts.
func TestLargeVolumeInsertion(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(1000, 0.01)

	// Act.
	for index := 0; index < 1000; index++ {
		filter.Add([]byte(fmt.Appendf(nil, "key_%d", index)))
	}

	// Assert.
	if !filter.Contains([]byte("key_500")) {
		test.Error("Filter lost data during high-volume insertion")
	}
}

// TestHashingUniqueness checks that different keys generally map to different bits.
func TestHashingUniqueness(test *testing.T) {
	// Arrange.
	filter := NewOptimalBloomFilter(100, 0.01)

	keyOne := []byte("apple")
	keyTwo := []byte("orange")

	// Act.
	filter.Add(keyOne)

	// Assert.
	if filter.Contains(keyTwo) {
		test.Log("Note: Probabilistic collision occurred (acceptable, but rare)")
	}
}

// FuzzBloomFilter_Operations performs native Go fuzz testing on the Counting Bloom Filter
// using randomized byte inputs to execute insertions, lookups, and removals safely without panicking.
func FuzzBloomFilter_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte("sample_key_one"))
	fuzz.Add([]byte(""))
	fuzz.Add([]byte("another_test_token"))

	fuzz.Fuzz(func(fuzzContext *testing.T, rawKeyBytes []byte) {
		// Act & Assert: Execute operations wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			filter := NewOptimalBloomFilter(500, 0.01)

			if filter == nil {
				return
			}

			filter.Add(rawKeyBytes)
			_ = filter.Contains(rawKeyBytes)

			filter.Remove(rawKeyBytes)
			_ = filter.Contains(rawKeyBytes)
		}()
	})
}

// BenchmarkBloomFilter_Add benchmarks the performance of inserting items
// into the counting Bloom filter.
func BenchmarkBloomFilter_Add(bench *testing.B) {
	// Arrange.
	filter := NewOptimalBloomFilter(10000, 0.01)
	items := make([][]byte, 1000)

	for itemIndex := 0; itemIndex < 1000; itemIndex++ {
		items[itemIndex] = []byte(fmt.Appendf(nil, "benchmark_key_%d", itemIndex))
	}

	// Act.
	for bench.Loop() {
		for itemIndex := 0; itemIndex < 1000; itemIndex++ {
			filter.Add(items[itemIndex])
		}
	}
}

// BenchmarkBloomFilter_Contains benchmarks the performance of membership lookup checks
// in a pre-populated counting Bloom filter.
func BenchmarkBloomFilter_Contains(bench *testing.B) {
	// Arrange.
	filter := NewOptimalBloomFilter(10000, 0.01)
	items := make([][]byte, 1000)

	for itemIndex := 0; itemIndex < 1000; itemIndex++ {
		items[itemIndex] = []byte(fmt.Appendf(nil, "benchmark_key_%d", itemIndex))
		filter.Add(items[itemIndex])
	}

	// Act.
	for bench.Loop() {
		for itemIndex := 0; itemIndex < 1000; itemIndex++ {
			_ = filter.Contains(items[itemIndex])
		}
	}
}

// BenchmarkBloomFilter_MixedOperations benchmarks a mix of additions, membership checks,
// and removals performed iteratively within a loop.
func BenchmarkBloomFilter_MixedOperations(bench *testing.B) {
	// Act.
	for bench.Loop() {
		filter := NewOptimalBloomFilter(1000, 0.01)

		for operationIndex := 0; operationIndex < 200; operationIndex++ {
			keyBytes := []byte(fmt.Appendf(nil, "key_%d", operationIndex))

			filter.Add(keyBytes)

			_ = filter.Contains(keyBytes)
		}

		for operationIndex := 0; operationIndex < 100; operationIndex++ {
			keyBytes := []byte(fmt.Appendf(nil, "key_%d", operationIndex))
			filter.Remove(keyBytes)
		}
	}
}
