package postgres

import (
	"context"
	"database/sql"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS currency_latest (
    code TEXT PRIMARY KEY,
    rate NUMERIC(18,6) NOT NULL,
    rate_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS currency_history (
    code TEXT NOT NULL,
    rate NUMERIC(18,6) NOT NULL,
    rate_date DATE NOT NULL,
    PRIMARY KEY (code, rate_date)
);
`
	_, err := db.ExecContext(ctx, query)
	return err
}
