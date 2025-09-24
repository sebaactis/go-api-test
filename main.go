package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users            = make([]User, 0)
	nextUserID int64 = 1
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// POST /v1/users
func createUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD NOT ALLOWED")
		return
	}

	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
	}

	if strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "Name requerido")
	}

	if !strings.Contains(in.Email, "@") {
		writeError(w, http.StatusBadRequest, "email inválido")
		return
	}

	u := User{ID: nextUserID, Name: in.Name, Email: in.Email}
	users = append(users, u)
	nextUserID++

	writeJSON(w, http.StatusCreated, u)

}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD NOT ALLOWED")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 3 || parts[0] != "v1" || parts[1] != "users" {
		writeError(w, http.StatusNotFound, "Ruta no encontrada")
		return
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	for _, u := range users {
		if u.ID == id {
			writeJSON(w, http.StatusOK, u)
			return
		}
	}

	writeError(w, http.StatusNotFound, "no encontrado")

}

func main() {
	port := 8080

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/v1/users", createUserHandler)

	http.HandleFunc("/v1/users/", getUserByIdHandler)

	log.Printf("Escuchando en el puerto %d", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}

}
