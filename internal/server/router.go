package server

import (
	"net/http"
	"simple-golang-crud/internal/handlers"
	"simple-golang-crud/internal/storage"

	"github.com/gorilla/mux"
)

// NewRouter - эндпоинты
func NewRouter(store *storage.Postgres) http.Handler {
	h := &handlers.Handler{Store: store}
	r := mux.NewRouter()

	r.HandleFunc("/users", h.ListUsersHandler).Methods("GET")
	r.HandleFunc("/users", h.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", h.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", h.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", h.DeleteUserHandler).Methods("DELETE")

	// Health-check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return r
}
