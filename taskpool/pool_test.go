package taskpool

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p := NewPool("test", 10)
	var n int32
	var now = time.Now()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		var s = i
		p.CtxGo(nil, func() {
			defer wg.Done()
			log.Printf("i=%v", s)
			time.Sleep(300 * time.Millisecond)
			atomic.AddInt32(&n, 1)
		})
	}
	wg.Wait()

	log.Printf("n=%v, timeStamp=%d", n, time.Since(now).Milliseconds())

}
