// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: client_service.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// The request message containing the command.包含命令的请求消息
type CommandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//string command = 1;
	LoginAccount     string `protobuf:"bytes,1,opt,name=login_account,json=loginAccount,proto3" json:"login_account,omitempty"`
	LoginPassword    string `protobuf:"bytes,2,opt,name=login_password,json=loginPassword,proto3" json:"login_password,omitempty"`
	RegisterAccount  string `protobuf:"bytes,3,opt,name=register_account,json=registerAccount,proto3" json:"register_account,omitempty"`
	RegisterPassword string `protobuf:"bytes,4,opt,name=register_password,json=registerPassword,proto3" json:"register_password,omitempty"`
	ClientIp         string `protobuf:"bytes,5,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
	ServerPassword   string `protobuf:"bytes,6,opt,name=server_password,json=serverPassword,proto3" json:"server_password,omitempty"`
	SearchContent    string `protobuf:"bytes,7,opt,name=search_content,json=searchContent,proto3" json:"search_content,omitempty"`
	LogoutIng        string `protobuf:"bytes,8,opt,name=logout_ing,json=logoutIng,proto3" json:"logout_ing,omitempty"`
	TabelName        string `protobuf:"bytes,9,opt,name=tabel_name,json=tabelName,proto3" json:"tabel_name,omitempty"`
	PermissionTable  string `protobuf:"bytes,10,opt,name=permissionTable,proto3" json:"permissionTable,omitempty"`
	Key              string `protobuf:"bytes,11,opt,name=key,proto3" json:"key,omitempty"`
	Value            string `protobuf:"bytes,12,opt,name=value,proto3" json:"value,omitempty"`
	ViewMyinfo       string `protobuf:"bytes,13,opt,name=view_myinfo,json=viewMyinfo,proto3" json:"view_myinfo,omitempty"`
	Uid              string `protobuf:"bytes,14,opt,name=uid,proto3" json:"uid,omitempty"`
	Creat            string `protobuf:"bytes,15,opt,name=creat,proto3" json:"creat,omitempty"`
	CreatPassword    string `protobuf:"bytes,16,opt,name=creat_password,json=creatPassword,proto3" json:"creat_password,omitempty"`
	Join             string `protobuf:"bytes,17,opt,name=join,proto3" json:"join,omitempty"`
	JoinIp           string `protobuf:"bytes,18,opt,name=join_ip,json=joinIp,proto3" json:"join_ip,omitempty"`
	JoinPassword     string `protobuf:"bytes,19,opt,name=join_password,json=joinPassword,proto3" json:"join_password,omitempty"`
	Set              string `protobuf:"bytes,20,opt,name=set,proto3" json:"set,omitempty"`
	PkgNum           string `protobuf:"bytes,21,opt,name=pkgNum,proto3" json:"pkgNum,omitempty"`
}

func (x *CommandRequest) Reset() {
	*x = CommandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandRequest) ProtoMessage() {}

func (x *CommandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandRequest.ProtoReflect.Descriptor instead.
func (*CommandRequest) Descriptor() ([]byte, []int) {
	return file_client_service_proto_rawDescGZIP(), []int{0}
}

func (x *CommandRequest) GetLoginAccount() string {
	if x != nil {
		return x.LoginAccount
	}
	return ""
}

func (x *CommandRequest) GetLoginPassword() string {
	if x != nil {
		return x.LoginPassword
	}
	return ""
}

func (x *CommandRequest) GetRegisterAccount() string {
	if x != nil {
		return x.RegisterAccount
	}
	return ""
}

func (x *CommandRequest) GetRegisterPassword() string {
	if x != nil {
		return x.RegisterPassword
	}
	return ""
}

func (x *CommandRequest) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

func (x *CommandRequest) GetServerPassword() string {
	if x != nil {
		return x.ServerPassword
	}
	return ""
}

func (x *CommandRequest) GetSearchContent() string {
	if x != nil {
		return x.SearchContent
	}
	return ""
}

func (x *CommandRequest) GetLogoutIng() string {
	if x != nil {
		return x.LogoutIng
	}
	return ""
}

func (x *CommandRequest) GetTabelName() string {
	if x != nil {
		return x.TabelName
	}
	return ""
}

func (x *CommandRequest) GetPermissionTable() string {
	if x != nil {
		return x.PermissionTable
	}
	return ""
}

func (x *CommandRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CommandRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *CommandRequest) GetViewMyinfo() string {
	if x != nil {
		return x.ViewMyinfo
	}
	return ""
}

func (x *CommandRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *CommandRequest) GetCreat() string {
	if x != nil {
		return x.Creat
	}
	return ""
}

func (x *CommandRequest) GetCreatPassword() string {
	if x != nil {
		return x.CreatPassword
	}
	return ""
}

func (x *CommandRequest) GetJoin() string {
	if x != nil {
		return x.Join
	}
	return ""
}

func (x *CommandRequest) GetJoinIp() string {
	if x != nil {
		return x.JoinIp
	}
	return ""
}

func (x *CommandRequest) GetJoinPassword() string {
	if x != nil {
		return x.JoinPassword
	}
	return ""
}

func (x *CommandRequest) GetSet() string {
	if x != nil {
		return x.Set
	}
	return ""
}

func (x *CommandRequest) GetPkgNum() string {
	if x != nil {
		return x.PkgNum
	}
	return ""
}

// The response message containing the execute results. 包含执行命令结果的响应消息
type CommandReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//string executeresult = 1;
	Login         string   `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Register      string   `protobuf:"bytes,2,opt,name=register,proto3" json:"register,omitempty"`
	ServerConnect string   `protobuf:"bytes,3,opt,name=server_connect,json=serverConnect,proto3" json:"server_connect,omitempty"`
	SearchResult  string   `protobuf:"bytes,4,opt,name=search_result,json=searchResult,proto3" json:"search_result,omitempty"`
	TabelNames    []string `protobuf:"bytes,5,rep,name=tabel_names,json=tabelNames,proto3" json:"tabel_names,omitempty"` //字符串数组
	Uaddress      string   `protobuf:"bytes,6,opt,name=uaddress,proto3" json:"uaddress,omitempty"`
	Value         string   `protobuf:"bytes,7,opt,name=value,proto3" json:"value,omitempty"`
	TabelOk       string   `protobuf:"bytes,8,opt,name=tabel_ok,json=tabelOk,proto3" json:"tabel_ok,omitempty"`
	PutOk         string   `protobuf:"bytes,9,opt,name=put_ok,json=putOk,proto3" json:"put_ok,omitempty"`
	GetOk         string   `protobuf:"bytes,10,opt,name=get_ok,json=getOk,proto3" json:"get_ok,omitempty"`
	LogoutOk      string   `protobuf:"bytes,11,opt,name=Logout_ok,json=LogoutOk,proto3" json:"Logout_ok,omitempty"`
	CreatOk       string   `protobuf:"bytes,12,opt,name=creat_ok,json=creatOk,proto3" json:"creat_ok,omitempty"`
	JoinOk        string   `protobuf:"bytes,13,opt,name=join_ok,json=joinOk,proto3" json:"join_ok,omitempty"`
}

