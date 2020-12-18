// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.7.0
// source: basic_arrays.proto

package protos

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

type BasicArrays struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ABool    []bool    `protobuf:"varint,1,rep,packed,name=a_bool,json=aBool,proto3" json:"a_bool,omitempty"`
	AByte    []byte    `protobuf:"bytes,2,opt,name=a_byte,json=aByte,proto3" json:"a_byte,omitempty"`
	AInt16   []int32   `protobuf:"varint,3,rep,packed,name=a_int16,json=aInt16,proto3" json:"a_int16,omitempty"`
	AUint16  []uint32  `protobuf:"varint,4,rep,packed,name=a_uint16,json=aUint16,proto3" json:"a_uint16,omitempty"`
	AInt32   []int32   `protobuf:"varint,5,rep,packed,name=a_int32,json=aInt32,proto3" json:"a_int32,omitempty"`
	AUint32  []uint32  `protobuf:"varint,6,rep,packed,name=a_uint32,json=aUint32,proto3" json:"a_uint32,omitempty"`
	AInt64   []int64   `protobuf:"varint,7,rep,packed,name=a_int64,json=aInt64,proto3" json:"a_int64,omitempty"`
	AUint64  []uint64  `protobuf:"varint,8,rep,packed,name=a_uint64,json=aUint64,proto3" json:"a_uint64,omitempty"`
	AFloat32 []float32 `protobuf:"fixed32,9,rep,packed,name=a_float32,json=aFloat32,proto3" json:"a_float32,omitempty"`
	AFloat64 []float64 `protobuf:"fixed64,10,rep,packed,name=a_float64,json=aFloat64,proto3" json:"a_float64,omitempty"`
	AString  []string  `protobuf:"bytes,11,rep,name=a_string,json=aString,proto3" json:"a_string,omitempty"`
	AGuid    []byte    `protobuf:"bytes,12,opt,name=a_guid,json=aGuid,proto3" json:"a_guid,omitempty"`
}

func (x *BasicArrays) Reset() {
	*x = BasicArrays{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_arrays_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicArrays) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicArrays) ProtoMessage() {}

func (x *BasicArrays) ProtoReflect() protoreflect.Message {
	mi := &file_basic_arrays_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicArrays.ProtoReflect.Descriptor instead.
func (*BasicArrays) Descriptor() ([]byte, []int) {
	return file_basic_arrays_proto_rawDescGZIP(), []int{0}
}

func (x *BasicArrays) GetABool() []bool {
	if x != nil {
		return x.ABool
	}
	return nil
}

func (x *BasicArrays) GetAByte() []byte {
	if x != nil {
		return x.AByte
	}
	return nil
}

func (x *BasicArrays) GetAInt16() []int32 {
	if x != nil {
		return x.AInt16
	}
	return nil
}

func (x *BasicArrays) GetAUint16() []uint32 {
	if x != nil {
		return x.AUint16
	}
	return nil
}

func (x *BasicArrays) GetAInt32() []int32 {
	if x != nil {
		return x.AInt32
	}
	return nil
}

func (x *BasicArrays) GetAUint32() []uint32 {
	if x != nil {
		return x.AUint32
	}
	return nil
}

func (x *BasicArrays) GetAInt64() []int64 {
	if x != nil {
		return x.AInt64
	}
	return nil
}

func (x *BasicArrays) GetAUint64() []uint64 {
	if x != nil {
		return x.AUint64
	}
	return nil
}

func (x *BasicArrays) GetAFloat32() []float32 {
	if x != nil {
		return x.AFloat32
	}
	return nil
}

func (x *BasicArrays) GetAFloat64() []float64 {
	if x != nil {
		return x.AFloat64
	}
	return nil
}

func (x *BasicArrays) GetAString() []string {
	if x != nil {
		return x.AString
	}
	return nil
}

func (x *BasicArrays) GetAGuid() []byte {
	if x != nil {
		return x.AGuid
	}
	return nil
}

type TestInt32Array struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A []int32 `protobuf:"varint,1,rep,packed,name=a,proto3" json:"a,omitempty"`
}

func (x *TestInt32Array) Reset() {
	*x = TestInt32Array{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_arrays_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestInt32Array) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestInt32Array) ProtoMessage() {}

