package main

import (
	"strconv"
	"testing"
)

func BenchmarkAllocs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x []byte
		for j := 0; j < 1000000; j += 10002 {
			x = append(x, strconv.Itoa(j)...)
			//x = strconv.AppendInt(x, int64(j), 10)
		}
	}
}
