package main

import "fmt"

// Queue is a generic FIFO queue
type Queue[T any] struct {
	items []T
}

// Create new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

// Add item to end of queue
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Remove and return first item
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// View first item without removing
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

func main() {
	// Create queue
	q := NewQueue[int]()

	// Enqueue some items
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(6)

	// Peek at first item
	if first, ok := q.Peek(); ok {
		fmt.Println("Peek:", first) // "first"
	}

	// Process all items
	fmt.Println("Processing queue:")
	for {
		item, ok := q.Dequeue()
		if !ok {
			break
		}
		fmt.Println(item)
	}
	q.Enqueue(9)
	value, exists := q.Peek()
	fmt.Printf("Peek value: %d, Exists: %t\n", value, exists)
	if _, ok := q.Dequeue(); !ok {
		fmt.Println("Queue is empty")
	}
}