package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/internal/persistence/entity"
)

type AccountsRepository struct {
	db *sql.DB
}

func NewAccountsRepository(db *sql.DB) *AccountsRepository {
	return &AccountsRepository{db}
}

func (a *AccountsRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/accounts.sql")

	if data == nil {
		panic("Error reading accounts schema file!")
	}

	_, err := a.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating accounts table!")
		panic(err)
	}
	fmt.Println("Accounts table created!")
}

func (a *AccountsRepository) DropTable() {
	_, err := a.db.Exec("DROP TABLE IF EXISTS accounts")
	if err != nil {
		fmt.Println("Error dropping accounts table!")
		panic(err)
	}
	fmt.Println("Accounts table dropped!")
}

func (a *AccountsRepository) Insert(account entity.Account) error {
	// check if the account already exists regardless of the balance
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM accounts
	WHERE name = $1`, account.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := a.db.Exec(`
		INSERT INTO accounts
		(created_by, group_id, currency_id, balance, name) VALUES ($1, $2, $3, $4, $5)`,
			account.CreatedBy, account.GroupID, account.CurrencyID, account.Balance, account.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountsRepository) GetAll() ([]entity.Account, error) {
	rows, err := a.db.Query(`
	SELECT
		id, 
		created_by, 
		group_id, 
		currency_id, 
		balance, 
		name, 
		created_at, 
		updated_at,
		disabled
	FROM accounts`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []entity.Account{}
	for rows.Next() {
		var account entity.Account
		err := rows.Scan(&account.ID,
			&account.CreatedBy,
			&account.GroupID,
			&account.CurrencyID,
			&account.Balance,
			&account.Name,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Disabled)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (a *AccountsRepository) GetByID(id int) (entity.Account, error) {
	var account entity.Account
	err := a.db.QueryRow(`
	SELECT
		id, 
		created_by, 
		group_id, 
		currency_id, 
		balance, 
		name, 
		created_at, 
		updated_at,
		disabled
	FROM accounts
	WHERE id = $1`, id).Scan(&account.ID,
		&account.CreatedBy,
		&account.GroupID,
		&account.CurrencyID,
		&account.Balance,
		&account.Name,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Disabled)

	if err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (a *AccountsRepository) GetByGroupID(groupID int) ([]entity.Account, error) {
	rows, err := a.db.Query(`
	SELECT
		id, 
		created_by, 
		group_id, 
		currency_id, 
		balance, 
		name, 
		created_at, 
		updated_at,
		disabled
	FROM accounts
	WHERE group_id = $1`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []entity.Account{}
	for rows.Next() {
		var account entity.Account
		err := rows.Scan(&account.ID,
			&account.CreatedBy,
			&account.GroupID,
			&account.CurrencyID,
			&account.Balance,
			&account.Name,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Disabled)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (a *AccountsRepository) Update(account entity.Account) error {
	_, err := a.db.Exec(`
	UPDATE accounts
	SET
		currency_id = $1,
		balance = $2,
		name = $3,
		updated_at = CURRENT_TIMESTAMP,
		disabled = $4
	WHERE id = $5`,
		account.CurrencyID,
		account.Balance,
		account.Name,
		account.Disabled,
		account.ID)
	if err != nil {
		return err
	}
	return nil
}
