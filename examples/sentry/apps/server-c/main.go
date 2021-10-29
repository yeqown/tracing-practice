package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	x2 "examples/opentracing/x"

	pb "github.com/yeqown/tracing-practice/api"
)

var (
	addr = "127.0.0.1:8083"
)

func bootstrap() {
	err := x2.BootTracerWrapper("service-c", addr)
	if err != nil {
		log.Fatalf("did not boot tracer: %v", err)
	}
}

func main() {
	bootstrap()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(x2.OpenTracingServerInterceptor(opentracing.GlobalTracer(), x2.LogPayloads())),
	)
	pb.RegisterPingCServer(s, &pingC{})

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingC struct {
	pb.UnimplementedPingCServer
}

func (p pingC) PingC(ctx context.Context, req *pb.PingCReq) (*pb.PingCResponse, error) {
	x2.LogWithContext(ctx, "PingC calling")
	if err := processInternalTrace3(ctx); err != nil {
		return nil, err
	}

	return &pb.PingCResponse{
		Now: time.Now().Unix(),
	}, nil
}

func processInternalTrace3(ctx context.Context) error {
	_, sp := x2.StartSpanFromContext(ctx)
	defer sp.Finish()

	// do some operation
	time.Sleep(3 * time.Millisecond)

	return nil
}
