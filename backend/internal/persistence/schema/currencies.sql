-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    symbol VARCHAR(255) NOT NULL
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_code ON currencies (code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_name ON currencies (name);