func (x *CommandReply) Reset() {
	*x = CommandReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandReply) ProtoMessage() {}

func (x *CommandReply) ProtoReflect() protoreflect.Message {
	mi := &file_client_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandReply.ProtoReflect.Descriptor instead.
func (*CommandReply) Descriptor() ([]byte, []int) {
	return file_client_service_proto_rawDescGZIP(), []int{1}
}

func (x *CommandReply) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *CommandReply) GetRegister() string {
	if x != nil {
		return x.Register
	}
	return ""
}

func (x *CommandReply) GetServerConnect() string {
	if x != nil {
		return x.ServerConnect
	}
	return ""
}

func (x *CommandReply) GetSearchResult() string {
	if x != nil {
		return x.SearchResult
	}
	return ""
}

func (x *CommandReply) GetTabelNames() []string {
	if x != nil {
		return x.TabelNames
	}
	return nil
}

func (x *CommandReply) GetUaddress() string {
	if x != nil {
		return x.Uaddress
	}
	return ""
}

func (x *CommandReply) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *CommandReply) GetTabelOk() string {
	if x != nil {
		return x.TabelOk
	}
	return ""
}

func (x *CommandReply) GetPutOk() string {
	if x != nil {
		return x.PutOk
	}
	return ""
}

func (x *CommandReply) GetGetOk() string {
	if x != nil {
		return x.GetOk
	}
	return ""
}

