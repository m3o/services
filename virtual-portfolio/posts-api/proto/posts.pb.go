// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/posts.proto

package posts

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

type Post struct {
	Uuid                    string     `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	FeedType                string     `protobuf:"bytes,2,opt,name=feed_type,json=feedType,proto3" json:"feed_type,omitempty"`
	FeedUuid                string     `protobuf:"bytes,3,opt,name=feed_uuid,json=feedUuid,proto3" json:"feed_uuid,omitempty"`
	Text                    string     `protobuf:"bytes,4,opt,name=text,proto3" json:"text,omitempty"`
	User                    *User      `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
	BullsCount              int32      `protobuf:"varint,6,opt,name=bulls_count,json=bullsCount,proto3" json:"bulls_count,omitempty"`
	BearsCount              int32      `protobuf:"varint,7,opt,name=bears_count,json=bearsCount,proto3" json:"bears_count,omitempty"`
	Comments                []*Comment `protobuf:"bytes,8,rep,name=comments,proto3" json:"comments,omitempty"`
	Opinion                 string     `protobuf:"bytes,9,opt,name=opinion,proto3" json:"opinion,omitempty"`
	EnhancedText            string     `protobuf:"bytes,10,opt,name=enhanced_text,json=enhancedText,proto3" json:"enhanced_text,omitempty"`
	Asset                   *Asset     `protobuf:"bytes,11,opt,name=asset,proto3" json:"asset,omitempty"`
	CreatedAt               string     `protobuf:"bytes,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	AttachmentPictureBase64 string     `protobuf:"bytes,13,opt,name=attachment_picture_base64,json=attachmentPictureBase64,proto3" json:"attachment_picture_base64,omitempty"`
	AttachmentPictureUrl    string     `protobuf:"bytes,14,opt,name=attachment_picture_url,json=attachmentPictureUrl,proto3" json:"attachment_picture_url,omitempty"`
	AttachmentLinkUrl       string     `protobuf:"bytes,15,opt,name=attachment_link_url,json=attachmentLinkUrl,proto3" json:"attachment_link_url,omitempty"`
	Title                   string     `protobuf:"bytes,16,opt,name=title,proto3" json:"title,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}   `json:"-"`
	XXX_unrecognized        []byte     `json:"-"`
	XXX_sizecache           int32      `json:"-"`
}

func (m *Post) Reset()         { *m = Post{} }
func (m *Post) String() string { return proto.CompactTextString(m) }
func (*Post) ProtoMessage()    {}
func (*Post) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{0}
}

func (m *Post) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Post.Unmarshal(m, b)
}
func (m *Post) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Post.Marshal(b, m, deterministic)
}
func (m *Post) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Post.Merge(m, src)
}
func (m *Post) XXX_Size() int {
	return xxx_messageInfo_Post.Size(m)
}
func (m *Post) XXX_DiscardUnknown() {
	xxx_messageInfo_Post.DiscardUnknown(m)
}

var xxx_messageInfo_Post proto.InternalMessageInfo

func (m *Post) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Post) GetFeedType() string {
	if m != nil {
		return m.FeedType
	}
	return ""
}

func (m *Post) GetFeedUuid() string {
	if m != nil {
		return m.FeedUuid
	}
	return ""
}

func (m *Post) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Post) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Post) GetBullsCount() int32 {
	if m != nil {
		return m.BullsCount
	}
	return 0
}

func (m *Post) GetBearsCount() int32 {
	if m != nil {
		return m.BearsCount
	}
	return 0
}

func (m *Post) GetComments() []*Comment {
	if m != nil {
		return m.Comments
	}
	return nil
}

func (m *Post) GetOpinion() string {
	if m != nil {
		return m.Opinion
	}
	return ""
}

func (m *Post) GetEnhancedText() string {
	if m != nil {
		return m.EnhancedText
	}
	return ""
}

func (m *Post) GetAsset() *Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *Post) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

func (m *Post) GetAttachmentPictureBase64() string {
	if m != nil {
		return m.AttachmentPictureBase64
	}
	return ""
}

func (m *Post) GetAttachmentPictureUrl() string {
	if m != nil {
		return m.AttachmentPictureUrl
	}
	return ""
}

func (m *Post) GetAttachmentLinkUrl() string {
	if m != nil {
		return m.AttachmentLinkUrl
	}
	return ""
}

func (m *Post) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

