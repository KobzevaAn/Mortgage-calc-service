package storage

import (
	"Mortgage-calc-service/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAddToCache проверяем добавление данных
func TestAddToCache(t *testing.T) {
	ClearCache()

	data := models.CreditData{
		Aggregates: models.Aggregates{Rate: 8},
		Params:     models.CreditParams{ObjectCost: 5000000, InitialPayment: 1000000, Months: 240},
	}

	AddToCache(data)
	cacheData := GetCache()
	assert.Equal(t, 1, len(cacheData))
	assert.Equal(t, data, cacheData[0].CreditData)
	assert.Equal(t, 0, cacheData[0].ID)

	AddToCache(data)
	cacheData = GetCache()
	assert.Equal(t, 2, len(cacheData))
	assert.Equal(t, 1, cacheData[1].ID)
}

// TestClearCache_Success проверяем очистку кэша
func TestClearCache_Success(t *testing.T) {
	AddToCache(models.CreditData{
		Aggregates: models.Aggregates{Rate: 8},
	})
	assert.NotEmpty(t, GetCache())
	ClearCache()
	assert.Empty(t, GetCache())
}
