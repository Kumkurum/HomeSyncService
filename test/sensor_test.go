package test

import (
	"HomeSyncService/internal/storage"
	grpc "HomeSyncService/internal/transport"
	"testing"
)

var sensor = storage.NewSensor(0, 10)

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
	if sensor.Get().BasicData.Value != 4.0 {
		t.Errorf("Wrong value %f was expected : %f", sensor.Get().BasicData.Value, 4.0)
	}
}

func TestGetProto(t *testing.T) {
	sensor.Clear()
	sensor.AddData(1.0)
	sensor.AddData(2.0)
	sensor.AddData(3.0)
	sensor.AddData(4.0)
	sensor.AddData(5.0)
	r := sensor.GetHistory()
	if len(r.GetSuccess().SensorData) != 5 {
		t.Errorf("Wrong len %d was expected : %d", len(r.GetSuccess().SensorData), 5)
	}
}

func TestChangeBoundary(t *testing.T) {
	sensor.Clear()
	sensor.AddData(1.0)
	var boundary = grpc.Boundary{Value1: 1, Value2: 2, Value3: 3, Value4: 4}
	sensor.UpdateBoundary(&boundary)
	r := sensor.Get()
	if r.Boundary.Value1 != 1 || r.Boundary.Value2 != 2 || r.Boundary.Value3 != 3 || r.Boundary.Value4 != 4 {
		t.Errorf("Wrong baundary1 %f was expected : %d", r.Boundary.Value1, 1)
		t.Errorf("Wrong baundary2 %f was expected : %d", r.Boundary.Value2, 2)
		t.Errorf("Wrong baundary3 %f was expected : %d", r.Boundary.Value3, 3)
		t.Errorf("Wrong baundary4 %f was expected : %d", r.Boundary.Value4, 4)
	}
}

func TestChangeName(t *testing.T) {
	sensor.Clear()
	sensor.AddData(1.0)
	sensor.UpdateName("newName")
	r := sensor.Get()
	if r.Name != "newName" {
		t.Errorf("Wrong Name %s was expected : %s", r.Name, "newName")
	}
}
