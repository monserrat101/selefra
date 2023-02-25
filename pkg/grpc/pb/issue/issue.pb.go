// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: issue.proto

package issue

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

type UploadIssueStream_Severity int32

const (
	UploadIssueStream_UNKNOWN       UploadIssueStream_Severity = 0
	UploadIssueStream_INFORMATIONAL UploadIssueStream_Severity = 1
	UploadIssueStream_LOW           UploadIssueStream_Severity = 2
	UploadIssueStream_MEDIUM        UploadIssueStream_Severity = 3
	UploadIssueStream_HIGH          UploadIssueStream_Severity = 4
	UploadIssueStream_CRITICAL      UploadIssueStream_Severity = 5
)

// Enum value maps for UploadIssueStream_Severity.
var (
	UploadIssueStream_Severity_name = map[int32]string{
		0: "UNKNOWN",
		1: "INFORMATIONAL",
		2: "LOW",
		3: "MEDIUM",
		4: "HIGH",
		5: "CRITICAL",
	}
	UploadIssueStream_Severity_value = map[string]int32{
		"UNKNOWN":       0,
		"INFORMATIONAL": 1,
		"LOW":           2,
		"MEDIUM":        3,
		"HIGH":          4,
		"CRITICAL":      5,
	}
)

func (x UploadIssueStream_Severity) Enum() *UploadIssueStream_Severity {
	p := new(UploadIssueStream_Severity)
	*p = x
	return p
}

func (x UploadIssueStream_Severity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UploadIssueStream_Severity) Descriptor() protoreflect.EnumDescriptor {
	return file_issue_proto_enumTypes[0].Descriptor()
}

func (UploadIssueStream_Severity) Type() protoreflect.EnumType {
	return &file_issue_proto_enumTypes[0]
}

func (x UploadIssueStream_Severity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UploadIssueStream_Severity.Descriptor instead.
func (UploadIssueStream_Severity) EnumDescriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 0}
}

type UploadIssueStream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadIssueStream) Reset() {
	*x = UploadIssueStream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream) ProtoMessage() {}

func (x *UploadIssueStream) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream.ProtoReflect.Descriptor instead.
func (*UploadIssueStream) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0}
}

type UploadIssueStream_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadIssueStream_Response) Reset() {
	*x = UploadIssueStream_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Response) ProtoMessage() {}

func (x *UploadIssueStream_Response) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Response.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Response) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 0}
}

type UploadIssueStream_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index    int32                       `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Rule     *UploadIssueStream_Rule     `protobuf:"bytes,2,opt,name=rule,proto3" json:"rule,omitempty"`
	Provider *UploadIssueStream_Provider `protobuf:"bytes,3,opt,name=provider,proto3" json:"provider,omitempty"`
	Module   *UploadIssueStream_Module   `protobuf:"bytes,4,opt,name=module,proto3" json:"module,omitempty"`
	// i do not know how to name it...
	Context *UploadIssueStream_Context `protobuf:"bytes,5,opt,name=context,proto3" json:"context,omitempty"`
}

func (x *UploadIssueStream_Request) Reset() {
	*x = UploadIssueStream_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Request) ProtoMessage() {}

func (x *UploadIssueStream_Request) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Request.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Request) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 1}
}

func (x *UploadIssueStream_Request) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *UploadIssueStream_Request) GetRule() *UploadIssueStream_Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

func (x *UploadIssueStream_Request) GetProvider() *UploadIssueStream_Provider {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *UploadIssueStream_Request) GetModule() *UploadIssueStream_Module {
	if x != nil {
		return x.Module
	}
	return nil
}

func (x *UploadIssueStream_Request) GetContext() *UploadIssueStream_Context {
	if x != nil {
		return x.Context
	}
	return nil
}

type UploadIssueStream_Context struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcTableNames []string `protobuf:"bytes,1,rep,name=src_table_names,json=srcTableNames,proto3" json:"src_table_names,omitempty"`
	// use which one pg db schema
	Schema string `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *UploadIssueStream_Context) Reset() {
	*x = UploadIssueStream_Context{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Context) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Context) ProtoMessage() {}

func (x *UploadIssueStream_Context) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Context.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Context) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 2}
}

func (x *UploadIssueStream_Context) GetSrcTableNames() []string {
	if x != nil {
		return x.SrcTableNames
	}
	return nil
}

