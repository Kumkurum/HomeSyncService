package storage

import (
	homeSyncGrpc "HomeSyncService/internal/transport"
	"fmt"
	"github.com/Kumkurum/LogService/pkg/log_client"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"sync"
)

// Storage - Структура для хранения информации о датчиках, разбитая на блоки, объединенная по каким-то принципам
type Storage struct {
	sync.RWMutex
	blocks     map[string]*Block
	maxSize    int                    //Максимальных размер хранения информации о датчике ( история изменений)
	lastUpdate *timestamppb.Timestamp //Время последнего обновления
	logger     *log_client.LoggingClient
}

// NewStorage - Создание нового хранилища
func NewStorage(maxSize int, logger *log_client.LoggingClient) *Storage {
	return &Storage{
		blocks:  make(map[string]*Block),
		maxSize: maxSize,
		logger:  logger,
	}
}

// UpdateSensorValue - Обновление или добавление нового датчика и определение его в какой-то блок
func (s *Storage) UpdateSensorValue(blockId string, sensorId string, typeSensor int, value float32) {
	_ = s.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "UpdateSensorValue"},
		log_client.KeyValue{Key: "blocId", Value: blockId},
		log_client.KeyValue{Key: "sensorId", Value: sensorId},
		log_client.KeyValue{Key: "typeSensor", Value: strconv.Itoa(typeSensor)},
	)
	s.Lock()
	defer s.Unlock()
	if _, ok := s.blocks[blockId]; ok == false {
		s.blocks[blockId] = NewBlock(blockId, s.maxSize)
	}
	s.blocks[blockId].UpdateSensor(sensorId, typeSensor, value)
}

func (s *Storage) GetHistoricSensorsData(blockId string, sensorId string) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "GetHistoricSensorsData"},
		log_client.KeyValue{Key: "blocId", Value: blockId},
		log_client.KeyValue{Key: "sensorId", Value: sensorId},
	)
	s.RLock()
	defer s.RUnlock()
	_, err := s.blocks[blockId].GetSensor(sensorId)
	if err != nil {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "Storage"},
			log_client.KeyValue{Key: "Function", Value: "GetHistoricSensorsData"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
		return nil, fmt.Errorf("not Found sensor with %s id", sensorId)
	}
	sensor, _ := s.blocks[blockId].GetSensor(sensorId)
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
		Time:       s.lastUpdate,
		GroupsData: make([]*homeSyncGrpc.GroupData, 0, len(s.blocks)),
	}
	for blockId, block := range s.blocks {
		blockData := block.GetBlockSensors()
		blockData.Id = blockId
		success.GroupsData = append(success.GroupsData, blockData)
	}
	result := &homeSyncGrpc.SensorsResponse{
		Response: &homeSyncGrpc.SensorsResponse_Success{
			Success: success,
		},
	}
	return result
}

func (s *Storage) SetSensorName(blockId string, sensorId string, sensorName string) *homeSyncGrpc.Error {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "SetSensorName"},
		log_client.KeyValue{Key: "blockId", Value: blockId},
		log_client.KeyValue{Key: "sensorId", Value: sensorId},
		log_client.KeyValue{Key: "sensorName", Value: sensorName},
	)
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.blocks[blockId]; ok == false {
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID}
	}
	sensor, err := s.blocks[blockId].GetSensor(sensorId)
	if err != nil {
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID}
	}
	sensor.Name = sensorName
	return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_OK}
}

func (s *Storage) SetBoundary(request *homeSyncGrpc.SetBoundaryRequest) *homeSyncGrpc.Error {
	_ = s.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "Storage"},
		log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
		log_client.KeyValue{Key: "blockId", Value: request.BlockId},
		log_client.KeyValue{Key: "sensorId", Value: request.SensorId},
	)
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.blocks[request.BlockId]; ok == false {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "Storage"},
			log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
			log_client.KeyValue{Key: "Error", Value: "Unknown Block Id"},
			log_client.KeyValue{Key: "BlockId", Value: request.BlockId},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID}
	}
	sensor, err := s.blocks[request.BlockId].GetSensor(request.SensorId)
	if err != nil {
		_ = s.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "Storage"},
			log_client.KeyValue{Key: "Function", Value: "SetBoundary"},
			log_client.KeyValue{Key: "Error", Value: "Unknown Sensor Id"},
			log_client.KeyValue{Key: "SensorId", Value: request.SensorId},
		)
		return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_UNKNOWN_ID}
	}
	sensor.Boundary.Value1 = request.Boundary.Value1
	sensor.Boundary.Value2 = request.Boundary.Value2
	sensor.Boundary.Value3 = request.Boundary.Value3
	sensor.Boundary.Value4 = request.Boundary.Value4
	return &homeSyncGrpc.Error{Code: homeSyncGrpc.Error_OK}
}
