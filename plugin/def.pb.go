// Code generated by protoc-gen-go. DO NOT EDIT.
// source: def.proto

package plugin

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type DefinedCmd int32

const (
	// 获取终端状态
	DefinedCmd_StatusPeer DefinedCmd = 0
	// 获取raft状态
	DefinedCmd_StatusRaft DefinedCmd = 1
	// leader切换
	DefinedCmd_LeaderSwift DefinedCmd = 2
)

var DefinedCmd_name = map[int32]string{
	0: "StatusPeer",
	1: "StatusRaft",
	2: "LeaderSwift",
}

var DefinedCmd_value = map[string]int32{
	"StatusPeer":  0,
	"StatusRaft":  1,
	"LeaderSwift": 2,
}

func (x DefinedCmd) String() string {
	return proto.EnumName(DefinedCmd_name, int32(x))
}

func (DefinedCmd) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{0}
}

// StatusPeer
type StatusPeerRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusPeerRequest) Reset()         { *m = StatusPeerRequest{} }
func (m *StatusPeerRequest) String() string { return proto.CompactTextString(m) }
func (*StatusPeerRequest) ProtoMessage()    {}
func (*StatusPeerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{0}
}

func (m *StatusPeerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusPeerRequest.Unmarshal(m, b)
}
func (m *StatusPeerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusPeerRequest.Marshal(b, m, deterministic)
}
func (m *StatusPeerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusPeerRequest.Merge(m, src)
}
func (m *StatusPeerRequest) XXX_Size() int {
	return xxx_messageInfo_StatusPeerRequest.Size(m)
}
func (m *StatusPeerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusPeerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusPeerRequest proto.InternalMessageInfo

type StatusPeerReply struct {
	Peers                []string `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusPeerReply) Reset()         { *m = StatusPeerReply{} }
func (m *StatusPeerReply) String() string { return proto.CompactTextString(m) }
func (*StatusPeerReply) ProtoMessage()    {}
func (*StatusPeerReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{1}
}

func (m *StatusPeerReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusPeerReply.Unmarshal(m, b)
}
func (m *StatusPeerReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusPeerReply.Marshal(b, m, deterministic)
}
func (m *StatusPeerReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusPeerReply.Merge(m, src)
}
func (m *StatusPeerReply) XXX_Size() int {
	return xxx_messageInfo_StatusPeerReply.Size(m)
}
func (m *StatusPeerReply) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusPeerReply.DiscardUnknown(m)
}

var xxx_messageInfo_StatusPeerReply proto.InternalMessageInfo

func (m *StatusPeerReply) GetPeers() []string {
	if m != nil {
		return m.Peers
	}
	return nil
}

// StatusRaft
type StatusRaftRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRaftRequest) Reset()         { *m = StatusRaftRequest{} }
func (m *StatusRaftRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRaftRequest) ProtoMessage()    {}
func (*StatusRaftRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{2}
}

func (m *StatusRaftRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRaftRequest.Unmarshal(m, b)
}
func (m *StatusRaftRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRaftRequest.Marshal(b, m, deterministic)
}
func (m *StatusRaftRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRaftRequest.Merge(m, src)
}
func (m *StatusRaftRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRaftRequest.Size(m)
}
func (m *StatusRaftRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRaftRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRaftRequest proto.InternalMessageInfo

type StatusRaftReply struct {
	FromRaftAddr         string   `protobuf:"bytes,1,opt,name=fromRaftAddr,proto3" json:"fromRaftAddr,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRaftReply) Reset()         { *m = StatusRaftReply{} }
func (m *StatusRaftReply) String() string { return proto.CompactTextString(m) }
func (*StatusRaftReply) ProtoMessage()    {}
func (*StatusRaftReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{3}
}

func (m *StatusRaftReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRaftReply.Unmarshal(m, b)
}
func (m *StatusRaftReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRaftReply.Marshal(b, m, deterministic)
}
func (m *StatusRaftReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRaftReply.Merge(m, src)
}
func (m *StatusRaftReply) XXX_Size() int {
	return xxx_messageInfo_StatusRaftReply.Size(m)
}
func (m *StatusRaftReply) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRaftReply.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRaftReply proto.InternalMessageInfo

func (m *StatusRaftReply) GetFromRaftAddr() string {
	if m != nil {
		return m.FromRaftAddr
	}
	return ""
}

func (m *StatusRaftReply) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

// LeaderSwift
type LeaderSwiftRequest struct {
	// raft.serverID
	SvrID                string   `protobuf:"bytes,1,opt,name=svrID,proto3" json:"svrID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeaderSwiftRequest) Reset()         { *m = LeaderSwiftRequest{} }
func (m *LeaderSwiftRequest) String() string { return proto.CompactTextString(m) }
func (*LeaderSwiftRequest) ProtoMessage()    {}
func (*LeaderSwiftRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{4}
}

func (m *LeaderSwiftRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaderSwiftRequest.Unmarshal(m, b)
}
func (m *LeaderSwiftRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaderSwiftRequest.Marshal(b, m, deterministic)
}
func (m *LeaderSwiftRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaderSwiftRequest.Merge(m, src)
}
func (m *LeaderSwiftRequest) XXX_Size() int {
	return xxx_messageInfo_LeaderSwiftRequest.Size(m)
}
func (m *LeaderSwiftRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaderSwiftRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LeaderSwiftRequest proto.InternalMessageInfo

func (m *LeaderSwiftRequest) GetSvrID() string {
	if m != nil {
		return m.SvrID
	}
	return ""
}

type LeaderSwiftReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeaderSwiftReply) Reset()         { *m = LeaderSwiftReply{} }
func (m *LeaderSwiftReply) String() string { return proto.CompactTextString(m) }
func (*LeaderSwiftReply) ProtoMessage()    {}
func (*LeaderSwiftReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{5}
}

func (m *LeaderSwiftReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaderSwiftReply.Unmarshal(m, b)
}
func (m *LeaderSwiftReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaderSwiftReply.Marshal(b, m, deterministic)
}
func (m *LeaderSwiftReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaderSwiftReply.Merge(m, src)
}
func (m *LeaderSwiftReply) XXX_Size() int {
	return xxx_messageInfo_LeaderSwiftReply.Size(m)
}
func (m *LeaderSwiftReply) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaderSwiftReply.DiscardUnknown(m)
}

var xxx_messageInfo_LeaderSwiftReply proto.InternalMessageInfo

type CallRequest struct {
	Cmd string `protobuf:"bytes,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Id  uint32 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// 数据存储区
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallRequest) Reset()         { *m = CallRequest{} }
func (m *CallRequest) String() string { return proto.CompactTextString(m) }
func (*CallRequest) ProtoMessage()    {}
func (*CallRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{6}
}

func (m *CallRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallRequest.Unmarshal(m, b)
}
func (m *CallRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallRequest.Marshal(b, m, deterministic)
}
func (m *CallRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallRequest.Merge(m, src)
}
func (m *CallRequest) XXX_Size() int {
	return xxx_messageInfo_CallRequest.Size(m)
}
func (m *CallRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CallRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CallRequest proto.InternalMessageInfo

func (m *CallRequest) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *CallRequest) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CallRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type CallReply struct {
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Desc string `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc,omitempty"`
	Cmd  string `protobuf:"bytes,3,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Id   uint32 `protobuf:"varint,4,opt,name=id,proto3" json:"id,omitempty"`
	// 数据存储区
	Data                 []byte   `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CallReply) Reset()         { *m = CallReply{} }
func (m *CallReply) String() string { return proto.CompactTextString(m) }
func (*CallReply) ProtoMessage()    {}
func (*CallReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_76fb0470a3b910d8, []int{7}
}

func (m *CallReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallReply.Unmarshal(m, b)
}
func (m *CallReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallReply.Marshal(b, m, deterministic)
}
func (m *CallReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallReply.Merge(m, src)
}
func (m *CallReply) XXX_Size() int {
	return xxx_messageInfo_CallReply.Size(m)
}
func (m *CallReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CallReply.DiscardUnknown(m)
}

var xxx_messageInfo_CallReply proto.InternalMessageInfo

func (m *CallReply) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *CallReply) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *CallReply) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *CallReply) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CallReply) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("plugin.DefinedCmd", DefinedCmd_name, DefinedCmd_value)
	proto.RegisterType((*StatusPeerRequest)(nil), "plugin.StatusPeerRequest")
	proto.RegisterType((*StatusPeerReply)(nil), "plugin.StatusPeerReply")
	proto.RegisterType((*StatusRaftRequest)(nil), "plugin.StatusRaftRequest")
	proto.RegisterType((*StatusRaftReply)(nil), "plugin.StatusRaftReply")
	proto.RegisterType((*LeaderSwiftRequest)(nil), "plugin.LeaderSwiftRequest")
	proto.RegisterType((*LeaderSwiftReply)(nil), "plugin.LeaderSwiftReply")
	proto.RegisterType((*CallRequest)(nil), "plugin.CallRequest")
	proto.RegisterType((*CallReply)(nil), "plugin.CallReply")
}

