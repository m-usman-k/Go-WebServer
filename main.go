package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Struct to decode JSON from client
type RequestData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Struct to send JSON back to client
type ResponseData struct {
	Message string `json:"message"`
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "abc123",
		Expires: time.Now().Add(24 * time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name: "Testing",
		Value: "123345",
	})
	w.Write([]byte("Cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No cookie found", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Got cookie: " + cookie.Value))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := ResponseData{
		Message: fmt.Sprintf("Hello %s, your email is %s", data.Name, data.Email),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/set-cookie", setCookieHandler)
	http.HandleFunc("/get-cookie", getCookieHandler)
	http.HandleFunc("/json", jsonHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
