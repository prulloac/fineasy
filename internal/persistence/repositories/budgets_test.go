package repositories

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prulloac/fineasy/internal/persistence/entity"
)

func TestInsertBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	budget := entity.Budget{
		AccountID:  1,
		CurrencyID: 1,
		CreatedBy:  1,
		Amount:     1.0,
		Name:       "Budget",
	}

	mock.ExpectQuery("SELECT id FROM budgets").
		WithArgs(budget.Name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO budgets").
		WithArgs(1, 1, 1.0, 1, budget.StartDate, budget.EndDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = BudgetsRepository{db}
	err = p.Insert(budget)

	if err != nil {
		t.Errorf("error was not expected while inserting budget: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllBudgets(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	budget := entity.Budget{
		AccountID:  1,
		CurrencyID: 1,
		CreatedBy:  1,
		Amount:     1.0,
		Name:       "Budget",
	}

	mock.ExpectQuery("SELECT id, name, account_id, currency_id, amount, created_by, start_date, end_date, created_at, updated_at FROM budgets").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_id", "currency_id", "amount", "created_by", "start_date", "end_date", "created_at", "updated_at"}).
			AddRow(1, budget.Name, budget.AccountID, budget.CurrencyID, budget.Amount, budget.CreatedBy, budget.StartDate, budget.EndDate, budget.CreatedAt, budget.UpdateAt))

	var p = BudgetsRepository{db}
	r, err := p.GetAll()

	if err != nil {
		t.Errorf("error was not expected while getting all budgets: %s", err)
	}

	if r[0].AccountID != budget.AccountID {
		t.Errorf("error was not expected while getting all budgets: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBudgetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	budget := entity.Budget{
		AccountID:  1,
		CurrencyID: 1,
		CreatedBy:  1,
		Amount:     1.0,
		Name:       "Budget",
	}

	mock.ExpectQuery("SELECT id, name, account_id, currency_id, amount, created_by, start_date, end_date, created_at, updated_at FROM budgets WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_id", "currency_id", "amount", "created_by", "start_date", "end_date", "created_at", "updated_at"}).
			AddRow(1, budget.Name, budget.AccountID, budget.CurrencyID, budget.Amount, budget.CreatedBy, budget.StartDate, budget.EndDate, budget.CreatedAt, budget.UpdateAt))

	var p = BudgetsRepository{db}
	r, err := p.GetByID(1)

	if err != nil {
		t.Errorf("error was not expected while getting budget by id: %s", err)
	}

	if r.AccountID != budget.AccountID {
		t.Errorf("error was not expected while getting budget by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByAccountID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	budget := entity.Budget{
		AccountID:  1,
		CurrencyID: 1,
		CreatedBy:  1,
		Amount:     1.0,
		Name:       "Budget",
	}

	mock.ExpectQuery("SELECT id, name, account_id, currency_id, amount, created_by, start_date, end_date, created_at, updated_at FROM budgets WHERE account_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_id", "currency_id", "amount", "created_by", "start_date", "end_date", "created_at", "updated_at"}).
			AddRow(1, budget.Name, budget.AccountID, budget.CurrencyID, budget.Amount, budget.CreatedBy, budget.StartDate, budget.EndDate, budget.CreatedAt, budget.UpdateAt))

	var p = BudgetsRepository{db}
	r, err := p.GetByAccountID(1)

	if err != nil {
		t.Errorf("error was not expected while getting budget by account id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, v := range r {
		if v.AccountID != budget.AccountID {
			t.Errorf("error was not expected while getting budget by account id: %s", err)
		}
	}

}

func TestUpdateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	budget := entity.Budget{
		ID:         1,
		AccountID:  1,
		CurrencyID: 1,
		CreatedBy:  1,
		Amount:     1.0,
		Name:       "Budget",
	}

	mock.ExpectExec("UPDATE budgets SET").
		WithArgs(budget.CurrencyID,
			budget.Amount,
			budget.StartDate,
			budget.EndDate,
			budget.Name,
			1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = BudgetsRepository{db}
	err = p.Update(budget)

	if err != nil {
		t.Errorf("error was not expected while updating budget: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
