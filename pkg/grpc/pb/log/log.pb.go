// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: log/log.proto

package log

import (
	"github.com/selefra/selefra/pkg/grpc/pb/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StageType int32

const (
	StageType_STAGE_TYPE_INITIALIZING            StageType = 0
	StageType_STAGE_TYPE_PULL_INFRASTRUCTURE     StageType = 1
	StageType_STAGE_TYPE_INFRASTRUCTURE_ANALYSIS StageType = 2
)

// Enum value maps for StageType.
var (
	StageType_name = map[int32]string{
		0: "STAGE_TYPE_INITIALIZING",
		1: "STAGE_TYPE_PULL_INFRASTRUCTURE",
		2: "STAGE_TYPE_INFRASTRUCTURE_ANALYSIS",
	}
	StageType_value = map[string]int32{
		"STAGE_TYPE_INITIALIZING":            0,
		"STAGE_TYPE_PULL_INFRASTRUCTURE":     1,
		"STAGE_TYPE_INFRASTRUCTURE_ANALYSIS": 2,
	}
)

func (x StageType) Enum() *StageType {
	p := new(StageType)
	*p = x
	return p
}

func (x StageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StageType) Descriptor() protoreflect.EnumDescriptor {
	return file_log_log_proto_enumTypes[0].Descriptor()
}

func (StageType) Type() protoreflect.EnumType {
	return &file_log_log_proto_enumTypes[0]
}

func (x StageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StageType.Descriptor instead.
func (StageType) EnumDescriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{0}
}

type Status int32

const (
	Status_STATUS_SUCCESS Status = 0
	Status_STATUS_FAILED  Status = 1
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "STATUS_SUCCESS",
		1: "STATUS_FAILED",
	}
	Status_value = map[string]int32{
		"STATUS_SUCCESS": 0,
		"STATUS_FAILED":  1,
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
	return file_log_log_proto_enumTypes[1].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_log_log_proto_enumTypes[1]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{1}
}

type Level int32

const (
	Level_LEVEL_DEBUG Level = 0
	Level_LEVEL_INFO  Level = 1
	Level_LEVEL_WARN  Level = 2
	Level_LEVEL_ERROR Level = 3
	Level_LEVEL_FATAL Level = 4
)

// Enum value maps for Level.
var (
	Level_name = map[int32]string{
		0: "LEVEL_DEBUG",
		1: "LEVEL_INFO",
		2: "LEVEL_WARN",
		3: "LEVEL_ERROR",
		4: "LEVEL_FATAL",
	}
	Level_value = map[string]int32{
		"LEVEL_DEBUG": 0,
		"LEVEL_INFO":  1,
		"LEVEL_WARN":  2,
		"LEVEL_ERROR": 3,
		"LEVEL_FATAL": 4,
	}
)

func (x Level) Enum() *Level {
	p := new(Level)
	*p = x
	return p
}

func (x Level) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Level) Descriptor() protoreflect.EnumDescriptor {
	return file_log_log_proto_enumTypes[2].Descriptor()
}

func (Level) Type() protoreflect.EnumType {
	return &file_log_log_proto_enumTypes[2]
}

func (x Level) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Level.Descriptor instead.
func (Level) EnumDescriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{2}
}

type UploadLogStream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadLogStream) Reset() {
	*x = UploadLogStream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadLogStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadLogStream) ProtoMessage() {}

func (x *UploadLogStream) ProtoReflect() protoreflect.Message {
	mi := &file_log_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadLogStream.ProtoReflect.Descriptor instead.
func (*UploadLogStream) Descriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{0}
}

type UploadLogStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadLogStatus) Reset() {
	*x = UploadLogStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadLogStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadLogStatus) ProtoMessage() {}

func (x *UploadLogStatus) ProtoReflect() protoreflect.Message {
	mi := &file_log_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadLogStatus.ProtoReflect.Descriptor instead.
func (*UploadLogStatus) Descriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{1}
}

type UploadLogStream_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stage StageType `protobuf:"varint,1,opt,name=stage,proto3,enum=log.StageType" json:"stage,omitempty"`
	// log id, task uniq
	Index uint64 `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Msg   string `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Level Level  `protobuf:"varint,4,opt,name=level,proto3,enum=log.Level" json:"level,omitempty"`
	// log product time
	Time *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *UploadLogStream_Request) Reset() {
	*x = UploadLogStream_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_log_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadLogStream_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadLogStream_Request) ProtoMessage() {}

