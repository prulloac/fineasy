-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(31) NOT NULL,
    symbol VARCHAR(15) NOT NULL
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER NOT NULL references currencies(id) ON DELETE CASCADE ON UPDATE CASCADE,
    base_currency_id INTEGER NOT NULL references currencies(id) ON DELETE CASCADE ON UPDATE CASCADE,
    rate DECIMAL(10, 4) NOT NULL,
    date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS currency_conversion_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL,
    params JSONB NOT NULL,
    run_at VARCHAR(63) NOT NULL
);

CREATE TABLE IF NOT EXISTS currency_conversion_providers_currencies (
    currency_conversion_provider_id INTEGER NOT NULL references currency_conversion_providers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    currency_id INTEGER NOT NULL references currencies(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_code ON currencies (code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_name ON currencies (name);

CREATE UNIQUE INDEX IF NOT EXISTS idx_exchange_rates_currency_id_date ON exchange_rates (currency_id, date);
