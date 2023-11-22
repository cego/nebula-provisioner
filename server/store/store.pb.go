// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.1
// source: store.proto

package store

import (
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

type CA_Status int32

const (
	CA_Active   CA_Status = 0
	CA_Expired  CA_Status = 1
	CA_Inactive CA_Status = 2
	CA_Next     CA_Status = 3
)

// Enum value maps for CA_Status.
var (
	CA_Status_name = map[int32]string{
		0: "Active",
		1: "Expired",
		2: "Inactive",
		3: "Next",
	}
	CA_Status_value = map[string]int32{
		"Active":   0,
		"Expired":  1,
		"Inactive": 2,
		"Next":     3,
	}
)

func (x CA_Status) Enum() *CA_Status {
	p := new(CA_Status)
	*p = x
	return p
}

func (x CA_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CA_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_store_proto_enumTypes[0].Descriptor()
}

func (CA_Status) Type() protoreflect.EnumType {
	return &file_store_proto_enumTypes[0]
}

func (x CA_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CA_Status.Descriptor instead.
func (CA_Status) EnumDescriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{0, 0}
}

type CA struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetworkName string    `protobuf:"bytes,1,opt,name=networkName,proto3" json:"networkName,omitempty"`
	PublicKey   []byte    `protobuf:"bytes,10,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	PrivateKey  []byte    `protobuf:"bytes,11,opt,name=privateKey,proto3" json:"privateKey,omitempty"`
	Sha256Sum   string    `protobuf:"bytes,12,opt,name=sha256sum,proto3" json:"sha256sum,omitempty"`
	Status      CA_Status `protobuf:"varint,20,opt,name=status,proto3,enum=store.CA_Status" json:"status,omitempty"`
}

func (x *CA) Reset() {
	*x = CA{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CA) ProtoMessage() {}

func (x *CA) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CA.ProtoReflect.Descriptor instead.
func (*CA) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{0}
}

func (x *CA) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

func (x *CA) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *CA) GetPrivateKey() []byte {
	if x != nil {
		return x.PrivateKey
	}
	return nil
}

func (x *CA) GetSha256Sum() string {
	if x != nil {
		return x.Sha256Sum
	}
	return ""
}

func (x *CA) GetStatus() CA_Status {
	if x != nil {
		return x.Status
	}
	return CA_Active
}

type EnrollmentToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token       string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	NetworkName string `protobuf:"bytes,2,opt,name=networkName,proto3" json:"networkName,omitempty"`
}

func (x *EnrollmentToken) Reset() {
	*x = EnrollmentToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnrollmentToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnrollmentToken) ProtoMessage() {}

func (x *EnrollmentToken) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnrollmentToken.ProtoReflect.Descriptor instead.
func (*EnrollmentToken) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{1}
}

func (x *EnrollmentToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *EnrollmentToken) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

type EnrollmentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fingerprint []byte                 `protobuf:"bytes,1,opt,name=fingerprint,proto3" json:"fingerprint,omitempty"`
	Created     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created,proto3" json:"created,omitempty"`
	Token       string                 `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	NetworkName string                 `protobuf:"bytes,4,opt,name=networkName,proto3" json:"networkName,omitempty"`
	CsrPEM      string                 `protobuf:"bytes,5,opt,name=csrPEM,proto3" json:"csrPEM,omitempty"`
	ClientIP    string                 `protobuf:"bytes,6,opt,name=clientIP,proto3" json:"clientIP,omitempty"`
	Groups      []string               `protobuf:"bytes,7,rep,name=groups,proto3" json:"groups,omitempty"`
	Name        string                 `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	RequestedIP string                 `protobuf:"bytes,9,opt,name=requestedIP,proto3" json:"requestedIP,omitempty"`
}

func (x *EnrollmentRequest) Reset() {
	*x = EnrollmentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnrollmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnrollmentRequest) ProtoMessage() {}

func (x *EnrollmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnrollmentRequest.ProtoReflect.Descriptor instead.
func (*EnrollmentRequest) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{2}
}

func (x *EnrollmentRequest) GetFingerprint() []byte {
	if x != nil {
		return x.Fingerprint
	}
	return nil
}

func (x *EnrollmentRequest) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *EnrollmentRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *EnrollmentRequest) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

func (x *EnrollmentRequest) GetCsrPEM() string {
	if x != nil {
		return x.CsrPEM
	}
	return ""
}

func (x *EnrollmentRequest) GetClientIP() string {
	if x != nil {
		return x.ClientIP
	}
	return ""
}

func (x *EnrollmentRequest) GetGroups() []string {
	if x != nil {
		return x.Groups
	}
	return nil
}

