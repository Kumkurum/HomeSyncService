package httpService

import (
	"HomeSyncService/internal/storage"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Kumkurum/LogService/pkg/log_client"
)

type Sensor struct {
	SensorId string  `json:"sId"`
	BlockId  string  `json:"bId"`
	Type     int     `json:"t"`
	Value    float32 `json:"v"`
}
type HttpService struct {
	storage storage.ImplStorage
	server  *http.Server
	logger  *log_client.LoggingClient
}

func NewHttpService(str storage.ImplStorage, logger *log_client.LoggingClient) *HttpService {
	service := &HttpService{storage: str, logger: logger}
	return service
}

// Run - запуск http сервера по адресу addr
func (h *HttpService) Run(addr string) {
	_ = h.logger.Info(
		log_client.KeyValue{Key: "Function", Value: "Run"},
		log_client.KeyValue{Key: "Listen on ", Value: addr},
	)
	http.HandleFunc("/", h.SetSensorData)
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}

// SetSensorData - хэндлер для установки значений через http get запрос
func (h *HttpService) SetSensorData(w http.ResponseWriter, r *http.Request) {
	_ = h.logger.Info(
		log_client.KeyValue{Key: "Function", Value: "SetSensorData"},
	)
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "set/json" {
		_ = h.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "SetSensorData"},
			log_client.KeyValue{Key: "Error", Value: "content type not set"},
		)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		_ = h.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "SetSensorData"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var s []Sensor
	err = json.Unmarshal(value, &s)

	for _, v := range s {
		h.storage.UpdateSensorValue(v.BlockId, v.SensorId, v.Type, v.Value)
	}
	if err != nil {
		_ = h.logger.Warn(
			log_client.KeyValue{Key: "Function", Value: "SetSensorData"},
			log_client.KeyValue{Key: "Error", Value: err.Error()},
		)
	}
}
