// Package storage предоставляет методы для работы с кэшем.
package storage

import (
	"Mortgage-calc-service/internal/models"
	"sync"
)

// CreditCache кэш для хранения расчета кредита.
//
//nolint:govet // Отключаем govet, чтобы сохранить порядок JSON-полей
type CreditCache struct {
	ID                int `json:"id"`
	models.CreditData `json:",inline"`
}

var (
	cache   []CreditCache
	mtx     sync.RWMutex
	cacheID int
)

// AddToCache добавляет новые данные в кэш.
func AddToCache(data models.CreditData) {
	mtx.Lock()
	defer mtx.Unlock()
	cache = append(cache, CreditCache{cacheID, data})
	cacheID++
}

// GetCache возвращает текущий кэш.
func GetCache() []CreditCache {
	mtx.RLock()
	defer mtx.RUnlock()
	return cache
}

// ClearCache очищает весь кеш.
func ClearCache() {
	mtx.Lock()
	defer mtx.Unlock()
	cache = nil
	cacheID = 0
}