func (x *EnrollmentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EnrollmentRequest) GetRequestedIP() string {
	if x != nil {
		return x.RequestedIP
	}
	return ""
}

type Agent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fingerprint   []byte                 `protobuf:"bytes,1,opt,name=fingerprint,proto3" json:"fingerprint,omitempty"`
	Created       *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created,proto3" json:"created,omitempty"`
	NetworkName   string                 `protobuf:"bytes,3,opt,name=networkName,proto3" json:"networkName,omitempty"`
	Groups        []string               `protobuf:"bytes,4,rep,name=groups,proto3" json:"groups,omitempty"`
	CsrPEM        string                 `protobuf:"bytes,5,opt,name=csrPEM,proto3" json:"csrPEM,omitempty"`
	AssignedIP    string                 `protobuf:"bytes,6,opt,name=assignedIP,proto3" json:"assignedIP,omitempty"`
	SignedPEM     string                 `protobuf:"bytes,10,opt,name=signedPEM,proto3" json:"signedPEM,omitempty"`
	IssuedAt      *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=issuedAt,proto3" json:"issuedAt,omitempty"`
	ExpiresAt     *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=expiresAt,proto3" json:"expiresAt,omitempty"`
	Name          string                 `protobuf:"bytes,13,opt,name=name,proto3" json:"name,omitempty"`
	OldSignedPEMs []string               `protobuf:"bytes,20,rep,name=oldSignedPEMs,proto3" json:"oldSignedPEMs,omitempty"`
}

func (x *Agent) Reset() {
	*x = Agent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Agent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Agent) ProtoMessage() {}

func (x *Agent) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Agent.ProtoReflect.Descriptor instead.
func (*Agent) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{3}
}

func (x *Agent) GetFingerprint() []byte {
	if x != nil {
		return x.Fingerprint
	}
	return nil
}

func (x *Agent) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Agent) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

func (x *Agent) GetGroups() []string {
	if x != nil {
		return x.Groups
	}
	return nil
}

func (x *Agent) GetCsrPEM() string {
	if x != nil {
		return x.CsrPEM
	}
	return ""
}

func (x *Agent) GetAssignedIP() string {
	if x != nil {
		return x.AssignedIP
	}
	return ""
}

func (x *Agent) GetSignedPEM() string {
	if x != nil {
		return x.SignedPEM
	}
	return ""
}

func (x *Agent) GetIssuedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.IssuedAt
	}
	return nil
}

func (x *Agent) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

func (x *Agent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Agent) GetOldSignedPEMs() []string {
	if x != nil {
		return x.OldSignedPEMs
	}
	return nil
}

type IPRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network []byte `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Netmask []byte `protobuf:"bytes,2,opt,name=netmask,proto3" json:"netmask,omitempty"`
}

func (x *IPRange) Reset() {
	*x = IPRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPRange) ProtoMessage() {}

func (x *IPRange) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPRange.ProtoReflect.Descriptor instead.
func (*IPRange) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{4}
}

func (x *IPRange) GetNetwork() []byte {
	if x != nil {
		return x.Network
	}
	return nil
}

func (x *IPRange) GetNetmask() []byte {
	if x != nil {
		return x.Netmask
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email    string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Created  *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created,proto3" json:"created,omitempty"`
	Approve  *UserApprove           `protobuf:"bytes,10,opt,name=approve,proto3" json:"approve,omitempty"`
	Disabled bool                   `protobuf:"varint,11,opt,name=disabled,proto3" json:"disabled,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{5}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *User) GetApprove() *UserApprove {
	if x != nil {
		return x.Approve
	}
	return nil
}

func (x *User) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

type UserApprove struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Approved   bool                   `protobuf:"varint,1,opt,name=approved,proto3" json:"approved,omitempty"`
	ApprovedBy string                 `protobuf:"bytes,2,opt,name=approvedBy,proto3" json:"approvedBy,omitempty"`
	ApprovedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=approvedAt,proto3" json:"approvedAt,omitempty"`
}

func (x *UserApprove) Reset() {
	*x = UserApprove{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserApprove) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserApprove) ProtoMessage() {}

func (x *UserApprove) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserApprove.ProtoReflect.Descriptor instead.
func (*UserApprove) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{6}
}

func (x *UserApprove) GetApproved() bool {
	if x != nil {
		return x.Approved
	}
	return false
}

func (x *UserApprove) GetApprovedBy() string {
	if x != nil {
		return x.ApprovedBy
	}
	return ""
}

func (x *UserApprove) GetApprovedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ApprovedAt
	}
	return nil
}

var File_store_proto protoreflect.FileDescriptor

