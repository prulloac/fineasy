package transactions

import (
	"time"

	p "github.com/prulloac/fineasy/internal/persistence"
)

type Repository struct {
	Persistence *p.Persistence
}

func NewRepository(persistence *p.Persistence) *Repository {
	return &Repository{persistence}
}

func (r *Repository) Close() {
	r.Persistence.Close()
}

func (r *Repository) CreateTables() error {
	_, err := r.Persistence.Exec(`
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

		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			category VARCHAR(255),
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

		CREATE UNIQUE INDEX IF NOT EXISTS idx_budgets_name_account_id ON budgets (name, account_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_name_group_id ON accounts (name, group_id);
	`)
	return err
}

func (r *Repository) DropTables() error {
	_, err := r.Persistence.Exec(`
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS budgets;
		DROP TABLE IF EXISTS accounts;
	`)
	return err
}

func (r *Repository) CreateAccount(name string, currency string, groupID, createdBy uint) (*Account, error) {
	account := &Account{}
	err := r.Persistence.QueryRow(`
		INSERT INTO accounts (name, group_id, currency, balance, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, name, group_id, currency, balance, created_by, created_at, updated_at
	`, name, groupID, currency, 0.0, createdBy).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *Repository) GetAccountsByUserID(uid uint) ([]Account, error) {
	rows, err := r.Persistence.Query(`
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

func (r *Repository) GetAccountByID(id uint) (*Account, error) {
	account := &Account{}
	err := r.Persistence.QueryRow(`
		SELECT id, name, group_id, currency, balance, created_by, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *Repository) UserHasAccessToAccount(uid, aid uint) (bool, error) {
	var count int
	err := r.Persistence.QueryRow(`
		SELECT COUNT(*)
		FROM accounts a
		JOIN user_groups ug ON a.group_id = ug.group_id
		WHERE a.id = $1 AND ug.user_id = $2
	`, aid, uid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) UpdateAccount(id uint, name, currency string, balance float64) (*Account, error) {
	account := &Account{}
	err := r.Persistence.QueryRow(`
		UPDATE accounts
		SET name = $1, currency = $2, balance = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, group_id, currency, balance, created_by, created_at, updated_at
	`, name, currency, balance, id).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *Repository) CreateBudget(name, currency string, amount float64, startDate, endDate time.Time, accountID, createdBy uint) (*Budget, error) {
	budget := &Budget{}
	err := r.Persistence.QueryRow(`
		INSERT INTO budgets (name, account_id, currency, amount, created_by, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, name, account_id, currency, amount, created_by, start_date, end_date, created_at, updated_at
	`, name, accountID, currency, amount, createdBy, startDate, endDate).Scan(&budget.ID, &budget.Name, &budget.AccountID, &budget.Currency, &budget.Amount, &budget.CreatedBy, &budget.StartDate, &budget.EndDate, &budget.CreatedAt, &budget.UpdatedAt)
	return budget, err
}

func (r *Repository) GetBudgetsByUserID(uid uint) ([]Budget, error) {
	rows, err := r.Persistence.Query(`
		SELECT id, name, account_id, currency, amount, created_by, start_date, end_date, created_at, updated_at
		FROM budgets
		WHERE created_by = $1
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var budget Budget
		if err := rows.Scan(&budget.ID, &budget.Name, &budget.AccountID, &budget.Currency, &budget.Amount, &budget.CreatedBy, &budget.StartDate, &budget.EndDate, &budget.CreatedAt, &budget.UpdatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}
