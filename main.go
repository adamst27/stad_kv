package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	e := NewEngine()

	http.HandleFunc("/set", authMiddleware(handleSet(e)))
	http.HandleFunc("/get", authMiddleware(handleGet(e)))
	http.HandleFunc("/delete", authMiddleware(handleDelete(e)))
	http.HandleFunc("/deleteAll", authMiddleware(handleDeleteAll(e)))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) 
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		authToken := os.Getenv("AUTH_TOKEN")
		if token != authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func handleSet(e *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for k, v := range data {
			if err := e.Set(k, v); err != nil {
				http.Error(w, fmt.Sprintf("Error setting %s: %v", k, err), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleGet(e *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}

		var value interface{}
		if err := e.Get(key, &value); err != nil {
			http.Error(w, fmt.Sprintf("Error getting %s: %v", key, err), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(value)
	}
}

func handleDelete(e *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}

		if err := e.Delete(key); err != nil {
			http.Error(w, fmt.Sprintf("Error deleting %s: %v", key, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleDeleteAll(e *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := e.DeleteAll(); err != nil {
			http.Error(w, fmt.Sprintf("Error deleting all: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
