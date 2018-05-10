package queues

import (
	"sync/atomic"
	"unsafe"
)

const cpuL1Size = 64

// apply the following optimisations:
// - avoid false sharing
type SPSC_3 struct {
	_    [cpuL1Size]byte
	head int64
	_    [cpuL1Size]byte
	tail int64
	_    [cpuL1Size]byte
	mask int64
	data []unsafe.Pointer
}

func NewSPSC_3(sz int) *SPSC_3 {
	c := nextPowerOf2(sz)
	return &SPSC_3{
		data: make([]unsafe.Pointer, c, c),
		mask: c - 1,
	}
}

func (q *SPSC_3) Pop(e *unsafe.Pointer) bool {
	head := atomic.LoadInt64(&q.head)
	tail := atomic.LoadInt64(&q.tail)
	if head >= tail {
		return false
	}
	i := head & q.mask
	*e = q.data[i]
	q.data[i] = nil
	atomic.StoreInt64(&q.head, head+1)
	return true
}

func (q *SPSC_3) Push(e unsafe.Pointer) bool {
	tail := atomic.LoadInt64(&q.tail)
	head := atomic.LoadInt64(&q.head)
	sz := int64(len(q.data))
	if head <= tail-sz {
		return false
	}
	q.data[tail&q.mask] = e
	atomic.StoreInt64(&q.tail, tail+1)
	return true
}
