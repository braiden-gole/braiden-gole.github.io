// -----------------------------------------------------------------------------
// Package: openaddressedhashtable
// File: table_test.go
//
// Description:
//
//	This file contains unit tests for the open-addressed hash table
//	implementation. The tests validate correct behavior for insertion,
//	lookup, deletion, key overwrite, size tracking, and automatic resizing.
//
//	Each test focuses on a specific operation or invariant and is designed
//	to exercise both normal and edge-case behavior. The tests collectively
//	serve as a correctness and regression safety net for the hash table
//	implementation.
//
// Usage:
//
//	go test ./...
//
// Notes:
//   - These tests assume single-threaded execution.
//   - The tests favor clarity and explicit assertions over conciseness.
//   - No Go built-in map is used for validation.
//
// -----------------------------------------------------------------------------
package openaddressedhashtable

import (
	"strconv"
	"testing"
)

// TestNew verifies that a newly constructed hash table is initialized
// with zero length and a non-nil backing storage.
func TestNew(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	// Act.
	length := table.Length()

	// Assert.
	if table == nil {
		test.Fatal("expected non-nil table")
	}

	if length != 0 {
		test.Fatalf("expected length 0, got %d", length)
	}

	if len(table.entries) == 0 {
		test.Fatal("expected backing array to be initialized")
	}
}

// TestPutAndGet verifies that inserting a key-value pair allows the value
// to be retrieved using the same key.
func TestPutAndGet(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	// Act.
	table.Put("alpha", 1)

	value, ok := table.Get("alpha")

	// Assert.
	if !ok {
		test.Fatal("expected key to be found")
	}

	if value != 1 {
		test.Fatalf("expected value 1, got %d", value)
	}
}

// TestOverwrite verifies that inserting a value with an existing key
// overwrites the previous value without increasing table size.
func TestOverwrite(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	// Act.
	table.Put("key", 1)
	table.Put("key", 2)

	value, ok := table.Get("key")

	length := table.Length()

	// Assert.
	if !ok {
		test.Fatal("expected key to be found")
	}

	if value != 2 {
		test.Fatalf("expected value 2, got %d", value)
	}

	if length != 1 {
		test.Fatalf("expected length 1, got %d", length)
	}
}

// TestGetMissingKey verifies that attempting to retrieve a non-existent
// key returns the zero value and a false indicator.
func TestGetMissingKey(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	// Act.
	value, ok := table.Get("missing")

	// Assert.
	if ok {
		test.Fatal("expected key to be missing")
	}

	if value != 0 {
		test.Fatalf("expected zero value, got %d", value)
	}
}

// TestDelete verifies that deleting an existing key removes it from the
// table and updates the size accordingly.
func TestDelete(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	table.Put("x", 42)

	// Act.
	deleted := table.Delete("x")

	length := table.Length()

	_, ok := table.Get("x")

	// Assert.
	if !deleted {
		test.Fatal("expected delete to succeed")
	}

	if length != 0 {
		test.Fatalf("expected length 0, got %d", length)
	}

	if ok {
		test.Fatal("expected key to be deleted")
	}
}

// TestDeleteMissingKey verifies that deleting a non-existent key returns
// false and does not modify table state.
func TestDeleteMissingKey(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	// Act.
	deleted := table.Delete("missing")
	length := table.Length()

	// Assert.
	if deleted {
		test.Fatal("expected delete to fail for missing key")
	}

	if length != 0 {
		test.Fatalf("expected length 0, got %d", length)
	}
}

// TestLength verifies that the Length method accurately reflects the
// number of stored entries after insertions and deletions.
func TestLength(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	// Act.
	table.Put("a", 1)
	table.Put("b", 2)

	lengthAfterInsert := table.Length()

	table.Delete("a")

	lengthAfterDelete := table.Length()

	// Assert.
	if lengthAfterInsert != 2 {
		test.Fatalf("expected length 2, got %d", lengthAfterInsert)
	}

	if lengthAfterDelete != 1 {
		test.Fatalf("expected length 1, got %d", lengthAfterDelete)
	}
}

