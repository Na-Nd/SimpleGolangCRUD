package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-golang-crud/internal/server"
	"simple-golang-crud/internal/storage"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Загрузка .env

	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	appPort := os.Getenv("APP_PORT")

	// Если значения не подгрузились из env
	if appPort == "" {
		appPort = "8080"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	// DSN для БД
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	// Подключение к БД
	store, err := storage.NewPostgres(dsn)
	if err != nil {
		log.Fatalf("Ошибка соединения с БД: %s", err)
	}
	defer func() {
		err := store.Close()
		if err != nil {
			return
		}
	}()

	// Роутер
	r := server.NewRouter(store)

	// Конфигурация и запуск сервера
	srv := &http.Server{
		Addr:         ":" + appPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("Запуск сервера на порту: %s", appPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
