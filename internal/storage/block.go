package storage

import (
	homeSyncGrpc "HomeSyncService/internal/transport"
	"fmt"
)

// Block - Группировка датчиков в блоки ( например по комнатам или иным признакам соответствия )
// Содержит в себе map из набора датчиком
type Block struct {
	Id      string
	sensors map[string]*Sensor
	MaxSize int
}

// NewBlock - Создание нового блока с датчиками
func NewBlock(id string, maxSize int) *Block {
	return &Block{
		Id:      id,
		sensors: make(map[string]*Sensor),
		MaxSize: maxSize,
	}
}

// UpdateSensor - Обновление инфы о датчике, если такого id нету, то добавляем его в map
func (b *Block) UpdateSensor(id string, typeSensor int, value float32) {
	_, ok := b.sensors[id]
	if !ok {
		b.sensors[id] = NewSensor(id, typeSensor, b.MaxSize)
	}
	b.sensors[id].AddData(value)
}

// GetSensor - получение указателя на датчик по его id, если его не будет, вернётся error
func (b *Block) GetSensor(id string) (*Sensor, error) {
	_, ok := b.sensors[id]
	if !ok {
		return nil, fmt.Errorf("sensor %s not found", id)
	}
	return b.sensors[id], nil
}

// GetBlockSensors - группировка в прото сообщение всех данных о датчиках в группе
func (b *Block) GetBlockSensors() *homeSyncGrpc.GroupData {
	result := &homeSyncGrpc.GroupData{
		Id:          b.Id,
		SensorsData: make([]*homeSyncGrpc.SensorData, 0, len(b.sensors)),
	}
	for sensorId, sensor := range b.sensors {
		sensorData := sensor.Get()
		sensorData.Id = sensorId
		result.SensorsData = append(result.SensorsData, sensorData)
	}
	return result
}
