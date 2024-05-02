// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package hdwallet

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

// HdWalletApiClient is the client API for HdWalletApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HdWalletApiClient interface {
	GenerateMnemonic(ctx context.Context, in *GenerateMnemonicRequest, opts ...grpc.CallOption) (*GenerateMnemonicResponse, error)
	EncryptMnemonic(ctx context.Context, in *EncryptMnemonicRequest, opts ...grpc.CallOption) (*EncryptMnemonicResponse, error)
	ValidateMnemonic(ctx context.Context, in *ValidateMnemonicRequest, opts ...grpc.CallOption) (*ValidateMnemonicResponse, error)
	LoadMnemonic(ctx context.Context, in *LoadMnemonicRequest, opts ...grpc.CallOption) (*LoadMnemonicResponse, error)
	UnLoadMnemonic(ctx context.Context, in *UnLoadMnemonicRequest, opts ...grpc.CallOption) (*UnLoadMnemonicResponse, error)
	UnLoadMultipleMnemonics(ctx context.Context, in *UnLoadMultipleMnemonicsRequest, opts ...grpc.CallOption) (*UnLoadMultipleMnemonicsResponse, error)
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	GetMultipleAccounts(ctx context.Context, in *GetMultipleAccountRequest, opts ...grpc.CallOption) (*GetMultipleAccountResponse, error)
	LoadAccount(ctx context.Context, in *LoadAccountRequest, opts ...grpc.CallOption) (*LoadAccountsResponse, error)
	SignData(ctx context.Context, in *SignDataRequest, opts ...grpc.CallOption) (*SignDataResponse, error)
}

type hdWalletApiClient struct {
	cc grpc.ClientConnInterface
}

func NewHdWalletApiClient(cc grpc.ClientConnInterface) HdWalletApiClient {
	return &hdWalletApiClient{cc}
}

