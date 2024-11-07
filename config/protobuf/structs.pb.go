// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.19.6
// source: config/protobuf/structs.proto

package protobuf

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

type RequestType int32

const (
	RequestType_ROUTINGTABLE RequestType = 0
	RequestType_RETRANSMIT   RequestType = 1
)

// Enum value maps for RequestType.
var (
	RequestType_name = map[int32]string{
		0: "ROUTINGTABLE",
		1: "RETRANSMIT",
	}
	RequestType_value = map[string]int32{
		"ROUTINGTABLE": 0,
		"RETRANSMIT":   1,
	}
)

func (x RequestType) Enum() *RequestType {
	p := new(RequestType)
	*p = x
	return p
}

func (x RequestType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RequestType) Descriptor() protoreflect.EnumDescriptor {
	return file_config_protobuf_structs_proto_enumTypes[0].Descriptor()
}

func (RequestType) Type() protoreflect.EnumType {
	return &file_config_protobuf_structs_proto_enumTypes[0]
}

func (x RequestType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RequestType.Descriptor instead.
func (RequestType) EnumDescriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{0}
}

type PlayerCommand int32

const (
	PlayerCommand_PLAY PlayerCommand = 0
	PlayerCommand_STOP PlayerCommand = 1
)

// Enum value maps for PlayerCommand.
var (
	PlayerCommand_name = map[int32]string{
		0: "PLAY",
		1: "STOP",
	}
	PlayerCommand_value = map[string]int32{
		"PLAY": 0,
		"STOP": 1,
	}
)

func (x PlayerCommand) Enum() *PlayerCommand {
	p := new(PlayerCommand)
	*p = x
	return p
}

func (x PlayerCommand) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PlayerCommand) Descriptor() protoreflect.EnumDescriptor {
	return file_config_protobuf_structs_proto_enumTypes[1].Descriptor()
}

func (PlayerCommand) Type() protoreflect.EnumType {
	return &file_config_protobuf_structs_proto_enumTypes[1]
}

func (x PlayerCommand) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PlayerCommand.Descriptor instead.
func (PlayerCommand) EnumDescriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{1}
}

type VideoFormat int32

const (
	VideoFormat_UNKNOWN VideoFormat = 0
	VideoFormat_MP4     VideoFormat = 1
	VideoFormat_MKV     VideoFormat = 2
	VideoFormat_AVI     VideoFormat = 3
	VideoFormat_WEBM    VideoFormat = 4
	VideoFormat_MJPEG   VideoFormat = 5
)

// Enum value maps for VideoFormat.
var (
	VideoFormat_name = map[int32]string{
		0: "UNKNOWN",
		1: "MP4",
		2: "MKV",
		3: "AVI",
		4: "WEBM",
		5: "MJPEG",
	}
	VideoFormat_value = map[string]int32{
		"UNKNOWN": 0,
		"MP4":     1,
		"MKV":     2,
		"AVI":     3,
		"WEBM":    4,
		"MJPEG":   5,
	}
)

func (x VideoFormat) Enum() *VideoFormat {
	p := new(VideoFormat)
	*p = x
	return p
}

func (x VideoFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VideoFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_config_protobuf_structs_proto_enumTypes[2].Descriptor()
}

func (VideoFormat) Type() protoreflect.EnumType {
	return &file_config_protobuf_structs_proto_enumTypes[2]
}

func (x VideoFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VideoFormat.Descriptor instead.
func (VideoFormat) EnumDescriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{2}
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender    string      `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Target    string      `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Type      RequestType `protobuf:"varint,3,opt,name=type,proto3,enum=video.RequestType" json:"type,omitempty"`
	Length    int32       `protobuf:"varint,4,opt,name=length,proto3" json:"length,omitempty"`
	Timestamp int32       `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Types that are assignable to Content:
	//
	//	*Header_ClientCommand
	//	*Header_ServerVideoChunk
	//	*Header_DistanceVectorRouting
	Content  isHeader_Content `protobuf_oneof:"content"`
	ClientIp string           `protobuf:"bytes,9,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	mi := &file_config_protobuf_structs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{0}
}

