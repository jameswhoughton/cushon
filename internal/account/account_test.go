package account_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jameswhoughton/cushon/internal/account"
)

func TestAccountValidation(t *testing.T) {
	type testCase struct {
		name           string
		account        account.Account
		isValid        bool
		expectedErrors []string
	}

	cases := []testCase{
		{
			name:           "CustomerID missing",
			account:        account.Account{AccountType: account.ACCOUNT_TYPE_ISA},
			isValid:        false,
			expectedErrors: []string{"customer_id"},
		},
		{
			name:           "AccountType missing",
			account:        account.Account{CustomerId: uuid.New()},
			isValid:        false,
			expectedErrors: []string{"account_type"},
		},
		{
			name:           "AccountType invalid",
			account:        account.Account{CustomerId: uuid.New(), AccountType: "AAA"},
			isValid:        false,
			expectedErrors: []string{"account_type"},
		},
		{
			name:           "Valid account",
			account:        account.Account{CustomerId: uuid.New(), AccountType: account.ACCOUNT_TYPE_ISA},
			isValid:        true,
			expectedErrors: []string{},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			isValid := testCase.account.Validate()

			if isValid != testCase.isValid {
				t.Errorf("Expected Validate to return %t, got %t", testCase.isValid, isValid)
			}

			if len(testCase.account.Errors) != len(testCase.expectedErrors) {
				t.Errorf("Expected %d validation errors, got %d", len(testCase.expectedErrors), len(testCase.account.Errors))
			}

			for _, field := range testCase.expectedErrors {
				if _, ok := testCase.account.Errors[field]; !ok {
					t.Errorf("expected validation error field '%s' missing", field)
				}
			}
		})
	}
}