func (x *UploadLogStream_Request) ProtoReflect() protoreflect.Message {
	mi := &file_log_log_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadLogStream_Request.ProtoReflect.Descriptor instead.
func (*UploadLogStream_Request) Descriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{0, 0}
}

func (x *UploadLogStream_Request) GetStage() StageType {
	if x != nil {
		return x.Stage
	}
	return StageType_STAGE_TYPE_INITIALIZING
}

func (x *UploadLogStream_Request) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *UploadLogStream_Request) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *UploadLogStream_Request) GetLevel() Level {
	if x != nil {
		return x.Level
	}
	return Level_LEVEL_DEBUG
}

func (x *UploadLogStream_Request) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type UploadLogStream_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadLogStream_Response) Reset() {
	*x = UploadLogStream_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_log_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadLogStream_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadLogStream_Response) ProtoMessage() {}

func (x *UploadLogStream_Response) ProtoReflect() protoreflect.Message {
	mi := &file_log_log_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadLogStream_Response.ProtoReflect.Descriptor instead.
func (*UploadLogStream_Response) Descriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{0, 1}
}

type UploadLogStatus_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stage  StageType `protobuf:"varint,1,opt,name=stage,proto3,enum=log.StageType" json:"stage,omitempty"`
	Status Status    `protobuf:"varint,2,opt,name=status,proto3,enum=log.Status" json:"status,omitempty"`
	// status change time
	Time *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *UploadLogStatus_Request) Reset() {
	*x = UploadLogStatus_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_log_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadLogStatus_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadLogStatus_Request) ProtoMessage() {}

func (x *UploadLogStatus_Request) ProtoReflect() protoreflect.Message {
	mi := &file_log_log_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadLogStatus_Request.ProtoReflect.Descriptor instead.
func (*UploadLogStatus_Request) Descriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{1, 0}
}

func (x *UploadLogStatus_Request) GetStage() StageType {
	if x != nil {
		return x.Stage
	}
	return StageType_STAGE_TYPE_INITIALIZING
}

func (x *UploadLogStatus_Request) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_STATUS_SUCCESS
}

func (x *UploadLogStatus_Request) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type UploadLogStatus_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Diagnosis *common.Diagnosis `protobuf:"bytes,1,opt,name=diagnosis,proto3" json:"diagnosis,omitempty"`
}

func (x *UploadLogStatus_Response) Reset() {
	*x = UploadLogStatus_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_log_log_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadLogStatus_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadLogStatus_Response) ProtoMessage() {}

func (x *UploadLogStatus_Response) ProtoReflect() protoreflect.Message {
	mi := &file_log_log_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadLogStatus_Response.ProtoReflect.Descriptor instead.
func (*UploadLogStatus_Response) Descriptor() ([]byte, []int) {
	return file_log_log_proto_rawDescGZIP(), []int{1, 1}
}

func (x *UploadLogStatus_Response) GetDiagnosis() *common.Diagnosis {
	if x != nil {
		return x.Diagnosis
	}
	return nil
}

var File_log_log_proto protoreflect.FileDescriptor

var file_log_log_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x6c, 0x6f, 0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc9, 0x01, 0x0a, 0x0f, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x1a, 0xa9,
	0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x6c, 0x6f, 0x67, 0x2e,
	0x53, 0x74, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x20, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x1a, 0x0a, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xd5, 0x01, 0x0a, 0x0f, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x84, 0x01, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x53, 0x74, 0x61, 0x67,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x6c,
	0x6f, 0x67, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x1a, 0x3b, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a,
	0x09, 0x64, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x69, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f,
	0x73, 0x69, 0x73, 0x52, 0x09, 0x64, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x69, 0x73, 0x2a, 0x74,
	0x0a, 0x09, 0x53, 0x74, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x53,
	0x54, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x49, 0x54, 0x49, 0x41,
	0x4c, 0x49, 0x5a, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x22, 0x0a, 0x1e, 0x53, 0x54, 0x41, 0x47,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x55, 0x4c, 0x4c, 0x5f, 0x49, 0x4e, 0x46, 0x52,
	0x41, 0x53, 0x54, 0x52, 0x55, 0x43, 0x54, 0x55, 0x52, 0x45, 0x10, 0x01, 0x12, 0x26, 0x0a, 0x22,
	0x53, 0x54, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x46, 0x52, 0x41,
	0x53, 0x54, 0x52, 0x55, 0x43, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x41, 0x4e, 0x41, 0x4c, 0x59, 0x53,
	0x49, 0x53, 0x10, 0x02, 0x2a, 0x2f, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12,
	0x0a, 0x0e, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53,
	0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x01, 0x2a, 0x5a, 0x0a, 0x05, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x0f,
	0x0a, 0x0b, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x44, 0x45, 0x42, 0x55, 0x47, 0x10, 0x00, 0x12,
	0x0e, 0x0a, 0x0a, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x01, 0x12,
	0x0e, 0x0a, 0x0a, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x10, 0x02, 0x12,
	0x0f, 0x0a, 0x0b, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03,
	0x12, 0x0f, 0x0a, 0x0b, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x46, 0x41, 0x54, 0x41, 0x4c, 0x10,
	0x04, 0x32, 0xab, 0x01, 0x0a, 0x03, 0x4c, 0x6f, 0x67, 0x12, 0x52, 0x0a, 0x0f, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e, 0x6c,
	0x6f, 0x67, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6c, 0x6f, 0x67,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x50, 0x0a,
	0x0f, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x1c, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x0b, 0x5a, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_log_log_proto_rawDescOnce sync.Once
	file_log_log_proto_rawDescData = file_log_log_proto_rawDesc
)

