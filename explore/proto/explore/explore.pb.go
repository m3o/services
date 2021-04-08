// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: proto/explore/explore.proto

package explore

import (
	proto "github.com/golang/protobuf/proto"
	registry "github.com/micro/micro/v3/proto/registry"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SearchTerm string `protobuf:"bytes,1,opt,name=searchTerm,proto3" json:"searchTerm,omitempty"`
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_explore_explore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_explore_explore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequest.ProtoReflect.Descriptor instead.
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return file_proto_explore_explore_proto_rawDescGZIP(), []int{0}
}

func (x *SearchRequest) GetSearchTerm() string {
	if x != nil {
		return x.SearchTerm
	}
	return ""
}

type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Readme      string            `protobuf:"bytes,1,opt,name=readme,proto3" json:"readme,omitempty"`
	OpenAPIJSON string            `protobuf:"bytes,2,opt,name=openAPIJSON,proto3" json:"openAPIJSON,omitempty"`
	Service     *registry.Service `protobuf:"bytes,3,opt,name=service,proto3" json:"service,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_explore_explore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_proto_explore_explore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_proto_explore_explore_proto_rawDescGZIP(), []int{1}
}

func (x *Service) GetReadme() string {
	if x != nil {
		return x.Readme
	}
	return ""
}

func (x *Service) GetOpenAPIJSON() string {
	if x != nil {
		return x.OpenAPIJSON
	}
	return ""
}

func (x *Service) GetService() *registry.Service {
	if x != nil {
		return x.Service
	}
	return nil
}

type SearchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Services []*Service `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
}

func (x *SearchResponse) Reset() {
	*x = SearchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_explore_explore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchResponse) ProtoMessage() {}

func (x *SearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_explore_explore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchResponse.ProtoReflect.Descriptor instead.
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return file_proto_explore_explore_proto_rawDescGZIP(), []int{2}
}

func (x *SearchResponse) GetServices() []*Service {
	if x != nil {
		return x.Services
	}
	return nil
}

type SaveMetaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=serviceName,proto3" json:"serviceName,omitempty"`
	// contents of the readme file
	Readme string `protobuf:"bytes,2,opt,name=readme,proto3" json:"readme,omitempty"`
	// marshalled JSON OpenAPI spec
	OpenAPIJSON string `protobuf:"bytes,3,opt,name=openAPIJSON,proto3" json:"openAPIJSON,omitempty"`
}

func (x *SaveMetaRequest) Reset() {
	*x = SaveMetaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_explore_explore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveMetaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveMetaRequest) ProtoMessage() {}

func (x *SaveMetaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_explore_explore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveMetaRequest.ProtoReflect.Descriptor instead.
func (*SaveMetaRequest) Descriptor() ([]byte, []int) {
	return file_proto_explore_explore_proto_rawDescGZIP(), []int{3}
}

func (x *SaveMetaRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *SaveMetaRequest) GetReadme() string {
	if x != nil {
		return x.Readme
	}
	return ""
}

func (x *SaveMetaRequest) GetOpenAPIJSON() string {
	if x != nil {
		return x.OpenAPIJSON
	}
	return ""
}

type SaveMetaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SaveMetaResponse) Reset() {
	*x = SaveMetaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_explore_explore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveMetaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveMetaResponse) ProtoMessage() {}

