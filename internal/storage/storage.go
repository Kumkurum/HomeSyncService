package storage

import (
	homeSyncGrpc "HomeSyncService/internal/transport"
	"github.com/Kumkurum/LogService/pkg/log_client"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"sync"
)

// Storage - Структура для хранения информации о датчиках, разбитая на блоки, объединенная по каким-то принципам
type Storage struct {
	sync.RWMutex
	sensors    map[string]*Sensor
	maxSize    int                    //Максимальных размер хранения информации о датчике ( история изменений)
	lastUpdate *timestamppb.Timestamp //Время последнего обновления
	logger     *log_client.LoggingClient
}

// NewStorage - Создание нового хранилища
func NewStorage(maxSize int, logger *log_client.LoggingClient) *Storage {
	return &Storage{
		sensors: make(map[string]*Sensor),
		maxSize: maxSize,
		logger:  logger,
	}
}

// UpdateSensorValue - Обновление или добавление нового датчика и определение его в какой-то блок
func (s *Storage) UpdateSensorValue(sensorId string, typeSensor int, value float32) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "UpdateSensorValue"},
		log_client.KeyValue{Key: "sensorId", Value: sensorId},
		log_client.KeyValue{Key: "typeSensor", Value: strconv.Itoa(typeSensor)},
	)
	s.Lock()
	defer s.Unlock()
	if _, ok := s.sensors[sensorId]; ok == false {
		s.sensors[sensorId] = NewSensor(typeSensor, s.maxSize)
	}
	s.sensors[sensorId].AddData(value)
}

func (s *Storage) GetHistoricSensorsData(sensorId string) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "GetHistoricSensorsData"},
		log_client.KeyValue{Key: "sensorId", Value: sensorId},
	)
	s.RLock()
	defer s.RUnlock()
	sensor, ok := s.sensors[sensorId]
	if !ok {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "Storage"},
			log_client.KeyValue{Key: "Function", Value: "GetHistoricSensorsData"},
			log_client.KeyValue{Key: "Error", Value: "sensor " + sensorId + " not found"},
		)
		return &homeSyncGrpc.HistorySensorsDataResponse{Response: &homeSyncGrpc.HistorySensorsDataResponse_Error{
			Error: &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID},
		}}, nil
	}
	return sensor.GetHistory(), nil
}

func (s *Storage) GetSensorsData() *homeSyncGrpc.SensorsResponse {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "GetSensorsData"},
	)
	s.RLock()
	defer s.RUnlock()
	success := &homeSyncGrpc.SensorsResponseSuccess{
		Time:    s.lastUpdate,
		Sensors: make([]*homeSyncGrpc.SensorData, 0, len(s.sensors)),
	}
	for sensorId, sensor := range s.sensors {
		sensorData := sensor.Get()
		sensorData.Id = sensorId
		success.Sensors = append(success.Sensors, sensorData)
	}
	result := &homeSyncGrpc.SensorsResponse{
		Response: &homeSyncGrpc.SensorsResponse_Success{
			Success: success,
		},
	}
	return result
}

func (s *Storage) SetBoundary(request *homeSyncGrpc.SetBoundaryRequest) *homeSyncGrpc.Error {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
		log_client.KeyValue{Key: "sensorId", Value: request.SensorId},
	)
	s.RLock()
	defer s.RUnlock()
	sensor, ok := s.sensors[request.SensorId]
	if ok == false {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "Storage"},
			log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
			log_client.KeyValue{Key: "Error", Value: "Unknown Sensor Id"},
			log_client.KeyValue{Key: "Sensor Id", Value: request.SensorId},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID}
	}
	sensor.Boundary.Value1 = request.Boundary.Value1
	sensor.Boundary.Value2 = request.Boundary.Value2
	sensor.Boundary.Value3 = request.Boundary.Value3
	sensor.Boundary.Value4 = request.Boundary.Value4
	return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_OK}
}

func (s *Storage) RemoveSensor(request *homeSyncGrpc.RemoveSensorRequest) *homeSyncGrpc.Error {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "RemoveSensorRequest"},
		log_client.KeyValue{Key: "sensorId", Value: request.SensorId},
	)
	s.RLock()
	defer s.RUnlock()
	_, ok := s.sensors[request.SensorId]
	if ok == false {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "Storage"},
			log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
			log_client.KeyValue{Key: "Error", Value: "Unknown Sensor Id"},
			log_client.KeyValue{Key: "SensorId", Value: request.SensorId},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID}
	}
	delete(s.sensors, request.SensorId)
	return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_OK}
}
