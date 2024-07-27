package repositories

import "github.com/prulloac/fineasy/internal/db/persistence"

type CoreRepository struct {
	Persistence *persistence.Persistence
}

func NewCoreRepository(persistence *persistence.Persistence) *CoreRepository {
	instance := &CoreRepository{}
	instance.Persistence = persistence
	err := instance.CreateTables()
	if err != nil {
		panic(err)
	}
	return instance
}

func (s *CoreRepository) Close() {
	s.Persistence.Close()
}

func (s *CoreRepository) CreateTables() error {
	_, err := s.Persistence.Exec(`
		CREATE TABLE IF NOT EXISTS user_data (
			user_id INT PRIMARY KEY,
			avatar_url TEXT NOT NULL,
			display_name VARCHAR(255) NOT NULL,
			currency VARCHAR(6) NOT NULL,
			language VARCHAR(6) NOT NULL,
			timezone VARCHAR(255) NOT NULL,
			upserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_preferences (
			user_id INT NOT NULL REFERENCES user_data(user_id),
			key VARCHAR(255) NOT NULL,
			value TEXT NOT NULL,
			upserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_preferences_user_id_key ON user_preferences (user_id, key);

		CREATE TABLE IF NOT EXISTS friendships (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES user_data(user_id),
			friend_id INT NOT NULL REFERENCES user_data(user_id),
			status INTEGER NOT NULL,
			relation_type INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_friendships_user_id ON friendships (user_id);
		CREATE INDEX IF NOT EXISTS idx_friendships_friend_id ON friendships (friend_id);

		CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_by INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_groups_created_by ON groups (created_by);

		CREATE TABLE IF NOT EXISTS user_groups (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL REFERENCES groups(id),
			user_id INT NOT NULL REFERENCES user_data(user_id),
			joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			left_at TIMESTAMP,
			status INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_user_groups_user_id ON user_groups (user_id);
		CREATE INDEX IF NOT EXISTS idx_user_groups_group_id ON user_groups (group_id);

		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			created_by INT NOT NULL REFERENCES user_data(user_id),
			group_id INT NOT NULL REFERENCES groups(id),
			currency VARCHAR(255) NOT NULL,
			balance FLOAT NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_accounts_group_id ON accounts (group_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_name_group_id ON accounts (name, group_id);

		CREATE TABLE IF NOT EXISTS budgets (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			account_id INT NOT NULL REFERENCES accounts(id),
			currency VARCHAR(255) NOT NULL,
			amount FLOAT NOT NULL,
			created_by INT NOT NULL REFERENCES user_data(user_id),
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_budgets_account_id ON budgets (account_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_budgets_name_account_id ON budgets (name, account_id);

		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			Icon VARCHAR(255) NOT NULL,
			Color VARCHAR(255) NOT NULL,
			Description TEXT,
			ord INTEGER NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories (name);

		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			category_id INT NOT NULL REFERENCES categories(id),
			currency VARCHAR(255) NOT NULL,
			currency_rate FLOAT NOT NULL,
			transaction_type smallint NOT NULL,
			budget_id INT NOT NULL REFERENCES budgets(id),
			amount FLOAT NOT NULL,
			date TIMESTAMP NOT NULL,
			executed_by_name VARCHAR(255) NOT NULL,
			executed_by_id INT REFERENCES user_data(user_id),
			description TEXT,
			receipt_url TEXT,
			registered_by INT NOT NULL REFERENCES user_data(user_id),
			registered_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions (category_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_budget_id ON transactions (budget_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_executed_by_id ON transactions (executed_by_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_registered_by ON transactions (registered_by);
		CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (date);
`)
	return err
}

func (s *CoreRepository) DropTables() error {
	_, err := s.Persistence.Exec(`
	DROP TABLE IF EXISTS transactions;
	DROP TABLE IF EXISTS categories;
	DROP TABLE IF EXISTS budgets;
	DROP TABLE IF EXISTS accounts;
	DROP TABLE IF EXISTS user_groups;
	DROP TABLE IF EXISTS groups;
	DROP TABLE IF EXISTS friendships;
	DROP TABLE IF EXISTS user_preferences;
	DROP TABLE IF EXISTS user_data;
	`)
	return err
}
