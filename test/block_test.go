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

	if result.SensorsData[0].BasicData.Value != 3.3 {
		t.Error("error in sensor data")
	}
	if result.SensorsData[0].Id != "first" {
		t.Error("error in sensor data")
	}
	if result.SensorsData[0].Type != 1 {
		t.Error("error in sensor data")
	}

	if result.SensorsData[1].BasicData.Value != 4.3 {
		t.Error("error in sensor data")
	}
	if result.SensorsData[1].Id != "second" {
		t.Error("error in sensor data")
	}
	if result.SensorsData[1].Type != 2 {
		t.Error("error in sensor data")
	}

	if result.SensorsData[2].BasicData.Value != 5.3 {
		t.Error("error in sensor data")
	}
	if result.SensorsData[2].Id != "third" {
		t.Error("error in sensor data")
	}
	if result.SensorsData[2].Type != 3 {
		t.Error("error in sensor data")
	}
}
