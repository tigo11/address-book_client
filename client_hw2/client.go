package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Record представляет запись для отправки на сервер.
type Record struct {
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
}

const serverURL = "http://localhost:8080"

func main() {
	for {
		// Вывод меню выбора операции
		printMenu()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			record := inputRecordData()
			sendRequest("/record/add", record)
		case 2:
			record := getRecordData()
			sendRequest("/records/get", record)
		case 3:
			record := inputRecordData()
			sendRequest("/record/update", record)
		case 4:
			record := getRecordData()
			sendRequest("/record/delete", record)
		default:
			fmt.Println("Неверный выбор операции")
		}
	}
}

// printMenu отображает меню выбора операции.
func printMenu() {
	fmt.Println("Выберите операцию:")
	fmt.Println("1. Добавить запись")
	fmt.Println("2. Получить записи")
	fmt.Println("3. Обновить запись")
	fmt.Println("4. Удалить запись")
}

func inputRecordData() Record {
	record := Record{}
	fmt.Print("Введите имя: ")
	fmt.Scanln(&record.Name)
	fmt.Print("Введите фамилию: ")
	fmt.Scanln(&record.LastName)
	fmt.Print("Введите отчество: ")
	fmt.Scanln(&record.MiddleName)
	fmt.Print("Введите номер телефона: ")
	fmt.Scanln(&record.Phone)
	fmt.Print("Введите адрес: ")
	fmt.Scanln(&record.Address)
	return record
}

func getRecordData() Record {
	record := Record{}
	fmt.Print("Введите номер телефона: ")
	fmt.Scanln(&record.Phone)
	return record
}

func sendRequest(endpoint string, data interface{}) {
	url := serverURL + endpoint

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Статус ответа:", resp.Status)

	// Обработка ответа от сервера
	var response interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Ошибка декодирования ответа:", err)
		return
	}

	// Вывод результата на клиенте
	fmt.Println("Результат:")
	switch responseData := response.(type) {
	case map[string]interface{}:
		printMap(responseData)
	case []interface{}:
		for _, item := range responseData {
			printMap(item.(map[string]interface{}))
		}
	default:
		fmt.Println("Неподдерживаемый тип данных в ответе")
	}
}

func printMap(m map[string]interface{}) {
	fieldOrder := []string{"id", "name", "last_name", "middle_name", "phone", "address"}

	for _, field := range fieldOrder {
		if value, ok := m[field]; ok {
			fmt.Printf("%s: %v\n", field, value)
		}
	}
	fmt.Println()
}
