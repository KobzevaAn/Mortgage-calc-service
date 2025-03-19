// Package handlers HTTP-обработчики запросов.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"Mortgage-calc-service/internal/models"
	"Mortgage-calc-service/internal/service"
	"Mortgage-calc-service/internal/storage"
)

// Execute обрабатывает HTTP-запрос на расчет ипотеки.
func Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Failed to close request body: %v", err)
		}
	}()

	var mortgage models.Mortgage
	err := json.NewDecoder(r.Body).Decode(&mortgage)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	count := 0
	if mortgage.Program.Salary {
		count++
	}
	if mortgage.Program.Military {
		count++
	}
	if mortgage.Program.Base {
		count++
	}
	if count != 1 {
		http.Error(w, "Choose only 1 program", http.StatusBadRequest)
		return
	}

	if mortgage.ObjectCost == 0 {
		http.Error(w, "Object cost cannot be zero", http.StatusBadRequest)
		return
	}
	percent := int(float64(mortgage.InitialPayment) / float64(mortgage.ObjectCost) * 100)
	if percent < 20 {
		http.Error(w, "The initial payment should be more", http.StatusBadRequest)
		return
	}

	results := service.Calculate(mortgage)

	storage.AddToCache(results)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Cache получение всех рассчитанных ипотек из кэша.
func Cache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := storage.GetCache()
	if len(data) == 0 {
		http.Error(w, "Cache is empty", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}
