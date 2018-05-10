package queues

import (
	"sync/atomic"
	"unsafe"
)

// apply the following optimisations:
// - use lazySet
type SPSC_5 struct {
	_         [cpuL1Size]byte
	head      int64
	_         [cpuL1Size]byte
	headCache int64
	_         [cpuL1Size]byte
	tail      int64
	_         [cpuL1Size]byte
	tailCache int64
	_         [cpuL1Size]byte
	mask      int64
	data      []unsafe.Pointer
}

func NewSPSC_5(sz int) *SPSC_5 {
	c := nextPowerOf2(sz)
	return &SPSC_5{
		data: make([]unsafe.Pointer, c, c),
		mask: c - 1,
	}
}

func (q *SPSC_5) Pop(e *unsafe.Pointer) bool {
	head := atomic.LoadInt64(&q.head)
	if head >= q.tailCache {
		q.tailCache = atomic.LoadInt64(&q.tail)
		if head >= q.tailCache {
			return false
		}
	}
	i := head & q.mask
	*e = q.data[i]
	q.data[i] = nil
	q.head = head + 1 // atomic.StoreInt64(&q.head, head+1)
	return true
}

func (q *SPSC_5) Push(e unsafe.Pointer) bool {
	tail := atomic.LoadInt64(&q.tail)
	sz := int64(len(q.data))
	if q.headCache <= tail-sz {
		q.headCache = atomic.LoadInt64(&q.head)
		if q.headCache <= tail-sz {
			return false
		}
	}
	q.data[tail&q.mask] = e
	q.tail = tail + 1 // atomic.StoreInt64(&q.tail, tail+1)
	return true
}
