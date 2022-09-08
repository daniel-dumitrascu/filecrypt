package utils

import (
	"errors"
)

type Queue[T any] struct {
	queue_data []T
}

func BuildQueue[T any]() *Queue[T] {
	q := new(Queue[T])
	q.queue_data = make([]T, 0)
	return q
}

// Adding on the bottom
func (q *Queue[T]) Add(new_data T) {
	q.queue_data = append(q.queue_data, new_data)
}

// Getting from top
func (q *Queue[T]) Get() (T, error) {
	if len(q.queue_data) == 0 {
		var def T
		return def, errors.New("cannot return the top element from queue because the queue is empty")
	}
	return q.queue_data[0], nil
}

// Removing from top
func (q *Queue[T]) Remove() {
	q.queue_data = q.queue_data[1:]
}

func (q *Queue[T]) Size() int {
	return len(q.queue_data)
}
