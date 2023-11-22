// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: server-command.proto

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

// ServerCommandClient is the client API for ServerCommand service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerCommandClient interface {
	IsInit(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*IsInitResponse, error)
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error)
	Unseal(ctx context.Context, in *UnsealRequest, opts ...grpc.CallOption) (*UnsealResponse, error)
	CreateNetwork(ctx context.Context, in *CreateNetworkRequest, opts ...grpc.CallOption) (*CreateNetworkResponse, error)
	ListNetwork(ctx context.Context, in *ListNetworkRequest, opts ...grpc.CallOption) (*ListNetworkResponse, error)
	ListCertificateAuthorityByNetwork(ctx context.Context, in *ListCertificateAuthorityByNetworkRequest, opts ...grpc.CallOption) (*ListCertificateAuthorityByNetworkResponse, error)
	GetEnrollmentTokenForNetwork(ctx context.Context, in *GetEnrollmentTokenForNetworkRequest, opts ...grpc.CallOption) (*GetEnrollmentTokenForNetworkResponse, error)
	ListEnrollmentRequests(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListEnrollmentRequestsResponse, error)
	ApproveEnrollmentRequest(ctx context.Context, in *ApproveEnrollmentRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListUsersWaitingForApproval(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersResponse, error)
	ApproveUserAccess(ctx context.Context, in *ApproveUserAccessRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type serverCommandClient struct {
	cc grpc.ClientConnInterface
}

func NewServerCommandClient(cc grpc.ClientConnInterface) ServerCommandClient {
	return &serverCommandClient{cc}
}

func (c *serverCommandClient) IsInit(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*IsInitResponse, error) {
	out := new(IsInitResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/IsInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	out := new(InitResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) Unseal(ctx context.Context, in *UnsealRequest, opts ...grpc.CallOption) (*UnsealResponse, error) {
	out := new(UnsealResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/Unseal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) CreateNetwork(ctx context.Context, in *CreateNetworkRequest, opts ...grpc.CallOption) (*CreateNetworkResponse, error) {
	out := new(CreateNetworkResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/CreateNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) ListNetwork(ctx context.Context, in *ListNetworkRequest, opts ...grpc.CallOption) (*ListNetworkResponse, error) {
	out := new(ListNetworkResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/ListNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) ListCertificateAuthorityByNetwork(ctx context.Context, in *ListCertificateAuthorityByNetworkRequest, opts ...grpc.CallOption) (*ListCertificateAuthorityByNetworkResponse, error) {
	out := new(ListCertificateAuthorityByNetworkResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/ListCertificateAuthorityByNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) GetEnrollmentTokenForNetwork(ctx context.Context, in *GetEnrollmentTokenForNetworkRequest, opts ...grpc.CallOption) (*GetEnrollmentTokenForNetworkResponse, error) {
	out := new(GetEnrollmentTokenForNetworkResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/GetEnrollmentTokenForNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) ListEnrollmentRequests(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListEnrollmentRequestsResponse, error) {
	out := new(ListEnrollmentRequestsResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/ListEnrollmentRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) ApproveEnrollmentRequest(ctx context.Context, in *ApproveEnrollmentRequestRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/ApproveEnrollmentRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) ListUsersWaitingForApproval(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/ListUsersWaitingForApproval", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverCommandClient) ApproveUserAccess(ctx context.Context, in *ApproveUserAccessRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/protocol.ServerCommand/ApproveUserAccess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerCommandServer is the server API for ServerCommand service.
// All implementations must embed UnimplementedServerCommandServer
// for forward compatibility
type ServerCommandServer interface {
	IsInit(context.Context, *emptypb.Empty) (*IsInitResponse, error)
	Init(context.Context, *InitRequest) (*InitResponse, error)
	Unseal(context.Context, *UnsealRequest) (*UnsealResponse, error)
	CreateNetwork(context.Context, *CreateNetworkRequest) (*CreateNetworkResponse, error)
	ListNetwork(context.Context, *ListNetworkRequest) (*ListNetworkResponse, error)
	ListCertificateAuthorityByNetwork(context.Context, *ListCertificateAuthorityByNetworkRequest) (*ListCertificateAuthorityByNetworkResponse, error)
	GetEnrollmentTokenForNetwork(context.Context, *GetEnrollmentTokenForNetworkRequest) (*GetEnrollmentTokenForNetworkResponse, error)
	ListEnrollmentRequests(context.Context, *emptypb.Empty) (*ListEnrollmentRequestsResponse, error)
	ApproveEnrollmentRequest(context.Context, *ApproveEnrollmentRequestRequest) (*emptypb.Empty, error)
	ListUsersWaitingForApproval(context.Context, *emptypb.Empty) (*ListUsersResponse, error)
	ApproveUserAccess(context.Context, *ApproveUserAccessRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedServerCommandServer()
}

// UnimplementedServerCommandServer must be embedded to have forward compatible implementations.
type UnimplementedServerCommandServer struct {
}

func (UnimplementedServerCommandServer) IsInit(context.Context, *emptypb.Empty) (*IsInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsInit not implemented")
}
func (UnimplementedServerCommandServer) Init(context.Context, *InitRequest) (*InitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedServerCommandServer) Unseal(context.Context, *UnsealRequest) (*UnsealResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unseal not implemented")
}
func (UnimplementedServerCommandServer) CreateNetwork(context.Context, *CreateNetworkRequest) (*CreateNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNetwork not implemented")
}
func (UnimplementedServerCommandServer) ListNetwork(context.Context, *ListNetworkRequest) (*ListNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNetwork not implemented")
}
func (UnimplementedServerCommandServer) ListCertificateAuthorityByNetwork(context.Context, *ListCertificateAuthorityByNetworkRequest) (*ListCertificateAuthorityByNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCertificateAuthorityByNetwork not implemented")
}
func (UnimplementedServerCommandServer) GetEnrollmentTokenForNetwork(context.Context, *GetEnrollmentTokenForNetworkRequest) (*GetEnrollmentTokenForNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnrollmentTokenForNetwork not implemented")
}
func (UnimplementedServerCommandServer) ListEnrollmentRequests(context.Context, *emptypb.Empty) (*ListEnrollmentRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEnrollmentRequests not implemented")
}
func (UnimplementedServerCommandServer) ApproveEnrollmentRequest(context.Context, *ApproveEnrollmentRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveEnrollmentRequest not implemented")
}
func (UnimplementedServerCommandServer) ListUsersWaitingForApproval(context.Context, *emptypb.Empty) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsersWaitingForApproval not implemented")
}
func (UnimplementedServerCommandServer) ApproveUserAccess(context.Context, *ApproveUserAccessRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveUserAccess not implemented")
}
func (UnimplementedServerCommandServer) mustEmbedUnimplementedServerCommandServer() {}

// UnsafeServerCommandServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerCommandServer will
// result in compilation errors.
type UnsafeServerCommandServer interface {
	mustEmbedUnimplementedServerCommandServer()
}

func RegisterServerCommandServer(s grpc.ServiceRegistrar, srv ServerCommandServer) {
	s.RegisterService(&ServerCommand_ServiceDesc, srv)
}

func _ServerCommand_IsInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).IsInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/IsInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).IsInit(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_Unseal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsealRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).Unseal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/Unseal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).Unseal(ctx, req.(*UnsealRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_CreateNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).CreateNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/CreateNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).CreateNetwork(ctx, req.(*CreateNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_ListNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).ListNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/ListNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).ListNetwork(ctx, req.(*ListNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_ListCertificateAuthorityByNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCertificateAuthorityByNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).ListCertificateAuthorityByNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/ListCertificateAuthorityByNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).ListCertificateAuthorityByNetwork(ctx, req.(*ListCertificateAuthorityByNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_GetEnrollmentTokenForNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEnrollmentTokenForNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).GetEnrollmentTokenForNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/GetEnrollmentTokenForNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).GetEnrollmentTokenForNetwork(ctx, req.(*GetEnrollmentTokenForNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_ListEnrollmentRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).ListEnrollmentRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/ListEnrollmentRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).ListEnrollmentRequests(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_ApproveEnrollmentRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveEnrollmentRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).ApproveEnrollmentRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/ApproveEnrollmentRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).ApproveEnrollmentRequest(ctx, req.(*ApproveEnrollmentRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_ListUsersWaitingForApproval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).ListUsersWaitingForApproval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/ListUsersWaitingForApproval",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).ListUsersWaitingForApproval(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerCommand_ApproveUserAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveUserAccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerCommandServer).ApproveUserAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.ServerCommand/ApproveUserAccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerCommandServer).ApproveUserAccess(ctx, req.(*ApproveUserAccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServerCommand_ServiceDesc is the grpc.ServiceDesc for ServerCommand service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServerCommand_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.ServerCommand",
	HandlerType: (*ServerCommandServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsInit",
			Handler:    _ServerCommand_IsInit_Handler,
		},
		{
			MethodName: "Init",
			Handler:    _ServerCommand_Init_Handler,
		},
		{
			MethodName: "Unseal",
			Handler:    _ServerCommand_Unseal_Handler,
		},
		{
			MethodName: "CreateNetwork",
			Handler:    _ServerCommand_CreateNetwork_Handler,
		},
		{
			MethodName: "ListNetwork",
			Handler:    _ServerCommand_ListNetwork_Handler,
		},
		{
			MethodName: "ListCertificateAuthorityByNetwork",
			Handler:    _ServerCommand_ListCertificateAuthorityByNetwork_Handler,
		},
		{
			MethodName: "GetEnrollmentTokenForNetwork",
			Handler:    _ServerCommand_GetEnrollmentTokenForNetwork_Handler,
		},
		{
			MethodName: "ListEnrollmentRequests",
			Handler:    _ServerCommand_ListEnrollmentRequests_Handler,
		},
		{
			MethodName: "ApproveEnrollmentRequest",
			Handler:    _ServerCommand_ApproveEnrollmentRequest_Handler,
		},
		{
			MethodName: "ListUsersWaitingForApproval",
			Handler:    _ServerCommand_ListUsersWaitingForApproval_Handler,
		},
		{
			MethodName: "ApproveUserAccess",
			Handler:    _ServerCommand_ApproveUserAccess_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server-command.proto",
}
