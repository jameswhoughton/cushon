package account

import "net/http"

// Handler to create an account
// POST /api/v1/account
//
// I have included this just to show where I would store the handlers
// and how I would structure them.
func PostAccountHandler(serviceFactory ServiceFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the post data.

		// create the correct account instance of the serviceFactory

		// Fetch the user from the session

		// call service.CreateAccount(...)

		// Return 200 on success or 422 on validation error.
	}
}

// Invest in a fund
// POST /api/v1/account/{account id}/invest

// Get account transactions
// GET /api/v1/account/{account id}
