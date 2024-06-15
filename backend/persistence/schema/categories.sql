-- postgresql
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(255),
    color VARCHAR(255),
    description TEXT,
    ord INTEGER NOT NULL DEFAULT 0
    group_id INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories (name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_order ON categories (ord);
CREATE INDEX IF NOT EXISTS idx_categories_group_id ON categories (group_id);
