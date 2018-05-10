package queues

import (
	"math/bits"
	"sync/atomic"
	"unsafe"
)

// apply the following optimisations:
// - make capacity power of 2 and use & instead of /
type SPSC_2 struct {
	head int64
	tail int64
	mask int64
	data []unsafe.Pointer
}

func NewSPSC_2(sz int) *SPSC_2 {
	c := nextPowerOf2(sz)
	return &SPSC_2{
		data: make([]unsafe.Pointer, c, c),
		mask: c - 1,
	}
}

func (q *SPSC_2) Pop(e *unsafe.Pointer) bool {
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

func (q *SPSC_2) Push(e unsafe.Pointer) bool {
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

// nextPowerOf2 ...
func nextPowerOf2(x int) int64 {
	return 1 << uint(64-bits.LeadingZeros64(uint64(x-1)))
}
