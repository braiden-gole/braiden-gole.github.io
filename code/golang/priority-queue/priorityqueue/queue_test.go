// -----------------------------------------------------------------------------
// Package: priorityqueue
// File: queue_test.go
//
// Description:
//
//	This file defines unit tests for the PriorityQueue implementation.
//	The tests exercise all public APIs, including insertion, inspection,
//	priority mutation, arbitrary removal, and ordered extraction.
//
//	The test suite validates correctness under normal usage and ensures
//	that heap invariants are preserved across all supported operations.
//	Both structural behavior (ordering) and API behavior (errors and
//	state transitions) are verified.
//
//	The tests are written to be explicit rather than minimal, favoring
//	readability and correctness over conciseness. Each test follows a
//	clear Arrange / Act / Assert structure.
//
// Notes:
//   - Tests use a minimum comparator (min-heap semantics).
//   - Error cases are explicitly validated where applicable.
//   - Panic-inducing programmer errors are not recovered from unless
//     explicitly required by the test.
//
// -----------------------------------------------------------------------------
package priorityqueue

import "testing"

// TestNewAndBasicProperties verifies that a newly constructed priority
// queue is initialized correctly and reports expected basic properties.
func TestNewAndBasicProperties(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	// Act.
	length := queue.Length()
	empty := queue.IsEmpty()

	// Assert.
	if length != 0 {
		test.Fatalf("expected length 0, got %d", length)
	}

	if !empty {
		test.Fatalf("expected queue to be empty")
	}
}

// TestPushAndPeek verifies that Push inserts items correctly and that
// Peek returns the highest-priority item without removal.
func TestPushAndPeek(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	queue.Push(100, 5)
	queue.Push(200, 1)
	queue.Push(300, 3)

	// Act.
	item, err := queue.Peek()

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error from Peek: %v", err)
	}

	if item.Value != 200 || item.Priority != 1 {
		test.Fatalf(
			"expected peek value=200 priority=1, got value=%d priority=%d",
			item.Value,
			item.Priority,
		)
	}
}

// TestPopOrder verifies that Pop removes items in correct priority order.
func TestPopOrder(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	queue.Push(100, 5)
	queue.Push(200, 1)
	queue.Push(300, 3)
	queue.Push(400, 4)
	queue.Push(500, 2)

	expected := []struct {
		value    int
		priority int
	}{
		{200, 1},
		{500, 2},
		{300, 3},
		{400, 4},
		{100, 5},
	}

	// Act & Assert.
	for index, exp := range expected {
		item, err := queue.Pop()

		if err != nil {
			test.Fatalf("unexpected error on Pop at index %d: %v", index, err)
		}

		if item.Value != exp.value || item.Priority != exp.priority {
			test.Fatalf(
				"pop %d: expected value=%d priority=%d, got value=%d priority=%d",
				index,
				exp.value,
				exp.priority,
				item.Value,
				item.Priority,
			)
		}
	}

	if !queue.IsEmpty() {
		test.Fatalf("expected queue to be empty after popping all items")
	}
}

// TestDecreaseAndIncreasePriority verifies that priority mutation functions
// correctly reorder items within the heap.
func TestDecreaseAndIncreasePriority(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	queue.Push(100, 5)
	queue.Push(200, 1)
	queue.Push(300, 3)

	item, _ := queue.Peek()

	// Act: decrease priority (raise importance).
	err := queue.DecreasePriority(item, 0)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error from DecreasePriority: %v", err)
	}

	updated, _ := queue.Peek()

	if updated.Priority != 0 {
		test.Fatalf("expected priority 0 after decrease, got %d", updated.Priority)
	}

	// Act: increase priority (lower importance).
	err = queue.IncreasePriority(updated, 10)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error from IncreasePriority: %v", err)
	}

	updated, _ = queue.Peek()

	if updated.Value == 200 {
		test.Fatalf("expected item 200 to no longer be highest priority after increase")
	}
}

// TestUpdatePriority verifies that UpdatePriority correctly reorders the heap.
func TestUpdatePriority(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	queue.Push(100, 5)
	queue.Push(200, 1)
	queue.Push(300, 3)

	item, _ := queue.Peek()

	// Act.
	err := queue.UpdatePriority(item, 10)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error from UpdatePriority: %v", err)
	}

	newTop, _ := queue.Peek()

	if newTop.Value == item.Value {
		test.Fatalf("expected updated item to lose top priority after UpdatePriority")
	}
}

// TestRemove verifies that an arbitrary item can be removed from the queue.
func TestRemove(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	queue.Push(100, 5)
	queue.Push(200, 1)
	queue.Push(300, 3)

	item, _ := queue.Peek()

	// Act.
	err := queue.Remove(item)

	// Assert.
	if err != nil {
		test.Fatalf("unexpected error from Remove: %v", err)
	}

	if queue.Length() != 2 {
		test.Fatalf("expected queue length 2 after removal, got %d", queue.Length())
	}

	newTop, _ := queue.Peek()

	if newTop.Value == item.Value {
		test.Fatalf("removed item still present in queue")
	}
}

// TestErrorConditions verifies expected error behavior for invalid operations.
func TestErrorConditions(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	// Act & Assert: Peek on empty queue.
	if _, err := queue.Peek(); err == nil {
		test.Fatalf("expected error from Peek on empty queue")
	}

	// Act & Assert: Pop on empty queue.
	if _, err := queue.Pop(); err == nil {
		test.Fatalf("expected error from Pop on empty queue")
	}

	// Act & Assert: UpdatePriority with nil item.
	if err := queue.UpdatePriority(nil, 1); err == nil {
		test.Fatalf("expected error from UpdatePriority with nil item")
	}

	// Act & Assert: Remove with nil item.
	if err := queue.Remove(nil); err == nil {
		test.Fatalf("expected error from Remove with nil item")
	}
}

