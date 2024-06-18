-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    transaction_type_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    date DATE NOT NULL,
    executed_by INTEGER NOT NULL,
    description TEXT NOT NULL,
    receipt_url TEXT,
    registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    registered_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- indexes
CREATE INDEX IF NOT EXISTS transactions_category_id_idx ON transactions (category_id);
CREATE INDEX IF NOT EXISTS transactions_currency_id_idx ON transactions (currency_id);
CREATE INDEX IF NOT EXISTS transactions_transaction_type_id_idx ON transactions (transaction_type_id);
CREATE INDEX IF NOT EXISTS transactions_account_id_idx ON transactions (account_id);
CREATE INDEX IF NOT EXISTS transactions_executed_by_idx ON transactions (executed_by);
CREATE INDEX IF NOT EXISTS transactions_registered_by_idx ON transactions (registered_by);

-- foreign key constraints
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_transactions_category_id'
    ) THEN
        ALTER TABLE transactions
            ADD CONSTRAINT fk_transactions_category_id
            FOREIGN KEY (category_id)
            REFERENCES categories (id)
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_transactions_currency_id'
    ) THEN
        ALTER TABLE transactions
            ADD CONSTRAINT fk_transactions_currency_id
            FOREIGN KEY (currency_id)
            REFERENCES currencies (id)
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_transactions_account_id'
    ) THEN
        ALTER TABLE transactions
            ADD CONSTRAINT fk_transactions_account_id
            FOREIGN KEY (account_id)
            REFERENCES accounts (id)
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_transactions_executed_by'
    ) THEN
        ALTER TABLE transactions
            ADD CONSTRAINT fk_transactions_executed_by
            FOREIGN KEY (executed_by)
            REFERENCES users (id)
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_transactions_registered_by'
    ) THEN
        ALTER TABLE transactions
            ADD CONSTRAINT fk_transactions_registered_by
            FOREIGN KEY (registered_by)
            REFERENCES users (id)
            ON DELETE CASCADE;
    END IF;
END $$;
