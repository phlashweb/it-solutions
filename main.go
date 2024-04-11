package main

import (
    "fmt"
    "net/http"
	"encoding/json"
	"bytes"
)

type FormData struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
	Subject           string `json:"subject"`
	Message           string `json:"message"`
	ServiceSelection  string `json:"service_selection"`
	PrivacyPolicy     string `json:"privacy_policy"`
}

func main() {
    // Serve static files
    http.Handle("/", http.FileServer(http.Dir(".")))
    
    // Define handler for POST request
    http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		formData := FormData{
			Name:              r.FormValue("your-name"),
			Email:             r.FormValue("your-email"),
			PhoneNumber:       r.FormValue("phone-number"),
			Subject:           r.FormValue("your-subject"),
			Message:           r.FormValue("message"),
			ServiceSelection:  r.FormValue("service-selection"),
			PrivacyPolicy:     r.FormValue("privacy-policy"),
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

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()
	})

	// Development
	http.ListenAndServe(":8080", nil)
	//Production
	// http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
}
