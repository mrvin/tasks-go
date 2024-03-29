// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: fibonacci_service.proto

package fibonacciapi

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

const (
	Fib_Get_FullMethodName = "/fibonacciapi.Fib/Get"
)

// FibClient is the client API for Fib service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FibClient interface {
	Get(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type fibClient struct {
	cc grpc.ClientConnInterface
}

func NewFibClient(cc grpc.ClientConnInterface) FibClient {
	return &fibClient{cc}
}

func (c *fibClient) Get(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Fib_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FibServer is the server API for Fib service.
// All implementations should embed UnimplementedFibServer
// for forward compatibility
type FibServer interface {
	Get(context.Context, *Request) (*Response, error)
}

// UnimplementedFibServer should be embedded to have forward compatible implementations.
type UnimplementedFibServer struct {
}

func (UnimplementedFibServer) Get(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

// UnsafeFibServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FibServer will
// result in compilation errors.
type UnsafeFibServer interface {
	mustEmbedUnimplementedFibServer()
}

func RegisterFibServer(s grpc.ServiceRegistrar, srv FibServer) {
	s.RegisterService(&Fib_ServiceDesc, srv)
}

func _Fib_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FibServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Fib_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FibServer).Get(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Fib_ServiceDesc is the grpc.ServiceDesc for Fib service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fib_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fibonacciapi.Fib",
	HandlerType: (*FibServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Fib_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fibonacci_service.proto",
}
