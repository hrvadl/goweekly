// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: v1/translator/translator.proto

package translator

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

type TranslateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TranslateRequest) Reset() {
	*x = TranslateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_translator_translator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TranslateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateRequest) ProtoMessage() {}

func (x *TranslateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_translator_translator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateRequest.ProtoReflect.Descriptor instead.
func (*TranslateRequest) Descriptor() ([]byte, []int) {
	return file_v1_translator_translator_proto_rawDescGZIP(), []int{0}
}

func (x *TranslateRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type TranslateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TranslateResponse) Reset() {
	*x = TranslateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_translator_translator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TranslateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateResponse) ProtoMessage() {}

func (x *TranslateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_translator_translator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateResponse.ProtoReflect.Descriptor instead.
func (*TranslateResponse) Descriptor() ([]byte, []int) {
	return file_v1_translator_translator_proto_rawDescGZIP(), []int{1}
}

func (x *TranslateResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_v1_translator_translator_proto protoreflect.FileDescriptor

var file_v1_translator_translator_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x22,
	0x2c, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2d, 0x0a,
	0x11, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x62, 0x0a, 0x10,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4e, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68,
	0x72, 0x76, 0x61, 0x64, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x31, 0x2f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_v1_translator_translator_proto_rawDescOnce sync.Once
	file_v1_translator_translator_proto_rawDescData = file_v1_translator_translator_proto_rawDesc
)

func file_v1_translator_translator_proto_rawDescGZIP() []byte {
	file_v1_translator_translator_proto_rawDescOnce.Do(func() {
		file_v1_translator_translator_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_translator_translator_proto_rawDescData)
	})
	return file_v1_translator_translator_proto_rawDescData
}

var file_v1_translator_translator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v1_translator_translator_proto_goTypes = []interface{}{
	(*TranslateRequest)(nil),  // 0: translator.v1.TranslateRequest
	(*TranslateResponse)(nil), // 1: translator.v1.TranslateResponse
}
var file_v1_translator_translator_proto_depIdxs = []int32{
	0, // 0: translator.v1.TranslateService.Translate:input_type -> translator.v1.TranslateRequest
	1, // 1: translator.v1.TranslateService.Translate:output_type -> translator.v1.TranslateResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_translator_translator_proto_init() }
func file_v1_translator_translator_proto_init() {
	if File_v1_translator_translator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_translator_translator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TranslateRequest); i {
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
		file_v1_translator_translator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TranslateResponse); i {
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
			RawDescriptor: file_v1_translator_translator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_translator_translator_proto_goTypes,
		DependencyIndexes: file_v1_translator_translator_proto_depIdxs,
		MessageInfos:      file_v1_translator_translator_proto_msgTypes,
	}.Build()
	File_v1_translator_translator_proto = out.File
	file_v1_translator_translator_proto_rawDesc = nil
	file_v1_translator_translator_proto_goTypes = nil
	file_v1_translator_translator_proto_depIdxs = nil
}
