package test

import (
	httpServ "HomeSyncService/internal/http_service"
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/Kumkurum/LogService/pkg/log_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHttp(t *testing.T) {
	logger, _ := log_client.NewLoggingClient("loggerAddr", "HomeSync")
	var str = storage.NewStorage(10, logger)
	service := httpServ.NewHttpService(str, logger)
	go service.Run("8080")

	client := http.Client{}
	testData := []httpServ.Sensor{{
		SensorId: "testSensor",
		BlockId:  "testBlock",
		Value:    1.0,
		Type:     0,
	},
		{
			SensorId: "testSensor2",
			BlockId:  "testBlock",
			Value:    2.0,
			Type:     1,
		},
	}

	data, err := json.Marshal(testData)
	if err != nil {
		t.Error(err)
	}
	req, _ := http.NewRequest("GET", "http://localhost:8080", bytes.NewBuffer(data))
	req.Header.Add("content-type", `set/json`)
	client.Do(req)

	if len(str.GetSensorsData().GetSuccess().GroupsData) != 1 {
		t.Errorf("expect %d groups, got %d", len(testData), len(str.GetSensorsData().GetSuccess().GroupsData))
	}
	if str.GetSensorsData().GetSuccess().GroupsData[0].Id != testData[0].BlockId {
		t.Errorf("error data in BlockId expect %s, got %s", testData[0].BlockId, str.GetSensorsData().GetSuccess().GroupsData[0].Id)
	}
	if str.GetSensorsData().GetSuccess().GroupsData[0].SensorsData[0].Id != testData[0].SensorId {
		t.Errorf("error data in SensorId expect %s, got %s", testData[0].SensorId, str.GetSensorsData().GetSuccess().GroupsData[0].SensorsData[0].Id)
	}
}

func TestHttpRequest(t *testing.T) {
	client := http.Client{}
	testData := []httpServ.Sensor{{
		SensorId: "testSensor",
		BlockId:  "testBlock",
		Value:    1.0,
		Type:     0,
	},
		{
			SensorId: "testSensor2",
			BlockId:  "testBlock",
			Value:    2.0,
			Type:     1,
		},
	}

	data, err := json.Marshal(testData)
	if err != nil {
		t.Error(err)
	}
	req, _ := http.NewRequest("GET", "http://ip:port", bytes.NewBuffer(data))
	req.Header.Add("content-type", `set/json`)
	client.Do(req)
}

func TestGrpc(t *testing.T) {
	conn, err := grpc.NewClient(
		"ip:port", // Адрес сервера
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создаем клиент (используя сгенерированный код)
	client := homeSyncGrpc.NewHomeSyncGrpcServiceClient(conn)

	requestH := &homeSyncGrpc.SensorsRequest{
		Token: "default",
	}
	responseH, err := client.GetSensors(context.Background(), requestH)
	if err != nil {
		t.Errorf("fail to dial: %v", err)
	}
	//responseH.GetSuccess().GroupsData
	response := responseH.GetSuccess().GroupsData[0].SensorsData[0]
	println(response.Name)
}
