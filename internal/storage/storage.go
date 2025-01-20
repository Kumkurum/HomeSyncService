package storage

import "sync"

type Storage struct {
	sync.RWMutex
	Blocks map[string]Block
}

func NewStorage() *Storage {
	return &Storage{
		Blocks: make(map[string]Block),
	}
}

func (s *Storage) AddBlock(block Block) {
	s.Lock()
	defer s.Unlock()

}
