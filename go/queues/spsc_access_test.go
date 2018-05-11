package queues

import (
	"runtime"
	"strconv"
	"testing"
	"time"
)

func testSPSCAccess(t *testing.T, sz int) {
	runtime.GC()

	sz = int(nextPowerOf2(sz))
	data := make([]int, sz, sz)
	q := NewSPSCAccess(sz)

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
			var p int
			for p = q.PreparePush(); p < 0; p = q.PreparePush() {
				retries++
				runtime.Gosched()
			}
			data[p] = i
			q.FinishPush()
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
			var p int
			for p = q.PreparePop(); p < 0; p = q.PreparePop() {
				retries++
				runtime.Gosched()
			}
			e := data[p]
			q.FinishPop()
			if e != i {
				t.Errorf("Want %d, got %d", i, e)
			}
			if slowConsumer {
				expensiveCalc()
			}

		}
		conRetries <- retries
	}()

	t.Logf("buffer size: %d", sz)
	t.Logf("messages: %d", N)
	t.Logf("publisher retries: %d", <-pubRetries)
	t.Logf("consumer retries: %d", <-conRetries)

	elapsed := time.Since(start)
	t.Logf("elapsed: %s", elapsed)
	t.Logf("throughput: %f million msgs/sec", N/float64(elapsed)*float64(time.Second)/1000000)
}

func TestSPSCAccess(t *testing.T) {
	for _, sz := range [...]int{1, 10, 100, 1000, 10000000} {
		t.Run(strconv.Itoa(sz), func(t *testing.T) { testSPSCAccess(t, sz) })
	}
}
