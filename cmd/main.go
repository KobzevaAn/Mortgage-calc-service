package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
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
	Salary   bool `json:"salary"`
	Military bool `json:"military"`
	Base     bool `json:"base"`
}

type Mortgage struct {
	ObjectCost     int     `json:"object_cost"`
	InitialPayment int     `json:"initial_payment"`
	Months         int     `json:"months"`
	Program        Program `json:"program"`
}

type Params struct {
	ObjectCost     int `json:"object_cost"`
	InitialPayment int `json:"initial_payment"`
	Months         int `json:"months"`
}

type Aggregates struct {
	Rate            int    `json:"rate"`
	LoanSum         int    `json:"loan_sum"`
	MonthlyPayment  int    `json:"monthly_payment"`
	Overpayment     int    `json:"overpayment"`
	LastPaymentDate string `json:"last_payment_date"`
}

type CreditData struct {
	Params     Params     `json:"params"`
	Program    Program    `json:"program"`
	Aggregates Aggregates `json:"aggregates"`
}

var ProgramCost = map[string]int{
	"salary":   8,
	"military": 9,
	"base":     10,
}

func Calculate(m Mortgage) CreditData {
	var result CreditData
	result.Program = m.Program
	result.Params = Params{
		ObjectCost:     m.ObjectCost,
		InitialPayment: m.InitialPayment,
		Months:         m.Months,
	}

	switch {
	case m.Program.Salary:
		result.Aggregates.Rate = 8
	case m.Program.Military:
		result.Aggregates.Rate = 9
	case m.Program.Base:
		result.Aggregates.Rate = 10
	}

	result.Aggregates.LoanSum = m.ObjectCost - m.InitialPayment
	ratePerMonth := float64(result.Aggregates.Rate) / 12 / 100
	sum := float64(result.Aggregates.LoanSum)
	payment := sum * (ratePerMonth * math.Pow(1+ratePerMonth, float64(m.Months))) / (math.Pow(1+ratePerMonth, float64(m.Months)) - 1)
	result.Aggregates.MonthlyPayment = int(math.Ceil(payment))
	result.Aggregates.Overpayment = result.Aggregates.MonthlyPayment*m.Months - result.Aggregates.LoanSum
	result.Aggregates.LastPaymentDate = time.Now().AddDate(0, m.Months, 0).Format("02.01.2006")
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
	v := reflect.ValueOf(mortgage.Program)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Bool() {
			count++
		}
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
	}

	results := Calculate(mortgage)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func Cache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	http.HandleFunc("/execute", Execute)
	http.HandleFunc("/cache", Cache)

	fmt.Println("Server started")
	http.ListenAndServe(":8080", nil)
}
