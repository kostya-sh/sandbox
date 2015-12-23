package main

import (
	"bytes"
	"runtime"
	"testing"
	"time"
)

func makeLongString() string {
	const N = 10000
	b := bytes.NewBuffer(make([]byte, N*3))
	for i := 0; i < N; i++ {
		b.Write([]byte{'x', 'y', 'z'})
	}
	return b.String()
}

func TestGC(t *testing.T) {
	var ss []string
	for i := 0; i < 10000; i++ {
		s := makeLongString()
		ss = append(ss, s[1:2])
	}
	for i := 0; i < 3; i++ {
		runtime.GC()
		time.Sleep(1 * time.Second)
	}
	for _, s := range ss {
		println(s)
	}
	for i := 0; i < 3; i++ {
		runtime.GC()
		time.Sleep(1 * time.Second)
	}
}
