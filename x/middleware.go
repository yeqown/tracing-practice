package x

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
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

		if err == nil && clientSpCtx != nil {
			log.Printf("traceMiddleware called 1: clientTraceId=%s\n", getTraceIdFromSpanContext(clientSpCtx))
			sp = tracer.StartSpan(c.Request.RequestURI,
				opentracing.ChildOf(clientSpCtx),
			)
		} else {
			// if context could not get from headers
			println("traceMiddleware called 2")
			sp = tracer.StartSpan(c.Request.RequestURI)
		}
		defer sp.Finish()

		// 记录annotations
		//sp.LogFields(
		//	opentracinglog.Object("call service a", ""),
		//	opentracinglog.Object("call service b", ""),
		//)

		sp.SetTag("Method", c.Request.Method)
		sp.SetTag("Path", c.Request.URL)
		sp.SetTag("Request", "todo add request data")
		sp.SetTag("Response", "todo add response body")

		ctx = opentracing.ContextWithSpan(c.Request.Context(), sp)
		c.Set(_traceContextKey, ctx)

		traceId := getTraceIdFromSpanContext(sp.Context())
		log.Printf("traceId=%s\n", traceId)
		c.Header("X-Trace-Id", traceId)

		// continue process request
		c.Next()
	}
}
