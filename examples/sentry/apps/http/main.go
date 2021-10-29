package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	pb "github.com/yeqown/tracing-practice/api"
	"github.com/yeqown/tracing-practice/examples/opentracing/x"
)

var (
	serverAAddr = "127.0.0.1:8081"
	addr        = "127.0.0.1:8080"

	serverAConn pb.PingAClient
)

func bootstrap() {
	err := x2.BootTracerWrapper("http-port", addr)
	if err != nil {
		log.Fatalf("did not boot tracer: %v", err)
	}

	// Set up a connection to the server-A.
	aConn, err := grpc.Dial(serverAAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(x2.OpenTracingClientInterceptor(opentracing.GlobalTracer(), x2.LogPayloads())),
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
	engi.Use(x.Opentracing())
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
	ctx := x.ExtractTraceContext(c)

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
	ctx2, sp := x2.StartSpanFromContext(ctx)
	defer sp.Finish()

	println("processInternalTrace1 called")
	// do some ops
	time.Sleep(10 * time.Millisecond)

	return processInternalTrace2(ctx2)
}

func processInternalTrace2(ctx context.Context) error {
	_, sp := x2.StartSpanFromContext(ctx)
	defer sp.Finish()

	println("processInternalTrace2 called")
	time.Sleep(5 * time.Millisecond)
	return nil
}
