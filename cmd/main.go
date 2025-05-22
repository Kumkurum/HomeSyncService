package main

import (
	"HomeSyncService/internal/grpc_service"
	httpServ "HomeSyncService/internal/http_service"
	"HomeSyncService/internal/storage"
	"flag"
	"fmt"
)

func main() {
	var hAddr, gAddr, token string
	var version, help bool
	flag.StringVar(&hAddr, "h_addr", "50050", "address of http service")
	flag.StringVar(&gAddr, "g_addr", "50051", "address of grpc service")
	flag.StringVar(&token, "token", "default", "verification token for grpc service")
	flag.BoolVar(&version, "version", false, "Version service")
	flag.BoolVar(&help, "help", false, "Help how to use service")
	flag.Parse()
	if version {
		fmt.Println("Version 0.0.1")
		return
	}
	if help {
		fmt.Println("This is a service for Home Sync Service")
		fmt.Println("flag h_addr to set port for http service, default = 50050 ")
		fmt.Println("flag g_addr to set port for grpc service, default = 50050 ")
		fmt.Println("flag token to set verification token, default = default")
		fmt.Println("just easy to use!")
		return
	}
	var str = storage.NewStorage(10)

	//Запуск HTTP -сервера для связи с ESP
	h := httpServ.NewHttpService(str)
	go h.Run(hAddr)

	// Создать сервер gRPC и зарегистрировать в нем наш KeyValueServer
	grpc_service.NewGrpcService(str, gAddr, token)
}
