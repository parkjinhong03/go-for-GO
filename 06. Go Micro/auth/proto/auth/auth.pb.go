// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.1
// 	protoc        v3.11.4
// source: proto/auth/auth.proto

package examples_blog_service_auth

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

type UserIdExistRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" validate:"required,minLength=4,maxLength=16"`
	Authorization string `protobuf:"bytes,2,opt,name=authorization,proto3" json:"authorization,omitempty"`
}

func (x *UserIdExistRequest) Reset() {
	*x = UserIdExistRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIdExistRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIdExistRequest) ProtoMessage() {}

func (x *UserIdExistRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIdExistRequest.ProtoReflect.Descriptor instead.
func (*UserIdExistRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_auth_proto_rawDescGZIP(), []int{0}
}

func (x *UserIdExistRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserIdExistRequest) GetAuthorization() string {
	if x != nil {
		return x.Authorization
	}
	return ""
}

type UserIdExistResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        int64  `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Message       string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Authorization string `protobuf:"bytes,3,opt,name=authorization,proto3" json:"authorization,omitempty"`
}

func (x *UserIdExistResponse) SetResponse(status int, msg string) {
	x.Status = int64(status)
	x.Message = msg
}

func (x *UserIdExistResponse) Reset() {
	*x = UserIdExistResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIdExistResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIdExistResponse) ProtoMessage() {}

func (x *UserIdExistResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIdExistResponse.ProtoReflect.Descriptor instead.
func (*UserIdExistResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_auth_proto_rawDescGZIP(), []int{1}
}

func (x *UserIdExistResponse) GetStatus() int64 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *UserIdExistResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UserIdExistResponse) GetAuthorization() string {
	if x != nil {
		return x.Authorization
	}
	return ""
}

type CreateAuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" validate:"required,minLength=4,maxLength=16"`
	UserPw       string `protobuf:"bytes,2,opt,name=user_pw,json=userPw,proto3" json:"user_pw,omitempty" validate:"required,minLength=4,maxLength=16"`
	Name         string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty" validate:"required,minLength=2,maxLength=4"`
	PhoneNumber  string `protobuf:"bytes,4,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty" validate:"required,strLength=11"`
	Email        string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty" validate:"required,email,maxLength=30"`
	Introduction string `protobuf:"bytes,6,opt,name=introduction,proto3" json:"introduction,omitempty" validate:"maxLength=100"`
	Authorization string `protobuf:"bytes,7,opt,name=authorization,proto3" json:"authorization,omitempty"`
}

func (x *CreateAuthRequest) Reset() {
	*x = CreateAuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAuthRequest) ProtoMessage() {}

func (x *CreateAuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAuthRequest.ProtoReflect.Descriptor instead.
func (*CreateAuthRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_auth_proto_rawDescGZIP(), []int{2}
}

func (x *CreateAuthRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *CreateAuthRequest) GetUserPw() string {
	if x != nil {
		return x.UserPw
	}
	return ""
}

func (x *CreateAuthRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAuthRequest) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *CreateAuthRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateAuthRequest) GetIntroduction() string {
	if x != nil {
		return x.Introduction
	}
	return ""
}

func (x *CreateAuthRequest) GetAuthorization() string {
	if x != nil {
		return x.Authorization
	}
	return ""
}

type CreateAuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  int64  `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CreateAuthResponse) Reset() {
	*x = CreateAuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAuthResponse) ProtoMessage() {}

func (x *CreateAuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAuthResponse.ProtoReflect.Descriptor instead.
func (*CreateAuthResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_auth_proto_rawDescGZIP(), []int{3}
}

func (x *CreateAuthResponse) GetStatus() int64 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *CreateAuthResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Say string `protobuf:"bytes,1,opt,name=say,proto3" json:"say,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_proto_auth_auth_proto_rawDescGZIP(), []int{4}
}

func (x *Message) GetSay() string {
	if x != nil {
		return x.Say
	}
	return ""
}

var File_proto_auth_auth_proto protoreflect.FileDescriptor

var file_proto_auth_auth_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x22, 0x53, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x6d, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x45, 0x78, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xdc, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70,
	0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x50, 0x77, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x22, 0x0a, 0x0c,
	0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x24, 0x0a, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x46, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1b,
	0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x61, 0x79, 0x32, 0xee, 0x01, 0x0a, 0x04,
	0x41, 0x75, 0x74, 0x68, 0x12, 0x6d, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75,
	0x74, 0x68, 0x12, 0x2d, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x62, 0x6c,
	0x6f, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2e, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x62, 0x6c, 0x6f,
	0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x77, 0x0a, 0x12, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x49, 0x66, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x45, 0x78, 0x69, 0x73, 0x74, 0x12, 0x2e, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_auth_auth_proto_rawDescOnce sync.Once
	file_proto_auth_auth_proto_rawDescData = file_proto_auth_auth_proto_rawDesc
)

func file_proto_auth_auth_proto_rawDescGZIP() []byte {
	file_proto_auth_auth_proto_rawDescOnce.Do(func() {
		file_proto_auth_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_auth_auth_proto_rawDescData)
	})
	return file_proto_auth_auth_proto_rawDescData
}

var file_proto_auth_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_auth_auth_proto_goTypes = []interface{}{
	(*UserIdExistRequest)(nil),  // 0: examples.blog.service.auth.UserIdExistRequest
	(*UserIdExistResponse)(nil), // 1: examples.blog.service.auth.UserIdExistResponse
	(*CreateAuthRequest)(nil),   // 2: examples.blog.service.auth.CreateAuthRequest
	(*CreateAuthResponse)(nil),  // 3: examples.blog.service.auth.CreateAuthResponse
	(*Message)(nil),             // 4: examples.blog.service.auth.Message
}
var file_proto_auth_auth_proto_depIdxs = []int32{
	2, // 0: examples.blog.service.auth.Auth.CreateAuth:input_type -> examples.blog.service.auth.CreateAuthRequest
	0, // 1: examples.blog.service.auth.Auth.CheckIfUserIdExist:input_type -> examples.blog.service.auth.UserIdExistRequest
	3, // 2: examples.blog.service.auth.Auth.CreateAuth:output_type -> examples.blog.service.auth.CreateAuthResponse
	1, // 3: examples.blog.service.auth.Auth.CheckIfUserIdExist:output_type -> examples.blog.service.auth.UserIdExistResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_auth_auth_proto_init() }
func file_proto_auth_auth_proto_init() {
	if File_proto_auth_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_auth_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIdExistRequest); i {
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
		file_proto_auth_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIdExistResponse); i {
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
		file_proto_auth_auth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAuthRequest); i {
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
		file_proto_auth_auth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAuthResponse); i {
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
		file_proto_auth_auth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_proto_auth_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_auth_auth_proto_goTypes,
		DependencyIndexes: file_proto_auth_auth_proto_depIdxs,
		MessageInfos:      file_proto_auth_auth_proto_msgTypes,
	}.Build()
	File_proto_auth_auth_proto = out.File
	file_proto_auth_auth_proto_rawDesc = nil
	file_proto_auth_auth_proto_goTypes = nil
	file_proto_auth_auth_proto_depIdxs = nil
}
