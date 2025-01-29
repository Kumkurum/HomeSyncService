package main

import (
	"HomeSyncService/internal/grpc_service"
	"HomeSyncService/internal/storage"
	homeSyncGrpc "HomeSyncService/internal/transport"
	"google.golang.org/grpc"
	"log"
	"net"
)

//homeSyncStorage:=homeSyncStorage.NewHomeSyncStorage(10)

func main() {
	var str = storage.NewStorage(10)
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
