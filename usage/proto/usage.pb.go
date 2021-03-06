// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.5
// source: proto/usage.proto

package usage

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

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Detail     bool   `protobuf:"varint,2,opt,name=detail,proto3" json:"detail,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{0}
}

func (x *ReadRequest) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

func (x *ReadRequest) GetDetail() bool {
	if x != nil {
		return x.Detail
	}
	return false
}

type UsageHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiName string         `protobuf:"bytes,1,opt,name=api_name,json=apiName,proto3" json:"api_name,omitempty"`
	Records []*UsageRecord `protobuf:"bytes,2,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *UsageHistory) Reset() {
	*x = UsageHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsageHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageHistory) ProtoMessage() {}

func (x *UsageHistory) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageHistory.ProtoReflect.Descriptor instead.
func (*UsageHistory) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{1}
}

func (x *UsageHistory) GetApiName() string {
	if x != nil {
		return x.ApiName
	}
	return ""
}

func (x *UsageHistory) GetRecords() []*UsageRecord {
	if x != nil {
		return x.Records
	}
	return nil
}

type UsageRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date     int64 `protobuf:"varint,1,opt,name=date,proto3" json:"date,omitempty"`
	Requests int64 `protobuf:"varint,2,opt,name=requests,proto3" json:"requests,omitempty"`
}

func (x *UsageRecord) Reset() {
	*x = UsageRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsageRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageRecord) ProtoMessage() {}

func (x *UsageRecord) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageRecord.ProtoReflect.Descriptor instead.
func (*UsageRecord) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{2}
}

func (x *UsageRecord) GetDate() int64 {
	if x != nil {
		return x.Date
	}
	return 0
}

func (x *UsageRecord) GetRequests() int64 {
	if x != nil {
		return x.Requests
	}
	return 0
}

type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Usage map[string]*UsageHistory `protobuf:"bytes,1,rep,name=usage,proto3" json:"usage,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{3}
}

func (x *ReadResponse) GetUsage() map[string]*UsageHistory {
	if x != nil {
		return x.Usage
	}
	return nil
}

type SweepRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SweepRequest) Reset() {
	*x = SweepRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SweepRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SweepRequest) ProtoMessage() {}

func (x *SweepRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SweepRequest.ProtoReflect.Descriptor instead.
func (*SweepRequest) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{4}
}

type SweepResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SweepResponse) Reset() {
	*x = SweepResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SweepResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SweepResponse) ProtoMessage() {}

func (x *SweepResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SweepResponse.ProtoReflect.Descriptor instead.
func (*SweepResponse) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{5}
}

type DeleteCustomerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteCustomerRequest) Reset() {
	*x = DeleteCustomerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCustomerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCustomerRequest) ProtoMessage() {}

func (x *DeleteCustomerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCustomerRequest.ProtoReflect.Descriptor instead.
func (*DeleteCustomerRequest) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteCustomerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteCustomerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteCustomerResponse) Reset() {
	*x = DeleteCustomerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_usage_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCustomerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCustomerResponse) ProtoMessage() {}

func (x *DeleteCustomerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_usage_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCustomerResponse.ProtoReflect.Descriptor instead.
func (*DeleteCustomerResponse) Descriptor() ([]byte, []int) {
	return file_proto_usage_proto_rawDescGZIP(), []int{7}
}

var File_proto_usage_proto protoreflect.FileDescriptor

var file_proto_usage_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x22, 0x46, 0x0a, 0x0b, 0x52, 0x65,
	0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x22, 0x57, 0x0a, 0x0c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x48, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x70, 0x69, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x70, 0x69, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a,
	0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x3d, 0x0a, 0x0b, 0x55,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x93, 0x01, 0x0a, 0x0c, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x1a, 0x4d, 0x0a, 0x0a, 0x55, 0x73, 0x61, 0x67, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x0e, 0x0a, 0x0c, 0x53, 0x77, 0x65, 0x65, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x77, 0x65, 0x65, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x27, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x18, 0x0a, 0x16, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc1, 0x01, 0x0a, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x31,
	0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x12, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x34, 0x0a, 0x05, 0x53, 0x77, 0x65, 0x65, 0x70, 0x12, 0x13, 0x2e, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x53, 0x77, 0x65, 0x65, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x77, 0x65, 0x65, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x12, 0x1c, 0x2e, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x3b, 0x75, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_usage_proto_rawDescOnce sync.Once
	file_proto_usage_proto_rawDescData = file_proto_usage_proto_rawDesc
)

func file_proto_usage_proto_rawDescGZIP() []byte {
	file_proto_usage_proto_rawDescOnce.Do(func() {
		file_proto_usage_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_usage_proto_rawDescData)
	})
	return file_proto_usage_proto_rawDescData
}

var file_proto_usage_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_usage_proto_goTypes = []interface{}{
	(*ReadRequest)(nil),            // 0: usage.ReadRequest
	(*UsageHistory)(nil),           // 1: usage.UsageHistory
	(*UsageRecord)(nil),            // 2: usage.UsageRecord
	(*ReadResponse)(nil),           // 3: usage.ReadResponse
	(*SweepRequest)(nil),           // 4: usage.SweepRequest
	(*SweepResponse)(nil),          // 5: usage.SweepResponse
	(*DeleteCustomerRequest)(nil),  // 6: usage.DeleteCustomerRequest
	(*DeleteCustomerResponse)(nil), // 7: usage.DeleteCustomerResponse
	nil,                            // 8: usage.ReadResponse.UsageEntry
}
var file_proto_usage_proto_depIdxs = []int32{
	2, // 0: usage.UsageHistory.records:type_name -> usage.UsageRecord
	8, // 1: usage.ReadResponse.usage:type_name -> usage.ReadResponse.UsageEntry
	1, // 2: usage.ReadResponse.UsageEntry.value:type_name -> usage.UsageHistory
	0, // 3: usage.Usage.Read:input_type -> usage.ReadRequest
	4, // 4: usage.Usage.Sweep:input_type -> usage.SweepRequest
	6, // 5: usage.Usage.DeleteCustomer:input_type -> usage.DeleteCustomerRequest
	3, // 6: usage.Usage.Read:output_type -> usage.ReadResponse
	5, // 7: usage.Usage.Sweep:output_type -> usage.SweepResponse
	7, // 8: usage.Usage.DeleteCustomer:output_type -> usage.DeleteCustomerResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_usage_proto_init() }
func file_proto_usage_proto_init() {
	if File_proto_usage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_usage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRequest); i {
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
		file_proto_usage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsageHistory); i {
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
		file_proto_usage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsageRecord); i {
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
		file_proto_usage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadResponse); i {
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
		file_proto_usage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SweepRequest); i {
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
		file_proto_usage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SweepResponse); i {
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
		file_proto_usage_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCustomerRequest); i {
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
		file_proto_usage_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCustomerResponse); i {
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
			RawDescriptor: file_proto_usage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_usage_proto_goTypes,
		DependencyIndexes: file_proto_usage_proto_depIdxs,
		MessageInfos:      file_proto_usage_proto_msgTypes,
	}.Build()
	File_proto_usage_proto = out.File
	file_proto_usage_proto_rawDesc = nil
	file_proto_usage_proto_goTypes = nil
	file_proto_usage_proto_depIdxs = nil
}
