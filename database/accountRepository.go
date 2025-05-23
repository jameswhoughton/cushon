package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jameswhoughton/cushon/internal/account"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(conn *sql.DB) *AccountRepository {
	return &AccountRepository{db: conn}
}

func (r *AccountRepository) Create(ctx context.Context, account *account.Account) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts
		(id, customer_id, account_type)
		VALUES (
			UUID_TO_BIN(?),
			UUID_TO_BIN(?),
			?
		)
	`, account.Id.String(), account.CustomerId.String(), account.AccountType)

	if err != nil {
		return fmt.Errorf("AccountRepository.Create: Unable to create account: %v", err)
	}

	return nil
}

func (r *AccountRepository) Invest(ctx context.Context, account *account.Account, investments []account.Investment) error {
	// Use a transaction to ensure tables are updated atomically
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("AccountRepository.Invest: Unable to start transaction: %v", err)
	}

	for _, investment := range investments {
		// As a minor optimisation we don't need to update the fund balance if the fund is new
		// as we can insert the balance directly.
		var newFund bool

		// If the fund is new, create an entry in account_funds
		if investment.AccountFundId == 0 {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO account_funds
				(account_id, fund_id, balance)
				VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)
			`, account.Id.String(), investment.FundId.String(), investment.Amount)

			if err != nil {
				return fmt.Errorf("AccountRepository.Invest: Unable to create an account fund: %v", err)
			}

			newFund = true
		}

		// Insert a new transaction
		_, err := tx.ExecContext(ctx, `
		INSERT INTO account_transactions
		(account_fund_id, trade_id, transaction_type_id, amount)
		VALUES (?, UUID_TO_BIN(trade_id), SELECT id FROM transaction_types WHERE code = ?), ?)
		`, investment.AccountFundId, investment.TradeId, investment.TransactionType, investment.Amount)

		if err != nil {
			return fmt.Errorf("AccountRepository.Invest: Unable to create an account transaction: %v", err)
		}

		// Update the account_funds table if the fund is not new
		if !newFund {
			_, err = tx.ExecContext(ctx, `
				UPDATE account_funds SET balance = balance + ?
			`, investment.Amount)

			if err != nil {
				return fmt.Errorf("AccountRepository.Invest: Unable to update fund balance: %v", err)
			}
		}
	}

	return nil
}
