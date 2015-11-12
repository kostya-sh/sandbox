package main

import (
	"fmt"

	"golang.org/x/sys/unix"

	"github.com/derekparker/delve/proc"
)

func main() {
	p := proc.New(10)
	var s *unix.WaitStatus
	s = p.Status()
	fmt.Println(s.Exited())
}
