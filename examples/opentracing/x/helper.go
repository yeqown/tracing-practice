package x

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func BootTracerWrapper(localServiceName string, hostPort string) error {
	// tracer, err := xzipkin.BootZipkinTracer(localServiceName, hostPort)
	tracer, err := BootJaegerTracer(localServiceName, hostPort)
	if err != nil {
		return errors.Wrap(err, "BootTracerWrapper.BootZipkinTracer")
	}
	opentracing.SetGlobalTracer(tracer)

	return nil
}

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
	format = fmt.Sprintf(_spanContextFormat, GetTraceIdFromSpanContext(sp.Context())) + format
	log.Printf(format, args...)
}
