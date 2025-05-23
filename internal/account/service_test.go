package account_test

import (
	"testing"
	"time"

	"github.com/jameswhoughton/cushon/internal/account"
)

func TestTransactionFilterValidation(t *testing.T) {
	type testCase struct {
		name           string
		filter         account.TransactionFilter
		isValid        bool
		expectedErrors []string
	}

	cases := []testCase{
		{
			name: "Start date is after end date",
			filter: account.TransactionFilter{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC),
			},
			isValid:        false,
			expectedErrors: []string{"start_date"},
		},
		{
			name: "Difference between start date and end date is greater than one year",
			filter: account.TransactionFilter{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			isValid:        false,
			expectedErrors: []string{"end_date"},
		},
		{
			name: "Valid filter",
			filter: account.TransactionFilter{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			isValid:        true,
			expectedErrors: []string{},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			isValid := testCase.filter.Validate()

			if isValid != testCase.isValid {
				t.Errorf("Expected Validate to return %t, got %t", testCase.isValid, isValid)
			}

			if len(testCase.filter.Errors) != len(testCase.expectedErrors) {
				t.Errorf("Expected %d validation errors, got %d", len(testCase.expectedErrors), len(testCase.filter.Errors))
			}

			for _, field := range testCase.expectedErrors {
				if _, ok := testCase.filter.Errors[field]; !ok {
					t.Errorf("expected validation error field '%s' missing", field)
				}
			}
		})
	}
}
