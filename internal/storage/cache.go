// Package storage предоставляет методы для работы с кэшем.
package storage

import (
	"Mortgage-calc-service/internal/models"
	"sync"
)

// CreditCache кэш для хранения расчета кредита.
type CreditCache struct {
	models.CreditData `json:",inline"`
	ID                int `json:"id"`
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
	cache = append(cache, CreditCache{data, cacheID})
	cacheID++
}

// GetCache возвращает текущий кэш.
func GetCache() []CreditCache {
	mtx.RLock()
	defer mtx.RUnlock()
	return cache
}
