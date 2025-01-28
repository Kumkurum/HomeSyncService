// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.3
// source: service.proto

package homeSyncGrpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x96,
	0x01, 0x0a, 0x13, 0x48, 0x6f, 0x6d, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x47, 0x72, 0x70, 0x63, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x53, 0x65, 0x6e,
	0x73, 0x6f, 0x72, 0x73, 0x12, 0x0f, 0x2e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x19, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1c, 0x5a, 0x1a, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x6d, 0x65, 0x53, 0x79, 0x6e,
	0x63, 0x47, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_proto_goTypes = []any{
	(*SensorsRequest)(nil),             // 0: SensorsRequest
	(*HistorySensorDataRequest)(nil),   // 1: HistorySensorDataRequest
	(*SensorsResponse)(nil),            // 2: SensorsResponse
	(*HistorySensorsDataResponse)(nil), // 3: HistorySensorsDataResponse
}
var file_service_proto_depIdxs = []int32{
	0, // 0: HomeSyncGrpcService.GetSensors:input_type -> SensorsRequest
	1, // 1: HomeSyncGrpcService.GetHistorySensorData:input_type -> HistorySensorDataRequest
	2, // 2: HomeSyncGrpcService.GetSensors:output_type -> SensorsResponse
	3, // 3: HomeSyncGrpcService.GetHistorySensorData:output_type -> HistorySensorsDataResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	file_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