func (x *Header) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Header) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Header) GetType() RequestType {
	if x != nil {
		return x.Type
	}
	return RequestType_ROUTINGTABLE
}

func (x *Header) GetLength() int32 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *Header) GetTimestamp() int32 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (m *Header) GetContent() isHeader_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *Header) GetClientCommand() *ClientCommand {
	if x, ok := x.GetContent().(*Header_ClientCommand); ok {
		return x.ClientCommand
	}
	return nil
}

func (x *Header) GetServerVideoChunk() *ServerVideoChunk {
	if x, ok := x.GetContent().(*Header_ServerVideoChunk); ok {
		return x.ServerVideoChunk
	}
	return nil
}

func (x *Header) GetDistanceVectorRouting() *DistanceVectorRouting {
	if x, ok := x.GetContent().(*Header_DistanceVectorRouting); ok {
		return x.DistanceVectorRouting
	}
	return nil
}

func (x *Header) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

type isHeader_Content interface {
	isHeader_Content()
}

type Header_ClientCommand struct {
	ClientCommand *ClientCommand `protobuf:"bytes,6,opt,name=client_command,json=clientCommand,proto3,oneof"`
}

type Header_ServerVideoChunk struct {
	ServerVideoChunk *ServerVideoChunk `protobuf:"bytes,7,opt,name=server_video_chunk,json=serverVideoChunk,proto3,oneof"`
}

type Header_DistanceVectorRouting struct {
	DistanceVectorRouting *DistanceVectorRouting `protobuf:"bytes,8,opt,name=distance_vector_routing,json=distanceVectorRouting,proto3,oneof"`
}

func (*Header_ClientCommand) isHeader_Content() {}

func (*Header_ServerVideoChunk) isHeader_Content() {}

func (*Header_DistanceVectorRouting) isHeader_Content() {}

type ClientCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command               PlayerCommand `protobuf:"varint,1,opt,name=command,proto3,enum=video.PlayerCommand" json:"command,omitempty"`
	AdditionalInformation string        `protobuf:"bytes,2,opt,name=AdditionalInformation,proto3" json:"AdditionalInformation,omitempty"`
}

func (x *ClientCommand) Reset() {
	*x = ClientCommand{}
	mi := &file_config_protobuf_structs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientCommand) ProtoMessage() {}

func (x *ClientCommand) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientCommand.ProtoReflect.Descriptor instead.
func (*ClientCommand) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{1}
}

func (x *ClientCommand) GetCommand() PlayerCommand {
	if x != nil {
		return x.Command
	}
	return PlayerCommand_PLAY
}

func (x *ClientCommand) GetAdditionalInformation() string {
	if x != nil {
		return x.AdditionalInformation
	}
	return ""
}

type ServerVideoChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SequenceNumber int32       `protobuf:"varint,1,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"` // Order of this chunk in the video stream
	Timestamp      int64       `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                                 // Timestamp of the chunk in milliseconds
	Format         VideoFormat `protobuf:"varint,3,opt,name=format,proto3,enum=video.VideoFormat" json:"format,omitempty"`                // Video format (MP4, MKV, etc.)
	Data           []byte      `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`                                            // The raw video chunk data (as byte array)
	IsLastChunk    bool        `protobuf:"varint,5,opt,name=is_last_chunk,json=isLastChunk,proto3" json:"is_last_chunk,omitempty"`        // Flag to indicate if this is the last chunk
}

func (x *ServerVideoChunk) Reset() {
	*x = ServerVideoChunk{}
	mi := &file_config_protobuf_structs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerVideoChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerVideoChunk) ProtoMessage() {}

func (x *ServerVideoChunk) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerVideoChunk.ProtoReflect.Descriptor instead.
func (*ServerVideoChunk) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{2}
}

func (x *ServerVideoChunk) GetSequenceNumber() int32 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *ServerVideoChunk) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ServerVideoChunk) GetFormat() VideoFormat {
	if x != nil {
		return x.Format
	}
	return VideoFormat_UNKNOWN
}

func (x *ServerVideoChunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ServerVideoChunk) GetIsLastChunk() bool {
	if x != nil {
		return x.IsLastChunk
	}
	return false
}

// //////////////////////////////////////////////////////////////////////////////////////////////////
// Presence
// //////////////////////////////////////////////////////////////////////////////////////////////////
type Interface struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Interface) Reset() {
	*x = Interface{}
	mi := &file_config_protobuf_structs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Interface) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Interface) ProtoMessage() {}

func (x *Interface) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Interface.ProtoReflect.Descriptor instead.
func (*Interface) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{3}
}

func (x *Interface) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Interface) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type NextHop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NextNode *Interface `protobuf:"bytes,1,opt,name=next_node,json=nextNode,proto3" json:"next_node,omitempty"`
	Distance int32      `protobuf:"varint,2,opt,name=distance,proto3" json:"distance,omitempty"`
}

func (x *NextHop) Reset() {
	*x = NextHop{}
	mi := &file_config_protobuf_structs_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NextHop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextHop) ProtoMessage() {}

func (x *NextHop) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextHop.ProtoReflect.Descriptor instead.
func (*NextHop) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{4}
}

func (x *NextHop) GetNextNode() *Interface {
	if x != nil {
		return x.NextNode
	}
	return nil
}

func (x *NextHop) GetDistance() int32 {
	if x != nil {
		return x.Distance
	}
	return 0
}

type DistanceVectorRouting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries map[string]*NextHop `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Source  *Interface          `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *DistanceVectorRouting) Reset() {
	*x = DistanceVectorRouting{}
	mi := &file_config_protobuf_structs_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DistanceVectorRouting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DistanceVectorRouting) ProtoMessage() {}

func (x *DistanceVectorRouting) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DistanceVectorRouting.ProtoReflect.Descriptor instead.
func (*DistanceVectorRouting) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{5}
}

func (x *DistanceVectorRouting) GetEntries() map[string]*NextHop {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *DistanceVectorRouting) GetSource() *Interface {
	if x != nil {
		return x.Source
	}
	return nil
}

var File_config_protobuf_structs_proto protoreflect.FileDescriptor

var file_config_protobuf_structs_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x22, 0x9e, 0x03, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x12, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x3d, 0x0a, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x48, 0x00, 0x52,
	0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x47,
	0x0a, 0x12, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x76, 0x69, 0x64,
	0x65, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x48, 0x00, 0x52, 0x10, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x69, 0x64,
	0x65, 0x6f, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x56, 0x0a, 0x17, 0x64, 0x69, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x5f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x67, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x2e, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52,
	0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x15, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x70, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x42, 0x09, 0x0a, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x75, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x76, 0x69, 0x64, 0x65,
	0x6f, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52,
	0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x34, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xbd,
	0x01, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x2a, 0x0a, 0x06, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x76, 0x69, 0x64,
	0x65, 0x6f, 0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x06,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x73,
	0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0b, 0x69, 0x73, 0x4c, 0x61, 0x73, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x2f,
	0x0a, 0x09, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22,
	0x54, 0x0a, 0x07, 0x4e, 0x65, 0x78, 0x74, 0x48, 0x6f, 0x70, 0x12, 0x2d, 0x0a, 0x09, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52,
	0x08, 0x6e, 0x65, 0x78, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x69, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0xd2, 0x01, 0x0a, 0x15, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x12,
	0x43, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x29, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x45,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x4a,
	0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x48, 0x6f, 0x70, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x2f, 0x0a, 0x0b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x4f, 0x55,
	0x54, 0x49, 0x4e, 0x47, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x52,
	0x45, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x4d, 0x49, 0x54, 0x10, 0x01, 0x2a, 0x23, 0x0a, 0x0d, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x08, 0x0a, 0x04,
	0x50, 0x4c, 0x41, 0x59, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x01,
	0x2a, 0x4a, 0x0a, 0x0b, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03,
	0x4d, 0x50, 0x34, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x4b, 0x56, 0x10, 0x02, 0x12, 0x07,
	0x0a, 0x03, 0x41, 0x56, 0x49, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x45, 0x42, 0x4d, 0x10,
	0x04, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x4a, 0x50, 0x45, 0x47, 0x10, 0x05, 0x42, 0x30, 0x5a, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x64, 0x72, 0x69,
	0x67, 0x6f, 0x30, 0x33, 0x34, 0x35, 0x2f, 0x65, 0x73, 0x72, 0x2d, 0x74, 0x70, 0x32, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_protobuf_structs_proto_rawDescOnce sync.Once
	file_config_protobuf_structs_proto_rawDescData = file_config_protobuf_structs_proto_rawDesc
)

func file_config_protobuf_structs_proto_rawDescGZIP() []byte {
	file_config_protobuf_structs_proto_rawDescOnce.Do(func() {
		file_config_protobuf_structs_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_protobuf_structs_proto_rawDescData)
	})
	return file_config_protobuf_structs_proto_rawDescData
}

var file_config_protobuf_structs_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_config_protobuf_structs_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_config_protobuf_structs_proto_goTypes = []any{
	(RequestType)(0),              // 0: video.RequestType
	(PlayerCommand)(0),            // 1: video.PlayerCommand
	(VideoFormat)(0),              // 2: video.VideoFormat
	(*Header)(nil),                // 3: video.Header
	(*ClientCommand)(nil),         // 4: video.ClientCommand
	(*ServerVideoChunk)(nil),      // 5: video.ServerVideoChunk
	(*Interface)(nil),             // 6: video.Interface
	(*NextHop)(nil),               // 7: video.NextHop
	(*DistanceVectorRouting)(nil), // 8: video.DistanceVectorRouting
	nil,                           // 9: video.DistanceVectorRouting.EntriesEntry
}
var file_config_protobuf_structs_proto_depIdxs = []int32{
	0,  // 0: video.Header.type:type_name -> video.RequestType
	4,  // 1: video.Header.client_command:type_name -> video.ClientCommand
	5,  // 2: video.Header.server_video_chunk:type_name -> video.ServerVideoChunk
	8,  // 3: video.Header.distance_vector_routing:type_name -> video.DistanceVectorRouting
	1,  // 4: video.ClientCommand.command:type_name -> video.PlayerCommand
	2,  // 5: video.ServerVideoChunk.format:type_name -> video.VideoFormat
	6,  // 6: video.NextHop.next_node:type_name -> video.Interface
	9,  // 7: video.DistanceVectorRouting.entries:type_name -> video.DistanceVectorRouting.EntriesEntry
	6,  // 8: video.DistanceVectorRouting.source:type_name -> video.Interface
	7,  // 9: video.DistanceVectorRouting.EntriesEntry.value:type_name -> video.NextHop
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_config_protobuf_structs_proto_init() }
func file_config_protobuf_structs_proto_init() {
	if File_config_protobuf_structs_proto != nil {
		return
	}
	file_config_protobuf_structs_proto_msgTypes[0].OneofWrappers = []any{
		(*Header_ClientCommand)(nil),
		(*Header_ServerVideoChunk)(nil),
		(*Header_DistanceVectorRouting)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_protobuf_structs_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_protobuf_structs_proto_goTypes,
		DependencyIndexes: file_config_protobuf_structs_proto_depIdxs,
		EnumInfos:         file_config_protobuf_structs_proto_enumTypes,
		MessageInfos:      file_config_protobuf_structs_proto_msgTypes,
	}.Build()
	File_config_protobuf_structs_proto = out.File
	file_config_protobuf_structs_proto_rawDesc = nil
	file_config_protobuf_structs_proto_goTypes = nil
	file_config_protobuf_structs_proto_depIdxs = nil
}