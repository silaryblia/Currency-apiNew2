DROP TABLE IF EXISTS currencies;

CREATE TABLE currencies (
                            code TEXT PRIMARY KEY,
                            rate NUMERIC(10,4) NOT NULL,
                            rate_date DATE NOT NULL DEFAULT CURRENT_DATE
);

-- Начальные данные
INSERT INTO currencies (code, rate, rate_date) VALUES
                                                   ('USD', 80.00, CURRENT_DATE),
                                                   ('EUR', 85.00, CURRENT_DATE),
                                                   ('AED', 20.00, CURRENT_DATE)
    ON CONFLICT (code) DO UPDATE SET
    rate = EXCLUDED.rate,
                              rate_date = EXCLUDED.rate_date;