// Package models содержит используемые структуры данных.
package models

// Program программы кредитования.
type Program struct {
	Salary   bool `json:"salary,omitempty"`
	Military bool `json:"military,omitempty"`
	Base     bool `json:"base,omitempty"`
}

// Mortgage входные параметры кредита.
type Mortgage struct {
	ObjectCost     int     `json:"object_cost"`
	InitialPayment int     `json:"initial_payment"`
	Months         int     `json:"months"`
	Program        Program `json:"program"`
}

// Aggregates рассчитываемые параметры.
type Aggregates struct {
	LastPaymentDate string `json:"last_payment_date"`
	Rate            int    `json:"rate"`
	LoanSum         int    `json:"loan_sum"`
	MonthlyPayment  int    `json:"monthly_payment"`
	Overpayment     int    `json:"overpayment"`
}

// CreditParams параметры для выходных данных.
type CreditParams struct {
	ObjectCost     int `json:"object_cost"`
	InitialPayment int `json:"initial_payment"`
	Months         int `json:"months"`
}

// CreditData рассчитанные данные кредита.
type CreditData struct {
	Aggregates Aggregates   `json:"aggregates"`
	Program    Program      `json:"program"`
	Params     CreditParams `json:"params"`
}
