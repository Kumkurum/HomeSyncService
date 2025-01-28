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
	Blocks     map[string]*Block
	maxSize    int                    //Максимальных размер хранения информации о датчике ( история изменений)
	LastUpdate *timestamppb.Timestamp //Время послденего обновления
}

// NewStorage - Создание нового хранилища
func NewStorage(maxSize int) *Storage {
	return &Storage{
		Blocks:  make(map[string]*Block),
		maxSize: maxSize,
	}
}

// UpdateSensorValue - Обновление или добавление нового датчика и определение его в какой-то блок
func (s *Storage) UpdateSensorValue(blockId string, sensorId string, typeSensor int, value float32) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Blocks[blockId]; ok == false {
		s.Blocks[blockId] = NewBlock(blockId, s.maxSize)
	}
	s.Blocks[blockId].UpdateSensor(sensorId, typeSensor, value)
}
func (s *Storage) GetHistoricSensorsData(blockId string, sensorId string) (error, *homeSyncGrpc.HistorySensorsDataResponse) {
	s.RLock()
	defer s.RUnlock()
	_, err := s.Blocks[blockId].GetSensor(sensorId)
	if err != nil {
		return fmt.Errorf("not Found sensor with %s id", sensorId), nil
	}
	sensor, _ := s.Blocks[blockId].GetSensor(sensorId)
	return nil, sensor.GetProto()
}

func (s *Storage) GetSensorsData() *homeSyncGrpc.SensorsResponse {
	s.RLock()
	defer s.RUnlock()
	result := &homeSyncGrpc.SensorsResponse{
		Time:       s.LastUpdate,
		GroupsData: make([]*homeSyncGrpc.GroupData, 0, len(s.Blocks)),
	}
	for blockId, block := range s.Blocks {
		blockData := block.GetBlockSensors()
		blockData.Id = blockId
		result.GroupsData = append(result.GroupsData, blockData)
	}
	return result
}
