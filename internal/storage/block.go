package storage

import (
	"fmt"
)

type Block struct {
	Id      string
	sensors map[string]*Sensor
	MaxSize int
}

func NewBlock(id string, maxSize int) *Block {
	return &Block{
		Id:      id,
		sensors: make(map[string]*Sensor),
		MaxSize: maxSize,
	}
}

func (b *Block) UpdateSensor(id string, typeSensor int, value float32) {
	_, ok := b.sensors[id]
	if !ok {
		b.sensors[id] = NewSensor(id, typeSensor, b.MaxSize)
	}
	b.sensors[id].AddData(value)
}

func (b *Block) GetSensor(id string) (*Sensor, error) {
	_, ok := b.sensors[id]
	fmt.Println("GetSensor")
	if !ok {
		return nil, fmt.Errorf("sensor %s not found", id)
	}
	return b.sensors[id], nil
}
