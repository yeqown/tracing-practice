package main

import (
	"context"
	"log"
	"net"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"google.golang.org/grpc"
)

var (
	addr = ":8082"
)

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
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
