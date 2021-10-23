// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package MiniProject2

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

// ChittyChatClient is the client API for ChittyChat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatClient interface {
	CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (ChittyChat_CreateStreamClient, error)
	BroadcastMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Empty, error)
}

type chittyChatClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatClient(cc grpc.ClientConnInterface) ChittyChatClient {
	return &chittyChatClient{cc}
}

func (c *chittyChatClient) CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (ChittyChat_CreateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChittyChat_ServiceDesc.Streams[0], "/proto.ChittyChat/CreateStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &chittyChatCreateStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChittyChat_CreateStreamClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chittyChatCreateStreamClient struct {
	grpc.ClientStream
}

func (x *chittyChatCreateStreamClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chittyChatClient) BroadcastMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.ChittyChat/BroadcastMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServer is the server API for ChittyChat service.
// All implementations must embed UnimplementedChittyChatServer
// for forward compatibility
type ChittyChatServer interface {
	CreateStream(*Connect, ChittyChat_CreateStreamServer) error
	BroadcastMessage(context.Context, *ChatMessage) (*Empty, error)
	mustEmbedUnimplementedChittyChatServer()
}

// UnimplementedChittyChatServer must be embedded to have forward compatible implementations.
type UnimplementedChittyChatServer struct {
}

func (UnimplementedChittyChatServer) CreateStream(*Connect, ChittyChat_CreateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}
func (UnimplementedChittyChatServer) BroadcastMessage(context.Context, *ChatMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastMessage not implemented")
}
func (UnimplementedChittyChatServer) mustEmbedUnimplementedChittyChatServer() {}

// UnsafeChittyChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServer will
// result in compilation errors.
type UnsafeChittyChatServer interface {
	mustEmbedUnimplementedChittyChatServer()
}

func RegisterChittyChatServer(s grpc.ServiceRegistrar, srv ChittyChatServer) {
	s.RegisterService(&ChittyChat_ServiceDesc, srv)
}

func _ChittyChat_CreateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Connect)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChittyChatServer).CreateStream(m, &chittyChatCreateStreamServer{stream})
}

type ChittyChat_CreateStreamServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type chittyChatCreateStreamServer struct {
	grpc.ServerStream
}

func (x *chittyChatCreateStreamServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _ChittyChat_BroadcastMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).BroadcastMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ChittyChat/BroadcastMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).BroadcastMessage(ctx, req.(*ChatMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChat_ServiceDesc is the grpc.ServiceDesc for ChittyChat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ChittyChat",
	HandlerType: (*ChittyChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BroadcastMessage",
			Handler:    _ChittyChat_BroadcastMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateStream",
			Handler:       _ChittyChat_CreateStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/chittychat.proto",
}
