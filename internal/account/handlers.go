package account

import "net/http"

// Create an account
func PostAccountHandler(serviceFactory ServiceFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// Invest in a fund

// Get account transactions
