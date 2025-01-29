package grpc_service

import (
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"context"
	"fmt"
)

type GrpcService struct {
	homeSyncGrpc.UnimplementedHomeSyncGrpcServiceServer
	sensorsStorage storage.ImplStorage
}

func NewGrpcService(sensorsStorage storage.ImplStorage) *GrpcService {
	return &GrpcService{sensorsStorage: sensorsStorage}
}

// GetSensors - возвращает массив сенсоров зарегестрированных на сервере
func (s *GrpcService) GetSensors(ctx context.Context, r *homeSyncGrpc.SensorsRequest) (*homeSyncGrpc.SensorsResponse, error) {
	fmt.Println("GetSensors")
	return s.sensorsStorage.GetSensorsData(), nil
}

func (s *GrpcService) GetHistorySensorData(ctx context.Context, r *homeSyncGrpc.HistorySensorDataRequest) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	fmt.Println("GetHistorySensorData")
	return s.sensorsStorage.GetHistoricSensorsData(r.GetBlockId(), r.GetSensorId())
}
