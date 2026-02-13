package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type CurrencyRepoPostgres struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewCurrencyRepoPostgres(db *sql.DB, logger *zap.Logger) *CurrencyRepoPostgres {
	return &CurrencyRepoPostgres{db: db, logger: logger}
}

func (r *CurrencyRepoPostgres) GetAll(
	ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT code, rate, rate_date
        FROM currency_latest
        ORDER BY code
    `)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Warn("failed to close rows", zap.Error(err))
		}
	}()
	//	defer rows.Close()

	res := make(map[domain.CurrencyCode]domain.Currency)

	for rows.Next() {
		var c domain.Currency
		var rateStr string
		var rateDate time.Time

		if err := rows.Scan(&c.Code, &rateStr, &rateDate); err != nil {
			return nil, err
		}

		// Конвертируем строку в Rate
		rate, err := domain.RateFromString(rateStr)
		if err != nil {
			r.logger.Error("failed to parse rate",
				zap.String("code", string(c.Code)),
				zap.String("rate", rateStr),
				zap.Error(err))
			continue
		}

		c.Rate = rate
		c.RateDate = rateDate
		res[c.Code] = c
	}

	return res, nil
}

func (r *CurrencyRepoPostgres) SaveRate(
	ctx context.Context,
	code domain.CurrencyCode,
	rate domain.Rate,
	date time.Time,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			r.logger.Warn("failed to rollback transaction", zap.Error(err))
		}
	}()
	//defer tx.Rollback()

	//history
	_, err = tx.ExecContext(ctx, `
		INSERT INTO currency_history (code, rate, rate_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (code, rate_date)
        DO UPDATE SET rate = EXCLUDED.rate
	`, code.String(), rate.Float64(), date)
	if err != nil {
		return err
	}

	// latest
	_, err = tx.ExecContext(ctx, `
		INSERT INTO currency_latest (code, rate, rate_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (code) 
		DO UPDATE
		SET rate = EXCLUDED.rate,
		    rate_date = EXCLUDED.rate_date
	`, code.String(), rate.Float64(), date)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CurrencyRepoPostgres) GetLatest(
	ctx context.Context,
	code domain.CurrencyCode,
) (domain.Currency, error) {

	var c domain.Currency
	var rateStr string

	err := r.db.QueryRowContext(ctx, `
		SELECT code, rate, rate_date
		FROM currency_latest
		WHERE code = $1
	`, code.String()).Scan(&c.Code, &rateStr, &c.RateDate)

	if err == sql.ErrNoRows {
		return domain.Currency{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Currency{}, err
	}

	rate, err := domain.RateFromString(rateStr)
	if err != nil {
		return domain.Currency{}, err
	}

	c.Rate = rate
	return c, nil
}

func (r *CurrencyRepoPostgres) GetAtDate(
	ctx context.Context,
	code domain.CurrencyCode,
	date time.Time,
) (domain.Currency, error) {

	var c domain.Currency
	var rateStr string

	err := r.db.QueryRowContext(ctx, `
		SELECT code, rate, rate_date
		FROM currency_history
		WHERE code = $1 AND rate_date = $2
	`, code.String(), date).Scan(&c.Code, &rateStr, &c.RateDate)

	if err == sql.ErrNoRows {
		return domain.Currency{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Currency{}, err
	}

	rate, err := domain.RateFromString(rateStr)
	if err != nil {
		return domain.Currency{}, err
	}

	c.Rate = rate
	return c, nil
}

func (r *CurrencyRepoPostgres) GetRange(
	ctx context.Context,
	code domain.CurrencyCode,
	from, to time.Time,
) ([]domain.Currency, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT code, rate, rate_date
		FROM currency_history
		WHERE code = $1 AND rate_date BETWEEN $2 AND $3
		ORDER BY rate_date
	`, code.String(), from, to)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Debug("failed to close rows", zap.Error(err))
		}
	}()

	var res []domain.Currency

	for rows.Next() {
		var c domain.Currency
		var rateStr string

		if err := rows.Scan(&c.Code, &rateStr, &c.RateDate); err != nil {
			return nil, err
		}

		rate, err := domain.RateFromString(rateStr)
		if err != nil {
			continue
		}

		c.Rate = rate
		res = append(res, c)
	}

	if len(res) == 0 {
		return nil, domain.ErrNotFound
	}

	return res, nil
}