func file_log_log_proto_rawDescGZIP() []byte {
	file_log_log_proto_rawDescOnce.Do(func() {
		file_log_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_log_log_proto_rawDescData)
	})
	return file_log_log_proto_rawDescData
}

var file_log_log_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_log_log_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_log_log_proto_goTypes = []interface{}{
	(StageType)(0),                   // 0: log.StageType
	(Status)(0),                      // 1: log.Status
	(Level)(0),                       // 2: log.Level
	(*UploadLogStream)(nil),          // 3: log.UploadLogStream
	(*UploadLogStatus)(nil),          // 4: log.UploadLogStatus
	(*UploadLogStream_Request)(nil),  // 5: log.UploadLogStream.Request
	(*UploadLogStream_Response)(nil), // 6: log.UploadLogStream.Response
	(*UploadLogStatus_Request)(nil),  // 7: log.UploadLogStatus.Request
	(*UploadLogStatus_Response)(nil), // 8: log.UploadLogStatus.Response
	(*timestamppb.Timestamp)(nil),    // 9: google.protobuf.Timestamp
	(*common.Diagnosis)(nil),         // 10: common.Diagnosis
}
var file_log_log_proto_depIdxs = []int32{
	0,  // 0: log.UploadLogStream.Request.stage:type_name -> log.StageType
	2,  // 1: log.UploadLogStream.Request.level:type_name -> log.Level
	9,  // 2: log.UploadLogStream.Request.time:type_name -> google.protobuf.Timestamp
	0,  // 3: log.UploadLogStatus.Request.stage:type_name -> log.StageType
	1,  // 4: log.UploadLogStatus.Request.status:type_name -> log.Status
	9,  // 5: log.UploadLogStatus.Request.time:type_name -> google.protobuf.Timestamp
	10, // 6: log.UploadLogStatus.Response.diagnosis:type_name -> common.Diagnosis
	5,  // 7: log.Log.UploadLogStream:input_type -> log.UploadLogStream.Request
	7,  // 8: log.Log.UploadLogStatus:input_type -> log.UploadLogStatus.Request
	6,  // 9: log.Log.UploadLogStream:output_type -> log.UploadLogStream.Response
	8,  // 10: log.Log.UploadLogStatus:output_type -> log.UploadLogStatus.Response
	9,  // [9:11] is the sub-list for method output_type
	7,  // [7:9] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_log_log_proto_init() }
func file_log_log_proto_init() {
	if File_log_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_log_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadLogStream); i {
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
		file_log_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadLogStatus); i {
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
		file_log_log_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadLogStream_Request); i {
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
		file_log_log_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadLogStream_Response); i {
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
		file_log_log_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadLogStatus_Request); i {
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
		file_log_log_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadLogStatus_Response); i {
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
			RawDescriptor: file_log_log_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_log_log_proto_goTypes,
		DependencyIndexes: file_log_log_proto_depIdxs,
		EnumInfos:         file_log_log_proto_enumTypes,
		MessageInfos:      file_log_log_proto_msgTypes,
	}.Build()
	File_log_log_proto = out.File
	file_log_log_proto_rawDesc = nil
	file_log_log_proto_goTypes = nil
	file_log_log_proto_depIdxs = nil
}
