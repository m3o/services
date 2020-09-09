// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.2
// source: proto/build.proto

package build

import (
	proto "github.com/golang/protobuf/proto"
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

// The status of the build:
type Status int32

const (
	Status_UNKNOWN        Status = 0  // Default / we don't know
	Status_REQUESTED      Status = 1  // Request has been received
	Status_REQUEST_FAILED Status = 2  // Something was wrong with the request
	Status_PENDING        Status = 3  // Request was accepted but hasn't been actioned yet
	Status_BUILDING       Status = 4  // Currently building
	Status_BUILD_FAILED   Status = 5  // Build was attempted but failed
	Status_BUILT          Status = 6  // Build succeeded
	Status_PUSHING        Status = 7  // Image is currently pushing
	Status_PUSHING_FAILED Status = 8  // Image failed to push
	Status_COMPLETE       Status = 9  // Image was successfully built and pushed
	Status_FAILED         Status = 10 // Failed for some other reason
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0:  "UNKNOWN",
		1:  "REQUESTED",
		2:  "REQUEST_FAILED",
		3:  "PENDING",
		4:  "BUILDING",
		5:  "BUILD_FAILED",
		6:  "BUILT",
		7:  "PUSHING",
		8:  "PUSHING_FAILED",
		9:  "COMPLETE",
		10: "FAILED",
	}
	Status_value = map[string]int32{
		"UNKNOWN":        0,
		"REQUESTED":      1,
		"REQUEST_FAILED": 2,
		"PENDING":        3,
		"BUILDING":       4,
		"BUILD_FAILED":   5,
		"BUILT":          6,
		"PUSHING":        7,
		"PUSHING_FAILED": 8,
		"COMPLETE":       9,
		"FAILED":         10,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_build_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_proto_build_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_build_proto_rawDescGZIP(), []int{0}
}

type CreateImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GitRepo   string `protobuf:"bytes,1,opt,name=gitRepo,proto3" json:"gitRepo,omitempty"`
	GitCommit string `protobuf:"bytes,2,opt,name=gitCommit,proto3" json:"gitCommit,omitempty"`
	ImageTag  string `protobuf:"bytes,3,opt,name=imageTag,proto3" json:"imageTag,omitempty"`
}

func (x *CreateImageRequest) Reset() {
	*x = CreateImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_build_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateImageRequest) ProtoMessage() {}

func (x *CreateImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_build_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateImageRequest.ProtoReflect.Descriptor instead.
func (*CreateImageRequest) Descriptor() ([]byte, []int) {
	return file_proto_build_proto_rawDescGZIP(), []int{0}
}

func (x *CreateImageRequest) GetGitRepo() string {
	if x != nil {
		return x.GitRepo
	}
	return ""
}

func (x *CreateImageRequest) GetGitCommit() string {
	if x != nil {
		return x.GitCommit
	}
	return ""
}

func (x *CreateImageRequest) GetImageTag() string {
	if x != nil {
		return x.ImageTag
	}
	return ""
}

type CreateImageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=build.Status" json:"status,omitempty"`
	Output string `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
	Error  string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CreateImageResponse) Reset() {
	*x = CreateImageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_build_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateImageResponse) ProtoMessage() {}

func (x *CreateImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_build_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateImageResponse.ProtoReflect.Descriptor instead.
func (*CreateImageResponse) Descriptor() ([]byte, []int) {
	return file_proto_build_proto_rawDescGZIP(), []int{1}
}

func (x *CreateImageResponse) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *CreateImageResponse) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

func (x *CreateImageResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_build_proto protoreflect.FileDescriptor

var file_proto_build_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x22, 0x68, 0x0a, 0x12, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x67, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x67, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x67, 0x69,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67,
	0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x54, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x54, 0x61, 0x67, 0x22, 0x6a, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x62, 0x75,
	0x69, 0x6c, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x2a, 0xab, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x52, 0x45, 0x51, 0x55,
	0x45, 0x53, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x45, 0x51, 0x55, 0x45,
	0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x50,
	0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x42, 0x55, 0x49, 0x4c,
	0x44, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x10, 0x0a, 0x0c, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x55, 0x49, 0x4c,
	0x54, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x55, 0x53, 0x48, 0x49, 0x4e, 0x47, 0x10, 0x07,
	0x12, 0x12, 0x0a, 0x0e, 0x50, 0x55, 0x53, 0x48, 0x49, 0x4e, 0x47, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x08, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45,
	0x10, 0x09, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0a, 0x32, 0x99,
	0x01, 0x0a, 0x05, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x46, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x19, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x48, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12,
	0x19, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x3b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_build_proto_rawDescOnce sync.Once
	file_proto_build_proto_rawDescData = file_proto_build_proto_rawDesc
)

func file_proto_build_proto_rawDescGZIP() []byte {
	file_proto_build_proto_rawDescOnce.Do(func() {
		file_proto_build_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_build_proto_rawDescData)
	})
	return file_proto_build_proto_rawDescData
}

var file_proto_build_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_build_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_build_proto_goTypes = []interface{}{
	(Status)(0),                 // 0: build.Status
	(*CreateImageRequest)(nil),  // 1: build.CreateImageRequest
	(*CreateImageResponse)(nil), // 2: build.CreateImageResponse
}
var file_proto_build_proto_depIdxs = []int32{
	0, // 0: build.CreateImageResponse.status:type_name -> build.Status
	1, // 1: build.Build.CreateImage:input_type -> build.CreateImageRequest
	1, // 2: build.Build.StreamImage:input_type -> build.CreateImageRequest
	2, // 3: build.Build.CreateImage:output_type -> build.CreateImageResponse
	2, // 4: build.Build.StreamImage:output_type -> build.CreateImageResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_build_proto_init() }
func file_proto_build_proto_init() {
	if File_proto_build_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_build_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateImageRequest); i {
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
		file_proto_build_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateImageResponse); i {
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
			RawDescriptor: file_proto_build_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_build_proto_goTypes,
		DependencyIndexes: file_proto_build_proto_depIdxs,
		EnumInfos:         file_proto_build_proto_enumTypes,
		MessageInfos:      file_proto_build_proto_msgTypes,
	}.Build()
	File_proto_build_proto = out.File
	file_proto_build_proto_rawDesc = nil
	file_proto_build_proto_goTypes = nil
	file_proto_build_proto_depIdxs = nil
}
