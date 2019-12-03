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
	FileName             string   `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileSize             int64    `protobuf:"varint,3,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	FileHash             string   `protobuf:"bytes,4,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
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

func (m *SubmitOrderRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *SubmitOrderRequest) GetFileSize() int64 {
	if m != nil {
		return m.FileSize
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

type CancelOrderRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelOrderRequest) Reset()         { *m = CancelOrderRequest{} }
func (m *CancelOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CancelOrderRequest) ProtoMessage()    {}
func (*CancelOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{6}
}

func (m *CancelOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelOrderRequest.Unmarshal(m, b)
}
func (m *CancelOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelOrderRequest.Marshal(b, m, deterministic)
}
func (m *CancelOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelOrderRequest.Merge(m, src)
}
func (m *CancelOrderRequest) XXX_Size() int {
	return xxx_messageInfo_CancelOrderRequest.Size(m)
}
func (m *CancelOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CancelOrderRequest proto.InternalMessageInfo

type CancelOrderResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelOrderResponse) Reset()         { *m = CancelOrderResponse{} }
func (m *CancelOrderResponse) String() string { return proto.CompactTextString(m) }
func (*CancelOrderResponse) ProtoMessage()    {}
func (*CancelOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd01338c35d87077, []int{7}
}

func (m *CancelOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelOrderResponse.Unmarshal(m, b)
}
func (m *CancelOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelOrderResponse.Marshal(b, m, deterministic)
}
func (m *CancelOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelOrderResponse.Merge(m, src)
}
func (m *CancelOrderResponse) XXX_Size() int {
	return xxx_messageInfo_CancelOrderResponse.Size(m)
}
func (m *CancelOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CancelOrderResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryBalanceRequest)(nil), "order.QueryBalanceRequest")
	proto.RegisterType((*QueryBalanceResponse)(nil), "order.QueryBalanceResponse")
	proto.RegisterType((*CreateOrderRequest)(nil), "order.CreateOrderRequest")
	proto.RegisterType((*CreateOrderResponse)(nil), "order.CreateOrderResponse")
	proto.RegisterType((*SubmitOrderRequest)(nil), "order.SubmitOrderRequest")
	proto.RegisterType((*SubmitOrderResponse)(nil), "order.SubmitOrderResponse")
	proto.RegisterType((*CancelOrderRequest)(nil), "order.CancelOrderRequest")
	proto.RegisterType((*CancelOrderResponse)(nil), "order.CancelOrderResponse")
}

func init() { proto.RegisterFile("order.proto", fileDescriptor_cd01338c35d87077) }

var fileDescriptor_cd01338c35d87077 = []byte{
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xdb, 0x4e, 0x02, 0x31,
	0x10, 0x65, 0xc1, 0x0b, 0x0c, 0x3c, 0x15, 0x4c, 0x96, 0x25, 0x26, 0xa4, 0x4f, 0x3c, 0x21, 0xd1,
	0x3f, 0x90, 0xc4, 0xc8, 0x83, 0x1a, 0xe1, 0x03, 0x48, 0xa1, 0x63, 0x68, 0x02, 0x2c, 0xb6, 0x8b,
	0x89, 0x7c, 0x80, 0xff, 0xe3, 0x1f, 0x9a, 0x5e, 0xd8, 0x6d, 0xb3, 0x9b, 0xf8, 0x38, 0x73, 0xe6,
	0x4c, 0xcf, 0x9c, 0x53, 0x68, 0xa7, 0x92, 0xa3, 0x1c, 0x1f, 0x64, 0x9a, 0xa5, 0xe4, 0xd2, 0x14,
	0xf4, 0x0e, 0xba, 0xef, 0x47, 0x94, 0xdf, 0x8f, 0x6c, 0xcb, 0xf6, 0x6b, 0x9c, 0xe3, 0xe7, 0x11,
	0x55, 0x46, 0x62, 0xb8, 0x66, 0x9c, 0x4b, 0x54, 0x2a, 0x8e, 0x86, 0xd1, 0xa8, 0x35, 0x3f, 0x97,
	0x74, 0x02, 0xbd, 0x90, 0xa0, 0x0e, 0xe9, 0x5e, 0xa1, 0x66, 0xac, 0x6c, 0x2b, 0x6e, 0x0c, 0xa3,
	0x51, 0x63, 0x7e, 0x2e, 0xe9, 0x0b, 0x90, 0xa9, 0x44, 0x96, 0xe1, 0x9b, 0x7e, 0xf1, 0xdf, 0x17,
	0xc8, 0x2d, 0x80, 0xb4, 0x43, 0x4b, 0xc1, 0xe3, 0xba, 0x01, 0x5b, 0xae, 0x33, 0xe3, 0x74, 0x02,
	0xdd, 0x60, 0x9d, 0x7b, 0xbf, 0x0f, 0x4d, 0x73, 0x91, 0xe6, 0x38, 0x01, 0xa6, 0x9e, 0x71, 0xfa,
	0x13, 0x01, 0x59, 0x1c, 0x57, 0x3b, 0x91, 0x05, 0x0a, 0x7c, 0x46, 0x14, 0x30, 0xc8, 0x00, 0x5a,
	0x1f, 0x62, 0x8b, 0xcb, 0x3d, 0xdb, 0xa1, 0x53, 0xd0, 0xd4, 0x8d, 0x57, 0xb6, 0xc3, 0x1c, 0x54,
	0xe2, 0x74, 0xbe, 0xd5, 0x80, 0x0b, 0x71, 0x2a, 0xc0, 0x0d, 0x53, 0x9b, 0xf8, 0xa2, 0x60, 0x3e,
	0x33, 0xb5, 0xa1, 0x37, 0xd0, 0x0d, 0x74, 0x58, 0xe9, 0xb4, 0x07, 0x64, 0xaa, 0x9d, 0xda, 0xfa,
	0xf2, 0xf4, 0x70, 0xd0, 0xb5, 0xc3, 0xf7, 0xbf, 0x75, 0xe8, 0x98, 0xce, 0x02, 0xe5, 0x97, 0x58,
	0x23, 0x99, 0x41, 0xc7, 0x0f, 0x84, 0x24, 0x63, 0x1b, 0x73, 0x45, 0xac, 0xc9, 0xa0, 0x12, 0x73,
	0x32, 0x6a, 0xe4, 0x09, 0xda, 0x9e, 0xb5, 0xa4, 0xef, 0xa6, 0xcb, 0xe9, 0x25, 0x49, 0x15, 0xe4,
	0xef, 0xf1, 0xee, 0xcc, 0xf7, 0x94, 0x33, 0xc8, 0xf7, 0x54, 0xd9, 0x62, 0xf5, 0x14, 0x16, 0x14,
	0x7a, 0x4a, 0x66, 0x15, 0x7a, 0xca, 0x8e, 0xd1, 0xda, 0xea, 0xca, 0x7c, 0xf9, 0x87, 0xbf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x2e, 0xd4, 0xd2, 0x00, 0x01, 0x03, 0x00, 0x00,
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
	CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*CancelOrderResponse, error)
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

func (c *orderServiceClient) CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*CancelOrderResponse, error) {
	out := new(CancelOrderResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/CancelOrder", in, out, opts...)
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
	CancelOrder(context.Context, *CancelOrderRequest) (*CancelOrderResponse, error)
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

func _OrderService_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.OrderService/CancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CancelOrder(ctx, req.(*CancelOrderRequest))
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
			MethodName: "CancelOrder",
			Handler:    _OrderService_CancelOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}