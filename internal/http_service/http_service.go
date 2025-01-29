package http_service

import (
	"HomeSyncService/internal/storage"
	"fmt"
	"io"
	"net/http"
)

type httpService struct {
	storage storage.ImplStorage
}

func NewHttpService(storage storage.ImplStorage) *httpService {

	return &httpService{
		storage: storage,
	}
}

func (h *httpService) SetSensorData(w http.ResponseWriter, r *http.Request) {
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(string(value))
}