func init() { proto.RegisterFile("def.proto", fileDescriptor_76fb0470a3b910d8) }

var fileDescriptor_76fb0470a3b910d8 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x31, 0x4f, 0xfb, 0x30,
	0x10, 0xc5, 0x9b, 0xa4, 0xad, 0x94, 0x6b, 0xfb, 0x6f, 0xea, 0xfe, 0x85, 0x22, 0xa6, 0xca, 0x0b,
	0x55, 0x87, 0x0e, 0x20, 0xb1, 0x31, 0xa0, 0x74, 0x41, 0x02, 0x09, 0xb9, 0x9f, 0x20, 0xc4, 0x17,
	0x64, 0xc9, 0x6d, 0x8c, 0xed, 0x80, 0xfa, 0xed, 0x91, 0x9d, 0x84, 0xb4, 0xb0, 0xdd, 0x7b, 0xc9,
	0xef, 0xde, 0xf9, 0x6c, 0x88, 0x39, 0x96, 0x5b, 0xa5, 0x2b, 0x5b, 0x91, 0xb1, 0x92, 0xf5, 0xbb,
	0x38, 0xd2, 0x25, 0x2c, 0xf6, 0x36, 0xb7, 0xb5, 0x79, 0x45, 0xd4, 0x0c, 0x3f, 0x6a, 0x34, 0x96,
	0xde, 0xc0, 0xfc, 0xdc, 0x54, 0xf2, 0x44, 0xfe, 0xc3, 0x48, 0x21, 0x6a, 0x93, 0x06, 0xab, 0x68,
	0x1d, 0xb3, 0x46, 0xf4, 0x34, 0xcb, 0x4b, 0xdb, 0xd1, 0x2f, 0x1d, 0xdd, 0x98, 0x8e, 0xa6, 0x30,
	0x2d, 0x75, 0x75, 0x70, 0xc6, 0x23, 0xe7, 0x3a, 0x0d, 0x56, 0xc1, 0x3a, 0x66, 0x17, 0x1e, 0xb9,
	0x82, 0xb1, 0xf1, 0x58, 0x1a, 0xfa, 0xaf, 0xad, 0xa2, 0x1b, 0x20, 0xcf, 0x98, 0x73, 0xd4, 0xfb,
	0x2f, 0xf1, 0x13, 0xe2, 0xe6, 0x31, 0x9f, 0xfa, 0x69, 0xd7, 0xb6, 0x6a, 0x04, 0x25, 0x90, 0x5c,
	0xfc, 0xab, 0xe4, 0x89, 0x66, 0x30, 0xc9, 0x72, 0x29, 0x3b, 0x30, 0x81, 0xa8, 0x38, 0xf0, 0x16,
	0x73, 0x25, 0xf9, 0x07, 0xa1, 0xe0, 0x3e, 0x74, 0xc6, 0x42, 0xc1, 0x09, 0x81, 0x21, 0xcf, 0x6d,
	0x9e, 0x46, 0xab, 0x60, 0x3d, 0x65, 0xbe, 0xa6, 0x02, 0xe2, 0xa6, 0x89, 0x3b, 0x0d, 0x81, 0x61,
	0x51, 0x71, 0xf4, 0x3d, 0x66, 0xcc, 0xd7, 0x1e, 0x42, 0x53, 0xb4, 0xb3, 0xfb, 0xba, 0x8b, 0x8a,
	0x7e, 0x47, 0x0d, 0xff, 0x44, 0x8d, 0xfa, 0xa8, 0xcd, 0x03, 0xc0, 0x0e, 0x4b, 0x71, 0x44, 0x9e,
	0x79, 0x02, 0xfa, 0xab, 0x48, 0x06, 0xbd, 0x76, 0x7b, 0x4b, 0x02, 0x32, 0x87, 0xc9, 0xd9, 0x89,
	0x93, 0xf0, 0x76, 0x07, 0x90, 0xc9, 0xda, 0x58, 0xd4, 0x4c, 0x15, 0xe4, 0x1e, 0x60, 0x2f, 0x0e,
	0x4a, 0xa2, 0x9b, 0x9e, 0x2c, 0xb7, 0xcd, 0xad, 0x6f, 0xcf, 0x16, 0x72, 0xbd, 0xb8, 0x34, 0xdd,
	0xca, 0x06, 0x6f, 0x63, 0xff, 0x4a, 0xee, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0xd8, 0x3d, 0xab,
	0xa1, 0x32, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClusterRpcClient is the client API for ClusterRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterRpcClient interface {
	SimpleCall(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallReply, error)
}

