package mqtt_service

import (
	"encoding/json"
	"fmt"
	"github.com/mochi-mqtt/server/v2/packets"
)

type BlockData struct {
	Sensors []SensorData `json:"sensors"`
}
type SensorData struct {
	Id    string  `json:"id"`
	Type  int     `json:"type"`
	Value float32 `json:"value"`
}

type MainBlockHandler struct {
	id        string
	topicName string
}

func NewMainBlockHandler(topicName string) *MainBlockHandler {
	return &MainBlockHandler{
		id:        "steadfastness",
		topicName: topicName,
	}
}

func (m *MainBlockHandler) TopicName() string {
	return m.topicName
}

func (m *MainBlockHandler) Id() string {
	return m.id
}

func (m *MainBlockHandler) Parse(pk packets.Packet) ([]SensorData, error) {
	block := &BlockData{}
	if err := json.Unmarshal(pk.Payload, &block); err != nil {
		return nil, fmt.Errorf("ошибка парсинга сенсорного массива: %v", err)
	}
	return block.Sensors, nil
}
