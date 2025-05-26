package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (r *AccountRepository) Invest(ctx context.Context, accountId uuid.UUID, investments []account.Investment) error {
	// Use a transaction to ensure tables are updated atomically
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("AccountRepository.Invest: Unable to start transaction: %v", err)
	}

	defer tx.Rollback()

	for _, investment := range investments {
		// As a minor optimisation we don't need to update the fund balance if the fund is new
		// as we can insert the balance directly.
		var newFund bool

		// If the fund is new, create an entry in account_funds
		if investment.AccountFundId == 0 {
			result, err := tx.ExecContext(ctx, `
				INSERT INTO account_funds
				(account_id, fund_id, balance)
				VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)
			`, accountId, investment.FundId.String(), investment.Amount)

			if err != nil {
				return fmt.Errorf("AccountRepository.Invest: Unable to create an account fund: %v", err)
			}

			investment.AccountFundId, err = result.LastInsertId()

			if err != nil {
				return fmt.Errorf("AccountRepository.Invest: Unable to fetch new account_funds Id: %v", err)
			}

			newFund = true
		}

		// Insert a new transaction
		_, err := tx.ExecContext(ctx, `
		INSERT INTO account_transactions
		(account_fund_id, trade_id, transaction_type, amount)
		VALUES (?, UUID_TO_BIN(?), ?, ?)
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

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("AccountRepository.Invest: Unable to commit transaction: %v", err)
	}

	return nil
}

func (r *AccountRepository) GetAccountTransactions(ctx context.Context, accountId uuid.UUID, filter account.TransactionFilter) ([]account.Transaction, error) {
	var transactions []account.Transaction

	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, BIN_TO_UUID(f.fund_id), t.transaction_type, t.amount
		FROM accounts a
		LEFT JOIN account_funds f
		ON a.id = f.account_id
		LEFT JOIN account_transactions t
		ON f.id = t.account_fund_id
		WHERE a.id = UUID_TO_BIN(?)
		AND t.created_at >= ?
		AND t.created_at <= ?
	`, accountId, filter.StartDate, filter.EndDate)

	if err != nil {
		return []account.Transaction{}, fmt.Errorf("AccountRepository.GetAccountTransactions: Unable to fetch transactions: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var transaction account.Transaction

		err := rows.Scan(&transaction.Id, &transaction.FundId, &transaction.TransactionType, &transaction.Amount)

		if err != nil {
			return []account.Transaction{}, fmt.Errorf("AccountRepository.GetAccountTransactions: Unable to fetch transactions: %v", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *AccountRepository) GetTotalInvestedToDate(ctx context.Context, accountId uuid.UUID, fromDate time.Time) (int, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT SUM(t.amount) AS total
		FROM accounts a
		LEFT JOIN account_funds f
		ON a.id = f.account_id
		LEFT JOIN account_transactions  t
		ON f.id = t.account_fund_id
		WHERE a.id = ?
		AND transaction_type = 'customer'
		AND t.amount > 0
	`, accountId)

	var total int

	row.Scan(&total)

	return total, nil
}
