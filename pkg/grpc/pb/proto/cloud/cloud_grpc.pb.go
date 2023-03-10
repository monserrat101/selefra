// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: cloud/cloud.proto

package cloud

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

// CloudNoAuthClient is the client API for CloudNoAuth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CloudNoAuthClient interface {
	Login(ctx context.Context, in *Login_Request, opts ...grpc.CallOption) (*Login_Response, error)
}

type cloudNoAuthClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudNoAuthClient(cc grpc.ClientConnInterface) CloudNoAuthClient {
	return &cloudNoAuthClient{cc}
}

func (c *cloudNoAuthClient) Login(ctx context.Context, in *Login_Request, opts ...grpc.CallOption) (*Login_Response, error) {
	out := new(Login_Response)
	err := c.cc.Invoke(ctx, "/cloud.CloudNoAuth/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloudNoAuthServer is the server API for CloudNoAuth service.
// All implementations must embed UnimplementedCloudNoAuthServer
// for forward compatibility
type CloudNoAuthServer interface {
	Login(context.Context, *Login_Request) (*Login_Response, error)
	mustEmbedUnimplementedCloudNoAuthServer()
}

// UnimplementedCloudNoAuthServer must be embedded to have forward compatible implementations.
type UnimplementedCloudNoAuthServer struct {
}

func (UnimplementedCloudNoAuthServer) Login(context.Context, *Login_Request) (*Login_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedCloudNoAuthServer) mustEmbedUnimplementedCloudNoAuthServer() {}

// UnsafeCloudNoAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloudNoAuthServer will
// result in compilation errors.
type UnsafeCloudNoAuthServer interface {
	mustEmbedUnimplementedCloudNoAuthServer()
}

func RegisterCloudNoAuthServer(s grpc.ServiceRegistrar, srv CloudNoAuthServer) {
	s.RegisterService(&CloudNoAuth_ServiceDesc, srv)
}

func _CloudNoAuth_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Login_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudNoAuthServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.CloudNoAuth/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudNoAuthServer).Login(ctx, req.(*Login_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// CloudNoAuth_ServiceDesc is the grpc.ServiceDesc for CloudNoAuth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloudNoAuth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.CloudNoAuth",
	HandlerType: (*CloudNoAuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _CloudNoAuth_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cloud/cloud.proto",
}

// CloudClient is the client API for Cloud service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CloudClient interface {
	FetchOrgDsn(ctx context.Context, in *RequestEmpty, opts ...grpc.CallOption) (*FetchOrgDsn_Response, error)
	Logout(ctx context.Context, in *RequestEmpty, opts ...grpc.CallOption) (*Logout_Response, error)
	CreateProject(ctx context.Context, in *CreateProject_Request, opts ...grpc.CallOption) (*CreateProject_Response, error)
	SyncWorkplace(ctx context.Context, in *SyncWorkplace_Request, opts ...grpc.CallOption) (*SyncWorkplace_Response, error)
	CreateTask(ctx context.Context, in *CreateTask_Request, opts ...grpc.CallOption) (*CreateTask_Response, error)
}

type cloudClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudClient(cc grpc.ClientConnInterface) CloudClient {
	return &cloudClient{cc}
}

func (c *cloudClient) FetchOrgDsn(ctx context.Context, in *RequestEmpty, opts ...grpc.CallOption) (*FetchOrgDsn_Response, error) {
	out := new(FetchOrgDsn_Response)
	err := c.cc.Invoke(ctx, "/cloud.Cloud/FetchOrgDsn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudClient) Logout(ctx context.Context, in *RequestEmpty, opts ...grpc.CallOption) (*Logout_Response, error) {
	out := new(Logout_Response)
	err := c.cc.Invoke(ctx, "/cloud.Cloud/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudClient) CreateProject(ctx context.Context, in *CreateProject_Request, opts ...grpc.CallOption) (*CreateProject_Response, error) {
	out := new(CreateProject_Response)
	err := c.cc.Invoke(ctx, "/cloud.Cloud/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudClient) SyncWorkplace(ctx context.Context, in *SyncWorkplace_Request, opts ...grpc.CallOption) (*SyncWorkplace_Response, error) {
	out := new(SyncWorkplace_Response)
	err := c.cc.Invoke(ctx, "/cloud.Cloud/SyncWorkplace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudClient) CreateTask(ctx context.Context, in *CreateTask_Request, opts ...grpc.CallOption) (*CreateTask_Response, error) {
	out := new(CreateTask_Response)
	err := c.cc.Invoke(ctx, "/cloud.Cloud/CreateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloudServer is the server API for Cloud service.
// All implementations must embed UnimplementedCloudServer
// for forward compatibility
type CloudServer interface {
	FetchOrgDsn(context.Context, *RequestEmpty) (*FetchOrgDsn_Response, error)
	Logout(context.Context, *RequestEmpty) (*Logout_Response, error)
	CreateProject(context.Context, *CreateProject_Request) (*CreateProject_Response, error)
	SyncWorkplace(context.Context, *SyncWorkplace_Request) (*SyncWorkplace_Response, error)
	CreateTask(context.Context, *CreateTask_Request) (*CreateTask_Response, error)
	mustEmbedUnimplementedCloudServer()
}

// UnimplementedCloudServer must be embedded to have forward compatible implementations.
type UnimplementedCloudServer struct {
}

func (UnimplementedCloudServer) FetchOrgDsn(context.Context, *RequestEmpty) (*FetchOrgDsn_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchOrgDsn not implemented")
}
func (UnimplementedCloudServer) Logout(context.Context, *RequestEmpty) (*Logout_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedCloudServer) CreateProject(context.Context, *CreateProject_Request) (*CreateProject_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedCloudServer) SyncWorkplace(context.Context, *SyncWorkplace_Request) (*SyncWorkplace_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncWorkplace not implemented")
}
func (UnimplementedCloudServer) CreateTask(context.Context, *CreateTask_Request) (*CreateTask_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedCloudServer) mustEmbedUnimplementedCloudServer() {}

// UnsafeCloudServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloudServer will
// result in compilation errors.
type UnsafeCloudServer interface {
	mustEmbedUnimplementedCloudServer()
}

func RegisterCloudServer(s grpc.ServiceRegistrar, srv CloudServer) {
	s.RegisterService(&Cloud_ServiceDesc, srv)
}

func _Cloud_FetchOrgDsn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudServer).FetchOrgDsn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.Cloud/FetchOrgDsn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudServer).FetchOrgDsn(ctx, req.(*RequestEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cloud_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.Cloud/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudServer).Logout(ctx, req.(*RequestEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cloud_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProject_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.Cloud/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudServer).CreateProject(ctx, req.(*CreateProject_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cloud_SyncWorkplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncWorkplace_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudServer).SyncWorkplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.Cloud/SyncWorkplace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudServer).SyncWorkplace(ctx, req.(*SyncWorkplace_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cloud_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTask_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.Cloud/CreateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudServer).CreateTask(ctx, req.(*CreateTask_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Cloud_ServiceDesc is the grpc.ServiceDesc for Cloud service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cloud_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.Cloud",
	HandlerType: (*CloudServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchOrgDsn",
			Handler:    _Cloud_FetchOrgDsn_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Cloud_Logout_Handler,
		},
		{
			MethodName: "CreateProject",
			Handler:    _Cloud_CreateProject_Handler,
		},
		{
			MethodName: "SyncWorkplace",
			Handler:    _Cloud_SyncWorkplace_Handler,
		},
		{
			MethodName: "CreateTask",
			Handler:    _Cloud_CreateTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cloud/cloud.proto",
}
