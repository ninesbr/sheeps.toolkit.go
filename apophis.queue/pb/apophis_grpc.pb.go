// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.24.4
// source: message.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PubSubService_About_FullMethodName          = "/pb.PubSubService/About"
	PubSubService_Ping_FullMethodName           = "/pb.PubSubService/Ping"
	PubSubService_Create_FullMethodName         = "/pb.PubSubService/Create"
	PubSubService_Drop_FullMethodName           = "/pb.PubSubService/Drop"
	PubSubService_Purge_FullMethodName          = "/pb.PubSubService/Purge"
	PubSubService_Info_FullMethodName           = "/pb.PubSubService/Info"
	PubSubService_Publish_FullMethodName        = "/pb.PubSubService/Publish"
	PubSubService_Subscribe_FullMethodName      = "/pb.PubSubService/Subscribe"
	PubSubService_MessageHistory_FullMethodName = "/pb.PubSubService/MessageHistory"
)

// PubSubServiceClient is the client API for PubSubService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PubSubServiceClient interface {
	About(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AboutResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Create(ctx context.Context, in *PubRequest, opts ...grpc.CallOption) (*PubResponse, error)
	Drop(ctx context.Context, in *DropRequest, opts ...grpc.CallOption) (*PubResponse, error)
	Purge(ctx context.Context, in *PurgeRequest, opts ...grpc.CallOption) (*PubResponse, error)
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*PubResponse, error)
	Publish(ctx context.Context, in *PubMessageRequest, opts ...grpc.CallOption) (*PubMessageResponse, error)
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[SubscribeMessage, SubscribeMessage], error)
	MessageHistory(ctx context.Context, in *MessageHistoryRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[MessageHistoryResponse], error)
}

type pubSubServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPubSubServiceClient(cc grpc.ClientConnInterface) PubSubServiceClient {
	return &pubSubServiceClient{cc}
}

func (c *pubSubServiceClient) About(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AboutResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AboutResponse)
	err := c.cc.Invoke(ctx, PubSubService_About_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, PubSubService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Create(ctx context.Context, in *PubRequest, opts ...grpc.CallOption) (*PubResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PubResponse)
	err := c.cc.Invoke(ctx, PubSubService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Drop(ctx context.Context, in *DropRequest, opts ...grpc.CallOption) (*PubResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PubResponse)
	err := c.cc.Invoke(ctx, PubSubService_Drop_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Purge(ctx context.Context, in *PurgeRequest, opts ...grpc.CallOption) (*PubResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PubResponse)
	err := c.cc.Invoke(ctx, PubSubService_Purge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*PubResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PubResponse)
	err := c.cc.Invoke(ctx, PubSubService_Info_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Publish(ctx context.Context, in *PubMessageRequest, opts ...grpc.CallOption) (*PubMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PubMessageResponse)
	err := c.cc.Invoke(ctx, PubSubService_Publish_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubServiceClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[SubscribeMessage, SubscribeMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PubSubService_ServiceDesc.Streams[0], PubSubService_Subscribe_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeMessage, SubscribeMessage]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PubSubService_SubscribeClient = grpc.BidiStreamingClient[SubscribeMessage, SubscribeMessage]

func (c *pubSubServiceClient) MessageHistory(ctx context.Context, in *MessageHistoryRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[MessageHistoryResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PubSubService_ServiceDesc.Streams[1], PubSubService_MessageHistory_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[MessageHistoryRequest, MessageHistoryResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PubSubService_MessageHistoryClient = grpc.ServerStreamingClient[MessageHistoryResponse]

// PubSubServiceServer is the server API for PubSubService service.
// All implementations must embed UnimplementedPubSubServiceServer
// for forward compatibility.
type PubSubServiceServer interface {
	About(context.Context, *emptypb.Empty) (*AboutResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Create(context.Context, *PubRequest) (*PubResponse, error)
	Drop(context.Context, *DropRequest) (*PubResponse, error)
	Purge(context.Context, *PurgeRequest) (*PubResponse, error)
	Info(context.Context, *InfoRequest) (*PubResponse, error)
	Publish(context.Context, *PubMessageRequest) (*PubMessageResponse, error)
	Subscribe(grpc.BidiStreamingServer[SubscribeMessage, SubscribeMessage]) error
	MessageHistory(*MessageHistoryRequest, grpc.ServerStreamingServer[MessageHistoryResponse]) error
	mustEmbedUnimplementedPubSubServiceServer()
}

// UnimplementedPubSubServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPubSubServiceServer struct{}

func (UnimplementedPubSubServiceServer) About(context.Context, *emptypb.Empty) (*AboutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method About not implemented")
}
func (UnimplementedPubSubServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedPubSubServiceServer) Create(context.Context, *PubRequest) (*PubResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPubSubServiceServer) Drop(context.Context, *DropRequest) (*PubResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Drop not implemented")
}
func (UnimplementedPubSubServiceServer) Purge(context.Context, *PurgeRequest) (*PubResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Purge not implemented")
}
func (UnimplementedPubSubServiceServer) Info(context.Context, *InfoRequest) (*PubResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedPubSubServiceServer) Publish(context.Context, *PubMessageRequest) (*PubMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedPubSubServiceServer) Subscribe(grpc.BidiStreamingServer[SubscribeMessage, SubscribeMessage]) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedPubSubServiceServer) MessageHistory(*MessageHistoryRequest, grpc.ServerStreamingServer[MessageHistoryResponse]) error {
	return status.Errorf(codes.Unimplemented, "method MessageHistory not implemented")
}
func (UnimplementedPubSubServiceServer) mustEmbedUnimplementedPubSubServiceServer() {}
func (UnimplementedPubSubServiceServer) testEmbeddedByValue()                       {}

// UnsafePubSubServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PubSubServiceServer will
// result in compilation errors.
type UnsafePubSubServiceServer interface {
	mustEmbedUnimplementedPubSubServiceServer()
}

func RegisterPubSubServiceServer(s grpc.ServiceRegistrar, srv PubSubServiceServer) {
	// If the following call pancis, it indicates UnimplementedPubSubServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PubSubService_ServiceDesc, srv)
}

func _PubSubService_About_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).About(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_About_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).About(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PubRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).Create(ctx, req.(*PubRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Drop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).Drop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_Drop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).Drop(ctx, req.(*DropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Purge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PurgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).Purge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_Purge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).Purge(ctx, req.(*PurgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_Info_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PubMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubSubService_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServiceServer).Publish(ctx, req.(*PubMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSubService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PubSubServiceServer).Subscribe(&grpc.GenericServerStream[SubscribeMessage, SubscribeMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PubSubService_SubscribeServer = grpc.BidiStreamingServer[SubscribeMessage, SubscribeMessage]

func _PubSubService_MessageHistory_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MessageHistoryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PubSubServiceServer).MessageHistory(m, &grpc.GenericServerStream[MessageHistoryRequest, MessageHistoryResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PubSubService_MessageHistoryServer = grpc.ServerStreamingServer[MessageHistoryResponse]

// PubSubService_ServiceDesc is the grpc.ServiceDesc for PubSubService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PubSubService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PubSubService",
	HandlerType: (*PubSubServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "About",
			Handler:    _PubSubService_About_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _PubSubService_Ping_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _PubSubService_Create_Handler,
		},
		{
			MethodName: "Drop",
			Handler:    _PubSubService_Drop_Handler,
		},
		{
			MethodName: "Purge",
			Handler:    _PubSubService_Purge_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _PubSubService_Info_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _PubSubService_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _PubSubService_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "MessageHistory",
			Handler:       _PubSubService_MessageHistory_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "message.proto",
}