func (x *CommandReply) GetLogoutOk() string {
	if x != nil {
		return x.LogoutOk
	}
	return ""
}

func (x *CommandReply) GetCreatOk() string {
	if x != nil {
		return x.CreatOk
	}
	return ""
}

func (x *CommandReply) GetJoinOk() string {
	if x != nil {
		return x.JoinOk
	}
	return ""
}

//流数据请求
type StreamReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *StreamReq) Reset() {
	*x = StreamReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamReq) ProtoMessage() {}

func (x *StreamReq) ProtoReflect() protoreflect.Message {
	mi := &file_client_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamReq.ProtoReflect.Descriptor instead.
func (*StreamReq) Descriptor() ([]byte, []int) {
	return file_client_service_proto_rawDescGZIP(), []int{2}
}

func (x *StreamReq) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

//流数据响应
type StreamRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *StreamRes) Reset() {
	*x = StreamRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRes) ProtoMessage() {}

func (x *StreamRes) ProtoReflect() protoreflect.Message {
	mi := &file_client_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRes.ProtoReflect.Descriptor instead.
func (*StreamRes) Descriptor() ([]byte, []int) {
	return file_client_service_proto_rawDescGZIP(), []int{3}
}

func (x *StreamRes) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_client_service_proto protoreflect.FileDescriptor

var file_client_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22, 0x9d, 0x05, 0x0a,
	0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x23, 0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x11, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x10, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x70,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70,
	0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x5f, 0x69, 0x6e, 0x67, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x49, 0x6e, 0x67, 0x12,
	0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x28,
	0x0a, 0x0f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x69, 0x65, 0x77, 0x5f, 0x6d, 0x79, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x76, 0x69, 0x65, 0x77, 0x4d, 0x79, 0x69, 0x6e, 0x66,
	0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x72, 0x65, 0x61, 0x74, 0x18, 0x0f, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x63, 0x72, 0x65, 0x61, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6a, 0x6f, 0x69, 0x6e, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6a, 0x6f, 0x69, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x69, 0x70, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6a, 0x6f, 0x69, 0x6e, 0x49, 0x70, 0x12, 0x23, 0x0a,
	0x0d, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x13,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6a, 0x6f, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x73, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6b, 0x67, 0x4e, 0x75, 0x6d, 0x18, 0x15,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6b, 0x67, 0x4e, 0x75, 0x6d, 0x22, 0xf9, 0x02, 0x0a,
	0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74,
	0x61, 0x62, 0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0a, 0x74, 0x61, 0x62, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x74, 0x61, 0x62, 0x65, 0x6c, 0x5f, 0x6f, 0x6b, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x74, 0x61, 0x62, 0x65, 0x6c, 0x4f, 0x6b, 0x12, 0x15, 0x0a, 0x06, 0x70, 0x75, 0x74,
	0x5f, 0x6f, 0x6b, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x75, 0x74, 0x4f, 0x6b,
	0x12, 0x15, 0x0a, 0x06, 0x67, 0x65, 0x74, 0x5f, 0x6f, 0x6b, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x67, 0x65, 0x74, 0x4f, 0x6b, 0x12, 0x1b, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x6f, 0x75,
	0x74, 0x5f, 0x6f, 0x6b, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4c, 0x6f, 0x67, 0x6f,
	0x75, 0x74, 0x4f, 0x6b, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x5f, 0x6f, 0x6b,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x4f, 0x6b, 0x12,
	0x17, 0x0a, 0x07, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x6f, 0x6b, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6a, 0x6f, 0x69, 0x6e, 0x4f, 0x6b, 0x22, 0x1f, 0x0a, 0x09, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x1f, 0x0a, 0x09, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xdc, 0x01, 0x0a, 0x06, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x14, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x0f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x22, 0x00, 0x30, 0x01, 0x12, 0x34,
	0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x0f,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x1a,
	0x0f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x22, 0x00, 0x28, 0x01, 0x12, 0x33, 0x0a, 0x09, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x54, 0x77,
	0x6f, 0x12, 0x0f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2e, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_service_proto_rawDescOnce sync.Once
	file_client_service_proto_rawDescData = file_client_service_proto_rawDesc
)

