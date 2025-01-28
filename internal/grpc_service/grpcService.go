package grpc_service

import (
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"context"
)

type GrpcService struct {
	homeSyncGrpc.UnimplementedHomeSyncGrpcServiceServer
	sensorsStorage storage.ImplStorage
}

// rpc GetHistorySensorData(HistorySensorDataRequest) returns (HistorySensorsDataResponse);
func (s *GrpcService) GetSensors(ctx context.Context, r *homeSyncGrpc.SensorsRequest) {

}