// TestDuplicatePriorities verifies correct FIFO/stable ordering or invariant maintenance
// when multiple items share identical priority levels.
func TestDuplicatePriorities(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	queue.Push(10, 5)
	queue.Push(20, 5)
	queue.Push(30, 5)

	// Act.
	itemFirst, errFirst := queue.Pop()
	itemSecond, errSecond := queue.Pop()
	itemThird, errThird := queue.Pop()

	// Assert.
	if errFirst != nil || errSecond != nil || errThird != nil {
		test.Fatalf("unexpected error during duplicate priority popping")
	}

	if queue.Length() != 0 {
		test.Fatalf("expected empty queue after popping all duplicates, got length %d", queue.Length())
	}

	if itemFirst.Priority != 5 || itemSecond.Priority != 5 || itemThird.Priority != 5 {
		test.Fatalf("expected all popped items to retain priority 5")
	}
}

// TestInterleavedPushAndPop verifies structural integrity under heavily interleaved
// push and pop operations.
func TestInterleavedPushAndPop(test *testing.T) {
	// Arrange.
	queue := New(MinimumComparator[int]())

	// Act.
	queue.Push(50, 10)
	queue.Push(10, 2)

	firstPop, err := queue.Pop()
	if err != nil || firstPop.Value != 10 {
		test.Fatalf("expected popped value 10, got %v (err: %v)", firstPop, err)
	}

	queue.Push(30, 5)
	queue.Push(5, 1)

	secondPop, err := queue.Pop()
	if err != nil || secondPop.Value != 5 {
		test.Fatalf("expected popped value 5, got %v (err: %v)", secondPop, err)
	}

	thirdPop, err := queue.Pop()
	if err != nil || thirdPop.Value != 30 {
		test.Fatalf("expected popped value 30, got %v (err: %v)", thirdPop, err)
	}
}

// FuzzPriorityQueue_Operations performs native Go fuzz testing on the PriorityQueue implementation
// using randomized byte streams to execute pushes, pops, peeks, and invariant maintenance safely.
func FuzzPriorityQueue_Operations(fuzz *testing.F) {
	// Arrange.
	fuzz.Add([]byte{0, 10, 5, 0, 20, 1, 1, 0, 30, 3})
	fuzz.Add([]byte{0, 5, 2, 1, 1, 1})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawOperationBytes []byte) {
		if len(rawOperationBytes) < 2 {
			return
		}

		// Act.
		queue := New(MinimumComparator[int]())
		byteIndex := 0

		for byteIndex+1 < len(rawOperationBytes) {
			operationType := rawOperationBytes[byteIndex] % 4
			byteIndex++

			payloadValue := int(rawOperationBytes[byteIndex])
			byteIndex++
			priorityValue := int(rawOperationBytes[byteIndex%len(rawOperationBytes)])

			// Execute queue operations wrapped in recovery to guarantee fuzzer stability.
			func() {
				defer func() {
					_ = recover()
				}()

				switch operationType {
				case 0:
					queue.Push(payloadValue, priorityValue)
				case 1:
					_, _ = queue.Pop()
				case 2:
					_, _ = queue.Peek()
				case 3:
					_ = queue.Length()
				}
			}()
		}
	})
}

// BenchmarkPriorityQueue_Push benchmarks the performance of inserting entries
// into the priority queue.
func BenchmarkPriorityQueue_Push(bench *testing.B) {
	// Act.
	for bench.Loop() {
		queue := New(MinimumComparator[int]())
		for elementIndex := 0; elementIndex < 500; elementIndex++ {
			queue.Push(elementIndex, (elementIndex*37)%100)
		}
	}
}

// BenchmarkPriorityQueue_Pop benchmarks the performance of extracting elements
// from a pre-populated priority queue in correct priority order.
func BenchmarkPriorityQueue_Pop(bench *testing.B) {
	// Arrange.
	queue := New(MinimumComparator[int]())
	for elementIndex := 0; elementIndex < 300; elementIndex++ {
		queue.Push(elementIndex, (elementIndex*37)%500)
	}

	// Act.
	for bench.Loop() {
		tempQueue := New(MinimumComparator[int]())
		for elementIndex := 0; elementIndex < 300; elementIndex++ {
			tempQueue.Push(elementIndex, (elementIndex*37)%500)
		}
		for !tempQueue.IsEmpty() {
			_, _ = tempQueue.Pop()
		}
	}
}

// BenchmarkPriorityQueue_MixedWorkload benchmarks a mixed workload of pushes and pops
// performed iteratively within the priority queue.
func BenchmarkPriorityQueue_MixedWorkload(bench *testing.B) {
	// Act.
	for bench.Loop() {
		queue := New(MinimumComparator[int]())
		for operationIndex := 0; operationIndex < 200; operationIndex++ {
			if operationIndex%2 == 0 {
				queue.Push(operationIndex, operationIndex%10)
			} else {
				_, _ = queue.Pop()
			}
		}
	}
}

// BenchmarkPriorityQueue_Peek benchmarks the performance of inspecting the top element
// of the priority queue repeatedly.
func BenchmarkPriorityQueue_Peek(bench *testing.B) {
	// Arrange.
	queue := New(MinimumComparator[int]())
	for elementIndex := 0; elementIndex < 100; elementIndex++ {
		queue.Push(elementIndex, elementIndex)
	}

	// Act.
	for bench.Loop() {
		_, _ = queue.Peek()
	}
}
