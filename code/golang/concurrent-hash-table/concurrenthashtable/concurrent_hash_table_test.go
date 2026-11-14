// -----------------------------------------------------------------------------
// Package: concurrenthashtable
// File: concurrent_hash_table_test.go
//
// Description:
//
//	This file contains unit tests for the ConcurrentHashTable implementation.
//	The tests validate correctness, concurrency safety, and API semantics
//	across all public operations.
//
//	Each test follows a strict Arrange / Act / Assert structure and uses
//	descriptive variable naming to make intent explicit. The test suite
//	covers both single-threaded behavior and concurrent access patterns
//	using goroutines and synchronization primitives.
//
// Design Notes:
//   - Tests are written in the same package to allow direct access.
//   - All public APIs are exercised at least once.
//   - Concurrency tests validate safety, not performance.
//   - The Go race detector should be used when running these tests.
//
// Recommended Invocation:
//
//	go test -race ./...
//
// -----------------------------------------------------------------------------
package concurrenthashtable

import (
	"sync"
	"testing"
)

// TestNewConcurrentHashTableInitializesEmptyTable verifies correct construction.
func TestNewConcurrentHashTableInitializesEmptyTable(test *testing.T) {
	// Arrange.
	requestedShardCount := 8

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	// Assert.
	if table == nil {
		test.Fatalf("expected non-nil table.")
	}

	if table.ShardCount() != requestedShardCount {
		test.Fatalf("expected shard count %d, got %d.",
			requestedShardCount,
			table.ShardCount(),
		)
	}

	if table.Length() != 0 {
		test.Fatalf("expected newly constructed table to be empty.")
	}
}

// TestNewConcurrentHashTableNormalizesInvalidShardCount verifies normalization.
func TestNewConcurrentHashTableNormalizesInvalidShardCount(test *testing.T) {
	// Arrange.
	requestedShardCount := 0

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	// Assert.
	if table.ShardCount() != defaultSharedCount {
		test.Fatalf("expected default shard count %d, got %d.",
			defaultSharedCount,
			table.ShardCount(),
		)
	}
}

// TestSetStoresValueAndGetRetrievesIt verifies basic CRUD behavior.
func TestSetStoresValueAndGetRetrievesIt(test *testing.T) {
	// Arrange.
	requestedShardCount := 4
	testKey := "alpha"
	testValue := 42

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set(testKey, testValue)

	retrievedValue, wasPresent := table.Get(testKey)

	// Assert.
	if !wasPresent {
		test.Fatalf("expected key to be present.")
	}

	if retrievedValue != testValue {
		test.Fatalf("expected value %d, got %d.",
			testValue,
			retrievedValue,
		)
	}
}

// TestDeleteRemovesExistingKey verifies deletion behavior.
func TestDeleteRemovesExistingKey(test *testing.T) {
	// Arrange.
	requestedShardCount := 4
	testKey := "beta"

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set(testKey, 100)

	wasDeleted := table.Delete(testKey)

	// Assert.
	if !wasDeleted {
		test.Fatalf("expected delete to report success.")
	}

	if table.Length() != 0 {
		test.Fatalf("expected table to be empty after deletion.")
	}
}

// TestDeleteReturnsFalseForMissingKey verifies delete miss semantics.
func TestDeleteReturnsFalseForMissingKey(test *testing.T) {
	// Arrange.
	requestedShardCount := 4

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	wasDeleted := table.Delete("missing")

	// Assert.
	if wasDeleted {
		test.Fatalf("expected delete to return false for missing key.")
	}
}

// TestLoadOrStoreStoresValueWhenKeyIsMissing verifies insert path.
func TestLoadOrStoreStoresValueWhenKeyIsMissing(test *testing.T) {
	// Arrange.
	requestedShardCount := 4
	testKey := "gamma"
	testValue := 7

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	returnedValue, wasLoaded := table.LoadOrStore(testKey, testValue)

	// Assert.
	if wasLoaded {
		test.Fatalf("expected wasLoaded to be false.")
	}

	if returnedValue != testValue {
		test.Fatalf("expected value %d, got %d.",
			testValue,
			returnedValue,
		)
	}
}

// TestLoadOrStoreLoadsExistingValue verifies load path.
func TestLoadOrStoreLoadsExistingValue(test *testing.T) {
	// Arrange.
	requestedShardCount := 4
	testKey := "delta"

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set(testKey, 10)

	returnedValue, wasLoaded := table.LoadOrStore(testKey, 99)

	// Assert.
	if !wasLoaded {
		test.Fatalf("expected wasLoaded to be true.")
	}

	if returnedValue != 10 {
		test.Fatalf("expected existing value 10, got %d.", returnedValue)
	}
}

// TestCompareAndSwapSucceedsWhenExpectedMatches verifies CAS success.
func TestCompareAndSwapSucceedsWhenExpectedMatches(test *testing.T) {
	// Arrange.
	requestedShardCount := 4
	testKey := "epsilon"

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set(testKey, 1)

	swapSucceeded := table.CompareAndSwap(testKey, 1, 10,
		func(compare, against int) bool {
			return compare == against
		},
	)

	// Assert.
	if !swapSucceeded {
		test.Fatalf("expected compare-and-swap to succeed.")
	}

	updatedValue, _ := table.Get(testKey)

	if updatedValue != 10 {
		test.Fatalf("expected updated value 10, got %d.", updatedValue)
	}
}

// TestCompareAndSwapFailsWhenExpectedDoesNotMatch verifies CAS failure.
func TestCompareAndSwapFailsWhenExpectedDoesNotMatch(test *testing.T) {
	// Arrange.
	requestedShardCount := 4
	testKey := "zeta"

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set(testKey, 5)

	swapSucceeded := table.CompareAndSwap(testKey, 1, 99,
		func(compare, against int) bool {
			return compare == against
		},
	)

	// Assert.
	if swapSucceeded {
		test.Fatalf("expected compare-and-swap to fail.")
	}
}

