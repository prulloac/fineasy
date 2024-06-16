-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS transaction_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_transaction_types_name ON transaction_types (name);
