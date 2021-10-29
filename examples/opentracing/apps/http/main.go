package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	opentracinglog "github.com/opentracing/opentracing-go/log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	pb "github.com/yeqown/tracing-practice/api"
	x "github.com/yeqown/tracing-practice/examples/opentracing/x"
)

var (
	serverAAddr = "127.0.0.1:8081"
	addr        = "127.0.0.1:8080"

	serverAConn pb.PingAClient
)

func bootstrap() {
	err := x.BootTracerWrapper("http-port", addr)
	if err != nil {
		log.Fatalf("did not boot tracer: %v", err)
	}

	// Set up a connection to the server-A.
	aConn, err := grpc.Dial(serverAAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(x.OpenTracingClientInterceptor(opentracing.GlobalTracer(), x.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	serverAConn = pb.NewPingAClient(aConn)
}

func main() {
	// prepare necessary data
	bootstrap()

	// prepare HTTP server
	engi := gin.New()

	// a middleware to generate a Context to pass by
	// it also parse trace info from client request header
	engi.Use(opentracingMiddleware())
	engi.GET("/trace", traceHdl)

	// running HTTP server
	if err := engi.Run(addr); err != nil {
		log.Fatal(err)
	}
}

// traceHdl is a trace handler from HTTP request
func traceHdl(c *gin.Context) {
	// get root Context from request
	// TODO: try to use c.Request.WithContext() to set context
	ctx := extractTraceContext(c)

	// process request call, remote and local process
	if err := clientCall(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"message": "traceHdl done"})
}

func clientCall(ctx context.Context) error {
	// first call remote servers
	_, err := serverAConn.PingA(ctx, &pb.PingAReq{
		Now:  time.Now().Unix(),
		From: "client",
	})

	if err != nil {
		return err
	}

	// then call internal process
	return processInternalTrace1(ctx)
}

// internal process trace examples 1
func processInternalTrace1(ctx context.Context) error {
	ctx2, sp := x.StartSpanFromContext(ctx)
	defer sp.Finish()

	println("processInternalTrace1 called")
	// do some ops
	time.Sleep(10 * time.Millisecond)

	return processInternalTrace2(ctx2)
}

func processInternalTrace2(ctx context.Context) error {
	_, sp := x.StartSpanFromContext(ctx)
	defer sp.Finish()

	println("processInternalTrace2 called")
	time.Sleep(5 * time.Millisecond)
	return nil
}

var ginTraceContextKey = "opentracing.gin"

type respBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w respBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// opentracingMiddleware get trace info from header, if not then create an new one
func opentracingMiddleware() gin.HandlerFunc {
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
		traceId := x.GetTraceIdFromSpanContext(sp.Context())
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

func headerToFields(header http.Header) []opentracinglog.Field {
	fields := make([]opentracinglog.Field, 0, len(http.Header{}))

	for k, v := range header {
		fields = append(fields, opentracinglog.String(k, strings.Join(v, ";")))
	}

	return fields
}

func injectIntoGinContext(c *gin.Context, ctx context.Context) {
	c.Set(ginTraceContextKey, ctx)
}

func extractTraceContext(c *gin.Context) context.Context {
	v, ok := c.Get(ginTraceContextKey)
	if !ok {
		return context.TODO()
	}

	ctx, ok := v.(context.Context)
	if !ok {
		return context.TODO()
	}

	return ctx
}
