package storage

import homeSyncGrpc "HomeSyncService/internal/transport"

type ImplStorage interface {
	UpdateSensorValue(blockId string, sensorId string, typeSensor int, value float32)
	GetSensorsData() *homeSyncGrpc.SensorsResponse
}
