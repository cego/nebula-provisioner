// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: agent-service.proto

package protocol

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

// AgentServiceClient is the client API for AgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentServiceClient interface {
	Enroll(ctx context.Context, in *EnrollRequest, opts ...grpc.CallOption) (*EnrollResponse, error)
	GetEnrollStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetEnrollStatusResponse, error)
	GetCertificateAuthorityByNetwork(ctx context.Context, in *GetCertificateAuthorityByNetworkRequest, opts ...grpc.CallOption) (*GetCertificateAuthorityByNetworkResponse, error)
	GetCRLByNetwork(ctx context.Context, in *GetCRLByNetworkRequest, opts ...grpc.CallOption) (*GetCRLByNetworkResponse, error)
}

type agentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentServiceClient(cc grpc.ClientConnInterface) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) Enroll(ctx context.Context, in *EnrollRequest, opts ...grpc.CallOption) (*EnrollResponse, error) {
	out := new(EnrollResponse)
	err := c.cc.Invoke(ctx, "/protocol.AgentService/Enroll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) GetEnrollStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetEnrollStatusResponse, error) {
	out := new(GetEnrollStatusResponse)
	err := c.cc.Invoke(ctx, "/protocol.AgentService/GetEnrollStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) GetCertificateAuthorityByNetwork(ctx context.Context, in *GetCertificateAuthorityByNetworkRequest, opts ...grpc.CallOption) (*GetCertificateAuthorityByNetworkResponse, error) {
	out := new(GetCertificateAuthorityByNetworkResponse)
	err := c.cc.Invoke(ctx, "/protocol.AgentService/GetCertificateAuthorityByNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) GetCRLByNetwork(ctx context.Context, in *GetCRLByNetworkRequest, opts ...grpc.CallOption) (*GetCRLByNetworkResponse, error) {
	out := new(GetCRLByNetworkResponse)
	err := c.cc.Invoke(ctx, "/protocol.AgentService/GetCRLByNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServiceServer is the server API for AgentService service.
// All implementations must embed UnimplementedAgentServiceServer
// for forward compatibility
type AgentServiceServer interface {
	Enroll(context.Context, *EnrollRequest) (*EnrollResponse, error)
	GetEnrollStatus(context.Context, *emptypb.Empty) (*GetEnrollStatusResponse, error)
	GetCertificateAuthorityByNetwork(context.Context, *GetCertificateAuthorityByNetworkRequest) (*GetCertificateAuthorityByNetworkResponse, error)
	GetCRLByNetwork(context.Context, *GetCRLByNetworkRequest) (*GetCRLByNetworkResponse, error)
	mustEmbedUnimplementedAgentServiceServer()
}

// UnimplementedAgentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (UnimplementedAgentServiceServer) Enroll(context.Context, *EnrollRequest) (*EnrollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enroll not implemented")
}
func (UnimplementedAgentServiceServer) GetEnrollStatus(context.Context, *emptypb.Empty) (*GetEnrollStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnrollStatus not implemented")
}
func (UnimplementedAgentServiceServer) GetCertificateAuthorityByNetwork(context.Context, *GetCertificateAuthorityByNetworkRequest) (*GetCertificateAuthorityByNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCertificateAuthorityByNetwork not implemented")
}
func (UnimplementedAgentServiceServer) GetCRLByNetwork(context.Context, *GetCRLByNetworkRequest) (*GetCRLByNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCRLByNetwork not implemented")
}
func (UnimplementedAgentServiceServer) mustEmbedUnimplementedAgentServiceServer() {}

// UnsafeAgentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServiceServer will
// result in compilation errors.
type UnsafeAgentServiceServer interface {
	mustEmbedUnimplementedAgentServiceServer()
}

func RegisterAgentServiceServer(s grpc.ServiceRegistrar, srv AgentServiceServer) {
	s.RegisterService(&AgentService_ServiceDesc, srv)
}

func _AgentService_Enroll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnrollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).Enroll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.AgentService/Enroll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).Enroll(ctx, req.(*EnrollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_GetEnrollStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).GetEnrollStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.AgentService/GetEnrollStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).GetEnrollStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_GetCertificateAuthorityByNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCertificateAuthorityByNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).GetCertificateAuthorityByNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.AgentService/GetCertificateAuthorityByNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).GetCertificateAuthorityByNetwork(ctx, req.(*GetCertificateAuthorityByNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_GetCRLByNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCRLByNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).GetCRLByNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.AgentService/GetCRLByNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).GetCRLByNetwork(ctx, req.(*GetCRLByNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AgentService_ServiceDesc is the grpc.ServiceDesc for AgentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enroll",
			Handler:    _AgentService_Enroll_Handler,
		},
		{
			MethodName: "GetEnrollStatus",
			Handler:    _AgentService_GetEnrollStatus_Handler,
		},
		{
			MethodName: "GetCertificateAuthorityByNetwork",
			Handler:    _AgentService_GetCertificateAuthorityByNetwork_Handler,
		},
		{
			MethodName: "GetCRLByNetwork",
			Handler:    _AgentService_GetCRLByNetwork_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent-service.proto",
}
