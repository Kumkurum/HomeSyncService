package grpc_service

import (
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcService struct {
	homeSyncGrpc.UnimplementedHomeSyncGrpcServiceServer
	sensorsStorage storage.ImplStorage
}

func NewGrpcService(sensorsStorage storage.ImplStorage, addr string) *GrpcService {
	fmt.Println("grpc service listening on " + addr)
	s := grpc.NewServer()
	grpcService := &GrpcService{sensorsStorage: sensorsStorage}
	homeSyncGrpc.RegisterHomeSyncGrpcServiceServer(s, grpcService)
	// Открыть порт 50051 для приема сообщений
	lis, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Начать цикл приема и обработку запросов
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return grpcService
}

// GetSensors - возвращает массив сенсоров зарегестрированных на сервере
func (s *GrpcService) GetSensors(ctx context.Context, r *homeSyncGrpc.SensorsRequest) (*homeSyncGrpc.SensorsResponse, error) {
	return s.sensorsStorage.GetSensorsData(), nil
}

// GetHistorySensorData ВОзвращает историю конкретного датчика, для построения графика
func (s *GrpcService) GetHistorySensorData(ctx context.Context, r *homeSyncGrpc.HistorySensorDataRequest) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	return s.sensorsStorage.GetHistoricSensorsData(r.GetBlockId(), r.GetSensorId())
}
