package x

import (
	"context"
	"fmt"
	"log"
	"runtime"

	xzipkin "github.com/yeqown/opentracing-practice/x/x-zipkin"

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

const (
	_spanContextFormat = "context{traceId: %s} "
)

func LogWithContext(ctx context.Context, format string, args ...interface{}) {
	sp := opentracing.SpanFromContext(ctx)
	format = fmt.Sprintf(_spanContextFormat, getTraceIdFromSpanContext(sp.Context())) + format
	log.Printf(format, args...)
}

func getTraceIdFromSpanContext(spCtx opentracing.SpanContext) string {
	return xzipkin.GetTraceIdFromSpanContext(spCtx)
	// return xjaeger.GetTraceIdFromSpanContext(spCtx)
}