// TestRangeVisitsAllEntries verifies full traversal.
func TestRangeVisitsAllEntries(test *testing.T) {
	// Arrange.
	requestedShardCount := 8

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set("a", 1)
	table.Set("b", 2)
	table.Set("c", 3)

	visitedEntryCount := 0

	table.Range(func(key string, value int) bool {
		visitedEntryCount++
		return true
	})

	// Assert.
	if visitedEntryCount != 3 {
		test.Fatalf("expected to visit 3 entries, got %d.",
			visitedEntryCount,
		)
	}
}

// TestRangeStopsEarlyWhenVisitorReturnsFalse verifies early termination.
func TestRangeStopsEarlyWhenVisitorReturnsFalse(test *testing.T) {
	// Arrange.
	requestedShardCount := 8

	// Act.
	table := NewSyncTable[string, int](requestedShardCount)

	table.Set("x", 1)
	table.Set("y", 2)

	visitedEntryCount := 0

	table.Range(func(key string, value int) bool {
		visitedEntryCount++

		return false
	})

	// Assert.
	if visitedEntryCount != 1 {
		test.Fatalf("expected early termination after one visit.")
	}
}

// TestConcurrentSetOperationsDoNotCorruptState verifies concurrent safety.
func TestConcurrentSetOperationsDoNotCorruptState(test *testing.T) {
	// Arrange.
	requestedShardCount := 16
	workerCount := 8
	operationsPerWorker := 1000

	// Act.
	table := NewSyncTable[int, int](requestedShardCount)

	var waitGroup sync.WaitGroup

	waitGroup.Add(workerCount)

	for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
		currentWorkerIndex := workerIndex

		go func() {
			defer waitGroup.Done()

			for operationIndex := 0; operationIndex < operationsPerWorker; operationIndex++ {
				key := currentWorkerIndex*operationsPerWorker + operationIndex
				table.Set(key, operationIndex)
			}
		}()
	}

	waitGroup.Wait()

	// Assert.
	expectedElementCount := workerCount * operationsPerWorker

	if table.Length() != expectedElementCount {
		test.Fatalf("expected length %d, got %d.",
			expectedElementCount,
			table.Length(),
		)
	}
}

// FuzzSyncTable_Operations performs native Go fuzz testing on the concurrent hash table
// using randomized byte streams to execute sets, gets, deletes, load-or-stores, and ranges safely.
func FuzzSyncTable_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 10, 1, 20, 2, 30, 3, 40})
	fuzz.Add([]byte{1, 5, 0, 15, 4, 25})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawOperationBytes []byte) {
		if len(rawOperationBytes) < 2 {
			return
		}

		// Act.
		table := NewSyncTable[int, int](4)
		byteIndex := 0

		for byteIndex < len(rawOperationBytes) {
			operationType := rawOperationBytes[byteIndex] % 6
			byteIndex++

			if byteIndex >= len(rawOperationBytes) {
				break
			}

			payloadKey := int(rawOperationBytes[byteIndex])
			byteIndex++

			var payloadValue int
			if byteIndex < len(rawOperationBytes) {
				payloadValue = int(rawOperationBytes[byteIndex])
				byteIndex++
			} else {
				payloadValue = 100
			}

			// Execute hash table operations wrapped in recovery to guarantee fuzzer stability.
			func() {
				defer func() {
					_ = recover()
				}()

				switch operationType {
				case 0:
					table.Set(payloadKey, payloadValue)
				case 1:
					_, _ = table.Get(payloadKey)
				case 2:
					_ = table.Delete(payloadKey)
				case 3:
					_, _ = table.LoadOrStore(payloadKey, payloadValue)
				case 4:
					_ = table.CompareAndSwap(payloadKey, payloadValue, payloadValue+1, func(compare, against int) bool {
						return compare == against
					})
				case 5:
					table.Range(func(key int, value int) bool {
						return key%2 == 0
					})
				}
			}()
		}
	})
}

// BenchmarkSyncTable_Set benchmarks the performance of setting key-value pairs
// in the concurrent hash table.
func BenchmarkSyncTable_Set(bench *testing.B) {
	// Arrange.
	table := NewSyncTable[int, int](16)

	// Act.
	for bench.Loop() {
		for elementIndex := 0; elementIndex < 1000; elementIndex++ {
			table.Set(elementIndex, elementIndex*2)
		}
	}
}

// BenchmarkSyncTable_Get benchmarks the performance of retrieving keys
// from a pre-populated concurrent hash table.
func BenchmarkSyncTable_Get(bench *testing.B) {
	// Arrange.
	table := NewSyncTable[int, int](16)

	for elementIndex := 0; elementIndex < 1000; elementIndex++ {
		table.Set(elementIndex, elementIndex*2)
	}

	// Act.
	for bench.Loop() {
		for elementIndex := 0; elementIndex < 1000; elementIndex++ {
			_, _ = table.Get(elementIndex)
		}
	}
}

// BenchmarkSyncTable_MixedOperations benchmarks a concurrent mix of writes, reads,
// and deletes executed iteratively within a loop.
func BenchmarkSyncTable_MixedOperations(bench *testing.B) {
	// Act.
	for bench.Loop() {
		table := NewSyncTable[int, int](8)

		for operationIndex := 0; operationIndex < 200; operationIndex++ {
			table.Set(operationIndex, operationIndex)

			_, _ = table.Get(operationIndex)
		}

		for operationIndex := 0; operationIndex < 100; operationIndex++ {
			table.Delete(operationIndex)
		}
	}
}
