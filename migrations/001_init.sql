CREATE TABLE currency_latest (
                                 code TEXT PRIMARY KEY,
                                 rate NUMERIC NOT NULL,
                                 rate_date DATE NOT NULL
);

CREATE TABLE currency_history (
                                  code TEXT NOT NULL,
                                  rate NUMERIC NOT NULL,
                                  rate_date DATE NOT NULL,
                                  PRIMARY KEY (code, rate_date)
);

CREATE INDEX idx_currency_history_code_date
    ON currency_history (code, rate_date);