-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS groups_name_idx ON groups (name);
CREATE INDEX IF NOT EXISTS groups_created_by_idx ON groups (created_by);

-- foreign key constraints
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'groups_created_by_fkey'
    ) THEN
        ALTER TABLE groups
            ADD CONSTRAINT groups_created_by_fkey
            FOREIGN KEY (created_by)
            REFERENCES users (id)
            ON DELETE CASCADE;
    END IF;
END $$;
