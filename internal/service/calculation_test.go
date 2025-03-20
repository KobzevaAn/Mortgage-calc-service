package service

import (
	"Mortgage-calc-service/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// / TestCalculate_Check проверяет корректность расчетов ипотеки.
func TestCalculate_Check(t *testing.T) {
	tests := []struct {
		name                    string
		mortgage                models.Mortgage
		expectedRate            int
		expectedLoanSum         int
		expectedPayment         int
		expectedOverpay         int
		expectedLastPaymentDate string
	}{
		{
			name: "Salary Program",
			mortgage: models.Mortgage{
				ObjectCost:     5000000,
				InitialPayment: 1000000,
				Months:         240,
				Program:        models.Program{Salary: true},
			},
			expectedRate:            8,
			expectedLoanSum:         4000000,
			expectedPayment:         33458,
			expectedOverpay:         4029920,
			expectedLastPaymentDate: "2044-02-18",
		},
		{
			name: "Military Program",
			mortgage: models.Mortgage{
				ObjectCost:     8000000,
				InitialPayment: 2000000,
				Months:         200,
				Program:        models.Program{Military: true},
			},
			expectedRate:            9,
			expectedLoanSum:         6000000,
			expectedPayment:         58019,
			expectedOverpay:         5603800,
			expectedLastPaymentDate: "2040-10-18",
		},
		{
			name: "Base Program",
			mortgage: models.Mortgage{
				ObjectCost:     12000000,
				InitialPayment: 3000000,
				Months:         120,
				Program:        models.Program{Base: true},
			},
			expectedRate:            10,
			expectedLoanSum:         9000000,
			expectedPayment:         118936,
			expectedOverpay:         5272320,
			expectedLastPaymentDate: "2034-02-18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Calculate(tt.mortgage)

			assert.Equal(t, tt.expectedRate, result.Aggregates.Rate)
			assert.Equal(t, tt.expectedLoanSum, result.Aggregates.LoanSum)
			assert.Equal(t, tt.expectedPayment, result.Aggregates.MonthlyPayment)
			assert.Equal(t, tt.expectedOverpay, result.Aggregates.Overpayment)
			expectedLastPaymentDate := time.Now().AddDate(0, tt.mortgage.Months, 0).Format("2006-01-02")
			assert.Equal(t, expectedLastPaymentDate, result.Aggregates.LastPaymentDate)
		})
	}
}
