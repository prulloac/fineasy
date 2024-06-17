-- persistence

-- tables
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    created_by INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    balance DECIMAL(10, 2) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    disabled BOOLEAN NOT NULL DEFAULT FALSE
);

-- indexes
CREATE INDEX IF NOT EXISTS accounts_name_idx ON accounts (group_id, name);
CREATE INDEX IF NOT EXISTS accounts_created_by_idx ON accounts (created_by);
CREATE INDEX IF NOT EXISTS accounts_currency_id_idx ON accounts (currency_id);

-- foreign key constraints
ALTER TABLE accounts
    ADD CONSTRAINT fk_accounts_created_by
    FOREIGN KEY (created_by)
    REFERENCES users (id)
    ON DELETE CASCADE;

ALTER TABLE accounts
    ADD CONSTRAINT fk_accounts_group_id
    FOREIGN KEY (group_id)
    REFERENCES groups (id)
    ON DELETE CASCADE;

ALTER TABLE accounts
    ADD CONSTRAINT fk_accounts_currency_id
    FOREIGN KEY (currency_id)
    REFERENCES currencies (id)
    ON DELETE CASCADE;
