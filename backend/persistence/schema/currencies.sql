-- postgresql
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    symbol VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_code ON currencies (code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_name ON currencies (name);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    rate DECIMAL(10, 4) NOT NULL,
    date DATE NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_exchange_rates_currency_id_date ON exchange_rates (currency_id, date);
CREATE INDEX IF NOT EXISTS idx_exchange_rates_group_id_date ON exchange_rates (group_id, date);
