-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS user_groups (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_groups_user_id_group_id ON user_groups (user_id, group_id);

-- foreign key constraints
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_user_groups_user_id'
    ) THEN
        ALTER TABLE user_groups
            ADD CONSTRAINT fk_user_groups_user_id
            FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_user_groups_group_id'
    ) THEN
        ALTER TABLE user_groups
            ADD CONSTRAINT fk_user_groups_group_id
            FOREIGN KEY (group_id)
            REFERENCES groups (id)
            ON DELETE CASCADE;
    END IF;
END $$;
