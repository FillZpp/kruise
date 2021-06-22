// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: apis/ctrlmesh/proto/ctrlmesh.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Status struct {
	SelfInfo             *SelfInfo `protobuf:"bytes,1,opt,name=selfInfo,proto3" json:"selfInfo,omitempty"`
	SpecHash             *SpecHash `protobuf:"bytes,2,opt,name=specHash,proto3" json:"specHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetSelfInfo() *SelfInfo {
	if m != nil {
		return m.SelfInfo
	}
	return nil
}

func (m *Status) GetSpecHash() *SpecHash {
	if m != nil {
		return m.SpecHash
	}
	return nil
}

type SelfInfo struct {
	Namespace            string   `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SelfInfo) Reset()         { *m = SelfInfo{} }
func (m *SelfInfo) String() string { return proto.CompactTextString(m) }
func (*SelfInfo) ProtoMessage()    {}
func (*SelfInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{1}
}
func (m *SelfInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfInfo.Unmarshal(m, b)
}
func (m *SelfInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfInfo.Marshal(b, m, deterministic)
}
func (m *SelfInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfInfo.Merge(m, src)
}
func (m *SelfInfo) XXX_Size() int {
	return xxx_messageInfo_SelfInfo.Size(m)
}
func (m *SelfInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SelfInfo proto.InternalMessageInfo

func (m *SelfInfo) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *SelfInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SpecHash struct {
	VopResourceVersion   string   `protobuf:"bytes,1,opt,name=vopResourceVersion,proto3" json:"vopResourceVersion,omitempty"`
	RouteHash            string   `protobuf:"bytes,2,opt,name=routeHash,proto3" json:"routeHash,omitempty"`
	RouteStrictHash      string   `protobuf:"bytes,3,opt,name=routeStrictHash,proto3" json:"routeStrictHash,omitempty"`
	EndpointsHash        string   `protobuf:"bytes,4,opt,name=endpointsHash,proto3" json:"endpointsHash,omitempty"`
	NamespacesHash       string   `protobuf:"bytes,5,opt,name=namespacesHash,proto3" json:"namespacesHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpecHash) Reset()         { *m = SpecHash{} }
func (m *SpecHash) String() string { return proto.CompactTextString(m) }
func (*SpecHash) ProtoMessage()    {}
func (*SpecHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{2}
}
func (m *SpecHash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpecHash.Unmarshal(m, b)
}
func (m *SpecHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpecHash.Marshal(b, m, deterministic)
}
func (m *SpecHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpecHash.Merge(m, src)
}
func (m *SpecHash) XXX_Size() int {
	return xxx_messageInfo_SpecHash.Size(m)
}
func (m *SpecHash) XXX_DiscardUnknown() {
	xxx_messageInfo_SpecHash.DiscardUnknown(m)
}

var xxx_messageInfo_SpecHash proto.InternalMessageInfo

func (m *SpecHash) GetVopResourceVersion() string {
	if m != nil {
		return m.VopResourceVersion
	}
	return ""
}

func (m *SpecHash) GetRouteHash() string {
	if m != nil {
		return m.RouteHash
	}
	return ""
}

func (m *SpecHash) GetRouteStrictHash() string {
	if m != nil {
		return m.RouteStrictHash
	}
	return ""
}

func (m *SpecHash) GetEndpointsHash() string {
	if m != nil {
		return m.EndpointsHash
	}
	return ""
}

func (m *SpecHash) GetNamespacesHash() string {
	if m != nil {
		return m.NamespacesHash
	}
	return ""
}

type Spec struct {
	Meta                 *SpecMeta   `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	Route                *Route      `protobuf:"bytes,2,opt,name=route,proto3" json:"route,omitempty"`
	Endpoints            []*Endpoint `protobuf:"bytes,3,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Spec) Reset()         { *m = Spec{} }
func (m *Spec) String() string { return proto.CompactTextString(m) }
func (*Spec) ProtoMessage()    {}
func (*Spec) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{3}
}
func (m *Spec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Spec.Unmarshal(m, b)
}
func (m *Spec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Spec.Marshal(b, m, deterministic)
}
func (m *Spec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Spec.Merge(m, src)
}
func (m *Spec) XXX_Size() int {
	return xxx_messageInfo_Spec.Size(m)
}
func (m *Spec) XXX_DiscardUnknown() {
	xxx_messageInfo_Spec.DiscardUnknown(m)
}

var xxx_messageInfo_Spec proto.InternalMessageInfo

func (m *Spec) GetMeta() *SpecMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Spec) GetRoute() *Route {
	if m != nil {
		return m.Route
	}
	return nil
}

func (m *Spec) GetEndpoints() []*Endpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

type SpecMeta struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	VopName              string   `protobuf:"bytes,2,opt,name=vopName,proto3" json:"vopName,omitempty"`
	VopResourceVersion   string   `protobuf:"bytes,3,opt,name=vopResourceVersion,proto3" json:"vopResourceVersion,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpecMeta) Reset()         { *m = SpecMeta{} }
func (m *SpecMeta) String() string { return proto.CompactTextString(m) }
func (*SpecMeta) ProtoMessage()    {}
func (*SpecMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{4}
}
func (m *SpecMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpecMeta.Unmarshal(m, b)
}
func (m *SpecMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpecMeta.Marshal(b, m, deterministic)
}
func (m *SpecMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpecMeta.Merge(m, src)
}
func (m *SpecMeta) XXX_Size() int {
	return xxx_messageInfo_SpecMeta.Size(m)
}
func (m *SpecMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_SpecMeta.DiscardUnknown(m)
}

var xxx_messageInfo_SpecMeta proto.InternalMessageInfo

func (m *SpecMeta) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SpecMeta) GetVopName() string {
	if m != nil {
		return m.VopName
	}
	return ""
}

func (m *SpecMeta) GetVopResourceVersion() string {
	if m != nil {
		return m.VopResourceVersion
	}
	return ""
}

type Route struct {
	Subset                  string                  `protobuf:"bytes,1,opt,name=subset,proto3" json:"subset,omitempty"`
	GlobalLimits            string                  `protobuf:"bytes,2,opt,name=globalLimits,proto3" json:"globalLimits,omitempty"`
	SubRules                string                  `protobuf:"bytes,3,opt,name=subRules,proto3" json:"subRules,omitempty"`
	Subsets                 string                  `protobuf:"bytes,4,opt,name=subsets,proto3" json:"subsets,omitempty"`
	GlobalExcludeNamespaces []string                `protobuf:"bytes,5,rep,name=globalExcludeNamespaces,proto3" json:"globalExcludeNamespaces,omitempty"`
	SensSubsetNamespaces    []*SensSubsetNamespaces `protobuf:"bytes,6,rep,name=sensSubsetNamespaces,proto3" json:"sensSubsetNamespaces,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}                `json:"-"`
	XXX_unrecognized        []byte                  `json:"-"`
	XXX_sizecache           int32                   `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{5}
}
func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}
func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}
func (m *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(m, src)
}
func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}
func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetSubset() string {
	if m != nil {
		return m.Subset
	}
	return ""
}

