package test

import (
	"encoding/json"
	"log"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMqtt(t *testing.T) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://ip:port")
	opts.SetClientID("quick-test-client")
	opts.SetUsername("default")
	opts.SetPassword("123")

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Ошибка подключения:", token.Error())
	}
	defer client.Disconnect(250)

	// Тестовые данные
	testData := map[string]interface{}{
		"sensors": []map[string]interface{}{
			{"id": "temperature", "value": 22.54, "type": 0},
			{"id": "pressure", "value": 60.2, "type": 1},
			{"id": "co2", "value": 666.55, "type": 2},
		},
	}

	jsonData, _ := json.Marshal(testData)

	// Отправка
	token := client.Publish("home/main_block", 0, false, jsonData)
	token.Wait()

	if token.Error() != nil {
		log.Fatal("Ошибка отправки:", token.Error())
	}

	log.Println("Тестовые данные успешно отправлены!")
	log.Println("JSON:", string(jsonData))

}
