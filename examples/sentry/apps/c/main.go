package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"

	pb "github.com/yeqown/tracing-practice/api"
	"github.com/yeqown/tracing-practice/examples/sentry/x"
)

var (
	addr = "127.0.0.1:8083"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://1c2d1ae347944688ae7593a33e40c0f2@sentry.example.com/33",
		ServerName:  "c",
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
	x.LogWithContext(ctx, "PingC calling")

	if r := rand.Intn(100); r <= 50 {
		println("paniced")
		panic(errors.New("random panic"))
	}

	if err := processInternalTrace3(ctx); err != nil {
		return nil, err
	}

	return &pb.PingCResponse{
		Now: time.Now().Unix(),
	}, nil
}

func processInternalTrace3(ctx context.Context) error {
	sp := sentry.StartSpan(ctx, "processInternalTrace3")
	defer sp.Finish()

	// do some operation
	time.Sleep(3 * time.Millisecond)

	return nil
}
