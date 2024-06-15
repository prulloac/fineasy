package persistence

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

	var p = Persistence{}
	p.db = db
	err = p.InsertCurrency(currency)

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

	var p = Persistence{}
	p.db = db
	r, err := p.GetCurrencies()

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

	var p = Persistence{}
	p.db = db
	r, err := p.GetCurrency(currency.ID)

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

	var p = Persistence{}
	p.db = db
	err = p.UpdateCurrency(currency)

	if err != nil {
		t.Errorf("error was not expected while updating currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertDefaultExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	exchangeRate := ExchangeRate{CurrencyID: 1, GroupID: 1, Rate: 1.0, Date: time.Now()}
	mock.ExpectExec("INSERT INTO default_exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.GroupID, exchangeRate.Rate, exchangeRate.Date).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = Persistence{}
	p.db = db
	err = p.InsertExchangeRate(exchangeRate)

	if err != nil {
		t.Errorf("error was not expected while inserting default exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultExchangeRates(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	exchangeRate := ExchangeRate{CurrencyID: currency.ID, Rate: 1.0, Date: time.Now(), GroupID: 1}
	mock.ExpectQuery("SELECT currency_id, rate, date, group_id FROM default_exchange_rates").
		WithArgs(currency.ID, exchangeRate.GroupID, exchangeRate.Date, exchangeRate.Date).
		WillReturnRows(sqlmock.NewRows([]string{"currency_id", "rate", "date", "group_id"}).
			AddRow(exchangeRate.CurrencyID, exchangeRate.Rate, exchangeRate.Date, exchangeRate.GroupID))

	var p = Persistence{}
	p.db = db
	r, err := p.GetExchangeRates(currency, exchangeRate.GroupID, exchangeRate.Date, exchangeRate.Date)

	if err != nil {
		t.Errorf("error was not expected while getting default exchange rates: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, e := range r {
		if e.CurrencyID != exchangeRate.CurrencyID {
			t.Errorf("expected currency ID %d, but got %d", exchangeRate.CurrencyID, e.CurrencyID)
		}
		if e.Rate != exchangeRate.Rate {
			t.Errorf("expected rate %f, but got %f", exchangeRate.Rate, e.Rate)
		}
		if e.Date != exchangeRate.Date {
			t.Errorf("expected date %s, but got %s", exchangeRate.Date, e.Date)
		}
	}
}