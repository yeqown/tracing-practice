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
	addr = "127.0.0.1:8082"
)

func bootstrap() {
	err := x2.BootTracerWrapper("service-b", addr)
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
	pb.RegisterPingBServer(s, &pingB{})

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingB struct {
	pb.UnimplementedPingBServer
}

func (p pingB) PingB(ctx context.Context, req *pb.PingBReq) (*pb.PingBResponse, error) {
	x2.LogWithContext(ctx, "PingB calling")
	return &pb.PingBResponse{
		Now: time.Now().Unix(),
	}, nil
}
