package x

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	opentracinglog "github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"
	xjaeger "github.com/yeqown/opentracing-practice/x/x-jaeger"
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
	//return xzipkin.GetTraceIdFromSpanContext(spCtx)
	return xjaeger.GetTraceIdFromSpanContext(spCtx)
}

func headerToFields(header http.Header) []opentracinglog.Field {
	fields := make([]opentracinglog.Field, 0, len(http.Header{}))

	for k, v := range header {
		fields = append(fields, opentracinglog.String(k, strings.Join(v, ";")))
	}

	return fields
}
