// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.19.6
// source: config/send_message/message.proto

package send_message

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

type TypeInteraction int32

const (
	TypeInteraction_NEIGHBOR    TypeInteraction = 0
	TypeInteraction_TRANSMITION TypeInteraction = 1
	TypeInteraction_IS_ALIVE    TypeInteraction = 2
	TypeInteraction_FAILED      TypeInteraction = 3
)

// Enum value maps for TypeInteraction.
var (
	TypeInteraction_name = map[int32]string{
		0: "NEIGHBOR",
		1: "TRANSMITION",
		2: "IS_ALIVE",
		3: "FAILED",
	}
	TypeInteraction_value = map[string]int32{
		"NEIGHBOR":    0,
		"TRANSMITION": 1,
		"IS_ALIVE":    2,
		"FAILED":      3,
	}
)

func (x TypeInteraction) Enum() *TypeInteraction {
	p := new(TypeInteraction)
	*p = x
	return p
}

func (x TypeInteraction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TypeInteraction) Descriptor() protoreflect.EnumDescriptor {
	return file_config_send_message_message_proto_enumTypes[0].Descriptor()
}

func (TypeInteraction) Type() protoreflect.EnumType {
	return &file_config_send_message_message_proto_enumTypes[0]
}

func (x TypeInteraction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TypeInteraction.Descriptor instead.
func (TypeInteraction) EnumDescriptor() ([]byte, []int) {
	return file_config_send_message_message_proto_rawDescGZIP(), []int{0}
}

// Message for sending a chunk of a video or other media
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string          `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"` // The actual content/message being sent
	Target  string          `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`   // The intended recipient or target of the message
	Type    TypeInteraction `protobuf:"varint,3,opt,name=type,proto3,enum=message.TypeInteraction" json:"type,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_config_send_message_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_config_send_message_message_proto_msgTypes[0]
	if x != nil {
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
	return file_config_send_message_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Message) GetType() TypeInteraction {
	if x != nil {
		return x.Type
	}
	return TypeInteraction_NEIGHBOR
}

var File_config_send_message_message_proto protoreflect.FileDescriptor

var file_config_send_message_message_proto_rawDesc = []byte{
	0x0a, 0x21, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x69, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x4a, 0x0a, 0x0f, 0x54, 0x79, 0x70, 0x65, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x45,
	0x49, 0x47, 0x48, 0x42, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x52, 0x41, 0x4e,
	0x53, 0x4d, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x53, 0x5f,
	0x41, 0x4c, 0x49, 0x56, 0x45, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x03, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x72, 0x6f, 0x64, 0x72, 0x69, 0x67, 0x6f, 0x30, 0x33, 0x34, 0x35, 0x2f, 0x65, 0x73,
	0x72, 0x2d, 0x74, 0x70, 0x32, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x73, 0x65, 0x6e,
	0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_config_send_message_message_proto_rawDescOnce sync.Once
	file_config_send_message_message_proto_rawDescData = file_config_send_message_message_proto_rawDesc
)

func file_config_send_message_message_proto_rawDescGZIP() []byte {
	file_config_send_message_message_proto_rawDescOnce.Do(func() {
		file_config_send_message_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_send_message_message_proto_rawDescData)
	})
	return file_config_send_message_message_proto_rawDescData
}

var file_config_send_message_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_config_send_message_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_config_send_message_message_proto_goTypes = []any{
	(TypeInteraction)(0), // 0: message.TypeInteraction
	(*Message)(nil),      // 1: message.Message
}
var file_config_send_message_message_proto_depIdxs = []int32{
	0, // 0: message.Message.type:type_name -> message.TypeInteraction
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_config_send_message_message_proto_init() }
func file_config_send_message_message_proto_init() {
	if File_config_send_message_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_send_message_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_send_message_message_proto_goTypes,
		DependencyIndexes: file_config_send_message_message_proto_depIdxs,
		EnumInfos:         file_config_send_message_message_proto_enumTypes,
		MessageInfos:      file_config_send_message_message_proto_msgTypes,
	}.Build()
	File_config_send_message_message_proto = out.File
	file_config_send_message_message_proto_rawDesc = nil
	file_config_send_message_message_proto_goTypes = nil
	file_config_send_message_message_proto_depIdxs = nil
}
