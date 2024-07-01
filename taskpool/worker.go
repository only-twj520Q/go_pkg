package taskpool

import (
	"sync"
	"sync/atomic"
)

var workerPool sync.Pool

type worker struct {
	pool *pool
}

func init() {
	workerPool.New = newWorker
}

func newWorker() interface{} {
	return &worker{}
}

func (w *worker) run() {
	go func() {
		for {
			var t *task

			// 加锁，获取任务
			w.pool.taskLock.Lock()

			// 取出队首的task
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}

			if t == nil {
				// 如果没有任务，则释放锁，退出
				w.close()
				w.pool.taskLock.Unlock()
				return
			}

			w.pool.taskLock.Unlock()

			func() {
				defer func() {
					if e := recover(); e != nil {
					}
				}()
				t.f()
			}()

			workerPool.Put(w)
		}
	}()
}

func (w *worker) close() {
	w.pool.decWorkerCount()
}
