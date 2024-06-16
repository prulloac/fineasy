-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS budgets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    account_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    created_by INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- indexes
CREATE INDEX IF NOT EXISTS budgets_name_idx ON budgets (account_id, name);
CREATE INDEX IF NOT EXISTS budgets_created_by_idx ON budgets (created_by);
CREATE INDEX IF NOT EXISTS budgets_currency_id_idx ON budgets (currency_id);

-- foreign key constraints
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_budgets_created_by'
    ) THEN
        ALTER TABLE budgets
            ADD CONSTRAINT fk_budgets_created_by
            FOREIGN KEY (created_by)
            REFERENCES users (id)
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_budgets_account_id'
    ) THEN
        ALTER TABLE budgets
            ADD CONSTRAINT fk_budgets_account_id
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
        WHERE constraint_name = 'fk_budgets_currency_id'
    ) THEN
        ALTER TABLE budgets
            ADD CONSTRAINT fk_budgets_currency_id
            FOREIGN KEY (currency_id)
            REFERENCES currencies (id)
            ON DELETE CASCADE;
    END IF;
END $$;
