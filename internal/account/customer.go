package account

import (
	"time"

	"github.com/google/uuid"
)

// Retail account service representation of a Customer
//
// For brevity I have only included fields required to verify
// a customer for an ISA account.
type Customer struct {
	Id           uuid.UUID
	NINumber     string
	TaxResidency string
	DateOfBirth  time.Time
}

// Ensure NI number is Valid
//
// This a a dummy function to make an external call to an external
// service to verify a National insurance number.
// Returns an error if the number is not valid.
func ValidateNINumber(ni string) error {

	// make external call to NI Validation service

	return nil
}
