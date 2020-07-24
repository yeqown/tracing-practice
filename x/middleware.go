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

// get trace info from header, if not then create an new one
func Opentracing() gin.HandlerFunc {
	tracer := opentracing.GlobalTracer()
	if tracer == nil {
		panic("tracer not set")
	}

	return func(c *gin.Context) {
		var (
			ctx context.Context
			sp  opentracing.Span
		)

		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		clientSpCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			log.Printf("could not extract trace data from http header, err=%v\n", err)
		}

		// derive a span or create an root span
		sp = tracer.StartSpan(
			c.Request.RequestURI,
			opentracing.ChildOf(clientSpCtx),
		)
		defer sp.Finish()

		// record and log traceId
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
