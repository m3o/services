// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: proto/status/status.proto

package api_status

import (
	proto "github.com/golang/protobuf/proto"
	proto1 "github.com/micro/go-micro/v2/api/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
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

var File_proto_status_status_proto protoreflect.FileDescriptor

var file_proto_status_status_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x70, 0x69,
	0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x32, 0x35, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2b, 0x0a, 0x04, 0x43, 0x61,
	0x6c, 0x6c, 0x12, 0x0f, 0x2e, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_status_status_proto_goTypes = []interface{}{
	(*proto1.Request)(nil),  // 0: go.api.Request
	(*proto1.Response)(nil), // 1: go.api.Response
}
var file_proto_status_status_proto_depIdxs = []int32{
	0, // 0: api.status.Status.Call:input_type -> go.api.Request
	1, // 1: api.status.Status.Call:output_type -> go.api.Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_status_status_proto_init() }
func file_proto_status_status_proto_init() {
	if File_proto_status_status_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_status_status_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_status_status_proto_goTypes,
		DependencyIndexes: file_proto_status_status_proto_depIdxs,
	}.Build()
	File_proto_status_status_proto = out.File
	file_proto_status_status_proto_rawDesc = nil
	file_proto_status_status_proto_goTypes = nil
	file_proto_status_status_proto_depIdxs = nil
}
