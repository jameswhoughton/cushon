CREATE TABLE account_transactions (
	id INT NOT NULL AUTO_INCREMENT,
	account_fund_id INT NOT NULL,
	trade_id BINARY(16) NOT NULL, -- Assumed to be UUID
	transaction_type VARCHAR(25) NOT NULL,
	amount INT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	FOREIGN KEY (account_fund_id)
		REFERENCES account_funds(id)
);