type Comment struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	User                 *User    `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Post                 *Post    `protobuf:"bytes,3,opt,name=post,proto3" json:"post,omitempty"`
	Text                 string   `protobuf:"bytes,4,opt,name=text,proto3" json:"text,omitempty"`
	EnhancedText         string   `protobuf:"bytes,5,opt,name=enhanced_text,json=enhancedText,proto3" json:"enhanced_text,omitempty"`
	BullsCount           int32    `protobuf:"varint,6,opt,name=bulls_count,json=bullsCount,proto3" json:"bulls_count,omitempty"`
	BearsCount           int32    `protobuf:"varint,7,opt,name=bears_count,json=bearsCount,proto3" json:"bears_count,omitempty"`
	Opinion              string   `protobuf:"bytes,8,opt,name=opinion,proto3" json:"opinion,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Comment) Reset()         { *m = Comment{} }
func (m *Comment) String() string { return proto.CompactTextString(m) }
func (*Comment) ProtoMessage()    {}
func (*Comment) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{1}
}

func (m *Comment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Comment.Unmarshal(m, b)
}
func (m *Comment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Comment.Marshal(b, m, deterministic)
}
func (m *Comment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Comment.Merge(m, src)
}
func (m *Comment) XXX_Size() int {
	return xxx_messageInfo_Comment.Size(m)
}
func (m *Comment) XXX_DiscardUnknown() {
	xxx_messageInfo_Comment.DiscardUnknown(m)
}

var xxx_messageInfo_Comment proto.InternalMessageInfo

func (m *Comment) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Comment) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Comment) GetPost() *Post {
	if m != nil {
		return m.Post
	}
	return nil
}

func (m *Comment) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Comment) GetEnhancedText() string {
	if m != nil {
		return m.EnhancedText
	}
	return ""
}

func (m *Comment) GetBullsCount() int32 {
	if m != nil {
		return m.BullsCount
	}
	return 0
}

func (m *Comment) GetBearsCount() int32 {
	if m != nil {
		return m.BearsCount
	}
	return 0
}

func (m *Comment) GetOpinion() string {
	if m != nil {
		return m.Opinion
	}
	return ""
}

type User struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	FirstName            string   `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string   `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	ProfilePictureUrl    string   `protobuf:"bytes,4,opt,name=profile_picture_url,json=profilePictureUrl,proto3" json:"profile_picture_url,omitempty"`
	Username             string   `protobuf:"bytes,5,opt,name=username,proto3" json:"username,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{2}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *User) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *User) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *User) GetProfilePictureUrl() string {
	if m != nil {
		return m.ProfilePictureUrl
	}
	return ""
}

func (m *User) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

type Error struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{3}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Response struct {
	Error                *Error   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Post                 *Post    `protobuf:"bytes,2,opt,name=post,proto3" json:"post,omitempty"`
	Comment              *Comment `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{4}
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

func (m *Response) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *Response) GetPost() *Post {
	if m != nil {
		return m.Post
	}
	return nil
}

func (m *Response) GetComment() *Comment {
	if m != nil {
		return m.Comment
	}
	return nil
}

