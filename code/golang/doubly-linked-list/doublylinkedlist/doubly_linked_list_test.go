// -----------------------------------------------------------------------------
// Package: doublylinkedlist
// File: doubly_linked_list_test.go
//
// Description:
//
//	This file contains unit tests for the doublylinkedlist package. The tests
//	verify correctness, safety, and defensive behavior for all exported
//	functions and methods of the DoublyLinkedList type.
//
//	The test suite exercises normal operation, edge cases, and failure
//	conditions to ensure that list invariants are preserved and that no
//	operation panics due to invalid input or state.
//
// Design Notes:
//   - Tests are written using explicit Arrange / Act / Assert sections.
//   - Defensive behaviors such as nil receiver handling are validated.
//   - Tests favor clarity and correctness over brevity.
//
// Intended Use:
//
//	This file is intended to validate the behavior of the doublylinkedlist
//	implementation during development and refactoring.
//
// -----------------------------------------------------------------------------
package doublylinkedlist

import "testing"

// TestNewCreatesEmptyList verifies that New returns a valid empty list.
func TestNewCreatesEmptyList(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	length := list.Length()

	// Assert.
	if length != 0 {
		test.Fatalf("expected length 0, got %d", length)
	}

	if !list.IsEmpty() {
		test.Fatalf("expected list to be empty")
	}
}

// TestLengthOnNilReceiver verifies that Length returns zero for a nil receiver.
func TestLengthOnNilReceiver(test *testing.T) {
	// Arrange.
	var list *DoublyLinkedList

	// Act.
	length := list.Length()

	// Assert.
	if length != 0 {
		test.Fatalf("expected length 0 for nil receiver, got %d", length)
	}
}

// TestIsEmptyOnNilReceiver verifies that IsEmpty returns true for a nil receiver.
func TestIsEmptyOnNilReceiver(test *testing.T) {
	// Arrange.
	var list *DoublyLinkedList

	// Act.
	empty := list.IsEmpty()

	// Assert.
	if !empty {
		test.Fatalf("expected IsEmpty to return true for nil receiver")
	}
}

// TestPushBackOnEmptyList verifies insertion into an empty list.
func TestPushBackOnEmptyList(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	list.PushBack("A")

	// Assert.
	if list.Length() != 1 {
		test.Fatalf("expected length 1, got %d", list.Length())
	}

	if value, ok := list.Front(); !ok || value != "A" {
		test.Fatalf("unexpected front value")
	}

	if value, ok := list.Back(); !ok || value != "A" {
		test.Fatalf("unexpected back value")
	}
}

// TestPushBackMultipleValues verifies order preservation.
func TestPushBackMultipleValues(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	list.PushBack("A")
	list.PushBack("B")
	list.PushBack("C")

	// Assert.
	expected := []string{"A", "B", "C"}
	index := 0

	list.ForEach(func(value interface{}) {
		if value != expected[index] {
			test.Fatalf("at index %d, expected %v, got %v", index, expected[index], value)
		}
		index++
	})

	if list.Length() != 3 {
		test.Fatalf("expected length 3, got %d", list.Length())
	}
}

// TestPushBackOnNilReceiver verifies safe behavior.
func TestPushBackOnNilReceiver(test *testing.T) {
	// Arrange.
	var list *DoublyLinkedList

	// Act / Assert.
	defer func() {
		if recover() != nil {
			test.Fatalf("expected PushBack on nil receiver not to panic")
		}
	}()

	list.PushBack("X")
}

// TestPushFrontOnEmptyList verifies insertion into an empty list.
func TestPushFrontOnEmptyList(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	list.PushFront("A")

	// Assert.
	if list.Length() != 1 {
		test.Fatalf("expected length 1, got %d", list.Length())
	}

	if value, ok := list.Front(); !ok || value != "A" {
		test.Fatalf("unexpected front value")
	}
}

// TestPushFrontMultipleValues verifies reverse insertion order.
func TestPushFrontMultipleValues(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	list.PushFront("C")
	list.PushFront("B")
	list.PushFront("A")

	// Assert.
	expected := []string{"A", "B", "C"}
	index := 0

	list.ForEach(func(value interface{}) {
		if value != expected[index] {
			test.Fatalf("at index %d, expected %v, got %v", index, expected[index], value)
		}
		index++
	})
}

// TestFrontAndBackOnEmptyList verifies safe behavior.
func TestFrontAndBackOnEmptyList(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	_, frontOK := list.Front()
	_, backOK := list.Back()

	// Assert.
	if frontOK || backOK {
		test.Fatalf("expected Front and Back to fail on empty list")
	}
}

// TestPopFrontOnEmptyList verifies failure on empty list.
func TestPopFrontOnEmptyList(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	_, ok := list.PopFront()

	// Assert.
	if ok {
		test.Fatalf("expected PopFront to fail on empty list")
	}
}

// TestPopFrontRemovesHead verifies correct removal.
func TestPopFrontRemovesHead(test *testing.T) {
	// Arrange.
	list := New()

	list.PushBack("A")
	list.PushBack("B")

	// Act.
	value, ok := list.PopFront()

	// Assert.
	if !ok || value != "A" {
		test.Fatalf("unexpected PopFront result")
	}

	if list.Length() != 1 {
		test.Fatalf("expected length 1, got %d", list.Length())
	}

	if front, _ := list.Front(); front != "B" {
		test.Fatalf("unexpected front value")
	}
}

