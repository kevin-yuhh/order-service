// Code generated by protoc-gen-go. DO NOT EDIT.
// source: order.proto

package order

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type QueryBalanceRequest struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryBalanceRequest) Reset()         { *m = QueryBalanceRequest{} }
func (m *QueryBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBalanceRequest) ProtoMessage()    {}
func (*QueryBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{0}
}

func (m *QueryBalanceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryBalanceRequest.Unmarshal(m, b)
}
func (m *QueryBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryBalanceRequest.Marshal(b, m, deterministic)
}
func (m *QueryBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBalanceRequest.Merge(m, src)
}
func (m *QueryBalanceRequest) XXX_Size() int {
	return xxx_messageInfo_QueryBalanceRequest.Size(m)
}
func (m *QueryBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBalanceRequest proto.InternalMessageInfo

func (m *QueryBalanceRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryBalanceResponse struct {
	Balance              int64    `protobuf:"varint,3,opt,name=balance,proto3" json:"balance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryBalanceResponse) Reset()         { *m = QueryBalanceResponse{} }
func (m *QueryBalanceResponse) String() string { return proto.CompactTextString(m) }
func (*QueryBalanceResponse) ProtoMessage()    {}
func (*QueryBalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{1}
}

func (m *QueryBalanceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryBalanceResponse.Unmarshal(m, b)
}
func (m *QueryBalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryBalanceResponse.Marshal(b, m, deterministic)
}
func (m *QueryBalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBalanceResponse.Merge(m, src)
}
func (m *QueryBalanceResponse) XXX_Size() int {
	return xxx_messageInfo_QueryBalanceResponse.Size(m)
}
func (m *QueryBalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBalanceResponse proto.InternalMessageInfo

func (m *QueryBalanceResponse) GetBalance() int64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

type CreateOrderRequest struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	RequestId            string   `protobuf:"bytes,2,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	FileName             string   `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileSize             int64    `protobuf:"varint,4,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrderRequest) Reset()         { *m = CreateOrderRequest{} }
func (m *CreateOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrderRequest) ProtoMessage()    {}
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{2}
}

func (m *CreateOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateOrderRequest.Unmarshal(m, b)
}
func (m *CreateOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateOrderRequest.Marshal(b, m, deterministic)
}
func (m *CreateOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateOrderRequest.Merge(m, src)
}
func (m *CreateOrderRequest) XXX_Size() int {
	return xxx_messageInfo_CreateOrderRequest.Size(m)
}
func (m *CreateOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateOrderRequest proto.InternalMessageInfo

func (m *CreateOrderRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *CreateOrderRequest) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *CreateOrderRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *CreateOrderRequest) GetFileSize() int64 {
	if m != nil {
		return m.FileSize
	}
	return 0
}

type CreateOrderResponse struct {
	OrderId              int64    `protobuf:"varint,3,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrderResponse) Reset()         { *m = CreateOrderResponse{} }
func (m *CreateOrderResponse) String() string { return proto.CompactTextString(m) }
func (*CreateOrderResponse) ProtoMessage()    {}
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{3}
}

func (m *CreateOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateOrderResponse.Unmarshal(m, b)
}
func (m *CreateOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateOrderResponse.Marshal(b, m, deterministic)
}
func (m *CreateOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateOrderResponse.Merge(m, src)
}
func (m *CreateOrderResponse) XXX_Size() int {
	return xxx_messageInfo_CreateOrderResponse.Size(m)
}
func (m *CreateOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateOrderResponse proto.InternalMessageInfo

func (m *CreateOrderResponse) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

type SubmitOrderRequest struct {
	OrderId              int64    `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	FileHash             string   `protobuf:"bytes,2,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubmitOrderRequest) Reset()         { *m = SubmitOrderRequest{} }
func (m *SubmitOrderRequest) String() string { return proto.CompactTextString(m) }
func (*SubmitOrderRequest) ProtoMessage()    {}
func (*SubmitOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{4}
}

func (m *SubmitOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitOrderRequest.Unmarshal(m, b)
}
func (m *SubmitOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitOrderRequest.Marshal(b, m, deterministic)
}
func (m *SubmitOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitOrderRequest.Merge(m, src)
}
func (m *SubmitOrderRequest) XXX_Size() int {
	return xxx_messageInfo_SubmitOrderRequest.Size(m)
}
func (m *SubmitOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitOrderRequest proto.InternalMessageInfo

func (m *SubmitOrderRequest) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *SubmitOrderRequest) GetFileHash() string {
	if m != nil {
		return m.FileHash
	}
	return ""
}

type SubmitOrderResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubmitOrderResponse) Reset()         { *m = SubmitOrderResponse{} }
func (m *SubmitOrderResponse) String() string { return proto.CompactTextString(m) }
func (*SubmitOrderResponse) ProtoMessage()    {}
func (*SubmitOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{5}
}

func (m *SubmitOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitOrderResponse.Unmarshal(m, b)
}
func (m *SubmitOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitOrderResponse.Marshal(b, m, deterministic)
}
func (m *SubmitOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitOrderResponse.Merge(m, src)
}
func (m *SubmitOrderResponse) XXX_Size() int {
	return xxx_messageInfo_SubmitOrderResponse.Size(m)
}
func (m *SubmitOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitOrderResponse proto.InternalMessageInfo

type CloseOrderRequest struct {
	OrderId              int64    `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CloseOrderRequest) Reset()         { *m = CloseOrderRequest{} }
func (m *CloseOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CloseOrderRequest) ProtoMessage()    {}
func (*CloseOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{6}
}

func (m *CloseOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CloseOrderRequest.Unmarshal(m, b)
}
func (m *CloseOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CloseOrderRequest.Marshal(b, m, deterministic)
}
func (m *CloseOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloseOrderRequest.Merge(m, src)
}
func (m *CloseOrderRequest) XXX_Size() int {
	return xxx_messageInfo_CloseOrderRequest.Size(m)
}
func (m *CloseOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CloseOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CloseOrderRequest proto.InternalMessageInfo

func (m *CloseOrderRequest) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

type CloseOrderResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CloseOrderResponse) Reset()         { *m = CloseOrderResponse{} }
func (m *CloseOrderResponse) String() string { return proto.CompactTextString(m) }
func (*CloseOrderResponse) ProtoMessage()    {}
func (*CloseOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{7}
}

func (m *CloseOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CloseOrderResponse.Unmarshal(m, b)
}
func (m *CloseOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CloseOrderResponse.Marshal(b, m, deterministic)
}
func (m *CloseOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloseOrderResponse.Merge(m, src)
}
func (m *CloseOrderResponse) XXX_Size() int {
	return xxx_messageInfo_CloseOrderResponse.Size(m)
}
func (m *CloseOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CloseOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CloseOrderResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryBalanceRequest)(nil), "order.QueryBalanceRequest")
	proto.RegisterType((*QueryBalanceResponse)(nil), "order.QueryBalanceResponse")
	proto.RegisterType((*CreateOrderRequest)(nil), "order.CreateOrderRequest")
	proto.RegisterType((*CreateOrderResponse)(nil), "order.CreateOrderResponse")
	proto.RegisterType((*SubmitOrderRequest)(nil), "order.SubmitOrderRequest")
	proto.RegisterType((*SubmitOrderResponse)(nil), "order.SubmitOrderResponse")
	proto.RegisterType((*CloseOrderRequest)(nil), "order.CloseOrderRequest")
	proto.RegisterType((*CloseOrderResponse)(nil), "order.CloseOrderResponse")
}

