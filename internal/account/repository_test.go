package account_test

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jameswhoughton/cushon/database"
	"github.com/jameswhoughton/cushon/internal/account"
)

// Helper function to connect to the testing database
//
// The database is migrated when this function is called.
// The database implementation of the repository is returned
// along with a deferrable closedown function that rolls back
// the database.
func NewTestRepository() (account.Repository, func()) {
	conn, err := sql.Open("mysql", "root@tcp(127.0.0.1:8002)/retail_accounts")

	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate(conn)

	if err != nil {
		log.Fatal(err)
	}

	closeDown := func() {
		err := database.Rollback(conn)

		if err != nil {
			log.Fatal(err)
		}
	}

	return database.NewAccountRepository(conn), closeDown
}
