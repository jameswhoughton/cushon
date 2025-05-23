package account

import (
	"context"

	"github.com/google/uuid"
)

// Representation of an investment into a given fund
//
// A positive amount represents a purchase whereas a negative amount represents a sale.
type Investment struct {
	FundId          uuid.UUID
	AccountFundId   int
	TradeId         uuid.UUID
	TransactionType string
	Amount          int
}

// Responsible for managing retail accounts, the repository is
// only responsible for updating the data store, any params
// passed in are assumed to be valid.
//
// Methods should be accessed through the service layer
type Repository interface {
	// Create a new account, the ID is populated on success
	//
	// Returns error if the account cannot be created
	Create(ctx context.Context, account *Account) error

	// Invests into one or more funds
	//
	// If the account is already invested in the fund, the total invested will be incremented
	// Returns an error if any of the investments fail, if any do fail non of the investments
	// will be processed.
	Invest(ctx context.Context, account *Account, investments []Investment) error

	GetAccountTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
}
