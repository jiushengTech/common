// Package queueutil 提供了一个基于 Dariusz Górecki 建议的快速环形缓冲队列。
// 使用这个包而不是其他更简单的队列实现（如切片+追加或链表）可以带来
// 实质性的内存和时间优势，并减少垃圾回收暂停的次数。
// 该队列的实现是线程安全的，使用读写锁来保护并发访问。
package queueutil

import (
	"sync"
)

// minQueueLen 是队列可能具有的最小容量。
// 必须是2的幂，以进行位掩码操作：x % n == x & (n - 1)。
const minQueueLen = 16

// Queue 表示队列数据结构的单个实例。
type Queue[V any] struct {
	buf               []*V
	head, tail, count int
	maxSize           int
	mutex             sync.RWMutex // 添加读写锁
}

// NewWithMaxSize New 构造并返回一个新的 Queue，可以指定最大容量。
func NewWithMaxSize[V any](maxSize int) *Queue[V] {
	if maxSize <= 0 {
		panic("queue: maxSize must be greater than 0")
	}
	return &Queue[V]{
		buf:     make([]*V, minQueueLen),
		maxSize: maxSize,
	}
}

// Add 将元素放在队列的末尾。
func (q *Queue[V]) Add(elem V) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// 如果设置了最大容量，且队列已满，则移除最旧的元素
	if q.maxSize > 0 && q.count == q.maxSize {
		_ = q.removeUnsafe()
	}
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = &elem
	// 位掩码操作
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// SetMaxSize 设置队列的最大容量。
func (q *Queue[V]) SetMaxSize(maxSize int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if maxSize <= 0 {
		panic("queue: maxSize must be greater than 0")
	}
	q.maxSize = maxSize
}

// Resize 将队列调整为恰好容纳两倍其当前内容
// 这可能导致缩小，如果队列不到一半满。
func (q *Queue[V]) resize() {
	newBuf := make([]*V, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}
	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// New 构造并返回一个新的 Queue。
func New[V any]() *Queue[V] {
	return &Queue[V]{
		buf: make([]*V, minQueueLen),
	}
}

// Length 返回当前存储在队列中的元素数量。
func (q *Queue[V]) Length() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.count
}

// Peek 返回队列头部的元素。如果队列为空，此调用将引发 panic。
func (q *Queue[V]) Peek() V {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return *(q.buf[q.head])
}

// Get 返回队列中索引为 i 的元素。如果索引无效，此调用将引发 panic。
// 此方法接受正索引和负索引值。索引 0 指的是第一个元素，索引 -1 指的是最后一个元素。
func (q *Queue[V]) Get(i int) V {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	// 如果索引为负数，则转换为正索引。
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}
	// 位掩码操作
	return *(q.buf[(q.head+i)&(len(q.buf)-1)])
}

// Remove 从队列的前端删除并返回元素。如果队列为空，此调用将引发 panic。
func (q *Queue[V]) Remove() V {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.removeUnsafe()
}

// removeUnsafe 是 Remove 的未加锁版本，用于内部调用
func (q *Queue[V]) removeUnsafe() V {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// 位掩码操作
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// 如果缓冲区 1/4 满，则缩小大小。
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return *ret
}
