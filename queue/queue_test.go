package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := New(1)
	q.Push(1)
	result := q.Pop()
	t.Log(result)
	q.Push(2)
	q.Push(3)
	result = q.Pop()
	t.Log(result)
	result = q.Pop()
	t.Log(result)
}
