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
	if sensor.Get().Value != 3.40 {
		t.Errorf("error in sensor value %v", sensor.Get().Value)
	}

	block.UpdateSensor("first", 1, 4.4)
	sensor, err = block.GetSensor("first")
	if err != nil {
		t.Error(err)
	}
	if sensor.Get().Value != 4.4 {
		t.Errorf("error in sensor value %v", sensor.Get().Value)
	}
	_, err = block.GetSensor("errorTest")

	if err == nil {
		t.Error("error in sensor does not exist")
	}
}
