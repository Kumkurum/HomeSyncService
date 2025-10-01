package mqtt_service

import (
	mqttservice2 "HomeSyncService/internal/mqtt_service/block_handlers"
	"HomeSyncService/internal/storage"
	"fmt"
	"github.com/Kumkurum/LogService/pkg/log_client"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type BlockManagerOptions struct {
	handlers map[string]mqttservice2.BlockHandler
	logger   *log_client.LoggingClient
	storage  storage.ImplStorage
}
type BlockManager struct {
	mqtt.HookBase
	config *BlockManagerOptions
}

func (b *BlockManager) AddBlock(handler mqttservice2.BlockHandler) {
	b.config.logger.Debug(
		log_client.KeyValue{Key: "Layer", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "AddBlock"},
		log_client.KeyValue{Key: "Topic", Value: handler.TopicName()},
	)
	b.config.handlers[handler.TopicName()] = handler
}

func (b *BlockManager) ID() string {
	return "universal-message-processor"
}

func (b *BlockManager) Provides(bt byte) bool {
	return true
	//return bt == mqtt.OnMessage // Используйте hooks.OnMessage вместо mqtt.OnMessage
}

func (b *BlockManager) Init(config any) error {
	if _, ok := config.(*BlockManagerOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	b.config = config.(*BlockManagerOptions)

	b.AddBlock(mqttservice2.NewMainBlockHandler("home/main_block"))
	b.config.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "BlockManager Init"})
	return nil
}

func (b *BlockManager) GetHandler(id string) mqttservice2.BlockHandler {
	if handler, ok := b.config.handlers[id]; ok {
		return handler
	}
	return nil
}
func (b *BlockManager) OnPublish(_ *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	handler := b.GetHandler(pk.TopicName)
	b.config.logger.Critical(
		log_client.KeyValue{Key: "Layer", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "OnMessage"},
	)
	if handler == nil {
		b.config.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "MQTT"},
			log_client.KeyValue{Key: "Action", Value: "OnMessage"},
			log_client.KeyValue{Key: "Topic", Value: pk.TopicName},
			log_client.KeyValue{Key: "Error", Value: "topic not found"},
		)
		return packets.Packet{}, fmt.Errorf("topic not found: %s", pk.TopicName)
	}
	b.config.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "OnMessage"},
		log_client.KeyValue{Key: "Topic", Value: pk.TopicName})

	sensors, err := handler.Parse(pk)
	if err != nil {
		b.config.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "MQTT"},
			log_client.KeyValue{Key: "Action", Value: "OnMessage"},
			log_client.KeyValue{Key: "Block", Value: pk.TopicName},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
		return packets.Packet{}, err
	}
	for _, sensor := range sensors {
		b.config.storage.UpdateSensorValue(handler.Id(), sensor.Id, sensor.Type, sensor.Value)
	}
	return pk, nil
}

func (b *BlockManager) OnMessage(_ *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	handler := b.GetHandler(pk.TopicName)
	b.config.logger.Critical(
		log_client.KeyValue{Key: "Layer", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "OnMessage"},
	)
	if handler == nil {
		b.config.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "MQTT"},
			log_client.KeyValue{Key: "Action", Value: "OnMessage"},
			log_client.KeyValue{Key: "Topic", Value: pk.TopicName},
			log_client.KeyValue{Key: "Error", Value: "topic not found"},
		)
		return packets.Packet{}, fmt.Errorf("topic not found: %s", pk.TopicName)
	}
	b.config.logger.Info(
		log_client.KeyValue{Key: "Layer", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "OnMessage"},
		log_client.KeyValue{Key: "Topic", Value: pk.TopicName})

	sensors, err := handler.Parse(pk)
	if err != nil {
		b.config.logger.Warn(
			log_client.KeyValue{Key: "Layer", Value: "MQTT"},
			log_client.KeyValue{Key: "Action", Value: "OnMessage"},
			log_client.KeyValue{Key: "Block", Value: pk.TopicName},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
		return packets.Packet{}, err
	}
	for _, sensor := range sensors {
		b.config.storage.UpdateSensorValue(handler.Id(), sensor.Id, sensor.Type, sensor.Value)
	}
	return pk, nil
}
