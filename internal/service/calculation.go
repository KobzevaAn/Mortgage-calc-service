// Package service содержит бизнес-логику расчета параметров кредита.
package service

import (
	"math"
	"time"

	"Mortgage-calc-service/internal/models"
)

// Calculate вычисляет параметры кредита.
func Calculate(m models.Mortgage) models.CreditData {
	var result models.CreditData

	result.Params.ObjectCost = m.ObjectCost
	result.Params.InitialPayment = m.InitialPayment
	result.Params.Months = m.Months

	result.Program = m.Program

	rate := 0
	if m.Program.Salary {
		rate = 8
	}
	if m.Program.Military {
		rate = 9
	}
	if m.Program.Base {
		rate = 10
	}
	result.Aggregates.Rate = rate

	result.Aggregates.LoanSum = m.ObjectCost - m.InitialPayment
	ratePerMonth := float64(rate) / 12 / 100
	sum := float64(result.Aggregates.LoanSum)
	payment := sum * (ratePerMonth * math.Pow(1+ratePerMonth, float64(m.Months))) / (math.Pow(1+ratePerMonth, float64(m.Months)) - 1)

	result.Aggregates.MonthlyPayment = int(math.Ceil(payment))
	result.Aggregates.Overpayment = result.Aggregates.MonthlyPayment*m.Months - result.Aggregates.LoanSum
	result.Aggregates.LastPaymentDate = time.Now().AddDate(0, m.Months, 0).Format("2006-01-02")

	return result
}
