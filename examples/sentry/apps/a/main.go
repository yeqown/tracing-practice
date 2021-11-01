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
	addr = "127.0.0.1:8081"

	serverBAddr = "127.0.0.1:8082"
	serverCAddr = "127.0.0.1:8083"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		//Dsn: "https://af1a4d56a9a349e08ff1581c0d1c8d5a@sentry.example.com/35",
		Dsn:         "https://1c2d1ae347944688ae7593a33e40c0f2@sentry.example.com/33",
		ServerName:  "a",
		Environment: "dev",
		Release:     "v1.0.0",
		SampleRate:  1.0,
	})
	defer sentry.Flush(2 * time.Second)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(x.UnaryServerInterceptor()),
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
		grpc.WithUnaryInterceptor(x.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cConn, err := grpc.Dial(serverCAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(x.UnaryClientInterceptor()),
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
