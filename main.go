package main

import (
	"encoding/json"
	"fmt"
)

type DeviceStructRequest struct {
	PaymentURLPhoto string `json:"payment_url_photo"`
	Amount          uint   `json:"amount"`
}

func main() {
	// Притер входящего JSON запроса
	incominRequest := `{"payment_url_photo": "https://example.com/photo.jpg", "amount": 100}`

	// Создаем экземпляр структуры и заполняем его данными из JSON запроса
	var deviceStructRequest DeviceStructRequest

	err := json.Unmarshal([]byte(incominRequest), &deviceStructRequest)
	if err != nil { // Обработка ошибки при разборе JSON
		panic(err)
	}

	// Теперь deviceStructRequest содержит данные из JSON запроса
	fmt.Println(deviceStructRequest.PaymentURLPhoto)
	fmt.Println(deviceStructRequest.Amount)

	// Теперь пакуем структуру обратно в JSON
	jsonData, err := json.Marshal(deviceStructRequest)
	if err != nil { // Обработка ошибки при упаковке в JSON
		panic(err)
	}

	fmt.Println(string(jsonData))
}
