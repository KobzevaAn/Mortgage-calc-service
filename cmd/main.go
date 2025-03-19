// Package main входная точка приложения ипотечный калькулятор.
package main

import (
	"Mortgage-calc-service/internal/middleware"
	"log"
	"net/http"
	"time"

	"Mortgage-calc-service/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/execute", handlers.Execute)
	mux.HandleFunc("/cache", handlers.Cache)
	loggedMux := middleware.LoggerMiddleware(mux)

	log.Println("Server started")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start", err)
	}
}
