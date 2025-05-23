package account_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jameswhoughton/cushon/internal/account"
)

func TestISAServiceCreatesAnIsaAccount(t *testing.T) {
	repo, closeDown := NewTestRepository()
	defer closeDown()

	service := account.NewISAService(&repo, 0)

	ctx := context.Background()

	customerId := uuid.New()

	newAccount, err := service.CreateAccount(ctx, customerId)

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

	service := account.NewISAService(&repo, 1000)

	ctx := context.Background()

	customerId := uuid.New()

	newAccount, err := service.CreateAccount(ctx, customerId)

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
		StartDate: time.Now().Add(24 * time.Hour),
		EndDate:   time.Now().Add(-24 * time.Hour),
	}

	transactions, err := service.Transactions(ctx, filter)

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

	service := account.NewISAService(&repo, 90)

	ctx := context.Background()

	customerId := uuid.New()

	newAccount, err := service.CreateAccount(ctx, customerId)

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
