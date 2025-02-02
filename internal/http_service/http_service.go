package httpService

import (
	"HomeSyncService/internal/storage"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
}

func NewHttpService(str storage.ImplStorage) *HttpService {
	service := &HttpService{storage: str}
	return service
}

// Run - запуск http сервера по адресу addr
func (h *HttpService) Run(addr string) {
	fmt.Println("http service listening on " + addr)
	http.HandleFunc("/", h.SetSensorData)
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}

// SetSensorData - хэндлер для установки значений через http get запрос
func (h *HttpService) SetSensorData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SetSensorData")
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "set/json" {
		fmt.Println("header content-type is not set")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var s []Sensor
	err = json.Unmarshal(value, &s)

	for _, v := range s {
		h.storage.UpdateSensorValue(v.BlockId, v.SensorId, v.Type, v.Value)
	}
	if err != nil {
		fmt.Println(err)
	}
}