var file_store_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe7, 0x01, 0x0a, 0x02, 0x43, 0x41, 0x12, 0x20, 0x0a, 0x0b,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x1e, 0x0a, 0x0a,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0a, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x73, 0x75, 0x6d, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x73, 0x75, 0x6d, 0x12, 0x28, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x43, 0x41, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x39, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0a,
	0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x64, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x6e, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x65, 0x78, 0x74, 0x10, 0x03, 0x22,
	0x49, 0x0a, 0x0f, 0x45, 0x6e, 0x72, 0x6f, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xa5, 0x02, 0x0a, 0x11, 0x45,
	0x6e, 0x72, 0x6f, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x20, 0x0a, 0x0b, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69,
	0x6e, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20,
	0x0a, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x63, 0x73, 0x72, 0x50, 0x45, 0x4d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x63, 0x73, 0x72, 0x50, 0x45, 0x4d, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x50, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x50, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x07,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x49, 0x50, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x49, 0x50, 0x22, 0x9b, 0x03, 0x0a, 0x05, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0b, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x34,
	0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x63, 0x73, 0x72, 0x50, 0x45, 0x4d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x73, 0x72, 0x50, 0x45, 0x4d, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x49, 0x50, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x49, 0x50, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x50, 0x45, 0x4d, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x50, 0x45, 0x4d, 0x12, 0x36, 0x0a, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38, 0x0a, 0x09,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x6f, 0x6c,
	0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x50, 0x45, 0x4d, 0x73, 0x18, 0x14, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0d, 0x6f, 0x6c, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x50, 0x45, 0x4d, 0x73,
	0x22, 0x3d, 0x0a, 0x07, 0x49, 0x50, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x73, 0x6b,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x73, 0x6b, 0x22,
	0xc0, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x07, 0x61, 0x70, 0x70, 0x72,
	0x6f, 0x76, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x52, 0x07, 0x61,
	0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c,
	0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c,
	0x65, 0x64, 0x22, 0x85, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x41, 0x70, 0x70, 0x72, 0x6f,
	0x76, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x42, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x42, 0x79, 0x12, 0x3a,
	0x0a, 0x0a, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a,
	0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x41, 0x74, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6c, 0x79, 0x6e, 0x67, 0x64, 0x6b,
	0x2f, 0x6e, 0x65, 0x62, 0x75, 0x6c, 0x61, 0x2d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_proto_rawDescOnce sync.Once
	file_store_proto_rawDescData = file_store_proto_rawDesc
)

func file_store_proto_rawDescGZIP() []byte {
	file_store_proto_rawDescOnce.Do(func() {
		file_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_proto_rawDescData)
	})
	return file_store_proto_rawDescData
}

var file_store_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_store_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_store_proto_goTypes = []interface{}{
	(CA_Status)(0),                // 0: store.CA.Status
	(*CA)(nil),                    // 1: store.CA
	(*EnrollmentToken)(nil),       // 2: store.EnrollmentToken
	(*EnrollmentRequest)(nil),     // 3: store.EnrollmentRequest
	(*Agent)(nil),                 // 4: store.Agent
	(*IPRange)(nil),               // 5: store.IPRange
	(*User)(nil),                  // 6: store.User
	(*UserApprove)(nil),           // 7: store.UserApprove
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_store_proto_depIdxs = []int32{
	0, // 0: store.CA.status:type_name -> store.CA.Status
	8, // 1: store.EnrollmentRequest.created:type_name -> google.protobuf.Timestamp
	8, // 2: store.Agent.created:type_name -> google.protobuf.Timestamp
	8, // 3: store.Agent.issuedAt:type_name -> google.protobuf.Timestamp
	8, // 4: store.Agent.expiresAt:type_name -> google.protobuf.Timestamp
	8, // 5: store.User.created:type_name -> google.protobuf.Timestamp
	7, // 6: store.User.approve:type_name -> store.UserApprove
	8, // 7: store.UserApprove.approvedAt:type_name -> google.protobuf.Timestamp
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_store_proto_init() }
func file_store_proto_init() {
	if File_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CA); i {
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
		file_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnrollmentToken); i {
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
		file_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnrollmentRequest); i {
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
		file_store_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Agent); i {
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
		file_store_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPRange); i {
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
		file_store_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_store_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserApprove); i {
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
			RawDescriptor: file_store_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_proto_goTypes,
		DependencyIndexes: file_store_proto_depIdxs,
		EnumInfos:         file_store_proto_enumTypes,
		MessageInfos:      file_store_proto_msgTypes,
	}.Build()
	File_store_proto = out.File
	file_store_proto_rawDesc = nil
	file_store_proto_goTypes = nil
	file_store_proto_depIdxs = nil
}
