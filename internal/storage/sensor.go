package storage

import (
	"HomeSyncService/internal/transport"
	"container/list"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type Sensor struct {
	sync.RWMutex
	Id         string
	Data       *list.List
	Type       grpc.SensorData_Type
	maxSize    int
	LastUpdate *timestamppb.Timestamp
}

func NewSensor(id string, typeSensor int, maxSize int) *Sensor {
	return &Sensor{
		Id:         id,
		Data:       list.New(),
		Type:       grpc.SensorData_Type(typeSensor),
		maxSize:    maxSize,
		LastUpdate: timestamppb.New(time.Now()),
	}
}

func (s *Sensor) Get() *grpc.BasicSensorData {
	s.RLock()
	defer s.RUnlock()
	r := s.Data.Back()
	return r.Value.(*grpc.BasicSensorData)
}

func (s *Sensor) AddData(value float32) {
	s.Lock()
	defer s.Unlock()
	s.Data.PushBack(&grpc.BasicSensorData{Time: timestamppb.Now(), Value: value})
	s.LastUpdate = timestamppb.New(time.Now())
	if s.Data.Len() > s.maxSize {
		s.Data.Remove(s.Data.Front())
	}
}
func (s *Sensor) Clear() {
	s.Lock()
	defer s.Unlock()
	e := s.Data.Front()
	for 0 < s.Data.Len() {
		e1 := e.Next()
		s.Data.Remove(e)
		e = e1
	}
}

func (s *Sensor) GetProto() *grpc.HistorySensorsDataResponse {
	s.RLock()
	defer s.RUnlock()
	sensorData := make([]*grpc.BasicSensorData, 0, s.Data.Len())
	for e := s.Data.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value.(*grpc.BasicSensorData).Value)
		sensorData = append(sensorData, e.Value.(*grpc.BasicSensorData))
	}
	resultData := grpc.HistorySensorsDataResponse{SensorData: sensorData}
	return &resultData
}
