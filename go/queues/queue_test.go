package queues

import (
	"runtime"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

func expensiveCalc() float64 {
	s := 0.0
	for i := 1; i < 10; i++ {
		s += (s / float64(i))
	}
	return s
}

func testQ(t *testing.T, q Q) {
	runtime.GC()

	const N = 10000000
	const lockThreads = false
	const slowConsumer = false
	const slowPublisher = false

	pubRetries := make(chan int, 1)
	conRetries := make(chan int, 1)
	start := time.Now()

	go func() {
		if lockThreads {
			runtime.LockOSThread()
		}
		retries := 0
		for i := 0; i < N; i++ {
			i := i
			for !q.Push(unsafe.Pointer(&i)) {
				retries++
				runtime.Gosched()
			}
			if slowPublisher {
				expensiveCalc()
			}
		}
		pubRetries <- retries
	}()

	go func() {
		if lockThreads {
			runtime.LockOSThread()
		}
		retries := 0
		for i := 0; i < N; i++ {
			var pe unsafe.Pointer
			for !q.Pop(&pe) {
				retries++
				runtime.Gosched()
			}
			e := *(*int)(pe)
			if e != i {
				t.Errorf("Want %d, got %d", i, e)
			}
			if slowConsumer {
				expensiveCalc()
			}

		}
		conRetries <- retries
	}()

	t.Logf("messages: %d", N)
	t.Logf("publisher retries: %d", <-pubRetries)
	t.Logf("consumer retries: %d", <-conRetries)

	elapsed := time.Since(start)
	t.Logf("elapsed: %s", elapsed)
	t.Logf("throughput: %f million msgs/sec", N/float64(elapsed)*float64(time.Second)/1000000)
}

func TestChanQ(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testQ(t, NewChanQ(sz)) })
	}
}

func TestSPSC_1(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testQ(t, NewSPSC_1(sz)) })
	}
}

func TestSPSC_2(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testQ(t, NewSPSC_2(sz)) })
	}
}

func TestSPSC_3(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testQ(t, NewSPSC_3(sz)) })
	}
}

func TestSPSC_4(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testQ(t, NewSPSC_4(sz)) })
	}
}

func TestSPSC_5(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testQ(t, NewSPSC_5(sz)) })
	}
}

func TestNextPowerOf2(t *testing.T) {
	if got := nextPowerOf2(1); got != 1 {
		t.Errorf("want 1, got %d", got)
	}
	if got := nextPowerOf2(2); got != 2 {
		t.Errorf("want 2, got %d", got)
	}
	if got := nextPowerOf2(17); got != 32 {
		t.Errorf("want 32, got %d", got)
	}
}