func file_client_service_proto_rawDescGZIP() []byte {
	file_client_service_proto_rawDescOnce.Do(func() {
		file_client_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_service_proto_rawDescData)
	})
	return file_client_service_proto_rawDescData
}

var file_client_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_client_service_proto_goTypes = []interface{}{
	(*CommandRequest)(nil), // 0: grpc.CommandRequest
	(*CommandReply)(nil),   // 1: grpc.CommandReply
	(*StreamReq)(nil),      // 2: grpc.StreamReq
	(*StreamRes)(nil),      // 3: grpc.StreamRes
}
var file_client_service_proto_depIdxs = []int32{
	0, // 0: grpc.Server.cmd:input_type -> grpc.CommandRequest
	2, // 1: grpc.Server.StreamServer:input_type -> grpc.StreamReq
	2, // 2: grpc.Server.StreamClient:input_type -> grpc.StreamReq
	2, // 3: grpc.Server.StreamTwo:input_type -> grpc.StreamReq
	1, // 4: grpc.Server.cmd:output_type -> grpc.CommandReply
	3, // 5: grpc.Server.StreamServer:output_type -> grpc.StreamRes
	3, // 6: grpc.Server.StreamClient:output_type -> grpc.StreamRes
	3, // 7: grpc.Server.StreamTwo:output_type -> grpc.StreamRes
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_client_service_proto_init() }
func file_client_service_proto_init() {
	if File_client_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_client_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandRequest); i {
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
		file_client_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandReply); i {
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
		file_client_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamReq); i {
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
		file_client_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamRes); i {
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
			RawDescriptor: file_client_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_client_service_proto_goTypes,
		DependencyIndexes: file_client_service_proto_depIdxs,
		MessageInfos:      file_client_service_proto_msgTypes,
	}.Build()
	File_client_service_proto = out.File
	file_client_service_proto_rawDesc = nil
	file_client_service_proto_goTypes = nil
	file_client_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerClient interface {
	//Sends a command information 发送命令信息
	Cmd(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandReply, error)
	//服务器流
	StreamServer(ctx context.Context, in *StreamReq, opts ...grpc.CallOption) (Server_StreamServerClient, error)
	//客户端流
	StreamClient(ctx context.Context, opts ...grpc.CallOption) (Server_StreamClientClient, error)
	//双向流
	StreamTwo(ctx context.Context, opts ...grpc.CallOption) (Server_StreamTwoClient, error)
}

type serverClient struct {
	cc grpc.ClientConnInterface
}

