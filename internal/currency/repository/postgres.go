package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
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
	code domain.CurrencyCode,
	rate domain.Rate,
	rateDate time.Time,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO currencies (code, rate, rate_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (code) DO UPDATE SET
		    rate = EXCLUDED.rate,
		    rate_date = EXCLUDED.rate_date
		
	`, code.String(), rate.Float64(), rateDate)

	return err
}

func (r *CurrencyRepoPostgres) GetOne(
	ctx context.Context,
	code domain.CurrencyCode) (domain.Currency, error) {

	var c domain.Currency
	var rate float64
	var codeStr string

	err := r.db.QueryRowContext(ctx, `
		SELECT code, rate, rate_date
		FROM currencies
		WHERE code = $1
	`, code.String()).Scan(&codeStr, &rate, &c.RateDate)

	if err == sql.ErrNoRows {
		return domain.Currency{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.Currency{}, err
	}

	c.Code = domain.CurrencyCode(codeStr)
	c.Rate = domain.Rate(rate)

	return c, nil
}

func (r *CurrencyRepoPostgres) GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT code, rate, rate_date
        FROM currencies
        ORDER BY code
    `)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	result := make(map[domain.CurrencyCode]domain.Currency)

	for rows.Next() {
		var c domain.Currency
		var codeStr string
		var rate float64

		if err := rows.Scan(&codeStr, &rate, &c.RateDate); err != nil {
			return nil, err
		}

		c.Code = domain.CurrencyCode(codeStr)
		c.Rate = domain.Rate(rate)
		result[c.Code] = c
	}

	return result, nil
}

func (r *CurrencyRepoPostgres) Create(
	ctx context.Context,
	code domain.CurrencyCode,
	rate domain.Rate, date time.Time) error {

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO currencies (code, rate) VALUES ($1, $2, $3)`,
		code.String(), rate.Float64(), date,
	)
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}
	}

	return fmt.Errorf("insert currency %s: %w", code, err)
}

func (r *CurrencyRepoPostgres) UpdateOne(
	ctx context.Context,
	code domain.CurrencyCode,
	rate domain.Rate,
	date time.Time) error {

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE currencies SET rate = $1, rate_date = $2 WHERE code = $3`,
		rate.Float64(), date, code.String(),
	)

	if err != nil {
		return fmt.Errorf("update currency %s: %w", code, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update currency %s: rows affected: %w", code, err)
	}
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
