package mqtt_service

import (
	mqttservice2 "HomeSyncService/internal/mqtt_service/block_handlers"
	"HomeSyncService/internal/storage"

	"github.com/Kumkurum/LogService/pkg/log_client"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type MQTTServiceConfig struct {
}

type MQTTService struct {
	storage storage.ImplStorage
	server  *mqtt.Server
	logger  *log_client.LoggingClient
}

func NewMQTTService(str storage.ImplStorage, logger *log_client.LoggingClient) *MQTTService {

	// Создаем новый MQTT-сервер
	server := mqtt.New(&mqtt.Options{
		InlineClient: true, // Включаем встроенного клиента для тестирования
	})
	service := &MQTTService{storage: str, server: server, logger: logger}
	return service
}

func (m *MQTTService) Run(port, userName, password string) {
	// Добавляем аутентификацию (опционально)
	_ = m.server.AddHook(new(auth.Hook), &auth.Options{
		Ledger: &auth.Ledger{
			Auth: auth.AuthRules{ // Правила аутентификации
				{Username: auth.RString(userName), Password: auth.RString(password), Allow: true},
			},
			ACL: auth.ACLRules{ // Правила доступа
				{Username: auth.RString(userName), Filters: auth.Filters{
					"home/#": auth.ReadWrite, // ДОБАВЬТЕ ЭТУ СТРОКУ
					"$SYS/#": auth.Deny,
					"#":      auth.Deny,
				}},
			},
		},
	})

	tcp := listeners.NewTCP(listeners.Config{
		Type:      "tcp",      // Тип листенера
		ID:        "tcp-main", // Уникальный ID
		Address:   ":" + port, // Адрес и порт
		TLSConfig: nil,        // Без TLS
	})

	err := m.server.AddListener(tcp)
	if err != nil {
		m.logger.Critical(
			log_client.KeyValue{Key: "Service", Value: "MQTT"},
			log_client.KeyValue{Key: "Action", Value: "AddListener"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
	}
	m.logger.Info(
		log_client.KeyValue{Key: "Service", Value: "MQTT"},
		log_client.KeyValue{Key: "Action", Value: "MQTT сервер запущен"},
		log_client.KeyValue{Key: "Port", Value: port})

	// ОБРАБОТКА СООБЩЕНИЙ - вот где magic happens!
	err = m.server.AddHook(new(BlockManager), &BlockManagerOptions{
		handlers: make(map[string]mqttservice2.BlockHandler),
		logger:   m.logger,
		storage:  m.storage,
	})
	if err != nil {
		m.logger.Critical(
			log_client.KeyValue{Key: "Service", Value: "MQTT"},
			log_client.KeyValue{Key: "Action", Value: "AddHook"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
	}
	// Запускаем сервер
	go func() {
		err := m.server.Serve()
		if err != nil {
			m.logger.Critical(
				log_client.KeyValue{Key: "Service", Value: "MQTT"},
				log_client.KeyValue{Key: "Action", Value: "Serve"},
				log_client.KeyValue{Key: "Error", Value: err.Error()},
			)
		}
	}()
}
