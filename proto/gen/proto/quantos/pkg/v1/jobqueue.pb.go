// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: quantos/pkg/v1/jobqueue.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
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

// Type of Job State
type State_Type int32

const (
	State_TYPE_INVALID_UNSPECIFIED State_Type = 0
	// This step corresponds to the PRE step of the executor callback.
	State_TYPE_PRE State_Type = 1
	// This step corresponds to the POST step of the executor callback.
	State_TYPE_POST State_Type = 2
	// This step corresponds to the POST_KEEP_RUNNING step of the executor
	// callback.
	State_TYPE_POST_KEEP_RUNNING State_Type = 3
	// This step indicates to the master that the worker has successfully
	// completed the graph execution and is ready to perist the computation
	// results.
	State_TYPE_EXECUTED_GRAPH State_Type = 4
	// This step indicates to the master that the worker has successfully
	// persisted the computation results.
	State_TYPE_PERSISTED_RESULTS State_Type = 5
	// This step indicates to the master that the worker has completed the job.
	State_TYPE_COMPLETED_JOB State_Type = 6
)

// Enum value maps for State_Type.
var (
	State_Type_name = map[int32]string{
		0: "TYPE_INVALID_UNSPECIFIED",
		1: "TYPE_PRE",
		2: "TYPE_POST",
		3: "TYPE_POST_KEEP_RUNNING",
		4: "TYPE_EXECUTED_GRAPH",
		5: "TYPE_PERSISTED_RESULTS",
		6: "TYPE_COMPLETED_JOB",
	}
	State_Type_value = map[string]int32{
		"TYPE_INVALID_UNSPECIFIED": 0,
		"TYPE_PRE":                 1,
		"TYPE_POST":                2,
		"TYPE_POST_KEEP_RUNNING":   3,
		"TYPE_EXECUTED_GRAPH":      4,
		"TYPE_PERSISTED_RESULTS":   5,
		"TYPE_COMPLETED_JOB":       6,
	}
)

func (x State_Type) Enum() *State_Type {
	p := new(State_Type)
	*p = x
	return p
}

func (x State_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_quantos_pkg_v1_jobqueue_proto_enumTypes[0].Descriptor()
}

func (State_Type) Type() protoreflect.EnumType {
	return &file_quantos_pkg_v1_jobqueue_proto_enumTypes[0]
}

func (x State_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State_Type.Descriptor instead.
func (State_Type) EnumDescriptor() ([]byte, []int) {
	return file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP(), []int{1, 0}
}

type JobStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Payload:
	//	*JobStreamRequest_State
	//	*JobStreamRequest_RelayMessage
	Payload isJobStreamRequest_Payload `protobuf_oneof:"payload"`
}

func (x *JobStreamRequest) Reset() {
	*x = JobStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobStreamRequest) ProtoMessage() {}

func (x *JobStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobStreamRequest.ProtoReflect.Descriptor instead.
func (*JobStreamRequest) Descriptor() ([]byte, []int) {
	return file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP(), []int{0}
}

func (m *JobStreamRequest) GetPayload() isJobStreamRequest_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *JobStreamRequest) GetState() *State {
	if x, ok := x.GetPayload().(*JobStreamRequest_State); ok {
		return x.State
	}
	return nil
}

func (x *JobStreamRequest) GetRelayMessage() *RelayMessage {
	if x, ok := x.GetPayload().(*JobStreamRequest_RelayMessage); ok {
		return x.RelayMessage
	}
	return nil
}

type isJobStreamRequest_Payload interface {
	isJobStreamRequest_Payload()
}

type JobStreamRequest_State struct {
	State *State `protobuf:"bytes,1,opt,name=state,proto3,oneof"`
}

type JobStreamRequest_RelayMessage struct {
	RelayMessage *RelayMessage `protobuf:"bytes,2,opt,name=relay_message,json=relayMessage,proto3,oneof"`
}

func (*JobStreamRequest_State) isJobStreamRequest_Payload() {}

