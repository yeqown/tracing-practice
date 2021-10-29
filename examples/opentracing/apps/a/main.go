package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"

	pb "github.com/yeqown/tracing-practice/api"
	"github.com/yeqown/tracing-practice/examples/opentracing/x"
)

var (
	addr = "127.0.0.1:8081"

	serverBAddr = "127.0.0.1:8082"
	serverCAddr = "127.0.0.1:8083"

	zipkinTracer *zipkin.Tracer
)

func bootstrap() {
	err := x.BootTracerWrapper("service-a", addr)
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
		grpc.UnaryInterceptor(x.OpenTracingServerInterceptor(
			opentracing.GlobalTracer(), x.LogPayloads())),
	)
	pb.RegisterPingAServer(s, newPingA())

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingA struct {
	pb.UnimplementedPingAServer

	serverBConn pb.PingBClient
	serverCConn pb.PingCClient
}

func newPingA() *pingA {
	// Set up a connection to the server.
	bConn, err := grpc.Dial(serverBAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(x.OpenTracingClientInterceptor(opentracing.GlobalTracer(), x.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cConn, err := grpc.Dial(serverCAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(x.OpenTracingClientInterceptor(opentracing.GlobalTracer(), x.LogPayloads())),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &pingA{
		serverBConn: pb.NewPingBClient(bConn),
		serverCConn: pb.NewPingCClient(cConn),
	}
}

func (p pingA) PingA(ctx context.Context, req *pb.PingAReq) (*pb.PingAResponse, error) {
	x.LogWithContext(ctx, "PingA calling")

	// call server-B and server-C
	_, err := p.serverBConn.PingB(ctx, &pb.PingBReq{
		Now:  req.Now,
		From: "a",
	})
	if err != nil {
		return nil, err
	}
	_, err = p.serverCConn.PingC(ctx, &pb.PingCReq{
		Now:  req.Now,
		From: "a",
	})
	if err != nil {
		return nil, err
	}

	return &pb.PingAResponse{
		Now: time.Now().Unix(),
	}, nil
}