func (c *hdWalletApiClient) GenerateMnemonic(ctx context.Context, in *GenerateMnemonicRequest, opts ...grpc.CallOption) (*GenerateMnemonicResponse, error) {
	out := new(GenerateMnemonicResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/GenerateMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) EncryptMnemonic(ctx context.Context, in *EncryptMnemonicRequest, opts ...grpc.CallOption) (*EncryptMnemonicResponse, error) {
	out := new(EncryptMnemonicResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/EncryptMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) ValidateMnemonic(ctx context.Context, in *ValidateMnemonicRequest, opts ...grpc.CallOption) (*ValidateMnemonicResponse, error) {
	out := new(ValidateMnemonicResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/ValidateMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) LoadMnemonic(ctx context.Context, in *LoadMnemonicRequest, opts ...grpc.CallOption) (*LoadMnemonicResponse, error) {
	out := new(LoadMnemonicResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/LoadMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) UnLoadMnemonic(ctx context.Context, in *UnLoadMnemonicRequest, opts ...grpc.CallOption) (*UnLoadMnemonicResponse, error) {
	out := new(UnLoadMnemonicResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/UnLoadMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) UnLoadMultipleMnemonics(ctx context.Context, in *UnLoadMultipleMnemonicsRequest, opts ...grpc.CallOption) (*UnLoadMultipleMnemonicsResponse, error) {
	out := new(UnLoadMultipleMnemonicsResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/UnLoadMultipleMnemonics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) GetMultipleAccounts(ctx context.Context, in *GetMultipleAccountRequest, opts ...grpc.CallOption) (*GetMultipleAccountResponse, error) {
	out := new(GetMultipleAccountResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/GetMultipleAccounts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) LoadAccount(ctx context.Context, in *LoadAccountRequest, opts ...grpc.CallOption) (*LoadAccountsResponse, error) {
	out := new(LoadAccountsResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/LoadAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hdWalletApiClient) SignData(ctx context.Context, in *SignDataRequest, opts ...grpc.CallOption) (*SignDataResponse, error) {
	out := new(SignDataResponse)
	err := c.cc.Invoke(ctx, "/hdwallet_api.HdWalletApi/SignData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HdWalletApiServer is the server API for HdWalletApi service.
// All implementations must embed UnimplementedHdWalletApiServer
// for forward compatibility
type HdWalletApiServer interface {
	GenerateMnemonic(context.Context, *GenerateMnemonicRequest) (*GenerateMnemonicResponse, error)
	EncryptMnemonic(context.Context, *EncryptMnemonicRequest) (*EncryptMnemonicResponse, error)
	ValidateMnemonic(context.Context, *ValidateMnemonicRequest) (*ValidateMnemonicResponse, error)
	LoadMnemonic(context.Context, *LoadMnemonicRequest) (*LoadMnemonicResponse, error)
	UnLoadMnemonic(context.Context, *UnLoadMnemonicRequest) (*UnLoadMnemonicResponse, error)
	UnLoadMultipleMnemonics(context.Context, *UnLoadMultipleMnemonicsRequest) (*UnLoadMultipleMnemonicsResponse, error)
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	GetMultipleAccounts(context.Context, *GetMultipleAccountRequest) (*GetMultipleAccountResponse, error)
	LoadAccount(context.Context, *LoadAccountRequest) (*LoadAccountsResponse, error)
	SignData(context.Context, *SignDataRequest) (*SignDataResponse, error)
	mustEmbedUnimplementedHdWalletApiServer()
}

// UnimplementedHdWalletApiServer must be embedded to have forward compatible implementations.
type UnimplementedHdWalletApiServer struct {
}

func (UnimplementedHdWalletApiServer) GenerateMnemonic(context.Context, *GenerateMnemonicRequest) (*GenerateMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateMnemonic not implemented")
}
func (UnimplementedHdWalletApiServer) EncryptMnemonic(context.Context, *EncryptMnemonicRequest) (*EncryptMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EncryptMnemonic not implemented")
}
func (UnimplementedHdWalletApiServer) ValidateMnemonic(context.Context, *ValidateMnemonicRequest) (*ValidateMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateMnemonic not implemented")
}
func (UnimplementedHdWalletApiServer) LoadMnemonic(context.Context, *LoadMnemonicRequest) (*LoadMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadMnemonic not implemented")
}
func (UnimplementedHdWalletApiServer) UnLoadMnemonic(context.Context, *UnLoadMnemonicRequest) (*UnLoadMnemonicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnLoadMnemonic not implemented")
}
func (UnimplementedHdWalletApiServer) UnLoadMultipleMnemonics(context.Context, *UnLoadMultipleMnemonicsRequest) (*UnLoadMultipleMnemonicsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnLoadMultipleMnemonics not implemented")
}
func (UnimplementedHdWalletApiServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedHdWalletApiServer) GetMultipleAccounts(context.Context, *GetMultipleAccountRequest) (*GetMultipleAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMultipleAccounts not implemented")
}
func (UnimplementedHdWalletApiServer) LoadAccount(context.Context, *LoadAccountRequest) (*LoadAccountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadAccount not implemented")
}
func (UnimplementedHdWalletApiServer) SignData(context.Context, *SignDataRequest) (*SignDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignData not implemented")
}
func (UnimplementedHdWalletApiServer) mustEmbedUnimplementedHdWalletApiServer() {}

// UnsafeHdWalletApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HdWalletApiServer will
// result in compilation errors.
type UnsafeHdWalletApiServer interface {
	mustEmbedUnimplementedHdWalletApiServer()
}

func RegisterHdWalletApiServer(s grpc.ServiceRegistrar, srv HdWalletApiServer) {
	s.RegisterService(&HdWalletApi_ServiceDesc, srv)
}

func _HdWalletApi_GenerateMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateMnemonicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).GenerateMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/GenerateMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).GenerateMnemonic(ctx, req.(*GenerateMnemonicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_EncryptMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncryptMnemonicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).EncryptMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/EncryptMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).EncryptMnemonic(ctx, req.(*EncryptMnemonicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_ValidateMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateMnemonicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).ValidateMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/ValidateMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).ValidateMnemonic(ctx, req.(*ValidateMnemonicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_LoadMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadMnemonicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).LoadMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/LoadMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).LoadMnemonic(ctx, req.(*LoadMnemonicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_UnLoadMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnLoadMnemonicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).UnLoadMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/UnLoadMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).UnLoadMnemonic(ctx, req.(*UnLoadMnemonicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_UnLoadMultipleMnemonics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnLoadMultipleMnemonicsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).UnLoadMultipleMnemonics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/UnLoadMultipleMnemonics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).UnLoadMultipleMnemonics(ctx, req.(*UnLoadMultipleMnemonicsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_GetMultipleAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMultipleAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).GetMultipleAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/GetMultipleAccounts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).GetMultipleAccounts(ctx, req.(*GetMultipleAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_LoadAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).LoadAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/LoadAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).LoadAccount(ctx, req.(*LoadAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HdWalletApi_SignData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HdWalletApiServer).SignData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hdwallet_api.HdWalletApi/SignData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HdWalletApiServer).SignData(ctx, req.(*SignDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HdWalletApi_ServiceDesc is the grpc.ServiceDesc for HdWalletApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HdWalletApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hdwallet_api.HdWalletApi",
	HandlerType: (*HdWalletApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateMnemonic",
			Handler:    _HdWalletApi_GenerateMnemonic_Handler,
		},
		{
			MethodName: "EncryptMnemonic",
			Handler:    _HdWalletApi_EncryptMnemonic_Handler,
		},
		{
			MethodName: "ValidateMnemonic",
			Handler:    _HdWalletApi_ValidateMnemonic_Handler,
		},
		{
			MethodName: "LoadMnemonic",
			Handler:    _HdWalletApi_LoadMnemonic_Handler,
		},
		{
			MethodName: "UnLoadMnemonic",
			Handler:    _HdWalletApi_UnLoadMnemonic_Handler,
		},
		{
			MethodName: "UnLoadMultipleMnemonics",
			Handler:    _HdWalletApi_UnLoadMultipleMnemonics_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _HdWalletApi_GetAccount_Handler,
		},
		{
			MethodName: "GetMultipleAccounts",
			Handler:    _HdWalletApi_GetMultipleAccounts_Handler,
		},
		{
			MethodName: "LoadAccount",
			Handler:    _HdWalletApi_LoadAccount_Handler,
		},
		{
			MethodName: "SignData",
			Handler:    _HdWalletApi_SignData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hdwallet_api.proto",
}
