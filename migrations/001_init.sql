-- Актуальные курсы (быстрые запросы GetAll / GetLatest)
CREATE TABLE currency_latest (
                                 code TEXT PRIMARY KEY,
                                 rate NUMERIC(18,6) NOT NULL,
                                 rate_date DATE NOT NULL
);

-- История курсов (source of truth)
CREATE TABLE currency_history (
                                  code TEXT NOT NULL,
                                  rate NUMERIC(18,6) NOT NULL,
                                  rate_date DATE NOT NULL,
                                  created_at TIMESTAMP NOT NULL DEFAULT now(),
                                  PRIMARY KEY (code, rate_date)
);

-- Для GetLatest из истории (если понадобится)
CREATE INDEX idx_currency_history_code_date_desc
    ON currency_history (code, rate_date DESC);

-- Для GetAtDate / GetRange
CREATE INDEX idx_currency_history_code_date
    ON currency_history (code, rate_date);