package taskpool

import (
	"context"
	"sync"
	"sync/atomic"
)

var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

type task struct {
	ctx  context.Context
	f    func()
	next *task
}

func newTask() interface{} {
	return &task{}
}

type pool struct {
	// 名字，用于打点和日志
	name string
	// pool的容量，实际运行的最大的goroutines数量
	cap int32
	// 任务链表
	taskHead *task
	taskTail *task

	taskLock sync.Mutex
	// 记录正在运行的 worker 数量
	taskCount int32
}

func NewPool(name string, cap int32) *pool {
	return &pool{
		name: name,
		cap:  cap,
	}
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f

	p.taskLock.Lock()

	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}

	p.taskLock.Unlock()

	atomic.AddInt32(&p.taskCount, 1)

	w := workerPool.Get().(*worker)
	w.pool = p
	w.run()
}
