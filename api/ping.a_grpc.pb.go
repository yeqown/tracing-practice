// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PingAClient is the client API for PingA service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PingAClient interface {
	PingA(ctx context.Context, in *PingAReq, opts ...grpc.CallOption) (*PingAResponse, error)
}

type pingAClient struct {
	cc grpc.ClientConnInterface
}

func NewPingAClient(cc grpc.ClientConnInterface) PingAClient {
	return &pingAClient{cc}
}

func (c *pingAClient) PingA(ctx context.Context, in *PingAReq, opts ...grpc.CallOption) (*PingAResponse, error) {
	out := new(PingAResponse)
	err := c.cc.Invoke(ctx, "/pb.PingA/PingA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PingAServer is the server API for PingA service.
// All implementations must embed UnimplementedPingAServer
// for forward compatibility
type PingAServer interface {
	PingA(context.Context, *PingAReq) (*PingAResponse, error)
	mustEmbedUnimplementedPingAServer()
}

// UnimplementedPingAServer must be embedded to have forward compatible implementations.
type UnimplementedPingAServer struct {
}

func (UnimplementedPingAServer) PingA(context.Context, *PingAReq) (*PingAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingA not implemented")
}
func (UnimplementedPingAServer) mustEmbedUnimplementedPingAServer() {}

// UnsafePingAServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PingAServer will
// result in compilation errors.
type UnsafePingAServer interface {
	mustEmbedUnimplementedPingAServer()
}

func RegisterPingAServer(s grpc.ServiceRegistrar, srv PingAServer) {
	s.RegisterService(&PingA_ServiceDesc, srv)
}

func _PingA_PingA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingAReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingAServer).PingA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PingA/PingA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingAServer).PingA(ctx, req.(*PingAReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PingA_ServiceDesc is the grpc.ServiceDesc for PingA service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PingA_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PingA",
	HandlerType: (*PingAServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingA",
			Handler:    _PingA_PingA_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ping.a.proto",
}
