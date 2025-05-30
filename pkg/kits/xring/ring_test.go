package xring

import (
	"testing"
)

func TestRing(t *testing.T) {
	queue := NewRing[int](3)
	t.Run("test", func(t *testing.T) {
		var popped int
		var ok bool
		t.Log(queue.All())
		queue.Push(1)
		t.Log(queue.All())
		queue.Push(2)
		t.Log(queue.All())
		queue.Push(3)
		t.Log(queue.All())
		queue.Push(4)
		t.Log(queue.All())
		popped, ok = queue.PushWithPopped(5)
		t.Log(queue.All(), ok, popped)
		popped, ok = queue.PushWithPopped(6)
		t.Log(queue.All(), ok, popped)
		queue.Push(7)
		t.Log(queue.All())
		popped, ok = queue.Pop()
		t.Log(queue.All(), ok, popped)
		popped, ok = queue.Pop()
		t.Log(queue.All(), ok, popped)
		popped, ok = queue.Pop()
		t.Log(queue.All(), ok, popped)
		popped, ok = queue.Pop()
		t.Log(queue.All(), ok, popped)
	})
}
