package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	tracegin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/yeqown/tracing-practice/api"
	"github.com/yeqown/tracing-practice/examples/sentry/x"
)

var (
	serverAAddr = "127.0.0.1:8081"
	addr        = "127.0.0.1:8080"

	serverAConn pb.PingAClient
)

func bootstrap() {
	// Set up a connection to the server-A.
	aConn, err := grpc.Dial(serverAAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(x.UnaryClientInterceptor()),
		// grpc.WithUnaryInterceptor(x2.OpenTracingClientInterceptor(opentracing.GlobalTracer(), x2.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	serverAConn = pb.NewPingAClient(aConn)
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://1c2d1ae347944688ae7593a33e40c0f2@sentry.example.com/33",
		ServerName:  "http-demo",
		Environment: "dev",
		Release:     "v1.0.0",
		SampleRate:  1.0,
	})
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// prepare necessary data
	bootstrap()
	//defer sentry.Flush(3 * time.Second)

	// prepare HTTP server
	engi := gin.New()

	// a middleware to generate a Context to pass by
	// it also parse trace info from client request header
	engi.Use(tracegin.New(tracegin.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         0,
	}))
	engi.Use(func(c *gin.Context) {

		//hub := tracegin.GetHubFromContext(c)
		sp := sentry.StartSpan(
			c.Request.Context(),
			c.FullPath(),
			sentry.ContinueFromRequest(c.Request),
			sentry.TransactionName(c.FullPath()),
		)
		//sp.Sampled = sentry.SampledTrue
		defer sp.Finish()
		println(sp.ToSentryTrace())

		c.Set("sentry.gin", sp.Context())
		c.Next()

		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("client", "i'm client")
			scope.SetExtra("http.status", c.Writer.Status())
		})
	})
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
	v, ok := c.Get("sentry.gin")
	if !ok {
		panic("invalid")
	}
	ctx := v.(context.Context)

	// process request call, remote and local process
	if err := clientCall(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// response to client
	c.JSON(http.StatusOK, gin.H{"message": "traceHdl done"})
}

func clientCall(ctx context.Context) error {
	if err := processInternalTrace1(ctx); err != nil {
		return err
	}

	// first call remote servers
	_, err := serverAConn.PingA(ctx, &pb.PingAReq{
		Now:  time.Now().Unix(),
		From: "client",
	})

	return err
}

// internal process trace examples 1
func processInternalTrace1(ctx context.Context) error {
	sp := sentry.StartSpan(ctx, "processInternalTrace1")
	defer sp.Finish()

	println("processInternalTrace1 called")
	// do some ops
	time.Sleep(10 * time.Millisecond)

	return processInternalTrace2(sp.Context())
}

func processInternalTrace2(ctx context.Context) error {
	sp := sentry.StartSpan(ctx, "processInternalTrace2")
	defer sp.Finish()

	println("processInternalTrace2 called")
	time.Sleep(5 * time.Millisecond)
	return nil
}
