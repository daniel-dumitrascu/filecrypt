package utils

import (
	"sync"
)

type SafeQueue[T any] struct {
	queue *Queue[T]
	mu    sync.Mutex
}

func BuildSafeQueue[T any](size int) *SafeQueue[T] {
	sq := new(SafeQueue[T])
	sq.queue = BuildQueue[T](size)
	return sq
}

// Adding at the bottom
func (sq *SafeQueue[T]) Add(new_data T) {
	sq.mu.Lock()
	defer sq.mu.Unlock()
	sq.queue.Add(new_data)
}

// Getting the top and also removing it from queue
func (sq *SafeQueue[T]) Next() (T, error) {
	sq.mu.Lock()
	defer sq.mu.Unlock()
	return sq.queue.Next()
}

// Get the size of the queue
func (sq *SafeQueue[T]) IsEmpty() bool {
	sq.mu.Lock()
	defer sq.mu.Unlock()
	return sq.queue.Size() == 0
}
