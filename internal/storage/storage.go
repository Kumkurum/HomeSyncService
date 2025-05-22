package storage

import (
	homeSyncGrpc "HomeSyncService/internal/transport"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

// Storage - Структура для хранения информации о датчиках, разбитая на блоки, объединенная по каким-то принципам
type Storage struct {
	sync.RWMutex
	blocks     map[string]*Block
	maxSize    int                    //Максимальных размер хранения информации о датчике ( история изменений)
	lastUpdate *timestamppb.Timestamp //Время послденего обновления
}

// NewStorage - Создание нового хранилища
func NewStorage(maxSize int) *Storage {
	return &Storage{
		blocks:  make(map[string]*Block),
		maxSize: maxSize,
	}
}

// UpdateSensorValue - Обновление или добавление нового датчика и определение его в какой-то блок
func (s *Storage) UpdateSensorValue(blockId string, sensorId string, typeSensor int, value float32) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.blocks[blockId]; ok == false {
		s.blocks[blockId] = NewBlock(blockId, s.maxSize)
	}
	s.blocks[blockId].UpdateSensor(sensorId, typeSensor, value)
}
func (s *Storage) GetHistoricSensorsData(blockId string, sensorId string) (*homeSyncGrpc.HistorySensorsDataResponse, error) {
	s.RLock()
	defer s.RUnlock()
	_, err := s.blocks[blockId].GetSensor(sensorId)
	if err != nil {
		return nil, fmt.Errorf("not Found sensor with %s id", sensorId)
	}
	sensor, _ := s.blocks[blockId].GetSensor(sensorId)
	return sensor.GetHistory(), nil
}

func (s *Storage) GetSensorsData() *homeSyncGrpc.SensorsResponse {
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
