// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.0
// source: internal/controller/gRPC/proto/order/order.proto

package order

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

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderServiceClient interface {
	// 根据日期区间和订单ID查询订单
	QueryOrders(ctx context.Context, in *QueryOrdersRequest, opts ...grpc.CallOption) (*QueryOrdersResponse, error)
	// 获取单个订单详情
	GetOrderDetail(ctx context.Context, in *GetOrderDetailRequest, opts ...grpc.CallOption) (*GetOrderDetailResponse, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) QueryOrders(ctx context.Context, in *QueryOrdersRequest, opts ...grpc.CallOption) (*QueryOrdersResponse, error) {
	out := new(QueryOrdersResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/QueryOrders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetOrderDetail(ctx context.Context, in *GetOrderDetailRequest, opts ...grpc.CallOption) (*GetOrderDetailResponse, error) {
	out := new(GetOrderDetailResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/GetOrderDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations must embed UnimplementedOrderServiceServer
// for forward compatibility
type OrderServiceServer interface {
	// 根据日期区间和订单ID查询订单
	QueryOrders(context.Context, *QueryOrdersRequest) (*QueryOrdersResponse, error)
	// 获取单个订单详情
	GetOrderDetail(context.Context, *GetOrderDetailRequest) (*GetOrderDetailResponse, error)
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServiceServer struct {
}

func (UnimplementedOrderServiceServer) QueryOrders(context.Context, *QueryOrdersRequest) (*QueryOrdersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryOrders not implemented")
}
func (UnimplementedOrderServiceServer) GetOrderDetail(context.Context, *GetOrderDetailRequest) (*GetOrderDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderDetail not implemented")
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_QueryOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).QueryOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/QueryOrders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).QueryOrders(ctx, req.(*QueryOrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetOrderDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrderDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetOrderDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/GetOrderDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetOrderDetail(ctx, req.(*GetOrderDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryOrders",
			Handler:    _OrderService_QueryOrders_Handler,
		},
		{
			MethodName: "GetOrderDetail",
			Handler:    _OrderService_GetOrderDetail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/controller/gRPC/proto/order/order.proto",
}
