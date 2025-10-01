package mqtt_service

import (
	"github.com/mochi-mqtt/server/v2/packets"
)

type BlockHandler interface {
	TopicName() string
	Id() string
	Parse(packets.Packet) ([]SensorData, error)
}
