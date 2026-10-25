package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitdb, downInitdb)
}

func upInitdb(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  type TEXT NOT NULL CHECK(type IN ('transfer','deposit','withdraw')),
  status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','completed','failed','cancelled')), 
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);`
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS transaction_entries (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    direction TEXT NOT NULL CHECK(direction IN ('DEBIT','CREDIT')),
    amount DECIMAL(20,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (transaction_id) REFERENCES transactions (id) ON DELETE CASCADE
)`
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	query = `
CREATE INDEX IF NOT EXISTS idx_transaction_entries_transaction_id ON transaction_entries(transaction_id);
CREATE INDEX IF NOT EXISTS idx_transaction_entries_account_id ON transaction_entries(account_id);
CREATE INDEX IF NOT EXISTS idx_transaction_entries_direction ON transaction_entries(direction);
`
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	query = `
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
`
	_, err = tx.Exec(query)
	return err

}

func downInitdb(ctx context.Context, tx *sql.Tx) error {
	queries := []string{
		`DROP INDEX idx_transaction_entries_transaction_id ON transaction_entries;`,
		`DROP INDEX idx_transaction_account_id ON transaction_entries;`,
		`DROP INDEX idx_transaction_direction ON transaction_entries;`,
		`DROP INDEX idx_transactions_user_id ON transactions;`,
		`DROP INDEX idx_transactions_status ON transactions;`,
		`DROP INDEX idx_transactions_created_at ON transactions;`,
		`ALTER TABLE transaction_entries DROP PRIMARY KEY;`,
		`ALTER TABLE transactions DROP PRIMARY KEY;`,
	}
	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			return err
		}
	}

	queries = []string{
		`DROP TABLE IF EXISTS transactions;`,
		`DROP TABLE IF EXISTS transaction_entries;`,
	}

	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
