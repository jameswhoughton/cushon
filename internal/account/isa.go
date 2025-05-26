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

// Service to manage ISA accounts
//
// ISAs must adhere to the following rules
// - Only available to UK tax residents over the age of 18
// - The account holder must have a valid NI number
// - The account holder is limited by how much they can deposit each tax year.
// - There are no limits on withdrawals
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
		return Account{}, ErrAccountCreatePermission{"Only UK tax residents can open an ISA"}
	}

	// Ensure the customer is over 18
	cutOff := time.Now().AddDate(-18, 0, 0)
	if customer.DateOfBirth.After(cutOff) {
		return Account{}, ErrAccountCreatePermission{"Only customers who are over the age of 18 can open an ISA"}
	}

	// Ensure the customer's NI number is valid
	if err := s.niValidator(customer.NINumber); err != nil {
		return Account{}, ErrAccountCreatePermission{"Customer NI number could not be verified: " + err.Error()}
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
