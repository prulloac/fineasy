package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/internal/persistence/entity"
)

type BudgetsRepository struct {
	db *sql.DB
}

func NewBudgetsRepository(db *sql.DB) *BudgetsRepository {
	return &BudgetsRepository{db}
}

func (b *BudgetsRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/budgets.sql")

	if data == nil {
		panic("Error reading budgets schema file!")
	}

	_, err := b.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating budgets table!")
		panic(err)
	}
	fmt.Println("Budgets table created!")
}

func (b *BudgetsRepository) DropTable() {
	_, err := b.db.Exec("DROP TABLE IF EXISTS budgets")
	if err != nil {
		fmt.Println("Error dropping budgets table!")
		panic(err)
	}
	fmt.Println("Budgets table dropped!")
}

func (b *BudgetsRepository) Insert(budget entity.Budget) error {
	// check if the budget already exists regardless of the amount
	var id int
	err := b.db.QueryRow(`
	SELECT
		id
	FROM budgets
	WHERE name = $1`, budget.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := b.db.Exec(`
		INSERT INTO budgets
		(account_id, currency_id, amount, created_by, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)`,
			budget.AccountID, budget.CurrencyID, budget.Amount, budget.CreatedBy, budget.StartDate, budget.EndDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BudgetsRepository) GetAll() ([]entity.Budget, error) {
	rows, err := b.db.Query(`
	SELECT
		id, 
		name, 
		account_id, 
		currency_id, 
		amount, 
		created_by, 
		start_date, 
		end_date, 
		created_at, 
		updated_at
	FROM budgets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []entity.Budget{}
	for rows.Next() {
		var budget entity.Budget
		err := rows.Scan(&budget.ID,
			&budget.Name,
			&budget.AccountID,
			&budget.CurrencyID,
			&budget.Amount,
			&budget.CreatedBy,
			&budget.StartDate,
			&budget.EndDate,
			&budget.CreatedAt,
			&budget.UpdateAt)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

func (b *BudgetsRepository) GetByID(id int) (entity.Budget, error) {
	var budget entity.Budget
	err := b.db.QueryRow(`
	SELECT
		id,
		name,
		account_id,
		currency_id,
		amount,
		created_by,
		start_date,
		end_date,
		created_at,
		updated_at
	FROM budgets
	WHERE id = $1`, id).Scan(&budget.ID,
		&budget.Name,
		&budget.AccountID,
		&budget.CurrencyID,
		&budget.Amount,
		&budget.CreatedBy,
		&budget.StartDate,
		&budget.EndDate,
		&budget.CreatedAt,
		&budget.UpdateAt)
	if err != nil {
		return entity.Budget{}, err
	}
	return budget, nil
}

func (b *BudgetsRepository) GetByAccountID(accountID int) ([]entity.Budget, error) {
	rows, err := b.db.Query(`
	SELECT
		id,
		name,
		account_id,
		currency_id,
		amount,
		created_by,
		start_date,
		end_date,
		created_at,
		updated_at
	FROM budgets
	WHERE account_id = $1`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []entity.Budget{}
	for rows.Next() {
		var budget entity.Budget
		err := rows.Scan(&budget.ID,
			&budget.Name,
			&budget.AccountID,
			&budget.CurrencyID,
			&budget.Amount,
			&budget.CreatedBy,
			&budget.StartDate,
			&budget.EndDate,
			&budget.CreatedAt,
			&budget.UpdateAt)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

func (b *BudgetsRepository) Update(budget entity.Budget) error {
	_, err := b.db.Exec(`
	UPDATE budgets
	SET currency_id = $1, 
		amount = $2, 
		start_date = $3, 
		end_date = $4,
		name = $5,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $6`,
		budget.CurrencyID,
		budget.Amount,
		budget.StartDate,
		budget.EndDate,
		budget.Name,
		budget.ID)
	if err != nil {
		return err
	}
	return nil
}
