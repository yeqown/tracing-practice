package main

import (
	"context"
	"log"
	"net"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"google.golang.org/grpc"
)

var (
	addr = ":8083"
)

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPingServer(s, &pingC{})

	log.Println("running on: ", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pingC struct{}

func (p pingC) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingResponse, error) {
	resp := new(pb.PingResponse)
	return resp, nil
}
