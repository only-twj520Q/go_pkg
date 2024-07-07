package taskpool

import (
	"context"
	"sync"
	"sync/atomic"
)

var globalPool sync.Pool

func init() {
	globalPool.New = newTask
}

type task struct {
	ctx  context.Context
	f    func()
	next *task
}

func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	globalPool.Put(t)
}

func newTask() interface{} {
	return &task{}
}

type pool struct {
	// 名字，用于打点和日志
	name string
	// pool的容量，实际运行的最大的goroutines数量
	cap int32
	// 配置信息
	config *Config
	// 任务链表
	taskHead *task
	taskTail *task

	taskLock sync.Mutex
	// 全部任务数量
	taskCount int32

	// 正在运行的 worker 数量
	workerCount int32
}

type Config struct {
	// 当任务数量超过这个值时，会创建新的worker
	ScaleThreshold int32
}

const (
	defaultScalaThreshold = 1
)

func NewConfig(scaleThreshold int32) *Config {
	c := &Config{
		ScaleThreshold: scaleThreshold,
	}
	return c
}

func NewDefaultConfig() *Config {
	c := &Config{
		ScaleThreshold: defaultScalaThreshold,
	}
	return c
}

func NewPool(name string, cap int32, config *Config) *pool {
	return &pool{
		name:   name,
		cap:    cap,
		config: config,
	}
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := globalPool.Get().(*task)
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

	if atomic.LoadInt32(&p.taskCount) < p.config.ScaleThreshold && p.WorkerCount() != 0 {
		return
	}

	if p.WorkerCount() >= atomic.LoadInt32(&p.cap) && p.WorkerCount() != 0 {
		return
	}

	p.incWorkerCount()
	w := workerPool.Get().(*worker)
	w.pool = p
	w.run()
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
