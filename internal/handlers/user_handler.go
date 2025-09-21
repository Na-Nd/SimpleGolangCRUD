package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"simple-golang-crud/internal/storage"
	"simple-golang-crud/pkg/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Handler - стурктура, в которую внедряем зависимости
type Handler struct {
	Store *storage.Postgres
}

// writeJSON - вспомогательная функиця для обертки и отправки ответов
func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)              // Статус
	_ = json.NewEncoder(w).Encode(v) // Сериализация
}

// writeError - вспомогательная функция для обертки и отправки ошибок
func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// CreateUserHandler - создание пользователя
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Невалидное тело запроса")
		return
	}

	// Валидация
	if input.Name == "" || input.Email == "" {
		writeError(w, http.StatusBadRequest, "Имя и почта не должны быть пустые")
		return
	}

	// Контекст с таймаутом для DB операции
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Создание пользователя
	if err := h.Store.CreateUser(ctx, &input); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при создании пользователя: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, input)
}

// GetUserHandler - получение пользователя по id
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                        // Достаем параметры из url
	idStr := vars["id"]                        // Достаем именно параметр по ключу id
	id, err := strconv.ParseInt(idStr, 10, 64) // Строку в число

	if err != nil {
		writeError(w, http.StatusBadRequest, "Невалидный идентификатор")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	u, err := h.Store.GetUser(ctx, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	writeJSON(w, http.StatusOK, u)
}

// ListUsersHandler — получение списка всех пользователей
func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := h.Store.GetAllUsers(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении списка пользователей: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// UpdateUserHandler — обновление пользователя по id
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Невалидный идентификатор")
		return
	}

	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Невалидное тело запроса")
		return
	}

	input.ID = id

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.Store.UpdateUser(ctx, &input); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при обновлении пользователя: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, input)
}

// DeleteUserHandler — удаление пользователя по id
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Невалидный идентификатор")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.Store.DeleteUser(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при удалении пользователя: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
