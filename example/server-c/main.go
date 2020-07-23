package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/yeqown/opentracing-practice/x"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"google.golang.org/grpc"
)

var (
	addr = "127.0.0.1:8083"
)

func bootstrap() {
	err := x.BootTracerWrapper("service-c", addr)
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
	pb.RegisterPingServer(s, &pingC{})

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingC struct{}

func (p pingC) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingResponse, error) {
	if err := processInternalTrace3(ctx); err != nil {
		return nil, err
	}

	resp := new(pb.PingResponse)
	return resp, nil
}

func processInternalTrace3(ctx context.Context) error {
	_, sp := x.DeriveFromContext(ctx)
	defer sp.Finish()

	// do some operation
	time.Sleep(3 * time.Millisecond)

	return nil
}
