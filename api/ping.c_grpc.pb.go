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

// PingCClient is the client API for PingC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PingCClient interface {
	PingC(ctx context.Context, in *PingCReq, opts ...grpc.CallOption) (*PingCResponse, error)
}

type pingCClient struct {
	cc grpc.ClientConnInterface
}

func NewPingCClient(cc grpc.ClientConnInterface) PingCClient {
	return &pingCClient{cc}
}

func (c *pingCClient) PingC(ctx context.Context, in *PingCReq, opts ...grpc.CallOption) (*PingCResponse, error) {
	out := new(PingCResponse)
	err := c.cc.Invoke(ctx, "/pb.PingC/PingC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PingCServer is the server API for PingC service.
// All implementations must embed UnimplementedPingCServer
// for forward compatibility
type PingCServer interface {
	PingC(context.Context, *PingCReq) (*PingCResponse, error)
	mustEmbedUnimplementedPingCServer()
}

// UnimplementedPingCServer must be embedded to have forward compatible implementations.
type UnimplementedPingCServer struct {
}

func (UnimplementedPingCServer) PingC(context.Context, *PingCReq) (*PingCResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingC not implemented")
}
func (UnimplementedPingCServer) mustEmbedUnimplementedPingCServer() {}

// UnsafePingCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PingCServer will
// result in compilation errors.
type UnsafePingCServer interface {
	mustEmbedUnimplementedPingCServer()
}

func RegisterPingCServer(s grpc.ServiceRegistrar, srv PingCServer) {
	s.RegisterService(&PingC_ServiceDesc, srv)
}

func _PingC_PingC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingCReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingCServer).PingC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PingC/PingC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingCServer).PingC(ctx, req.(*PingCReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PingC_ServiceDesc is the grpc.ServiceDesc for PingC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PingC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PingC",
	HandlerType: (*PingCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingC",
			Handler:    _PingC_PingC_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ping.c.proto",
}
