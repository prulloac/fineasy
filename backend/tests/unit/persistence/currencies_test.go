package persistence__test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prulloac/fineasy/persistence"
	. "github.com/prulloac/fineasy/persistence/entity"
)

func TestInsertCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id FROM currencies").
		WithArgs(currency.Code).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO currencies").
		WithArgs(currency.Code, currency.Symbol, currency.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = persistence.NewPersistence(db)
	err = p.GetCurrencyRepository().InsertCurrency(currency)

	if err != nil {
		t.Errorf("error was not expected while inserting currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCurrencies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id, code, symbol, name FROM currencies").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "symbol", "name"}).
			AddRow(1, currency.Code, currency.Symbol, currency.Name))

	var p = persistence.NewPersistence(db)
	r, err := p.GetCurrencyRepository().GetCurrencies()

	if err != nil {
		t.Errorf("error was not expected while getting currencies: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, c := range r {
		if c.Code != currency.Code {
			t.Errorf("expected code %s, but got %s", currency.Code, c.Code)
		}
		if c.Symbol != currency.Symbol {
			t.Errorf("expected symbol %s, but got %s", currency.Symbol, c.Symbol)
		}
		if c.Name != currency.Name {
			t.Errorf("expected name %s, but got %s", currency.Name, c.Name)
		}
	}
}

func TestGetCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id, code, symbol, name FROM currencies").
		WithArgs(currency.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "symbol", "name"}).
			AddRow(currency.ID, currency.Code, currency.Symbol, currency.Name))

	var p = persistence.NewPersistence(db)
	r, err := p.GetCurrencyRepository().GetCurrency(currency.ID)

	if err != nil {
		t.Errorf("error was not expected while getting currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.ID != currency.ID {
		t.Errorf("expected ID %d, but got %d", currency.ID, r.ID)
	}
	if r.Code != currency.Code {
		t.Errorf("expected code %s, but got %s", currency.Code, r.Code)
	}
	if r.Symbol != currency.Symbol {
		t.Errorf("expected symbol %s, but got %s", currency.Symbol, r.Symbol)
	}
	if r.Name != currency.Name {
		t.Errorf("expected name %s, but got %s", currency.Name, r.Name)
	}
}

func TestUpdateCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectExec("UPDATE currencies").
		WithArgs(currency.Code, currency.Symbol, currency.Name, currency.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = persistence.NewPersistence(db)
	err = p.GetCurrencyRepository().UpdateCurrency(currency)

	if err != nil {
		t.Errorf("error was not expected while updating currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
