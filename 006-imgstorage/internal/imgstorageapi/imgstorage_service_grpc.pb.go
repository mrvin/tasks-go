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
	UploadImg(ctx context.Context, opts ...grpc.CallOption) (ImgStorage_UploadImgClient, error)
	DownloadImg(ctx context.Context, in *NameImg, opts ...grpc.CallOption) (ImgStorage_DownloadImgClient, error)
	GetListImg(ctx context.Context, in *Null, opts ...grpc.CallOption) (*ListImg, error)
}

type imgStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewImgStorageClient(cc grpc.ClientConnInterface) ImgStorageClient {
	return &imgStorageClient{cc}
}

func (c *imgStorageClient) UploadImg(ctx context.Context, opts ...grpc.CallOption) (ImgStorage_UploadImgClient, error) {
	stream, err := c.cc.NewStream(ctx, &ImgStorage_ServiceDesc.Streams[0], "/imgstorage.ImgStorage/UploadImg", opts...)
	if err != nil {
		return nil, err
	}
	x := &imgStorageUploadImgClient{stream}
	return x, nil
}

type ImgStorage_UploadImgClient interface {
	Send(*Img) error
	CloseAndRecv() (*Null, error)
	grpc.ClientStream
}

type imgStorageUploadImgClient struct {
	grpc.ClientStream
}

func (x *imgStorageUploadImgClient) Send(m *Img) error {
	return x.ClientStream.SendMsg(m)
}

func (x *imgStorageUploadImgClient) CloseAndRecv() (*Null, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Null)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *imgStorageClient) DownloadImg(ctx context.Context, in *NameImg, opts ...grpc.CallOption) (ImgStorage_DownloadImgClient, error) {
	stream, err := c.cc.NewStream(ctx, &ImgStorage_ServiceDesc.Streams[1], "/imgstorage.ImgStorage/DownloadImg", opts...)
	if err != nil {
		return nil, err
	}
	x := &imgStorageDownloadImgClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ImgStorage_DownloadImgClient interface {
	Recv() (*Img, error)
	grpc.ClientStream
}

type imgStorageDownloadImgClient struct {
	grpc.ClientStream
}

func (x *imgStorageDownloadImgClient) Recv() (*Img, error) {
	m := new(Img)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *imgStorageClient) GetListImg(ctx context.Context, in *Null, opts ...grpc.CallOption) (*ListImg, error) {
	out := new(ListImg)
	err := c.cc.Invoke(ctx, "/imgstorage.ImgStorage/GetListImg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImgStorageServer is the server API for ImgStorage service.
// All implementations should embed UnimplementedImgStorageServer
// for forward compatibility
type ImgStorageServer interface {
	UploadImg(ImgStorage_UploadImgServer) error
	DownloadImg(*NameImg, ImgStorage_DownloadImgServer) error
	GetListImg(context.Context, *Null) (*ListImg, error)
}

// UnimplementedImgStorageServer should be embedded to have forward compatible implementations.
type UnimplementedImgStorageServer struct {
}

func (UnimplementedImgStorageServer) UploadImg(ImgStorage_UploadImgServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadImg not implemented")
}
func (UnimplementedImgStorageServer) DownloadImg(*NameImg, ImgStorage_DownloadImgServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadImg not implemented")
}
func (UnimplementedImgStorageServer) GetListImg(context.Context, *Null) (*ListImg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListImg not implemented")
}

// UnsafeImgStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ImgStorageServer will
// result in compilation errors.
type UnsafeImgStorageServer interface {
	mustEmbedUnimplementedImgStorageServer()
}

func RegisterImgStorageServer(s grpc.ServiceRegistrar, srv ImgStorageServer) {
	s.RegisterService(&ImgStorage_ServiceDesc, srv)
}

func _ImgStorage_UploadImg_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ImgStorageServer).UploadImg(&imgStorageUploadImgServer{stream})
}

type ImgStorage_UploadImgServer interface {
	SendAndClose(*Null) error
	Recv() (*Img, error)
	grpc.ServerStream
}

type imgStorageUploadImgServer struct {
	grpc.ServerStream
}

func (x *imgStorageUploadImgServer) SendAndClose(m *Null) error {
	return x.ServerStream.SendMsg(m)
}

func (x *imgStorageUploadImgServer) Recv() (*Img, error) {
	m := new(Img)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ImgStorage_DownloadImg_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NameImg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ImgStorageServer).DownloadImg(m, &imgStorageDownloadImgServer{stream})
}

type ImgStorage_DownloadImgServer interface {
	Send(*Img) error
	grpc.ServerStream
}

type imgStorageDownloadImgServer struct {
	grpc.ServerStream
}

func (x *imgStorageDownloadImgServer) Send(m *Img) error {
	return x.ServerStream.SendMsg(m)
}

func _ImgStorage_GetListImg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImgStorageServer).GetListImg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imgstorage.ImgStorage/GetListImg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImgStorageServer).GetListImg(ctx, req.(*Null))
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
			MethodName: "GetListImg",
			Handler:    _ImgStorage_GetListImg_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadImg",
			Handler:       _ImgStorage_UploadImg_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadImg",
			Handler:       _ImgStorage_DownloadImg_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "imgstorage_service.proto",
}
