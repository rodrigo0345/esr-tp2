// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
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
	// just to check if the client is still alive
	RequestType_HEARTBEAT   RequestType = 2
	RequestType_BOOTSTRAPER RequestType = 3
	RequestType_FAILED      RequestType = 4
	RequestType_CLIENT_PING RequestType = 5
)

// Enum value maps for RequestType.
var (
	RequestType_name = map[int32]string{
		0: "ROUTINGTABLE",
		1: "RETRANSMIT",
		2: "HEARTBEAT",
		3: "BOOTSTRAPER",
		4: "FAILED",
		5: "CLIENT_PING",
	}
	RequestType_value = map[string]int32{
		"ROUTINGTABLE": 0,
		"RETRANSMIT":   1,
		"HEARTBEAT":    2,
		"BOOTSTRAPER":  3,
		"FAILED":       4,
		"CLIENT_PING":  5,
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
	Target    []string    `protobuf:"bytes,2,rep,name=target,proto3" json:"target,omitempty"`
	Type      RequestType `protobuf:"varint,3,opt,name=type,proto3,enum=video.RequestType" json:"type,omitempty"`
	Length    int32       `protobuf:"varint,4,opt,name=length,proto3" json:"length,omitempty"`
	Timestamp int64       `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Types that are assignable to Content:
	//
	//	*Header_ClientCommand
	//	*Header_ServerVideoChunk
	//	*Header_DistanceVectorRouting
	//	*Header_BootstraperResult
	Content    isHeader_Content `protobuf_oneof:"content"`
	ClientPort string           `protobuf:"bytes,9,opt,name=client_port,json=clientPort,proto3" json:"client_port,omitempty"`
	// coloquei no header para n ter de haver tanto processamento de header e body
	RequestedVideo string `protobuf:"bytes,10,opt,name=RequestedVideo,proto3" json:"RequestedVideo,omitempty"`
	Path           string `protobuf:"bytes,12,opt,name=path,proto3" json:"path,omitempty"`
	SeqId          int64  `protobuf:"varint,13,opt,name=seq_id,json=seqId,proto3" json:"seq_id,omitempty"`
	Hops           int32  `protobuf:"varint,14,opt,name=hops,proto3" json:"hops,omitempty"`
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

func (x *Header) GetTarget() []string {
	if x != nil {
		return x.Target
	}
	return nil
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

func (x *Header) GetTimestamp() int64 {
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

func (x *Header) GetBootstraperResult() *BootstraperResult {
	if x, ok := x.GetContent().(*Header_BootstraperResult); ok {
		return x.BootstraperResult
	}
	return nil
}

func (x *Header) GetClientPort() string {
	if x != nil {
		return x.ClientPort
	}
	return ""
}

func (x *Header) GetRequestedVideo() string {
	if x != nil {
		return x.RequestedVideo
	}
	return ""
}

func (x *Header) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Header) GetSeqId() int64 {
	if x != nil {
		return x.SeqId
	}
	return 0
}

func (x *Header) GetHops() int32 {
	if x != nil {
		return x.Hops
	}
	return 0
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

type Header_BootstraperResult struct {
	BootstraperResult *BootstraperResult `protobuf:"bytes,11,opt,name=bootstraper_result,json=bootstraperResult,proto3,oneof"`
}

func (*Header_ClientCommand) isHeader_Content() {}

func (*Header_ServerVideoChunk) isHeader_Content() {}

func (*Header_DistanceVectorRouting) isHeader_Content() {}

func (*Header_BootstraperResult) isHeader_Content() {}

type BootstraperResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Neighbors []string `protobuf:"bytes,1,rep,name=neighbors,proto3" json:"neighbors,omitempty"`
}

func (x *BootstraperResult) Reset() {
	*x = BootstraperResult{}
	mi := &file_config_protobuf_structs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BootstraperResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstraperResult) ProtoMessage() {}

func (x *BootstraperResult) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use BootstraperResult.ProtoReflect.Descriptor instead.
func (*BootstraperResult) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{1}
}

func (x *BootstraperResult) GetNeighbors() []string {
	if x != nil {
		return x.Neighbors
	}
	return nil
}

type ClientCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command               PlayerCommand `protobuf:"varint,1,opt,name=command,proto3,enum=video.PlayerCommand" json:"command,omitempty"`
	AdditionalInformation string        `protobuf:"bytes,2,opt,name=AdditionalInformation,proto3" json:"AdditionalInformation,omitempty"`
}

func (x *ClientCommand) Reset() {
	*x = ClientCommand{}
	mi := &file_config_protobuf_structs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientCommand) ProtoMessage() {}

func (x *ClientCommand) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ClientCommand.ProtoReflect.Descriptor instead.
func (*ClientCommand) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{2}
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

	SequenceNumber      int32       `protobuf:"varint,1,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"`                   // Order of this chunk in the video stream
	Timestamp           int64       `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                                                   // Timestamp of the chunk in milliseconds
	Format              VideoFormat `protobuf:"varint,3,opt,name=format,proto3,enum=video.VideoFormat" json:"format,omitempty"`                                  // Video format (MP4, MKV, etc.)
	Data                []byte      `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`                                                              // The raw video chunk data (as byte array)
	IsLastChunk         bool        `protobuf:"varint,5,opt,name=is_last_chunk,json=isLastChunk,proto3" json:"is_last_chunk,omitempty"`                          // Flag to indicate if this is the last chunk
	BitrateKbps         int32       `protobuf:"varint,6,opt,name=bitrate_kbps,json=bitrateKbps,proto3" json:"bitrate_kbps,omitempty"`                            // Bitrate in kilobits per second
	Width               int32       `protobuf:"varint,7,opt,name=width,proto3" json:"width,omitempty"`                                                           // Width of the video in pixels
	Height              int32       `protobuf:"varint,8,opt,name=height,proto3" json:"height,omitempty"`                                                         // Height of the video in pixels
	Fps                 float64     `protobuf:"fixed64,9,opt,name=fps,proto3" json:"fps,omitempty"`                                                              // Frames per second
	LatencyMs           int64       `protobuf:"varint,10,opt,name=latency_ms,json=latencyMs,proto3" json:"latency_ms,omitempty"`                                 // Latency in milliseconds (if measurable)
	CurrentBitrateKbps  int32       `protobuf:"varint,11,opt,name=current_bitrate_kbps,json=currentBitrateKbps,proto3" json:"current_bitrate_kbps,omitempty"`    // Current bitrate in kilobits per second
	PreviousBitrateKbps int32       `protobuf:"varint,12,opt,name=previous_bitrate_kbps,json=previousBitrateKbps,proto3" json:"previous_bitrate_kbps,omitempty"` // Previous bitrate before adaptation
}

func (x *ServerVideoChunk) Reset() {
	*x = ServerVideoChunk{}
	mi := &file_config_protobuf_structs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerVideoChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerVideoChunk) ProtoMessage() {}

func (x *ServerVideoChunk) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ServerVideoChunk.ProtoReflect.Descriptor instead.
func (*ServerVideoChunk) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{3}
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

func (x *ServerVideoChunk) GetBitrateKbps() int32 {
	if x != nil {
		return x.BitrateKbps
	}
	return 0
}

func (x *ServerVideoChunk) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *ServerVideoChunk) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *ServerVideoChunk) GetFps() float64 {
	if x != nil {
		return x.Fps
	}
	return 0
}

