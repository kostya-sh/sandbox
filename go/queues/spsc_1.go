package queues

import (
	"sync/atomic"
	"unsafe"
)

type SPSC_1 struct {
	head int64
	tail int64
	data []unsafe.Pointer
}

func NewSPSC_1(sz int) *SPSC_1 {
	return &SPSC_1{
		data: make([]unsafe.Pointer, sz, sz),
	}
}

func (q *SPSC_1) Pop(e *unsafe.Pointer) bool {
	head := atomic.LoadInt64(&q.head)
	tail := atomic.LoadInt64(&q.tail)
	if head >= tail {
		return false
	}
	i := head % int64(len(q.data))
	*e = q.data[i]
	q.data[i] = nil
	atomic.StoreInt64(&q.head, head+1)
	return true
}

func (q *SPSC_1) Push(e unsafe.Pointer) bool {
	tail := atomic.LoadInt64(&q.tail)
	head := atomic.LoadInt64(&q.head)
	sz := int64(len(q.data))
	if head <= tail-sz {
		return false
	}
	q.data[tail%sz] = e
	atomic.StoreInt64(&q.tail, tail+1)
	return true
}
