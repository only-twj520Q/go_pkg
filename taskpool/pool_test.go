package taskpool

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

const benchmarkTimes = 10000

func DoCopyStack(a, b int) int {
	if b < 100 {
		return DoCopyStack(0, b+1)
	}
	return 0
}

func testFunc() {
	DoCopyStack(0, 0)
}

func TestPool(t *testing.T) {
	p := NewPool("test", 100, NewConfig(2))
	var n int32

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		var s = i
		p.CtxGo(nil, func() {
			defer wg.Done()
			log.Printf("i=%v", s)
			atomic.AddInt32(&n, 1)
		})
	}

	wg.Wait()

	if n != 1000 {
		t.Error(n)
	}
}

func BenchmarkPool(b *testing.B) {
	config := NewDefaultConfig()
	config.ScaleThreshold = 1
	p := NewPool("benchmark", int32(runtime.GOMAXPROCS(0)), config)
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			p.Go(func() {
				testFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
}

func BenchmarkGo(b *testing.B) {
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			go func() {
				testFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
