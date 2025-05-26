package account_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jameswhoughton/cushon/internal/account"
)

func TestCannotCreateAnAccountIfCustomerDoesNotReachRequirements(t *testing.T) {
	repo, closeDown := NewTestRepository()
	defer closeDown()

	type testCase struct {
		name     string
		customer account.Customer
	}

	testCases := []testCase{
		{
			name: "Non-uk tax resident",
			customer: account.Customer{
				Id:           uuid.New(),
				TaxResidency: "fr",
				DateOfBirth:  time.Now().AddDate(-19, 0, 0),
				NINumber:     "SD000000B",
			},
		},
		{
			name: "Customer under 18",
			customer: account.Customer{
				Id:           uuid.New(),
				TaxResidency: "uk",
				DateOfBirth:  time.Now().AddDate(-17, 0, 0),
				NINumber:     "SD000000B",
			},
		},
		{
			name: "NI Number is invalid",
			customer: account.Customer{
				Id:           uuid.New(),
				TaxResidency: "uk",
				DateOfBirth:  time.Now().AddDate(-19, 0, 0),
				NINumber:     "INVALID",
			},
		},
	}

	niValidator := func(ni string) error {
		if ni != "SD000000B" {
			return fmt.Errorf("NI '%s' is invalid", ni)
		}

		return nil
	}

	service := account.NewISAService(&repo, 0, account.StartOfTaxYear{1, 1}, niValidator)

	ctx := context.Background()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := service.CreateAccount(ctx, testCase.customer)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !errors.As(err, &account.ErrAccountCreatePermission{}) {
				t.Errorf("Expected error of type %T, got %T: %v", account.ErrAccountCreatePermission{}, err, err)
			}
		})
	}
}

func TestISAServiceCreatesAnIsaAccount(t *testing.T) {
	repo, closeDown := NewTestRepository()
	defer closeDown()

	passingNiValidator := func(_ string) error {
		return nil
	}

	service := account.NewISAService(&repo, 0, account.StartOfTaxYear{1, 1}, passingNiValidator)

	ctx := context.Background()

	customer := account.Customer{
		Id:           uuid.New(),
		TaxResidency: "uk",
		NINumber:     "SD000000A",
		DateOfBirth:  time.Now().AddDate(-20, 0, 0),
	}

	newAccount, err := service.CreateAccount(ctx, customer)

	if err != nil {
		t.Errorf("unexpected error when creating ISA account: %v", err)
	}

	if newAccount.AccountType != account.ACCOUNT_TYPE_ISA {
		t.Errorf("expected new account to have the type %s, got %s", account.ACCOUNT_TYPE_ISA, newAccount.AccountType)
	}
}

func TestISAServiceCanInvestMoneyInAFund(t *testing.T) {
	repo, closeDown := NewTestRepository()
	defer closeDown()

	passingNiValidator := func(_ string) error {
		return nil
	}

	service := account.NewISAService(&repo, 200, account.StartOfTaxYear{1, 1}, passingNiValidator)

	ctx := context.Background()

	customer := account.Customer{
		Id:           uuid.New(),
		TaxResidency: "uk",
		NINumber:     "SD000000A",
		DateOfBirth:  time.Now().AddDate(-20, 0, 0),
	}

	newAccount, err := service.CreateAccount(ctx, customer)

	if err != nil {
		t.Errorf("unexpected error when creating ISA account: %v", err)
	}

	investments := []account.Investment{
		{
			FundId:          uuid.New(),
			TradeId:         uuid.New(),
			TransactionType: account.TRANSACTION_TYPE_CUSTOMER,
			Amount:          100,
		},
	}

	err = service.Invest(ctx, newAccount.Id, investments)

	if err != nil {
		t.Errorf("unexpected error when investing in fund: %v", err)
	}

	filter := account.TransactionFilter{
		StartDate: time.Now().Add(-24 * time.Hour),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	transactions, err := service.AccountTransactions(ctx, newAccount.Id, filter)

	if err != nil {
		t.Errorf("unexpected error fetching transactions: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, found %d", len(transactions))
	}
}

func TestICannotInvestInAFundIfIHaveReachedMyAnnualLimit(t *testing.T) {
	repo, closeDown := NewTestRepository()
	defer closeDown()

	passingNiValidator := func(_ string) error {
		return nil
	}

	service := account.NewISAService(&repo, 50, account.StartOfTaxYear{1, 1}, passingNiValidator)

	ctx := context.Background()

	customer := account.Customer{
		Id:           uuid.New(),
		TaxResidency: "uk",
		NINumber:     "SD000000A",
		DateOfBirth:  time.Now().AddDate(-20, 0, 0),
	}

	newAccount, err := service.CreateAccount(ctx, customer)

	if err != nil {
		t.Errorf("unexpected error when creating ISA account: %v", err)
	}

	investments := []account.Investment{
		{
			FundId:          uuid.New(),
			TradeId:         uuid.New(),
			TransactionType: account.TRANSACTION_TYPE_CUSTOMER,
			Amount:          100,
		},
	}

	err = service.Invest(ctx, newAccount.Id, investments)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if !errors.Is(err, account.ErrExceededISALimit) {
		t.Errorf("Expected error %v, got %T - %v", account.ErrExceededISALimit, err, err)
	}
}
