-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    rate DECIMAL(10, 4) NOT NULL,
    date DATE NOT NULL
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_exchange_rates_currency_id_date ON exchange_rates (currency_id, date);
CREATE INDEX IF NOT EXISTS idx_exchange_rates_group_id_date ON exchange_rates (group_id, date);

-- foreign key constraints
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_exchange_rates_currency_id'
    ) THEN
        ALTER TABLE exchange_rates
            ADD CONSTRAINT fk_exchange_rates_currency_id
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
        WHERE constraint_name = 'fk_exchange_rates_group_id'
    ) THEN
        ALTER TABLE exchange_rates
            ADD CONSTRAINT fk_exchange_rates_group_id
            FOREIGN KEY (group_id)
            REFERENCES groups (id)
            ON DELETE CASCADE;
    END IF;
END $$;