func (*JobStreamRequest_RelayMessage) isJobStreamRequest_Payload() {}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type State_Type `protobuf:"varint,1,opt,name=type,proto3,enum=proto.v1.State_Type" json:"type,omitempty"`
	// Workers use this field to submit their local aggregator delta values wen
	// reaching the POST step. The master collects the deltas, aggregates them to
	// its own aggregator values and broadcasts the global aggregator values in
	// the response. Workers must then *overwrite* their local aggregator values
	// with the values provided by the master.
	AggregatorValues map[string]*anypb.Any `protobuf:"bytes,2,rep,name=aggregator_values,json=aggregatorValues,proto3" json:"aggregator_values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Workers use this field to submit their local active-in-step count when
	// reaching the POST_KEEP_RUNNING step. The step response broadcasted by
	// the master uses the same field to specify the global active-in-step count
	// that the workers should pass to the graph executor callbacks.
	ActiveInState int64 `protobuf:"varint,3,opt,name=active_in_state,json=activeInState,proto3" json:"active_in_state,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP(), []int{1}
}

func (x *State) GetType() State_Type {
	if x != nil {
		return x.Type
	}
	return State_TYPE_INVALID_UNSPECIFIED
}

func (x *State) GetAggregatorValues() map[string]*anypb.Any {
	if x != nil {
		return x.AggregatorValues
	}
	return nil
}

func (x *State) GetActiveInState() int64 {
	if x != nil {
		return x.ActiveInState
	}
	return 0
}

type JobStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Payload:
	//	*JobStreamResponse_JobDetails
	//	*JobStreamResponse_State
	//	*JobStreamResponse_RelayMessage
	Payload isJobStreamResponse_Payload `protobuf_oneof:"payload"`
}

func (x *JobStreamResponse) Reset() {
	*x = JobStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobStreamResponse) ProtoMessage() {}

func (x *JobStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobStreamResponse.ProtoReflect.Descriptor instead.
func (*JobStreamResponse) Descriptor() ([]byte, []int) {
	return file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP(), []int{2}
}

func (m *JobStreamResponse) GetPayload() isJobStreamResponse_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *JobStreamResponse) GetJobDetails() *JobDetails {
	if x, ok := x.GetPayload().(*JobStreamResponse_JobDetails); ok {
		return x.JobDetails
	}
	return nil
}

func (x *JobStreamResponse) GetState() *State {
	if x, ok := x.GetPayload().(*JobStreamResponse_State); ok {
		return x.State
	}
	return nil
}

func (x *JobStreamResponse) GetRelayMessage() *RelayMessage {
	if x, ok := x.GetPayload().(*JobStreamResponse_RelayMessage); ok {
		return x.RelayMessage
	}
	return nil
}

type isJobStreamResponse_Payload interface {
	isJobStreamResponse_Payload()
}

type JobStreamResponse_JobDetails struct {
	JobDetails *JobDetails `protobuf:"bytes,1,opt,name=job_details,json=jobDetails,proto3,oneof"`
}

type JobStreamResponse_State struct {
	State *State `protobuf:"bytes,2,opt,name=state,proto3,oneof"`
}

type JobStreamResponse_RelayMessage struct {
	RelayMessage *RelayMessage `protobuf:"bytes,3,opt,name=relay_message,json=relayMessage,proto3,oneof"`
}

func (*JobStreamResponse_JobDetails) isJobStreamResponse_Payload() {}

func (*JobStreamResponse_State) isJobStreamResponse_Payload() {}

func (*JobStreamResponse_RelayMessage) isJobStreamResponse_Payload() {}

type RelayMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The message destination UUID.
	Destination string `protobuf:"bytes,1,opt,name=destination,proto3" json:"destination,omitempty"`
	// The serialized message contents.
	Message *anypb.Any `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RelayMessage) Reset() {
	*x = RelayMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelayMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelayMessage) ProtoMessage() {}

func (x *RelayMessage) ProtoReflect() protoreflect.Message {
	mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelayMessage.ProtoReflect.Descriptor instead.
func (*RelayMessage) Descriptor() ([]byte, []int) {
	return file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP(), []int{3}
}

func (x *RelayMessage) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *RelayMessage) GetMessage() *anypb.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

//JobDetails describes a job assigned by a master node to a worker.
type JobDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique ID for the job.
	JobId string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	// The creation time for the job.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The [from, to) UUID range assigned to the worker. Note that from is
	// inclusive and to is exclusive.
	PartitionFromUuid []byte `protobuf:"bytes,3,opt,name=partition_from_uuid,json=partitionFromUuid,proto3" json:"partition_from_uuid,omitempty"`
	PartitionToUuid   []byte `protobuf:"bytes,4,opt,name=partition_to_uuid,json=partitionToUuid,proto3" json:"partition_to_uuid,omitempty"`
}

func (x *JobDetails) Reset() {
	*x = JobDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobDetails) ProtoMessage() {}

func (x *JobDetails) ProtoReflect() protoreflect.Message {
	mi := &file_quantos_pkg_v1_jobqueue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobDetails.ProtoReflect.Descriptor instead.
func (*JobDetails) Descriptor() ([]byte, []int) {
	return file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP(), []int{4}
}

func (x *JobDetails) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

func (x *JobDetails) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *JobDetails) GetPartitionFromUuid() []byte {
	if x != nil {
		return x.PartitionFromUuid
	}
	return nil
}

func (x *JobDetails) GetPartitionToUuid() []byte {
	if x != nil {
		return x.PartitionToUuid
	}
	return nil
}

var File_quantos_pkg_v1_jobqueue_proto protoreflect.FileDescriptor

var file_quantos_pkg_v1_jobqueue_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x76, 0x31,
	0x2f, 0x6a, 0x6f, 0x62, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x01, 0x0a, 0x10, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x3d, 0x0a, 0x0d, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0xb5, 0x03,
	0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x52, 0x0a, 0x11, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x41, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x10, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f,
	0x69, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x49, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x59, 0x0a,
	0x15, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xaa, 0x01, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1c, 0x0a, 0x18, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49,
	0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x0c, 0x0a, 0x08, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x52, 0x45, 0x10, 0x01, 0x12, 0x0d, 0x0a,
	0x09, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x4f, 0x53, 0x54, 0x5f, 0x4b, 0x45, 0x45, 0x50, 0x5f, 0x52,
	0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x45, 0x44, 0x5f, 0x47, 0x52, 0x41, 0x50, 0x48, 0x10,
	0x04, 0x12, 0x1a, 0x0a, 0x16, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x45, 0x52, 0x53, 0x49, 0x53,
	0x54, 0x45, 0x44, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x10, 0x05, 0x12, 0x16, 0x0a,
	0x12, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x5f,
	0x4a, 0x4f, 0x42, 0x10, 0x06, 0x22, 0xbf, 0x01, 0x0a, 0x11, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x6a,
	0x6f, 0x62, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x48, 0x00, 0x52, 0x0a, 0x6a, 0x6f, 0x62, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x27, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3d, 0x0a,
	0x0d, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x6c, 0x61, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0c,
	0x72, 0x65, 0x6c, 0x61, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x09, 0x0a, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x60, 0x0a, 0x0c, 0x52, 0x65, 0x6c, 0x61, 0x79,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xba, 0x01, 0x0a, 0x0a, 0x4a, 0x6f,
	0x62, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x72, 0x6f, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x6f, 0x55, 0x75, 0x69, 0x64, 0x32, 0x5b, 0x0a, 0x0f, 0x4a, 0x6f, 0x62, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x09, 0x4a, 0x6f, 0x62,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f,
	0x62, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28,
	0x01, 0x30, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6f, 0x73, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x2f, 0x64, 0x65, 0x76, 0x2d, 0x30, 0x2e, 0x31, 0x2e, 0x30, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_quantos_pkg_v1_jobqueue_proto_rawDescOnce sync.Once
	file_quantos_pkg_v1_jobqueue_proto_rawDescData = file_quantos_pkg_v1_jobqueue_proto_rawDesc
)

func file_quantos_pkg_v1_jobqueue_proto_rawDescGZIP() []byte {
	file_quantos_pkg_v1_jobqueue_proto_rawDescOnce.Do(func() {
		file_quantos_pkg_v1_jobqueue_proto_rawDescData = protoimpl.X.CompressGZIP(file_quantos_pkg_v1_jobqueue_proto_rawDescData)
	})
	return file_quantos_pkg_v1_jobqueue_proto_rawDescData
}

var file_quantos_pkg_v1_jobqueue_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_quantos_pkg_v1_jobqueue_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_quantos_pkg_v1_jobqueue_proto_goTypes = []interface{}{
	(State_Type)(0),               // 0: proto.v1.State.Type
	(*JobStreamRequest)(nil),      // 1: proto.v1.JobStreamRequest
	(*State)(nil),                 // 2: proto.v1.State
	(*JobStreamResponse)(nil),     // 3: proto.v1.JobStreamResponse
	(*RelayMessage)(nil),          // 4: proto.v1.RelayMessage
	(*JobDetails)(nil),            // 5: proto.v1.JobDetails
	nil,                           // 6: proto.v1.State.AggregatorValuesEntry
	(*anypb.Any)(nil),             // 7: google.protobuf.Any
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_quantos_pkg_v1_jobqueue_proto_depIdxs = []int32{
	2,  // 0: proto.v1.JobStreamRequest.state:type_name -> proto.v1.State
	4,  // 1: proto.v1.JobStreamRequest.relay_message:type_name -> proto.v1.RelayMessage
	0,  // 2: proto.v1.State.type:type_name -> proto.v1.State.Type
	6,  // 3: proto.v1.State.aggregator_values:type_name -> proto.v1.State.AggregatorValuesEntry
	5,  // 4: proto.v1.JobStreamResponse.job_details:type_name -> proto.v1.JobDetails
	2,  // 5: proto.v1.JobStreamResponse.state:type_name -> proto.v1.State
	4,  // 6: proto.v1.JobStreamResponse.relay_message:type_name -> proto.v1.RelayMessage
	7,  // 7: proto.v1.RelayMessage.message:type_name -> google.protobuf.Any
	8,  // 8: proto.v1.JobDetails.created_at:type_name -> google.protobuf.Timestamp
	7,  // 9: proto.v1.State.AggregatorValuesEntry.value:type_name -> google.protobuf.Any
	1,  // 10: proto.v1.JobQueueService.JobStream:input_type -> proto.v1.JobStreamRequest
	3,  // 11: proto.v1.JobQueueService.JobStream:output_type -> proto.v1.JobStreamResponse
	11, // [11:12] is the sub-list for method output_type
	10, // [10:11] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_quantos_pkg_v1_jobqueue_proto_init() }
func file_quantos_pkg_v1_jobqueue_proto_init() {
	if File_quantos_pkg_v1_jobqueue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_quantos_pkg_v1_jobqueue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobStreamRequest); i {
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
		file_quantos_pkg_v1_jobqueue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
		file_quantos_pkg_v1_jobqueue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobStreamResponse); i {
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
		file_quantos_pkg_v1_jobqueue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelayMessage); i {
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
		file_quantos_pkg_v1_jobqueue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobDetails); i {
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
	file_quantos_pkg_v1_jobqueue_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*JobStreamRequest_State)(nil),
		(*JobStreamRequest_RelayMessage)(nil),
	}
	file_quantos_pkg_v1_jobqueue_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*JobStreamResponse_JobDetails)(nil),
		(*JobStreamResponse_State)(nil),
		(*JobStreamResponse_RelayMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_quantos_pkg_v1_jobqueue_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_quantos_pkg_v1_jobqueue_proto_goTypes,
		DependencyIndexes: file_quantos_pkg_v1_jobqueue_proto_depIdxs,
		EnumInfos:         file_quantos_pkg_v1_jobqueue_proto_enumTypes,
		MessageInfos:      file_quantos_pkg_v1_jobqueue_proto_msgTypes,
	}.Build()
	File_quantos_pkg_v1_jobqueue_proto = out.File
	file_quantos_pkg_v1_jobqueue_proto_rawDesc = nil
	file_quantos_pkg_v1_jobqueue_proto_goTypes = nil
	file_quantos_pkg_v1_jobqueue_proto_depIdxs = nil
}
