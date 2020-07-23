package main

import (
	"context"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	endpointUrl = "http://localhost:9411/api/v2/spans"
	name        = "http"
)

// 先创建“HTTP Collector”(the agent)用来收集跟踪数据并将其发送到Zipkin-UI，endpointUrl是Zipkin UI的URL
// 其次创建了一个记录器(recorder)来记录端口上的信息，“hostUrl”是gRPC(客户端)呼叫的URL
// 第三，用我们新建的记录器创建了一个新的跟踪器(tracer)
// 最后，为“OpenTracing”设置了“GlobalTracer”，这样你可以在程序中的任何地方访问它。
func bootTracer() (*zipkin.Tracer, error) {
	reporter := zipkinhttp.NewReporter(endpointUrl)
	localEndpoint, err := zipkin.NewEndpoint("http", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	nativeTracer, err := zipkin.NewTracer(reporter,
		zipkin.WithTraceID128Bit(false), // TODO: diff between 128 and 64bit
		zipkin.WithSharedSpans(true),    // TODO: what effect
		zipkin.WithLocalEndpoint(localEndpoint),
		// TODO: more options
	)
	if err != nil {
		return nil, err
	}

	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)

	return nativeTracer, nil
}

func getTraceAndSetSpan(ctx context.Context) (context.Context, opentracing.Span) {
	var opName = "notset"

	// log ctx traceId and maybe other data
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			if opName == "notset" {
				opName = whoami()
			}
			sp := tracer.StartSpan(opName, opentracing.ChildOf(parent.Context()))
			ctx = opentracing.ContextWithSpan(ctx, sp)
			return ctx, sp
		}
	}

	return nil, nil
}

const (
	_traceContextKey = "traceContext"
)

// get trace info from header, if not then create an new one
func traceMiddleware() gin.HandlerFunc {
	tracer := opentracing.GlobalTracer()
	if tracer == nil {
		panic("tracer not set")
	}

	return func(c *gin.Context) {
		var (
			clientSpCtx opentracing.SpanContext
			sp          opentracing.Span
			ctx         context.Context
		)

		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		clientSpCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

		if err == nil && clientSpCtx != nil {
			println("traceMiddleware called 1")
			sp = tracer.StartSpan(c.Request.RequestURI,
				opentracing.ChildOf(clientSpCtx),
			)
		} else {
			// if context could not get from headers
			println("traceMiddleware called 2")
			sp = tracer.StartSpan(c.Request.RequestURI)
		}
		defer sp.Finish()

		sp.LogFields(
			opentracinglog.Object("request", "TODO: add form data"),
			opentracinglog.Object("response", "TODO: add response"),
		)

		ctx = opentracing.ContextWithSpan(c.Request.Context(), sp)
		c.Set(_traceContextKey, ctx)

		traceId := getTraceIdFromSpanContext(sp.Context())
		log.Printf("traceId=%s\n", traceId)
		c.Header("X-Trace-Id", traceId)

		// continue process request
		c.Next()
	}
}

func getTraceIdFromSpanContext(spanCtx opentracing.SpanContext) string {
	return spanCtx.(zipkinot.SpanContext).TraceID.String()
}

func childContextOfSpan(ctx context.Context, sp opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, sp)
}

func whoami() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}
