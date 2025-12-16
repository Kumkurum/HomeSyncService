package storage

import homeSyncGrpc "HomeSyncService/internal/transport"

type ImplStorage interface {
	UpdateSensorValue(sensorId string, typeSensor int, value float32)
	GetSensorsData() *homeSyncGrpc.SensorsResponse
	GetHistoricSensorsData(sensorId string) (*homeSyncGrpc.HistorySensorsDataResponse, error)
	SetBoundary(request *homeSyncGrpc.SetBoundaryRequest) *homeSyncGrpc.Error
	RemoveSensor(request *homeSyncGrpc.RemoveSensorRequest) *homeSyncGrpc.Error
}
