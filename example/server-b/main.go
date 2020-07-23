package main

import (
	"context"
	"log"
	"net"

	"github.com/yeqown/opentracing-practice/x"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yeqown/opentracing-practice/protogen"

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
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())),
	)
	pb.RegisterPingServer(s, &pingB{})

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingB struct{}

func (p pingB) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingResponse, error) {
	resp := new(pb.PingResponse)
	return resp, nil
}
