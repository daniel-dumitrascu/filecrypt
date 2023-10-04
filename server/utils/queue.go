package utils

import (
	"errors"
	"fmt"
)

type Queue[T any] struct {
	queueData    []T
	queueMaxSize int
}

func BuildQueue[T any](size int) *Queue[T] {
	q := new(Queue[T])
	q.queueData = make([]T, 0, size)
	q.queueMaxSize = size
	return q
}

// Adding at the bottom
func (q *Queue[T]) Add(new_data T) {
	if q.queueMaxSize == len(q.queueData) {
		fmt.Printf("Queue is maxed out, cannot add anymore items\n")
		return
	}
	q.queueData = append(q.queueData, new_data)
}

// Getting the top and also removing it from queue
func (q *Queue[T]) Next() (T, error) {
	if len(q.queueData) == 0 {
		var def T
		return def, errors.New("[ERROR] cannot return the top element from queue because the queue is empty")
	}
	v := q.queueData[0]
	q.queueData = q.queueData[1:]
	return v, nil
}

// Get the size of the queue
func (q *Queue[T]) Size() int {
	return len(q.queueData)
}
