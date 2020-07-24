package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/yeqown/opentracing-practice/protogen"
	"github.com/yeqown/opentracing-practice/x"
	opentracingrpc "github.com/yeqown/opentracing-practice/x/grpc-interceptor"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var (
	addr = "127.0.0.1:8082"
)

func bootstrap() {
	err := x.BootTracerWrapper("service-b", addr)
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
		grpc.UnaryInterceptor(opentracingrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), opentracingrpc.LogPayloads())),
	)
	pb.RegisterPingBServer(s, &pingB{})

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingB struct{}

func (p pingB) PingB(ctx context.Context, req *pb.PingBReq) (*pb.PingBResponse, error) {
	x.LogWithContext(ctx, "PingB calling")
	return &pb.PingBResponse{
		Now: time.Now().Unix(),
	}, nil
}
