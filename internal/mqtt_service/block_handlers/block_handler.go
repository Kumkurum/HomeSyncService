package mqtt_service

import (
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

type BlockHandler interface {
	TopicName() string
	Id() string
	Parse(packets.Packet) ([]SensorData, error)
}
