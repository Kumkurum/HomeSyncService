package test

import (
	"HomeSyncService/internal/grpc_service"
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"testing"
	"time"
)

var testData = []homeSyncGrpc.SensorData{
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
	},
}

func runServer() {
	var str = storage.NewStorage(10)
	for _, data := range testData {
		str.UpdateSensorValue("testBlock", data.Id, int(data.Type), data.BasicData.Value)
	}

	// Создать сервер gRPC и зарегистрировать в нем наш KeyValueServer
	s := grpc.NewServer()
	homeSyncGrpc.RegisterHomeSyncGrpcServiceServer(s, grpc_service.NewGrpcService(str))
	// Открыть порт 50051 для приема сообщений
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Начать цикл приема и обработку запросов
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestRequest(t *testing.T) {
	go runServer()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("127.0.0.1:50051", opts...)

	if err != nil {
		t.Errorf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := homeSyncGrpc.NewHomeSyncGrpcServiceClient(conn)
	request := &homeSyncGrpc.SensorsRequest{}
	response, err := client.GetSensors(context.Background(), request)

	if err != nil {
		t.Errorf("fail to dial: %v", err)
	}
	if len(response.GroupsData) != 1 {
		t.Errorf("expect response groups data, but got %v", response.GroupsData)
	}

	if response.GroupsData[0].SensorsData[0].Id != testData[0].Id {
		t.Errorf("expect response groups data, but got %v", response.GroupsData)
	}
	if response.GroupsData[0].SensorsData[0].BasicData.Value != testData[len(testData)-1].BasicData.Value {
		t.Errorf("expect response groups data, but got %v", response.GroupsData)
	}

	requestH := &homeSyncGrpc.HistorySensorDataRequest{
		BlockId:  "testBlock",
		SensorId: "testSensor0",
	}
	responseH, err := client.GetHistorySensorData(context.Background(), requestH)
	if err != nil {
		t.Errorf("fail to dial: %v", err)
	}
	if len(responseH.SensorData) != len(testData) {
		t.Errorf("error len %d, but reale len is %d", len(responseH.SensorData), len(testData))
	}

	if responseH.SensorData[0].Value != testData[0].BasicData.Value {
		t.Errorf("expect response data %v, but got %v", responseH.SensorData[0].Value, testData[0].BasicData.Value)
	}

}
