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
	addr           string
	logger         *log_client.LoggingClient
}

func NewGrpcService(sensorsStorage storage.ImplStorage, addr string, token string, logger *log_client.LoggingClient) (*GrpcService, error) {
	_ = logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Action", Value: "Initialize"},
	)
	return &GrpcService{
		sensorsStorage: sensorsStorage,
		_token:         token,
		logger:         logger,
		addr:           addr,
	}, nil
}

func (s *GrpcService) Run() error {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Action", Value: "Start"},
	)

	serv := grpc.NewServer()
	homeSyncGrpc.RegisterHomeSyncGrpcServiceServer(serv, s)

	// Открыть порт для приема сообщений
	lis, err := net.Listen("tcp", ":"+s.addr)
	if err != nil {
		_ = s.logger.Critical(
			log_client.KeyValue{Key: "Layer", Value: "GRPC"},
			log_client.KeyValue{Key: "Function", Value: "Start"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
			log_client.KeyValue{Key: "Failed to listen port", Value: s.addr},
		)
		return fmt.Errorf("failed to listen on port %s: %w", s.addr, err)
	}

	// Начать цикл приема и обработку запросов
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Action", Value: "Server started"},
		log_client.KeyValue{Key: "Port", Value: s.addr},
	)

	if err := serv.Serve(lis); err != nil {
		_ = s.logger.Critical(
			log_client.KeyValue{Key: "Layer", Value: "GRPC"},
			log_client.KeyValue{Key: "Function", Value: "Start"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
			log_client.KeyValue{Key: "Failed to Serve on port", Value: s.addr},
		)
		return fmt.Errorf("failed to serve on port %s: %w", s.addr, err)
	}

	return nil
}

// GetSensors - возвращает массив сенсоров зарегестрированных на сервере
func (s *GrpcService) GetSensors(_ context.Context, r *homeSyncGrpc.SensorsRequest) (*homeSyncGrpc.SensorsResponse, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Function", Value: "GetSensors"},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "GRPC"},
			log_client.KeyValue{Key: "Function", Value: "GetSensors"},
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
func (s *GrpcService) GetHistorySensorData(_ context.Context, r *homeSyncGrpc.HistorySensorDataRequest) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Function", Value: "GetHistorySensorData"},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "GRPC"},
			log_client.KeyValue{Key: "Function", Value: "GetHistorySensorData"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.HistorySensorsDataResponse{Response: &homeSyncGrpc.HistorySensorsDataResponse_Error{
			Error: &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
		},
		}, fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.GetHistoricSensorsData(r.GetSensorId())
}

// SetBoundary Устанавливает имя для конкретного сенсора
func (s *GrpcService) SetBoundary(_ context.Context, r *homeSyncGrpc.SetBoundaryRequest) (*homeSyncGrpc.Error, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
		log_client.KeyValue{Key: "Sensor", Value: r.SensorId},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "GRPC"},
			log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
			fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.SetBoundary(r), nil
}

func (s *GrpcService) RemoveSensor(_ context.Context, r *homeSyncGrpc.RemoveSensorRequest) (*homeSyncGrpc.Error, error) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "GRPC"},
		log_client.KeyValue{Key: "Function", Value: "RemoveSensor"},
		log_client.KeyValue{Key: "Sensor", Value: r.SensorId},
	)
	if r.GetToken() != s._token {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "GRPC"},
			log_client.KeyValue{Key: "Function", Value: "RemoveSensor"},
			log_client.KeyValue{Key: "Error", Value: "invalid token"},
			log_client.KeyValue{Key: "Token", Value: r.GetToken()},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_TOKEN},
			fmt.Errorf("invalid token")
	}
	return s.sensorsStorage.RemoveSensor(r), nil
}
