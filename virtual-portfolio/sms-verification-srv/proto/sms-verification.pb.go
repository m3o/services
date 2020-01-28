// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/sms-verification.proto

package sms_verification

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

type Verification struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	PhoneNumber          string   `protobuf:"bytes,2,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Code                 string   `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	Verified             bool     `protobuf:"varint,4,opt,name=verified,proto3" json:"verified,omitempty"`
	Expired              bool     `protobuf:"varint,5,opt,name=expired,proto3" json:"expired,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Verification) Reset()         { *m = Verification{} }
func (m *Verification) String() string { return proto.CompactTextString(m) }
func (*Verification) ProtoMessage()    {}
func (*Verification) Descriptor() ([]byte, []int) {
	return fileDescriptor_d05c580cca936553, []int{0}
}

func (m *Verification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Verification.Unmarshal(m, b)
}
func (m *Verification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Verification.Marshal(b, m, deterministic)
}
func (m *Verification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Verification.Merge(m, src)
}
func (m *Verification) XXX_Size() int {
	return xxx_messageInfo_Verification.Size(m)
}
func (m *Verification) XXX_DiscardUnknown() {
	xxx_messageInfo_Verification.DiscardUnknown(m)
}

var xxx_messageInfo_Verification proto.InternalMessageInfo

func (m *Verification) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Verification) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *Verification) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Verification) GetVerified() bool {
	if m != nil {
		return m.Verified
	}
	return false
}

func (m *Verification) GetExpired() bool {
	if m != nil {
		return m.Expired
	}
	return false
}

func init() {
	proto.RegisterType((*Verification)(nil), "Verification")
}

func init() { proto.RegisterFile("proto/sms-verification.proto", fileDescriptor_d05c580cca936553) }

var fileDescriptor_d05c580cca936553 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x2f, 0xce, 0x2d, 0xd6, 0x2d, 0x4b, 0x2d, 0xca, 0x4c, 0xcb, 0x4c, 0x4e, 0x2c, 0xc9,
	0xcc, 0xcf, 0xd3, 0x03, 0x0b, 0x2b, 0xf5, 0x33, 0x72, 0xf1, 0x84, 0x21, 0x09, 0x0b, 0x09, 0x71,
	0xb1, 0x94, 0x96, 0x66, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42, 0x8a,
	0x5c, 0x3c, 0x05, 0x19, 0xf9, 0x79, 0xa9, 0xf1, 0x79, 0xa5, 0xb9, 0x49, 0xa9, 0x45, 0x12, 0x4c,
	0x60, 0x39, 0x6e, 0xb0, 0x98, 0x1f, 0x58, 0x08, 0xa4, 0x2d, 0x39, 0x3f, 0x25, 0x55, 0x82, 0x19,
	0xa2, 0x0d, 0xc4, 0x16, 0x92, 0xe2, 0xe2, 0x80, 0xd8, 0x98, 0x9a, 0x22, 0xc1, 0xa2, 0xc0, 0xa8,
	0xc1, 0x11, 0x04, 0xe7, 0x0b, 0x49, 0x70, 0xb1, 0xa7, 0x56, 0x14, 0x64, 0x16, 0xa5, 0xa6, 0x48,
	0xb0, 0x82, 0xa5, 0x60, 0x5c, 0xa3, 0x5e, 0x46, 0x2e, 0xfe, 0x60, 0xdf, 0x60, 0x14, 0x47, 0x69,
	0x72, 0xb1, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0xf1, 0xea, 0x21, 0xcb, 0x48, 0xa1,
	0x72, 0x95, 0x18, 0x84, 0x34, 0xb8, 0xd8, 0xc0, 0x22, 0x95, 0x04, 0x55, 0xaa, 0x72, 0x31, 0xbb,
	0xa7, 0x12, 0x34, 0x30, 0x89, 0x0d, 0x1c, 0x50, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x85,
	0xde, 0x67, 0x77, 0x48, 0x01, 0x00, 0x00,
}
