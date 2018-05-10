package queues

import "unsafe"

type ChanQ chan unsafe.Pointer

func NewChanQ(sz int) ChanQ {
	return make(chan unsafe.Pointer, sz)
}

func (q ChanQ) Pop(e *unsafe.Pointer) bool {
	select {
	case *e = <-q:
		return true
	default:
		return false
	}
}

func (q ChanQ) Push(e unsafe.Pointer) bool {
	select {
	case q <- e:
		return true
	default:
		return false
	}
}
