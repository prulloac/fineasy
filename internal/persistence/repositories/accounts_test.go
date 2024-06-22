package repositories

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prulloac/fineasy/internal/persistence/entity"
)

func TestInsertAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	account := entity.Account{Name: "Bank",
		CurrencyID: 1,
		GroupID:    1,
		CreatedBy:  1,
		Balance:    1.0,
	}
	mock.ExpectQuery("SELECT id FROM accounts").
		WithArgs(account.Name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(1, 1, 1, 1.0, account.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AccountsRepository{db}
	err = p.Insert(account)

	if err != nil {
		t.Errorf("error was not expected while inserting account: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllAccounts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	account := entity.Account{Name: "Bank",
		CurrencyID: 1,
		GroupID:    1,
		CreatedBy:  1,
		Balance:    1.0,
	}
	mock.ExpectQuery("SELECT id, created_by, group_id, currency_id, balance, name, created_at, updated_at, disabled FROM accounts").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_by", "group_id", "currency_id", "balance", "name", "created_at", "updated_at", "disabled"}).
			AddRow(1, account.CreatedBy, account.GroupID, account.CurrencyID, account.Balance, account.Name, account.CreatedAt, account.UpdatedAt, account.Disabled))

	var p = AccountsRepository{db}
	r, err := p.GetAll()

	if err != nil {
		t.Errorf("error was not expected while getting accounts: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, a := range r {
		if a.Name != account.Name {
			t.Errorf("expected: %s, got: %s", account.Name, a.Name)
		}
	}
}

func TestGetAccountByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	account := entity.Account{Name: "Bank",
		CurrencyID: 1,
		GroupID:    1,
		CreatedBy:  1,
		Balance:    1.0,
	}
	mock.ExpectQuery("SELECT id, created_by, group_id, currency_id, balance, name, created_at, updated_at, disabled FROM accounts").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_by", "group_id", "currency_id", "balance", "name", "created_at", "updated_at", "disabled"}).
			AddRow(1, account.CreatedBy, account.GroupID, account.CurrencyID, account.Balance, account.Name, account.CreatedAt, account.UpdatedAt, account.Disabled))

	var p = AccountsRepository{db}
	r, err := p.GetByID(1)

	if err != nil {
		t.Errorf("error was not expected while getting account: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.Name != account.Name {
		t.Errorf("expected: %s, got: %s", account.Name, r.Name)
	}
}

func TestGetAccountByGroupID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	account := entity.Account{Name: "Bank",
		CurrencyID: 1,
		GroupID:    1,
		CreatedBy:  1,
		Balance:    1.0,
	}
	mock.ExpectQuery("SELECT id, created_by, group_id, currency_id, balance, name, created_at, updated_at, disabled FROM accounts").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_by", "group_id", "currency_id", "balance", "name", "created_at", "updated_at", "disabled"}).
			AddRow(1, account.CreatedBy, account.GroupID, account.CurrencyID, account.Balance, account.Name, account.CreatedAt, account.UpdatedAt, account.Disabled))

	var p = AccountsRepository{db}
	r, err := p.GetByGroupID(1)

	if err != nil {
		t.Errorf("error was not expected while getting account: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, a := range r {
		if a.Name != account.Name {
			t.Errorf("expected: %s, got: %s", account.Name, a.Name)
		}
	}
}

func TestUpdateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	account := entity.Account{Name: "Bank",
		CurrencyID: 1,
		Balance:    1.0,
		Disabled:   false,
	}
	mock.ExpectExec("UPDATE accounts").
		WithArgs(account.CurrencyID, account.Balance, account.Name, account.Disabled, account.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AccountsRepository{db}
	err = p.Update(account)

	if err != nil {
		t.Errorf("error was not expected while updating account: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
