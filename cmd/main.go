package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

/*
Требуется реализовать 2 эндпоинта:
1. `/execute` - расчет ипотеки (POST).
```json
{
    "object_cost": 5000000,
    "initial_payment": 1000000,
    "months": 240,
    "program": {
        "salary": true
    }
}

{
    "object_cost": 5000000,     // стоимость объекта
    "initial_payment": 1000000, // первоначальный взнос
    "months": 240,              // срок
    "program": {                // блок программы кредита
        "salary": true,         // программа для корпоративных клиентов
        "military": true,       // военная ипотека
        "base": true            // базовая программа
    }
}
``
```
2. `/cache` - получение всех рассчитанных ипотек из кэша (GET).

*/

type Program struct {
	Salary   bool `json:"salary,omitempty"`
	Military bool `json:"military,omitempty"`
	Base     bool `json:"base,omitempty"`
}

type Mortgage struct {
	ObjectCost     int     `json:"object_cost"`
	InitialPayment int     `json:"initial_payment"`
	Months         int     `json:"months"`
	Program        Program `json:"program"`
}

type Aggregates struct {
	Rate            int    `json:"rate"`
	LoanSum         int    `json:"loan_sum"`
	MonthlyPayment  int    `json:"monthly_payment"`
	Overpayment     int    `json:"overpayment"`
	LastPaymentDate string `json:"last_payment_date"`
}

type CreditData struct {
	Params struct {
		ObjectCost     int `json:"object_cost"`
		InitialPayment int `json:"initial_payment"`
		Months         int `json:"months"`
	} `json:"params"`
	Program    Program    `json:"program"`
	Aggregates Aggregates `json:"aggregates"`
}

type CreditCache struct {
	ID         int `json:"id"`
	CreditData `json:",inline"`
}

var (
	cache   []CreditCache
	mtx     sync.RWMutex
	cacheID int
)

func Calculate(m Mortgage) CreditData {
	var result CreditData
	result.Params.ObjectCost = m.ObjectCost
	result.Params.InitialPayment = m.InitialPayment
	result.Params.Months = m.Months

	result.Program = m.Program

	rate := 0
	switch {
	case m.Program.Salary:
		rate = 8
	case m.Program.Military:
		rate = 9
	case m.Program.Base:
		rate = 10
	}
	result.Aggregates.Rate = rate

	result.Aggregates.LoanSum = m.ObjectCost - m.InitialPayment
	ratePerMonth := float64(result.Aggregates.Rate) / 12 / 100
	sum := float64(result.Aggregates.LoanSum)
	payment := sum * (ratePerMonth * math.Pow(1+ratePerMonth, float64(m.Months))) / (math.Pow(1+ratePerMonth, float64(m.Months)) - 1)

	result.Aggregates.MonthlyPayment = int(math.Ceil(payment))
	result.Aggregates.Overpayment = result.Aggregates.MonthlyPayment*m.Months - result.Aggregates.LoanSum
	result.Aggregates.LastPaymentDate = time.Now().AddDate(0, m.Months, 0).Format("2006-01-02")
	return result
}

func Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var mortgage Mortgage
	err := json.NewDecoder(r.Body).Decode(&mortgage)
	if err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	// проверка что выбрана одна программа
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
	if count > 1 {
		http.Error(w, "choose only 1 program", http.StatusBadRequest)
		return
	}
	if count == 0 {
		http.Error(w, "choose program", http.StatusBadRequest)
		return
	}

	//проверка первоначального взноса
	if mortgage.ObjectCost == 0 {
		http.Error(w, "object_cost cannot be zero", http.StatusBadRequest)
		return
	}
	percent := int(float64(mortgage.InitialPayment) / float64(mortgage.ObjectCost) * 100)
	if percent < 20 {
		http.Error(w, "the initial payment should be more", http.StatusBadRequest)
		return
	}

	results := Calculate(mortgage)
	mtx.Lock()
	cache = append(cache, CreditCache{cacheID, results})
	cacheID++
	mtx.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func Cache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	mtx.RLock()
	defer mtx.RUnlock()
	if len(cache) == 0 {
		http.Error(w, "empty cache", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cache)
}

func main() {
	http.HandleFunc("/execute", Execute)
	http.HandleFunc("/cache", Cache)

	fmt.Println("Server started")
	http.ListenAndServe(":8080", nil)
}
