package queues

import (
	"sync/atomic"
)

type SPSCAccess struct {
	_         [cpuL1Size]byte
	head      int64
	_         [cpuL1Size]byte
	headCache int64
	_         [cpuL1Size]byte
	tail      int64
	_         [cpuL1Size]byte
	tailCache int64
	_         [cpuL1Size]byte
	size      int64
	mask      int64
}

func NewSPSCAccess(sz int) *SPSCAccess {
	if sz&(sz-1) != 0 {
		panic("sz must be power of 2")
	}
	return &SPSCAccess{
		size: int64(sz),
		mask: int64(sz - 1),
	}
}

func (q *SPSCAccess) PreparePop() int {
	head := atomic.LoadInt64(&q.head)
	if head >= q.tailCache {
		q.tailCache = atomic.LoadInt64(&q.tail)
		if head >= q.tailCache {
			return -1
		}
	}
	return int(head & q.mask)
}

func (q *SPSCAccess) FinishPop() {
	atomic.StoreInt64(&q.head, q.head+1)
}

func (q *SPSCAccess) PreparePush() int {
	tail := atomic.LoadInt64(&q.tail)
	if q.headCache <= tail-q.size {
		q.headCache = atomic.LoadInt64(&q.head)
		if q.headCache <= tail-q.size {
			return -1
		}
	}
	return int(tail & q.mask)
}

func (q *SPSCAccess) FinishPush() {
	atomic.StoreInt64(&q.tail, q.tail+1)
}