type clusterRpcClient struct {
	cc *grpc.ClientConn
}

func NewClusterRpcClient(cc *grpc.ClientConn) ClusterRpcClient {
	return &clusterRpcClient{cc}
}

func (c *clusterRpcClient) SimpleCall(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallReply, error) {
	out := new(CallReply)
	err := c.cc.Invoke(ctx, "/plugin.ClusterRpc/SimpleCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterRpcServer is the server API for ClusterRpc service.
type ClusterRpcServer interface {
	SimpleCall(context.Context, *CallRequest) (*CallReply, error)
}

// UnimplementedClusterRpcServer can be embedded to have forward compatible implementations.
type UnimplementedClusterRpcServer struct {
}

func (*UnimplementedClusterRpcServer) SimpleCall(ctx context.Context, req *CallRequest) (*CallReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SimpleCall not implemented")
}

func RegisterClusterRpcServer(s *grpc.Server, srv ClusterRpcServer) {
	s.RegisterService(&_ClusterRpc_serviceDesc, srv)
}

func _ClusterRpc_SimpleCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterRpcServer).SimpleCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.ClusterRpc/SimpleCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterRpcServer).SimpleCall(ctx, req.(*CallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClusterRpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.ClusterRpc",
	HandlerType: (*ClusterRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SimpleCall",
			Handler:    _ClusterRpc_SimpleCall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "def.proto",
}
