package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var Users []string
var mu sync.Mutex

func ifUserExist(user string) bool {
	for _, u := range Users {
		if u == user {
			return true
		}
	}
	return false
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}
func addHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "User parameter is missing", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if !ifUserExist(user) {
		Users = append(Users, user)
	}

	response := struct {
		Users []string `json:"users"`
	}{
		Users: Users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/add", addHandler)

	fmt.Println("Starting server at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
