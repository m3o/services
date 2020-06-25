// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.6.1
// source: proto/onboarding/onboarding.proto

package go_micro_service_onboarding

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

type SendVerificationEmailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *SendVerificationEmailRequest) Reset() {
	*x = SendVerificationEmailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_onboarding_onboarding_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendVerificationEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendVerificationEmailRequest) ProtoMessage() {}

func (x *SendVerificationEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_onboarding_onboarding_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendVerificationEmailRequest.ProtoReflect.Descriptor instead.
func (*SendVerificationEmailRequest) Descriptor() ([]byte, []int) {
	return file_proto_onboarding_onboarding_proto_rawDescGZIP(), []int{0}
}

func (x *SendVerificationEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type SendVerificationEmailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendVerificationEmailResponse) Reset() {
	*x = SendVerificationEmailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_onboarding_onboarding_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendVerificationEmailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendVerificationEmailResponse) ProtoMessage() {}

func (x *SendVerificationEmailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_onboarding_onboarding_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendVerificationEmailResponse.ProtoReflect.Descriptor instead.
func (*SendVerificationEmailResponse) Descriptor() ([]byte, []int) {
	return file_proto_onboarding_onboarding_proto_rawDescGZIP(), []int{1}
}

type VerifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	// Email token that was received in an email.
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *VerifyRequest) Reset() {
	*x = VerifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_onboarding_onboarding_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyRequest) ProtoMessage() {}

func (x *VerifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_onboarding_onboarding_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyRequest.ProtoReflect.Descriptor instead.
func (*VerifyRequest) Descriptor() ([]byte, []int) {
	return file_proto_onboarding_onboarding_proto_rawDescGZIP(), []int{2}
}

func (x *VerifyRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *VerifyRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type VerifyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Auth token to be saved into `~/.micro`
	// For users who are already registered and paid,
	// the flow stops here.
	// For users who are yet to be registered
	// the token will be acquired in the `FinishOnboarding` step.
	AuthToken string `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
}

func (x *VerifyResponse) Reset() {
	*x = VerifyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_onboarding_onboarding_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyResponse) ProtoMessage() {}

func (x *VerifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_onboarding_onboarding_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyResponse.ProtoReflect.Descriptor instead.
func (*VerifyResponse) Descriptor() ([]byte, []int) {
	return file_proto_onboarding_onboarding_proto_rawDescGZIP(), []int{3}
}

func (x *VerifyResponse) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

type FinishOnboardingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	// The token has to be passed here too for identification purposes.
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	// This payment method ID is the one we got back from Stripe on the frontend (ie. `m3o.com/subscribe.html`)
	PaymentMethodId string `protobuf:"bytes,3,opt,name=paymentMethodId,proto3" json:"paymentMethodId,omitempty"`
}

func (x *FinishOnboardingRequest) Reset() {
	*x = FinishOnboardingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_onboarding_onboarding_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishOnboardingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishOnboardingRequest) ProtoMessage() {}

func (x *FinishOnboardingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_onboarding_onboarding_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishOnboardingRequest.ProtoReflect.Descriptor instead.
func (*FinishOnboardingRequest) Descriptor() ([]byte, []int) {
	return file_proto_onboarding_onboarding_proto_rawDescGZIP(), []int{4}
}

func (x *FinishOnboardingRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *FinishOnboardingRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *FinishOnboardingRequest) GetPaymentMethodId() string {
	if x != nil {
		return x.PaymentMethodId
	}
	return ""
}

type FinishOnboardingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken string `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	NameSpace string `protobuf:"bytes,2,opt,name=nameSpace,proto3" json:"nameSpace,omitempty"`
}

func (x *FinishOnboardingResponse) Reset() {
	*x = FinishOnboardingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_onboarding_onboarding_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishOnboardingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishOnboardingResponse) ProtoMessage() {}

func (x *FinishOnboardingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_onboarding_onboarding_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishOnboardingResponse.ProtoReflect.Descriptor instead.
func (*FinishOnboardingResponse) Descriptor() ([]byte, []int) {
	return file_proto_onboarding_onboarding_proto_rawDescGZIP(), []int{5}
}

func (x *FinishOnboardingResponse) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *FinishOnboardingResponse) GetNameSpace() string {
	if x != nil {
		return x.NameSpace
	}
	return ""
}

var File_proto_onboarding_onboarding_proto protoreflect.FileDescriptor

var file_proto_onboarding_onboarding_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69,
	0x6e, 0x67, 0x2f, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x67, 0x6f, 0x2e, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67,
	0x22, 0x34, 0x0a, 0x1c, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x1f, 0x0a, 0x1d, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3b, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2e, 0x0a, 0x0e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6f, 0x0a, 0x17, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4f, 0x6e,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x28, 0x0a, 0x0f, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x18, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4f,
	0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x32, 0x81, 0x03,
	0x0a, 0x0a, 0x4f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x8e, 0x01, 0x0a,
	0x15, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x39, 0x2e, 0x67, 0x6f, 0x2e, 0x6d, 0x69, 0x63, 0x72,
	0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x3a, 0x2e, 0x67, 0x6f, 0x2e, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a,
	0x06, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x2a, 0x2e, 0x67, 0x6f, 0x2e, 0x6d, 0x69, 0x63,
	0x72, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x67, 0x6f, 0x2e, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e,
	0x67, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x7f, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4f, 0x6e, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x69, 0x6e, 0x67, 0x12, 0x34, 0x2e, 0x67, 0x6f, 0x2e, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69,
	0x6e, 0x67, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4f, 0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x67, 0x6f, 0x2e,
	0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x6f, 0x6e,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4f,
	0x6e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_onboarding_onboarding_proto_rawDescOnce sync.Once
	file_proto_onboarding_onboarding_proto_rawDescData = file_proto_onboarding_onboarding_proto_rawDesc
)

func file_proto_onboarding_onboarding_proto_rawDescGZIP() []byte {
	file_proto_onboarding_onboarding_proto_rawDescOnce.Do(func() {
		file_proto_onboarding_onboarding_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_onboarding_onboarding_proto_rawDescData)
	})
	return file_proto_onboarding_onboarding_proto_rawDescData
}

