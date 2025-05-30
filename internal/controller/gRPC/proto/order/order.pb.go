// internal/controller/gRPC/proto/order/order.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.11.0
// source: internal/controller/gRPC/proto/order/order.proto

package order

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 请求消息
type QueryOrdersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartDate string `protobuf:"bytes,1,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty" validate:"required,datetime=2006-01-02"` // 开始日期，格式：YYYY-MM-DD
	EndDate   string `protobuf:"bytes,2,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty" validate:"required,datetime=2006-01-02"`       // 结束日期，格式：YYYY-MM-DD
	OrderId   string `protobuf:"bytes,3,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty" validate:"omitempty,uuid"`       // 订单ID，可选
	Page      int32  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty" validate:"required,min=1"`                           // 页码，从1开始
	PageSize  int32  `protobuf:"varint,5,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty" validate:"required,min=1"`   // 每页数量
}

func (x *QueryOrdersRequest) Reset() {
	*x = QueryOrdersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryOrdersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryOrdersRequest) ProtoMessage() {}

func (x *QueryOrdersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryOrdersRequest.ProtoReflect.Descriptor instead.
func (*QueryOrdersRequest) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{0}
}

func (x *QueryOrdersRequest) GetStartDate() string {
	if x != nil {
		return x.StartDate
	}
	return ""
}

func (x *QueryOrdersRequest) GetEndDate() string {
	if x != nil {
		return x.EndDate
	}
	return ""
}

func (x *QueryOrdersRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *QueryOrdersRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *QueryOrdersRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

// 获取单个订单详情请求
type GetOrderDetailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"` // 订单ID
}

func (x *GetOrderDetailRequest) Reset() {
	*x = GetOrderDetailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOrderDetailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderDetailRequest) ProtoMessage() {}

func (x *GetOrderDetailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderDetailRequest.ProtoReflect.Descriptor instead.
func (*GetOrderDetailRequest) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{1}
}

func (x *GetOrderDetailRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

// 获取单个订单详情响应
type GetOrderDetailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Order *Order `protobuf:"bytes,1,opt,name=order,proto3" json:"order,omitempty"` // 订单详情
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"` // 错误信息，如果为空则表示成功
}

func (x *GetOrderDetailResponse) Reset() {
	*x = GetOrderDetailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOrderDetailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderDetailResponse) ProtoMessage() {}

func (x *GetOrderDetailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderDetailResponse.ProtoReflect.Descriptor instead.
func (*GetOrderDetailResponse) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{2}
}

func (x *GetOrderDetailResponse) GetOrder() *Order {
	if x != nil {
		return x.Order
	}
	return nil
}

func (x *GetOrderDetailResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// 响应消息
type QueryOrdersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Orders   []*OrderDetail `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`                      // 订单列表
	Total    int32          `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`                       // 总记录数
	Page     int32          `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`                         // 当前页码
	PageSize int32          `protobuf:"varint,4,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"` // 每页数量
}

func (x *QueryOrdersResponse) Reset() {
	*x = QueryOrdersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryOrdersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryOrdersResponse) ProtoMessage() {}

func (x *QueryOrdersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryOrdersResponse.ProtoReflect.Descriptor instead.
func (*QueryOrdersResponse) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{3}
}

func (x *QueryOrdersResponse) GetOrders() []*OrderDetail {
	if x != nil {
		return x.Orders
	}
	return nil
}

func (x *QueryOrdersResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *QueryOrdersResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *QueryOrdersResponse) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

// 订单信息
type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId    string       `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`          // 订单ID
	UserId     string       `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`             // 用户ID
	Amount     float64      `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`                         // 订单金额
	Status     string       `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`                           // 订单状态
	CreateTime string       `protobuf:"bytes,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"` // 创建时间
	UpdateTime string       `protobuf:"bytes,6,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"` // 更新时间
	Items      []*OrderItem `protobuf:"bytes,7,rep,name=items,proto3" json:"items,omitempty"`                             // 订单项
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{4}
}

func (x *Order) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *Order) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Order) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Order) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Order) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *Order) GetUpdateTime() string {
	if x != nil {
		return x.UpdateTime
	}
	return ""
}

func (x *Order) GetItems() []*OrderItem {
	if x != nil {
		return x.Items
	}
	return nil
}

// 订单项
type OrderItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemId   string  `protobuf:"bytes,1,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`       // 商品ID
	ItemName string  `protobuf:"bytes,2,opt,name=item_name,json=itemName,proto3" json:"item_name,omitempty"` // 商品名称
	Quantity int32   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`                // 数量
	Price    float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`                     // 单价
}

func (x *OrderItem) Reset() {
	*x = OrderItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderItem) ProtoMessage() {}

func (x *OrderItem) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderItem.ProtoReflect.Descriptor instead.
func (*OrderItem) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{5}
}

func (x *OrderItem) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *OrderItem) GetItemName() string {
	if x != nil {
		return x.ItemName
	}
	return ""
}

func (x *OrderItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *OrderItem) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