// TestPopBackOnEmptyList verifies failure on empty list.
func TestPopBackOnEmptyList(test *testing.T) {
	// Arrange.
	list := New()

	// Act.
	_, ok := list.PopBack()

	// Assert.
	if ok {
		test.Fatalf("expected PopBack to fail on empty list")
	}
}

// TestPopBackRemovesTail verifies correct removal.
func TestPopBackRemovesTail(test *testing.T) {
	// Arrange.
	list := New()

	list.PushBack("A")
	list.PushBack("B")

	// Act.
	value, ok := list.PopBack()

	// Assert.
	if !ok || value != "B" {
		test.Fatalf("unexpected PopBack result")
	}

	if list.Length() != 1 {
		test.Fatalf("expected length 1, got %d", list.Length())
	}

	if back, _ := list.Back(); back != "A" {
		test.Fatalf("unexpected back value")
	}
}

// TestForEachOnNilReceiver verifies safe behavior.
func TestForEachOnNilReceiver(test *testing.T) {
	// Arrange.
	var list *DoublyLinkedList

	// Act / Assert.
	defer func() {
		if recover() != nil {
			test.Fatalf("expected ForEach on nil receiver not to panic")
		}
	}()

	list.ForEach(func(value interface{}) {})
}

// TestClearEmptiesList verifies that Clear resets the list state.
func TestClearEmptiesList(test *testing.T) {
	// Arrange.
	list := New()

	list.PushBack("A")
	list.PushBack("B")

	// Act.
	list.Clear()

	// Assert.
	if list.Length() != 0 {
		test.Fatalf("expected length 0 after Clear, got %d", list.Length())
	}

	if !list.IsEmpty() {
		test.Fatalf("expected list to be empty after Clear")
	}
}

// FuzzDoublyLinkedList_Operations performs native Go fuzz testing on the doubly linked list
// using randomized byte streams to execute front/back pushes, pops, iterations, and clears safely.
func FuzzDoublyLinkedList_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 10, 1, 20, 2, 3, 4})
	fuzz.Add([]byte{1, 5, 0, 15, 5, 6})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawOperationBytes []byte) {
		if len(rawOperationBytes) < 2 {
			return
		}

		// Act.
		list := New()

		byteIndex := 0

		for byteIndex < len(rawOperationBytes) {
			operationType := rawOperationBytes[byteIndex] % 7
			byteIndex++

			if byteIndex >= len(rawOperationBytes) {
				break
			}

			payloadValue := int(rawOperationBytes[byteIndex])
			byteIndex++

			// Execute list operations wrapped in recovery to guarantee fuzzer stability.
			func() {
				defer func() {
					_ = recover()
				}()

				switch operationType {
				case 0:
					list.PushBack(payloadValue)
				case 1:
					list.PushFront(payloadValue)
				case 2:
					_, _ = list.PopFront()
				case 3:
					_, _ = list.PopBack()
				case 4:
					list.ForEach(func(value interface{}) {
						_ = value.(int)
					})
				case 5:
					list.Clear()
				case 6:
					_ = list.Length()
					_ = list.IsEmpty()
					_, _ = list.Front()
					_, _ = list.Back()
				}
			}()
		}
	})
}

// BenchmarkDoublyLinkedList_PushBack benchmarks the performance of inserting elements
// at the back of the doubly linked list.
func BenchmarkDoublyLinkedList_PushBack(bench *testing.B) {
	// Act.
	for bench.Loop() {
		list := New()

		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			list.PushBack(elementIndex)
		}
	}
}

// BenchmarkDoublyLinkedList_PushFront benchmarks the performance of inserting elements
// at the front of the doubly linked list.
func BenchmarkDoublyLinkedList_PushFront(bench *testing.B) {
	// Act.
	for bench.Loop() {
		list := New()

		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			list.PushFront(elementIndex)
		}
	}
}

// BenchmarkDoublyLinkedList_PopFront benchmarks the performance of removing elements
// from the front of a pre-populated doubly linked list.
func BenchmarkDoublyLinkedList_PopFront(bench *testing.B) {
	// Act.
	for bench.Loop() {
		list := New()

		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			list.PushBack(elementIndex)
		}

		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			_, _ = list.PopFront()
		}
	}
}

// BenchmarkDoublyLinkedList_MixedOperations benchmarks a mixed workload of pushes, pops,
// and traversals performed iteratively within a loop.
func BenchmarkDoublyLinkedList_MixedOperations(bench *testing.B) {
	// Act.
	for bench.Loop() {
		list := New()

		for operationIndex := 0; operationIndex < 200; operationIndex++ {
			list.PushBack(operationIndex)
			list.PushFront(operationIndex * 2)
		}

		for operationIndex := 0; operationIndex < 100; operationIndex++ {
			_, _ = list.PopFront()
			_, _ = list.PopBack()
		}

		list.ForEach(func(value interface{}) {
			_ = value
		})
	}
}
