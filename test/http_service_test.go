package test

import (
	httpServ "HomeSyncService/internal/http_service"
	"HomeSyncService/internal/storage"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	var str = storage.NewStorage(10)
	service := httpServ.NewHttpService(str)
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

	if len(str.GetSensorsData().GroupsData) != 1 {
		t.Errorf("expect %d groups, got %d", len(testData), len(str.GetSensorsData().GroupsData))
	}
	if str.GetSensorsData().GroupsData[0].Id != testData[0].BlockId {
		t.Errorf("error data in BlockId expect %s, got %s", testData[0].BlockId, str.GetSensorsData().GroupsData[0].Id)
	}
	if str.GetSensorsData().GroupsData[0].SensorsData[0].Id != testData[0].SensorId {
		t.Errorf("error data in SensorId expect %s, got %s", testData[0].SensorId, str.GetSensorsData().GroupsData[0].SensorsData[0].Id)
	}
}