var file_proto_onboarding_onboarding_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_onboarding_onboarding_proto_goTypes = []interface{}{
	(*SendVerificationEmailRequest)(nil),  // 0: go.micro.service.onboarding.SendVerificationEmailRequest
	(*SendVerificationEmailResponse)(nil), // 1: go.micro.service.onboarding.SendVerificationEmailResponse
	(*VerifyRequest)(nil),                 // 2: go.micro.service.onboarding.VerifyRequest
	(*VerifyResponse)(nil),                // 3: go.micro.service.onboarding.VerifyResponse
	(*FinishOnboardingRequest)(nil),       // 4: go.micro.service.onboarding.FinishOnboardingRequest
	(*FinishOnboardingResponse)(nil),      // 5: go.micro.service.onboarding.FinishOnboardingResponse
}
var file_proto_onboarding_onboarding_proto_depIdxs = []int32{
	0, // 0: go.micro.service.onboarding.Onboarding.SendVerificationEmail:input_type -> go.micro.service.onboarding.SendVerificationEmailRequest
	2, // 1: go.micro.service.onboarding.Onboarding.Verify:input_type -> go.micro.service.onboarding.VerifyRequest
	4, // 2: go.micro.service.onboarding.Onboarding.FinishOnboarding:input_type -> go.micro.service.onboarding.FinishOnboardingRequest
	1, // 3: go.micro.service.onboarding.Onboarding.SendVerificationEmail:output_type -> go.micro.service.onboarding.SendVerificationEmailResponse
	3, // 4: go.micro.service.onboarding.Onboarding.Verify:output_type -> go.micro.service.onboarding.VerifyResponse
	5, // 5: go.micro.service.onboarding.Onboarding.FinishOnboarding:output_type -> go.micro.service.onboarding.FinishOnboardingResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_onboarding_onboarding_proto_init() }
func file_proto_onboarding_onboarding_proto_init() {
	if File_proto_onboarding_onboarding_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_onboarding_onboarding_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendVerificationEmailRequest); i {
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
		file_proto_onboarding_onboarding_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendVerificationEmailResponse); i {
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
		file_proto_onboarding_onboarding_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyRequest); i {
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
		file_proto_onboarding_onboarding_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyResponse); i {
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
		file_proto_onboarding_onboarding_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishOnboardingRequest); i {
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
		file_proto_onboarding_onboarding_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishOnboardingResponse); i {
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
			RawDescriptor: file_proto_onboarding_onboarding_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_onboarding_onboarding_proto_goTypes,
		DependencyIndexes: file_proto_onboarding_onboarding_proto_depIdxs,
		MessageInfos:      file_proto_onboarding_onboarding_proto_msgTypes,
	}.Build()
	File_proto_onboarding_onboarding_proto = out.File
	file_proto_onboarding_onboarding_proto_rawDesc = nil
	file_proto_onboarding_onboarding_proto_goTypes = nil
	file_proto_onboarding_onboarding_proto_depIdxs = nil
}
