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
	_token         string
}

func NewGrpcService(sensorsStorage storage.ImplStorage, addr string, token string) *GrpcService {
	fmt.Println("grpc service listening on " + addr)
	s := grpc.NewServer()
	grpcService := &GrpcService{sensorsStorage: sensorsStorage, _token: token}
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
	if r.GetToken() != s._token {
		return &homeSyncGrpc.SensorsResponse{Response: &homeSyncGrpc.SensorsResponse_Error{
			Error: &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
		},
		}, fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.GetSensorsData(), nil
}

// GetHistorySensorData Возвращает историю конкретного датчика, для построения графика
func (s *GrpcService) GetHistorySensorData(ctx context.Context, r *homeSyncGrpc.HistorySensorDataRequest) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	if r.GetToken() != s._token {
		return &homeSyncGrpc.HistorySensorsDataResponse{Response: &homeSyncGrpc.HistorySensorsDataResponse_Error{
			Error: &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
		},
		}, fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.GetHistoricSensorsData(r.GetBlockId(), r.GetSensorId())
}
