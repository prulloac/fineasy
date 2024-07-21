package repositories

import (
	"encoding/json"
	"fmt"
	"time"
)

type Budget struct {
	ID        uint
	Name      string
	AccountID uint
	Currency  string
	Amount    float64
	CreatedBy uint
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Budget) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("%+v", b.Name)
	}
	return string(out)
}

func (r *CoreRepository) CreateBudget(name, currency string, amount float64, startDate, endDate time.Time, accountID, createdBy uint) (*Budget, error) {
	budget := &Budget{}
	err := r.Persistence.QueryRow(`
		INSERT INTO budgets (name, account_id, currency, amount, created_by, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, name, account_id, currency, amount, created_by, start_date, end_date, created_at, updated_at
	`, name, accountID, currency, amount, createdBy, startDate, endDate).Scan(&budget.ID, &budget.Name, &budget.AccountID, &budget.Currency, &budget.Amount, &budget.CreatedBy, &budget.StartDate, &budget.EndDate, &budget.CreatedAt, &budget.UpdatedAt)
	return budget, err
}

func (r *CoreRepository) GetBudgetsByUserID(uid uint) ([]Budget, error) {
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
