// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/importer/importer.proto

package go_micro_api_apps_importer

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type ImportRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportRequest) Reset()         { *m = ImportRequest{} }
func (m *ImportRequest) String() string { return proto.CompactTextString(m) }
func (*ImportRequest) ProtoMessage()    {}
func (*ImportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbec88664242c3b9, []int{0}
}

func (m *ImportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportRequest.Unmarshal(m, b)
}
func (m *ImportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportRequest.Marshal(b, m, deterministic)
}
func (m *ImportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportRequest.Merge(m, src)
}
func (m *ImportRequest) XXX_Size() int {
	return xxx_messageInfo_ImportRequest.Size(m)
}
func (m *ImportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ImportRequest proto.InternalMessageInfo

type ImportResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportResponse) Reset()         { *m = ImportResponse{} }
func (m *ImportResponse) String() string { return proto.CompactTextString(m) }
func (*ImportResponse) ProtoMessage()    {}
func (*ImportResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbec88664242c3b9, []int{1}
}

func (m *ImportResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportResponse.Unmarshal(m, b)
}
func (m *ImportResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportResponse.Marshal(b, m, deterministic)
}
func (m *ImportResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportResponse.Merge(m, src)
}
func (m *ImportResponse) XXX_Size() int {
	return xxx_messageInfo_ImportResponse.Size(m)
}
func (m *ImportResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ImportResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ImportRequest)(nil), "go.micro.api.apps.importer.ImportRequest")
	proto.RegisterType((*ImportResponse)(nil), "go.micro.api.apps.importer.ImportResponse")
}

func init() { proto.RegisterFile("proto/importer/importer.proto", fileDescriptor_cbec88664242c3b9) }

var fileDescriptor_cbec88664242c3b9 = []byte{
	// 125 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0xcc, 0x2d, 0xc8, 0x2f, 0x2a, 0x49, 0x2d, 0x82, 0x33, 0xf4, 0xc0, 0xe2, 0x42,
	0x52, 0xe9, 0xf9, 0x7a, 0xb9, 0x99, 0xc9, 0x45, 0xf9, 0x7a, 0x89, 0x05, 0x99, 0x7a, 0x89, 0x05,
	0x05, 0xc5, 0x7a, 0x30, 0x15, 0x4a, 0xfc, 0x5c, 0xbc, 0x9e, 0x60, 0x76, 0x50, 0x6a, 0x61, 0x69,
	0x6a, 0x71, 0x89, 0x92, 0x00, 0x17, 0x1f, 0x4c, 0xa0, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0xd5, 0x28,
	0x9b, 0x8b, 0xc3, 0x13, 0xaa, 0x5c, 0x28, 0x9e, 0x8b, 0x0d, 0xc2, 0x16, 0xd2, 0xd4, 0xc3, 0x6d,
	0xaa, 0x1e, 0x8a, 0x91, 0x52, 0x5a, 0xc4, 0x28, 0x85, 0x58, 0x96, 0xc4, 0x06, 0x76, 0xb2, 0x31,
	0x20, 0x00, 0x00, 0xff, 0xff, 0xbd, 0x4f, 0x20, 0x59, 0xd3, 0x00, 0x00, 0x00,
}
