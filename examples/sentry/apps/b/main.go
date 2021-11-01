package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"

	pb "github.com/yeqown/tracing-practice/api"
	"github.com/yeqown/tracing-practice/examples/sentry/x"
)

var (
	addr = "127.0.0.1:8082"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://1c2d1ae347944688ae7593a33e40c0f2@sentry.example.com/33",
		ServerName:  "b",
		Environment: "dev",
		Release:     "v1.0.0",
		SampleRate:  1.0,
	})
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(x.UnaryServerInterceptor()),
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
	x.LogWithContext(ctx, "PingB calling")
	return &pb.PingBResponse{
		Now: time.Now().Unix(),
	}, nil
}
