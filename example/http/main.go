package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	// "github.com/opentracing/opentracing-go"
)

var (
	serverAAddr = "127.0.0.1:8081"
	addr        = ":8080"

	serverAConn pb.PingClient
)

func init() {
	// Set up a connection to the server.
	aConn, err := grpc.Dial(serverAAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	serverAConn = pb.NewPingClient(aConn)
}

func main() {
	engi := gin.New()
	// TODO: writing a middleware to generate a Context to pass by

	engi.GET("/trace", traceHdl)

	if err := engi.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func traceHdl(c *gin.Context) {
	// TODO: generate root Context and tracer
	rootCtx := context.Background()

	clientCall(rootCtx)

	c.JSON(http.StatusOK, gin.H{"traceId": "todo-trace-id"})
}

func clientCall(ctx context.Context) {
	// first call remote servers
	_, _ = serverAConn.Ping(ctx, &pb.PingReq{})

	// then call internal process
	processInternalTrace(ctx)
}

// 应用内部的追踪
func processInternalTrace(ctx context.Context) {
	// TODO: log ctx traceId and maybe other data
	return
}