func init() { proto.RegisterFile("order.proto", fileDescriptor_cd01338c35d87077) }

var fileDescriptor_cd01338c35d87077 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xed, 0x4a, 0x02, 0x41,
	0x14, 0x75, 0xb5, 0x0f, 0xf7, 0xea, 0x9f, 0xae, 0x06, 0xeb, 0x4a, 0x20, 0xf3, 0xcb, 0x5f, 0x26,
	0xf5, 0x06, 0x09, 0x91, 0x10, 0x45, 0xeb, 0x03, 0xc8, 0xe8, 0xde, 0x70, 0x40, 0x5d, 0x9b, 0x59,
	0x83, 0x7c, 0x81, 0xde, 0xa6, 0x67, 0x8c, 0x9d, 0x99, 0xdd, 0x66, 0x73, 0xa1, 0x7e, 0xde, 0x7b,
	0xe6, 0x9e, 0x73, 0x38, 0x87, 0x81, 0x56, 0x22, 0x63, 0x92, 0xa3, 0x9d, 0x4c, 0xd2, 0x04, 0x4f,
	0xf5, 0xc0, 0xae, 0xa1, 0xf3, 0xb2, 0x27, 0xf9, 0x71, 0xc7, 0xd7, 0x7c, 0xbb, 0xa4, 0x88, 0xde,
	0xf6, 0xa4, 0x52, 0x0c, 0xe0, 0x9c, 0xc7, 0xb1, 0x24, 0xa5, 0x02, 0x6f, 0xe0, 0x0d, 0xfd, 0x28,
	0x1f, 0xd9, 0x18, 0xba, 0xe5, 0x03, 0xb5, 0x4b, 0xb6, 0x8a, 0xb2, 0x8b, 0x85, 0x59, 0x05, 0x8d,
	0x81, 0x37, 0x6c, 0x44, 0xf9, 0xc8, 0x3e, 0x3d, 0xc0, 0x89, 0x24, 0x9e, 0xd2, 0x73, 0x26, 0xf9,
	0xa7, 0x04, 0x5e, 0x01, 0x48, 0xf3, 0x68, 0x2e, 0xe2, 0xa0, 0xae, 0x41, 0xdf, 0x6e, 0xa6, 0x31,
	0xf6, 0xc1, 0x7f, 0x15, 0x6b, 0x9a, 0x6f, 0xf9, 0xc6, 0x68, 0xf9, 0x51, 0x33, 0x5b, 0x3c, 0xf1,
	0x0d, 0x15, 0xa0, 0x12, 0x07, 0x0a, 0x4e, 0xb4, 0x11, 0x0d, 0xce, 0xc4, 0x81, 0xd8, 0x18, 0x3a,
	0x25, 0x23, 0xd6, 0x7a, 0x0f, 0x9a, 0x3a, 0x8c, 0x4c, 0xcd, 0x7a, 0xd7, 0xf3, 0x34, 0x66, 0x8f,
	0x80, 0xb3, 0xfd, 0x62, 0x23, 0xd2, 0x92, 0x75, 0xf7, 0xc0, 0x2b, 0x1d, 0x14, 0xfa, 0x2b, 0xae,
	0x56, 0xd6, 0xba, 0xd6, 0x7f, 0xe0, 0x6a, 0xc5, 0x2e, 0xa1, 0x53, 0x62, 0x33, 0xfa, 0x6c, 0x04,
	0x17, 0x93, 0x75, 0xa2, 0xe8, 0x9f, 0x1a, 0xac, 0x0b, 0xe8, 0xbe, 0x37, 0x2c, 0x37, 0x5f, 0x75,
	0x68, 0xeb, 0xcd, 0x8c, 0xe4, 0xbb, 0x58, 0x12, 0x4e, 0xa1, 0xed, 0x36, 0x85, 0xe1, 0xc8, 0xf4,
	0x5f, 0xd1, 0x77, 0xd8, 0xaf, 0xc4, 0xac, 0xbf, 0x1a, 0xde, 0x43, 0xcb, 0x09, 0x0e, 0x7b, 0xf6,
	0xf5, 0x71, 0xab, 0x61, 0x58, 0x05, 0xb9, 0x3c, 0x4e, 0x00, 0x05, 0xcf, 0x71, 0xc4, 0x05, 0x4f,
	0x55, 0x5e, 0x35, 0x9c, 0x00, 0xfc, 0x24, 0x80, 0x41, 0xae, 0xf9, 0x3b, 0xc4, 0xb0, 0x57, 0x81,
	0xe4, 0x24, 0x8b, 0x33, 0xfd, 0x11, 0x6e, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xd4, 0x02, 0x8b,
	0x56, 0x17, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrderServiceClient interface {
	QueryBalance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error)
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	SubmitOrder(ctx context.Context, in *SubmitOrderRequest, opts ...grpc.CallOption) (*SubmitOrderResponse, error)
	CloseOrder(ctx context.Context, in *CloseOrderRequest, opts ...grpc.CallOption) (*CloseOrderResponse, error)
}

type orderServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrderServiceClient(cc *grpc.ClientConn) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) QueryBalance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error) {
	out := new(QueryBalanceResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/QueryBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) SubmitOrder(ctx context.Context, in *SubmitOrderRequest, opts ...grpc.CallOption) (*SubmitOrderResponse, error) {
	out := new(SubmitOrderResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/SubmitOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CloseOrder(ctx context.Context, in *CloseOrderRequest, opts ...grpc.CallOption) (*CloseOrderResponse, error) {
	out := new(CloseOrderResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/CloseOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServiceServer is the server API for OrderService service.
type OrderServiceServer interface {
	QueryBalance(context.Context, *QueryBalanceRequest) (*QueryBalanceResponse, error)
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	SubmitOrder(context.Context, *SubmitOrderRequest) (*SubmitOrderResponse, error)
	CloseOrder(context.Context, *CloseOrderRequest) (*CloseOrderResponse, error)
}

func RegisterOrderServiceServer(s *grpc.Server, srv OrderServiceServer) {
	s.RegisterService(&_OrderService_serviceDesc, srv)
}

func _OrderService_QueryBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).QueryBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/QueryBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).QueryBalance(ctx, req.(*QueryBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_SubmitOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).SubmitOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/SubmitOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).SubmitOrder(ctx, req.(*SubmitOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CloseOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CloseOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/CloseOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CloseOrder(ctx, req.(*CloseOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrderService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "order.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryBalance",
			Handler:    _OrderService_QueryBalance_Handler,
		},
		{
			MethodName: "CreateOrder",
			Handler:    _OrderService_CreateOrder_Handler,
		},
		{
			MethodName: "SubmitOrder",
			Handler:    _OrderService_SubmitOrder_Handler,
		},
		{
			MethodName: "CloseOrder",
			Handler:    _OrderService_CloseOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
