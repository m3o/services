// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/notifications.proto

package notifications

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

type PushToken struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushToken) Reset()         { *m = PushToken{} }
func (m *PushToken) String() string { return proto.CompactTextString(m) }
func (*PushToken) ProtoMessage()    {}
func (*PushToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_c97242c58cf45775, []int{0}
}

func (m *PushToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushToken.Unmarshal(m, b)
}
func (m *PushToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushToken.Marshal(b, m, deterministic)
}
func (m *PushToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushToken.Merge(m, src)
}
func (m *PushToken) XXX_Size() int {
	return xxx_messageInfo_PushToken.Size(m)
}
func (m *PushToken) XXX_DiscardUnknown() {
	xxx_messageInfo_PushToken.DiscardUnknown(m)
}

var xxx_messageInfo_PushToken proto.InternalMessageInfo

func (m *PushToken) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type Query struct {
	Page                 int64    `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit                int64    `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_c97242c58cf45775, []int{1}
}

func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (m *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(m, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetPage() int64 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *Query) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type Response struct {
	Notifications        []*Notification `protobuf:"bytes,1,rep,name=notifications,proto3" json:"notifications,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_c97242c58cf45775, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetNotifications() []*Notification {
	if m != nil {
		return m.Notifications
	}
	return nil
}

type Notification struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	CreatedAt            string   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Seen                 bool     `protobuf:"varint,3,opt,name=seen,proto3" json:"seen,omitempty"`
	Title                string   `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Description          string   `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	ResourceType         string   `protobuf:"bytes,6,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	ResourceUuid         string   `protobuf:"bytes,7,opt,name=resource_uuid,json=resourceUuid,proto3" json:"resource_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Notification) Reset()         { *m = Notification{} }
func (m *Notification) String() string { return proto.CompactTextString(m) }
func (*Notification) ProtoMessage()    {}
func (*Notification) Descriptor() ([]byte, []int) {
	return fileDescriptor_c97242c58cf45775, []int{3}
}

func (m *Notification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Notification.Unmarshal(m, b)
}
func (m *Notification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Notification.Marshal(b, m, deterministic)
}
func (m *Notification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Notification.Merge(m, src)
}
func (m *Notification) XXX_Size() int {
	return xxx_messageInfo_Notification.Size(m)
}
func (m *Notification) XXX_DiscardUnknown() {
	xxx_messageInfo_Notification.DiscardUnknown(m)
}

var xxx_messageInfo_Notification proto.InternalMessageInfo

func (m *Notification) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Notification) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

func (m *Notification) GetSeen() bool {
	if m != nil {
		return m.Seen
	}
	return false
}

func (m *Notification) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Notification) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Notification) GetResourceType() string {
	if m != nil {
		return m.ResourceType
	}
	return ""
}

func (m *Notification) GetResourceUuid() string {
	if m != nil {
		return m.ResourceUuid
	}
	return ""
}

func init() {
	proto.RegisterType((*PushToken)(nil), "PushToken")
	proto.RegisterType((*Query)(nil), "Query")
	proto.RegisterType((*Response)(nil), "Response")
	proto.RegisterType((*Notification)(nil), "Notification")
}

func init() { proto.RegisterFile("proto/notifications.proto", fileDescriptor_c97242c58cf45775) }

var fileDescriptor_c97242c58cf45775 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0x4d, 0x4f, 0xc2, 0x40,
	0x10, 0xb5, 0x16, 0x90, 0x0e, 0xf4, 0xe0, 0xc6, 0xc3, 0x8a, 0x31, 0xa9, 0xf5, 0xc2, 0xc1, 0xd4,
	0x08, 0x3f, 0xc0, 0x78, 0xf2, 0x66, 0x74, 0xc5, 0x33, 0xa9, 0xed, 0x88, 0x1b, 0xb1, 0xbb, 0xd9,
	0x0f, 0x13, 0x7e, 0xa9, 0x7f, 0xc7, 0xec, 0x00, 0xda, 0x26, 0xde, 0xde, 0xbc, 0xf7, 0x26, 0xf3,
	0x66, 0x06, 0x4e, 0xb5, 0x51, 0x4e, 0x5d, 0x37, 0xca, 0xc9, 0x37, 0x59, 0x95, 0x4e, 0xaa, 0xc6,
	0x16, 0xc4, 0xe5, 0x17, 0x90, 0x3c, 0x7a, 0xfb, 0xbe, 0x50, 0x1f, 0xd8, 0xb0, 0x13, 0xe8, 0xbb,
	0x00, 0x78, 0x94, 0x45, 0xd3, 0x44, 0x6c, 0x8b, 0xfc, 0x06, 0xfa, 0x4f, 0x1e, 0xcd, 0x86, 0x31,
	0xe8, 0xe9, 0x72, 0x85, 0xa4, 0xc6, 0x82, 0x70, 0x68, 0x59, 0xcb, 0x4f, 0xe9, 0xf8, 0x21, 0x91,
	0xdb, 0x22, 0xbf, 0x85, 0xa1, 0x40, 0xab, 0x55, 0x63, 0x91, 0xcd, 0x21, 0xed, 0x0c, 0xe6, 0x51,
	0x16, 0x4f, 0x47, 0xb3, 0xb4, 0x78, 0x68, 0xb1, 0xa2, 0xeb, 0xc9, 0xbf, 0x23, 0x18, 0xb7, 0xf5,
	0x30, 0xdb, 0x7b, 0x59, 0xef, 0x92, 0x11, 0x66, 0xe7, 0x00, 0x95, 0xc1, 0xd2, 0x61, 0xbd, 0x2c,
	0xb7, 0x01, 0x12, 0x91, 0xec, 0x98, 0x3b, 0x17, 0x5a, 0x2c, 0x62, 0xc3, 0xe3, 0x2c, 0x9a, 0x0e,
	0x05, 0x61, 0xda, 0x50, 0xba, 0x35, 0xf2, 0xde, 0x6e, 0xc3, 0x50, 0xb0, 0x0c, 0x46, 0x35, 0xda,
	0xca, 0x48, 0x1d, 0x66, 0xf1, 0x3e, 0x69, 0x6d, 0x8a, 0x5d, 0x42, 0x6a, 0xd0, 0x2a, 0x6f, 0x2a,
	0x5c, 0xba, 0x8d, 0x46, 0x3e, 0x20, 0xcf, 0x78, 0x4f, 0x2e, 0x36, 0x1a, 0x3b, 0x26, 0x0a, 0x7b,
	0xd4, 0x35, 0xbd, 0x78, 0x59, 0xcf, 0xbe, 0x20, 0x6d, 0x2f, 0x66, 0xd9, 0x04, 0xe2, 0x7b, 0x74,
	0x6c, 0x50, 0xd0, 0x91, 0x27, 0x49, 0xb1, 0xbf, 0x5c, 0x7e, 0xc0, 0xce, 0xa0, 0xf7, 0x1c, 0x62,
	0xff, 0x2b, 0x5e, 0xc1, 0xb1, 0xc0, 0x95, 0xb4, 0x0e, 0xcd, 0xdf, 0x0b, 0xa1, 0xf8, 0xc5, 0x1d,
	0xf7, 0xeb, 0x80, 0xfe, 0x3d, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x63, 0xa5, 0x24, 0x99, 0x0c,
	0x02, 0x00, 0x00,
}
