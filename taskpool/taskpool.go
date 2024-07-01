package taskpool

import (
	"context"
)

var defaultPool = new(pool)

func init() {
	defaultPool = NewPool("defaultPool", 100, NewDefaultConfig())
}

func GoWithoutCtx(f func()) {
	GoCtx(context.Background(), f)
}

func GoCtx(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}