func (m *Route) GetGlobalLimits() string {
	if m != nil {
		return m.GlobalLimits
	}
	return ""
}

func (m *Route) GetSubRules() string {
	if m != nil {
		return m.SubRules
	}
	return ""
}

func (m *Route) GetSubsets() string {
	if m != nil {
		return m.Subsets
	}
	return ""
}

func (m *Route) GetGlobalExcludeNamespaces() []string {
	if m != nil {
		return m.GlobalExcludeNamespaces
	}
	return nil
}

func (m *Route) GetSensSubsetNamespaces() []*SensSubsetNamespaces {
	if m != nil {
		return m.SensSubsetNamespaces
	}
	return nil
}

type SensSubsetNamespaces struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespaces           []string `protobuf:"bytes,2,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SensSubsetNamespaces) Reset()         { *m = SensSubsetNamespaces{} }
func (m *SensSubsetNamespaces) String() string { return proto.CompactTextString(m) }
func (*SensSubsetNamespaces) ProtoMessage()    {}
func (*SensSubsetNamespaces) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{6}
}
func (m *SensSubsetNamespaces) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SensSubsetNamespaces.Unmarshal(m, b)
}
func (m *SensSubsetNamespaces) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SensSubsetNamespaces.Marshal(b, m, deterministic)
}
func (m *SensSubsetNamespaces) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SensSubsetNamespaces.Merge(m, src)
}
func (m *SensSubsetNamespaces) XXX_Size() int {
	return xxx_messageInfo_SensSubsetNamespaces.Size(m)
}
func (m *SensSubsetNamespaces) XXX_DiscardUnknown() {
	xxx_messageInfo_SensSubsetNamespaces.DiscardUnknown(m)
}

var xxx_messageInfo_SensSubsetNamespaces proto.InternalMessageInfo

func (m *SensSubsetNamespaces) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SensSubsetNamespaces) GetNamespaces() []string {
	if m != nil {
		return m.Namespaces
	}
	return nil
}

type Endpoint struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Subset               string   `protobuf:"bytes,2,opt,name=subset,proto3" json:"subset,omitempty"`
	Ip                   string   `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_015120ea8d6d1ea6, []int{7}
}
func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (m *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(m, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Endpoint) GetSubset() string {
	if m != nil {
		return m.Subset
	}
	return ""
}

func (m *Endpoint) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func init() {
	proto.RegisterType((*Status)(nil), "proto.Status")
	proto.RegisterType((*SelfInfo)(nil), "proto.SelfInfo")
	proto.RegisterType((*SpecHash)(nil), "proto.SpecHash")
	proto.RegisterType((*Spec)(nil), "proto.Spec")
	proto.RegisterType((*SpecMeta)(nil), "proto.SpecMeta")
	proto.RegisterType((*Route)(nil), "proto.Route")
	proto.RegisterType((*SensSubsetNamespaces)(nil), "proto.SensSubsetNamespaces")
	proto.RegisterType((*Endpoint)(nil), "proto.Endpoint")
}

func init() {
	proto.RegisterFile("apis/ctrlmesh/proto/ctrlmesh.proto", fileDescriptor_015120ea8d6d1ea6)
}

var fileDescriptor_015120ea8d6d1ea6 = []byte{
	// 518 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0x41, 0x6f, 0xd3, 0x4c,
	0x10, 0xfd, 0xec, 0xc4, 0xa9, 0x3d, 0x6d, 0x53, 0x69, 0x14, 0x7d, 0x58, 0x05, 0xa1, 0x68, 0x41,
	0x28, 0x12, 0x90, 0xa2, 0x72, 0xe1, 0x80, 0x38, 0x80, 0x8a, 0x00, 0xd1, 0x22, 0x6d, 0x24, 0x0e,
	0xdc, 0x6c, 0x77, 0xda, 0x18, 0x1c, 0xaf, 0xe5, 0x5d, 0x57, 0x9c, 0x38, 0xf1, 0x03, 0xf9, 0x49,
	0xc8, 0xe3, 0xb5, 0x9d, 0x84, 0xe4, 0xe4, 0xbc, 0xb7, 0x6f, 0xdf, 0xcc, 0xbe, 0xd9, 0x2c, 0x88,
	0xa8, 0x48, 0xf5, 0x59, 0x62, 0xca, 0x6c, 0x45, 0x7a, 0x79, 0x56, 0x94, 0xca, 0xa8, 0x0e, 0xce,
	0x19, 0xa2, 0xc7, 0x1f, 0x11, 0xc3, 0x68, 0x61, 0x22, 0x53, 0x69, 0x7c, 0x0a, 0xbe, 0xa6, 0xec,
	0xe6, 0x63, 0x7e, 0xa3, 0x42, 0x67, 0xea, 0xcc, 0x0e, 0xcf, 0x4f, 0x1a, 0xe9, 0x7c, 0x61, 0x69,
	0xd9, 0x09, 0x58, 0x5c, 0x50, 0xf2, 0x21, 0xd2, 0xcb, 0xd0, 0xdd, 0x14, 0x5b, 0x5a, 0x76, 0x02,
	0xf1, 0x1a, 0xfc, 0xd6, 0x02, 0x1f, 0x40, 0x90, 0x47, 0x2b, 0xd2, 0x45, 0x94, 0x10, 0x97, 0x09,
	0x64, 0x4f, 0x20, 0xc2, 0xb0, 0x06, 0x6c, 0x19, 0x48, 0xfe, 0x2d, 0xfe, 0x38, 0xe0, 0xb7, 0xa6,
	0x38, 0x07, 0xbc, 0x53, 0x85, 0x24, 0xad, 0xaa, 0x32, 0xa1, 0xaf, 0x54, 0xea, 0x54, 0xe5, 0xd6,
	0x67, 0xc7, 0x4a, 0x5d, 0xae, 0x54, 0x95, 0xa1, 0xae, 0xd1, 0x40, 0xf6, 0x04, 0xce, 0xe0, 0x84,
	0xc1, 0xc2, 0x94, 0x69, 0x62, 0x58, 0x33, 0x60, 0xcd, 0x36, 0x8d, 0x8f, 0xe1, 0x98, 0xf2, 0xeb,
	0x42, 0xa5, 0xb9, 0xd1, 0xac, 0x1b, 0xb2, 0x6e, 0x93, 0xc4, 0x27, 0x30, 0xee, 0xce, 0xd2, 0xc8,
	0x3c, 0x96, 0x6d, 0xb1, 0xe2, 0x17, 0x0c, 0xeb, 0x13, 0xe1, 0x23, 0x18, 0xae, 0xc8, 0x44, 0xdb,
	0x71, 0x17, 0x94, 0x5c, 0x92, 0x89, 0x24, 0x2f, 0xa2, 0x00, 0x8f, 0xbb, 0xb1, 0x39, 0x1f, 0x59,
	0x95, 0xac, 0x39, 0xd9, 0x2c, 0xe1, 0x73, 0x08, 0xba, 0x4e, 0xc2, 0xc1, 0x74, 0xb0, 0xe6, 0x76,
	0x61, 0x79, 0xd9, 0x2b, 0xc4, 0xf7, 0x26, 0xd1, 0xba, 0x08, 0x4e, 0xc0, 0x33, 0xea, 0x07, 0xb5,
	0x21, 0x36, 0x00, 0x43, 0x38, 0xb8, 0x53, 0xc5, 0x55, 0x3f, 0x8b, 0x16, 0xee, 0x99, 0xc0, 0x60,
	0xdf, 0x04, 0xc4, 0x6f, 0x17, 0x3c, 0xee, 0x15, 0xff, 0x87, 0x91, 0xae, 0x62, 0x4d, 0xc6, 0x96,
	0xb2, 0x08, 0x05, 0x1c, 0xdd, 0x66, 0x2a, 0x8e, 0xb2, 0xcf, 0xe9, 0x2a, 0x35, 0xda, 0x16, 0xdc,
	0xe0, 0xf0, 0x14, 0x7c, 0x5d, 0xc5, 0xb2, 0xca, 0x48, 0xdb, 0x5a, 0x1d, 0xae, 0x7b, 0x6d, 0x9c,
	0xb4, 0x9d, 0x4a, 0x0b, 0xf1, 0x15, 0xdc, 0x6b, 0x5c, 0x2e, 0x7e, 0x26, 0x59, 0x75, 0x4d, 0x57,
	0xdd, 0x18, 0x42, 0x6f, 0x3a, 0x98, 0x05, 0x72, 0xdf, 0x32, 0x7e, 0x81, 0x89, 0xa6, 0x5c, 0x2f,
	0xd8, 0x68, 0x6d, 0xdb, 0x88, 0xb3, 0xbd, 0xdf, 0xfd, 0x31, 0xfe, 0x95, 0xc8, 0x9d, 0x1b, 0xc5,
	0x27, 0x98, 0xec, 0x52, 0x77, 0x37, 0xde, 0xe9, 0x6f, 0x3c, 0x3e, 0x04, 0xe8, 0x2f, 0x4c, 0xe8,
	0x72, 0xa7, 0x6b, 0x8c, 0x78, 0x0f, 0x7e, 0x3b, 0xd5, 0x9d, 0xfb, 0xfb, 0xa0, 0xdd, 0x8d, 0xa0,
	0xc7, 0xe0, 0xa6, 0x85, 0x8d, 0xcf, 0x4d, 0x8b, 0xf3, 0x37, 0x30, 0x7e, 0xa7, 0x72, 0x53, 0xaa,
	0x2c, 0xa3, 0xf2, 0x92, 0xf4, 0x12, 0x9f, 0x81, 0x2f, 0xe9, 0x36, 0xd5, 0x86, 0x4a, 0x3c, 0x6e,
	0x0f, 0xc9, 0xcf, 0xc3, 0xe9, 0xe1, 0xda, 0xed, 0x14, 0xff, 0xcd, 0x9c, 0x17, 0xce, 0xdb, 0x83,
	0x6f, 0xcd, 0x23, 0x12, 0x8f, 0xf8, 0xf3, 0xf2, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6a, 0xcf,
	0xca, 0xee, 0x78, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ControllerMeshClient is the client API for ControllerMesh service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ControllerMeshClient interface {
	Register(ctx context.Context, opts ...grpc.CallOption) (ControllerMesh_RegisterClient, error)
}

type controllerMeshClient struct {
	cc *grpc.ClientConn
}

func NewControllerMeshClient(cc *grpc.ClientConn) ControllerMeshClient {
	return &controllerMeshClient{cc}
}

func (c *controllerMeshClient) Register(ctx context.Context, opts ...grpc.CallOption) (ControllerMesh_RegisterClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ControllerMesh_serviceDesc.Streams[0], "/proto.ControllerMesh/Register", opts...)
	if err != nil {
		return nil, err
	}
	x := &controllerMeshRegisterClient{stream}
	return x, nil
}

type ControllerMesh_RegisterClient interface {
	Send(*Status) error
	Recv() (*Spec, error)
	grpc.ClientStream
}

type controllerMeshRegisterClient struct {
	grpc.ClientStream
}

func (x *controllerMeshRegisterClient) Send(m *Status) error {
	return x.ClientStream.SendMsg(m)
}

func (x *controllerMeshRegisterClient) Recv() (*Spec, error) {
	m := new(Spec)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ControllerMeshServer is the server API for ControllerMesh service.
type ControllerMeshServer interface {
	Register(ControllerMesh_RegisterServer) error
}

// UnimplementedControllerMeshServer can be embedded to have forward compatible implementations.
type UnimplementedControllerMeshServer struct {
}

func (*UnimplementedControllerMeshServer) Register(srv ControllerMesh_RegisterServer) error {
	return status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func RegisterControllerMeshServer(s *grpc.Server, srv ControllerMeshServer) {
	s.RegisterService(&_ControllerMesh_serviceDesc, srv)
}

func _ControllerMesh_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ControllerMeshServer).Register(&controllerMeshRegisterServer{stream})
}

type ControllerMesh_RegisterServer interface {
	Send(*Spec) error
	Recv() (*Status, error)
	grpc.ServerStream
}

type controllerMeshRegisterServer struct {
	grpc.ServerStream
}

func (x *controllerMeshRegisterServer) Send(m *Spec) error {
	return x.ServerStream.SendMsg(m)
}

func (x *controllerMeshRegisterServer) Recv() (*Status, error) {
	m := new(Status)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ControllerMesh_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ControllerMesh",
	HandlerType: (*ControllerMeshServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _ControllerMesh_Register_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/ctrlmesh/proto/ctrlmesh.proto",
}
