package transactions

import (
	"log"

	p "github.com/prulloac/fineasy/internal/persistence"
)

type TransactionsRepository struct {
	Persistence *p.Persistence
}

func NewTransactionsRepository(persistence *p.Persistence) *TransactionsRepository {
	return &TransactionsRepository{persistence}
}

func (r *TransactionsRepository) Close() {
	r.Persistence.Close()
}

func (r *TransactionsRepository) CreateTables() error {
	_, err := r.Persistence.SQL().Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_by INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			user_count INT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS user_groups (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			group_id INT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			created_by INT NOT NULL,
			group_id INT NOT NULL,
			currency VARCHAR(255) NOT NULL,
			balance FLOAT NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS budgets (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			account_id INT NOT NULL,
			currency VARCHAR(255) NOT NULL,
			amount FLOAT NOT NULL,
			created_by INT NOT NULL,
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			icon VARCHAR(255) NOT NULL,
			color VARCHAR(255) NOT NULL,
			description TEXT,
			ord INT NOT NULL,
			group_id INT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			category_id INT NOT NULL,
			currency VARCHAR(255) NOT NULL,
			currency_rate FLOAT NOT NULL,
			transaction_type smallint NOT NULL,
			budget_id INT NOT NULL,
			amount FLOAT NOT NULL,
			date TIMESTAMP NOT NULL,
			executed_by_name VARCHAR(255) NOT NULL,
			executed_by_id INT NOT NULL,
			description TEXT,
			receipt_url TEXT,
			registered_by INT NOT NULL,
			registered_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories (name);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_budgets_name_account_id ON budgets (name, account_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_name_group_id ON accounts (name, group_id);
	`)
	return err
}

func (r *TransactionsRepository) DropTable() error {
	_, err := r.Persistence.SQL().Exec(`
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS categories;
		DROP TABLE IF EXISTS budgets;
		DROP TABLE IF EXISTS accounts;
		DROP TABLE IF EXISTS user_groups;
		DROP TABLE IF EXISTS groups;
	`)
	return err
}

func (r *TransactionsRepository) CreateAccount(name string, currency string, groupID, createdBy uint) (*Account, error) {
	account := &Account{}
	err := r.Persistence.SQL().QueryRow(`
		INSERT INTO accounts (name, group_id, currency, balance, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, name, group_id, currency, balance, created_by, created_at, updated_at
	`, name, groupID, currency, 0.0, createdBy).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *TransactionsRepository) GetAccountsByUserID(uid uint) ([]Account, error) {
	rows, err := r.Persistence.SQL().Query(`
		SELECT id, name, group_id, currency, balance, created_by, created_at, updated_at
		FROM accounts
		WHERE created_by = $1
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *TransactionsRepository) GetAccountByID(id int) (*Account, error) {
	account := &Account{}
	err := r.Persistence.SQL().QueryRow(`
		SELECT id, name, group_id, currency, balance, created_by, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *TransactionsRepository) UserHasAccessToAccount(uid, aid int) (bool, error) {
	var count int
	err := r.Persistence.SQL().QueryRow(`
		SELECT COUNT(*)
		FROM accounts a
		JOIN user_groups ug ON a.group_id = ug.group_id
		WHERE a.id = $1 AND ug.user_id = $2
	`, aid, uid).Scan(&count)
	if err != nil {
		log.Printf("⚠️ Error checking user access to account: %s", err)
		return false, err
	}
	log.Printf("🔒 User %d has access to account %d: %t", uid, aid, count > 0)
	return count > 0, nil
}

func (r *TransactionsRepository) UpdateAccount(id int, name, currency string, balance float64) (*Account, error) {
	account := &Account{}
	err := r.Persistence.SQL().QueryRow(`
		UPDATE accounts
		SET name = $1, currency = $2, balance = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, group_id, currency, balance, created_by, created_at, updated_at
	`, name, currency, balance, id).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}
