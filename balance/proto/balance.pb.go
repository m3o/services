// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.5
// source: proto/balance.proto

package balance

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

type IncrementRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId     string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Delta          int64  `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"`
	IdempotencyKey string `protobuf:"bytes,3,opt,name=idempotency_key,json=idempotencyKey,proto3" json:"idempotency_key,omitempty"`
	Visible        bool   `protobuf:"varint,4,opt,name=visible,proto3" json:"visible,omitempty"`
	Reference      string `protobuf:"bytes,5,opt,name=reference,proto3" json:"reference,omitempty"`
}

func (x *IncrementRequest) Reset() {
	*x = IncrementRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncrementRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncrementRequest) ProtoMessage() {}

func (x *IncrementRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncrementRequest.ProtoReflect.Descriptor instead.
func (*IncrementRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{0}
}

func (x *IncrementRequest) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

func (x *IncrementRequest) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *IncrementRequest) GetIdempotencyKey() string {
	if x != nil {
		return x.IdempotencyKey
	}
	return ""
}

func (x *IncrementRequest) GetVisible() bool {
	if x != nil {
		return x.Visible
	}
	return false
}

func (x *IncrementRequest) GetReference() string {
	if x != nil {
		return x.Reference
	}
	return ""
}

type IncrementResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewBalance int64 `protobuf:"varint,1,opt,name=new_balance,json=newBalance,proto3" json:"new_balance,omitempty"`
}

func (x *IncrementResponse) Reset() {
	*x = IncrementResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncrementResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncrementResponse) ProtoMessage() {}

func (x *IncrementResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncrementResponse.ProtoReflect.Descriptor instead.
func (*IncrementResponse) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{1}
}

func (x *IncrementResponse) GetNewBalance() int64 {
	if x != nil {
		return x.NewBalance
	}
	return 0
}

type DecrementRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId     string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Delta          int64  `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"`
	IdempotencyKey string `protobuf:"bytes,3,opt,name=idempotency_key,json=idempotencyKey,proto3" json:"idempotency_key,omitempty"`
	Visible        bool   `protobuf:"varint,4,opt,name=visible,proto3" json:"visible,omitempty"`
	Reference      string `protobuf:"bytes,5,opt,name=reference,proto3" json:"reference,omitempty"`
}

func (x *DecrementRequest) Reset() {
	*x = DecrementRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecrementRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecrementRequest) ProtoMessage() {}

func (x *DecrementRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecrementRequest.ProtoReflect.Descriptor instead.
func (*DecrementRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{2}
}

func (x *DecrementRequest) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

func (x *DecrementRequest) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *DecrementRequest) GetIdempotencyKey() string {
	if x != nil {
		return x.IdempotencyKey
	}
	return ""
}

func (x *DecrementRequest) GetVisible() bool {
	if x != nil {
		return x.Visible
	}
	return false
}

func (x *DecrementRequest) GetReference() string {
	if x != nil {
		return x.Reference
	}
	return ""
}

type DecrementResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewBalance int64 `protobuf:"varint,1,opt,name=new_balance,json=newBalance,proto3" json:"new_balance,omitempty"`
}

func (x *DecrementResponse) Reset() {
	*x = DecrementResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecrementResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecrementResponse) ProtoMessage() {}

func (x *DecrementResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecrementResponse.ProtoReflect.Descriptor instead.
func (*DecrementResponse) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{3}
}

func (x *DecrementResponse) GetNewBalance() int64 {
	if x != nil {
		return x.NewBalance
	}
	return 0
}

type CurrentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
}

func (x *CurrentRequest) Reset() {
	*x = CurrentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CurrentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CurrentRequest) ProtoMessage() {}

func (x *CurrentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CurrentRequest.ProtoReflect.Descriptor instead.
func (*CurrentRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{4}
}

func (x *CurrentRequest) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

type CurrentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentBalance int64 `protobuf:"varint,1,opt,name=current_balance,json=currentBalance,proto3" json:"current_balance,omitempty"`
}

func (x *CurrentResponse) Reset() {
	*x = CurrentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CurrentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CurrentResponse) ProtoMessage() {}

func (x *CurrentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CurrentResponse.ProtoReflect.Descriptor instead.
func (*CurrentResponse) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{5}
}

func (x *CurrentResponse) GetCurrentBalance() int64 {
	if x != nil {
		return x.CurrentBalance
	}
	return 0
}

type ListAdjustmentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
}

func (x *ListAdjustmentsRequest) Reset() {
	*x = ListAdjustmentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAdjustmentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAdjustmentsRequest) ProtoMessage() {}

func (x *ListAdjustmentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAdjustmentsRequest.ProtoReflect.Descriptor instead.
func (*ListAdjustmentsRequest) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{6}
}

func (x *ListAdjustmentsRequest) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

type Adjustment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Created   int64  `protobuf:"varint,2,opt,name=created,proto3" json:"created,omitempty"`
	Delta     int64  `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`
	Reference string `protobuf:"bytes,4,opt,name=reference,proto3" json:"reference,omitempty"`
}

func (x *Adjustment) Reset() {
	*x = Adjustment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Adjustment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Adjustment) ProtoMessage() {}

func (x *Adjustment) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Adjustment.ProtoReflect.Descriptor instead.
func (*Adjustment) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{7}
}

func (x *Adjustment) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Adjustment) GetCreated() int64 {
	if x != nil {
		return x.Created
	}
	return 0
}

