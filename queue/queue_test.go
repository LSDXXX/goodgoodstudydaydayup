package queue

import (
	"sync"
	"sync/atomic"
	"testing"
)

func testSync(t *testing.T, q ThreadSafeQueue) {
	tmp := []int{}
	for i := 0; i < 1000; i++ {
		tmp = append(tmp, i)
	}
	for _, n := range tmp {
		q.Enqueue(n)
	}
	for _, n := range tmp {
		if n != q.Dequeue() {
			t.Fail()
		}
	}
}

func TestLockFreeSync(t *testing.T) {
	q := NewLockfreeQueue()
	testSync(t, q)
}

func TestLockedSync(t *testing.T) {
	q := NewLockedQueue()
	testSync(t, q)
}

func testEnqueueAsync(t *testing.T, q ThreadSafeQueue, nNum, nThread int) {
	var wg sync.WaitGroup
	for i := 0; i < nThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < nNum; i++ {
				q.Enqueue(i)
			}
		}()
	}
	wg.Wait()
	count := 0
	for {
		v := q.Dequeue()
		if v == nil {
			break
		}
		count++
	}
	if count != nNum*nThread {
		t.Fail()
	}
}

func TestLockFreeEnqueueAsync(t *testing.T) {
	q := NewLockfreeQueue()
	nNum := 1000
	nThread := 8
	testEnqueueAsync(t, q, nNum, nThread)
}
func TestLockedEnqueueAsync(t *testing.T) {
	q := NewLockedQueue()
	nNum := 1000
	nThread := 8
	testEnqueueAsync(t, q, nNum, nThread)
}

func testEnqueueDequeueAsync(t *testing.T, q ThreadSafeQueue, nNum, nThread int) {
	var wg sync.WaitGroup

	for i := 0; i < nThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < nNum; i++ {
				q.Enqueue(i)
			}
		}()
	}

	var dequeueCount int32

	for i := 0; i < nThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < nNum; i++ {
				if q.Dequeue() != nil {
					atomic.AddInt32(&dequeueCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	var count int32
	for {
		v := q.Dequeue()
		if v == nil {
			break
		}
		count++
	}
	if count+dequeueCount != int32(nNum)*int32(nThread) {
		t.Fail()
	}
}

func TestLockFreeEnqueueDequeueAsync(t *testing.T) {
	q := NewLockfreeQueue()
	nNum := 100000
	nThread := 8
	testEnqueueDequeueAsync(t, q, nNum, nThread)
}

func TestLockedEnqueueDequeueAsync(t *testing.T) {
	q := NewLockedQueue()
	nNum := 100000
	nThread := 8
	testEnqueueDequeueAsync(t, q, nNum, nThread)
}

func benchmarkEnqueue(b *testing.B, q ThreadSafeQueue) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue(1)
		}
	})
}

func BenchmarkLockfreeEnqueue(b *testing.B) {
	q := NewLockfreeQueue()
	benchmarkEnqueue(b, q)
}

func BenchmarkLockedEnqueue(b *testing.B) {
	q := NewLockedQueue()
	benchmarkEnqueue(b, q)
}

func benchmarkEnqueueDequeue(b *testing.B, q ThreadSafeQueue) {
	var n int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if atomic.AddInt32(&n, 1)%2 == 0 {
				q.Dequeue()
			} else {
				q.Enqueue(1)
			}
		}
	})
}

func BenchmarkLockfreeEnqueueDequeue(b *testing.B) {
	q := NewLockfreeQueue()
	benchmarkEnqueueDequeue(b, q)
}

func BenchmarkLockedEnqueueDequeue(b *testing.B) {
	q := NewLockedQueue()
	benchmarkEnqueueDequeue(b, q)
}