func (x *ServerVideoChunk) GetLatencyMs() int64 {
	if x != nil {
		return x.LatencyMs
	}
	return 0
}

func (x *ServerVideoChunk) GetCurrentBitrateKbps() int32 {
	if x != nil {
		return x.CurrentBitrateKbps
	}
	return 0
}

func (x *ServerVideoChunk) GetPreviousBitrateKbps() int32 {
	if x != nil {
		return x.PreviousBitrateKbps
	}
	return 0
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
	mi := &file_config_protobuf_structs_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Interface) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Interface) ProtoMessage() {}

func (x *Interface) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Interface.ProtoReflect.Descriptor instead.
func (*Interface) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{4}
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
	Distance int64      `protobuf:"varint,2,opt,name=distance,proto3" json:"distance,omitempty"`
}

func (x *NextHop) Reset() {
	*x = NextHop{}
	mi := &file_config_protobuf_structs_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NextHop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextHop) ProtoMessage() {}

func (x *NextHop) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use NextHop.ProtoReflect.Descriptor instead.
func (*NextHop) Descriptor() ([]byte, []int) {
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{5}
}

func (x *NextHop) GetNextNode() *Interface {
	if x != nil {
		return x.NextNode
	}
	return nil
}

func (x *NextHop) GetDistance() int64 {
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
	mi := &file_config_protobuf_structs_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DistanceVectorRouting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DistanceVectorRouting) ProtoMessage() {}

