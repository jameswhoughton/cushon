package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrExceededISALimit = errors.New("ISA limit will be exceeded by transaction")

type StartOfTaxYear struct {
	Day   int
	Month int
}

type ISAService struct {
	repository     Repository
	annualLimit    int
	startOfTaxYear StartOfTaxYear
	niValidator    func(string) error
}

func NewISAService(repository *Repository, annualLimit int, startOfTaxYear StartOfTaxYear, niValidator func(string) error) *ISAService {
	return &ISAService{
		repository:     *repository,
		annualLimit:    annualLimit,
		startOfTaxYear: startOfTaxYear,
		niValidator:    niValidator,
	}
}

func (s *ISAService) CreateAccount(ctx context.Context, customer Customer) (Account, error) {

	// Ensure the customer is a UK tax resident
	if customer.TaxResidency != "uk" {
		return Account{}, fmt.Errorf("Only UK tax residents can open an ISA")
	}

	// Ensure the customer is over 18
	cutOff := time.Now().AddDate(-18, 0, 0)
	if customer.DateOfBirth.After(cutOff) {
		return Account{}, fmt.Errorf("Only customers who are over the age of 18 can open an ISA")
	}

	// Ensure the customer's NI number is valid
	if err := s.niValidator(customer.NINumber); err != nil {
		return Account{}, fmt.Errorf("Customer NI number could not be verified: %v", err)
	}

	account := Account{
		AccountType: ACCOUNT_TYPE_ISA,
		CustomerId:  customer.Id,
	}

	return createAccount(ctx, s.repository, account)
}

func (s *ISAService) Invest(ctx context.Context, accountId uuid.UUID, investments []Investment) error {
	var totalToInvest int

	for _, investment := range investments {
		totalToInvest += investment.Amount
	}

	startOfTaxYear := time.Date(time.Now().Year(), time.Month(s.startOfTaxYear.Month), s.startOfTaxYear.Day, 0, 0, 0, 0, time.UTC)

	totalInvested, _ := s.repository.GetTotalInvestedToDate(ctx, accountId, startOfTaxYear)

	if totalToInvest+totalInvested > s.annualLimit {
		return ErrExceededISALimit
	}

	err := s.repository.Invest(ctx, accountId, investments)

	if err != nil {
		return fmt.Errorf("Unable to complete investment: %w", err)
	}

	return nil
}

func (s *ISAService) AccountTransactions(ctx context.Context, accountId uuid.UUID, filter TransactionFilter) ([]Transaction, error) {
	return getAccountTransactions(ctx, s.repository, accountId, filter)
}
