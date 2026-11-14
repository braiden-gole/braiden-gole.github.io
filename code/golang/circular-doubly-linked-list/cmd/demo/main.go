// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file contains a runnable demonstration of the
//	circulardoublylinkedlist package. The program exercises core circular
//	doubly linked list operations, including insertion at both ends,
//	removal from both ends, predicate-based search, state inspection,
//	and defensive nil-receiver handling, while printing list state and
//	size after each operation.
//
//	The demo is intended to provide a clear, step-by-step illustration of
//	circular doubly linked list behavior and to serve as a manual
//	verification tool for correctness, invariant preservation, and
//	edge-case handling.
//
// Notes:
//   - This file is not part of the circulardoublylinkedlist library API.
//   - Output is intentionally verbose to make state transitions explicit.
//   - Traversal and inspection are bounded to prevent infinite loops
//     inherent to circular data structures.
//   - The demo favors clarity and correctness over conciseness.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/braiden-gole/circular-doubly-linked-list/circulardoublylinkedlist"
)

func main() {
	list := circulardoublylinkedlist.New()

	fmt.Println("Initial state:")
	fmt.Println("\tIsEmpty:", list.IsEmpty())
	fmt.Println("\tSize:", list.Size())
	fmt.Println("\tHead:", list.Head())
	fmt.Println("\tTail:", list.Tail())
	fmt.Println()

	fmt.Println("PushFront:")

	list.PushFront(10)
	list.PushFront(20)
	list.PushFront(30)

	fmt.Println("\tHead:", list.Head().Value)
	fmt.Println("\tTail:", list.Tail().Value)
	fmt.Println("\tSize:", list.Size())
	fmt.Println()

	fmt.Println("PushBack:")

	list.PushBack(40)
	list.PushBack(50)

	fmt.Println("\tHead:", list.Head().Value)
	fmt.Println("\tTail:", list.Tail().Value)
	fmt.Println("\tSize:", list.Size())
	fmt.Println()

	fmt.Print("Find: ")

	found := list.Find(func(value interface{}) bool {
		return value == 20
	})

	if found != nil {
		fmt.Println(found.Value)
	} else {
		fmt.Println("Value not found")
	}

	fmt.Println()

	fmt.Println("PopFront:")

	value, ok := list.PopFront()

	fmt.Println("\tPopped:", value, "Success:", ok)

	fmt.Println("\tNew Head:", list.Head().Value)
	fmt.Println("\tSize:", list.Size())
	fmt.Println()

	fmt.Println("PopBack:")

	value, ok = list.PopBack()

	fmt.Println("\tPopped:", value, "Success:", ok)

	fmt.Println("\tNew Tail:", list.Tail().Value)
	fmt.Println("\tSize:", list.Size())
	fmt.Println()

	fmt.Println("Clear:")

	list.Clear()

	fmt.Println("\tIsEmpty:", list.IsEmpty())
	fmt.Println("\tSize:", list.Size())
	fmt.Println("\tHead:", list.Head())
	fmt.Println("\tTail:", list.Tail())
	fmt.Println()

	fmt.Println("Nil receiver behavior:")

	var nilList *circulardoublylinkedlist.List

	fmt.Println("\tNil Size:", nilList.Size())
	fmt.Println("\tNil IsEmpty:", nilList.IsEmpty())
	fmt.Println("\tNil Head:", nilList.Head())
	fmt.Println("\tNil Tail:", nilList.Tail())

	nilList.PushFront(1)
	nilList.PushBack(2)

	fmt.Println()

	value, ok = nilList.PopFront()

	fmt.Println("Nil PopFront:", value, ok)

	value, ok = nilList.PopBack()

	fmt.Println("Nil PopBack:", value, ok)
}
