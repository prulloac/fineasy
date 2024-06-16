package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/prulloac/fineasy/persistence/entity"
)

func TestInsertExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	exchangeRate := ExchangeRate{CurrencyID: 1, GroupID: 1, Rate: 1.0, Date: time.Now()}
	mock.ExpectQuery("SELECT id FROM exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.GroupID, exchangeRate.Date).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.GroupID, exchangeRate.Rate, exchangeRate.Date).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = ExchangeRateRepository{db}
	err = p.InsertExchangeRate(exchangeRate)

	if err != nil {
		t.Errorf("error was not expected while inserting default exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetExchangeRates(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	exchangeRate := ExchangeRate{ID: 1, CurrencyID: currency.ID, Rate: 1.0, Date: time.Now(), GroupID: 1}
	mock.ExpectQuery("SELECT id, currency_id, rate, date, group_id FROM default_exchange_rates").
		WithArgs(currency.ID, exchangeRate.GroupID, exchangeRate.Date, exchangeRate.Date).
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_id", "rate", "date", "group_id"}).
			AddRow(exchangeRate.ID, exchangeRate.CurrencyID, exchangeRate.Rate, exchangeRate.Date, exchangeRate.GroupID))

	var p = ExchangeRateRepository{db}
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

func TestGetExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	exchangeRate := ExchangeRate{CurrencyID: currency.ID, Rate: 1.0, Date: time.Now(), GroupID: 1}
	mock.ExpectQuery("SELECT currency_id, group_id, rate, date FROM exchange_rates").
		WithArgs(currency.ID, exchangeRate.GroupID, exchangeRate.Date).
		WillReturnRows(sqlmock.NewRows([]string{"currency_id", "group_id", "rate", "date"}).
			AddRow(exchangeRate.CurrencyID, exchangeRate.GroupID, exchangeRate.Rate, exchangeRate.Date))

	var p = ExchangeRateRepository{db}
	r, err := p.GetExchangeRate(currency, exchangeRate.GroupID, exchangeRate.Date)

	if err != nil {
		t.Errorf("error was not expected while getting default exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.CurrencyID != exchangeRate.CurrencyID {
		t.Errorf("expected currency ID %d, but got %d", exchangeRate.CurrencyID, r.CurrencyID)
	}
	if r.GroupID != exchangeRate.GroupID {
		t.Errorf("expected group ID %d, but got %d", exchangeRate.GroupID, r.GroupID)
	}
	if r.Rate != exchangeRate.Rate {
		t.Errorf("expected rate %f, but got %f", exchangeRate.Rate, r.Rate)
	}
	if r.Date != exchangeRate.Date {
		t.Errorf("expected date %s, but got %s", exchangeRate.Date, r.Date)
	}
}

func TestUpdateExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	exchangeRate := ExchangeRate{ID: 1, CurrencyID: 1, GroupID: 1, Rate: 1.0, Date: time.Now()}
	mock.ExpectExec("UPDATE exchange_rates").
		WithArgs(exchangeRate.Rate, exchangeRate.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = ExchangeRateRepository{db}
	err = p.UpdateExchangeRate(exchangeRate)

	if err != nil {
		t.Errorf("error was not expected while updating default exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
