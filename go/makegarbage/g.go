package main

import (
	"math/rand"
	"os"
	"strconv"
)

func main() {
	nnnn, _ := strconv.Atoi(os.Args[1])
	s, _ := strconv.Atoi(os.Args[2])
	for i := 0; i < nnnn; i++ {
		_ = make([]int32, rand.Intn(s)+1)
	}
}