func (x *SaveMetaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_explore_explore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveMetaResponse.ProtoReflect.Descriptor instead.
func (*SaveMetaResponse) Descriptor() ([]byte, []int) {
	return file_proto_explore_explore_proto_rawDescGZIP(), []int{4}
}

var File_proto_explore_explore_proto protoreflect.FileDescriptor

var file_proto_explore_explore_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x70, 0x6c, 0x6f, 0x72, 0x65, 0x2f,
	0x65, 0x78, 0x70, 0x6c, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x65,
	0x78, 0x70, 0x6c, 0x6f, 0x72, 0x65, 0x1a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6d,
	0x70, 0x6f, 0x72, 0x74, 0x73, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54,
	0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x70, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x70, 0x65, 0x6e,
	0x41, 0x50, 0x49, 0x4a, 0x53, 0x4f, 0x4e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f,
	0x70, 0x65, 0x6e, 0x41, 0x50, 0x49, 0x4a, 0x53, 0x4f, 0x4e, 0x12, 0x2b, 0x0a, 0x07, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x07,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x3e, 0x0a, 0x0e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x08, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x78,
	0x70, 0x6c, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x08, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x6d, 0x0a, 0x0f, 0x53, 0x61, 0x76, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65,
	0x61, 0x64, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x70, 0x65, 0x6e, 0x41, 0x50, 0x49, 0x4a,
	0x53, 0x4f, 0x4e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x70, 0x65, 0x6e, 0x41,
	0x50, 0x49, 0x4a, 0x53, 0x4f, 0x4e, 0x22, 0x12, 0x0a, 0x10, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65,
	0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x89, 0x01, 0x0a, 0x07, 0x45,
	0x78, 0x70, 0x6c, 0x6f, 0x72, 0x65, 0x12, 0x3b, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x12, 0x16, 0x2e, 0x65, 0x78, 0x70, 0x6c, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x65, 0x78, 0x70, 0x6c, 0x6f,
	0x72, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x08, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x12,
	0x18, 0x2e, 0x65, 0x78, 0x70, 0x6c, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65,
	0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x6c,
	0x6f, 0x72, 0x65, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_explore_explore_proto_rawDescOnce sync.Once
	file_proto_explore_explore_proto_rawDescData = file_proto_explore_explore_proto_rawDesc
)

func file_proto_explore_explore_proto_rawDescGZIP() []byte {
	file_proto_explore_explore_proto_rawDescOnce.Do(func() {
		file_proto_explore_explore_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_explore_explore_proto_rawDescData)
	})
	return file_proto_explore_explore_proto_rawDescData
}

var file_proto_explore_explore_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_explore_explore_proto_goTypes = []interface{}{
	(*SearchRequest)(nil),    // 0: explore.SearchRequest
	(*Service)(nil),          // 1: explore.Service
	(*SearchResponse)(nil),   // 2: explore.SearchResponse
	(*SaveMetaRequest)(nil),  // 3: explore.SaveMetaRequest
	(*SaveMetaResponse)(nil), // 4: explore.SaveMetaResponse
	(*registry.Service)(nil), // 5: registry.Service
}
var file_proto_explore_explore_proto_depIdxs = []int32{
	5, // 0: explore.Service.service:type_name -> registry.Service
	1, // 1: explore.SearchResponse.services:type_name -> explore.Service
	0, // 2: explore.Explore.Search:input_type -> explore.SearchRequest
	3, // 3: explore.Explore.SaveMeta:input_type -> explore.SaveMetaRequest
	2, // 4: explore.Explore.Search:output_type -> explore.SearchResponse
	4, // 5: explore.Explore.SaveMeta:output_type -> explore.SaveMetaResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_explore_explore_proto_init() }
func file_proto_explore_explore_proto_init() {
	if File_proto_explore_explore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_explore_explore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchRequest); i {
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
		file_proto_explore_explore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service); i {
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
		file_proto_explore_explore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchResponse); i {
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
		file_proto_explore_explore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveMetaRequest); i {
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
		file_proto_explore_explore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveMetaResponse); i {
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
			RawDescriptor: file_proto_explore_explore_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_explore_explore_proto_goTypes,
		DependencyIndexes: file_proto_explore_explore_proto_depIdxs,
		MessageInfos:      file_proto_explore_explore_proto_msgTypes,
	}.Build()
	File_proto_explore_explore_proto = out.File
	file_proto_explore_explore_proto_rawDesc = nil
	file_proto_explore_explore_proto_goTypes = nil
	file_proto_explore_explore_proto_depIdxs = nil
}
