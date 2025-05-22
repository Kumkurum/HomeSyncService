package test

import (
	"HomeSyncService/internal/storage"
	"testing"
)

var block = storage.NewBlock("testid", 5)

func TestUpdateSensor(t *testing.T) {
	block.UpdateSensor("first", 1, 3.40)
	sensor, err := block.GetSensor("first")
	if err != nil {
		t.Error(err)
	}
	if sensor.Get().BasicData.Value != 3.40 {
		t.Errorf("error in sensor value %v", sensor.Get().BasicData.Value)
	}

	block.UpdateSensor("first", 1, 4.4)
	sensor, err = block.GetSensor("first")
	if err != nil {
		t.Error(err)
	}
	if sensor.Get().BasicData.Value != 4.4 {
		t.Errorf("error in sensor value %v", sensor.Get().BasicData.Value)
	}
	_, err = block.GetSensor("errorTest")

	if err == nil {
		t.Error("error in sensor does not exist")
	}
}

// Может быть ошибка потому как это map, но пофиг на этот тест
func TestGetBlockSensors(t *testing.T) {
	block.UpdateSensor("first", 1, 3.1)
	block.UpdateSensor("first", 1, 3.2)
	block.UpdateSensor("first", 1, 3.3)

	block.UpdateSensor("second", 2, 4.1)
	block.UpdateSensor("second", 2, 4.2)
	block.UpdateSensor("second", 2, 4.3)

	block.UpdateSensor("third", 3, 5.1)
	block.UpdateSensor("third", 3, 5.2)
	block.UpdateSensor("third", 3, 5.3)

	result := block.GetBlockSensors()

	first := false
	second := false
	third := false
	for _, sensor := range result.SensorsData {
		if sensor.Id == "first" {
			first = true
			if sensor.BasicData.Value != 3.3 {
				t.Errorf("error in sensor data, expect  %v, but got %v ", 3.3, sensor.BasicData.Value)
			}
			if sensor.Type != 1 {
				t.Errorf("error in sensor data, expect  %v, but got %v ", 1, sensor.Type)
			}
		}

		if sensor.Id == "second" {
			second = true
			if sensor.BasicData.Value != 4.3 {
				t.Errorf("error in sensor data, expect  %v, but got %v ", 4.3, sensor.BasicData.Value)
			}
			if sensor.Type != 2 {
				t.Errorf("error in sensor data, expect  %v, but got %v ", 2, sensor.Type)
			}
		}

		if sensor.Id == "third" {
			third = true
			if sensor.BasicData.Value != 5.3 {
				t.Errorf("error in sensor data, expect  %v, but got %v ", 3.3, sensor.BasicData.Value)
			}
			if sensor.Type != 3 {
				t.Errorf("error in sensor data, expect  %v, but got %v ", 3, sensor.Type)
			}
		}

	}
	if !(first && second && third) {
		t.Error("error in sensor data")
	}
}
