// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: static.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Features []*Feature `protobuf:"bytes,1,rep,name=features,proto3" json:"features,omitempty"`
}

func (x *AddRequest) Reset() {
	*x = AddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRequest) ProtoMessage() {}

func (x *AddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRequest.ProtoReflect.Descriptor instead.
func (*AddRequest) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{0}
}

func (x *AddRequest) GetFeatures() []*Feature {
	if x != nil {
		return x.Features
	}
	return nil
}

type AddReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *AddReply) Reset() {
	*x = AddReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddReply) ProtoMessage() {}

func (x *AddReply) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddReply.ProtoReflect.Descriptor instead.
func (*AddReply) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{1}
}

func (x *AddReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// The request message containing the user's name.
type SearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SearchType int32      `protobuf:"varint,1,opt,name=search_type,json=searchType,proto3" json:"search_type,omitempty"`
	Features   []*Feature `protobuf:"bytes,2,rep,name=features,proto3" json:"features,omitempty"`
	Topk       int32      `protobuf:"varint,3,opt,name=topk,proto3" json:"topk,omitempty"`
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[2]
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
	return file_static_proto_rawDescGZIP(), []int{2}
}

func (x *SearchRequest) GetSearchType() int32 {
	if x != nil {
		return x.SearchType
	}
	return 0
}

func (x *SearchRequest) GetFeatures() []*Feature {
	if x != nil {
		return x.Features
	}
	return nil
}

func (x *SearchRequest) GetTopk() int32 {
	if x != nil {
		return x.Topk
	}
	return 0
}

type Feature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blob string `protobuf:"bytes,1,opt,name=blob,proto3" json:"blob,omitempty"`
}

func (x *Feature) Reset() {
	*x = Feature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Feature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Feature) ProtoMessage() {}

func (x *Feature) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Feature.ProtoReflect.Descriptor instead.
func (*Feature) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{3}
}

func (x *Feature) GetBlob() string {
	if x != nil {
		return x.Blob
	}
	return ""
}

type Features struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Features []*Feature `protobuf:"bytes,1,rep,name=Features,proto3" json:"Features,omitempty"`
}

func (x *Features) Reset() {
	*x = Features{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Features) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Features) ProtoMessage() {}

func (x *Features) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Features.ProtoReflect.Descriptor instead.
func (*Features) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{4}
}

func (x *Features) GetFeatures() []*Feature {
	if x != nil {
		return x.Features
	}
	return nil
}

type SearchReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Featuresgroup []*Features `protobuf:"bytes,1,rep,name=featuresgroup,proto3" json:"featuresgroup,omitempty"`
	Distancetopk  *Vectors    `protobuf:"bytes,2,opt,name=distancetopk,proto3" json:"distancetopk,omitempty"`
}

func (x *SearchReply) Reset() {
	*x = SearchReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchReply) ProtoMessage() {}

func (x *SearchReply) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchReply.ProtoReflect.Descriptor instead.
func (*SearchReply) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{5}
}

func (x *SearchReply) GetFeaturesgroup() []*Features {
	if x != nil {
		return x.Featuresgroup
	}
	return nil
}

