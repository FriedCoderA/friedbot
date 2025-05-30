package xring

import "sync"

// Ring 泛型环形队列结构体
type Ring[T any] struct {
	mu       sync.RWMutex // 互斥锁保证并发安全
	data     []T          // 底层数据存储
	head     int          // 队列头指针
	tail     int          // 队列尾指针
	size     int          // 当前元素数量
	capacity int          // 队列容量
}

// NewRing 创建新环形队列
func NewRing[T any](capacity int) *Ring[T] {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	return &Ring[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

// Push 入队操作
func (q *Ring[T]) Push(item T) {
	_, _ = q.PushWithPopped(item)
}

// PushWithPopped 入队操作，返回被挤出的元素及其存在标志
func (q *Ring[T]) PushWithPopped(item T) (popped T, exists bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == q.capacity {
		// 队列已满时挤出头部元素
		popped = q.data[q.head]
		exists = true

		// 移动头指针并插入新元素
		q.head = (q.head + 1) % q.capacity
		q.data[q.tail] = item
		q.tail = (q.tail + 1) % q.capacity
	} else {
		// 队列未满时直接插入
		q.data[q.tail] = item
		q.tail = (q.tail + 1) % q.capacity
		q.size++
	}
	return
}

// Pop 出队操作，返回元素及其存在标志
func (q *Ring[T]) Pop() (item T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return
	}

	item = q.data[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return item, true
}

// Size 获取当前队列元素数量
func (q *Ring[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.size
}

// Capacity 获取队列容量
func (q *Ring[T]) Capacity() int {
	return q.capacity
}

// Reset 重置队列
func (q *Ring[T]) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.head = 0
	q.tail = 0
	q.size = 0
	var zero T
	for i := range q.data {
		q.data[i] = zero
	}
}

// All 返回当前队列所有元素的副本（按入队顺序）
func (q *Ring[T]) All() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.size == 0 {
		return nil
	}

	result := make([]T, 0, q.size)

	// 计算第一部分长度
	firstPart := q.capacity - q.head
	if firstPart > q.size {
		firstPart = q.size
	}

	// 添加两部分元素
	result = append(result, q.data[q.head:q.head+firstPart]...)
	if remaining := q.size - firstPart; remaining > 0 {
		result = append(result, q.data[:remaining]...)
	}

	return result
}
