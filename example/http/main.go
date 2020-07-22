package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
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

// 先创建“HTTP Collector”(the agent)用来收集跟踪数据并将其发送到Zipkin-UI，endpointUrl是Zipkin UI的URL
// 其次创建了一个记录器(recorder)来记录端口上的信息，“hostUrl”是gRPC(客户端)呼叫的URL
// 第三，用我们新建的记录器创建了一个新的跟踪器(tracer)
// 最后，为“OpenTracing”设置了“GlobalTracer”，这样你可以在程序中的任何地方访问它。
var (
	endpointUrl            = "http://localhost:9411/api/v1/spans"
	hostUrl                = "localhost:5051"
	serviceNameCacheClient = "cache service client"
)

func newTracer() (opentracing.Tracer, zipkintracer.Collector, error) {
	collector, err := openzipkin.NewHTTPCollector(endpointUrl)
	if err != nil {
		return nil, nil, err
	}
	recorder := openzipkin.NewRecorder(collector, true, hostUrl, serviceNameCacheClient)
	tracer, err := openzipkin.NewTracer(
		recorder,
		openzipkin.ClientServerSameSpan(true))

	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)

	return tracer, collector, nil
}
