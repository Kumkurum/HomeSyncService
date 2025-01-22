package storage

import (
	"HomeSyncService/internal/transport"
	"container/list"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

type Sensor struct {
	sync.RWMutex
	Name    string
	Data    *list.List
	maxSize int
}

func NewSensor(name string, maxSize int) *Sensor {
	return &Sensor{
		Name:    name,
		Data:    list.New(),
		maxSize: maxSize,
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
	s.Data.PushBack(grpc.BasicSensorData{Time: timestamppb.Now(), Value: value})
	if s.Data.Len() > s.maxSize {
		s.Data.Remove(s.Data.Front())
	}
}

func (s *Sensor) getProto() grpc.HistorySensorsDataResponse {
	resultData := grpc.HistorySensorsDataResponse{}
	for e := s.Data.Front(); e != nil; e = e.Next() {
		//element:= grpc.SensorDataPoint{Value: e.Value.(SensorData).Value, Time: e.Value.(SensorData).TimeStamp.()}
		//append(resultData.SensorData, element)
	}
	return resultData
}