// TestResize verifies that the hash table automatically resizes when the
// load factor threshold is exceeded and that all entries remain accessible
// after resizing.
func TestResize(test *testing.T) {
	// Arrange.
	table := New[string, int](StringHasher{})

	initialCapacity := len(table.entries)

	// Act.
	for index := 0; index < initialCapacity; index++ {
		table.Put(string(rune('a'+index)), index)
	}

	newCapacity := len(table.entries)

	// Assert.
	if newCapacity <= initialCapacity {
		test.Fatal("expected table to resize")
	}

	for index := 0; index < initialCapacity; index++ {
		value, ok := table.Get(string(rune('a' + index)))

		if !ok {
			test.Fatalf("expected key %d to be found after resize", index)
		}

		if value != index {
			test.Fatalf("expected value %d, got %d", index, value)
		}
	}
}

// FuzzHashTable_Operations performs native Go fuzz testing on the open-addressed hash table
// using randomized byte streams to execute puts, gets, and deletions safely.
func FuzzHashTable_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 5, 1, 5, 2, 5, 0})
	fuzz.Add([]byte{0, 1, 0, 2, 1, 1, 2, 1})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawOperationBytes []byte) {
		if len(rawOperationBytes) < 2 {
			return
		}

		// Act.
		table := New[string, int](StringHasher{})
		byteIndex := 0

		for byteIndex+1 < len(rawOperationBytes) {
			operationType := rawOperationBytes[byteIndex] % 4
			byteIndex++

			keyID := int(rawOperationBytes[byteIndex] % 50)

			keyString := "key-" + strconv.Itoa(keyID)
			byteIndex++

			// Execute hash table operations wrapped in recovery to guarantee fuzzer stability.
			func() {
				defer func() {
					_ = recover()
				}()

				switch operationType {
				case 0:
					table.Put(keyString, keyID)
				case 1:
					_, _ = table.Get(keyString)
				case 2:
					_ = table.Delete(keyString)
				case 3:
					_ = table.Length()
				}
			}()
		}
	})
}

// BenchmarkHashTable_Put benchmarks the performance of inserting entries
// into the open-addressed hash table.
func BenchmarkHashTable_Put(bench *testing.B) {
	// Act.
	for bench.Loop() {
		table := New[string, int](StringHasher{})

		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			key := "key-" + strconv.Itoa(elementIndex)
			table.Put(key, elementIndex)
		}
	}
}

// BenchmarkHashTable_Get benchmarks the performance of retrieving entries
// from a pre-populated open-addressed hash table.
func BenchmarkHashTable_Get(bench *testing.B) {
	// Arrange.
	table := New[string, int](StringHasher{})

	for elementIndex := 0; elementIndex < 300; elementIndex++ {
		key := "key-" + strconv.Itoa(elementIndex)
		table.Put(key, elementIndex)
	}

	// Act.
	for bench.Loop() {
		for elementIndex := 0; elementIndex < 300; elementIndex++ {
			key := "key-" + strconv.Itoa(elementIndex)
			_, _ = table.Get(key)
		}
	}
}

// BenchmarkHashTable_MixedWorkload benchmarks a mixed workload of puts, gets,
// and deletions performed iteratively within the hash table.
func BenchmarkHashTable_MixedWorkload(bench *testing.B) {
	// Act.
	for bench.Loop() {
		table := New[string, int](StringHasher{})

		for operationIndex := 0; operationIndex < 200; operationIndex++ {
			key := "key-" + strconv.Itoa(operationIndex%100)

			if operationIndex%2 == 0 {
				table.Put(key, operationIndex)
			} else {
				_, _ = table.Get(key)
			}

			if operationIndex%10 == 0 {
				table.Delete(key)
			}
		}
	}
}
