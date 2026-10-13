package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddNamedMigrationContext("20260906103017_initdb.go", upCreateUserTable, downCreateUserTable)
	goose.AddNamedMigrationContext("20260906103017_add_account_index.go", upAddAccountIndex, downAddAccountIndex)
	goose.AddNamedMigrationContext("20261013215048_add_balance_to_users", upAddBalanceToUsers, downAddBalanceToUsers)

}

func upCreateUserTable(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    first_name TEXT,
    last_name TEXT,
    middle_name TEXT,
    age INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := tx.Exec(query)

	return err
}

func downCreateUserTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS users;`)
	return err
}
