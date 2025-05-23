package account

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExceededISALimit = errors.New("ISA limit will be exceeded by transaction")

type ISAService struct {
	repository     Repository
	annualLimit    int
	startOfTaxYear time.Time
}

func NewISAService(repository *Repository, annualLimit int, startOfTaxYear time.Time) *ISAService {
	return &ISAService{
		repository:     *repository,
		annualLimit:    annualLimit,
		startOfTaxYear: startOfTaxYear,
	}
}

func (s *ISAService) CreateAccount(ctx context.Context, customerId uuid.UUID) (Account, error) {
	account := Account{
		AccountType: ACCOUNT_TYPE_ISA,
		CustomerId:  customerId,
	}

	return createAccount(ctx, s.repository, account)
}

func (s *ISAService) Invest(ctx context.Context, accountId uuid.UUID, investments []Investment) error {
	var totalToInvest int

	for _, investment := range investments {
		totalToInvest += investment.Amount
	}

	return nil
}

func (s *ISAService) AccountTransactions(ctx context.Context, accountId uuid.UUID, filter TransactionFilter) ([]Transaction, error) {
	return getAccountTransactions(ctx, s.repository, accountId, filter)
}
