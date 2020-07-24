package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/yeqown/opentracing-practice/protogen"
	"github.com/yeqown/opentracing-practice/x"
	opentracingrpc "github.com/yeqown/opentracing-practice/x/grpc-interceptor"
	xzipkin "github.com/yeqown/opentracing-practice/x/x-zipkin"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var (
	serverAAddr = "127.0.0.1:8081"
	addr        = "127.0.0.1:8080"

	serverAConn pb.PingClient
)

func bootstrap() {
	err := x.BootTracerWrapper("http-port", addr)
	if err != nil {
		log.Fatalf("did not boot tracer: %v", err)
	}

	// Set up a connection to the server-A.
	aConn, err := grpc.Dial(serverAAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(opentracingrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(), opentracingrpc.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	serverAConn = pb.NewPingClient(aConn)
}

func main() {
	// prepare necessary data
	bootstrap()

	// prepare HTTP server
	engi := gin.New()

	// a middleware to generate a Context to pass by
	// it also parse trace info from client request header
	engi.Use(x.Opentracing(xzipkin.GetTraceIdFromSpanContext))
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
	ctx, ok := c.Get(x.GetTraceContextKey())
	if !ok {
		panic("impossible")
	}

	// process request call, remote and local process
	if err := clientCall(ctx.(context.Context)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"message": "traceHdl done"})
}

func clientCall(ctx context.Context) error {
	// first call remote servers
	_, err := serverAConn.Ping(ctx, &pb.PingReq{})

	if err != nil {
		return err
	}

	// then call internal process
	return processInternalTrace1(ctx)
}

// internal process trace example 1
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
