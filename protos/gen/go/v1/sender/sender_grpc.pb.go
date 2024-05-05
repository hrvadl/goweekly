// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: v1/sender/sender.proto

package sender

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

// SenderServiceClient is the client API for SenderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SenderServiceClient interface {
	Send(ctx context.Context, opts ...grpc.CallOption) (SenderService_SendClient, error)
}

type senderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSenderServiceClient(cc grpc.ClientConnInterface) SenderServiceClient {
	return &senderServiceClient{cc}
}

func (c *senderServiceClient) Send(ctx context.Context, opts ...grpc.CallOption) (SenderService_SendClient, error) {
	stream, err := c.cc.NewStream(ctx, &SenderService_ServiceDesc.Streams[0], "/sender.v1.SenderService/Send", opts...)
	if err != nil {
		return nil, err
	}
	x := &senderServiceSendClient{stream}
	return x, nil
}

type SenderService_SendClient interface {
	Send(*SendRequest) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type senderServiceSendClient struct {
	grpc.ClientStream
}

func (x *senderServiceSendClient) Send(m *SendRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *senderServiceSendClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SenderServiceServer is the server API for SenderService service.
// All implementations must embed UnimplementedSenderServiceServer
// for forward compatibility
type SenderServiceServer interface {
	Send(SenderService_SendServer) error
	mustEmbedUnimplementedSenderServiceServer()
}

// UnimplementedSenderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSenderServiceServer struct {
}

func (UnimplementedSenderServiceServer) Send(SenderService_SendServer) error {
	return status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedSenderServiceServer) mustEmbedUnimplementedSenderServiceServer() {}

// UnsafeSenderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SenderServiceServer will
// result in compilation errors.
type UnsafeSenderServiceServer interface {
	mustEmbedUnimplementedSenderServiceServer()
}

func RegisterSenderServiceServer(s grpc.ServiceRegistrar, srv SenderServiceServer) {
	s.RegisterService(&SenderService_ServiceDesc, srv)
}

func _SenderService_Send_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SenderServiceServer).Send(&senderServiceSendServer{stream})
}

type SenderService_SendServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*SendRequest, error)
	grpc.ServerStream
}

type senderServiceSendServer struct {
	grpc.ServerStream
}

func (x *senderServiceSendServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *senderServiceSendServer) Recv() (*SendRequest, error) {
	m := new(SendRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SenderService_ServiceDesc is the grpc.ServiceDesc for SenderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SenderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sender.v1.SenderService",
	HandlerType: (*SenderServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Send",
			Handler:       _SenderService_Send_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "v1/sender/sender.proto",
}
