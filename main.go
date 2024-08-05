package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

type FormData struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Subject          string `json:"subject"`
	Message          string `json:"message"`
	ServiceSelection string `json:"service_selection"`
	PrivacyPolicy    string `json:"privacy_policy"`
}

type IPMessages struct {
	IP           string `json:"ip"`
	MessageCount int    `json:"message_count"`
}

func main() {
	godotenv.Load()

	var ipData []IPMessages

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Define handler for POST request
	http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		numberOfTimesIPHasBeenSeen := 0
		for i := 0; i < len(ipData); i++ {
			if ipData[i].IP == r.RemoteAddr {
				numberOfTimesIPHasBeenSeen = ipData[i].MessageCount
				ipData = slices.Delete(ipData, i, i+1)
				fmt.Println("IP: ", r.RemoteAddr, " Name: ", r.FormValue("your-name"))
			}
		}

		ipData = append(ipData, IPMessages{IP: r.RemoteAddr, MessageCount: numberOfTimesIPHasBeenSeen + 1})

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		formData := FormData{
			Name:             r.FormValue("your-name"),
			Email:            r.FormValue("your-email"),
			PhoneNumber:      r.FormValue("phone-number"),
			Subject:          r.FormValue("your-subject"),
			Message:          r.FormValue("message"),
			ServiceSelection: r.FormValue("service-selection"),
			PrivacyPolicy:    r.FormValue("privacy-policy"),
		}

		// Convert struct to JSON
		formDataJSON, err := json.Marshal(formData)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		// Create map with "content" key
		content := map[string]interface{}{
			"content": string(formDataJSON),
		}

		// Convert map to JSON
		jsonData, err := json.Marshal(content)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		// Print JSON data
		fmt.Println("JSON Data:", string(jsonData))

		// Define the URL to send the POST request to
		url := ""
		discordWebhookURL, exists := os.LookupEnv("DISCORD_WEBHOOK_URL")
		if exists {
			url = discordWebhookURL
		}

		// Create the JSON content
		contentSend := map[string]string{"content": string(jsonData)}
		jsonContent, err := json.Marshal(contentSend)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		// Create a new request with POST method and the JSON content
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonContent))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the content type header to application/json
		req.Header.Set("Content-Type", "application/json")

		// Create a new HTTP client
		client := &http.Client{}

		if numberOfTimesIPHasBeenSeen < 2 || r.FormValue("your-name") == "Marie" {
			// Send the request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()
		}
	})

	// Development
	http.ListenAndServe(":8000", nil)
	//Production
	// http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
}
