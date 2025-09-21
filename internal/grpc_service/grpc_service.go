package grpc_service

import (
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"context"
	"fmt"
	"github.com/Kumkurum/LogService/pkg/log_client"
	"google.golang.org/grpc"
	"net"
)

type GrpcService struct {
	homeSyncGrpc.UnimplementedHomeSyncGrpcServiceServer
	sensorsStorage storage.ImplStorage
	_token         string
	logger         *log_client.LoggingClient
}

func NewGrpcService(sensorsStorage storage.ImplStorage, addr string, token string, logger *log_client.LoggingClient) *GrpcService {
	_ = logger.Info(
		log_client.KeyValue{Key: "Service", Value: "HomeSync"},
		log_client.KeyValue{Key: "Class", Value: "GrpcService"},
		log_client.KeyValue{Key: "Action", Value: "Start"},
	)
	s := grpc.NewServer()
	grpcService := &GrpcService{sensorsStorage: sensorsStorage, _token: token, logger: logger}
	homeSyncGrpc.RegisterHomeSyncGrpcServiceServer(s, grpcService)
	// Открыть порт 50051 для приема сообщений
	lis, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		_ = logger.Critical(
			log_client.KeyValue{Key: "Function", Value: "NewGrpcService"},
			log_client.KeyValue{Key: "Class", Value: "GrpcService"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
			log_client.KeyValue{Key: "Failed to listen port", Value: addr},
		)
	}
	// Начать цикл приема и обработку запросов
	if err := s.Serve(lis); err != nil {
		_ = logger.Critical(
			log_client.KeyValue{Key: "Function", Value: "NewGrpcService"},
			log_client.KeyValue{Key: "Class", Value: "GrpcService"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
			log_client.KeyValue{Key: "Failed to Serve on port", Value: addr},
		)
	}
	return grpcService
}

// GetSensors - возвращает массив сенсоров зарегестрированных на сервере
func (s *GrpcService) GetSensors(ctx context.Context, r *homeSyncGrpc.SensorsRequest) (*homeSyncGrpc.SensorsResponse, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Function", Value: "GetSensors"},
		log_client.KeyValue{Key: "Class", Value: "GrpcService"},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "GetSensors"},
			log_client.KeyValue{Key: "Class", Value: "GrpcService"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.SensorsResponse{Response: &homeSyncGrpc.SensorsResponse_Error{
			Error: &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
		},
		}, fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.GetSensorsData(), nil
}

// GetHistorySensorData Возвращает историю конкретного датчика, для построения графика
func (s *GrpcService) GetHistorySensorData(ctx context.Context, r *homeSyncGrpc.HistorySensorDataRequest) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Function", Value: "GetHistorySensorData"},
		log_client.KeyValue{Key: "Class", Value: "GrpcService"},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "GetSensors"},
			log_client.KeyValue{Key: "Class", Value: "GrpcService"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.HistorySensorsDataResponse{Response: &homeSyncGrpc.HistorySensorsDataResponse_Error{
			Error: &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
		},
		}, fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.GetHistoricSensorsData(r.GetBlockId(), r.GetSensorId())
}

// SetName Устанавливает имя для конкретного сенсора
func (s *GrpcService) SetName(ctx context.Context, r *homeSyncGrpc.SetSensorData) (*homeSyncGrpc.Error, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Function", Value: "SetName"},
		log_client.KeyValue{Key: "Class", Value: "GrpcService"},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "GetSensors"},
			log_client.KeyValue{Key: "Class", Value: "GrpcService"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
			fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.SetSensorName(r.BlockId, r.SensorId, r.Name), nil
}

// SetBoundary Устанавливает имя для конкретного сенсора
func (s *GrpcService) SetBoundary(ctx context.Context, r *homeSyncGrpc.SetBoundaryRequest) (*homeSyncGrpc.Error, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
		log_client.KeyValue{Key: "Class", Value: "GrpcService"},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "GetSensors"},
			log_client.KeyValue{Key: "Class", Value: "GrpcService"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
			fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.SetBoundary(r), nil
}
