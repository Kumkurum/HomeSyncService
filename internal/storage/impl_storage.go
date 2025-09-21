package storage

import homeSyncGrpc "HomeSyncService/internal/transport"

type ImplStorage interface {
	UpdateSensorValue(blockId string, sensorId string, typeSensor int, value float32)
	GetSensorsData() *homeSyncGrpc.SensorsResponse
	GetHistoricSensorsData(blockId string, sensorId string) (*homeSyncGrpc.HistorySensorsDataResponse, error)
	SetSensorName(blockId string, sensorId string, sensorName string) *homeSyncGrpc.Error
	SetBoundary(request *homeSyncGrpc.SetBoundaryRequest) *homeSyncGrpc.Error
}
