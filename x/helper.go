package x

import (
	"context"
	"runtime"

	"github.com/opentracing/opentracing-go"
)

func WhoCalling() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func StartSpanFromContext(ctx context.Context) (context.Context, opentracing.Span) {
	opName := WhoCalling()
	sp, ctx := opentracing.StartSpanFromContext(ctx, opName)
	return ctx, sp
}
