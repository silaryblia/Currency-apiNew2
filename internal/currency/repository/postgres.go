package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type CurrencyRepoPostgres struct {
	mu     sync.RWMutex
	db     *sql.DB
	logger *zap.Logger
}

func NewCurrencyRepoPostgres(db *sql.DB, logger *zap.Logger) *CurrencyRepoPostgres {
	return &CurrencyRepoPostgres{db: db, logger: logger}
}

func (r *CurrencyRepoPostgres) Upsert(
	ctx context.Context,
	code string,
	rate float64,
	rateDate time.Time,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO currencies (code, rate, rate_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (code) DO UPDATE SET
		    rate = EXCLUDED.rate,
		    rate_date = EXCLUDED.rate_date
		
	`, code, rate, rateDate)

	return err
}

func (r *CurrencyRepoPostgres) GetOne(
	ctx context.Context,
	code string) (domain.Currency, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	var c domain.Currency

	err := r.db.QueryRowContext(ctx, `
		SELECT code, rate, rate_date
		FROM currencies
		WHERE code = $1
	`, code).Scan(&c.Code, &c.Rate, &c.RateDate)

	if err == sql.ErrNoRows {

		return domain.Currency{}, domain.ErrNotFound
	}

	return c, err
}

func (r *CurrencyRepoPostgres) GetAll(ctx context.Context) (map[string]domain.Currency, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT code, rate, rate_date
		FROM currencies
		ORDER BY code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]domain.Currency)

	for rows.Next() {
		var c domain.Currency
		if err := rows.Scan(&c.Code, &c.Rate, &c.RateDate); err != nil {
			return nil, err
		}
		result[c.Code] = c
	}

	return result, nil
}

func (r *CurrencyRepoPostgres) Create(ctx context.Context, code string, rate float64, date time.Time) error {
	code = strings.ToUpper(strings.TrimSpace(code))

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO currencies (code, rate) VALUES ($1, $2, $3)`,
		code,
		rate,
		date,
	)

	if err != nil {
		return domain.ErrAlreadyExists
	}

	return nil
}

func (r *CurrencyRepoPostgres) UpdateOne(ctx context.Context, code string, rate float64, date time.Time) error {
	code = strings.ToUpper(strings.TrimSpace(code))

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE currencies SET rate = $1, rate_date = $2 WHERE code = $3`,
		rate, date, code,
	)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *CurrencyRepoPostgres) UpdateAll(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE currencies SET rate_date = $1`,
		time.Now(),
	)

	return err
}

func (r *CurrencyRepoPostgres) DeleteAll(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM currencies`)
	return err
}
