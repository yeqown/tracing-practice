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

// TODO: ensure tracer cannot be nil
func DeriveFromContext(ctx context.Context) (context.Context, opentracing.Span) {
	tracer := opentracing.GlobalTracer()
	if tracer == nil {
		// panic("tracer not set")
		return nil, nil
	}

	var opName = "notset"

	// parent span could be parsed from `ctx`
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		opName = WhoCalling()
		println("deliver a child span", opName)
		sp := tracer.StartSpan(opName, opentracing.ChildOf(parent.Context()))
		ctx = opentracing.ContextWithSpan(ctx, sp)
		return ctx, sp

	}

	// could not get parent span from `ctx`
	opName = WhoCalling()
	println("a root span", opName)
	sp := tracer.StartSpan(opName)
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp
}