func (x *SearchReply) GetDistancetopk() *Vectors {
	if x != nil {
		return x.Distancetopk
	}
	return nil
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Row int32 `protobuf:"varint,1,opt,name=row,proto3" json:"row,omitempty"`
	Col int32 `protobuf:"varint,2,opt,name=col,proto3" json:"col,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteRequest) GetRow() int32 {
	if x != nil {
		return x.Row
	}
	return 0
}

func (x *DeleteRequest) GetCol() int32 {
	if x != nil {
		return x.Col
	}
	return 0
}

type DeleteReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *DeleteReply) Reset() {
	*x = DeleteReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReply) ProtoMessage() {}

func (x *DeleteReply) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReply.ProtoReflect.Descriptor instead.
func (*DeleteReply) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type Vector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vector []float32 `protobuf:"fixed32,1,rep,packed,name=vector,proto3" json:"vector,omitempty"`
}

func (x *Vector) Reset() {
	*x = Vector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vector) ProtoMessage() {}

func (x *Vector) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vector.ProtoReflect.Descriptor instead.
func (*Vector) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{8}
}

func (x *Vector) GetVector() []float32 {
	if x != nil {
		return x.Vector
	}
	return nil
}

type Vectors struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vectors []*Vector `protobuf:"bytes,1,rep,name=vectors,proto3" json:"vectors,omitempty"`
}

func (x *Vectors) Reset() {
	*x = Vectors{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vectors) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vectors) ProtoMessage() {}

func (x *Vectors) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vectors.ProtoReflect.Descriptor instead.
func (*Vectors) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{9}
}

func (x *Vectors) GetVectors() []*Vector {
	if x != nil {
		return x.Vectors
	}
	return nil
}

type VectorsGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vectorgroup []*Vectors `protobuf:"bytes,1,rep,name=vectorgroup,proto3" json:"vectorgroup,omitempty"`
}

func (x *VectorsGroup) Reset() {
	*x = VectorsGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_static_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VectorsGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VectorsGroup) ProtoMessage() {}

func (x *VectorsGroup) ProtoReflect() protoreflect.Message {
	mi := &file_static_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VectorsGroup.ProtoReflect.Descriptor instead.
func (*VectorsGroup) Descriptor() ([]byte, []int) {
	return file_static_proto_rawDescGZIP(), []int{10}
}

func (x *VectorsGroup) GetVectorgroup() []*Vectors {
	if x != nil {
		return x.Vectorgroup
	}
	return nil
}

var File_static_proto protoreflect.FileDescriptor

var file_static_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x0a, 0x41, 0x64,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x22, 0x0a, 0x08, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x77, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x31, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x6f, 0x70, 0x6b, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x74, 0x6f, 0x70, 0x6b, 0x22, 0x1d, 0x0a, 0x07, 0x46, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x22, 0x3d, 0x0a, 0x08, 0x46, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x08, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x08, 0x46, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x86, 0x01, 0x0a, 0x0b, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x3c, 0x0a, 0x0d, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x73, 0x52, 0x0d, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x39, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x74, 0x6f, 0x70, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x73, 0x52, 0x0c, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x74, 0x6f, 0x70, 0x6b, 0x22,
	0x33, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x72, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x72,
	0x6f, 0x77, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x63, 0x6f, 0x6c, 0x22, 0x25, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x20, 0x0a, 0x06, 0x56,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x06, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x22, 0x39, 0x0a,
	0x07, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x2e, 0x0a, 0x07, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52,
	0x07, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x22, 0x47, 0x0a, 0x0c, 0x56, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x73, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x37, 0x0a, 0x0b, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x52, 0x0b, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x32, 0x82, 0x02, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x58, 0x0a, 0x0a,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x49, 0x6e, 0x44, 0x42, 0x12, 0x1b, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x22, 0x07, 0x2f, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x3a, 0x01, 0x2a, 0x12, 0x48, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x18, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x0f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09, 0x22, 0x04, 0x2f, 0x61, 0x64, 0x64, 0x3a, 0x01, 0x2a,
	0x12, 0x54, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x22, 0x07, 0x2f, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_static_proto_rawDescOnce sync.Once
	file_static_proto_rawDescData = file_static_proto_rawDesc
)

func file_static_proto_rawDescGZIP() []byte {
	file_static_proto_rawDescOnce.Do(func() {
		file_static_proto_rawDescData = protoimpl.X.CompressGZIP(file_static_proto_rawDescData)
	})
	return file_static_proto_rawDescData
}

var file_static_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_static_proto_goTypes = []interface{}{
	(*AddRequest)(nil),    // 0: static_proto.AddRequest
	(*AddReply)(nil),      // 1: static_proto.AddReply
	(*SearchRequest)(nil), // 2: static_proto.SearchRequest
	(*Feature)(nil),       // 3: static_proto.Feature
	(*Features)(nil),      // 4: static_proto.Features
	(*SearchReply)(nil),   // 5: static_proto.SearchReply
	(*DeleteRequest)(nil), // 6: static_proto.DeleteRequest
	(*DeleteReply)(nil),   // 7: static_proto.DeleteReply
	(*Vector)(nil),        // 8: static_proto.Vector
	(*Vectors)(nil),       // 9: static_proto.Vectors
	(*VectorsGroup)(nil),  // 10: static_proto.VectorsGroup
}
var file_static_proto_depIdxs = []int32{
	3,  // 0: static_proto.AddRequest.features:type_name -> static_proto.Feature
	3,  // 1: static_proto.SearchRequest.features:type_name -> static_proto.Feature
	3,  // 2: static_proto.Features.Features:type_name -> static_proto.Feature
	4,  // 3: static_proto.SearchReply.featuresgroup:type_name -> static_proto.Features
	9,  // 4: static_proto.SearchReply.distancetopk:type_name -> static_proto.Vectors
	8,  // 5: static_proto.Vectors.vectors:type_name -> static_proto.Vector
	9,  // 6: static_proto.VectorsGroup.vectorgroup:type_name -> static_proto.Vectors
	2,  // 7: static_proto.Search.SearchInDB:input_type -> static_proto.SearchRequest
	0,  // 8: static_proto.Search.Add:input_type -> static_proto.AddRequest
	6,  // 9: static_proto.Search.Delete:input_type -> static_proto.DeleteRequest
	5,  // 10: static_proto.Search.SearchInDB:output_type -> static_proto.SearchReply
	1,  // 11: static_proto.Search.Add:output_type -> static_proto.AddReply
	7,  // 12: static_proto.Search.Delete:output_type -> static_proto.DeleteReply
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_static_proto_init() }
func file_static_proto_init() {
	if File_static_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_static_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRequest); i {
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
		file_static_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddReply); i {
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
		file_static_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_static_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Feature); i {
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
		file_static_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Features); i {
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
		file_static_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchReply); i {
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
		file_static_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_static_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteReply); i {
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
		file_static_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vector); i {
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
		file_static_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vectors); i {
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
		file_static_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VectorsGroup); i {
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
			RawDescriptor: file_static_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_static_proto_goTypes,
		DependencyIndexes: file_static_proto_depIdxs,
		MessageInfos:      file_static_proto_msgTypes,
	}.Build()
	File_static_proto = out.File
	file_static_proto_rawDesc = nil
	file_static_proto_goTypes = nil
	file_static_proto_depIdxs = nil
}
