package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"sync"
)

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

var Users []User
var mu sync.Mutex

func getUserIndex(name string) int {
	for i, u := range Users {
		if u.Name == name {
			return i
		}
	}
	return -1
}
func ifUserExist(name string) bool {
	for _, u := range Users {
		if u.Name == name {
			return true
		}
	}
	return false
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("user")
	if name == "" {
		http.Error(w, "User parameter is missing", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	index := -1
	for i, u := range Users {
		if u.Name == name {
			index = i
			break
		}
	}

	if index != -1 {
		Users = append(Users[:index], Users[index+1:]...)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("user")
	email := r.URL.Query().Get("email")
	avatar := r.URL.Query().Get("avatar")

	if name == "" {
		http.Error(w, "User parameter is missing", http.StatusBadRequest)
		return
	}

	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}
	}

	if avatar != "" {
		if _, err := url.ParseRequestURI(avatar); err != nil {
			http.Error(w, "Invalid avatar URL", http.StatusBadRequest)
			return
		}
	}

	mu.Lock()
	defer mu.Unlock()

	if ifUserExist(name) {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	Users = append(Users, User{Name: name, Email: email, Avatar: avatar})
	w.WriteHeader(http.StatusOK)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	response := struct {
		Users []User `json:"users"`
	}{
		Users: Users,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("user")
	email := r.URL.Query().Get("email")
	avatar := r.URL.Query().Get("avatar")

	if name == "" {
		http.Error(w, "User parameter is missing", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	index := getUserIndex(name)
	if index == -1 {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}
		Users[index].Email = email
	}

	if avatar != "" {
		if _, err := url.ParseRequestURI(avatar); err != nil {
			http.Error(w, "Invalid avatar URL", http.StatusBadRequest)
			return
		}
		Users[index].Avatar = avatar
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/remove", removeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/update", updateHandler)

	fmt.Println("Starting server at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
