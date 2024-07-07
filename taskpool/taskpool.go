package taskpool

import (
	"context"
	"fmt"
	"sync"
)

var poolMap sync.Map

var defaultPool = new(pool)

func init() {
	defaultPool = NewPool("defaultPool", 1000, NewDefaultConfig())
}

func GoWithoutCtx(f func()) {
	GoCtx(context.Background(), f)
}

func GoCtx(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}

func SetCap(cap int32) {
	defaultPool.SetCap(cap)
}

func RegisterPool(p *pool) error {
	_, loaded := poolMap.LoadOrStore(p.Name(), p)
	if loaded {
		return fmt.Errorf("name: %s already registered", p.Name())
	}
	return nil
}

func GetPool(name string) *pool {
	p, ok := poolMap.Load(name)
	if !ok {
		return nil
	}
	pInt, ok := p.(*pool)
	if !ok {
		return nil
	}
	return pInt
}