func (x *UploadIssueStream_Context) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

type UploadIssueStream_Module struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name             string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Source           string   `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	DependenciesPath []string `protobuf:"bytes,3,rep,name=dependencies_path,json=dependenciesPath,proto3" json:"dependencies_path,omitempty"`
}

func (x *UploadIssueStream_Module) Reset() {
	*x = UploadIssueStream_Module{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Module) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Module) ProtoMessage() {}

func (x *UploadIssueStream_Module) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Module.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Module) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 3}
}

func (x *UploadIssueStream_Module) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UploadIssueStream_Module) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *UploadIssueStream_Module) GetDependenciesPath() []string {
	if x != nil {
		return x.DependenciesPath
	}
	return nil
}

type UploadIssueStream_Provider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Provider string `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Version  string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *UploadIssueStream_Provider) Reset() {
	*x = UploadIssueStream_Provider{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Provider) ProtoMessage() {}

func (x *UploadIssueStream_Provider) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Provider.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Provider) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 4}
}

func (x *UploadIssueStream_Provider) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UploadIssueStream_Provider) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *UploadIssueStream_Provider) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// rule's file block
type UploadIssueStream_Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// sql
	Query  string            `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
	Labels map[string]string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// rule's metadata
	Metadata *UploadIssueStream_Metadata `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Output   string                      `protobuf:"bytes,5,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *UploadIssueStream_Rule) Reset() {
	*x = UploadIssueStream_Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Rule) ProtoMessage() {}

func (x *UploadIssueStream_Rule) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Rule.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Rule) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 5}
}

func (x *UploadIssueStream_Rule) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UploadIssueStream_Rule) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *UploadIssueStream_Rule) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *UploadIssueStream_Rule) GetMetadata() *UploadIssueStream_Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *UploadIssueStream_Rule) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

// rule's metadata
type UploadIssueStream_Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Author      string                     `protobuf:"bytes,1,opt,name=author,proto3" json:"author,omitempty"`
	Description string                     `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Id          string                     `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Provider    string                     `protobuf:"bytes,4,opt,name=provider,proto3" json:"provider,omitempty"`
	Remediation string                     `protobuf:"bytes,5,opt,name=remediation,proto3" json:"remediation,omitempty"`
	Severity    UploadIssueStream_Severity `protobuf:"varint,6,opt,name=severity,proto3,enum=issue.UploadIssueStream_Severity" json:"severity,omitempty"`
	Tags        []string                   `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	Title       string                     `protobuf:"bytes,8,opt,name=title,proto3" json:"title,omitempty"`
}

func (x *UploadIssueStream_Metadata) Reset() {
	*x = UploadIssueStream_Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_issue_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIssueStream_Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIssueStream_Metadata) ProtoMessage() {}

func (x *UploadIssueStream_Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_issue_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIssueStream_Metadata.ProtoReflect.Descriptor instead.
func (*UploadIssueStream_Metadata) Descriptor() ([]byte, []int) {
	return file_issue_proto_rawDescGZIP(), []int{0, 6}
}

func (x *UploadIssueStream_Metadata) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *UploadIssueStream_Metadata) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UploadIssueStream_Metadata) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UploadIssueStream_Metadata) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *UploadIssueStream_Metadata) GetRemediation() string {
	if x != nil {
		return x.Remediation
	}
	return ""
}

func (x *UploadIssueStream_Metadata) GetSeverity() UploadIssueStream_Severity {
	if x != nil {
		return x.Severity
	}
	return UploadIssueStream_UNKNOWN
}

func (x *UploadIssueStream_Metadata) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *UploadIssueStream_Metadata) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

var File_issue_proto protoreflect.FileDescriptor

