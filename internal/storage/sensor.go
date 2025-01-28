package storage

import (
	homeSyncGrpc "HomeSyncService/internal/transport"
	"container/list"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

// Sensor - Структура для хранения инфы об Датчике или девайсе
type Sensor struct {
	sync.RWMutex
	Data    *list.List                   //Data - стэк данных с изменениями значений Sensor по времени
	Type    homeSyncGrpc.SensorData_Type //Type - тип датчика
	maxSize int                          //Максимальный размер стэка с данными
}

// NewSensor - Создание нового датчика
func NewSensor(id string, typeSensor int, maxSize int) *Sensor {
	return &Sensor{
		Data:    list.New(),
		Type:    homeSyncGrpc.SensorData_Type(typeSensor),
		maxSize: maxSize,
	}
}

// Get - Получение последней инфы о состоянии датчика (Mutex)
func (s *Sensor) Get() *homeSyncGrpc.SensorData {
	s.RLock()
	defer s.RUnlock()
	r := s.Data.Back()
	return &homeSyncGrpc.SensorData{
		BasicData: r.Value.(*homeSyncGrpc.BasicSensorData),
		Type:      s.Type,
	}
}

// AddData - Обновление данных на датчике (Mutex)
func (s *Sensor) AddData(value float32) {
	s.Lock()
	defer s.Unlock()
	s.Data.PushBack(&homeSyncGrpc.BasicSensorData{Time: timestamppb.Now(), Value: value})
	if s.Data.Len() > s.maxSize {
		s.Data.Remove(s.Data.Front())
	}
}

// Clear - Отчистка всей истории о датчике
func (s *Sensor) Clear() {
	s.Lock()
	defer s.Unlock()
	s.Data = list.New().Init()
}

// GetProto - Получение прото сообщения со всей историей обновления данных датчика
func (s *Sensor) GetProto() *homeSyncGrpc.HistorySensorsDataResponse {
	s.RLock()
	defer s.RUnlock()
	sensorData := make([]*homeSyncGrpc.BasicSensorData, 0, s.Data.Len())
	for e := s.Data.Front(); e != nil; e = e.Next() {
		sensorData = append(sensorData, e.Value.(*homeSyncGrpc.BasicSensorData))
	}
	resultData := &homeSyncGrpc.HistorySensorsDataResponse{SensorData: sensorData}
	return resultData
}