func NewServerClient(cc grpc.ClientConnInterface) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) Cmd(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandReply, error) {
	out := new(CommandReply)
	err := c.cc.Invoke(ctx, "/grpc.Server/cmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) StreamServer(ctx context.Context, in *StreamReq, opts ...grpc.CallOption) (Server_StreamServerClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Server_serviceDesc.Streams[0], "/grpc.Server/StreamServer", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverStreamServerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Server_StreamServerClient interface {
	Recv() (*StreamRes, error)
	grpc.ClientStream
}

type serverStreamServerClient struct {
	grpc.ClientStream
}

func (x *serverStreamServerClient) Recv() (*StreamRes, error) {
	m := new(StreamRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serverClient) StreamClient(ctx context.Context, opts ...grpc.CallOption) (Server_StreamClientClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Server_serviceDesc.Streams[1], "/grpc.Server/StreamClient", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverStreamClientClient{stream}
	return x, nil
}

type Server_StreamClientClient interface {
	Send(*StreamReq) error
	CloseAndRecv() (*StreamRes, error)
	grpc.ClientStream
}

type serverStreamClientClient struct {
	grpc.ClientStream
}

func (x *serverStreamClientClient) Send(m *StreamReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serverStreamClientClient) CloseAndRecv() (*StreamRes, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StreamRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serverClient) StreamTwo(ctx context.Context, opts ...grpc.CallOption) (Server_StreamTwoClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Server_serviceDesc.Streams[2], "/grpc.Server/StreamTwo", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverStreamTwoClient{stream}
	return x, nil
}

type Server_StreamTwoClient interface {
	Send(*StreamReq) error
	Recv() (*StreamRes, error)
	grpc.ClientStream
}

type serverStreamTwoClient struct {
	grpc.ClientStream
}

func (x *serverStreamTwoClient) Send(m *StreamReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serverStreamTwoClient) Recv() (*StreamRes, error) {
	m := new(StreamRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServerServer is the server API for Server service.
type ServerServer interface {
	//Sends a command information 发送命令信息
	Cmd(context.Context, *CommandRequest) (*CommandReply, error)
	//服务器流
	StreamServer(*StreamReq, Server_StreamServerServer) error
	//客户端流
	StreamClient(Server_StreamClientServer) error
	//双向流
	StreamTwo(Server_StreamTwoServer) error
}

// UnimplementedServerServer can be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (*UnimplementedServerServer) Cmd(context.Context, *CommandRequest) (*CommandReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cmd not implemented")
}
func (*UnimplementedServerServer) StreamServer(*StreamReq, Server_StreamServerServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamServer not implemented")
}
func (*UnimplementedServerServer) StreamClient(Server_StreamClientServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamClient not implemented")
}
func (*UnimplementedServerServer) StreamTwo(Server_StreamTwoServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamTwo not implemented")
}

func RegisterServerServer(s *grpc.Server, srv ServerServer) {
	s.RegisterService(&_Server_serviceDesc, srv)
}

func _Server_Cmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).Cmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Server/Cmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).Cmd(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_StreamServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServerServer).StreamServer(m, &serverStreamServerServer{stream})
}

type Server_StreamServerServer interface {
	Send(*StreamRes) error
	grpc.ServerStream
}

type serverStreamServerServer struct {
	grpc.ServerStream
}

func (x *serverStreamServerServer) Send(m *StreamRes) error {
	return x.ServerStream.SendMsg(m)
}

func _Server_StreamClient_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServerServer).StreamClient(&serverStreamClientServer{stream})
}

type Server_StreamClientServer interface {
	SendAndClose(*StreamRes) error
	Recv() (*StreamReq, error)
	grpc.ServerStream
}

type serverStreamClientServer struct {
	grpc.ServerStream
}

func (x *serverStreamClientServer) SendAndClose(m *StreamRes) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serverStreamClientServer) Recv() (*StreamReq, error) {
	m := new(StreamReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Server_StreamTwo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServerServer).StreamTwo(&serverStreamTwoServer{stream})
}

type Server_StreamTwoServer interface {
	Send(*StreamRes) error
	Recv() (*StreamReq, error)
	grpc.ServerStream
}

type serverStreamTwoServer struct {
	grpc.ServerStream
}

func (x *serverStreamTwoServer) Send(m *StreamRes) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serverStreamTwoServer) Recv() (*StreamReq, error) {
	m := new(StreamReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Server_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "cmd",
			Handler:    _Server_Cmd_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamServer",
			Handler:       _Server_StreamServer_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamClient",
			Handler:       _Server_StreamClient_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamTwo",
			Handler:       _Server_StreamTwo_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "client_service.proto",
}
