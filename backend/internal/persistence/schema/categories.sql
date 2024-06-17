-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(255),
    color VARCHAR(255),
    description TEXT,
    ord INTEGER NOT NULL DEFAULT 0,
    group_id INTEGER NOT NULL
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories (name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_order ON categories (ord);
CREATE INDEX IF NOT EXISTS idx_categories_group_id ON categories (group_id);

-- foreign key constraints
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_categories_group_id'
    ) THEN
        ALTER TABLE categories
            ADD CONSTRAINT fk_categories_group_id
            FOREIGN KEY (group_id)
            REFERENCES groups (id)
            ON DELETE CASCADE;
    END IF;
END $$;
