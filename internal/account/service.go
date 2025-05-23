package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrTransactionFilterInValid = errors.New("Filter values are not valid")

type TransactionFilter struct {
	StartDate time.Time         `json:"start_date"`
	EndDate   time.Time         `json:"end_date"`
	Errors    map[string]string `json:"errors"`
}

func (f *TransactionFilter) Validate() bool {
	if f.Errors == nil {
		f.Errors = make(map[string]string, 1)
	}

	if f.StartDate.After(f.EndDate) {
		f.Errors["start_date"] = "Start date must come before the End date"
	}

	if f.EndDate.After(f.StartDate.AddDate(1, 0, 0)) {
		f.Errors["end_date"] = "End date cannot be more than a year after start date"
	}

	return len(f.Errors) == 0
}

const (
	TRANSACTION_TYPE_CUSTOMER string = "customer_purchase"
	TRANSACTION_TYPE_FUND     string = "fund_purchase"
)

type Transaction struct {
	Id              int
	FundId          uuid.UUID
	TransactionType string
	Amount          int
}

// Implementations of the Service interface should
type Service interface {
	CreateAccount(ctx context.Context, customerId uuid.UUID) (Account, error)
	Invest(ctx context.Context, accountId uuid.UUID, investments []Investment) error
	AccountTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
}

type ServiceFactory struct {
	isa *ISAService
	// Other account types can be added here
}

func (f *ServiceFactory) Service(accountType string) Service {
	switch accountType {
	case ACCOUNT_TYPE_ISA:
		return f.isa
	default:
		return nil
	}
}

// Generic function to create a new account
//
// This function is designed to be used across all different account types.
// If an account is invalid it will return a ErrorAccountInvalid error
func createAccount(ctx context.Context, repo Repository, account Account) (Account, error) {
	if !account.Validate() {
		return account, ErrorAccountInvalid{}
	}

	err := repo.Create(ctx, &account)

	if err != nil {
		return account, fmt.Errorf("unable to create an account: %v", err)
	}

	return account, nil
}

func getAccountTransactions(ctx context.Context, repo Repository, accountId uuid.UUID, filter TransactionFilter) ([]Transaction, error) {
	if !filter.Validate() {
		return []Transaction{}, ErrTransactionFilterInValid
	}

	return repo.GetAccountTransactions(ctx, filter)
}
