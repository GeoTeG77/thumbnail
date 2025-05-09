// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.1
// source: thumbnail.proto

package proto

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
	ThumbnailService_GetThumbnail_FullMethodName  = "/proto.ThumbnailService/GetThumbnail"
	ThumbnailService_GetThumbnails_FullMethodName = "/proto.ThumbnailService/GetThumbnails"
)

// ThumbnailServiceClient is the client API for ThumbnailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ThumbnailServiceClient interface {
	GetThumbnail(ctx context.Context, in *ThumbnailRequest, opts ...grpc.CallOption) (*ThumbnailResponse, error)
	GetThumbnails(ctx context.Context, in *ThumbnailsRequest, opts ...grpc.CallOption) (*ThumbnailsResponse, error)
}

type thumbnailServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewThumbnailServiceClient(cc grpc.ClientConnInterface) ThumbnailServiceClient {
	return &thumbnailServiceClient{cc}
}

func (c *thumbnailServiceClient) GetThumbnail(ctx context.Context, in *ThumbnailRequest, opts ...grpc.CallOption) (*ThumbnailResponse, error) {
	out := new(ThumbnailResponse)
	err := c.cc.Invoke(ctx, ThumbnailService_GetThumbnail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thumbnailServiceClient) GetThumbnails(ctx context.Context, in *ThumbnailsRequest, opts ...grpc.CallOption) (*ThumbnailsResponse, error) {
	out := new(ThumbnailsResponse)
	err := c.cc.Invoke(ctx, ThumbnailService_GetThumbnails_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ThumbnailServiceServer is the server API for ThumbnailService service.
// All implementations must embed UnimplementedThumbnailServiceServer
// for forward compatibility
type ThumbnailServiceServer interface {
	GetThumbnail(context.Context, *ThumbnailRequest) (*ThumbnailResponse, error)
	GetThumbnails(context.Context, *ThumbnailsRequest) (*ThumbnailsResponse, error)
	mustEmbedUnimplementedThumbnailServiceServer()
}

// UnimplementedThumbnailServiceServer must be embedded to have forward compatible implementations.
type UnimplementedThumbnailServiceServer struct {
}

func (UnimplementedThumbnailServiceServer) GetThumbnail(context.Context, *ThumbnailRequest) (*ThumbnailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetThumbnail not implemented")
}
func (UnimplementedThumbnailServiceServer) GetThumbnails(context.Context, *ThumbnailsRequest) (*ThumbnailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetThumbnails not implemented")
}
func (UnimplementedThumbnailServiceServer) mustEmbedUnimplementedThumbnailServiceServer() {}

// UnsafeThumbnailServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ThumbnailServiceServer will
// result in compilation errors.
type UnsafeThumbnailServiceServer interface {
	mustEmbedUnimplementedThumbnailServiceServer()
}

func RegisterThumbnailServiceServer(s grpc.ServiceRegistrar, srv ThumbnailServiceServer) {
	s.RegisterService(&ThumbnailService_ServiceDesc, srv)
}

func _ThumbnailService_GetThumbnail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThumbnailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThumbnailServiceServer).GetThumbnail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ThumbnailService_GetThumbnail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThumbnailServiceServer).GetThumbnail(ctx, req.(*ThumbnailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ThumbnailService_GetThumbnails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThumbnailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThumbnailServiceServer).GetThumbnails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ThumbnailService_GetThumbnails_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThumbnailServiceServer).GetThumbnails(ctx, req.(*ThumbnailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ThumbnailService_ServiceDesc is the grpc.ServiceDesc for ThumbnailService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ThumbnailService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ThumbnailService",
	HandlerType: (*ThumbnailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetThumbnail",
			Handler:    _ThumbnailService_GetThumbnail_Handler,
		},
		{
			MethodName: "GetThumbnails",
			Handler:    _ThumbnailService_GetThumbnails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "thumbnail.proto",
}
