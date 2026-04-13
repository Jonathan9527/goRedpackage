package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func Migrate(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			account VARCHAR(64) NOT NULL UNIQUE,
			username VARCHAR(64) NOT NULL,
			balance NUMERIC(12, 2) NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
		INSERT INTO users (account, username, balance)
		SELECT
			'user' || LPAD(n::text, 4, '0'),
			'用户' || LPAD(n::text, 4, '0'),
			100.00
		FROM generate_series(1, 1000) AS n
		ON CONFLICT (account) DO NOTHING;
	`); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	return nil
}
