// Package main входная точка приложения ипотечный калькулятор.
package main

import (
	"Mortgage-calc-service/internal/config"
	"Mortgage-calc-service/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"time"

	"Mortgage-calc-service/internal/handlers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	addr := fmt.Sprintf(":%s", cfg.Port)

	mux := http.NewServeMux()

	mux.HandleFunc("/execute", handlers.Execute)
	mux.HandleFunc("/cache", handlers.Cache)
	loggedMux := middleware.LoggerMiddleware(mux)

	log.Println("Server started")
	server := &http.Server{
		Addr:         addr,
		Handler:      loggedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start", err)
	}
}