func (x *Adjustment) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Adjustment) GetReference() string {
	if x != nil {
		return x.Reference
	}
	return ""
}

type ListAdjustmentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Adjustments []*Adjustment `protobuf:"bytes,1,rep,name=adjustments,proto3" json:"adjustments,omitempty"`
}

func (x *ListAdjustmentsResponse) Reset() {
	*x = ListAdjustmentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_balance_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAdjustmentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAdjustmentsResponse) ProtoMessage() {}

func (x *ListAdjustmentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balance_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAdjustmentsResponse.ProtoReflect.Descriptor instead.
func (*ListAdjustmentsResponse) Descriptor() ([]byte, []int) {
	return file_proto_balance_proto_rawDescGZIP(), []int{8}
}

func (x *ListAdjustmentsResponse) GetAdjustments() []*Adjustment {
	if x != nil {
		return x.Adjustments
	}
	return nil
}

var File_proto_balance_proto protoreflect.FileDescriptor

var file_proto_balance_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0xaa,
	0x01, 0x0a, 0x10, 0x49, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x64,
	0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79,
	0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x69, 0x73, 0x69, 0x62, 0x6c, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x76, 0x69, 0x73, 0x69, 0x62, 0x6c, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x34, 0x0a, 0x11, 0x49,
	0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6e, 0x65, 0x77, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x22, 0xaa, 0x01, 0x0a, 0x10, 0x44, 0x65, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x27, 0x0a,
	0x0f, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x69, 0x73, 0x69, 0x62, 0x6c,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x76, 0x69, 0x73, 0x69, 0x62, 0x6c, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x34,
	0x0a, 0x11, 0x44, 0x65, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6e, 0x65, 0x77, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x22, 0x31, 0x0a, 0x0e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3a, 0x0a, 0x0f, 0x43, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x22, 0x39, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x64, 0x6a, 0x75, 0x73,
	0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x22, 0x6a,
	0x0a, 0x0a, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x50, 0x0a, 0x17, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0b, 0x61, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x62, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x2e, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x0b, 0x61, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0xad, 0x02, 0x0a,
	0x07, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x09, 0x49, 0x6e, 0x63, 0x72,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e,
	0x49, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x63, 0x72, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x44,
	0x0a, 0x09, 0x44, 0x65, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x62, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x2e, 0x44, 0x65, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x07, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x12,
	0x17, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x64, 0x6a, 0x75,
	0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1f, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f,
	0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_balance_proto_rawDescOnce sync.Once
	file_proto_balance_proto_rawDescData = file_proto_balance_proto_rawDesc
)

func file_proto_balance_proto_rawDescGZIP() []byte {
	file_proto_balance_proto_rawDescOnce.Do(func() {
		file_proto_balance_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_balance_proto_rawDescData)
	})
	return file_proto_balance_proto_rawDescData
}

var file_proto_balance_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_balance_proto_goTypes = []interface{}{
	(*IncrementRequest)(nil),        // 0: balance.IncrementRequest
	(*IncrementResponse)(nil),       // 1: balance.IncrementResponse
	(*DecrementRequest)(nil),        // 2: balance.DecrementRequest
	(*DecrementResponse)(nil),       // 3: balance.DecrementResponse
	(*CurrentRequest)(nil),          // 4: balance.CurrentRequest
	(*CurrentResponse)(nil),         // 5: balance.CurrentResponse
	(*ListAdjustmentsRequest)(nil),  // 6: balance.ListAdjustmentsRequest
	(*Adjustment)(nil),              // 7: balance.Adjustment
	(*ListAdjustmentsResponse)(nil), // 8: balance.ListAdjustmentsResponse
}
var file_proto_balance_proto_depIdxs = []int32{
	7, // 0: balance.ListAdjustmentsResponse.adjustments:type_name -> balance.Adjustment
	0, // 1: balance.Balance.Increment:input_type -> balance.IncrementRequest
	2, // 2: balance.Balance.Decrement:input_type -> balance.DecrementRequest
	4, // 3: balance.Balance.Current:input_type -> balance.CurrentRequest
	6, // 4: balance.Balance.ListAdjustments:input_type -> balance.ListAdjustmentsRequest
	1, // 5: balance.Balance.Increment:output_type -> balance.IncrementResponse
	3, // 6: balance.Balance.Decrement:output_type -> balance.DecrementResponse
	5, // 7: balance.Balance.Current:output_type -> balance.CurrentResponse
	8, // 8: balance.Balance.ListAdjustments:output_type -> balance.ListAdjustmentsResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_balance_proto_init() }
func file_proto_balance_proto_init() {
	if File_proto_balance_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_balance_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IncrementRequest); i {
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
		file_proto_balance_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IncrementResponse); i {
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
		file_proto_balance_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DecrementRequest); i {
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
		file_proto_balance_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DecrementResponse); i {
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
		file_proto_balance_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CurrentRequest); i {
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
		file_proto_balance_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CurrentResponse); i {
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
		file_proto_balance_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAdjustmentsRequest); i {
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
		file_proto_balance_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Adjustment); i {
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
		file_proto_balance_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAdjustmentsResponse); i {
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
			RawDescriptor: file_proto_balance_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_balance_proto_goTypes,
		DependencyIndexes: file_proto_balance_proto_depIdxs,
		MessageInfos:      file_proto_balance_proto_msgTypes,
	}.Build()
	File_proto_balance_proto = out.File
	file_proto_balance_proto_rawDesc = nil
	file_proto_balance_proto_goTypes = nil
	file_proto_balance_proto_depIdxs = nil
}