func (x *DistanceVectorRouting) ProtoReflect() protoreflect.Message {
	mi := &file_config_protobuf_structs_proto_msgTypes[6]
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
	return file_config_protobuf_structs_proto_rawDescGZIP(), []int{6}
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
	0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x22, 0xd4, 0x04, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x12, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
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
	0x49, 0x0a, 0x12, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x65, 0x72, 0x5f, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x48, 0x00, 0x52, 0x11, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72,
	0x61, 0x70, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x56, 0x69,
	0x64, 0x65, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x65, 0x71, 0x5f, 0x69,
	0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x65, 0x71, 0x49, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x6f, 0x70, 0x73, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x68, 0x6f,
	0x70, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x31, 0x0a,
	0x11, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73,
	0x22, 0x75, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x14, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x12, 0x34, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49,
	0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x15, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xa5, 0x03, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x27, 0x0a, 0x0f,
	0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x2a, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x56, 0x69, 0x64, 0x65,
	0x6f, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x4c, 0x61,
	0x73, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x69, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x5f, 0x6b, 0x62, 0x70, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x62,
	0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x4b, 0x62, 0x70, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69,
	0x64, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68,
	0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x70, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x66, 0x70, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61,
	0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6d, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x4d, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6b, 0x62, 0x70,
	0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x12, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x42, 0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x4b, 0x62, 0x70, 0x73, 0x12, 0x32, 0x0a, 0x15, 0x70,
	0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x62, 0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x5f,
	0x6b, 0x62, 0x70, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x70, 0x72, 0x65, 0x76,
	0x69, 0x6f, 0x75, 0x73, 0x42, 0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x4b, 0x62, 0x70, 0x73, 0x22,
	0x2f, 0x0a, 0x09, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x22, 0x54, 0x0a, 0x07, 0x4e, 0x65, 0x78, 0x74, 0x48, 0x6f, 0x70, 0x12, 0x2d, 0x0a, 0x09, 0x6e,
	0x65, 0x78, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x52, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x69,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0xd2, 0x01, 0x0a, 0x15, 0x44, 0x69, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x43, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x29, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a,
	0x4a, 0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x48, 0x6f, 0x70,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x6c, 0x0a, 0x0b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x4f,
	0x55, 0x54, 0x49, 0x4e, 0x47, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a,
	0x52, 0x45, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x4d, 0x49, 0x54, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09,
	0x48, 0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x42,
	0x4f, 0x4f, 0x54, 0x53, 0x54, 0x52, 0x41, 0x50, 0x45, 0x52, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4c, 0x49, 0x45,
	0x4e, 0x54, 0x5f, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x05, 0x2a, 0x23, 0x0a, 0x0d, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4c,
	0x41, 0x59, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x01, 0x2a, 0x4a,
	0x0a, 0x0b, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x0b, 0x0a,
	0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x50,
	0x34, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x4b, 0x56, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03,
	0x41, 0x56, 0x49, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x45, 0x42, 0x4d, 0x10, 0x04, 0x12,
	0x09, 0x0a, 0x05, 0x4d, 0x4a, 0x50, 0x45, 0x47, 0x10, 0x05, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x64, 0x72, 0x69, 0x67, 0x6f,
	0x30, 0x33, 0x34, 0x35, 0x2f, 0x65, 0x73, 0x72, 0x2d, 0x74, 0x70, 0x32, 0x2f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
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
var file_config_protobuf_structs_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_config_protobuf_structs_proto_goTypes = []any{
	(RequestType)(0),              // 0: video.RequestType
	(PlayerCommand)(0),            // 1: video.PlayerCommand
	(VideoFormat)(0),              // 2: video.VideoFormat
	(*Header)(nil),                // 3: video.Header
	(*BootstraperResult)(nil),     // 4: video.BootstraperResult
	(*ClientCommand)(nil),         // 5: video.ClientCommand
	(*ServerVideoChunk)(nil),      // 6: video.ServerVideoChunk
	(*Interface)(nil),             // 7: video.Interface
	(*NextHop)(nil),               // 8: video.NextHop
	(*DistanceVectorRouting)(nil), // 9: video.DistanceVectorRouting
	nil,                           // 10: video.DistanceVectorRouting.EntriesEntry
}
var file_config_protobuf_structs_proto_depIdxs = []int32{
	0,  // 0: video.Header.type:type_name -> video.RequestType
	5,  // 1: video.Header.client_command:type_name -> video.ClientCommand
	6,  // 2: video.Header.server_video_chunk:type_name -> video.ServerVideoChunk
	9,  // 3: video.Header.distance_vector_routing:type_name -> video.DistanceVectorRouting
	4,  // 4: video.Header.bootstraper_result:type_name -> video.BootstraperResult
	1,  // 5: video.ClientCommand.command:type_name -> video.PlayerCommand
	2,  // 6: video.ServerVideoChunk.format:type_name -> video.VideoFormat
	7,  // 7: video.NextHop.next_node:type_name -> video.Interface
	10, // 8: video.DistanceVectorRouting.entries:type_name -> video.DistanceVectorRouting.EntriesEntry
	7,  // 9: video.DistanceVectorRouting.source:type_name -> video.Interface
	8,  // 10: video.DistanceVectorRouting.EntriesEntry.value:type_name -> video.NextHop
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
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
		(*Header_BootstraperResult)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_protobuf_structs_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
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
