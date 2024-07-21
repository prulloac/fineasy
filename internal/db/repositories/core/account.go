package repositories

import (
	"encoding/json"
	"fmt"
	"time"
)

type Account struct {
	ID        uint
	CreatedBy uint
	GroupID   uint
	Currency  string
	Balance   float64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Account) String() string {
	out, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%+v", a.Name)
	}
	return string(out)
}

func (r *CoreRepository) CreateAccount(name string, currency string, groupID, createdBy uint) (*Account, error) {
	account := &Account{}
	err := r.Persistence.QueryRow(`
		INSERT INTO accounts (name, group_id, currency, balance, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, name, group_id, currency, balance, created_by, created_at, updated_at
	`, name, groupID, currency, 0.0, createdBy).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *CoreRepository) GetAccountsByUserID(uid uint) ([]Account, error) {
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

func (r *CoreRepository) GetAccountByID(id uint) (*Account, error) {
	account := &Account{}
	err := r.Persistence.QueryRow(`
		SELECT id, name, group_id, currency, balance, created_by, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}

func (r *CoreRepository) UserHasAccessToAccount(uid, aid uint) (bool, error) {
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

func (r *CoreRepository) UpdateAccount(id uint, name, currency string, balance float64) (*Account, error) {
	account := &Account{}
	err := r.Persistence.QueryRow(`
		UPDATE accounts
		SET name = $1, currency = $2, balance = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, group_id, currency, balance, created_by, created_at, updated_at
	`, name, currency, balance, id).Scan(&account.ID, &account.Name, &account.GroupID, &account.Currency, &account.Balance, &account.CreatedBy, &account.CreatedAt, &account.UpdatedAt)
	return account, err
}
