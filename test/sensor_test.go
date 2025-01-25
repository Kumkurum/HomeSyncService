package test

import (
	"HomeSyncService/internal/storage"
	grpc "HomeSyncService/internal/transport"
	"testing"
)

var sensor = storage.NewSensor("12345678", 0, 10)

func TestAddData(t *testing.T) {
	sensor.AddData(1.0)
	sensor.AddData(2.0)
	sensor.AddData(3.0)
	sensor.AddData(4.0)
	sensor.AddData(5.0)
	sensor.AddData(6.0)
	sensor.AddData(7.0)
	sensor.AddData(8.0)
	sensor.AddData(9.0)
	sensor.AddData(10.0)
	sensor.AddData(11.0)
	sensor.AddData(12.0)
	if sensor.Data.Len() > 10 {
		t.Errorf("Sensor List dont remove after overflow maxSize : %d len List : %d", 10, sensor.Data.Len())
	}
	if sensor.Data.Back().Value.(*grpc.BasicSensorData).Value != 12 {
		t.Errorf("Wrong value was expected :%d", 12)
	}
	if sensor.Data.Front().Value.(*grpc.BasicSensorData).Value != 3 {
		t.Errorf("Wrong value was expected :%d", 3)
	}
}

func TestGetData(t *testing.T) {
	sensor.AddData(1.0)
	sensor.AddData(2.0)
	sensor.AddData(3.0)
	sensor.AddData(4.0)
	if sensor.Get().Value != 4.0 {
		t.Errorf("Wrong value %f was expected : %f", sensor.Get().Value, 4.0)
	}
}

func TestGetProto(t *testing.T) {
	sensor.Clear()
	sensor.AddData(1.0)
	sensor.AddData(2.0)
	sensor.AddData(3.0)
	sensor.AddData(4.0)
	sensor.AddData(5.0)
	r := sensor.GetProto()
	if len(r.SensorData) != 5 {
		t.Errorf("Wrong len %d was expected : %d", len(r.SensorData), 5)
	}
}