type Asset struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	ProfilePictureUrl    string   `protobuf:"bytes,4,opt,name=profile_picture_url,json=profilePictureUrl,proto3" json:"profile_picture_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Asset) Reset()         { *m = Asset{} }
func (m *Asset) String() string { return proto.CompactTextString(m) }
func (*Asset) ProtoMessage()    {}
func (*Asset) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{5}
}

func (m *Asset) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Asset.Unmarshal(m, b)
}
func (m *Asset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Asset.Marshal(b, m, deterministic)
}
func (m *Asset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Asset.Merge(m, src)
}
func (m *Asset) XXX_Size() int {
	return xxx_messageInfo_Asset.Size(m)
}
func (m *Asset) XXX_DiscardUnknown() {
	xxx_messageInfo_Asset.DiscardUnknown(m)
}

var xxx_messageInfo_Asset proto.InternalMessageInfo

func (m *Asset) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Asset) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Asset) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Asset) GetProfilePictureUrl() string {
	if m != nil {
		return m.ProfilePictureUrl
	}
	return ""
}

func init() {
	proto.RegisterType((*Post)(nil), "Post")
	proto.RegisterType((*Comment)(nil), "Comment")
	proto.RegisterType((*User)(nil), "User")
	proto.RegisterType((*Error)(nil), "Error")
	proto.RegisterType((*Response)(nil), "Response")
	proto.RegisterType((*Asset)(nil), "Asset")
}

func init() { proto.RegisterFile("proto/posts.proto", fileDescriptor_e93dc7d934d9dc10) }

var fileDescriptor_e93dc7d934d9dc10 = []byte{
	// 640 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0xfe, 0x25, 0xb5, 0x13, 0x67, 0xd2, 0xfe, 0x20, 0xdb, 0x0a, 0x36, 0x25, 0x15, 0x91, 0x41,
	0x28, 0xe2, 0x60, 0xa4, 0x52, 0x38, 0x70, 0x2b, 0x05, 0x71, 0x41, 0x50, 0xb9, 0xcd, 0x39, 0x72,
	0x9d, 0x69, 0x6b, 0xd5, 0xb1, 0x2d, 0xef, 0x5a, 0xa2, 0xaf, 0xc2, 0x05, 0x89, 0xa7, 0xe3, 0x31,
	0xd0, 0xcc, 0x7a, 0x93, 0x54, 0x71, 0x85, 0x10, 0xb7, 0xd9, 0xef, 0x9b, 0x6f, 0x3c, 0x7f, 0x0d,
	0x83, 0xa2, 0xcc, 0x75, 0xfe, 0xaa, 0xc8, 0x95, 0x56, 0x01, 0xdb, 0xfe, 0x0f, 0x07, 0x9c, 0xd3,
	0x5c, 0x69, 0x21, 0xc0, 0xa9, 0xaa, 0x64, 0x2e, 0x5b, 0xe3, 0xd6, 0xa4, 0x17, 0xb2, 0x2d, 0x9e,
	0x40, 0xef, 0x12, 0x71, 0x3e, 0xd3, 0xb7, 0x05, 0xca, 0x36, 0x13, 0x1e, 0x01, 0xe7, 0xb7, 0x05,
	0x2e, 0x49, 0x56, 0x6d, 0xad, 0xc8, 0x29, 0x29, 0x05, 0x38, 0x1a, 0xbf, 0x69, 0xe9, 0x98, 0x68,
	0x64, 0x8b, 0x21, 0x38, 0x95, 0xc2, 0x52, 0xba, 0xe3, 0xd6, 0xa4, 0x7f, 0xe8, 0x06, 0x53, 0x85,
	0x65, 0xc8, 0x90, 0x78, 0x0a, 0xfd, 0x8b, 0x2a, 0x4d, 0xd5, 0x2c, 0xce, 0xab, 0x4c, 0xcb, 0xce,
	0xb8, 0x35, 0x71, 0x43, 0x60, 0xe8, 0x84, 0x10, 0x76, 0xc0, 0xa8, 0xb4, 0x0e, 0xdd, 0xda, 0x81,
	0x20, 0xe3, 0xf0, 0x1c, 0xbc, 0x38, 0x5f, 0x2c, 0x30, 0xd3, 0x4a, 0x7a, 0xe3, 0xad, 0x49, 0xff,
	0xd0, 0x0b, 0x4e, 0x0c, 0x10, 0x2e, 0x19, 0x21, 0xa1, 0x9b, 0x17, 0x49, 0x96, 0xe4, 0x99, 0xec,
	0x71, 0x66, 0xf6, 0x29, 0x9e, 0xc1, 0x0e, 0x66, 0xd7, 0x51, 0x16, 0x53, 0xb9, 0x94, 0x39, 0x30,
	0xbf, 0x6d, 0xc1, 0x73, 0xaa, 0x60, 0x04, 0x6e, 0xa4, 0x14, 0x6a, 0xd9, 0xe7, 0x12, 0x3a, 0xc1,
	0x31, 0xbd, 0x42, 0x03, 0x8a, 0x03, 0x80, 0xb8, 0xc4, 0x48, 0xe3, 0x7c, 0x16, 0x69, 0xb9, 0xcd,
	0xfa, 0x5e, 0x8d, 0x1c, 0x6b, 0xf1, 0x0e, 0x86, 0x91, 0xd6, 0x51, 0x7c, 0x4d, 0xa9, 0xcc, 0x8a,
	0x24, 0xd6, 0x55, 0x89, 0xb3, 0x8b, 0x48, 0xe1, 0xdb, 0x23, 0xb9, 0xc3, 0xde, 0x8f, 0x57, 0x0e,
	0xa7, 0x86, 0x7f, 0xcf, 0xb4, 0x38, 0x82, 0x47, 0x0d, 0xda, 0xaa, 0x4c, 0xe5, 0xff, 0x2c, 0xdc,
	0xdb, 0x10, 0x4e, 0xcb, 0x54, 0x04, 0xb0, 0xbb, 0xa6, 0x4a, 0x93, 0xec, 0x86, 0x25, 0x0f, 0x58,
	0x32, 0x58, 0x51, 0x9f, 0x93, 0xec, 0x86, 0xfc, 0xf7, 0xc0, 0xd5, 0x89, 0x4e, 0x51, 0x3e, 0x64,
	0x0f, 0xf3, 0xf0, 0x7f, 0xb5, 0xa0, 0x5b, 0x77, 0xb2, 0x71, 0x49, 0xec, 0x58, 0xdb, 0x9b, 0x63,
	0x1d, 0x82, 0x43, 0xbb, 0xc6, 0xdb, 0x41, 0x14, 0x2d, 0x5a, 0xc8, 0x50, 0xe3, 0x82, 0x6c, 0xcc,
	0xc0, 0x6d, 0x98, 0xc1, 0xbf, 0xaf, 0xca, 0xda, 0x12, 0x78, 0x77, 0x96, 0xc0, 0xff, 0xd9, 0x02,
	0x87, 0xd2, 0x6f, 0xac, 0xf3, 0x00, 0xe0, 0x32, 0x29, 0x95, 0x9e, 0x65, 0xd1, 0xc2, 0x5e, 0x43,
	0x8f, 0x91, 0x2f, 0xd1, 0x82, 0xcf, 0x21, 0x8d, 0x2c, 0x5b, 0x9f, 0x03, 0x01, 0x4c, 0x06, 0xb0,
	0x5b, 0x94, 0xf9, 0x65, 0x92, 0xe2, 0x9d, 0xe1, 0x99, 0xe2, 0x07, 0x35, 0xb5, 0x36, 0xb9, 0x7d,
	0xf0, 0xa8, 0x81, 0x1c, 0xcb, 0x34, 0x61, 0xf9, 0xf6, 0xdf, 0x80, 0xfb, 0xb1, 0x2c, 0x73, 0x4e,
	0x32, 0xce, 0xe7, 0xc8, 0x49, 0xba, 0x21, 0xdb, 0x54, 0xdb, 0x02, 0x95, 0x8a, 0xae, 0x6c, 0x86,
	0xf6, 0xe9, 0x5f, 0x81, 0x17, 0xa2, 0x2a, 0xf2, 0x4c, 0x21, 0xed, 0x31, 0x52, 0x08, 0x96, 0xd2,
	0x1e, 0x73, 0xc0, 0xd0, 0x80, 0xcb, 0xa9, 0xb5, 0x37, 0xa7, 0xe6, 0x43, 0xb7, 0xbe, 0xa5, 0x7a,
	0xa6, 0xab, 0x23, 0xb3, 0x84, 0xaf, 0xc0, 0xe5, 0xb3, 0xe0, 0x11, 0xd3, 0x8f, 0xa3, 0x6e, 0x22,
	0xd9, 0xcb, 0xc6, 0xb6, 0xd7, 0x1a, 0x2b, 0xc0, 0x59, 0x6b, 0x1a, 0xdb, 0x7f, 0xdb, 0xb0, 0xc3,
	0xef, 0x6d, 0x70, 0x29, 0x4f, 0x25, 0x46, 0xd0, 0x39, 0xe1, 0x9b, 0x13, 0x26, 0xf3, 0xfd, 0x5e,
	0x60, 0xeb, 0xf6, 0xff, 0x13, 0x43, 0xd8, 0xfa, 0x84, 0xba, 0x91, 0x1a, 0x41, 0x67, 0x5a, 0xcc,
	0xef, 0x13, 0x8e, 0xa0, 0xf3, 0x01, 0x53, 0xbc, 0x87, 0x7d, 0x01, 0x3b, 0xe6, 0xa3, 0xf6, 0x50,
	0x96, 0x7d, 0xd9, 0xf0, 0x33, 0x51, 0xfe, 0xe0, 0x37, 0x06, 0x38, 0x43, 0xfd, 0xb5, 0xfe, 0x37,
	0x35, 0x7d, 0xf1, 0x25, 0x0c, 0xce, 0x50, 0xd7, 0x62, 0xeb, 0xd8, 0x1c, 0xed, 0xa2, 0xc3, 0xbf,
	0xfa, 0xd7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xdb, 0x65, 0xa3, 0xe0, 0xff, 0x05, 0x00, 0x00,
}
