package storage

import "sync"

type Storage struct {
	sync.RWMutex
	Blocks  map[string]*Block
	maxSize int
}

func NewStorage(maxSize int) *Storage {
	return &Storage{
		Blocks:  make(map[string]*Block),
		maxSize: maxSize,
	}
}

func (s *Storage) AddBlock(block Block) {
	s.Lock()
	defer s.Unlock()

}

func (s *Storage) UpdateSensorValue(blockId string, sensorId string, typeSensor int, value float32) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Blocks[blockId]; ok == false {
		s.Blocks[blockId] = NewBlock(blockId, s.maxSize)
	}
	s.Blocks[blockId].UpdateSensor(sensorId, typeSensor, value)
}
