package x

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

const (
	_traceContextKey = "traceContext"
)

type respBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w respBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

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
		rbw := &respBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = rbw
		body, err := c.GetRawData()
		if err == nil && len(body) != 0 {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		// try to parse context from HTTP request header
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		clientSpCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			log.Printf("could not extract trace data from http header, err=%v\n", err)
		}

		// derive a span or create an root span
		operation := c.FullPath()
		sp := tracer.StartSpan(
			operation,
			opentracing.ChildOf(clientSpCtx),
		)
		defer sp.Finish()

		// restful tags to for searching
		sp.SetTag("method", c.Request.Method)
		if len(c.Params) != 0 {
			for _, v := range c.Params {
				sp.SetTag("http.params."+v.Key, v.Value)
			}
		}

		// record and log traceId
		traceId := getTraceIdFromSpanContext(sp.Context())
		c.Header("X-Trace-Id", traceId)
		// log.Println("request with traceId:", traceId)

		// fields recorded
		sp.LogFields(
			opentracinglog.String("request.query", c.Request.URL.RawQuery),
			opentracinglog.String("request.body", string(body)),
		)
		sp.LogFields(headerToFields(c.Request.Header)...)

		// inject into gin.Context so it can be propagate into downstream servers.
		injectIntoGinContext(c, opentracing.ContextWithSpan(c.Request.Context(), sp))

		// continue process request
		c.Next()

		// all handlers are finished, so record response message those may be needed.
		// status code into tag
		sp.SetTag("http.status", c.Writer.Status())
		fields := make([]opentracinglog.Field, 0, 1)
		if c.Writer.Status() >= http.StatusBadRequest {
			fields = append(fields, opentracinglog.String("response.body", rbw.body.String()))
		}

		if len(fields) > 0 {
			sp.LogFields(fields...)
		}
	}
}

func injectIntoGinContext(c *gin.Context, ctx context.Context) {
	c.Set(_traceContextKey, ctx)
}

func ExtractTraceContext(c *gin.Context) context.Context {
	v, ok := c.Get(_traceContextKey)
	if !ok {
		return context.TODO()
	}

	ctx, ok := v.(context.Context)
	if !ok {
		return context.TODO()
	}

	return ctx
}
