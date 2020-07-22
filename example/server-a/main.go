package main

import (
	"context"
	"log"
	"net"

	pb "github.com/yeqown/opentracing-practice/protogen"

	"google.golang.org/grpc"
)

var (
	addr = ":8081"

	serverBAddr = "127.0.0.1:8082"
	serverCAddr = "127.0.0.1:8083"
)

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
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
	bConn, err := grpc.Dial(serverBAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cConn, err := grpc.Dial(serverBAddr, grpc.WithInsecure())
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

	resp := new(pb.PingResponse)
	return resp, nil
}
