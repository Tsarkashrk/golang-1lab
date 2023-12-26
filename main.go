package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestBody struct {
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const expectedMessage = "Привет, сервер! Это JSON-данные из Postman."

func main() {
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestBody RequestBody
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if requestBody.Message != expectedMessage {
			response := Response{
				Status:  "400",
				Message: "Некорректное содержание поля 'message'",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		fmt.Printf("Received message: %s\n", requestBody.Message)

		response := Response{
			Status:  "success",
			Message: "Данные успешно приняты",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