func (x *TestInt32Array) ProtoReflect() protoreflect.Message {
	mi := &file_basic_arrays_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestInt32Array.ProtoReflect.Descriptor instead.
func (*TestInt32Array) Descriptor() ([]byte, []int) {
	return file_basic_arrays_proto_rawDescGZIP(), []int{1}
}

func (x *TestInt32Array) GetA() []int32 {
	if x != nil {
		return x.A
	}
	return nil
}

var File_basic_arrays_proto protoreflect.FileDescriptor

var file_basic_arrays_proto_rawDesc = []byte{
	0x0a, 0x12, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x72, 0x72, 0x61, 0x79, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3, 0x02, 0x0a, 0x0b,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x72, 0x72, 0x61, 0x79, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x61,
	0x5f, 0x62, 0x6f, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x03, 0x28, 0x08, 0x52, 0x05, 0x61, 0x42, 0x6f,
	0x6f, 0x6c, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x61, 0x42, 0x79, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x5f, 0x69,
	0x6e, 0x74, 0x31, 0x36, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x06, 0x61, 0x49, 0x6e, 0x74,
	0x31, 0x36, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x5f, 0x75, 0x69, 0x6e, 0x74, 0x31, 0x36, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x55, 0x69, 0x6e, 0x74, 0x31, 0x36, 0x12, 0x17, 0x0a,
	0x07, 0x61, 0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0x05, 0x20, 0x03, 0x28, 0x05, 0x52, 0x06,
	0x61, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x5f, 0x75, 0x69, 0x6e, 0x74,
	0x33, 0x32, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x55, 0x69, 0x6e, 0x74, 0x33,
	0x32, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x03, 0x52, 0x06, 0x61, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x5f,
	0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x08, 0x20, 0x03, 0x28, 0x04, 0x52, 0x07, 0x61, 0x55,
	0x69, 0x6e, 0x74, 0x36, 0x34, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74,
	0x33, 0x32, 0x18, 0x09, 0x20, 0x03, 0x28, 0x02, 0x52, 0x08, 0x61, 0x46, 0x6c, 0x6f, 0x61, 0x74,
	0x33, 0x32, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x18,
	0x0a, 0x20, 0x03, 0x28, 0x01, 0x52, 0x08, 0x61, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x12,
	0x19, 0x0a, 0x08, 0x61, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x0b, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x5f,
	0x67, 0x75, 0x69, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x61, 0x47, 0x75, 0x69,
	0x64, 0x22, 0x1e, 0x0a, 0x0e, 0x54, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x41, 0x72,
	0x72, 0x61, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x01,
	0x61, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x32, 0x30, 0x30, 0x73, 0x63, 0x2f, 0x62, 0x65, 0x62, 0x6f, 0x70, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x64, 0x61, 0x74, 0x61, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_basic_arrays_proto_rawDescOnce sync.Once
	file_basic_arrays_proto_rawDescData = file_basic_arrays_proto_rawDesc
)

func file_basic_arrays_proto_rawDescGZIP() []byte {
	file_basic_arrays_proto_rawDescOnce.Do(func() {
		file_basic_arrays_proto_rawDescData = protoimpl.X.CompressGZIP(file_basic_arrays_proto_rawDescData)
	})
	return file_basic_arrays_proto_rawDescData
}

var file_basic_arrays_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_basic_arrays_proto_goTypes = []interface{}{
	(*BasicArrays)(nil),    // 0: proto.BasicArrays
	(*TestInt32Array)(nil), // 1: proto.TestInt32Array
}
var file_basic_arrays_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_basic_arrays_proto_init() }
func file_basic_arrays_proto_init() {
	if File_basic_arrays_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_basic_arrays_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicArrays); i {
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
		file_basic_arrays_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestInt32Array); i {
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
			RawDescriptor: file_basic_arrays_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_basic_arrays_proto_goTypes,
		DependencyIndexes: file_basic_arrays_proto_depIdxs,
		MessageInfos:      file_basic_arrays_proto_msgTypes,
	}.Build()
	File_basic_arrays_proto = out.File
	file_basic_arrays_proto_rawDesc = nil
	file_basic_arrays_proto_goTypes = nil
	file_basic_arrays_proto_depIdxs = nil
}
