package test

import (
	homeSyncGrpc "HomeSyncService/internal/transport"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Boundary: &homeSyncGrpc.Boundary{
			Value1: 1,
			Value2: 2,
			Value3: 3,
			Value4: 4,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
		Boundary: &homeSyncGrpc.Boundary{
			Value1: 1,
			Value2: 2,
			Value3: 3,
			Value4: 4,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
		Boundary: &homeSyncGrpc.Boundary{
			Value1: 1,
			Value2: 2,
			Value3: 3,
			Value4: 4,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
		Boundary: &homeSyncGrpc.Boundary{
			Value1: 1,
			Value2: 2,
			Value3: 3,
			Value4: 4,
		},
	},
	homeSyncGrpc.SensorData{
		Id:   "testSensor0",
		Type: 0,
		BasicData: &homeSyncGrpc.BasicSensorData{
			Time:  timestamppb.New(time.Now()),
			Value: 1,
		},
		Boundary: &homeSyncGrpc.Boundary{
			Value1: 1,
			Value2: 2,
			Value3: 3,
			Value4: 4,
		},
	},
}

//func runServer() {
//var str = storage.NewStorage(10)
//for _, data := range testData {
//	str.UpdateSensorValue("testBlock", data.Id, int(data.Type), data.BasicData.Value)
//}
//
//// Создать сервер gRPC и зарегистрировать в нем наш KeyValueServer
//grpc_service.NewGrpcService(str, "50051", "kum")
//}

func TestRequest(t *testing.T) {
	//go runServer()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("127.0.0.1:50051", opts...)

	if err != nil {
		t.Errorf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := homeSyncGrpc.NewHomeSyncGrpcServiceClient(conn)
	request := &homeSyncGrpc.SensorsRequest{
		Token: "default",
	}
	response, err := client.GetSensors(context.Background(), request)
	fmt.Println(response)

}
