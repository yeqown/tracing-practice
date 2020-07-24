package x

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

const (
	_traceContextKey = "traceContext"
)

func GetTraceContextKey() string {
	return _traceContextKey
}

type getTraceID func(spCtx opentracing.SpanContext) string

// get trace info from header, if not then create an new one
func Opentracing(getTraceIdFromSpanContext getTraceID) gin.HandlerFunc {
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

		// FIXME: err always be nil, why ?
		if err == nil && clientSpCtx != nil {
			log.Printf("Opentracing get client span: clientTraceId=%s\n", getTraceIdFromSpanContext(clientSpCtx))
			sp = tracer.StartSpan(
				c.Request.RequestURI,             // TODO: regenerate request operationName
				opentracing.ChildOf(clientSpCtx), // TODO: more options
			)
		} else {
			sp = tracer.StartSpan(c.Request.RequestURI)
		}
		defer sp.Finish()

		traceId := getTraceIdFromSpanContext(sp.Context())
		c.Header("X-Trace-Id", traceId)
		log.Println("request with traceId:", traceId)

		start := time.Now()
		sp.LogFields(opentracinglog.Int64("start", start.Unix()))
		sp.SetTag("Method", c.Request.Method)
		sp.SetTag("Path", c.Request.URL)
		sp.SetTag("Request", "todo add request data")
		sp.SetTag("Response", "todo add response body")

		ctx = opentracing.ContextWithSpan(c.Request.Context(), sp)
		c.Set(_traceContextKey, ctx)

		// continue process request
		c.Next()

		end := time.Now()
		sp.SetTag("latency (ms)", end.Sub(start).Milliseconds())
		sp.LogFields(opentracinglog.Int64("finish", end.Unix()))
	}
}
