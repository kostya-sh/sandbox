package queues

import "unsafe"

type Q interface {
	Push(e unsafe.Pointer) bool
	Pop(e *unsafe.Pointer) bool
}
