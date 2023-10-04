package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateQueue(t *testing.T) {
	queue := BuildQueue[int](2)

	assert.Equal(t, 2, queue.queueMaxSize, "The max size of the queue should be 2.")
	assert.Equal(t, 0, queue.Size(), "The len of the queue should be 0.")
}

func TestAddToQueue(t *testing.T) {
	maxSize := 2
	queue := BuildQueue[int](maxSize)

	queue.Add(1)
	queue.Add(2)
	queue.Add(3)
	queue.Add(1)

	assert.Equal(t, maxSize, queue.Size(), "The max size of the queue should be %d.", maxSize)
}

func TestRemoveFromQueue(t *testing.T) {
	maxSize := 3
	queue := BuildQueue[int](maxSize)
	queue.Add(5)
	queue.Add(6)
	queue.Add(1)
	v, _ := queue.Next()

	assert.Equal(t, 2, queue.Size(), "The max size of the queue should be 2.")
	assert.Equal(t, 5, v, "The top value should be 5.")

	v, _ = queue.Next()
	assert.Equal(t, 6, v, "The top value should be 6.")
	assert.Equal(t, 1, queue.Size(), "The max size of the queue should be 1.")

	v, _ = queue.Next()
	assert.Equal(t, 1, v, "The top value should be 1.")
	assert.Equal(t, 0, queue.Size(), "The max size of the queue should be 0.")

	queue.Add(6)
	v, _ = queue.Next()
	assert.Equal(t, 6, v, "The top value should be 6.")
	assert.Equal(t, 0, queue.Size(), "The max size of the queue should be 0.")
}
