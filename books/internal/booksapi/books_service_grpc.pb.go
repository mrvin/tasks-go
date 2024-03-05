// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: books_service.proto

package booksapi

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BookService_CreateBook_FullMethodName        = "/books.BookService/CreateBook"
	BookService_GetBooksByAuthor_FullMethodName  = "/books.BookService/GetBooksByAuthor"
	BookService_GetAuthorsByTitle_FullMethodName = "/books.BookService/GetAuthorsByTitle"
)

// BookServiceClient is the client API for BookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookServiceClient interface {
	CreateBook(ctx context.Context, in *CreateBookRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetBooksByAuthor(ctx context.Context, in *GetBooksByAuthorRequest, opts ...grpc.CallOption) (*GetBooksByAuthorResponse, error)
	GetAuthorsByTitle(ctx context.Context, in *GetAuthorsByTitleRequest, opts ...grpc.CallOption) (*GetAuthorsByTitleResponse, error)
}

type bookServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookServiceClient(cc grpc.ClientConnInterface) BookServiceClient {
	return &bookServiceClient{cc}
}

func (c *bookServiceClient) CreateBook(ctx context.Context, in *CreateBookRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BookService_CreateBook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) GetBooksByAuthor(ctx context.Context, in *GetBooksByAuthorRequest, opts ...grpc.CallOption) (*GetBooksByAuthorResponse, error) {
	out := new(GetBooksByAuthorResponse)
	err := c.cc.Invoke(ctx, BookService_GetBooksByAuthor_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookServiceClient) GetAuthorsByTitle(ctx context.Context, in *GetAuthorsByTitleRequest, opts ...grpc.CallOption) (*GetAuthorsByTitleResponse, error) {
	out := new(GetAuthorsByTitleResponse)
	err := c.cc.Invoke(ctx, BookService_GetAuthorsByTitle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookServiceServer is the server API for BookService service.
// All implementations should embed UnimplementedBookServiceServer
// for forward compatibility
type BookServiceServer interface {
	CreateBook(context.Context, *CreateBookRequest) (*emptypb.Empty, error)
	GetBooksByAuthor(context.Context, *GetBooksByAuthorRequest) (*GetBooksByAuthorResponse, error)
	GetAuthorsByTitle(context.Context, *GetAuthorsByTitleRequest) (*GetAuthorsByTitleResponse, error)
}

// UnimplementedBookServiceServer should be embedded to have forward compatible implementations.
type UnimplementedBookServiceServer struct {
}

func (UnimplementedBookServiceServer) CreateBook(context.Context, *CreateBookRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBook not implemented")
}
func (UnimplementedBookServiceServer) GetBooksByAuthor(context.Context, *GetBooksByAuthorRequest) (*GetBooksByAuthorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBooksByAuthor not implemented")
}
func (UnimplementedBookServiceServer) GetAuthorsByTitle(context.Context, *GetAuthorsByTitleRequest) (*GetAuthorsByTitleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthorsByTitle not implemented")
}

// UnsafeBookServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookServiceServer will
// result in compilation errors.
type UnsafeBookServiceServer interface {
	mustEmbedUnimplementedBookServiceServer()
}

func RegisterBookServiceServer(s grpc.ServiceRegistrar, srv BookServiceServer) {
	s.RegisterService(&BookService_ServiceDesc, srv)
}

func _BookService_CreateBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).CreateBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookService_CreateBook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).CreateBook(ctx, req.(*CreateBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_GetBooksByAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBooksByAuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).GetBooksByAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookService_GetBooksByAuthor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).GetBooksByAuthor(ctx, req.(*GetBooksByAuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookService_GetAuthorsByTitle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthorsByTitleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookServiceServer).GetAuthorsByTitle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookService_GetAuthorsByTitle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookServiceServer).GetAuthorsByTitle(ctx, req.(*GetAuthorsByTitleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BookService_ServiceDesc is the grpc.ServiceDesc for BookService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "books.BookService",
	HandlerType: (*BookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBook",
			Handler:    _BookService_CreateBook_Handler,
		},
		{
			MethodName: "GetBooksByAuthor",
			Handler:    _BookService_GetBooksByAuthor_Handler,
		},
		{
			MethodName: "GetAuthorsByTitle",
			Handler:    _BookService_GetAuthorsByTitle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "books_service.proto",
}