// 订单详情
type OrderDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId    string       `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`          // 订单ID
	UserId     string       `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`             // 用户ID
	Amount     float64      `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`                         // 订单金额
	Status     string       `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`                           // 订单状态
	CreateTime string       `protobuf:"bytes,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"` // 创建时间
	UpdateTime string       `protobuf:"bytes,6,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"` // 更新时间
	Items      []*OrderItem `protobuf:"bytes,7,rep,name=items,proto3" json:"items,omitempty"`                             // 订单项, repeated定义数组
}

func (x *OrderDetail) Reset() {
	*x = OrderDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderDetail) ProtoMessage() {}

func (x *OrderDetail) ProtoReflect() protoreflect.Message {
	mi := &file_internal_controller_gRPC_proto_order_order_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderDetail.ProtoReflect.Descriptor instead.
func (*OrderDetail) Descriptor() ([]byte, []int) {
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP(), []int{6}
}

func (x *OrderDetail) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *OrderDetail) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *OrderDetail) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *OrderDetail) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *OrderDetail) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *OrderDetail) GetUpdateTime() string {
	if x != nil {
		return x.UpdateTime
	}
	return ""
}

func (x *OrderDetail) GetItems() []*OrderItem {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_internal_controller_gRPC_proto_order_order_proto protoreflect.FileDescriptor

var file_internal_controller_gRPC_proto_order_order_proto_rawDesc = []byte{
	0x0a, 0x30, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x67, 0x52, 0x50, 0x43, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x9a, 0x01, 0x0a, 0x12, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x32, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x16, 0x47, 0x65,
	0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x88,
	0x01, 0x0a, 0x13, 0x51, 0x75, 0x65, 0x72, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x06, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0xd5, 0x01, 0x0a, 0x05, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x73, 0x0a, 0x09, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x17,
	0x0a, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0xdb, 0x01, 0x0a, 0x0b, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x32, 0xa7, 0x01, 0x0a, 0x0c, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x73, 0x12, 0x19, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1a, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12,
	0x1c, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2d,
	0x5a, 0x2b, 0x54, 0x61, 0x75, 0x72, 0x75, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x67, 0x52, 0x50,
	0x43, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_controller_gRPC_proto_order_order_proto_rawDescOnce sync.Once
	file_internal_controller_gRPC_proto_order_order_proto_rawDescData = file_internal_controller_gRPC_proto_order_order_proto_rawDesc
)

func file_internal_controller_gRPC_proto_order_order_proto_rawDescGZIP() []byte {
	file_internal_controller_gRPC_proto_order_order_proto_rawDescOnce.Do(func() {
		file_internal_controller_gRPC_proto_order_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_controller_gRPC_proto_order_order_proto_rawDescData)
	})
	return file_internal_controller_gRPC_proto_order_order_proto_rawDescData
}

var file_internal_controller_gRPC_proto_order_order_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_controller_gRPC_proto_order_order_proto_goTypes = []interface{}{
	(*QueryOrdersRequest)(nil),     // 0: order.QueryOrdersRequest
	(*GetOrderDetailRequest)(nil),  // 1: order.GetOrderDetailRequest
	(*GetOrderDetailResponse)(nil), // 2: order.GetOrderDetailResponse
	(*QueryOrdersResponse)(nil),    // 3: order.QueryOrdersResponse
	(*Order)(nil),                  // 4: order.Order
	(*OrderItem)(nil),              // 5: order.OrderItem
	(*OrderDetail)(nil),            // 6: order.OrderDetail
}
var file_internal_controller_gRPC_proto_order_order_proto_depIdxs = []int32{
	4, // 0: order.GetOrderDetailResponse.order:type_name -> order.Order
	6, // 1: order.QueryOrdersResponse.orders:type_name -> order.OrderDetail
	5, // 2: order.Order.items:type_name -> order.OrderItem
	5, // 3: order.OrderDetail.items:type_name -> order.OrderItem
	0, // 4: order.OrderService.QueryOrders:input_type -> order.QueryOrdersRequest
	1, // 5: order.OrderService.GetOrderDetail:input_type -> order.GetOrderDetailRequest
	3, // 6: order.OrderService.QueryOrders:output_type -> order.QueryOrdersResponse
	2, // 7: order.OrderService.GetOrderDetail:output_type -> order.GetOrderDetailResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_internal_controller_gRPC_proto_order_order_proto_init() }
func file_internal_controller_gRPC_proto_order_order_proto_init() {
	if File_internal_controller_gRPC_proto_order_order_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryOrdersRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOrderDetailRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOrderDetailResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryOrdersResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Order); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_controller_gRPC_proto_order_order_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderDetail); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_controller_gRPC_proto_order_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_controller_gRPC_proto_order_order_proto_goTypes,
		DependencyIndexes: file_internal_controller_gRPC_proto_order_order_proto_depIdxs,
		MessageInfos:      file_internal_controller_gRPC_proto_order_order_proto_msgTypes,
	}.Build()
	File_internal_controller_gRPC_proto_order_order_proto = out.File
	file_internal_controller_gRPC_proto_order_order_proto_rawDesc = nil
	file_internal_controller_gRPC_proto_order_order_proto_goTypes = nil
	file_internal_controller_gRPC_proto_order_order_proto_depIdxs = nil
}