var file_issue_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x69,
	0x73, 0x73, 0x75, 0x65, 0x22, 0x8b, 0x09, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49,
	0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x1a, 0x0a, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x86, 0x02, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x31, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x12, 0x3d, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75,
	0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x06, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x06, 0x6d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x1a,
	0x49, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x73, 0x72,
	0x63, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x72, 0x63, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x1a, 0x61, 0x0a, 0x06, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x2b, 0x0a, 0x11, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x64, 0x65, 0x70,
	0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x50, 0x61, 0x74, 0x68, 0x1a, 0x54, 0x0a,
	0x08, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x1a, 0x85, 0x02, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x41, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x3d, 0x0a, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x73,
	0x73, 0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0xfb, 0x01, 0x0a, 0x08,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x20,
	0x0a, 0x0b, 0x72, 0x65, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x3d, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x21, 0x2e, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x53, 0x65, 0x76,
	0x65, 0x72, 0x69, 0x74, 0x79, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x22, 0x57, 0x0a, 0x08, 0x53, 0x65, 0x76,
	0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x4e, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x4c, 0x4f, 0x57, 0x10, 0x02, 0x12, 0x0a,
	0x0a, 0x06, 0x4d, 0x45, 0x44, 0x49, 0x55, 0x4d, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x49,
	0x47, 0x48, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c,
	0x10, 0x05, 0x32, 0x65, 0x0a, 0x05, 0x49, 0x73, 0x73, 0x75, 0x65, 0x12, 0x5c, 0x0a, 0x11, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x12, 0x20, 0x2e, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49,
	0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_issue_proto_rawDescOnce sync.Once
	file_issue_proto_rawDescData = file_issue_proto_rawDesc
)

func file_issue_proto_rawDescGZIP() []byte {
	file_issue_proto_rawDescOnce.Do(func() {
		file_issue_proto_rawDescData = protoimpl.X.CompressGZIP(file_issue_proto_rawDescData)
	})
	return file_issue_proto_rawDescData
}

var file_issue_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_issue_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_issue_proto_goTypes = []interface{}{
	(UploadIssueStream_Severity)(0),    // 0: issue.UploadIssueStream.Severity
	(*UploadIssueStream)(nil),          // 1: issue.UploadIssueStream
	(*UploadIssueStream_Response)(nil), // 2: issue.UploadIssueStream.Response
	(*UploadIssueStream_Request)(nil),  // 3: issue.UploadIssueStream.Request
	(*UploadIssueStream_Context)(nil),  // 4: issue.UploadIssueStream.Context
	(*UploadIssueStream_Module)(nil),   // 5: issue.UploadIssueStream.Module
	(*UploadIssueStream_Provider)(nil), // 6: issue.UploadIssueStream.Provider
	(*UploadIssueStream_Rule)(nil),     // 7: issue.UploadIssueStream.Rule
	(*UploadIssueStream_Metadata)(nil), // 8: issue.UploadIssueStream.Metadata
	nil,                                // 9: issue.UploadIssueStream.Rule.LabelsEntry
}
var file_issue_proto_depIdxs = []int32{
	7, // 0: issue.UploadIssueStream.Request.rule:type_name -> issue.UploadIssueStream.Rule
	6, // 1: issue.UploadIssueStream.Request.provider:type_name -> issue.UploadIssueStream.Provider
	5, // 2: issue.UploadIssueStream.Request.module:type_name -> issue.UploadIssueStream.Module
	4, // 3: issue.UploadIssueStream.Request.context:type_name -> issue.UploadIssueStream.Context
	9, // 4: issue.UploadIssueStream.Rule.labels:type_name -> issue.UploadIssueStream.Rule.LabelsEntry
	8, // 5: issue.UploadIssueStream.Rule.metadata:type_name -> issue.UploadIssueStream.Metadata
	0, // 6: issue.UploadIssueStream.Metadata.severity:type_name -> issue.UploadIssueStream.Severity
	3, // 7: issue.Issue.UploadIssueStream:input_type -> issue.UploadIssueStream.Request
	2, // 8: issue.Issue.UploadIssueStream:output_type -> issue.UploadIssueStream.Response
	8, // [8:9] is the sub-list for method output_type
	7, // [7:8] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_issue_proto_init() }
func file_issue_proto_init() {
	if File_issue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_issue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream); i {
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
		file_issue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Response); i {
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
		file_issue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Request); i {
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
		file_issue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Context); i {
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
		file_issue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Module); i {
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
		file_issue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Provider); i {
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
		file_issue_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Rule); i {
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
		file_issue_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIssueStream_Metadata); i {
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
			RawDescriptor: file_issue_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_issue_proto_goTypes,
		DependencyIndexes: file_issue_proto_depIdxs,
		EnumInfos:         file_issue_proto_enumTypes,
		MessageInfos:      file_issue_proto_msgTypes,
	}.Build()
	File_issue_proto = out.File
	file_issue_proto_rawDesc = nil
	file_issue_proto_goTypes = nil
	file_issue_proto_depIdxs = nil
}
