package account

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

const (
	ACCOUNT_TYPE_ISA string = "isa"
)

type Account struct {
	Id          uuid.UUID         `json:"id"`
	CustomerId  uuid.UUID         `json:"customer_id"`
	AccountType string            `json:"account_type"`
	CreatedAt   time.Time         `json:"created_at"`
	Errors      map[string]string `json:"errors"`
}

// Validate a new Account entity
//
// Any errors are stored in a map using the json struct tag
// so that they can be returned straight back to the UI.
func (a *Account) Validate() bool {
	if a.Errors == nil {
		a.Errors = make(map[string]string, 2)
	}

	if a.CustomerId == (uuid.UUID{}) {
		a.Errors["customer_id"] = "Customer ID missing"
	}

	if !slices.Contains([]string{ACCOUNT_TYPE_ISA}, a.AccountType) {
		a.Errors["account_type"] = "Account type invalid or missing"
	}

	return len(a.Errors) == 0
}

var ErrAccountInvalid = errors.New("Account invalid")
