package handlers

import (
	"Mortgage-calc-service/internal/models"
	"Mortgage-calc-service/internal/storage"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExecute_MethodNotAllowed проверка невалидного метода
func TestExecute_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/execute", nil)
	assert.Equal(t, err, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Execute)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Contains(t, rr.Body.String(), "Method not allowed")
}

// TestExecute_InvalidJSON проверка невалидного входного json
func TestExecute_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer([]byte(`{invalid json`)))
	assert.Equal(t, err, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Execute)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request format")
}

// TestExecute_MultipleProgramSelected выбрано больше 1 программы
func TestExecute_MultipleProgramSelected(t *testing.T) {
	requestData := models.Mortgage{
		ObjectCost:     5000000,
		InitialPayment: 1000000,
		Months:         240,
		Program:        models.Program{Salary: true, Military: true},
	}
	body, err := json.Marshal(requestData)
	assert.Equal(t, err, nil)

	req, err := http.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	assert.Equal(t, err, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Execute)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Choose only 1 program")
}

// TestExecute_ZeroObjectCost проверка ненулевой стоимости объекта
func TestExecute_ZeroObjectCost(t *testing.T) {
	requestData := models.Mortgage{
		ObjectCost:     0,
		InitialPayment: 1000000,
		Months:         240,
		Program:        models.Program{Salary: true},
	}
	body, err := json.Marshal(requestData)
	assert.Equal(t, err, nil)

	req, err := http.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	assert.Equal(t, err, nil)
	req.Header.Set("Cotent-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Execute)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Object cost cannot be zero")
}

// TestExecute_InitialPaymentTooLow проверка первоначального взноса
func TestExecute_InitialPaymentTooLow(t *testing.T) {
	requestData := models.Mortgage{
		ObjectCost:     5000000,
		InitialPayment: 10000,
		Months:         240,
		Program:        models.Program{Salary: true},
	}
	body, _ := json.Marshal(requestData)

	req, _ := http.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Execute)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "The initial payment should be more")
}

// TestExecute_Success проверяем успешный кейс
func TestExecute_Success(t *testing.T) {
	storage.ClearCache()

	requestData := models.Mortgage{
		ObjectCost:     5000000,
		InitialPayment: 1000000,
		Months:         240,
		Program:        models.Program{Salary: true},
	}
	body, err := json.Marshal(requestData)
	assert.Equal(t, err, nil)

	req, err := http.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	assert.Equal(t, err, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Execute)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.CreditData
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, err, nil)

	expectedResponse := models.CreditData{
		Aggregates: models.Aggregates{
			Rate:            8,
			LoanSum:         4000000,
			MonthlyPayment:  33458,
			Overpayment:     4029920,
			LastPaymentDate: "2045-03-20",
		},
		Program: models.Program{Salary: true},
		Params:  models.CreditParams{ObjectCost: 5000000, InitialPayment: 1000000, Months: 240},
	}
	assert.Equal(t, expectedResponse, response)

	cacheData := storage.GetCache()
	assert.Equal(t, len(cacheData), 1)
}

// TestCache_MethodNotAllowed проверка невалидного метода
func TestCache_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/cache", nil)
	assert.Equal(t, err, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Cache)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Contains(t, rr.Body.String(), "Method not allowed")
}

// TestCache_EmptyCache проверка для пустого кэша
func TestCache_EmptyCache(t *testing.T) {
	storage.ClearCache()

	req, err := http.NewRequest(http.MethodGet, "/cache", nil)
	assert.Equal(t, err, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Cache)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Cache is empty")
}

// TestCache_Success проверяем успешный кейс
func TestCache_Success(t *testing.T) {
	storage.ClearCache()
	storage.AddToCache(
		models.CreditData{
			Aggregates: models.Aggregates{
				Rate:            8,
				LoanSum:         4000000,
				MonthlyPayment:  33458,
				Overpayment:     4029920,
				LastPaymentDate: "2044-02-18",
			},
			Program: models.Program{Salary: true},
			Params: models.CreditParams{
				ObjectCost:     5000000,
				InitialPayment: 1000000,
				Months:         240,
			},
		},
	)
	req, err := http.NewRequest(http.MethodGet, "/cache", nil)
	assert.Equal(t, err, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Cache)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.CreditData
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, response)

}
