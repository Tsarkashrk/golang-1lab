package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const expectedMessage = "Привет, сервер! Это JSON-данные из Postman."

const htmlPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Page</title>
</head>
<body>
    <h1>Тестовая страница</h1>
    <form id="myForm">
        <label for="message">Message:</label>
        <input type="text" id="message" name="message" value="Привет, сервер! Это JSON-данные из Postman.">
        <button type="button" onclick="sendData()">Отправить</button>
    </form>

    <script>
        function sendData() {
            var message = document.getElementById("message").value;

            var data = {
                "message": message
            };

            fetch("/", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            })
            .then(response => response.json())
            .then(data => {
                console.log("Response:", data);
            })
            .catch(error => {
                console.error("Error:", error);
            });
        }
    </script>
</body>
</html>
`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, htmlPage)
		case http.MethodPost:
			var requestBody struct{ Message string `json:"message"` }
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if requestBody.Message != expectedMessage {
				response := struct {
					Status  string `json:"status"`
					Message string `json:"message"`
				}{
					Status:  "400",
					Message: "Некорректное содержание поля 'message'",
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			fmt.Printf("Received message: %s\n", requestBody.Message)

			response := struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Status:  "success",
				Message: "Данные успешно приняты",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is listening on port 8081...")
	http.ListenAndServe(":8081", nil)
}
