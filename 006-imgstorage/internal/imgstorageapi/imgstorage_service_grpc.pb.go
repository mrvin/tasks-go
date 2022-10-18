// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: imgstorage_service.proto

package imgstorageapi

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

// ImgStorageClient is the client API for ImgStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ImgStorageClient interface {
	UploadImg(ctx context.Context, in *Img, opts ...grpc.CallOption) (*Null, error)
	DownloadImg(ctx context.Context, in *NameImg, opts ...grpc.CallOption) (*Img, error)
}

type imgStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewImgStorageClient(cc grpc.ClientConnInterface) ImgStorageClient {
	return &imgStorageClient{cc}
}

func (c *imgStorageClient) UploadImg(ctx context.Context, in *Img, opts ...grpc.CallOption) (*Null, error) {
	out := new(Null)
	err := c.cc.Invoke(ctx, "/imgstorage.ImgStorage/UploadImg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imgStorageClient) DownloadImg(ctx context.Context, in *NameImg, opts ...grpc.CallOption) (*Img, error) {
	out := new(Img)
	err := c.cc.Invoke(ctx, "/imgstorage.ImgStorage/DownloadImg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImgStorageServer is the server API for ImgStorage service.
// All implementations must embed UnimplementedImgStorageServer
// for forward compatibility
type ImgStorageServer interface {
	UploadImg(context.Context, *Img) (*Null, error)
	DownloadImg(context.Context, *NameImg) (*Img, error)
}

// UnimplementedImgStorageServer must be embedded to have forward compatible implementations.
type UnimplementedImgStorageServer struct {
}

func (UnimplementedImgStorageServer) UploadImg(context.Context, *Img) (*Null, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadImg not implemented")
}
func (UnimplementedImgStorageServer) DownloadImg(context.Context, *NameImg) (*Img, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadImg not implemented")
}

func RegisterImgStorageServer(s grpc.ServiceRegistrar, srv ImgStorageServer) {
	s.RegisterService(&ImgStorage_ServiceDesc, srv)
}

func _ImgStorage_UploadImg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Img)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImgStorageServer).UploadImg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imgstorage.ImgStorage/UploadImg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImgStorageServer).UploadImg(ctx, req.(*Img))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImgStorage_DownloadImg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NameImg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImgStorageServer).DownloadImg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imgstorage.ImgStorage/DownloadImg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImgStorageServer).DownloadImg(ctx, req.(*NameImg))
	}
	return interceptor(ctx, in, info, handler)
}

// ImgStorage_ServiceDesc is the grpc.ServiceDesc for ImgStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ImgStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "imgstorage.ImgStorage",
	HandlerType: (*ImgStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadImg",
			Handler:    _ImgStorage_UploadImg_Handler,
		},
		{
			MethodName: "DownloadImg",
			Handler:    _ImgStorage_DownloadImg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "imgstorage_service.proto",
}
