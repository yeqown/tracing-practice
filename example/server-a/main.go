package main

import (
	"context"
	"log"
	"net"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
)

var (
	addr = ":8081"

	serverBAddr = "127.0.0.1:8082"
	serverCAddr = "127.0.0.1:8083"

	zipkinTracer *zipkin.Tracer
)

func bootstrap() {
	var err error
	// Set up opentracing tracer
	zipkinTracer, err = bootTracer()
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
	pb.RegisterPingServer(s, newPingA())

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingA struct {
	serverBConn pb.PingClient
	serverCConn pb.PingClient
}

func newPingA() *pingA {
	// Set up a connection to the server.
	bConn, err := grpc.Dial(serverBAddr,
		grpc.WithInsecure(),
		// grpc.WithStatsHandler(zipkingrpc.NewClientHandler(zipkinTracer)),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cConn, err := grpc.Dial(serverCAddr,
		grpc.WithInsecure(),
		// grpc.WithStatsHandler(zipkingrpc.NewClientHandler(zipkinTracer)),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &pingA{
		serverBConn: pb.NewPingClient(bConn),
		serverCConn: pb.NewPingClient(cConn),
	}
}

func (p pingA) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingResponse, error) {
	// TODO: call server-B and server-C
	_, err := p.serverBConn.Ping(ctx, req)
	if err != nil {
		return nil, err
	}
	_, err = p.serverCConn.Ping(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := new(pb.PingResponse)
	return resp, nil
}
