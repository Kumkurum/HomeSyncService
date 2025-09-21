package main

import (
	"HomeSyncService/internal/grpc_service"
	httpServ "HomeSyncService/internal/http_service"
	mqttServ "HomeSyncService/internal/mqtt_service"
	"HomeSyncService/internal/storage"
	"flag"
	"fmt"
	"github.com/Kumkurum/LogService/pkg/log_client"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var hAddr, gAddr, loggerAddr, token string
	var version, help bool
	flag.StringVar(&hAddr, "h_addr", "50050", "address of http service")
	flag.StringVar(&gAddr, "g_addr", "50051", "address of grpc service")
	flag.StringVar(&loggerAddr, "logger_addr", "/tmp/grpc.sock", "address of grpc logger client")
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

	logger, _ := log_client.NewLoggingClient(loggerAddr, "HomeSync")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan // Ждем сигнал
		fmt.Printf("\nПолучен сигнал: %v\n", sig)
		_ = logger.Info(
			log_client.KeyValue{Key: "Service", Value: "HomeSync"},
			log_client.KeyValue{Key: "Action", Value: "Stop"},
			log_client.KeyValue{Key: "Signal", Value: sig.String()},
		)
		err := logger.Close()
		if err != nil {
			fmt.Println("Error closing logger")
		}
		os.Exit(0)
	}()

	defer func(logger *log_client.LoggingClient) {
		_ = logger.Info(
			log_client.KeyValue{Key: "Service", Value: "HomeSync"},
			log_client.KeyValue{Key: "Action", Value: "Stop"},
		)
		err := logger.Close()
		if err != nil {
			fmt.Println("Error closing logger")
		}
	}(logger)

	var str = storage.NewStorage(10, logger)

	//Запуск HTTP -сервера для связи с ESP
	h := httpServ.NewHttpService(str, logger)
	go h.Run(hAddr)

	mqtt := mqttServ.NewMQTTService(str, logger)
	mqtt.Run("6668", "default", "123")

	// Создать сервер gRPC и зарегистрировать в нем наш KeyValueServer
	grpc_service.NewGrpcService(str, gAddr, token, logger)
}
